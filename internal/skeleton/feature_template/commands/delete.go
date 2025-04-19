package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure command implements the marker interfaces
var _ core.Command = (*DeleteCommand)(nil)
var _ Validator = (*DeleteCommand)(nil)

// DeleteCommand holds data needed to delete an entity.
type DeleteCommand struct {
	ID        string    `json:"id"` // ID of the entity to delete
	DeletedBy string    `json:"deleted_by"`
	DeletedAt time.Time `json:"deleted_at"` // Timestamp when deleted
}

// CommandName returns the unique kebab-case name for this command type.
func (c *DeleteCommand) CommandName() string {
	return "petrock_example_feature_name/delete" // Removed suffix
}

// Validate implements the Validator interface for DeleteCommand.
func (c *DeleteCommand) Validate(state *State) error {
	// Trim all string fields
	trimmedID := strings.TrimSpace(c.ID)

	// Basic stateless validation
	if trimmedID == "" {
		return errors.New("item ID cannot be empty for deletion")
	}

	// Example stateful validation: Check if the item exists
	_, found := state.GetItem(trimmedID) // GetItem handles locking
	if !found {
		// Decide if deleting a non-existent item is an error or idempotent success
		return fmt.Errorf("item with ID %q not found", trimmedID) // Return error
		// return nil // Alternative: Treat as success
	}
	// Add other validation rules (e.g., check if item is deletable based on status)
	return nil
}

// HandleDelete applies state changes for DeleteCommand.
func (e *Executor) HandleDelete(ctx context.Context, command core.Command, msg *core.Message) error {
	// Type assertion for pointer type
	cmd, ok := command.(*DeleteCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleDelete, expected *DeleteCommand", command)
		slog.Error("Type assertion failed in HandleDelete", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// Check if the item exists first
	_, found := e.state.GetItem(cmd.ID)
	if !found {
		err := fmt.Errorf("item with ID %s not found for deletion", cmd.ID)
		slog.Error("Delete failed", "error", err, "id", cmd.ID)
		return err
	}
	
	// Delete the item
	if err := e.state.DeleteItem(cmd.ID); err != nil {
		slog.Error("Failed to delete item from state", "error", err, "id", cmd.ID)
		return fmt.Errorf("failed to delete item from state: %w", err)
	}

	slog.Debug("State change applied successfully for DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}
