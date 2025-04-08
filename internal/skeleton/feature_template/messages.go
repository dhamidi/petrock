package petrock_example_feature_name

import (
	"fmt"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// --- Commands (Implement core.Command) ---
// Commands represent intentions to change the system state.

// CreateCommand holds data needed to create a new entity.
type CreateCommand struct {
	// Example fields - replace with actual data needed for creation
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"` // e.g., User ID
}

// CommandName returns the unique kebab-case name for this command type.
func (c CreateCommand) CommandName() string {
	return "petrock_example_feature_name/create" // Removed suffix
}

// Validate implements core.Validator interface to enable self-validation.
func (c CreateCommand) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	// Add other validation rules as needed
	return nil
}

// UpdateCommand holds data needed to update an existing entity.
type UpdateCommand struct {
	ID          string `json:"id"` // ID of the entity to update
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedBy   string `json:"updated_by"`
}

// CommandName returns the unique kebab-case name for this command type.
func (c UpdateCommand) CommandName() string {
	return "petrock_example_feature_name/update" // Removed suffix
}

// Validate implements core.Validator interface to enable self-validation.
func (c UpdateCommand) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	// Add other validation rules as needed
	return nil
}

// DeleteCommand holds data needed to delete an entity.
type DeleteCommand struct {
	ID        string `json:"id"` // ID of the entity to delete
	DeletedBy string `json:"deleted_by"`
}

// CommandName returns the unique kebab-case name for this command type.
func (c DeleteCommand) CommandName() string {
	return "petrock_example_feature_name/delete" // Removed suffix
}

// Validate implements core.Validator interface to enable self-validation.
func (c DeleteCommand) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}
	// Add other validation rules as needed
	return nil
}

// --- Queries (Implement core.Query) ---
// Queries represent requests to read system state.

// GetQuery holds data needed to retrieve a single entity.
type GetQuery struct {
	ID string // ID of the entity to retrieve
}

// QueryName returns the unique kebab-case name for this query type.
func (q GetQuery) QueryName() string {
	return "petrock_example_feature_name/get" // Removed suffix
}

// Validate implements validation for the query.
func (q GetQuery) Validate() error {
	if q.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}
	return nil
}

// ListQuery holds data needed to retrieve a list of entities, possibly filtered or paginated.
type ListQuery struct {
	Page     int    `json:"page"`      // For pagination
	PageSize int    `json:"page_size"` // For pagination
	Filter   string `json:"filter"`    // Example filter criteria
}

// QueryName returns the unique kebab-case name for this query type.
func (q ListQuery) QueryName() string {
	return "petrock_example_feature_name/list" // Removed suffix
}

// Validate implements validation for the query.
func (q ListQuery) Validate() error {
	if q.Page < 0 {
		return fmt.Errorf("page cannot be negative")
	}
	if q.PageSize < 1 {
		return fmt.Errorf("page size must be positive")
	}
	return nil
}

// --- Query Results (Implement core.QueryResult) ---
// QueryResults represent the data returned by query handlers.

// Result holds the data for a single entity returned by a query.
type Result struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version"` // Example version field
}

// ListResult holds a list of entities and pagination details.
type ListResult struct {
	Items      []Result `json:"items"`
	TotalCount int      `json:"total_count"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
}

// Ensure query results implement the marker interface (optional but good practice)
var _ core.QueryResult = (*Result)(nil)
var _ core.QueryResult = (*ListResult)(nil)

// Ensure commands implement the marker interface (optional)
var _ core.Command = (*CreateCommand)(nil)
var _ core.Command = (*UpdateCommand)(nil)
var _ core.Command = (*DeleteCommand)(nil)

// Ensure queries implement the marker interface (optional)
var _ core.Query = (*GetQuery)(nil)
var _ core.Query = (*ListQuery)(nil)
