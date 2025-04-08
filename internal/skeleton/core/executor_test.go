package core

import (
	"context"
	"errors"
	"testing"
)

// MockCommand implements a test command with validation
type MockCommand struct {
	ShouldValidate bool
	ValidateError  error
}

func (c MockCommand) CommandName() string {
	return "mock/command"
}

func (c MockCommand) Validate() error {
	return c.ValidateError
}

// MockMessageLog implements MessageLog for testing
type MockMessageLog struct {
	AppendCalled bool
	AppendError  error
	LastCommand  Command
}

func (l *MockMessageLog) Append(ctx context.Context, cmd interface{}) error {
	l.AppendCalled = true
	l.LastCommand = cmd.(Command)
	return l.AppendError
}

// Stub the other methods that MessageLog would have
func (l *MockMessageLog) RegisterType(instance interface{}) {}
func (l *MockMessageLog) Load(ctx context.Context) ([]interface{}, error) { return nil, nil }

// MockCommandRegistry implements CommandRegistry for testing
type MockCommandRegistry struct {
	DispatchCalled bool
	DispatchError  error
	LastCommand    Command
}

func (r *MockCommandRegistry) Dispatch(ctx context.Context, cmd Command) error {
	r.DispatchCalled = true
	r.LastCommand = cmd
	return r.DispatchError
}

// Stub the other methods that CommandRegistry would have
func (r *MockCommandRegistry) Register(cmd Command, handler CommandHandler) {}
func (r *MockCommandRegistry) RegisteredCommandNames() []string { return nil }
func (r *MockCommandRegistry) GetCommandType(name string) (interface{}, bool) { return nil, false }

func TestBaseExecutor_Execute(t *testing.T) {
	tests := []struct {
		name            string
		cmd             MockCommand
		logError        error
		dispatchError   error
		expectLogCalled bool
		expectDispCalled bool
		wantErr         bool
		validationErr   bool
	}{
		{
			name:             "successful execution",
			cmd:              MockCommand{ValidateError: nil},
			logError:         nil,
			dispatchError:    nil,
			expectLogCalled:  true,
			expectDispCalled: true,
			wantErr:          false,
			validationErr:    false,
		},
		{
			name:             "validation fails",
			cmd:              MockCommand{ValidateError: errors.New("validation error")},
			logError:         nil,
			dispatchError:    nil,
			expectLogCalled:  false,
			expectDispCalled: false,
			wantErr:          true,
			validationErr:    true,
		},
		{
			name:             "logging fails",
			cmd:              MockCommand{ValidateError: nil},
			logError:         errors.New("log error"),
			dispatchError:    nil,
			expectLogCalled:  true,
			expectDispCalled: false,
			wantErr:          true,
			validationErr:    false,
		},
		{
			name:             "dispatch fails",
			cmd:              MockCommand{ValidateError: nil},
			logError:         nil,
			dispatchError:    errors.New("dispatch error"),
			expectLogCalled:  true,
			expectDispCalled: true,
			wantErr:          true,
			validationErr:    false,
		},
	}

	for _, tt := range tests {
		t.cmd.ShouldValidate = true
		
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			log := &MockMessageLog{AppendError: tt.logError}
			reg := &MockCommandRegistry{DispatchError: tt.dispatchError}
			
			// Create executor
			executor := NewBaseExecutor(log, reg)
			
			// Execute the command
			err := executor.Execute(context.Background(), tt.cmd)
			
			// Check results
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if tt.validationErr && !IsValidationError(err) {
				t.Error("Expected ValidationError but got different error type")
			}
			
			if log.AppendCalled != tt.expectLogCalled {
				t.Errorf("Log.Append called = %v, want %v", log.AppendCalled, tt.expectLogCalled)
			}
			
			if reg.DispatchCalled != tt.expectDispCalled {
				t.Errorf("Registry.Dispatch called = %v, want %v", reg.DispatchCalled, tt.expectDispCalled)
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	// Create a validation error
	cmd := MockCommand{}
	valErr := NewValidationError(cmd, errors.New("test error"), nil)
	
	// Test with validation error
	if !IsValidationError(valErr) {
		t.Error("IsValidationError failed to recognize ValidationError")
	}
	
	// Test with regular error
	regErr := errors.New("regular error")
	if IsValidationError(regErr) {
		t.Error("IsValidationError incorrectly identified regular error as ValidationError")
	}
}

func TestNewValidationError(t *testing.T) {
	cmd := MockCommand{}
	underlyingErr := errors.New("test error")
	fields := map[string]string{"field1": "error1"}
	
	valErr := NewValidationError(cmd, underlyingErr, fields)
	
	if valErr.CommandName != cmd.CommandName() {
		t.Errorf("Expected command name %s, got %s", cmd.CommandName(), valErr.CommandName)
	}
	
	if valErr.Err != underlyingErr {
		t.Error("Underlying error not properly stored")
	}
	
	if valErr.Fields["field1"] != "error1" {
		t.Error("Fields not properly stored")
	}
}