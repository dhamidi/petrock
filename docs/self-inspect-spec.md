# `self inspect` Command Specification

## Overview
This document specifies a new command for projects generated with petrock: `self inspect`. This command will initialize the application and dump information about it in a structured format (JSON by default).

## Command Structure

```
<project> self inspect [flags]
```

For example: `blog self inspect`

### Flags
- `--format`: Control the output format. Default is `json`. Future formats might include `yaml` or `text`.

## Core Design Changes

To make self-inspection easier and more maintainable, we'll make several improvements to the core design:

1. Enhance the `App` struct to track features and routes directly
2. Move feature registration logic into the `App` struct as a method
3. Have features register directly with the App instance

### App Struct Changes

```go
// App is the central struct that holds all application dependencies and state
type App struct {
    DB              *sql.DB
    MessageLog      *MessageLog
    CommandRegistry *CommandRegistry
    QueryRegistry   *QueryRegistry
    Executor        *Executor
    Features        []string           // Track registered feature names
    Routes          []string           // Track registered routes
    Mux             *http.ServeMux     // Store the HTTP mux
    AppState        interface{}        // Generic application state interface
}

// RegisterFeature registers a feature with the application
// This method tracks the feature name and calls the provided registration function
func (a *App) RegisterFeature(name string, registerFn func(app *App, state interface{})) {
    a.Features = append(a.Features, name)
    registerFn(a, a.AppState)
}

// RegisterRoute registers an HTTP route with the application
// This is a wrapper around mux.HandleFunc that tracks the route
func (a *App) RegisterRoute(pattern string, handler http.HandlerFunc) {
    a.Routes = append(a.Routes, pattern)
    a.Mux.HandleFunc(pattern, handler)
}
```

## Implementation Details

### 1. Core Functionality
The core logic for gathering application metadata should be implemented in the `internal/skeleton/core` package, allowing it to be reused later for HTTP endpoints.

#### Create a new file: `internal/skeleton/core/inspect.go`

This file should contain:

- A `InspectResult` struct that holds application metadata
- Functions to gather metadata from App components

```go
package core

import (
    "net/http"
)

// InspectResult represents the application metadata
type InspectResult struct {
    Commands []string `json:"commands"` // List of all registered command names
    Queries  []string `json:"queries"`  // List of all registered query names
    Routes   []string `json:"routes"`   // List of all registered HTTP routes
    Features []string `json:"features"` // List of all registered features
}

// GetInspectResult gathers metadata about the application
func (a *App) GetInspectResult() *InspectResult {
    return &InspectResult{
        Commands: a.CommandRegistry.RegisteredCommandNames(),
        Queries:  a.QueryRegistry.RegisteredQueryNames(),
        Routes:   a.Routes,
        Features: a.Features,
    }
}
```

### 2. Command Implementation

Implement a new cobra command in the generated application's cmd directory.

#### Create a new file: `cmd/<project>/self.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "github.com/<user>/<project>/core"
    "github.com/spf13/cobra"
)

// NewSelfCmd creates the 'self' parent command for introspection commands
func NewSelfCmd() *cobra.Command {
    selfCmd := &cobra.Command{
        Use:   "self",
        Short: "Commands for application self-inspection",
        Long:  `Commands that provide information about the application itself.`,
    }

    // Add subcommands
    selfCmd.AddCommand(NewSelfInspectCmd())

    return selfCmd
}

// NewSelfInspectCmd creates the 'self inspect' command
func NewSelfInspectCmd() *cobra.Command {
    inspectCmd := &cobra.Command{
        Use:   "inspect",
        Short: "Inspect the application structure",
        Long:  `Initializes the application and dumps information about its structure in the specified format.`,
        RunE:  runSelfInspect,
    }

    // Add flags
    inspectCmd.Flags().String("format", "json", "Output format: json")
    inspectCmd.Flags().String("db-path", "app.db", "Path to the SQLite database file")

    return inspectCmd
}

func runSelfInspect(cmd *cobra.Command, args []string) error {
    format, _ := cmd.Flags().GetString("format")
    dbPath, _ := cmd.Flags().GetString("db-path")

    // Only support JSON for now
    if format != "json" {
        return fmt.Errorf("unsupported format: %s (only 'json' is currently supported)", format)
    }

    // Initialize the application (similar to serve.go but without starting HTTP server)
    app, err := core.NewApp(dbPath)
    if err != nil {
        return fmt.Errorf("failed to initialize application: %w", err)
    }
    defer app.Close()

    // Initialize Application State
    appState := NewAppState()
    app.AppState = appState

    // Create HTTP mux for capturing routes
    app.Mux = http.NewServeMux()

    // Register features using the updated app pattern
    RegisterAllFeatures(app)

    // We don't need to replay the log - just inspecting the structure

    // Register core HTTP routes to ensure they're captured
    app.RegisterRoute("GET /", core.HandleIndex(app.CommandRegistry, app.QueryRegistry))
    app.RegisterRoute("GET /commands", handleListCommands(app.CommandRegistry))
    app.RegisterRoute("POST /commands", handleExecuteCommand(app.Executor, app.CommandRegistry))
    app.RegisterRoute("GET /queries", handleListQueries(app.QueryRegistry))
    app.RegisterRoute("GET /queries/{feature}/{queryName}", handleExecuteQuery(app.QueryRegistry))

    // Gather application metadata
    result := app.GetInspectResult()

    // Output as JSON
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(result); err != nil {
        return fmt.Errorf("failed to encode result as JSON: %w", err)
    }

    return nil
}
```

### 3. Update main.go

Update the main.go file to add the new self command:

```go
func init() {
    // Add subcommands here
    rootCmd.AddCommand(NewServeCmd())
    rootCmd.AddCommand(NewBuildCmd())
    rootCmd.AddCommand(NewDeployCmd())
    rootCmd.AddCommand(NewSelfCmd()) // Add the new self command
}
```

### 4. Update RegisterAllFeatures

Modify the `RegisterAllFeatures` function to work with the new App registration pattern:

```go
// RegisterAllFeatures registers handlers and types for all compiled-in features
func RegisterAllFeatures(app *core.App) {
    // The `petrock feature <n>` command will insert code below this line
    // to initialize each feature's state and call its RegisterFeature function
    postsState := posts.NewState()
    app.RegisterFeature("posts", func(a *core.App, appState interface{}) {
        posts.RegisterFeature(a.Mux, a.CommandRegistry, a.QueryRegistry, a.MessageLog, a.Executor, postsState, a.DB)
    })
    // petrock:register-feature - Do not remove or modify this line
}
```

### 5. Update Feature Templates

Update the feature template to work with the new App registration pattern. This will require modifying how petrock generates new features.

## Feature Template Changes

The `RegisterFeature` function in feature templates should be updated to work with the App directly:

```go
// RegisterFeature registers all command handlers, query handlers, and routes for this feature
func RegisterFeature(app *core.App, state *State) {
    // Register command handlers
    app.CommandRegistry.Register("feature/command", handleCommand, reflect.TypeOf(Command{}))
    
    // Register query handlers
    app.QueryRegistry.Register("feature/query", handleQuery, reflect.TypeOf(Query{}))
    
    // Register routes
    app.RegisterRoute("GET /feature", handleFeatureIndex(app.Executor, state))
}
```

## Implementation Challenges

### Feature Registration Refactoring

Changing the feature registration pattern will require updating:
1. The core App struct in `internal/skeleton/core/app.go`
2. The feature registration logic in `cmd/<project>/features.go`
3. The template code used by petrock to generate new features

## Future Enhancements

1. Support additional output formats (YAML, text)
2. Add more inspection details (middleware, dependencies)
3. Add an HTTP endpoint for the same information
4. Add subcommands for more targeted inspection (e.g., `self inspect routes`, `self inspect features`)

## Example Output

```json
{
  "commands": [
    "posts/create",
    "posts/update",
    "posts/delete"
  ],
  "queries": [
    "posts/get",
    "posts/list"
  ],
  "routes": [
    "GET /",
    "GET /commands",
    "POST /commands",
    "GET /queries",
    "GET /queries/{feature}/{queryName}",
    "GET /posts",
    "GET /posts/{id}"
  ],
  "features": [
    "posts"
  ]
}
```