# Plan for core/app.go

This file defines the central application initialization logic, handling dependency setup, feature registration, and state replay.

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
  - Initializes application state

- `(a *App) RegisterFeatures()`: Registers all application features
  - Calls feature registration functions
  - Registers message types with the message log
  - Sets up command handlers, query handlers
  - MUST be called before replay

- `(a *App) ReplayLog() error`: Replays the message log to build application state
  - Iterates through all messages starting from version 0
  - Applies commands via registered handlers
  - Panics if a command handler fails during replay (indicates state inconsistency)

- `(a *App) SetupHTTPHandlers(mux *http.ServeMux)`: Registers HTTP handlers for core routes
  - Sets up /commands, /queries endpoints
  - Calls feature handler registration functions

- `(a *App) Close() error`: Gracefully closes all resources
  - Closes database connection
  - Any other cleanup needed