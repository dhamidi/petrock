# core/worker.go

## Overview

The `worker.go` file defines the Worker interface and related types for background processes that react to events in the message log. Workers are long-running processes managed by the App that maintain their own internal state and perform operations that may span multiple events.

## Key Components

### Worker Interface

The core `Worker` interface defines the contract for all background workers:

```go
type Worker interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Work() error
    Replay(ctx context.Context) error
    WorkerInfo() *WorkerInfo  // Optional method
}
```

#### Methods

- **`Start(ctx context.Context) error`**: Initializes the worker. Should be idempotent and return an error if the worker is already started.

- **`Stop(ctx context.Context) error`**: Gracefully shuts down the worker, allowing it to clean up resources and finish any in-progress work. Should be idempotent and respect the provided context's deadline or cancellation.

- **`Work() error`**: Performs a single processing cycle of the worker, handling new messages from the message log since the last processed position, updating internal state, and performing any required actions. Called periodically by the App's worker scheduler and should be quick and non-blocking when possible.

- **`Replay(ctx context.Context) error`**: Processes all messages from the beginning to reconstruct worker state. Called during startup after Start() but before regular Work() cycles. Should process messages for state reconstruction only, avoiding side effects.

- **`WorkerInfo() *WorkerInfo`** (Optional): Provides self-description information for introspection and debugging purposes. If not implemented, information will be extracted via reflection.

### WorkerInfo

Provides optional self-description for workers:

```go
type WorkerInfo struct {
    Name        string // Name of the worker
    Description string // Description of the worker's purpose
}
```

### Error Types

#### WorkerError

Represents an error that occurred during worker operations:

```go
type WorkerError struct {
    Op  string // Operation that failed
    Err error  // Underlying error
}
```

Implements the standard error interface with wrapped error support.

#### Predefined Errors

- **`ErrWorkerStopped`**: Returned when an operation is attempted on a stopped worker
- **`ErrWorkerAlreadyStarted`**: Returned when attempting to start an already running worker

## Usage Patterns

### Implementing a Worker

```go
type MyWorker struct {
    started bool
    // other fields...
}

func (w *MyWorker) Start(ctx context.Context) error {
    if w.started {
        return ErrWorkerAlreadyStarted
    }
    // Initialize worker state by processing message log
    w.started = true
    return nil
}

func (w *MyWorker) Stop(ctx context.Context) error {
    if !w.started {
        return ErrWorkerStopped
    }
    // Clean up resources
    w.started = false
    return nil
}

func (w *MyWorker) Work() error {
    if !w.started {
        return ErrWorkerStopped
    }
    // Perform one cycle of work
    return nil
}

func (w *MyWorker) WorkerInfo() *WorkerInfo {
    return &WorkerInfo{
        Name:        "MyWorker",
        Description: "Processes background tasks for my feature",
    }
}
```

### Worker Lifecycle

1. **Registration**: Workers are registered with the App during feature registration
2. **Startup**: App calls `Start()` on all workers during application initialization  
3. **Replay**: App calls `Replay()` on all workers to reconstruct state from historical messages
4. **Execution**: App periodically calls `Work()` on all workers (typically every 1-2 seconds with jitter)
5. **Shutdown**: App calls `Stop()` on all workers during application shutdown

### Message Log Integration

Workers typically:
- Use `LogFollower` to track their position in the message log with persistent storage via `KVStore`
- React to specific message types relevant to their functionality
- Distinguish between replay (state-only) and normal processing using `ProcessingContext`
- Dispatch commands through the central `core.Executor` when they need to update application state
- Maintain their own internal state by processing messages chronologically

## Design Principles

- **Idempotency**: Start and Stop methods should be safe to call multiple times
- **Non-blocking**: Work method should complete quickly to allow regular scheduling
- **State Isolation**: Workers maintain their own internal state separate from feature state
- **Event-driven**: Workers react to events in the message log rather than polling
- **Graceful Shutdown**: Workers should respect context cancellation and clean up properly
