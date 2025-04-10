# Plan for core/app.go

This file defines the central application initialization logic, handling dependency setup, feature registration, and state replay. The App struct has no HTTP/web server concerns; it focuses solely on core business logic and application state management.

## Types

- `App`: Central struct holding all application dependencies and state
  - Fields: db, messageLog, commandRegistry, queryRegistry, executor, appState
  - Methods for initialization, registration, and replay

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

- `(a *App) Close() error`: Gracefully closes all resources
  - Closes database connection
  - Any other cleanup needed