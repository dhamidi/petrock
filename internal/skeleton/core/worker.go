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
	name            string
	description     string
	lastProcessedID uint64
	state           interface{}
	handlers        map[string]CommandHandler
	periodicWork    func(context.Context) error
	log             *MessageLog
	executor        *Executor
	ctx             context.Context
	cancel          context.CancelFunc
	started         bool
}

// NewWorker creates a new command-based worker with the given name, description, and initial state
func NewWorker(name, description string, initialState interface{}) *CommandWorker {
	return &CommandWorker{
		name:         name,
		description:  description,
		state:        initialState,
		handlers:     make(map[string]CommandHandler),
		periodicWork: nil,
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

// SetDependencies sets the core dependencies needed by the worker
func (w *CommandWorker) SetDependencies(log *MessageLog, executor *Executor) {
	w.log = log
	w.executor = executor
}

// Start initializes the worker
func (w *CommandWorker) Start(ctx context.Context) error {
	if w.started {
		return ErrWorkerAlreadyStarted
	}

	w.ctx, w.cancel = context.WithCancel(ctx)
	w.lastProcessedID = 0
	w.started = true

	slog.Info("Worker started", "name", w.name)
	return nil
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

	// Process new messages
	if err := w.processNewMessages(); err != nil {
		return fmt.Errorf("failed to process messages: %w", err)
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

// processNewMessages handles new messages from the message log
func (w *CommandWorker) processNewMessages() error {
	if w.log == nil {
		return nil
	}

	// Process messages after lastProcessedID
	for msg := range w.log.After(w.ctx, w.lastProcessedID) {
		if err := w.processMessage(msg); err != nil {
			slog.Error("Failed to process message", "id", msg.ID, "error", err)
			continue
		}
		w.lastProcessedID = msg.ID
	}

	return nil
}

// processMessage processes a single message
func (w *CommandWorker) processMessage(msg PersistedMessage) error {
	// Only process commands
	cmd, ok := msg.DecodedPayload.(Command)
	if !ok {
		return nil
	}

	// Find handler for this command
	handler, found := w.handlers[cmd.CommandName()]
	if !found {
		return nil
	}

	// Execute the handler
	return handler(w.ctx, cmd, &msg.Message)
}