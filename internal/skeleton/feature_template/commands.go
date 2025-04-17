package petrock_example_feature_name

import (
	"errors"  // Added for validation errors
	"fmt"     // Added for validation errors
	"strings" // Added for string trimming
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Validator defines an interface for commands that require stateful validation.
// The feature's Executor will call this method if implemented by a command.
type Validator interface {
	Validate(state *State) error
}

// --- Commands (Implement core.Command) ---
// Commands represent intentions to change the system state.
// All commands MUST be created and used as pointer types (*CommandType).
// Example usage: cmd := &CreateCommand{Name: "example"}

// CreateCommand holds data needed to create a new entity.
type CreateCommand struct {
	// Example fields - replace with actual data needed for creation
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"` // e.g., User ID
	CreatedAt   time.Time `json:"created_at"` // Timestamp when created
}

// CommandName returns the unique kebab-case name for this command type.
func (c *CreateCommand) CommandName() string {
	return "petrock_example_feature_name/create" // Removed suffix
}

// Validate implements the Validator interface for CreateCommand.
// It performs validation checks, potentially using the current state.
func (c *CreateCommand) Validate(state *State) error {
	// Trim all string fields
	trimmedName := strings.TrimSpace(c.Name)
	trimmedDescription := strings.TrimSpace(c.Description)

	// Basic stateless validation
	if trimmedName == "" {
		return errors.New("item name cannot be empty")
	}

	if trimmedDescription == "" {
		return errors.New("item description cannot be empty")
	}

	// Example stateful validation: Check if an item with the same name already exists
	// Note: state.GetItem currently uses ID, not name. If using name as ID on create,
	// this check is relevant. Adjust based on actual ID strategy.
	// Assuming state.GetItem uses the ID field from the Item struct.
	// If CreateCommand implies using Name as the potential ID:
	state.mu.RLock() // Read lock for checking existence
	_, exists := state.Items[trimmedName]
	state.mu.RUnlock()
	if exists {
		return fmt.Errorf("item with name %q already exists", trimmedName)
	}

	// Add other validation rules...
	return nil
}

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

// RequestSummaryGenerationCommand requests a summary be generated for an item
type RequestSummaryGenerationCommand struct {
	ID        string `json:"id"`         // ID of the item to summarize
	RequestID string `json:"request_id"` // Unique ID for this summary request
}

// CommandName returns the unique kebab-case name for this command type
func (c *RequestSummaryGenerationCommand) CommandName() string {
	return "petrock_example_feature_name/request-summary-generation"
}

// Validate implements the Validator interface
func (c *RequestSummaryGenerationCommand) Validate(state *State) error {
	if strings.TrimSpace(c.ID) == "" {
		return errors.New("item ID cannot be empty")
	}
	if strings.TrimSpace(c.RequestID) == "" {
		return errors.New("request ID cannot be empty")
	}

	// Verify the item exists
	_, found := state.GetItem(c.ID)
	if !found {
		return fmt.Errorf("item with ID %q not found", c.ID)
	}
	return nil
}

// FailSummaryGenerationCommand indicates a summary generation request failed
type FailSummaryGenerationCommand struct {
	ID        string `json:"id"`         // ID of the item
	RequestID string `json:"request_id"` // References the original request
	Reason    string `json:"reason"`     // Reason for failure
}

// CommandName returns the unique kebab-case name for this command type
func (c *FailSummaryGenerationCommand) CommandName() string {
	return "petrock_example_feature_name/fail-summary-generation"
}

// SetGeneratedSummaryCommand sets the generated summary for an item
type SetGeneratedSummaryCommand struct {
	ID        string `json:"id"`         // ID of the item
	RequestID string `json:"request_id"` // References the original request
	Summary   string `json:"summary"`    // The generated summary text
}

// CommandName returns the unique kebab-case name for this command type
func (c *SetGeneratedSummaryCommand) CommandName() string {
	return "petrock_example_feature_name/set-generated-summary"
}

// Ensure commands implement the marker interface (optional) and Validator where applicable
var _ core.Command = (*CreateCommand)(nil)
var _ Validator = (*CreateCommand)(nil)
var _ core.Command = (*UpdateCommand)(nil)
var _ Validator = (*UpdateCommand)(nil)
var _ core.Command = (*DeleteCommand)(nil)
var _ Validator = (*DeleteCommand)(nil)
var _ core.Command = (*RequestSummaryGenerationCommand)(nil)
var _ Validator = (*RequestSummaryGenerationCommand)(nil)
var _ core.Command = (*FailSummaryGenerationCommand)(nil)
var _ core.Command = (*SetGeneratedSummaryCommand)(nil)
