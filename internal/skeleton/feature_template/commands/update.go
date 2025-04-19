package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
	"github.com/petrock/example_module_path/petrock_example_feature_name/state" // Import state package
)

// Ensure command implements the marker interfaces
var _ core.Command = (*UpdateCommand)(nil)
var _ Validator = (*UpdateCommand)(nil)

// UpdateCommand holds data needed to update an existing entity.
type UpdateCommand struct {
	ID          string    `json:"id"` // ID of the entity to update
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"` // Timestamp when updated
}

// CommandName returns the unique kebab-case name for this command type.
func (c *UpdateCommand) CommandName() string {
	return "petrock_example_feature_name/update" // Removed suffix
}

// Validate implements the Validator interface for UpdateCommand.
func (c *UpdateCommand) Validate(state *State) error {
	// Trim all string fields
	trimmedID := strings.TrimSpace(c.ID)
	trimmedName := strings.TrimSpace(c.Name)
	trimmedDescription := strings.TrimSpace(c.Description)

	// Basic stateless validation
	if trimmedID == "" {
		return errors.New("item ID cannot be empty for update")
	}

	if trimmedName == "" {
		return errors.New("item name cannot be empty")
	}

	if trimmedDescription == "" {
		return errors.New("item description cannot be empty")
	}

	// Example stateful validation: Check if the item exists
	_, found := state.GetItem(trimmedID) // GetItem handles locking
	if !found {
		return fmt.Errorf("item with ID %q not found", trimmedID)
	}
	// Example: Check if updating the name conflicts with another existing item's name
	// state.mu.RLock()
	// for id, item := range state.Items {
	//     if id != c.ID && item.Name == c.Name {
	//         state.mu.RUnlock()
	//         return fmt.Errorf("another item with name %q already exists", c.Name)
	//     }
	// }
	// state.mu.RUnlock()

	// Add other validation rules...
	return nil
}

// HandleUpdate applies state changes for UpdateCommand.
func (e *Executor) HandleUpdate(ctx context.Context, command core.Command, msg *core.Message) error {
	// Type assertion for pointer type
	cmd, ok := command.(*UpdateCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleUpdate, expected *UpdateCommand", command)
		slog.Error("Type assertion failed in HandleUpdate", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// Get the existing item
	existingItem, found := e.state.GetItem(cmd.ID)
	if !found {
		err := fmt.Errorf("item with ID %s not found for update", cmd.ID)
		slog.Error("Update failed", "error", err, "id", cmd.ID)
		return err
	}
	
	// Update the item properties
	existingItem.Name = cmd.Name
	existingItem.Description = cmd.Description
	existingItem.UpdatedAt = getTimestamp(msg)
	existingItem.Version++
	
	// Save the updated item
	if err := e.state.UpdateItem(existingItem); err != nil {
		slog.Error("Failed to update item in state", "error", err, "id", cmd.ID)
		return fmt.Errorf("failed to update item in state: %w", err)
	}

	slog.Debug("State change applied successfully for UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}
