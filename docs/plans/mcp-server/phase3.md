# Phase 3: Petrock Generator Tools Integration

## Overview

Replace the demo addition tool with real petrock generators, making the MCP server a powerful development assistant for petrock applications. This phase integrates the command, query, and worker generators as MCP tools.

## Goals

1. Replace demo tool with petrock component generators
2. Support field customization for command and query generators
3. Provide comprehensive parameter validation and help
4. Enable seamless component generation from MCP clients (Claude Desktop, etc.)

## Generator Tools Integration

### 1. Command Generator Tool

```go
func (s *MCPServer) registerCommandGeneratorTool() {
    commandTool := Tool{
        Name:        "generate_command",
        Description: "Generate a new command component with optional custom fields",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]interface{}{
                "feature_name": map[string]interface{}{
                    "type":        "string",
                    "description": "Name of the feature (e.g., 'posts', 'users')",
                },
                "entity_name": map[string]interface{}{
                    "type":        "string", 
                    "description": "Name of the entity/command (e.g., 'create', 'delete')",
                },
                "fields": map[string]interface{}{
                    "type":        "array",
                    "description": "Optional custom fields for the command struct",
                    "items": map[string]interface{}{
                        "type": "object",
                        "properties": map[string]interface{}{
                            "name": map[string]interface{}{
                                "type": "string",
                                "description": "Field name",
                            },
                            "type": map[string]interface{}{
                                "type": "string", 
                                "description": "Go type (e.g., 'string', 'int', 'time.Time')",
                            },
                        },
                        "required": []string{"name", "type"},
                    },
                },
            },
            Required: []string{"feature_name", "entity_name"},
        },
        Handler: s.handleGenerateCommand,
    }
    s.RegisterTool(commandTool)
}

func (s *MCPServer) handleGenerateCommand(params map[string]interface{}) (interface{}, error) {
    // Extract parameters
    featureName, ok := params["feature_name"].(string)
    if !ok {
        return nil, fmt.Errorf("feature_name must be a string")
    }
    
    entityName, ok := params["entity_name"].(string) 
    if !ok {
        return nil, fmt.Errorf("entity_name must be a string")
    }
    
    // Extract optional fields
    var fields []generator.CommandField
    if fieldsParam, exists := params["fields"]; exists {
        // Parse custom fields array
        fields = parseCommandFields(fieldsParam)
    }
    
    // Generate command using the existing generator
    cmdGen := generator.NewCommandGenerator(".")
    err := cmdGen.GenerateCommandComponentWithFields(featureName, entityName, ".", s.modulePath, fields)
    if err != nil {
        return nil, fmt.Errorf("failed to generate command: %w", err)
    }
    
    return map[string]interface{}{
        "success": true,
        "message": fmt.Sprintf("Generated command %s/%s", featureName, entityName),
        "files_created": []string{
            fmt.Sprintf("%s/commands/%s.go", featureName, entityName),
        },
    }, nil
}
```

### 2. Query Generator Tool

```go
func (s *MCPServer) registerQueryGeneratorTool() {
    queryTool := Tool{
        Name:        "generate_query", 
        Description: "Generate a new query component with optional custom fields",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]interface{}{
                "feature_name": map[string]interface{}{
                    "type":        "string",
                    "description": "Name of the feature (e.g., 'posts', 'users')",
                },
                "entity_name": map[string]interface{}{
                    "type":        "string",
                    "description": "Name of the query entity (e.g., 'summary', 'details')",
                },
                "fields": map[string]interface{}{
                    "type":        "array", 
                    "description": "Optional custom fields for the query struct",
                    "items": map[string]interface{}{
                        "type": "object",
                        "properties": map[string]interface{}{
                            "name": map[string]interface{}{
                                "type": "string",
                                "description": "Field name",
                            },
                            "type": map[string]interface{}{
                                "type": "string",
                                "description": "Go type (e.g., 'string', 'int', 'time.Time')",
                            },
                        },
                        "required": []string{"name", "type"},
                    },
                },
            },
            Required: []string{"feature_name", "entity_name"},
        },
        Handler: s.handleGenerateQuery,
    }
    s.RegisterTool(queryTool)
}
```

### 3. Worker Generator Tool

```go
func (s *MCPServer) registerWorkerGeneratorTool() {
    workerTool := Tool{
        Name:        "generate_worker",
        Description: "Generate a new worker component for background processing",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]interface{}{
                "feature_name": map[string]interface{}{
                    "type":        "string",
                    "description": "Name of the feature (e.g., 'notifications', 'exports')",
                },
                "worker_name": map[string]interface{}{
                    "type":        "string",
                    "description": "Name of the worker (e.g., 'email', 'csv_export')",
                },
            },
            Required: []string{"feature_name", "worker_name"},
        },
        Handler: s.handleGenerateWorker,
    }
    s.RegisterTool(workerTool)
}
```

### 4. Universal Component Generator Tool

```go
func (s *MCPServer) registerComponentGeneratorTool() {
    componentTool := Tool{
        Name:        "generate_component",
        Description: "Universal component generator supporting all types (command, query, worker)",
        InputSchema: ToolInputSchema{
            Type: "object", 
            Properties: map[string]interface{}{
                "component_type": map[string]interface{}{
                    "type":        "string",
                    "description": "Type of component to generate",
                    "enum":        []string{"command", "query", "worker"},
                },
                "feature_name": map[string]interface{}{
                    "type":        "string",
                    "description": "Name of the feature",
                },
                "entity_name": map[string]interface{}{
                    "type":        "string", 
                    "description": "Name of the entity/component",
                },
                "fields": map[string]interface{}{
                    "type":        "array",
                    "description": "Optional custom fields (for commands and queries)",
                },
            },
            Required: []string{"component_type", "feature_name", "entity_name"},
        },
        Handler: s.handleGenerateComponent,
    }
    s.RegisterTool(componentTool)
}
```

## Integration with Existing Infrastructure

### 1. Module Path Detection

```go
// MCPServer needs to know the current module path for generation
type MCPServer struct {
    jsonrpcServer *JSONRPCServer
    capabilities  ServerCapabilities
    tools         map[string]Tool
    modulePath    string  // Detected from go.mod
    projectRoot   string  // Working directory
}

func (s *MCPServer) detectModulePath() error {
    // Read go.mod to determine module path
    // Set s.modulePath for use in generators
}
```

### 2. File Path Validation

```go
func (s *MCPServer) validateProjectStructure() error {
    // Ensure we're in a petrock project
    // Check for required directories and files
    // Validate that generators can run successfully
}
```

### 3. Error Handling and User Feedback

```go
func (s *MCPServer) handleGenerateCommand(params map[string]interface{}) (interface{}, error) {
    // Validate parameters
    // Check for existing components (avoid collisions)
    // Run generator with proper error handling
    // Return detailed success/failure information
    
    if err != nil {
        return map[string]interface{}{
            "success": false,
            "error":   err.Error(),
            "suggestion": "Make sure you're in a petrock project root and the feature name is valid",
        }, nil // Return as successful MCP response, error in content
    }
    
    return map[string]interface{}{
        "success": true,
        "message": "Component generated successfully",
        "files_created": generatedFiles,
        "next_steps": []string{
            "Run './build.sh' to compile the project",
            "Check the generated files for any needed customizations",
        },
    }, nil
}
```

## Advanced Features

### 1. Component Introspection Tools

```go
// Tool to list existing components
func (s *MCPServer) registerListComponentsTool() {
    // List all existing commands, queries, workers in the project
}

// Tool to show component details  
func (s *MCPServer) registerInspectComponentTool() {
    // Show details about a specific component
}
```

### 2. Template Customization Support

```go
// Support for custom field types and validation
func parseCommandFields(fieldsParam interface{}) []generator.CommandField {
    // Parse and validate field definitions
    // Support common Go types and custom types
    // Provide helpful error messages for invalid types
}
```

### 3. Project Status Tool

```go
func (s *MCPServer) registerProjectStatusTool() {
    // Show project structure
    // List features and their components  
    // Show build status and any issues
}
```

## Testing Strategy

1. **Generator integration tests** - Test each tool with real parameters
2. **Field parsing tests** - Validate custom field parsing logic  
3. **Error handling tests** - Test invalid parameters and project states
4. **End-to-end tests** - Generate components and verify they compile
5. **MCP client tests** - Test with actual MCP clients like Claude Desktop

## Files to Modify/Create

- `internal/skeleton/core/mcp.go` - Add generator tools
- `internal/skeleton/core/mcp_generators.go` - Generator tool implementations
- `internal/skeleton/core/mcp_test.go` - Updated tests
- `internal/skeleton/cmd/petrock_example_project_name/mcp.go` - Update command setup

## Success Criteria

1. All four generator tools work correctly from MCP clients
2. Custom fields are properly parsed and applied
3. Generated components compile successfully  
4. Error messages are helpful and actionable
5. Integration with existing petrock workflow is seamless
6. Performance is acceptable for interactive use

## Future Enhancements

- Feature scaffolding tool (generate complete feature with multiple components)
- Database migration generator integration
- UI component generator integration  
- Project template customization tools
- Code analysis and refactoring tools
