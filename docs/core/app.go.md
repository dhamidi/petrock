# Plan for core/app.go

This file defines the central application initialization logic, handling dependency setup, feature registration, and state replay. The App struct has no HTTP/web server concerns; it focuses solely on core business logic and application state management.

## Types

- `Worker`: Interface for background workers
  - `Start(context.Context)`: Called to initialize worker state from the message log
  - `Stop(context.Context)`: Called to clean up resources when shutting down
  - `Work() error`: Called periodically to process events and perform work

- `App`: Central struct holding all application dependencies and state
  - Fields: db, messageLog, commandRegistry, queryRegistry, executor, appState, workers
  - Methods for initialization, registration, replay, and worker management

## Functions

- `NewApp(dbPath string) (*App, error)`: Creates and initializes all core dependencies
  - Sets up database connection
  - Initializes encoder
  - Creates message log, command registry, query registry
  - Creates executor
  - Returns the initialized app (appState will be set by caller)

- `(a *App) RegisterFeatures(appState interface{})`: Registers all application features
  - Takes appState as a parameter for features to access
  - Registers message types with the message log
  - Sets up command handlers, query handlers
  - MUST be called before replay for proper message deserialization

- `(a *App) ReplayLog() error`: Replays the message log to build application state
  - Iterates through all messages starting from version 0
  - Applies commands via registered handlers
  - Panics if a command handler fails during replay (indicates state inconsistency)

- `(a *App) RegisterWorker(worker Worker)`: Registers a worker with the application
  - Adds the worker to the list of managed workers

- `(a *App) StartWorkers(ctx context.Context)`: Starts all registered workers
  - Initializes each worker by calling its Start method in a separate goroutine
  - Sets up periodic execution of the Work method with jitter

- `(a *App) StopWorkers(ctx context.Context)`: Stops all workers gracefully
  - Calls the Stop method on each worker
  - Waits for all worker goroutines to finish

- `(a *App) Close() error`: Gracefully closes all resources
  - Stops all workers
  - Closes database connection
  - Any other cleanup needed