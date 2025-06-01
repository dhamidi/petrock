# core/inspect.go

## Overview

The `inspect.go` file provides runtime inspection and debugging tools for Petrock applications. It enables introspection of the application's registered commands, queries, routes, features, and workers through reflection-based schema generation.

## Key Components

### Schema Types

#### CommandSchema

Represents the JSON schema for a command:

```go
type CommandSchema struct {
    Name        string                 `json:"name"`        // Command name (e.g., "posts/create")
    Description string                 `json:"description"` // Command description if available
    Type        string                 `json:"type"`        // Go type name
    Properties  map[string]PropertyDef `json:"properties"`  // Field definitions
    Required    []string               `json:"required"`    // Required field names
}
```

#### QuerySchema

Represents the JSON schema for a query:

```go
type QuerySchema struct {
    Name        string                 `json:"name"`        // Query name (e.g., "posts/list")
    Description string                 `json:"description"` // Query description if available
    Type        string                 `json:"type"`        // Go type name
    Properties  map[string]PropertyDef `json:"properties"`  // Field definitions
    Required    []string               `json:"required"`    // Required field names
    Result      ResultDef              `json:"result"`      // Schema of the query result
}
```

#### PropertyDef

Represents a property in a command or query:

```go
type PropertyDef struct {
    Type        string      `json:"type"`                  // JSON schema type
    Description string      `json:"description,omitempty"` // Field description
    Format      string      `json:"format,omitempty"`      // Format hint (e.g., date-time)
    Enum        []string    `json:"enum,omitempty"`        // Enum values if applicable
    Default     interface{} `json:"default,omitempty"`     // Default value if any
}
```

#### WorkerSchema

Represents metadata about a registered worker:

```go
type WorkerSchema struct {
    Name        string   `json:"name"`        // Worker name or type
    Description string   `json:"description"` // Worker description if available
    Type        string   `json:"type"`        // Go type name
    Methods     []string `json:"methods"`     // Available methods
}
```

#### ResultDef

Represents a query result:

```go
type ResultDef struct {
    Type        string                 `json:"type"`                  // Usually "object"
    Description string                 `json:"description,omitempty"` // Result description
    Properties  map[string]PropertyDef `json:"properties"`            // Result fields
}
```

### InspectResult

The main result structure that holds application metadata:

```go
type InspectResult struct {
    Commands []CommandSchema `json:"commands"` // Schema of all registered commands
    Queries  []QuerySchema   `json:"queries"`  // Schema of all registered queries
    Routes   []string        `json:"routes"`   // List of all registered HTTP routes
    Features []string        `json:"features"` // List of all registered features
    Workers  []WorkerSchema  `json:"workers"`  // Schema of all registered workers
}
```

## Main Functions

### GetInspectResult

The primary method for gathering application metadata:

```go
func (a *App) GetInspectResult() *InspectResult
```

This method:
1. Collects all registered command names and builds schemas using reflection
2. Collects all registered query names and builds schemas using reflection
3. Gathers all registered routes from the app
4. Lists all registered features
5. Builds schemas for all registered workers

### Schema Building Functions

#### buildCommandSchema

Creates a JSON schema from a command's `reflect.Type`:
- Extracts struct fields using reflection
- Maps Go types to JSON Schema types
- Handles JSON tags for field naming
- Marks all fields as required (simplified approach)

#### buildQuerySchema

Creates a JSON schema from a query's `reflect.Type`:
- Similar to command schema building
- Includes placeholder result schema (to be enhanced in real implementations)

#### buildWorkerSchema

Creates a schema from a worker instance:
- Attempts to get worker info via the `WorkerInfo()` method
- Falls back to extracting type name via reflection
- Lists all available methods including standard Worker interface methods

#### buildPropertyDef

Maps Go struct fields to JSON Schema property definitions:
- Handles basic Go types (string, int, float, bool)
- Special handling for `time.Time` (mapped to string with date-time format)
- Maps slices/arrays to JSON array type
- Maps structs to object type

## Type Mapping

The inspection system maps Go types to JSON Schema types as follows:

| Go Type | JSON Schema Type | Notes |
|---------|------------------|-------|
| string | string | |
| int, int8, int16, int32, int64 | integer | |
| uint, uint8, uint16, uint32, uint64 | integer | |
| float32, float64 | number | |
| bool | boolean | |
| time.Time | string | format: "date-time" |
| []T, [N]T | array | |
| struct | object | |

## Usage Examples

### Getting Application Metadata

```go
app := &App{...} // initialized app
result := app.GetInspectResult()

// Access command schemas
for _, cmd := range result.Commands {
    fmt.Printf("Command: %s (%s)\n", cmd.Name, cmd.Type)
    for fieldName, prop := range cmd.Properties {
        fmt.Printf("  %s: %s\n", fieldName, prop.Type)
    }
}

// Access worker information
for _, worker := range result.Workers {
    fmt.Printf("Worker: %s - %s\n", worker.Name, worker.Description)
}
```

### JSON Output

The `InspectResult` can be easily serialized to JSON for API responses or debugging:

```go
import "encoding/json"

result := app.GetInspectResult()
jsonData, err := json.MarshalIndent(result, "", "  ")
if err != nil {
    // handle error
}
fmt.Println(string(jsonData))
```

## Design Principles

- **Reflection-based**: Uses Go's reflection capabilities to introspect types at runtime
- **JSON Schema Compatible**: Generates schemas compatible with JSON Schema specification
- **Extensible**: Easy to add new metadata types and schema enhancements
- **Self-describing**: Workers can provide their own metadata via the `WorkerInfo()` method
- **Debugging-friendly**: Provides comprehensive information for runtime debugging and documentation
