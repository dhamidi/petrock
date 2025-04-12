# Implementation Plan for `self inspect` Command

## Overview

This document outlines the step-by-step implementation plan for adding the `self inspect` command to Petrock-generated projects. This feature will enable applications to inspect their internal structure, commands, queries, routes, and features.

## Implementation Steps

### 1. Update Core `App` Struct

**Files to modify:**
- `internal/skeleton/core/app.go`

**Changes:**
1. Add new fields to the `App` struct to track features and routes
2. Add methods for registering features and routes

```go
// Update App struct
type App struct {
    DB              *sql.DB
    MessageLog      *MessageLog
    CommandRegistry *CommandRegistry
    QueryRegistry   *QueryRegistry
    Executor        *Executor
    Features        []string           // NEW: Track registered feature names
    Routes          []string           // NEW: Track registered routes
    Mux             *http.ServeMux     // NEW: Store the HTTP mux
    AppState        interface{}        // Generic application state interface
}

// NEW: Add RegisterFeature method
func (a *App) RegisterFeature(name string, registerFn func(app *App, state interface{})) {
    a.Features = append(a.Features, name)
    registerFn(a, a.AppState)
}

// NEW: Add RegisterRoute method
func (a *App) RegisterRoute(pattern string, handler http.HandlerFunc) {
    a.Routes = append(a.Routes, pattern)
    if a.Mux != nil {
        a.Mux.HandleFunc(pattern, handler)
    }
}
```

**Definition of Done:**
- [x] `App` struct has new fields: `Features`, `Routes`, and `Mux`
- [x] `RegisterFeature` method is implemented and captures feature name
- [x] `RegisterRoute` method is implemented for route tracking
- [x] `NewApp` constructor initializes the new fields with empty slices

### 2. Create Core Inspection Logic

**Files to create:**
- `internal/skeleton/core/inspect.go`

**Implementation:**

```go
package core

import (
    "reflect"
    "strings"
    "time"
)

// CommandSchema represents the JSON schema for a command
type CommandSchema struct {
    Name        string                 `json:"name"`        // Command name (e.g., "posts/create")
    Description string                 `json:"description"` // Command description if available
    Type        string                 `json:"type"`        // Go type name
    Properties  map[string]PropertyDef `json:"properties"` // Field definitions
    Required    []string               `json:"required"`    // Required field names
}

// QuerySchema represents the JSON schema for a query
type QuerySchema struct {
    Name        string                 `json:"name"`        // Query name (e.g., "posts/list")
    Description string                 `json:"description"` // Query description if available
    Type        string                 `json:"type"`        // Go type name
    Properties  map[string]PropertyDef `json:"properties"` // Field definitions
    Required    []string               `json:"required"`    // Required field names
    Result      ResultDef              `json:"result"`      // Schema of the query result
}

// PropertyDef represents a property in a command or query
type PropertyDef struct {
    Type        string      `json:"type"`                  // JSON schema type
    Description string      `json:"description,omitempty"` // Field description
    Format      string      `json:"format,omitempty"`      // Format hint (e.g., date-time)
    Enum        []string    `json:"enum,omitempty"`        // Enum values if applicable
    Default     interface{} `json:"default,omitempty"`     // Default value if any
}

// ResultDef represents a query result
type ResultDef struct {
    Type        string                 `json:"type"`                  // Usually "object"
    Description string                 `json:"description,omitempty"` // Result description
    Properties  map[string]PropertyDef `json:"properties"`            // Result fields
}

// InspectResult holds application metadata
type InspectResult struct {
    Commands []CommandSchema `json:"commands"` // Schema of all registered commands
    Queries  []QuerySchema  `json:"queries"`  // Schema of all registered queries
    Routes   []string       `json:"routes"`   // List of all registered HTTP routes
    Features []string       `json:"features"` // List of all registered features
}

// GetInspectResult gathers metadata about the application
func (a *App) GetInspectResult() *InspectResult {
    result := &InspectResult{
        Routes:   a.Routes,
        Features: a.Features,
    }
    
    // Build command schemas
    commandNames := a.CommandRegistry.RegisteredCommandNames()
    result.Commands = make([]CommandSchema, 0, len(commandNames))
    for _, name := range commandNames {
        cmdType, _ := a.CommandRegistry.GetCommandType(name)
        schema := buildCommandSchema(name, cmdType)
        result.Commands = append(result.Commands, schema)
    }
    
    // Build query schemas
    queryNames := a.QueryRegistry.RegisteredQueryNames()
    result.Queries = make([]QuerySchema, 0, len(queryNames))
    for _, name := range queryNames {
        queryType, _ := a.QueryRegistry.GetQueryType(name)
        schema := buildQuerySchema(name, queryType)
        result.Queries = append(result.Queries, schema)
    }
    
    return result
}

// buildCommandSchema creates a JSON schema from a command's reflect.Type
func buildCommandSchema(name string, cmdType reflect.Type) CommandSchema {
    schema := CommandSchema{
        Name:       name,
        Type:       cmdType.String(),
        Properties: make(map[string]PropertyDef),
        Required:   []string{},
    }
    
    // Extract fields using reflection
    for i := 0; i < cmdType.NumField(); i++ {
        field := cmdType.Field(i)
        if field.PkgPath != "" { // Skip unexported fields
            continue
        }
        
        // Get field name from JSON tag if available
        fieldName := field.Name
        if jsonTag := field.Tag.Get("json"); jsonTag != "" {
            parts := strings.Split(jsonTag, ",")
            if parts[0] != "-" {
                fieldName = parts[0]
            } else {
                continue // Skip fields with json:"-"
            }
        }
        
        // Extract property definition
        propDef := buildPropertyDef(field)
        schema.Properties[fieldName] = propDef
        
        // Mark required fields (simplified approach - all fields are required)
        schema.Required = append(schema.Required, fieldName)
    }
    
    return schema
}

// buildQuerySchema creates a JSON schema from a query's reflect.Type
func buildQuerySchema(name string, queryType reflect.Type) QuerySchema {
    schema := QuerySchema{
        Name:       name,
        Type:       queryType.String(),
        Properties: make(map[string]PropertyDef),
        Required:   []string{},
    }
    
    // Extract fields using reflection
    for i := 0; i < queryType.NumField(); i++ {
        field := queryType.Field(i)
        if field.PkgPath != "" { // Skip unexported fields
            continue
        }
        
        fieldName := field.Name
        if jsonTag := field.Tag.Get("json"); jsonTag != "" {
            parts := strings.Split(jsonTag, ",")
            if parts[0] != "-" {
                fieldName = parts[0]
            } else {
                continue
            }
        }
        
        propDef := buildPropertyDef(field)
        schema.Properties[fieldName] = propDef
        schema.Required = append(schema.Required, fieldName)
    }
    
    // For now, we'll add a placeholder result schema
    // In a real implementation, this would be derived from the query result type
    schema.Result = ResultDef{
        Type:       "object",
        Properties: make(map[string]PropertyDef),
    }
    
    return schema
}

// buildPropertyDef builds a property definition from a struct field
func buildPropertyDef(field reflect.StructField) PropertyDef {
    propDef := PropertyDef{
        Description: field.Tag.Get("description"),
    }
    
    // Map Go types to JSON Schema types
    switch field.Type.Kind() {
    case reflect.String:
        propDef.Type = "string"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
         reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        propDef.Type = "integer"
    case reflect.Float32, reflect.Float64:
        propDef.Type = "number"
    case reflect.Bool:
        propDef.Type = "boolean"
    case reflect.Struct:
        if field.Type == reflect.TypeOf(time.Time{}) {
            propDef.Type = "string"
            propDef.Format = "date-time"
        } else {
            propDef.Type = "object"
        }
    case reflect.Slice, reflect.Array:
        propDef.Type = "array"
    default:
        propDef.Type = "object"
    }
    
    return propDef
}
```

**Definition of Done:**
- [x] `CommandSchema`, `QuerySchema`, `PropertyDef`, and `ResultDef` structs are defined
- [x] `InspectResult` struct to hold complete application metadata
- [x] `App.GetInspectResult()` method to gather and return inspection data
- [x] Helper functions for building command/query schemas using reflection

### 3. Create Self Command Templates

**Files to create:**
- `internal/skeleton/cmd_templates/self.go.tmpl`

**Implementation:**

```go
package {{.Package}}

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "{{.ModulePath}}/core"
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

    // Initialize the application
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

    // Register features
    RegisterAllFeatures(app)

    // We don't need to replay the log since we're only inspecting structure

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

**Definition of Done:**
- [x] Self command template created with proper template variables
- [x] Parent `self` command and child `inspect` subcommand defined
- [x] Command flags defined (`--format`, `--db-path`)
- [x] Implementation of `runSelfInspect` to gather and output metadata

### 4. Update Feature Registration Pattern in Templates

**Files to modify:**
- `internal/skeleton/cmd_templates/features.go.tmpl`

**Original code:**
```go
func RegisterAllFeatures(
    mux *http.ServeMux,
    commands *core.CommandRegistry,
    queries *core.QueryRegistry,
    messageLog *core.MessageLog,
    executor *core.Executor,
    appState *AppState,
    db *sql.DB,
) {
    // Feature registrations...
}
```

**Updated code:**
```go
func RegisterAllFeatures(app *core.App) {
    // Feature registrations will be updated by Petrock
    // petrock:register-feature - Do not remove or modify this line
}
```

**Definition of Done:**
- [x] Updated `RegisterAllFeatures` function to accept `app *core.App` instead of individual dependencies
- [x] Maintained the marker comment for Petrock code generation

### 5. Update Serve Command Template

**Files to modify:**
- `internal/skeleton/cmd_templates/serve.go.tmpl`

**Changes:** 
1. Update how mux is created and passed to features
2. Update core routes registration to use `app.RegisterRoute`

**Key code changes:**
```go
// Initialize the application
app, err := core.NewApp(dbPath)
if err != nil {
    return fmt.Errorf("failed to initialize application: %w", err)
}
defer app.Close()

// Initialize Application State
appState := NewAppState()
app.AppState = appState

// Create HTTP mux
app.Mux = http.NewServeMux()

// Register features
RegisterAllFeatures(app)

// Replay the message log
if err := app.ReplayLog(); err != nil {
    return fmt.Errorf("failed to replay message log: %w", err)
}

// Register core HTTP routes
app.RegisterRoute("GET /", core.HandleIndex(app.CommandRegistry, app.QueryRegistry))
app.RegisterRoute("GET /commands", handleListCommands(app.CommandRegistry))
// ...other routes...
```

**Definition of Done:**
- [x] Updated initialization of the HTTP mux in the `App` struct
- [x] Updated feature registration to use the `App`-based pattern
- [x] Changed route registration to use `app.RegisterRoute`

### 6. Update Feature Template Registration

**Files to modify:**
- `internal/skeleton/feature_template/register.go`

**Original code:**
```go
func RegisterFeature(
    mux *http.ServeMux,
    commands *core.CommandRegistry,
    queries *core.QueryRegistry,
    messageLog *core.MessageLog,
    executor *core.Executor,
    state *State,
    db *sql.DB,
) {
    // Registration logic...
}
```

**Updated code:**
```go
func RegisterFeature(app *core.App, state *State) {
    // Register command types for deserialization
    app.MessageLog.RegisterType("petrock_example_feature_name/create", reflect.TypeOf(CreateCommand{}))
    app.MessageLog.RegisterType("petrock_example_feature_name/update", reflect.TypeOf(UpdateCommand{}))
    app.MessageLog.RegisterType("petrock_example_feature_name/delete", reflect.TypeOf(DeleteCommand{}))

    // Register command handlers
    app.CommandRegistry.Register("petrock_example_feature_name/create", handleCreate(state), reflect.TypeOf(CreateCommand{}))
    app.CommandRegistry.Register("petrock_example_feature_name/update", handleUpdate(state), reflect.TypeOf(UpdateCommand{}))
    app.CommandRegistry.Register("petrock_example_feature_name/delete", handleDelete(state), reflect.TypeOf(DeleteCommand{}))

    // Register query handlers
    app.QueryRegistry.Register("petrock_example_feature_name/get", handleGet(state), reflect.TypeOf(GetQuery{}))
    app.QueryRegistry.Register("petrock_example_feature_name/list", handleList(state), reflect.TypeOf(ListQuery{}))

    // Register HTTP routes
    app.RegisterRoute("GET /petrock_example_feature_name", handleIndex(app.Executor, state))
    app.RegisterRoute("GET /petrock_example_feature_name/{id}", handleGetHTTP(app.Executor, state))
}
```

**Definition of Done:**
- [x] Updated `RegisterFeature` function to accept `app *core.App` instead of individual dependencies
- [x] Modified command/query/message registrations to use the app object
- [x] Changed route registration to use `app.RegisterRoute`

### 7. Update Petrock Code Generation

**Files to modify:**
- `cmd/petrock/new.go` (or equivalent generator file)

**Changes:**
1. Update the generation of `main.go` to include the `NewSelfCmd()`
2. Update how feature registration code is generated

**Code snippet for feature registration generation:**
```go
featureRegistration := fmt.Sprintf(
    "app.RegisterFeature(\"%s\", func(a *core.App, appState interface{}) {\n"+
    "    %sState := %s.NewState()\n"+
    "    %s.RegisterFeature(a, %sState)\n"+
    "})\n",
    featureName, featureName, featureName, featureName, featureName)
```

**Definition of Done:**
- [x] Updated code generation for `RegisterAllFeatures` to use the new App pattern
- [x] Added code generation for the `self` command in the main package
- [x] Ensured that the new `App` approach is properly used in generated code

### 8. Update Documentation

**Files to modify/create:**
- `docs/self-inspect.md` (new file documenting the feature)
- `docs/core/app.md` (update to reflect new App structure)

**Content for `docs/self-inspect.md`:**
```markdown
# Self Inspection

Petrock applications include a `self inspect` command that provides detailed information about the application structure, including:

- All registered commands with their schema
- All registered queries with their schema
- All HTTP routes
- All features

## Usage

```shell
# Basic usage (outputs JSON to stdout)
$ myapp self inspect

# Specify database path
$ myapp self inspect --db-path=custom.db
```

## Output Format

The command outputs structured JSON that describes the application. Here's an example snippet:

```json
{
  "commands": [
    {
      "name": "posts/create",
      "type": "posts.CreatePostCommand",
      "properties": {
        "title": { "type": "string" },
        "content": { "type": "string" },
        "authorID": { "type": "string" }
      },
      "required": ["title", "content"]
    }
  ],
  "queries": [...],
  "routes": [...],
  "features": [...]
}
```

## Use Cases

- API discovery and documentation
- Client code generation
- Integration testing
- Infrastructure validation
```

**Definition of Done:**
- [x] New documentation file for the `self inspect` command
- [x] Updated documentation for the App struct and feature registration pattern
- [x] Examples of command usage and output format

### 9. Integration and Testing

**Testing Steps:**
1. Generate a new project using the updated Petrock
2. Add a feature to the project
3. Run the `self inspect` command
4. Verify the output includes all expected information

**Definition of Done:**
- [x] Command works in a newly generated project
- [x] Command accurately displays all registered commands, queries, routes, and features
- [x] JSON output is properly formatted and contains schema information
- [x] Command handles errors gracefully (e.g., invalid format, missing database)

## Summary

This implementation plan outlines the steps required to add the `self inspect` command to Petrock-generated applications. The feature will allow applications to introspect their structure, providing detailed information about commands, queries, routes, and features in a structured JSON format.

The plan includes both code changes and documentation updates, with clear definitions of done for each step.