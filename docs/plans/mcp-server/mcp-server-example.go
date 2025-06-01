package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime"
)

// MCP Protocol Types

type MCPCapabilities struct {
	Resources *ResourceCapabilities `json:"resources,omitempty"`
	Tools     *ToolCapabilities     `json:"tools,omitempty"`
	Prompts   *PromptCapabilities   `json:"prompts,omitempty"`
	Logging   *LoggingCapabilities  `json:"logging,omitempty"`
}

type ResourceCapabilities struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type ToolCapabilities struct{}

type PromptCapabilities struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type LoggingCapabilities struct{}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Request/Response types

type InitializeRequest struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    MCPCapabilities `json:"capabilities"`
	ClientInfo      ServerInfo      `json:"clientInfo"`
}

type InitializeResponse struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    MCPCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo     `json:"serverInfo"`
	Instructions    string         `json:"instructions,omitempty"`
}

type Resource struct {
	URI         string      `json:"uri"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	MimeType    string      `json:"mimeType,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema interface{} `json:"inputSchema"`
}

type Prompt struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Arguments   []PromptArgument       `json:"arguments,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// MCP Server implementation

type MCPServer struct {
	app    *App
	server *JSONRPCServer
	logger *slog.Logger
}

func NewMCPServer(app *App, logger *slog.Logger) *MCPServer {
	server := NewJSONRPCServer(logger)
	mcp := &MCPServer{
		app:    app,
		server: server,
		logger: logger,
	}
	
	// Register MCP protocol handlers
	mcp.registerHandlers()
	
	return mcp
}

func (m *MCPServer) registerHandlers() {
	m.server.RegisterHandler("initialize", m.handleInitialize)
	m.server.RegisterHandler("ping", m.handlePing)
	m.server.RegisterHandler("resources/list", m.handleResourcesList)
	m.server.RegisterHandler("resources/read", m.handleResourcesRead)
	m.server.RegisterHandler("tools/list", m.handleToolsList)
	m.server.RegisterHandler("tools/call", m.handleToolsCall)
	m.server.RegisterHandler("prompts/list", m.handlePromptsList)
	m.server.RegisterHandler("prompts/get", m.handlePromptsGet)
}

func (m *MCPServer) Serve(ctx context.Context) error {
	return m.server.Serve(ctx)
}

// Protocol handlers

func (m *MCPServer) handleInitialize(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req InitializeRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, NewJSONRPCError(InvalidParams, "Invalid initialize parameters", nil)
	}
	
	// Validate protocol version
	if req.ProtocolVersion != "2024-11-05" && req.ProtocolVersion != "2025-03-26" {
		return nil, NewJSONRPCError(InvalidRequest, "Unsupported protocol version", nil)
	}
	
	response := InitializeResponse{
		ProtocolVersion: "2024-11-05", // Use stable version
		Capabilities: MCPCapabilities{
			Resources: &ResourceCapabilities{
				Subscribe:   false,
				ListChanged: false,
			},
			Tools: &ToolCapabilities{},
			Prompts: &PromptCapabilities{
				ListChanged: false,
			},
			Logging: &LoggingCapabilities{},
		},
		ServerInfo: ServerInfo{
			Name:    "petrock_example_project_name", // Will be replaced by template
			Version: "1.0.0",
		},
		Instructions: "This server provides access to the petrock application data and tools.",
	}
	
	return response, nil
}

func (m *MCPServer) handlePing(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{}, nil
}

func (m *MCPServer) handleResourcesList(ctx context.Context, params json.RawMessage) (interface{}, error) {
	resources := []Resource{
		{
			URI:         "app://info",
			Name:        "Application Info",
			Description: "Basic information about this application",
			MimeType:    "application/json",
		},
		{
			URI:         "app://features",
			Name:        "Features List",
			Description: "List of registered application features",
			MimeType:    "application/json",
		},
		{
			URI:         "app://stats",
			Name:        "Runtime Statistics",
			Description: "Application runtime statistics and health",
			MimeType:    "application/json",
		},
	}
	
	return map[string]interface{}{
		"resources": resources,
	}, nil
}

func (m *MCPServer) handleResourcesRead(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, NewJSONRPCError(InvalidParams, "Invalid read parameters", nil)
	}
	
	switch req.URI {
	case "app://info":
		return map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"uri":      "app://info",
					"mimeType": "application/json",
					"text": fmt.Sprintf(`{
  "name": "petrock_example_project_name",
  "version": "1.0.0",
  "runtime": "%s",
  "architecture": "%s"
}`, runtime.Version(), runtime.GOARCH),
				},
			},
		}, nil
		
	case "app://features":
		features := make([]string, 0, len(m.app.Commands.handlers))
		for name := range m.app.Commands.handlers {
			features = append(features, name)
		}
		
		featuresJSON, _ := json.MarshalIndent(map[string]interface{}{
			"commands": features,
			"count":    len(features),
		}, "", "  ")
		
		return map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"uri":      "app://features",
					"mimeType": "application/json",
					"text":     string(featuresJSON),
				},
			},
		}, nil
		
	case "app://stats":
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		statsJSON, _ := json.MarshalIndent(map[string]interface{}{
			"goroutines":    runtime.NumGoroutine(),
			"memory_alloc": m.Alloc,
			"memory_sys":   m.Sys,
			"gc_cycles":    m.NumGC,
		}, "", "  ")
		
		return map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"uri":      "app://stats",
					"mimeType": "application/json",
					"text":     string(statsJSON),
				},
			},
		}, nil
		
	default:
		return nil, NewJSONRPCError(InvalidParams, "Resource not found", nil)
	}
}

func (m *MCPServer) handleToolsList(ctx context.Context, params json.RawMessage) (interface{}, error) {
	tools := []Tool{
		{
			Name:        "query_database",
			Description: "Execute a read-only SQL query against the application database",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "SQL SELECT query to execute",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "get_kv",
			Description: "Retrieve a value from the key-value store",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"key": map[string]interface{}{
						"type":        "string",
						"description": "Key to retrieve",
					},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "list_kv_keys",
			Description: "List all keys in the key-value store",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
	
	return map[string]interface{}{
		"tools": tools,
	}, nil
}

func (m *MCPServer) handleToolsCall(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, NewJSONRPCError(InvalidParams, "Invalid tool call parameters", nil)
	}
	
	switch req.Name {
	case "query_database":
		return m.handleQueryDatabase(ctx, req.Arguments)
	case "get_kv":
		return m.handleGetKV(ctx, req.Arguments)
	case "list_kv_keys":
		return m.handleListKVKeys(ctx, req.Arguments)
	default:
		return nil, NewJSONRPCError(MethodNotFound, "Tool not found", nil)
	}
}

func (m *MCPServer) handleQueryDatabase(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, NewJSONRPCError(InvalidParams, "Invalid query parameters", nil)
	}
	
	// Simple SQL injection protection - only allow SELECT
	if len(params.Query) < 6 || params.Query[:6] != "SELECT" {
		return nil, NewJSONRPCError(InvalidParams, "Only SELECT queries are allowed", nil)
	}
	
	rows, err := m.app.DB.QueryContext(ctx, params.Query)
	if err != nil {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Query failed: %v", err),
				},
			},
			"isError": true,
		}, nil
	}
	defer rows.Close()
	
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, NewJSONRPCError(InternalError, "Failed to get columns", nil)
	}
	
	// Read all rows
	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, NewJSONRPCError(InternalError, "Failed to scan row", nil)
		}
		
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}
	
	resultJSON, _ := json.MarshalIndent(results, "", "  ")
	
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": string(resultJSON),
			},
		},
	}, nil
}

func (m *MCPServer) handleGetKV(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, NewJSONRPCError(InvalidParams, "Invalid key parameters", nil)
	}
	
	value, err := m.app.KV.Get(ctx, params.Key)
	if err != nil {
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Key not found or error: %v", err),
				},
			},
			"isError": true,
		}, nil
	}
	
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": value,
			},
		},
	}, nil
}

func (m *MCPServer) handleListKVKeys(ctx context.Context, args json.RawMessage) (interface{}, error) {
	// This would need to be implemented based on your KV store
	// For now, return a placeholder
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "KV key listing not yet implemented",
			},
		},
	}, nil
}

func (m *MCPServer) handlePromptsList(ctx context.Context, params json.RawMessage) (interface{}, error) {
	prompts := []Prompt{
		{
			Name:        "analyze_app",
			Description: "Analyze the application structure and provide insights",
			Arguments: []PromptArgument{
				{
					Name:        "focus",
					Description: "What aspect to focus on (features, performance, security)",
					Required:    false,
				},
			},
		},
		{
			Name:        "debug_feature",
			Description: "Debug a specific application feature",
			Arguments: []PromptArgument{
				{
					Name:        "feature_name",
					Description: "Name of the feature to debug",
					Required:    true,
				},
			},
		},
	}
	
	return map[string]interface{}{
		"prompts": prompts,
	}, nil
}

func (m *MCPServer) handlePromptsGet(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments,omitempty"`
	}
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, NewJSONRPCError(InvalidParams, "Invalid prompt parameters", nil)
	}
	
	switch req.Name {
	case "analyze_app":
		focus := "general"
		if f, ok := req.Arguments["focus"].(string); ok {
			focus = f
		}
		
		return map[string]interface{}{
			"description": "Analyze this petrock application",
			"messages": []map[string]interface{}{
				{
					"role": "user",
					"content": map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf("Please analyze this petrock application with focus on: %s\n\nUse the available MCP tools to explore the application structure, database, and features.", focus),
					},
				},
			},
		}, nil
		
	case "debug_feature":
		featureName, ok := req.Arguments["feature_name"].(string)
		if !ok || featureName == "" {
			return nil, NewJSONRPCError(InvalidParams, "feature_name is required", nil)
		}
		
		return map[string]interface{}{
			"description": fmt.Sprintf("Debug the %s feature", featureName),
			"messages": []map[string]interface{}{
				{
					"role": "user",
					"content": map[string]interface{}{
						"type": "text",
						"text": fmt.Sprintf("Help me debug the '%s' feature in this petrock application.\n\nPlease use the available MCP tools to:\n1. Check if the feature is registered\n2. Query related database tables\n3. Examine the application state\n4. Suggest debugging steps", featureName),
					},
				},
			},
		}, nil
		
	default:
		return nil, NewJSONRPCError(MethodNotFound, "Prompt not found", nil)
	}
}
