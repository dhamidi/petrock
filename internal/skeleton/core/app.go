package core

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// App is the central struct that holds all application dependencies and state
type App struct {
	DB              *sql.DB
	MessageLog      *MessageLog
	CommandRegistry *CommandRegistry
	QueryRegistry   *QueryRegistry
	Executor        *Executor
	Features        []string       // Track registered feature names
	Routes          []string       // Track registered routes
	Mux             *http.ServeMux // Store the HTTP mux
	AppState        interface{}    // Generic application state interface

	// Worker management
	workers      []Worker           // Registered background workers
	workerCtx    context.Context    // Context for worker goroutines
	workerCancel context.CancelFunc // Function to cancel worker context
	workerWg     sync.WaitGroup     // WaitGroup for worker goroutines
}

// NewApp creates and initializes all core dependencies
func NewApp(dbPath string) (*App, error) {
	slog.Info("Initializing application...")

	// 1. Initialize Core Registries
	commandRegistry := NewCommandRegistry()
	queryRegistry := NewQueryRegistry()
	slog.Debug("Initialized command and query registries")

	// 2. Initialize Encoder
	encoder := &JSONEncoder{} // Using JSON encoder
	slog.Debug("Initialized JSON encoder")

	// 3. Initialize Database Connection
	slog.Debug("Setting up database connection", "path", dbPath)
	db, err := SetupDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database at %s: %w", dbPath, err)
	}

	// 4. Initialize Message Log
	slog.Debug("Initializing message log")
	messageLog, err := NewMessageLog(db, encoder)
	if err != nil {
		// Close the database connection since we're returning an error
		db.Close()
		return nil, fmt.Errorf("failed to initialize message log: %w", err)
	}

	// 5. Initialize Central Command Executor
	slog.Debug("Initializing central command executor")
	executor := NewExecutor(messageLog, commandRegistry)

	// 6. Return the App struct with all dependencies
	return &App{
		DB:              db,
		MessageLog:      messageLog,
		CommandRegistry: commandRegistry,
		QueryRegistry:   queryRegistry,
		Executor:        executor,
		Features:        []string{},
		Routes:          []string{},
		// AppState will be initialized by the caller
	}, nil
}

// RegisterFeatures registers all application features
// This MUST be called before replay
func (a *App) RegisterFeatures(appState interface{}) {
	slog.Debug("Registering features...")
	// The caller should provide this function to register all features
	// RegisterAllFeatures(a)
	a.AppState = appState
	slog.Info("Features registered")
}

// RegisterFeature registers a feature with the application
// This method tracks the feature name and calls the provided registration function
func (a *App) RegisterFeature(name string, registerFn func(app *App, state interface{})) {
	slog.Debug("Registering feature", "name", name)
	a.Features = append(a.Features, name)
	registerFn(a, a.AppState)
}

// RegisterRoute registers an HTTP route with the application
// This is a wrapper around mux.HandleFunc that tracks the route
func (a *App) RegisterRoute(pattern string, handler http.HandlerFunc) {
	slog.Debug("Registering route", "pattern", pattern)
	a.Routes = append(a.Routes, pattern)
	if a.Mux != nil {
		a.Mux.HandleFunc(pattern, handler)
	}
}

// RegisterWorker registers a background worker with the application
// Workers are started when StartWorkers is called and stopped during shutdown
func (a *App) RegisterWorker(worker Worker) {
	slog.Debug("Registering worker", "type", fmt.Sprintf("%T", worker))
	
	// If it's a CommandWorker, set up its dependencies
	if cmdWorker, ok := worker.(*CommandWorker); ok {
		cmdWorker.SetDependencies(a.MessageLog, a.Executor)
	}
	
	a.workers = append(a.workers, worker)
}

// StartWorkers initializes and starts all registered workers
// Each worker is started in its own goroutine where it first replays all messages
// and then periodically calls Work()
func (a *App) StartWorkers(ctx context.Context) error {
	slog.Info("Starting workers...")

	// Create a cancelable context for worker operations
	a.workerCtx, a.workerCancel = context.WithCancel(ctx)

	// Start each worker in its own goroutine
	for i, worker := range a.workers {
		a.workerWg.Add(1)

		// Start worker goroutine
		go func(index int, w Worker) {
			defer a.workerWg.Done()

			// Initialize worker
			slog.Debug("Initializing worker", "index", index, "type", fmt.Sprintf("%T", w))
			if err := w.Start(a.workerCtx); err != nil {
				slog.Error("Worker initialization failed", "index", index, "error", err)
				return
			}

			// Replay all messages first (synchronously within this goroutine)
			slog.Debug("Worker replaying messages", "index", index)

			// Worker's lastProcessedID is initialized to 0, so calling Work() will process
			// all messages from the beginning. Messages are processed in the order they appear
			// in the log and state is updated.
			if err := w.Work(); err != nil {
				slog.Error("Worker message replay failed", "index", index, "error", err)
			}

			slog.Info("Worker message replay completed", "index", index)

			// Create ticker with jitter (1-2 seconds)
			baseTick := time.Second
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			ticker := time.NewTicker(baseTick + jitter)
			defer ticker.Stop()

			slog.Debug("Worker started", "index", index, "interval", baseTick+jitter)

			// Periodically call Work()
			for {
				select {
				case <-a.workerCtx.Done():
					// Context was cancelled, exit the goroutine
					slog.Debug("Worker stopping due to context cancellation", "index", index)
					return

				case <-ticker.C:
					// Time to do work
					if err := w.Work(); err != nil {
						// Log error but don't stop worker on work errors
						slog.Error("Worker cycle failed", "index", index, "error", err)
					}
				}
			}
		}(i, worker)
	}

	slog.Info("All workers started", "count", len(a.workers))
	return nil
}

// StopWorkers gracefully shuts down all workers
// This method signals workers to stop and waits for them to finish with timeout
func (a *App) StopWorkers(ctx context.Context) error {
	if a.workerCancel == nil {
		// No workers were started
		return nil
	}

	slog.Info("Stopping workers...")

	// Signal all workers to stop by canceling the worker context
	a.workerCancel()

	// Create a channel to signal when all workers have stopped
	done := make(chan struct{})

	// Wait for workers to finish in a goroutine
	go func() {
		a.workerWg.Wait()
		close(done)
	}()

	// Wait for workers to finish or timeout
	select {
	case <-done:
		// All workers finished successfully
		slog.Info("All workers stopped successfully")
		return nil

	case <-ctx.Done():
		// Timeout or parent context canceled
		slog.Warn("Timed out waiting for workers to stop")
		return &WorkerError{Op: "stop", Err: ctx.Err()}
	}

	// Explicitly call Stop() on each worker with a separate context
	var stopErrors []error
	for i, worker := range a.workers {
		// Create a short timeout for stopping each worker
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := worker.Stop(stopCtx)
		cancel()

		if err != nil {
			slog.Error("Failed to stop worker", "index", i, "type", fmt.Sprintf("%T", worker), "error", err)
			stopErrors = append(stopErrors, fmt.Errorf("worker %d (%T): %w", i, worker, err))
		}
	}

	if len(stopErrors) > 0 {
		// Return a combined error if any workers failed to stop
		return &WorkerError{Op: "stop", Err: fmt.Errorf("%d worker(s) failed to stop properly", len(stopErrors))}
	}

	return nil
}

// ReplayLog replays the message log to build application state
func (a *App) ReplayLog() error {
	slog.Info("Replaying message log to build application state...")

	// Get the starting version
	startVersion := uint64(0)         // Start from the beginning
	replayCtx := context.Background() // Use a background context for replay
	replayErrors := 0                 // Count errors during replay

	// Use the iterator to process messages one by one
	messageCount := 0
	// Use a longer timeout specifically for replay operations
	replayTimeoutCtx, cancel := context.WithTimeout(replayCtx, 60*time.Second)
	defer cancel()
	
	for msg := range a.MessageLog.After(replayTimeoutCtx, startVersion) {
		messageCount++
		// DecodedPayload contains the decoded command or query
		decodedMsg := msg.DecodedPayload

		// Check if the message is a command
		cmd, isCommand := decodedMsg.(Command)
		if !isCommand {
			// If it's not a command, skip it during replay
			slog.Debug("Skipping non-command message during handler replay", "id", msg.ID, "type", fmt.Sprintf("%T", decodedMsg))
			continue
		}

		// Get the state update handler for the command
		handler, found := a.CommandRegistry.GetHandler(cmd.CommandName())
		if !found {
			// This indicates a potential issue: a command was logged but no handler is registered.
			slog.Error("Log replay: No state handler found for logged command", "id", msg.ID, "name", cmd.CommandName())
			replayErrors++
			continue // Skip this command
		}

		// Execute ONLY the state update handler. DO NOT VALIDATE OR LOG AGAIN.
		// Pass the message metadata to provide context like timestamp during replay
		slog.Debug("Log replay: Applying state handler", "id", msg.ID, "name", cmd.CommandName())
		handlerErr := handler(replayCtx, cmd, &msg.Message)
		if handlerErr != nil {
			// PANIC! If a state handler fails during replay, the state logic is
			// inconsistent with the previously validated and logged command.
			slog.Error("Log replay: State update handler failed! PANICKING.", "id", msg.ID, "name", cmd.CommandName(), "error", handlerErr)
			panic(fmt.Sprintf("unrecoverable state inconsistency during log replay: handler for %q failed: %v", cmd.CommandName(), handlerErr))
		}
	}

	slog.Info("State replay completed", "message_count", messageCount, "replay_errors", replayErrors)
	if replayErrors > 0 {
		slog.Warn("Some messages were skipped during state replay due to missing handlers.")
	}

	return nil
}

// Close gracefully closes all resources
func (a *App) Close() error {
	slog.Debug("Closing application resources")

	// First stop all workers with a timeout
	if len(a.workers) > 0 {
		slog.Debug("Stopping workers")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.StopWorkers(ctx); err != nil {
			slog.Warn("Error stopping workers", "error", err)
			// Continue with cleanup despite worker errors
		}
	}

	// Close the database connection
	if a.DB != nil {
		slog.Debug("Closing database connection")
		if err := a.DB.Close(); err != nil {
			slog.Error("Error closing database", "error", err)
			return err
		}
	}

	// Add any other cleanup needed

	return nil
}
