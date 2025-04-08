# Plan for core/executor.go

This file defines the centralized command execution system that standardizes how commands are processed across the application.

## Types

- `Executor`: An interface defining the contract for command execution.
  - `Execute(ctx context.Context, cmd Command) error`: Executes a command following standard validation, logging, and state application workflow.

- `BaseExecutor`: A concrete implementation of `Executor` that provides a standard execution flow.
  - `log *MessageLog`: For persisting commands
  - `cmdRegistry *CommandRegistry`: For accessing command handlers

- `ValidationError`: A specialized error type for command validation failures.
  - `Wrapped error`: The underlying error
  - `CommandName string`: The name of the command that failed validation
  - `Fields map[string]string`: Map of field names to error messages

## Functions

- `NewBaseExecutor(log *MessageLog, cmdRegistry *CommandRegistry) *BaseExecutor`: Constructor for the standard executor implementation.

- `(e *BaseExecutor) Execute(ctx context.Context, cmd Command) error`: Implements the standard execution flow:
  1. Validates the command if it implements `Validator` interface
  2. Persists the command to the message log
  3. Dispatches the command to its registered handler
  4. Returns appropriate errors for validation or execution failures

- `IsValidationError(err error) bool`: Helper function to check if an error is a validation error.

## Interfaces

- `Validator`: An interface that commands can implement to enable self-validation.
  - `Validate() error`: Validates the command's fields and returns an error if invalid.