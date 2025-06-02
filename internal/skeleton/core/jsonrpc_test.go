package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewJSONRPCServer(t *testing.T) {
	server := NewJSONRPCServer()
	if server == nil {
		t.Fatal("NewJSONRPCServer returned nil")
	}
	if server.requestHandlers == nil {
		t.Fatal("requestHandlers map not initialized")
	}
	if server.notificationHandlers == nil {
		t.Fatal("notificationHandlers map not initialized")
	}
}

func TestRegisterRequestHandler(t *testing.T) {
	server := NewJSONRPCServer()
	called := false
	
	handler := func(params interface{}) (interface{}, error) {
		called = true
		return "test result", nil
	}
	
	server.RegisterRequestHandler("test_method", handler)
	
	// Verify handler was registered
	if _, exists := server.requestHandlers["test_method"]; !exists {
		t.Fatal("Handler was not registered")
	}
	
	// Test that handler can be called
	request := Request{
		JSONRpc: "2.0",
		ID:      1,
		Method:  "test_method",
	}
	
	response := server.HandleRequest(request)
	if !called {
		t.Fatal("Handler was not called")
	}
	if response.Result != "test result" {
		t.Fatalf("Expected 'test result', got %v", response.Result)
	}
}

func TestRegisterNotificationHandler(t *testing.T) {
	server := NewJSONRPCServer()
	called := false
	var receivedParams interface{}
	
	handler := func(params interface{}) {
		called = true
		receivedParams = params
	}
	
	server.RegisterNotificationHandler("test_notification", handler)
	
	// Verify handler was registered
	if _, exists := server.notificationHandlers["test_notification"]; !exists {
		t.Fatal("Notification handler was not registered")
	}
	
	// Test that handler can be called
	notification := Notification{
		JSONRpc: "2.0",
		Method:  "test_notification",
		Params:  "test params",
	}
	
	server.ReceiveNotification(notification)
	if !called {
		t.Fatal("Notification handler was not called")
	}
	if receivedParams != "test params" {
		t.Fatalf("Expected 'test params', got %v", receivedParams)
	}
}

func TestHandleRequest_Success(t *testing.T) {
	server := NewJSONRPCServer()
	
	server.RegisterRequestHandler("add", func(params interface{}) (interface{}, error) {
		paramMap := params.(map[string]interface{})
		a := paramMap["a"].(float64)
		b := paramMap["b"].(float64)
		return a + b, nil
	})
	
	request := Request{
		JSONRpc: "2.0",
		ID:      1,
		Method:  "add",
		Params:  map[string]interface{}{"a": 5.0, "b": 3.0},
	}
	
	response := server.HandleRequest(request)
	
	if response.JSONRpc != "2.0" {
		t.Fatalf("Expected jsonrpc '2.0', got %s", response.JSONRpc)
	}
	if response.ID != 1 {
		t.Fatalf("Expected ID 1, got %v", response.ID)
	}
	if response.Error != nil {
		t.Fatalf("Expected no error, got %v", response.Error)
	}
	if response.Result != 8.0 {
		t.Fatalf("Expected result 8.0, got %v", response.Result)
	}
}

func TestHandleRequest_MethodNotFound(t *testing.T) {
	server := NewJSONRPCServer()
	
	request := Request{
		JSONRpc: "2.0",
		ID:      1,
		Method:  "nonexistent_method",
	}
	
	response := server.HandleRequest(request)
	
	if response.Error == nil {
		t.Fatal("Expected error for nonexistent method")
	}
	if response.Error.Code != JSONRPCMethodNotFound {
		t.Fatalf("Expected error code %d, got %d", JSONRPCMethodNotFound, response.Error.Code)
	}
	if response.Error.Message != "Method not found" {
		t.Fatalf("Expected 'Method not found', got %s", response.Error.Message)
	}
	if response.Error.Data != "nonexistent_method" {
		t.Fatalf("Expected method name in error data, got %v", response.Error.Data)
	}
}

func TestHandleRequest_HandlerError(t *testing.T) {
	server := NewJSONRPCServer()
	
	server.RegisterRequestHandler("error_method", func(params interface{}) (interface{}, error) {
		return nil, fmt.Errorf("test error")
	})
	
	request := Request{
		JSONRpc: "2.0",
		ID:      1,
		Method:  "error_method",
	}
	
	response := server.HandleRequest(request)
	
	if response.Error == nil {
		t.Fatal("Expected error from handler")
	}
	if response.Error.Code != JSONRPCInternalError {
		t.Fatalf("Expected error code %d, got %d", JSONRPCInternalError, response.Error.Code)
	}
}

func TestReceiveNotification_Success(t *testing.T) {
	server := NewJSONRPCServer()
	called := false
	
	server.RegisterNotificationHandler("log", func(params interface{}) {
		called = true
	})
	
	notification := Notification{
		JSONRpc: "2.0",
		Method:  "log",
		Params:  "test message",
	}
	
	server.ReceiveNotification(notification)
	
	if !called {
		t.Fatal("Notification handler was not called")
	}
}

func TestReceiveNotification_NoHandler(t *testing.T) {
	server := NewJSONRPCServer()
	
	notification := Notification{
		JSONRpc: "2.0",
		Method:  "nonexistent_notification",
	}
	
	// Should not panic when no handler exists
	server.ReceiveNotification(notification)
}

func TestProcessMessage_ValidRequest(t *testing.T) {
	server := NewJSONRPCServer()
	
	server.RegisterRequestHandler("echo", func(params interface{}) (interface{}, error) {
		return params, nil
	})
	
	requestJSON := `{"jsonrpc":"2.0","id":1,"method":"echo","params":"hello"}`
	
	responseData, err := server.ProcessMessage([]byte(requestJSON))
	if err != nil {
		t.Fatalf("ProcessMessage returned error: %v", err)
	}
	
	var response Response
	if err := json.Unmarshal(responseData, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if response.JSONRpc != "2.0" {
		t.Fatalf("Expected jsonrpc '2.0', got %s", response.JSONRpc)
	}
	if response.ID != float64(1) { // JSON unmarshals numbers as float64
		t.Fatalf("Expected ID 1, got %v", response.ID)
	}
	if response.Result != "hello" {
		t.Fatalf("Expected result 'hello', got %v", response.Result)
	}
}

func TestProcessMessage_ValidNotification(t *testing.T) {
	server := NewJSONRPCServer()
	called := false
	
	server.RegisterNotificationHandler("notify", func(params interface{}) {
		called = true
	})
	
	notificationJSON := `{"jsonrpc":"2.0","method":"notify","params":"test"}`
	
	responseData, err := server.ProcessMessage([]byte(notificationJSON))
	if err != nil {
		t.Fatalf("ProcessMessage returned error: %v", err)
	}
	
	// Notifications should not return a response
	if responseData != nil {
		t.Fatalf("Expected no response for notification, got: %s", responseData)
	}
	
	if !called {
		t.Fatal("Notification handler was not called")
	}
}

func TestProcessMessage_ParseError(t *testing.T) {
	server := NewJSONRPCServer()
	
	invalidJSON := `{"jsonrpc":"2.0","id":1,"method":"test"`
	
	responseData, err := server.ProcessMessage([]byte(invalidJSON))
	if err != nil {
		t.Fatalf("ProcessMessage returned error: %v", err)
	}
	
	var response Response
	if err := json.Unmarshal(responseData, &response); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}
	
	if response.Error == nil {
		t.Fatal("Expected parse error")
	}
	if response.Error.Code != JSONRPCParseError {
		t.Fatalf("Expected error code %d, got %d", JSONRPCParseError, response.Error.Code)
	}
}

func TestProcessMessage_InvalidRequest_MissingJSONRPC(t *testing.T) {
	server := NewJSONRPCServer()
	
	invalidJSON := `{"id":1,"method":"test"}`
	
	responseData, err := server.ProcessMessage([]byte(invalidJSON))
	if err != nil {
		t.Fatalf("ProcessMessage returned error: %v", err)
	}
	
	var response Response
	if err := json.Unmarshal(responseData, &response); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}
	
	if response.Error == nil {
		t.Fatal("Expected invalid request error")
	}
	if response.Error.Code != JSONRPCInvalidRequest {
		t.Fatalf("Expected error code %d, got %d", JSONRPCInvalidRequest, response.Error.Code)
	}
}

func TestProcessMessage_InvalidRequest_MissingMethod(t *testing.T) {
	server := NewJSONRPCServer()
	
	invalidJSON := `{"jsonrpc":"2.0","id":1}`
	
	responseData, err := server.ProcessMessage([]byte(invalidJSON))
	if err != nil {
		t.Fatalf("ProcessMessage returned error: %v", err)
	}
	
	var response Response
	if err := json.Unmarshal(responseData, &response); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}
	
	if response.Error == nil {
		t.Fatal("Expected invalid request error")
	}
	if response.Error.Code != JSONRPCInvalidRequest {
		t.Fatalf("Expected error code %d, got %d", JSONRPCInvalidRequest, response.Error.Code)
	}
}

func TestNewError(t *testing.T) {
	err := NewError(JSONRPCInvalidParams, "Invalid parameters", map[string]string{"field": "value"})
	
	if err.Code != JSONRPCInvalidParams {
		t.Fatalf("Expected code %d, got %d", JSONRPCInvalidParams, err.Code)
	}
	if err.Message != "Invalid parameters" {
		t.Fatalf("Expected message 'Invalid parameters', got %s", err.Message)
	}
	if err.Data == nil {
		t.Fatal("Expected data to be set")
	}
}

func TestCreateErrorResponse(t *testing.T) {
	response := CreateErrorResponse(123, JSONRPCInvalidParams, "Test error", "test data")
	
	if response.JSONRpc != "2.0" {
		t.Fatalf("Expected jsonrpc '2.0', got %s", response.JSONRpc)
	}
	if response.ID != 123 {
		t.Fatalf("Expected ID 123, got %v", response.ID)
	}
	if response.Error == nil {
		t.Fatal("Expected error to be set")
	}
	if response.Error.Code != JSONRPCInvalidParams {
		t.Fatalf("Expected error code %d, got %d", JSONRPCInvalidParams, response.Error.Code)
	}
	if response.Result != nil {
		t.Fatal("Expected result to be nil in error response")
	}
}

func TestCreateSuccessResponse(t *testing.T) {
	result := map[string]string{"status": "success"}
	response := CreateSuccessResponse(456, result)
	
	if response.JSONRpc != "2.0" {
		t.Fatalf("Expected jsonrpc '2.0', got %s", response.JSONRpc)
	}
	if response.ID != 456 {
		t.Fatalf("Expected ID 456, got %v", response.ID)
	}
	if response.Error != nil {
		t.Fatal("Expected error to be nil in success response")
	}
	if response.Result == nil {
		t.Fatal("Expected result to be set")
	}
}

func TestRequest_IsValidRequest(t *testing.T) {
	// Valid request
	validRequest := Request{
		JSONRpc: "2.0",
		Method:  "test",
		ID:      1,
	}
	if err := validRequest.IsValidRequest(); err != nil {
		t.Fatalf("Valid request failed validation: %v", err)
	}
	
	// Invalid jsonrpc version
	invalidVersionRequest := Request{
		JSONRpc: "1.0",
		Method:  "test",
		ID:      1,
	}
	if err := invalidVersionRequest.IsValidRequest(); err == nil {
		t.Fatal("Expected error for invalid jsonrpc version")
	}
	
	// Missing method
	missingMethodRequest := Request{
		JSONRpc: "2.0",
		ID:      1,
	}
	if err := missingMethodRequest.IsValidRequest(); err == nil {
		t.Fatal("Expected error for missing method")
	}
}

func TestNotification_IsValidNotification(t *testing.T) {
	// Valid notification
	validNotification := Notification{
		JSONRpc: "2.0",
		Method:  "test",
	}
	if err := validNotification.IsValidNotification(); err != nil {
		t.Fatalf("Valid notification failed validation: %v", err)
	}
	
	// Invalid jsonrpc version
	invalidVersionNotification := Notification{
		JSONRpc: "1.0",
		Method:  "test",
	}
	if err := invalidVersionNotification.IsValidNotification(); err == nil {
		t.Fatal("Expected error for invalid jsonrpc version")
	}
	
	// Missing method
	missingMethodNotification := Notification{
		JSONRpc: "2.0",
	}
	if err := missingMethodNotification.IsValidNotification(); err == nil {
		t.Fatal("Expected error for missing method")
	}
}

func TestProcessMessage_NotificationParseError_SilentlyIgnored(t *testing.T) {
	server := NewJSONRPCServer()
	
	// Invalid JSON that looks like a notification (no ID)
	invalidNotificationJSON := `{"jsonrpc":"2.0","method":"test"`
	
	responseData, err := server.ProcessMessage([]byte(invalidNotificationJSON))
	
	// Should not return an error or response for malformed notifications
	if err != nil {
		t.Fatalf("ProcessMessage should not return error for malformed notification: %v", err)
	}
	if responseData != nil {
		t.Fatalf("ProcessMessage should not return response for malformed notification: %s", responseData)
	}
}
