package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MCPServer implements the Model Context Protocol
type MCPServer struct {
	jsonrpcServer *Server
	capabilities  ServerCapabilities
	tools         map[string]Tool
	initialized   bool
	modulePath    string // Current module path for generator commands
	projectRoot   string // Project root directory
}

// ServerCapabilities describes what the server supports
type ServerCapabilities struct {
	Tools   *ToolsCapability `json:"tools,omitempty"`
	Logging *LoggingCapability `json:"logging,omitempty"`
}

// ToolsCapability indicates the server supports tools
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// LoggingCapability indicates the server supports logging
type LoggingCapability struct{}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema ToolInputSchema `json:"inputSchema"`
	Handler     func(params map[string]interface{}) (interface{}, error)
}

// ToolInputSchema defines the expected parameters for a tool
type ToolInputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

// ToolResponseContent represents content in a tool response
type ToolResponseContent struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// MCP Protocol Message Types

// InitializeRequest represents the MCP initialize request
type InitializeRequest struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

// ClientCapabilities describes what the client supports
type ClientCapabilities struct {
	Roots    *RootsCapability    `json:"roots,omitempty"`
	Sampling *SamplingCapability `json:"sampling,omitempty"`
}

// RootsCapability indicates client support for roots
type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// SamplingCapability indicates client support for sampling
type SamplingCapability struct{}

// ClientInfo provides information about the client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResponse represents the MCP initialize response
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ServerInfo provides information about the server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ListToolsRequest represents the tools/list request
type ListToolsRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListToolsResponse represents the tools/list response
type ListToolsResponse struct {
	Tools      []ToolDefinition `json:"tools"`
	NextCursor string           `json:"nextCursor,omitempty"`
}

// ToolDefinition represents a tool definition in the list response
type ToolDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema ToolInputSchema `json:"inputSchema"`
}

// CallToolRequest represents the tools/call request
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResponse represents the tools/call response
type CallToolResponse struct {
	Content []ToolResponseContent `json:"content"`
	IsError bool                  `json:"isError,omitempty"`
}

// NewMCPServer creates a new MCP server
func NewMCPServer() *MCPServer {
	server := &MCPServer{
		jsonrpcServer: NewJSONRPCServer(),
		capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Logging: &LoggingCapability{},
		},
		tools:       make(map[string]Tool),
		initialized: false,
		projectRoot: ".", // Default to current directory
	}

	// Detect module path
	if err := server.detectModulePath(); err != nil {
		// Log error but continue - tools can provide helpful error messages
		fmt.Fprintf(os.Stderr, "Warning: Could not detect module path: %v\n", err)
	}

	server.setupMCPHandlers()
	server.RegisterGeneratorTools()

	return server
}

// setupMCPHandlers registers the MCP protocol handlers
func (s *MCPServer) setupMCPHandlers() {
	// Initialize handlers
	s.jsonrpcServer.RegisterRequestHandler("initialize", s.handleInitialize)
	s.jsonrpcServer.RegisterNotificationHandler("initialized", s.handleInitialized)

	// Tool handlers
	s.jsonrpcServer.RegisterRequestHandler("tools/list", s.handleListTools)
	s.jsonrpcServer.RegisterRequestHandler("tools/call", s.handleCallTool)

	// Utility handlers
	s.jsonrpcServer.RegisterRequestHandler("ping", s.handlePing)
}

// RegisterTool registers a new tool with the server
func (s *MCPServer) RegisterTool(tool Tool) {
	s.tools[tool.Name] = tool
}

// handleInitialize processes the MCP initialize request
func (s *MCPServer) handleInitialize(params interface{}) (interface{}, error) {
	var req InitializeRequest
	if params != nil {
		// Convert params to JSON and back to properly unmarshal
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal initialize params: %w", err)
		}
		if err := json.Unmarshal(jsonData, &req); err != nil {
			return nil, fmt.Errorf("failed to unmarshal initialize request: %w", err)
		}
	}



	response := InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo: ServerInfo{
			Name:    "petrock-mcp-server",
			Version: "1.0.0",
		},
	}

	return response, nil
}

// handleInitialized processes the MCP initialized notification
func (s *MCPServer) handleInitialized(params interface{}) {
	s.initialized = true
}

// handleListTools processes the tools/list request
func (s *MCPServer) handleListTools(params interface{}) (interface{}, error) {

	var req ListToolsRequest
	if params != nil {
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal list tools params: %w", err)
		}
		if err := json.Unmarshal(jsonData, &req); err != nil {
			return nil, fmt.Errorf("failed to unmarshal list tools request: %w", err)
		}
	}

	var tools []ToolDefinition
	for _, tool := range s.tools {
		tools = append(tools, ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	response := ListToolsResponse{
		Tools: tools,
	}

	return response, nil
}

// handleCallTool processes the tools/call request
func (s *MCPServer) handleCallTool(params interface{}) (interface{}, error) {

	var req CallToolRequest
	if params != nil {
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal call tool params: %w", err)
		}
		if err := json.Unmarshal(jsonData, &req); err != nil {
			return nil, fmt.Errorf("failed to unmarshal call tool request: %w", err)
		}
	}

	tool, exists := s.tools[req.Name]
	if !exists {
		return CallToolResponse{
			Content: []ToolResponseContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Tool '%s' not found", req.Name),
				},
			},
			IsError: true,
		}, nil
	}

	result, err := tool.Handler(req.Arguments)
	if err != nil {
		return CallToolResponse{
			Content: []ToolResponseContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Tool execution failed: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil
	}

	// Convert result to text representation
	var content []ToolResponseContent
	if result != nil {
		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			content = append(content, ToolResponseContent{
				Type: "text",
				Text: fmt.Sprintf("Result: %v", result),
			})
		} else {
			content = append(content, ToolResponseContent{
				Type: "text",
				Text: string(jsonResult),
			})
		}
	} else {
		content = append(content, ToolResponseContent{
			Type: "text",
			Text: "Tool executed successfully",
		})
	}

	return CallToolResponse{
		Content: content,
		IsError: false,
	}, nil
}

// handlePing processes ping requests
func (s *MCPServer) handlePing(params interface{}) (interface{}, error) {
	return map[string]interface{}{}, nil
}

// detectModulePath reads go.mod to determine the current module path
func (s *MCPServer) detectModulePath() error {
	goModPath := filepath.Join(s.projectRoot, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			s.modulePath = strings.TrimSpace(strings.TrimPrefix(line, "module"))
			return nil
		}
	}

	return fmt.Errorf("module directive not found in go.mod")
}

// validateProjectStructure ensures we're in a petrock project
func (s *MCPServer) validateProjectStructure() error {
	// Check for go.mod
	if _, err := os.Stat(filepath.Join(s.projectRoot, "go.mod")); os.IsNotExist(err) {
		return fmt.Errorf("not in a Go project directory (go.mod not found)")
	}

	// Check for petrock project indicators
	coreDir := filepath.Join(s.projectRoot, "core")
	if _, err := os.Stat(coreDir); os.IsNotExist(err) {
		return fmt.Errorf("not in a petrock project directory (core/ directory not found)")
	}

	return nil
}

// RegisterGeneratorTools adds the petrock generator tools
func (s *MCPServer) RegisterGeneratorTools() {
	s.RegisterTool(s.createCommandGeneratorTool())
	s.RegisterTool(s.createQueryGeneratorTool())
	s.RegisterTool(s.createWorkerGeneratorTool())
	s.RegisterTool(s.createComponentGeneratorTool())
}

// createCommandGeneratorTool creates the command generator tool
func (s *MCPServer) createCommandGeneratorTool() Tool {
	return Tool{
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
								"type":        "string",
								"description": "Field name",
							},
							"type": map[string]interface{}{
								"type":        "string",
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
}

// createQueryGeneratorTool creates the query generator tool
func (s *MCPServer) createQueryGeneratorTool() Tool {
	return Tool{
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
					"description": "Name of the query entity (e.g., 'get', 'list')",
				},
				"fields": map[string]interface{}{
					"type":        "array",
					"description": "Optional custom fields for the query struct",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type":        "string",
								"description": "Field name",
							},
							"type": map[string]interface{}{
								"type":        "string",
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
}

// createWorkerGeneratorTool creates the worker generator tool
func (s *MCPServer) createWorkerGeneratorTool() Tool {
	return Tool{
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
}

// createComponentGeneratorTool creates the universal component generator tool
func (s *MCPServer) createComponentGeneratorTool() Tool {
	return Tool{
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
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type":        "string",
								"description": "Field name",
							},
							"type": map[string]interface{}{
								"type":        "string",
								"description": "Go type (e.g., 'string', 'int', 'time.Time')",
							},
						},
						"required": []string{"name", "type"},
					},
				},
			},
			Required: []string{"component_type", "feature_name", "entity_name"},
		},
		Handler: s.handleGenerateComponent,
	}
}

// Handler methods for generator tools

// handleGenerateCommand processes the generate_command tool request
func (s *MCPServer) handleGenerateCommand(params map[string]interface{}) (interface{}, error) {
	// Validate project structure
	if err := s.validateProjectStructure(); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"suggestion": "Make sure you're in a petrock project root directory",
		}, nil
	}

	// Extract required parameters
	featureName, ok := params["feature_name"].(string)
	if !ok {
		return nil, fmt.Errorf("feature_name must be a string")
	}

	entityName, ok := params["entity_name"].(string)
	if !ok {
		return nil, fmt.Errorf("entity_name must be a string")
	}

	// Build command arguments
	args := []string{"new", "command", fmt.Sprintf("%s/%s", featureName, entityName)}

	// Extract optional fields
	if fieldsParam, exists := params["fields"]; exists {
		fields, err := s.parseFields(fieldsParam)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Invalid fields: %v", err),
			}, nil
		}
		// Add field definitions to command
		for _, field := range fields {
			args = append(args, fmt.Sprintf("%s:%s", field.Name, field.Type))
		}
	}

	// Execute petrock command
	result, err := s.executePetrockCommand(args)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"suggestion": "Check that the feature and entity names are valid",
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Generated command %s/%s", featureName, entityName),
		"files_created": []string{
			fmt.Sprintf("%s/commands/%s.go", featureName, entityName),
		},
		"output": result,
		"next_steps": []string{
			"Run './build.sh' to compile the project",
			"Check the generated files for any needed customizations",
		},
	}, nil
}

// handleGenerateQuery processes the generate_query tool request
func (s *MCPServer) handleGenerateQuery(params map[string]interface{}) (interface{}, error) {
	// Validate project structure
	if err := s.validateProjectStructure(); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"suggestion": "Make sure you're in a petrock project root directory",
		}, nil
	}

	// Extract required parameters
	featureName, ok := params["feature_name"].(string)
	if !ok {
		return nil, fmt.Errorf("feature_name must be a string")
	}

	entityName, ok := params["entity_name"].(string)
	if !ok {
		return nil, fmt.Errorf("entity_name must be a string")
	}

	// Build command arguments
	args := []string{"new", "query", fmt.Sprintf("%s/%s", featureName, entityName)}

	// Note: Query generation with custom fields would need CLI support
	// For now, we'll generate basic queries and mention this limitation
	if fieldsParam, exists := params["fields"]; exists && fieldsParam != nil {
		if fields, ok := fieldsParam.([]interface{}); ok && len(fields) > 0 {
			return map[string]interface{}{
				"success": false,
				"error":   "Query generation with custom fields is not yet supported via CLI",
				"suggestion": "Generate the basic query first, then customize manually",
			}, nil
		}
	}

	// Execute petrock command
	result, err := s.executePetrockCommand(args)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"suggestion": "Check that the feature and entity names are valid",
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Generated query %s/%s", featureName, entityName),
		"files_created": []string{
			fmt.Sprintf("%s/queries/%s.go", featureName, entityName),
		},
		"output": result,
		"next_steps": []string{
			"Run './build.sh' to compile the project",
			"Check the generated files for any needed customizations",
		},
	}, nil
}

// handleGenerateWorker processes the generate_worker tool request
func (s *MCPServer) handleGenerateWorker(params map[string]interface{}) (interface{}, error) {
	// Validate project structure
	if err := s.validateProjectStructure(); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"suggestion": "Make sure you're in a petrock project root directory",
		}, nil
	}

	// Extract required parameters
	featureName, ok := params["feature_name"].(string)
	if !ok {
		return nil, fmt.Errorf("feature_name must be a string")
	}

	workerName, ok := params["worker_name"].(string)
	if !ok {
		return nil, fmt.Errorf("worker_name must be a string")
	}

	// Build command arguments
	args := []string{"new", "worker", fmt.Sprintf("%s/%s", featureName, workerName)}

	// Execute petrock command
	result, err := s.executePetrockCommand(args)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"suggestion": "Check that the feature and worker names are valid",
		}, nil
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Generated worker %s/%s", featureName, workerName),
		"files_created": []string{
			fmt.Sprintf("%s/workers/%s_worker.go", featureName, workerName),
		},
		"output": result,
		"next_steps": []string{
			"Run './build.sh' to compile the project",
			"Check the generated files for any needed customizations",
		},
	}, nil
}

// handleGenerateComponent processes the generate_component tool request
func (s *MCPServer) handleGenerateComponent(params map[string]interface{}) (interface{}, error) {
	// Extract component type
	componentType, ok := params["component_type"].(string)
	if !ok {
		return nil, fmt.Errorf("component_type must be a string")
	}

	// Route to specific handler based on component type
	switch componentType {
	case "command":
		return s.handleGenerateCommand(params)
	case "query":
		return s.handleGenerateQuery(params)
	case "worker":
		// Map entity_name to worker_name for worker handler
		if entityName, exists := params["entity_name"]; exists {
			params["worker_name"] = entityName
		}
		return s.handleGenerateWorker(params)
	default:
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Unknown component type: %s", componentType),
			"suggestion": "Use 'command', 'query', or 'worker'",
		}, nil
	}
}

// Helper methods

// Field represents a field definition
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// parseFields parses field definitions from MCP parameters
func (s *MCPServer) parseFields(fieldsParam interface{}) ([]Field, error) {
	fieldsArray, ok := fieldsParam.([]interface{})
	if !ok {
		return nil, fmt.Errorf("fields must be an array")
	}

	var fields []Field
	for i, fieldInterface := range fieldsArray {
		fieldMap, ok := fieldInterface.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field %d must be an object", i)
		}

		name, ok := fieldMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("field %d name must be a string", i)
		}

		fieldType, ok := fieldMap["type"].(string)
		if !ok {
			return nil, fmt.Errorf("field %d type must be a string", i)
		}

		if name == "" {
			return nil, fmt.Errorf("field %d name cannot be empty", i)
		}

		if fieldType == "" {
			return nil, fmt.Errorf("field %d type cannot be empty", i)
		}

		fields = append(fields, Field{Name: name, Type: fieldType})
	}

	return fields, nil
}

// executePetrockCommand executes a petrock CLI command
func (s *MCPServer) executePetrockCommand(args []string) (string, error) {
	// Find petrock binary
	petrockBinary, err := exec.LookPath("petrock")
	if err != nil {
		return "", fmt.Errorf("petrock command not found in PATH: %w", err)
	}

	// Create command
	cmd := exec.Command(petrockBinary, args...)
	cmd.Dir = s.projectRoot

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("petrock command failed: %w\nOutput: %s", err, string(output))
	}

	return string(output), nil
}

// StdioTransport handles communication over stdin/stdout
type StdioTransport struct {
	server *MCPServer
	reader *bufio.Scanner
}

// NewStdioTransport creates a new stdio transport
func NewStdioTransport(server *MCPServer) *StdioTransport {
	return &StdioTransport{
		server: server,
		reader: bufio.NewScanner(os.Stdin),
	}
}

// Run starts the stdio transport loop
func (t *StdioTransport) Run() error {
	for t.reader.Scan() {
		line := t.reader.Bytes()
		if len(line) == 0 {
			continue
		}

		responseBytes, err := t.server.jsonrpcServer.ProcessMessage(line)
		if err != nil {
			// Log error but continue processing
			fmt.Fprintf(os.Stderr, "Error processing message: %v\n", err)
			continue
		}

		// Only send response if there is one (requests get responses, notifications don't)
		if responseBytes != nil {
			fmt.Printf("%s\n", responseBytes)
		}
	}

	if err := t.reader.Err(); err != nil {
		return fmt.Errorf("error reading from stdin: %w", err)
	}

	return nil
}

// StartStdioServer runs the MCP server on stdio
func StartStdioServer(server *MCPServer) error {
	transport := NewStdioTransport(server)
	return transport.Run()
}
