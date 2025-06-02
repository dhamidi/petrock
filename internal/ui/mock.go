package ui

import (
	"context"
	"fmt"
	"sync"
)

// CapturedMessage represents a message captured by the MockUI
type CapturedMessage struct {
	Type    MessageType
	Message string
	Args    []interface{}
}

// MockUI implements the UI interface for testing purposes
type MockUI struct {
	mu         sync.RWMutex
	messages   []CapturedMessage
	prompts    []string
	responses  []string
	errors     []error
	progressStates []ProgressState
}

// NewMockUI creates a new mock UI implementation for testing
func NewMockUI() *MockUI {
	return &MockUI{
		messages:   make([]CapturedMessage, 0),
		prompts:    make([]string, 0),
		responses:  make([]string, 0),
		errors:     make([]error, 0),
		progressStates: make([]ProgressState, 0),
	}
}

// Present captures the message for later verification in tests
func (m *MockUI) Present(ctx context.Context, msgType MessageType, message string, args ...interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.messages = append(m.messages, CapturedMessage{
		Type:    msgType,
		Message: message,
		Args:    args,
	})
	
	return nil
}

// Prompt captures the question and returns a pre-configured response
func (m *MockUI) Prompt(ctx context.Context, question string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.prompts = append(m.prompts, question)
	
	if len(m.responses) > 0 {
		response := m.responses[0]
		m.responses = m.responses[1:]
		return response, nil
	}
	
	return "", fmt.Errorf("no response configured for prompt: %s", question)
}

// ShowProgress captures the progress state for later verification
func (m *MockUI) ShowProgress(ctx context.Context, state ProgressState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.progressStates = append(m.progressStates, state)
	return nil
}

// ShowError captures the error for later verification
func (m *MockUI) ShowError(ctx context.Context, err error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.errors = append(m.errors, err)
	return nil
}

// ShowSuccess captures the success message for later verification
func (m *MockUI) ShowSuccess(ctx context.Context, message string, args ...interface{}) error {
	return m.Present(ctx, MessageTypeSuccess, message, args...)
}

// ShowHeader captures the header for later verification
func (m *MockUI) ShowHeader(ctx context.Context, title string) error {
	return m.Present(ctx, MessageTypeInfo, title)
}

// ShowFileOperation captures the file operation for later verification
func (m *MockUI) ShowFileOperation(ctx context.Context, operation, filePath string) error {
	return m.Present(ctx, MessageTypeInfo, "%s  %s", operation, filePath)
}

// GetMessages returns all captured messages
func (m *MockUI) GetMessages() []CapturedMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	messages := make([]CapturedMessage, len(m.messages))
	copy(messages, m.messages)
	return messages
}

// GetPrompts returns all captured prompts
func (m *MockUI) GetPrompts() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	prompts := make([]string, len(m.prompts))
	copy(prompts, m.prompts)
	return prompts
}

// GetErrors returns all captured errors
func (m *MockUI) GetErrors() []error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	errors := make([]error, len(m.errors))
	copy(errors, m.errors)
	return errors
}

// GetProgressStates returns all captured progress states
func (m *MockUI) GetProgressStates() []ProgressState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	states := make([]ProgressState, len(m.progressStates))
	copy(states, m.progressStates)
	return states
}

// ClearMessages clears all captured messages
func (m *MockUI) ClearMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.messages = m.messages[:0]
	m.prompts = m.prompts[:0]
	m.errors = m.errors[:0]
	m.progressStates = m.progressStates[:0]
}

// SetResponses configures responses for future Prompt calls
func (m *MockUI) SetResponses(responses ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.responses = make([]string, len(responses))
	copy(m.responses, responses)
}

// LastMessage returns the last captured message, or nil if none
func (m *MockUI) LastMessage() *CapturedMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if len(m.messages) == 0 {
		return nil
	}
	
	return &m.messages[len(m.messages)-1]
}

// MessageCount returns the number of captured messages
func (m *MockUI) MessageCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return len(m.messages)
}
