package petrock_example_feature_name

import (
	"context"
	"errors" // Import errors package
	"fmt"
	"log/slog"
	// "time" // Removed unused import

	// "github.com/google/uuid"                       // Removed unused import
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// FeatureExecutor handles domain-specific command processing for the feature.
// It focuses on applying commands to the feature's state.
type FeatureExecutor struct {
	state *State // Dependency on the feature's state
}

// NewFeatureExecutor creates a new FeatureExecutor instance.
func NewFeatureExecutor(state *State) *FeatureExecutor {
	if state == nil {
		panic("state cannot be nil for FeatureExecutor") // Or return an error
	}
	return &FeatureExecutor{
		state: state,
	}
}

// HandleCreate processes the CreateCommand.
// This function signature matches core.CommandHandler.
func (e *FeatureExecutor) HandleCreate(ctx context.Context, command core.Command) error {
	cmd, ok := command.(CreateCommand)
	if !ok {
		return fmt.Errorf("invalid command type for HandleCreate: expected CreateCommand, got %T", command)
	}

	slog.Debug("Handling CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)

	// Apply to State - validation and logging already handled by core.Executor
	if err := e.state.Apply(cmd); err != nil {
		slog.Error("Failed to apply CreateCommand to state", "error", err, "name", cmd.Name)
		return fmt.Errorf("failed to update state: %w", err)
	}

	slog.Debug("Successfully processed CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)
	return nil
}

// HandleUpdate processes the UpdateCommand.
// This function signature matches core.CommandHandler.
func (e *FeatureExecutor) HandleUpdate(ctx context.Context, command core.Command) error {
	cmd, ok := command.(UpdateCommand)
	if !ok {
		return fmt.Errorf("invalid command type for HandleUpdate: expected UpdateCommand, got %T", command)
	}

	slog.Debug("Handling UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// Domain checks (business rule validation, not input validation)
	_, found := e.state.GetItem(cmd.ID)
	if !found {
		slog.Warn("Attempted to update non-existent item", "id", cmd.ID)
		return fmt.Errorf("item with ID %s not found", cmd.ID) // Return a "not found" error
	}

	// Apply to State - validation and logging already handled by core.Executor
	if err := e.state.Apply(cmd); err != nil {
		slog.Error("Failed to apply UpdateCommand to state", "error", err, "id", cmd.ID)
		return fmt.Errorf("failed to update state: %w", err)
	}

	slog.Debug("Successfully processed UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}

// HandleDelete processes the DeleteCommand.
// This function signature matches core.CommandHandler.
func (e *FeatureExecutor) HandleDelete(ctx context.Context, command core.Command) error {
	cmd, ok := command.(DeleteCommand)
	if !ok {
		return fmt.Errorf("invalid command type for HandleDelete: expected DeleteCommand, got %T", command)
	}

	slog.Debug("Handling DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// Domain checks (business rule validation, not input validation)
	_, found := e.state.GetItem(cmd.ID)
	if !found {
		slog.Warn("Attempted to delete non-existent item", "id", cmd.ID)
		// Decide if this is an error or idempotent success
		return fmt.Errorf("item with ID %s not found", cmd.ID) // Return error
		// return nil // Alternative: Treat deletion of non-existent item as success
	}

	// Apply to State - validation and logging already handled by core.Executor
	if err := e.state.Apply(cmd); err != nil {
		slog.Error("Failed to apply DeleteCommand to state", "error", err, "id", cmd.ID)
		return fmt.Errorf("failed to update state: %w", err)
	}

	slog.Debug("Successfully processed DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}

// Add more command handlers here...
