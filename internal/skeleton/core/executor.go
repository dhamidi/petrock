package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// Executor defines the interface for executing commands with standardized flow
type Executor interface {
	// Execute processes a command through validation, logging, and state application
	Execute(ctx context.Context, cmd Command) error
}

// Validator is the interface that commands should implement to enable self-validation
type Validator interface {
	// Validate checks if the command is well-formed and returns an error if not
	Validate() error
}

// ValidationError represents an error that occurred during command validation
type ValidationError struct {
	CommandName string            // The name of the command that failed validation
	Fields      map[string]string // Field-specific error messages
	Err         error             // The underlying error
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for command %s: %v", e.CommandName, e.Err)
}

// Unwrap returns the underlying error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new ValidationError
func NewValidationError(cmd Command, err error, fields map[string]string) *ValidationError {
	return &ValidationError{
		CommandName: cmd.CommandName(),
		Err:         err,
		Fields:      fields,
	}
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	var valErr *ValidationError
	return errors.As(err, &valErr)
}

// BaseExecutor provides a standard implementation of command execution
type BaseExecutor struct {
	log         *MessageLog
	cmdRegistry *CommandRegistry
}

// NewBaseExecutor creates a new executor with dependencies
func NewBaseExecutor(log *MessageLog, cmdRegistry *CommandRegistry) *BaseExecutor {
	if log == nil {
		panic("message log cannot be nil")
	}
	if cmdRegistry == nil {
		panic("command registry cannot be nil")
	}

	return &BaseExecutor{
		log:         log,
		cmdRegistry: cmdRegistry,
	}
}

// Execute implements the standard command execution flow:
// 1. Validate the command if it implements Validator
// 2. Log the command to the message log
// 3. Dispatch the command to its registered handler
func (e *BaseExecutor) Execute(ctx context.Context, cmd Command) error {
	// 1. Validate the command if it implements Validator
	if validator, ok := cmd.(Validator); ok {
		if err := validator.Validate(); err != nil {
			// Create a validation error with the underlying error
			slog.Debug("Command validation failed", 
				"command", cmd.CommandName(), 
				"error", err)

			return NewValidationError(cmd, err, nil)
		}
	}

	// 2. Log the command to the message log
	if err := e.log.Append(ctx, cmd); err != nil {
		slog.Error("Failed to append command to log", 
			"command", cmd.CommandName(), 
			"error", err)
		return fmt.Errorf("failed to persist command: %w", err)
	}

	slog.Debug("Command appended to log", "command", cmd.CommandName())

	// 3. Dispatch the command to its registered handler
	if err := e.cmdRegistry.Dispatch(ctx, cmd); err != nil {
		slog.Error("Failed to dispatch command", 
			"command", cmd.CommandName(), 
			"error", err)
		return fmt.Errorf("failed to execute command handler: %w", err)
	}

	slog.Debug("Command executed successfully", "command", cmd.CommandName())
	return nil
}