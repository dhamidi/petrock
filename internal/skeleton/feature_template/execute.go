package petrock_example_feature_name

import (
	"context"
	"fmt"
	"log/slog"

	// "time" // Removed unused import

	// "github.com/google/uuid"                       // Removed unused import
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Executor implements the core.FeatureExecutor interface for this feature.
// It holds the feature's state and provides state update handlers.
// It also bridges validation calls from the central core.Executor to
// command structs that implement the feature's Validator interface.
type Executor struct {
	state *State // Dependency on the feature's state
}

// NewExecutor creates a new feature-specific Executor instance.
func NewExecutor(state *State) *Executor {
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
	slog.Debug("Feature executor validating command", "feature", "petrock_example_feature_name", "command_type", cmd.CommandName())

	// Check if the command implements the stateful validator interface defined in messages.go
	if validator, ok := cmd.(Validator); ok {
		slog.Debug("Command implements Validator, calling Validate(state)", "feature", "petrock_example_feature_name", "command_type", cmd.CommandName())
		// If yes, call the command's Validate method with the feature state
		return validator.Validate(e.state)
	}

	// If the command doesn't implement Validator, assume no stateful validation needed by the feature executor.
	// Basic stateless validation might have happened elsewhere (e.g., in the command struct itself, or HTTP handler).
	slog.Debug("Command does not implement Validator, skipping stateful validation", "feature", "petrock_example_feature_name", "command_type", cmd.CommandName())
	return nil
}

// --- State Update Handlers (Match core.CommandHandler signature) ---
// These methods are registered with the core.CommandRegistry and are called
// by the core.Executor *after* validation and logging.
// They should ONLY contain the logic to apply the state change.
// Returning an error here will cause the core.Executor to PANIC.

// HandleCreate applies state changes for CreateCommand.
func (e *Executor) HandleCreate(ctx context.Context, command core.Command, msg *core.Message) error {
	cmd, ok := command.(CreateCommand)
	if !ok {
		// This should ideally not happen if registration is correct, but check defensively.
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleCreate", command)
		slog.Error("Type assertion failed in HandleCreate", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)

	// Apply the change using the state's Apply method or direct state modification.
	// The state.Apply method (in state.go) contains the actual logic.
	// Pass through the message metadata if available (from replay)
	if err := e.state.Apply(cmd, msg); err != nil {
		// Log the error, but return it to trigger panic in core.Executor
		slog.Error("State Apply failed for CreateCommand", "error", err, "name", cmd.Name)
		return fmt.Errorf("state.Apply failed for CreateCommand: %w", err)
	}

	slog.Debug("State change applied successfully for CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)
	return nil
}

// HandleUpdate applies state changes for UpdateCommand.
func (e *Executor) HandleUpdate(ctx context.Context, command core.Command, msg *core.Message) error {
	cmd, ok := command.(UpdateCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleUpdate", command)
		slog.Error("Type assertion failed in HandleUpdate", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// Apply the change using the state's Apply method.
	// Pass through the message metadata if available (from replay)
	if err := e.state.Apply(cmd, msg); err != nil {
		slog.Error("State Apply failed for UpdateCommand", "error", err, "id", cmd.ID)
		return fmt.Errorf("state.Apply failed for UpdateCommand: %w", err)
	}

	slog.Debug("State change applied successfully for UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}

// HandleDelete applies state changes for DeleteCommand.
func (e *Executor) HandleDelete(ctx context.Context, command core.Command, msg *core.Message) error {
	cmd, ok := command.(DeleteCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleDelete", command)
		slog.Error("Type assertion failed in HandleDelete", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// Apply the change using the state's Apply method.
	// Pass through the message metadata if available (from replay)
	if err := e.state.Apply(cmd, msg); err != nil {
		slog.Error("State Apply failed for DeleteCommand", "error", err, "id", cmd.ID)
		return fmt.Errorf("state.Apply failed for DeleteCommand: %w", err)
	}

	slog.Debug("State change applied successfully for DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}

// Add more command handlers here...
