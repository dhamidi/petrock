package core

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewMCPServer(t *testing.T) {
	server := NewMCPServer()

	if server == nil {
		t.Fatal("NewMCPServer() returned nil")
	}

	if server.jsonrpcServer == nil {
		t.Fatal("JSON-RPC server not initialized")
	}

	if server.tools == nil {
		t.Fatal("Tools map not initialized")
	}

	if server.initialized {
		t.Fatal("Server should not be initialized at creation")
	}

	// Check that demo tool is registered
	if _, exists := server.tools["add"]; !exists {
		t.Fatal("Demo 'add' tool not registered")
	}
}

func TestMCPServer_HandleInitialize(t *testing.T) {
	server := NewMCPServer()

	// Test valid initialize request
	params := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "test-client",
			"version": "1.0.0",
		},
	}

	result, err := server.handleInitialize(params)
	if err != nil {
		t.Fatalf("handleInitialize() failed: %v", err)
	}

	response, ok := result.(InitializeResponse)
	if !ok {
		t.Fatalf("Expected InitializeResponse, got %T", result)
	}

	if response.ProtocolVersion != "2024-11-05" {
		t.Errorf("Expected protocol version '2024-11-05', got '%s'", response.ProtocolVersion)
	}

	if response.ServerInfo.Name != "petrock-mcp-server" {
		t.Errorf("Expected server name 'petrock-mcp-server', got '%s'", response.ServerInfo.Name)
	}

	// Test invalid protocol version
	params["protocolVersion"] = "invalid-version"
	_, err = server.handleInitialize(params)
	if err == nil {
		t.Fatal("Expected error for invalid protocol version")
	}
}

func TestMCPServer_HandleInitialized(t *testing.T) {
	server := NewMCPServer()

	if server.initialized {
		t.Fatal("Server should not be initialized initially")
	}

	server.handleInitialized(nil)

	if !server.initialized {
		t.Fatal("Server should be initialized after handleInitialized")
	}
}

func TestMCPServer_HandleListTools(t *testing.T) {
	server := NewMCPServer()

	// Test before initialization
	_, err := server.handleListTools(nil)
	if err == nil {
		t.Fatal("Expected error when server not initialized")
	}

	// Initialize server
	server.initialized = true

	result, err := server.handleListTools(nil)
	if err != nil {
		t.Fatalf("handleListTools() failed: %v", err)
	}

	response, ok := result.(ListToolsResponse)
	if !ok {
		t.Fatalf("Expected ListToolsResponse, got %T", result)
	}

	if len(response.Tools) == 0 {
		t.Fatal("Expected at least one tool (demo add tool)")
	}

	// Check for the demo add tool
	found := false
	for _, tool := range response.Tools {
		if tool.Name == "add" {
			found = true
			if tool.Description == "" {
				t.Error("Tool description is empty")
			}
			if tool.InputSchema.Type != "object" {
				t.Errorf("Expected input schema type 'object', got '%s'", tool.InputSchema.Type)
			}
			break
		}
	}

	if !found {
		t.Fatal("Demo 'add' tool not found in tools list")
	}
}

func TestMCPServer_HandleCallTool(t *testing.T) {
	server := NewMCPServer()

	// Test before initialization
	_, err := server.handleCallTool(nil)
	if err == nil {
		t.Fatal("Expected error when server not initialized")
	}

	// Initialize server
	server.initialized = true

	// Test calling non-existent tool
	params := map[string]interface{}{
		"name":      "nonexistent",
		"arguments": map[string]interface{}{},
	}

	result, err := server.handleCallTool(params)
	if err != nil {
		t.Fatalf("handleCallTool() should not return error for tool execution failures: %v", err)
	}

	response, ok := result.(CallToolResponse)
	if !ok {
		t.Fatalf("Expected CallToolResponse, got %T", result)
	}

	if !response.IsError {
		t.Fatal("Expected error response for non-existent tool")
	}

	// Test calling demo add tool with valid arguments
	params = map[string]interface{}{
		"name": "add",
		"arguments": map[string]interface{}{
			"a": 5.0,
			"b": 3.0,
		},
	}

	result, err = server.handleCallTool(params)
	if err != nil {
		t.Fatalf("handleCallTool() failed: %v", err)
	}

	response, ok = result.(CallToolResponse)
	if !ok {
		t.Fatalf("Expected CallToolResponse, got %T", result)
	}

	if response.IsError {
		t.Fatal("Tool execution should not have failed")
	}

	if len(response.Content) == 0 {
		t.Fatal("Expected response content")
	}

	// Parse the JSON response to check the result
	var resultData map[string]interface{}
	err = json.Unmarshal([]byte(response.Content[0].Text), &resultData)
	if err != nil {
		t.Fatalf("Failed to parse tool response: %v", err)
	}

	if resultData["result"] != 8.0 {
		t.Errorf("Expected result 8.0, got %v", resultData["result"])
	}

	// Test calling add tool with invalid arguments
	params = map[string]interface{}{
		"name": "add",
		"arguments": map[string]interface{}{
			"a": "invalid",
			"b": 3.0,
		},
	}

	result, err = server.handleCallTool(params)
	if err != nil {
		t.Fatalf("handleCallTool() should not return error: %v", err)
	}

	response, ok = result.(CallToolResponse)
	if !ok {
		t.Fatalf("Expected CallToolResponse, got %T", result)
	}

	if !response.IsError {
		t.Fatal("Expected error response for invalid arguments")
	}
}

func TestMCPServer_HandlePing(t *testing.T) {
	server := NewMCPServer()

	result, err := server.handlePing(nil)
	if err != nil {
		t.Fatalf("handlePing() failed: %v", err)
	}

	expected := map[string]interface{}{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected empty map, got %v", result)
	}
}

func TestMCPServer_RegisterTool(t *testing.T) {
	server := NewMCPServer()

	testTool := Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: ToolInputSchema{
			Type: "object",
		},
		Handler: func(params map[string]interface{}) (interface{}, error) {
			return "test result", nil
		},
	}

	server.RegisterTool(testTool)

	if _, exists := server.tools["test_tool"]; !exists {
		t.Fatal("Tool not registered")
	}

	// Test that the tool can be called
	server.initialized = true
	params := map[string]interface{}{
		"name":      "test_tool",
		"arguments": map[string]interface{}{},
	}

	result, err := server.handleCallTool(params)
	if err != nil {
		t.Fatalf("Tool call failed: %v", err)
	}

	response, ok := result.(CallToolResponse)
	if !ok {
		t.Fatalf("Expected CallToolResponse, got %T", result)
	}

	if response.IsError {
		t.Fatal("Tool execution should not have failed")
	}
}

func TestDemoAddTool(t *testing.T) {
	server := NewMCPServer()
	addTool := server.tools["add"]

	// Test valid addition
	params := map[string]interface{}{
		"a": 10.0,
		"b": 5.0,
	}

	result, err := addTool.Handler(params)
	if err != nil {
		t.Fatalf("Add tool failed: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["result"] != 15.0 {
		t.Errorf("Expected result 15.0, got %v", resultMap["result"])
	}

	// Test missing parameter
	params = map[string]interface{}{
		"a": 10.0,
	}

	_, err = addTool.Handler(params)
	if err == nil {
		t.Fatal("Expected error for missing parameter")
	}

	// Test invalid parameter type
	params = map[string]interface{}{
		"a": "invalid",
		"b": 5.0,
	}

	_, err = addTool.Handler(params)
	if err == nil {
		t.Fatal("Expected error for invalid parameter type")
	}
}
