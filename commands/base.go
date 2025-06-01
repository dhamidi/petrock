package commands

import (
	"context"
	// "fmt" - Will be used in implementation
	"log/slog"
	"time"

	"github.com/dhamidi/petrock/core" // Placeholder for target project's core package
	"github.com/dhamidi/petrock/test/state" // Import state package
)

// Validator defines an interface for commands that require stateful validation.
// The feature's Executor will call this method if implemented by a command.
type Validator interface {
	Validate(state *state.State) error
}

// Executor implements the core.FeatureExecutor interface for this feature.
// It holds the feature's state and provides state update handlers.
// It also bridges validation calls from the central core.Executor to
// command structs that implement the feature's Validator interface.
type Executor struct {
	state *state.State // Dependency on the feature's state
}

// NewExecutor creates a new feature-specific Executor instance.
func NewExecutor(state *state.State) *Executor {
	if state == nil {
		panic("state cannot be nil for feature Executor")
	}
	return &Executor{
		state: state,
	}
}

// --- core.FeatureExecutor Implementation ---

// ValidateCommand is called by the central core.Executor.
// It checks if the command implements the feature's Validator interface
// and calls its Validate method with the feature's state if it does.
func (e *Executor) ValidateCommand(ctx context.Context, cmd core.Command) error {
	slog.Debug("Feature executor validating command", "feature", "test", "command_type", cmd.CommandName())

	// Check if the command implements the stateful validator interface defined in commands.go
	if validator, ok := cmd.(Validator); ok {
		slog.Debug("Command implements Validator, calling Validate(state)", "feature", "test", "command_type", cmd.CommandName())
		// If yes, call the command's Validate method with the feature state
		return validator.Validate(e.state)
	}

	// If the command doesn't implement Validator, assume no stateful validation needed by the feature executor.
	// Basic stateless validation might have happened elsewhere (e.g., in the command struct itself, or HTTP handler).
	slog.Debug("Command does not implement Validator, skipping stateful validation", "feature", "test", "command_type", cmd.CommandName())
	return nil
}

// getTimestamp returns the timestamp from the message metadata if available, otherwise current time
func getTimestamp(msg *core.Message) time.Time {
	if msg != nil {
		return msg.Timestamp
	}
	return time.Now().UTC()
}
