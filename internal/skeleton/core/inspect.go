package core

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// CommandSchema represents the JSON schema for a command
type CommandSchema struct {
	Name        string                 `json:"name"`        // Command name (e.g., "posts/create")
	Description string                 `json:"description"` // Command description if available
	Type        string                 `json:"type"`        // Go type name
	Properties  map[string]PropertyDef `json:"properties"`  // Field definitions
	Required    []string               `json:"required"`    // Required field names
}

// QuerySchema represents the JSON schema for a query
type QuerySchema struct {
	Name        string                 `json:"name"`        // Query name (e.g., "posts/list")
	Description string                 `json:"description"` // Query description if available
	Type        string                 `json:"type"`        // Go type name
	Properties  map[string]PropertyDef `json:"properties"`  // Field definitions
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

// WorkerSchema represents metadata about a registered worker
type WorkerSchema struct {
	Name        string   `json:"name"`        // Worker name or type
	Description string   `json:"description"` // Worker description if available
	Type        string   `json:"type"`        // Go type name
	Methods     []string `json:"methods"`     // Available methods
}

// InspectResult holds application metadata
type InspectResult struct {
	Commands []CommandSchema `json:"commands"` // Schema of all registered commands
	Queries  []QuerySchema   `json:"queries"`  // Schema of all registered queries
	Routes   []string        `json:"routes"`   // List of all registered HTTP routes
	Features []string        `json:"features"` // List of all registered features
	Workers  []WorkerSchema  `json:"workers"`  // Schema of all registered workers
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
		cmdType, found := a.CommandRegistry.GetCommandType(name)
		if found {
			schema := buildCommandSchema(name, cmdType)
			result.Commands = append(result.Commands, schema)
		}
	}

	// Build query schemas
	queryNames := a.QueryRegistry.RegisteredQueryNames()
	result.Queries = make([]QuerySchema, 0, len(queryNames))
	for _, name := range queryNames {
		queryType, found := a.QueryRegistry.GetQueryType(name)
		if found {
			schema := buildQuerySchema(name, queryType)
			result.Queries = append(result.Queries, schema)
		}
	}

	// Build worker schemas
	result.Workers = make([]WorkerSchema, 0)
	for _, worker := range a.workers {
		schemas := buildWorkerSchemas(worker)
		result.Workers = append(result.Workers, schemas...)
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

// buildWorkerSchemas creates schemas from a worker instance, one for each registered command
func buildWorkerSchemas(worker Worker) []WorkerSchema {
	schemas := []WorkerSchema{}
	
	// Check if this is a CommandWorker that can report its registered commands
	type commandRegistrar interface {
		RegisteredCommands() []string
	}
	
	if cmdWorker, ok := worker.(commandRegistrar); ok {
		// Create a schema for each registered command
		commands := cmdWorker.RegisteredCommands()
		
		if len(commands) == 0 {
			// If no commands registered, fall back to the original behavior
			return []WorkerSchema{buildSingleWorkerSchema(worker)}
		}
		
		for _, commandName := range commands {
			schema := WorkerSchema{
				Name:        commandName, // Use command name instead of worker name
				Description: fmt.Sprintf("Worker handler for %s command", commandName),
				Type:        fmt.Sprintf("%T", worker),
			}
			
			// Add standard worker methods
			schema.Methods = []string{"Start", "Stop", "Work", "OnCommand", "Replay", "SetDependencies", "SetPeriodicWork", "State", "WorkerInfo"}
			
			schemas = append(schemas, schema)
		}
	} else {
		// Fallback for workers that don't implement RegisteredCommands
		schemas = append(schemas, buildSingleWorkerSchema(worker))
	}
	
	return schemas
}

// buildSingleWorkerSchema creates a schema from a worker instance (original behavior)
func buildSingleWorkerSchema(worker Worker) WorkerSchema {
	schema := WorkerSchema{
		Type: fmt.Sprintf("%T", worker),
	}

	// Try to get WorkerInfo if implemented
	if info := worker.WorkerInfo(); info != nil {
		schema.Name = info.Name
		schema.Description = info.Description
	}

	// Try to get worker name from type if not provided via WorkerInfo
	if schema.Name == "" {
		// Extract type name from the fully qualified type
		typeName := schema.Type
		// Find the last dot in the type name (package separator)
		if lastDot := strings.LastIndex(typeName, "."); lastDot != -1 {
			// Extract the type name after the last dot
			typeName = typeName[lastDot+1:]
		}
		// Remove any *pointer prefix
		typeName = strings.TrimPrefix(typeName, "*")
		schema.Name = typeName
	}

	// Extract methods using reflection
	methods := []string{}
	workerType := reflect.TypeOf(worker)
	workerValue := reflect.ValueOf(worker)

	// If it's a pointer, get the element type
	if workerType.Kind() == reflect.Ptr {
		workerType = workerType.Elem()
	}

	// Add the standard Worker interface methods
	methods = append(methods, "Start", "Stop", "Work")

	// Look for additional exported methods
	for i := 0; i < workerValue.Type().NumMethod(); i++ {
		method := workerValue.Type().Method(i)
		// Skip the standard Worker interface methods we already added
		if method.Name == "Start" || method.Name == "Stop" || method.Name == "Work" {
			continue
		}
		// Only include exported methods
		if method.PkgPath == "" {
			methods = append(methods, method.Name)
		}
	}

	schema.Methods = methods
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
