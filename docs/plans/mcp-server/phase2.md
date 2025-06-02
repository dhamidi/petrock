# Phase 2: Basic MCP Server Implementation

## Overview

Build the Model Context Protocol (MCP) server on top of the JSON-RPC 2.0 foundation. This phase implements core MCP protocol messages and provides a demo "addition" tool to validate the complete flow.

## Goals

1. Implement MCP initialization and capability negotiation
2. Support basic MCP protocol messages (initialize, ping, tools/list, tools/call)
3. Provide a working demo tool for end-to-end testing
4. Add stdio transport for communication with MCP clients

## Core Components

### 1. `core/mcp.go` - MCP Server Implementation

```go
// MCPServer implements the Model Context Protocol
type MCPServer struct {
    jsonrpcServer *JSONRPCServer
    capabilities  ServerCapabilities
    tools         map[string]Tool
}

// ServerCapabilities describes what the server supports
type ServerCapabilities struct {
    Tools       *ToolsCapability       `json:"tools,omitempty"`
    Resources   *ResourcesCapability   `json:"resources,omitempty"`
    Prompts     *PromptsCapability     `json:"prompts,omitempty"`
    Logging     *LoggingCapability     `json:"logging,omitempty"`
}

// Tool represents an MCP tool definition
type Tool struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    InputSchema ToolInputSchema        `json:"inputSchema"`
    Handler     func(params map[string]interface{}) (interface{}, error)
}

// ToolInputSchema defines the expected parameters for a tool
type ToolInputSchema struct {
    Type       string                 `json:"type"`
    Properties map[string]interface{} `json:"properties"`
    Required   []string               `json:"required,omitempty"`
}
```

### 2. MCP Protocol Messages

#### Initialize Request/Response
```go
type InitializeRequest struct {
    ProtocolVersion string            `json:"protocolVersion"`
    Capabilities    ClientCapabilities `json:"capabilities"`
    ClientInfo      ClientInfo        `json:"clientInfo"`
}

type InitializeResponse struct {
    ProtocolVersion string            `json:"protocolVersion"`
    Capabilities    ServerCapabilities `json:"capabilities"`
    ServerInfo      ServerInfo        `json:"serverInfo"`
}
```

#### Tools Messages
```go
type ListToolsRequest struct {
    Cursor string `json:"cursor,omitempty"`
}

type ListToolsResponse struct {
    Tools      []Tool `json:"tools"`
    NextCursor string `json:"nextCursor,omitempty"`
}

type CallToolRequest struct {
    Name      string                 `json:"name"`
    Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type CallToolResponse struct {
    Content []ToolResponseContent `json:"content"`
    IsError bool                  `json:"isError,omitempty"`
}
```

### 3. Demo Addition Tool

```go
// RegisterDemoTools adds the addition tool for testing
func (s *MCPServer) RegisterDemoTools() {
    addTool := Tool{
        Name:        "add",
        Description: "Add two numbers together",
        InputSchema: ToolInputSchema{
            Type: "object",
            Properties: map[string]interface{}{
                "a": map[string]interface{}{
                    "type": "number",
                    "description": "First number",
                },
                "b": map[string]interface{}{
                    "type": "number", 
                    "description": "Second number",
                },
            },
            Required: []string{"a", "b"},
        },
        Handler: func(params map[string]interface{}) (interface{}, error) {
            a, ok := params["a"].(float64)
            if !ok {
                return nil, fmt.Errorf("parameter 'a' must be a number")
            }
            b, ok := params["b"].(float64)
            if !ok {
                return nil, fmt.Errorf("parameter 'b' must be a number")
            }
            return map[string]interface{}{
                "result": a + b,
            }, nil
        },
    }
    s.RegisterTool(addTool)
}
```

### 4. Stdio Transport

```go
// StdioTransport handles communication over stdin/stdout
type StdioTransport struct {
    server *MCPServer
    reader *bufio.Scanner
    writer *bufio.Writer
}

// StartStdioServer runs the MCP server on stdio
func StartStdioServer(server *MCPServer) error {
    transport := &StdioTransport{
        server: server,
        reader: bufio.NewScanner(os.Stdin),
        writer: bufio.NewWriter(os.Stdout),
    }
    
    return transport.Run()
}
```

## MCP Protocol Flow

### 1. Initialization
```
Client -> Server: initialize request
Server -> Client: initialize response (capabilities)
Client -> Server: initialized notification  
```

### 2. Tool Discovery  
```
Client -> Server: tools/list request
Server -> Client: tools/list response (available tools)
```

### 3. Tool Execution
```
Client -> Server: tools/call request (tool name + arguments)  
Server -> Client: tools/call response (result or error)
```

## Implementation Details

### MCP Message Handlers

```go
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
```

### Error Handling

- MCP-specific error codes for tool execution failures
- Proper JSON-RPC error responses for protocol violations
- Graceful handling of malformed requests

### Testing Strategy

1. **Unit tests** for individual MCP handlers
2. **Integration tests** with mock stdio transport
3. **End-to-end tests** with actual MCP clients
4. **Tool execution tests** for the demo addition tool

## Files to Create

- `internal/skeleton/core/mcp.go` - Core MCP server implementation
- `internal/skeleton/core/mcp_test.go` - Unit tests
- `internal/skeleton/cmd/petrock_example_project_name/mcp.go` - MCP command

## Success Criteria

1. MCP server initializes and negotiates capabilities correctly
2. Demo addition tool can be discovered via tools/list
3. Demo addition tool executes correctly via tools/call
4. Server responds to ping requests
5. Stdio transport works with real MCP clients
6. All MCP protocol messages properly formatted

## Next Phase

Phase 3 will replace the demo addition tool with actual petrock generators (command, query, worker generation) to provide real development utility.
