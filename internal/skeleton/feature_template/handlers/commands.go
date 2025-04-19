package handlers

import (
	"time"
)

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