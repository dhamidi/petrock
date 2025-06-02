# Phase 1: JSON-RPC 2.0 Foundation

## Overview

Implement a transport-agnostic JSON-RPC 2.0 server that can handle requests and notifications without being tied to specific I/O mechanisms. The server uses a notification pattern where incoming messages trigger registered handlers.

## Goals

1. Create a minimal JSON-RPC 2.0 implementation with no external dependencies
2. Support both request/response and notification patterns
3. Implement error handling according to JSON-RPC 2.0 specification
4. Design for easy testing with mock transports

## Core Components

### 1. `core/jsonrpc.go` - JSON-RPC 2.0 Types and Protocol

```go
// Request represents a JSON-RPC 2.0 request
type Request struct {
    JSONRpc string      `json:"jsonrpc"`
    ID      interface{} `json:"id,omitempty"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response
type Response struct {
    JSONRpc string      `json:"jsonrpc"`
    ID      interface{} `json:"id"`
    Result  interface{} `json:"result,omitempty"`
    Error   *Error      `json:"error,omitempty"`
}

// Notification represents a JSON-RPC 2.0 notification (no ID)
type Notification struct {
    JSONRpc string      `json:"jsonrpc"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params,omitempty"`
}

// Error represents a JSON-RPC 2.0 error
type Error struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 2. Server Interface

```go
// Server handles JSON-RPC 2.0 requests and notifications
type Server struct {
    requestHandlers      map[string]RequestHandler
    notificationHandlers map[string]NotificationHandler
}

// RequestHandler processes requests that expect responses
type RequestHandler func(params interface{}) (interface{}, error)

// NotificationHandler processes notifications (no response)
type NotificationHandler func(params interface{})

// RegisterRequestHandler registers a handler for a specific method
func (s *Server) RegisterRequestHandler(method string, handler RequestHandler)

// RegisterNotificationHandler registers a handler for a specific method
func (s *Server) RegisterNotificationHandler(method string, handler NotificationHandler)

// ReceiveNotification processes an incoming notification
func (s *Server) ReceiveNotification(notification Notification)

// HandleRequest processes an incoming request and returns a response
func (s *Server) HandleRequest(request Request) Response
```

### 3. Message Processing

```go
// ProcessMessage handles incoming JSON-RPC messages
func (s *Server) ProcessMessage(data []byte) ([]byte, error) {
    // Parse message to determine if it's a request or notification
    // Route to appropriate handler
    // Return response for requests, nil for notifications
}
```

## Implementation Details

### Request vs Notification Detection

- **Request**: Contains an `id` field (can be string, number, or null)
- **Notification**: No `id` field present

### Error Codes (JSON-RPC 2.0 Standard)

- `-32700`: Parse error
- `-32600`: Invalid Request  
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error
- `-32000` to `-32099`: Server error (reserved)

### Handler Registration

```go
server := NewJSONRPCServer()

// Register request handler that returns a response
server.RegisterRequestHandler("add", func(params interface{}) (interface{}, error) {
    // Parse params and return result
    return result, nil
})

// Register notification handler (no response)
server.RegisterNotificationHandler("log", func(params interface{}) {
    // Process notification, no return value
})
```

### Testing Strategy

1. **Unit tests** for message parsing and routing
2. **Handler tests** with mock request/notification data
3. **Error handling tests** for all error conditions
4. **Integration tests** with mock transport layer

## Files to Create

- `internal/skeleton/core/jsonrpc.go` - Core JSON-RPC implementation
- `internal/skeleton/core/jsonrpc_test.go` - Comprehensive unit tests

## Success Criteria

1. All JSON-RPC 2.0 message types properly parsed
2. Request/notification routing works correctly
3. Error responses follow JSON-RPC 2.0 specification
4. Handlers can be registered and invoked
5. No external dependencies (standard library only)
6. 100% test coverage for core functionality

## Next Phase

Phase 2 will build the MCP server on top of this JSON-RPC foundation, implementing MCP-specific message handlers and protocol negotiation.
