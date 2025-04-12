# App Structure

The `App` struct is the central dependency container for Petrock applications. It holds all core components and provides methods for registering features and routes.

## Structure

```go
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
```

## Key Components

- **DB**: Database connection used by the application
- **MessageLog**: Records all commands for event sourcing
- **CommandRegistry**: Maps command types to their handlers
- **QueryRegistry**: Maps query types to their handlers
- **Executor**: Central component for executing commands
- **Features**: Tracks all registered feature names
- **Routes**: Tracks all registered HTTP routes
- **Mux**: HTTP router for the application
- **AppState**: Generic container for application state

## Registration Methods

### RegisterFeature

Registers a feature with the application:

```go
func (a *App) RegisterFeature(name string, registerFn func(app *App, state interface{})) {
    a.Features = append(a.Features, name)
    registerFn(a, a.AppState)
}
```

This method:
1. Adds the feature name to the tracked features list
2. Calls the registration function with app and state

### RegisterRoute

Registers an HTTP route with the application:

```go
func (a *App) RegisterRoute(pattern string, handler http.HandlerFunc) {
    a.Routes = append(a.Routes, pattern)
    if a.Mux != nil {
        a.Mux.HandleFunc(pattern, handler)
    }
}
```

This method:
1. Adds the route pattern to the tracked routes list
2. Registers the handler with the HTTP mux

## Initialization

The `App` is initialized using the `NewApp` function, which:

1. Creates registries for commands and queries
2. Sets up the database connection
3. Initializes the message log
4. Creates the central executor
5. Initializes empty slices for features and routes

```go
func NewApp(dbPath string) (*App, error) {
    // Initialize components...
    
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
```

## Self Inspection

The App struct provides a method for self-inspection that returns metadata about the application:

```go
func (a *App) GetInspectResult() *InspectResult {
    // Return metadata about commands, queries, routes, and features
}
```

See the [Self Inspection](../self-inspect.md) documentation for details.