package core

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

// JSON-RPC 2.0 Types

type JSONRPCRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// Handler function type
type Handler func(ctx context.Context, params json.RawMessage) (interface{}, error)

// JSONRPCServer handles JSON-RPC 2.0 protocol over stdio
type JSONRPCServer struct {
	handlers map[string]Handler
	mu       sync.RWMutex
	logger   *slog.Logger
}

func NewJSONRPCServer(logger *slog.Logger) *JSONRPCServer {
	return &JSONRPCServer{
		handlers: make(map[string]Handler),
		logger:   logger,
	}
}

// RegisterHandler registers a method handler
func (s *JSONRPCServer) RegisterHandler(method string, handler Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[method] = handler
}

// Serve starts the JSON-RPC server on stdio
func (s *JSONRPCServer) Serve(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		line := scanner.Text()
		if line == "" {
			continue
		}
		
		s.handleRequest(ctx, line)
	}
	
	return scanner.Err()
}

func (s *JSONRPCServer) handleRequest(ctx context.Context, line string) {
	var req JSONRPCRequest
	if err := json.Unmarshal([]byte(line), &req); err != nil {
		s.sendError(nil, ParseError, "Parse error", nil)
		return
	}
	
	// Validate JSON-RPC version
	if req.Jsonrpc != "2.0" {
		s.sendError(req.ID, InvalidRequest, "Invalid Request", nil)
		return
	}
	
	// Look up handler
	s.mu.RLock()
	handler, exists := s.handlers[req.Method]
	s.mu.RUnlock()
	
	if !exists {
		// Only send error response for requests (not notifications)
		if req.ID != nil {
			s.sendError(req.ID, MethodNotFound, "Method not found", nil)
		}
		return
	}
	
	// Execute handler
	result, err := handler(ctx, req.Params)
	
	// Only send response for requests (not notifications)
	if req.ID != nil {
		if err != nil {
			var code int
			var message string
			
			// Try to extract JSON-RPC error details
			if rpcErr, ok := err.(*JSONRPCError); ok {
				code = rpcErr.Code
				message = rpcErr.Message
			} else {
				code = InternalError
				message = err.Error()
			}
			
			s.sendError(req.ID, code, message, nil)
		} else {
			s.sendResult(req.ID, result)
		}
	}
}

func (s *JSONRPCServer) sendResult(id *json.RawMessage, result interface{}) {
	response := JSONRPCResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
	}
	
	s.writeResponse(response)
}

func (s *JSONRPCServer) sendError(id *json.RawMessage, code int, message string, data interface{}) {
	response := JSONRPCResponse{
		Jsonrpc: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	
	s.writeResponse(response)
}

func (s *JSONRPCServer) writeResponse(response JSONRPCResponse) {
	data, err := json.Marshal(response)
	if err != nil {
		s.logger.Error("Failed to marshal response", "error", err)
		return
	}
	
	fmt.Println(string(data))
}

// SendNotification sends a notification to the client
func (s *JSONRPCServer) SendNotification(method string, params interface{}) error {
	notification := struct {
		Jsonrpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
	}{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
	
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	
	fmt.Println(string(data))
	return nil
}

// Helper to create JSON-RPC errors
func NewJSONRPCError(code int, message string, data interface{}) error {
	return &JSONRPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
