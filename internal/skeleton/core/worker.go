package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// WorkerError represents an error that occurred during worker operations.
type WorkerError struct {
	Op  string // Operation that failed
	Err error  // Underlying error
}

func (e *WorkerError) Error() string {
	return fmt.Sprintf("worker %s: %v", e.Op, e.Err)
}

func (e *WorkerError) Unwrap() error {
	return e.Err
}

// ErrWorkerStopped is returned when an operation is attempted on a stopped worker.
var ErrWorkerStopped = errors.New("worker has been stopped")

// ErrWorkerAlreadyStarted is returned when attempting to start an already running worker.
var ErrWorkerAlreadyStarted = errors.New("worker is already started")

// WorkerInfo provides optional self-description for workers
type WorkerInfo struct {
	Name        string // Name of the worker
	Description string // Description of the worker's purpose
}

// Worker defines the interface for background processes that react to events
// in the message log. Workers are responsible for maintaining their own internal
// state by tracking events and performing operations that may span multiple events,
// often interacting with external systems.
//
// Workers are managed by the App, which starts them in separate goroutines,
// schedules their Work() method periodically, and stops them during shutdown.
type Worker interface {
	// Start initializes the worker and rebuilds its internal state by
	// processing existing messages from the message log. This method
	// should be idempotent and return an error if the worker is already started.
	//
	// Example usage:
	//   err := worker.Start(ctx)
	//   if err != nil {
	//     return fmt.Errorf("failed to start worker: %w", err)
	//   }
	Start(ctx context.Context) error

	// Stop gracefully shuts down the worker, allowing it to clean up resources
	// and finish any in-progress work. This method should be idempotent and
	// should respect the provided context's deadline or cancellation.
	//
	// Example usage:
	//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//   defer cancel()
	//   err := worker.Stop(ctx)
	//   if err != nil {
	//     return fmt.Errorf("failed to stop worker cleanly: %w", err)
	//   }
	Stop(ctx context.Context) error

	// Work performs a single processing cycle of the worker, handling new messages
	// from the message log, updating internal state, and performing any required actions.
	// This method is called periodically by the App's worker scheduler and should be
	// designed to be quick and non-blocking when possible.
	//
	// Example usage:
	//   err := worker.Work()
	//   if err != nil {
	//     log.Printf("worker cycle failed: %v", err)
	//   }
	Work() error

	// Replay processes all messages from the beginning to reconstruct worker state.
	// This method is called during startup after Start() but before regular Work() cycles.
	// It should process messages for state reconstruction only, avoiding side effects.
	//
	// Example usage:
	//   err := worker.Replay(ctx)
	//   if err != nil {
	//     return fmt.Errorf("failed to replay messages: %w", err)
	//   }
	Replay(ctx context.Context) error

	// WorkerInfo is an optional method that provides self-description information
	// for introspection and debugging purposes. If not implemented, information
	// will be extracted via reflection.
	//
	// Example implementation:
	//   func (w *MyWorker) WorkerInfo() *WorkerInfo {
	//     return &WorkerInfo{
	//       Name: "MyWorker",
	//       Description: "Processes background tasks for my feature",
	//     }
	//   }
	//
	// Note: This method is optional. Workers don't need to implement it.
	WorkerInfo() *WorkerInfo
}

// CommandWorker provides a command-based worker abstraction that handles
// the common infrastructure for message processing, command routing, and state management.
type CommandWorker struct {
	name         string
	description  string
	state        interface{}
	handlers     map[string]CommandHandler
	periodicWork func(context.Context) error
	log          *MessageLog
	executor     *Executor
	follower     LogFollower
	kvStore      KVStore
	ctx          context.Context
	cancel       context.CancelFunc
	started      bool
}

// NewWorker creates a new command-based worker with the given name, description, and initial state
func NewWorker(name, description string, initialState interface{}) *CommandWorker {
	return &CommandWorker{
		name:         name,
		description:  description,
		state:        initialState,
		handlers:     make(map[string]CommandHandler),
		periodicWork: nil,
		follower:     NewLogFollower(),
	}
}

// OnCommand registers a command handler for the specified command name
func (w *CommandWorker) OnCommand(commandName string, handler CommandHandler) {
	w.handlers[commandName] = handler
}

// SetPeriodicWork sets the periodic work function that gets called during Work() cycles
func (w *CommandWorker) SetPeriodicWork(fn func(context.Context) error) {
	w.periodicWork = fn
}

// State returns the worker's current state
func (w *CommandWorker) State() interface{} {
	return w.state
}

// RegisteredCommands returns a list of command names this worker handles
func (w *CommandWorker) RegisteredCommands() []string {
	commands := make([]string, 0, len(w.handlers))
	for commandName := range w.handlers {
		commands = append(commands, commandName)
	}
	return commands
}

// SetDependencies sets the core dependencies needed by the worker
func (w *CommandWorker) SetDependencies(log *MessageLog, executor *Executor, kvStore KVStore) {
	w.log = log
	w.executor = executor
	w.kvStore = kvStore
}

// positionKey returns the KVStore key for this worker's position
func (w *CommandWorker) positionKey() string {
	return fmt.Sprintf("worker:%s:position", w.name)
}

// Start initializes the worker
func (w *CommandWorker) Start(ctx context.Context) error {
	if w.started {
		return ErrWorkerAlreadyStarted
	}

	w.ctx, w.cancel = context.WithCancel(ctx)
	w.started = true

	slog.Info("Worker started", "name", w.name)
	return nil
}

// Replay processes all messages from the beginning to reconstruct worker state
func (w *CommandWorker) Replay(ctx context.Context) error {
	if w.kvStore == nil {
		return fmt.Errorf("kvStore not set for worker %s", w.name)
	}

	// Load last known position from KV store - cast the follower to access methods
	if follower, ok := w.follower.(*SimpleLogFollower); ok {
		if err := follower.LoadPosition(w.kvStore, w.positionKey()); err != nil {
			// If no position found, start from beginning
			w.follower.LogSeek(0)
		}
	}

	startPosition := w.follower.LogPosition()
	slog.Info("Starting message replay", "worker", w.name, "fromPosition", startPosition)

	// Replay ALL messages from beginning to reconstruct state
	messageCount := 0
	var lastMessageID uint64
	for msg := range w.log.After(ctx, 0) { // Always start from beginning for state reconstruction
		if err := w.replayMessage(msg); err != nil {
			slog.Error("Failed to replay message", "worker", w.name, "id", msg.ID, "error", err)
			continue
		}

		lastMessageID = msg.ID
		messageCount++
	}

	// Set position to latest message ID to avoid reprocessing
	if messageCount > 0 {
		w.follower.LogSeek(lastMessageID)
	}

	// Save final position
	if follower, ok := w.follower.(*SimpleLogFollower); ok {
		if err := follower.SavePosition(w.kvStore, w.positionKey()); err != nil {
			return fmt.Errorf("failed to save final position: %w", err)
		}
	}

	slog.Info("Message replay completed", "worker", w.name, "messagesProcessed", messageCount, "finalPosition", w.follower.LogPosition())
	return nil
}

// replayMessage processes a message for state reconstruction only (no side effects)
func (w *CommandWorker) replayMessage(msg PersistedMessage) error {
	cmd, ok := msg.DecodedPayload.(Command)
	if !ok {
		return nil
	}

	handler, found := w.handlers[cmd.CommandName()]
	if !found {
		return nil
	}

	// Create replay context to indicate this is state-only processing
	replayCtx := &ProcessingContext{IsReplay: true}

	return handler(w.ctx, cmd, &msg.Message, replayCtx)
}

// Stop gracefully shuts down the worker
func (w *CommandWorker) Stop(ctx context.Context) error {
	if !w.started {
		return nil
	}

	if w.cancel != nil {
		w.cancel()
	}
	w.started = false

	slog.Info("Worker stopped", "name", w.name)
	return nil
}

// Work processes new messages and executes periodic work
func (w *CommandWorker) Work() error {
	if !w.started {
		return ErrWorkerStopped
	}

	// Process only new messages (after current position)
	currentPosition := w.follower.LogPosition()
	messagesProcessed := 0

	for msg := range w.log.After(w.ctx, currentPosition) {
		if err := w.processMessage(msg); err != nil {
			slog.Error("Failed to process message", "worker", w.name, "id", msg.ID, "error", err)
			continue
		}

		w.follower.LogSeek(msg.ID)
		messagesProcessed++
	}

	// Persist position if we processed any messages
	if messagesProcessed > 0 && w.kvStore != nil {
		if follower, ok := w.follower.(*SimpleLogFollower); ok {
			if err := follower.SavePosition(w.kvStore, w.positionKey()); err != nil {
				slog.Error("Failed to save worker position", "worker", w.name, "error", err)
			}
		}
	}

	// Execute periodic work if defined
	if w.periodicWork != nil {
		if err := w.periodicWork(w.ctx); err != nil {
			return fmt.Errorf("periodic work failed: %w", err)
		}
	}

	return nil
}

// WorkerInfo returns worker information
func (w *CommandWorker) WorkerInfo() *WorkerInfo {
	return &WorkerInfo{
		Name:        w.name,
		Description: w.description,
	}
}

// processMessage handles normal message processing with side effects
func (w *CommandWorker) processMessage(msg PersistedMessage) error {
	cmd, ok := msg.DecodedPayload.(Command)
	if !ok {
		return nil
	}

	handler, found := w.handlers[cmd.CommandName()]
	if !found {
		return nil
	}

	// Create normal processing context (allows side effects)
	normalCtx := &ProcessingContext{IsReplay: false}

	return handler(w.ctx, cmd, &msg.Message, normalCtx)
}