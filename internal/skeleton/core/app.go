package core

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
)

// App is the central struct that holds all application dependencies and state
type App struct {
	DB              *sql.DB
	MessageLog      *MessageLog
	CommandRegistry *CommandRegistry
	QueryRegistry   *QueryRegistry
	Executor        *Executor
	Features        []string         // Track registered feature names
	Routes          []string         // Track registered routes
	Mux             *http.ServeMux   // Store the HTTP mux
	AppState        interface{}      // Generic application state interface
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

// ReplayLog replays the message log to build application state
func (a *App) ReplayLog() error {
	slog.Info("Replaying message log to build application state...")
	
	// Get the starting version
	startVersion := uint64(0) // Start from the beginning
	replayCtx := context.Background() // Use a background context for replay
	replayErrors := 0 // Count errors during replay

	// Use the iterator to process messages one by one
	messageCount := 0
	for msg := range a.MessageLog.After(replayCtx, startVersion) {
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