package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSON-RPC 2.0 standard error codes
const (
	JSONRPCParseError     = -32700
	JSONRPCInvalidRequest = -32600
	JSONRPCMethodNotFound = -32601
	JSONRPCInvalidParams  = -32602
	JSONRPCInternalError  = -32603
)

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

// RequestHandler processes requests that expect responses
type RequestHandler func(params interface{}) (interface{}, error)

// NotificationHandler processes notifications (no response)
type NotificationHandler func(params interface{})

// Server handles JSON-RPC 2.0 requests and notifications
type Server struct {
	requestHandlers      map[string]RequestHandler
	notificationHandlers map[string]NotificationHandler
}

// NewJSONRPCServer creates a new JSON-RPC server
func NewJSONRPCServer() *Server {
	return &Server{
		requestHandlers:      make(map[string]RequestHandler),
		notificationHandlers: make(map[string]NotificationHandler),
	}
}

// RegisterRequestHandler registers a handler for a specific method
func (s *Server) RegisterRequestHandler(method string, handler RequestHandler) {
	s.requestHandlers[method] = handler
}

// RegisterNotificationHandler registers a handler for a specific method
func (s *Server) RegisterNotificationHandler(method string, handler NotificationHandler) {
	s.notificationHandlers[method] = handler
}

// ReceiveNotification processes an incoming notification
func (s *Server) ReceiveNotification(notification Notification) {
	if handler, exists := s.notificationHandlers[notification.Method]; exists {
		handler(notification.Params)
	}
	// Notifications are silently ignored if no handler exists
}

// HandleRequest processes an incoming request and returns a response
func (s *Server) HandleRequest(request Request) Response {
	response := Response{
		JSONRpc: "2.0",
		ID:      request.ID,
	}

	handler, exists := s.requestHandlers[request.Method]
	if !exists {
		response.Error = &Error{
			Code:    JSONRPCMethodNotFound,
			Message: "Method not found",
			Data:    request.Method,
		}
		return response
	}

	result, err := handler(request.Params)
	if err != nil {
		response.Error = &Error{
			Code:    JSONRPCInternalError,
			Message: err.Error(),
		}
		return response
	}

	response.Result = result
	return response
}

// ProcessMessage handles incoming JSON-RPC messages
func (s *Server) ProcessMessage(data []byte) ([]byte, error) {
	// Try to parse as a generic message first to determine type
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		// If parse fails and the JSON doesn't contain "id" field, treat as malformed notification
		// and silently ignore per JSON-RPC spec
		dataStr := string(data)
		if !strings.Contains(dataStr, `"id"`) {
			return nil, nil // Silently ignore malformed notifications
		}
		
		errorResponse := Response{
			JSONRpc: "2.0",
			ID:      nil,
			Error: &Error{
				Code:    JSONRPCParseError,
				Message: "Parse error",
				Data:    err.Error(),
			},
		}
		return json.Marshal(errorResponse)
	}

	// Check for required jsonrpc field
	if jsonrpc, ok := msg["jsonrpc"].(string); !ok || jsonrpc != "2.0" {
		errorResponse := Response{
			JSONRpc: "2.0",
			ID:      msg["id"],
			Error: &Error{
				Code:    JSONRPCInvalidRequest,
				Message: "Invalid Request",
				Data:    "Missing or invalid jsonrpc field",
			},
		}
		return json.Marshal(errorResponse)
	}

	// Check for required method field
	_, ok := msg["method"].(string)
	if !ok {
		errorResponse := Response{
			JSONRpc: "2.0",
			ID:      msg["id"],
			Error: &Error{
				Code:    JSONRPCInvalidRequest,
				Message: "Invalid Request",
				Data:    "Missing method field",
			},
		}
		return json.Marshal(errorResponse)
	}

	// Determine if this is a request or notification based on presence of ID
	if _, hasID := msg["id"]; hasID {
		// This is a request
		var request Request
		if err := json.Unmarshal(data, &request); err != nil {
			errorResponse := Response{
				JSONRpc: "2.0",
				ID:      msg["id"],
				Error: &Error{
					Code:    JSONRPCParseError,
					Message: "Parse error",
					Data:    err.Error(),
				},
			}
			return json.Marshal(errorResponse)
		}

		response := s.HandleRequest(request)
		return json.Marshal(response)
	} else {
		// This is a notification
		var notification Notification
		if err := json.Unmarshal(data, &notification); err != nil {
			// Notifications with parse errors are silently ignored per JSON-RPC spec
			return nil, nil
		}

		s.ReceiveNotification(notification)
		return nil, nil // No response for notifications
	}
}

// NewError creates a new JSON-RPC error with the given code and message
func NewError(code int, message string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// CreateErrorResponse creates a JSON-RPC error response
func CreateErrorResponse(id interface{}, code int, message string, data interface{}) Response {
	return Response{
		JSONRpc: "2.0",
		ID:      id,
		Error:   NewError(code, message, data),
	}
}

// CreateSuccessResponse creates a JSON-RPC success response
func CreateSuccessResponse(id interface{}, result interface{}) Response {
	return Response{
		JSONRpc: "2.0",
		ID:      id,
		Result:  result,
	}
}

// IsValidRequest checks if a request has all required fields
func (r *Request) IsValidRequest() error {
	if r.JSONRpc != "2.0" {
		return fmt.Errorf("invalid jsonrpc version: %s", r.JSONRpc)
	}
	if r.Method == "" {
		return fmt.Errorf("missing method field")
	}
	return nil
}

// IsValidNotification checks if a notification has all required fields
func (n *Notification) IsValidNotification() error {
	if n.JSONRpc != "2.0" {
		return fmt.Errorf("invalid jsonrpc version: %s", n.JSONRpc)
	}
	if n.Method == "" {
		return fmt.Errorf("missing method field")
	}
	return nil
}
