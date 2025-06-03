package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/petrock/example_module_path/core"                               // Placeholder for target project's core package
	"github.com/petrock/example_module_path/petrock_example_feature_name/state" // Import state package
)

// Ensure command implements the marker interfaces
var (
	_ core.Command = (*CreateCommand)(nil)
	_ Validator    = (*CreateCommand)(nil)
)

// CreateCommand holds data needed to create a new entity.
type CreateCommand struct {
	// Example fields - replace with actual data needed for creation
	Name        string    `json:"name" validate:"required,minlen=2,maxlen=100"`
	Description string    `json:"description" validate:"required,minlen=5,maxlen=500"`
	Content     string    `json:"content" validate:"required,minlen=10"`
	CreatedBy   string    `json:"created_by"` // e.g., User ID
	CreatedAt   time.Time `json:"created_at"` // Timestamp when created
}

// CommandName returns the unique kebab-case name for this command type.
func (c *CreateCommand) CommandName() string {
	return "petrock_example_feature_name/create" // Removed suffix
}

// Validate implements the Validator interface for CreateCommand.
// It performs validation checks, potentially using the current state.
func (c *CreateCommand) Validate(state *state.State) error {
	// Trim all string fields
	trimmedName := strings.TrimSpace(c.Name)
	trimmedDescription := strings.TrimSpace(c.Description)
	trimmedContent := strings.TrimSpace(c.Content)

	// Basic stateless validation
	if trimmedName == "" {
		return errors.New("item name cannot be empty")
	}

	if trimmedDescription == "" {
		return errors.New("item description cannot be empty")
	}

	if trimmedContent == "" {
		return errors.New("item content cannot be empty")
	}

	// Example stateful validation: Check if an item with the same name already exists
	// This simplistic approach just looks for items with the same name
	items, _ := state.ListItems(1, 1000, "")
	for _, item := range items {
		if item.Name == trimmedName {
			return fmt.Errorf("item with name %q already exists", trimmedName)
		}
	}

	// Add other validation rules...
	return nil
}

// HandleCreate applies state changes for CreateCommand.
func (e *Executor) HandleCreate(ctx context.Context, command core.Command, msg *core.Message, pctx *core.ProcessingContext) error {
	// Type assertion for pointer type
	cmd, ok := command.(*CreateCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleCreate, expected *CreateCommand", command)
		slog.Error("Type assertion failed in HandleCreate", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)

	// Create a new item
	newItem := &state.Item{
		ID:          cmd.Name, // Use Name as ID for simplicity, replace with generated ID if needed
		Name:        cmd.Name,
		Description: cmd.Description,
		Content:     cmd.Content,
		CreatedAt:   getTimestamp(msg), // Use message timestamp if available, otherwise current time
		UpdatedAt:   getTimestamp(msg),
		Version:     1,
	}

	// Add item to state
	if err := e.state.AddItem(newItem); err != nil {
		// Log the error, but return it to trigger panic in core.Executor
		slog.Error("Failed to add item to state", "error", err, "name", cmd.Name)
		return fmt.Errorf("failed to add item to state: %w", err)
	}

	slog.Debug("State change applied successfully for CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)
	return nil
}
