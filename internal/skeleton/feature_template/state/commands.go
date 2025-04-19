package state

import (
	"time"
)

// Command types used in the state.Apply method

// CreateCommand holds data needed to create a new entity.
type CreateCommand struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"` // e.g., User ID
	CreatedAt   time.Time `json:"created_at"` // Timestamp when created
}

// CommandName returns the unique kebab-case name for this command type.
func (c *CreateCommand) CommandName() string {
	return "petrock_example_feature_name/create"
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
	return "petrock_example_feature_name/update"
}

// DeleteCommand holds data needed to delete an entity.
type DeleteCommand struct {
	ID        string    `json:"id"` // ID of the entity to delete
	DeletedBy string    `json:"deleted_by"`
	DeletedAt time.Time `json:"deleted_at"` // Timestamp when deleted
}

// CommandName returns the unique kebab-case name for this command type.
func (c *DeleteCommand) CommandName() string {
	return "petrock_example_feature_name/delete"
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
