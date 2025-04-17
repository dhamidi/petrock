package core

import (
	"context"
	"errors"
	"fmt"
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
}