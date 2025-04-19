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

	// Apply the change using the state's Apply method
	if err := e.state.Apply(cmd, msg); err != nil {
		slog.Error("State Apply failed for UpdateCommand", "error", err, "id", cmd.ID)
		return fmt.Errorf("state.Apply failed for UpdateCommand: %w", err)
	}

	slog.Debug("State change applied successfully for UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}
