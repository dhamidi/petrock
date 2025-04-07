package petrock_example_feature_name

import (
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

// UpdateCommand holds data needed to update an existing entity.
type UpdateCommand struct {
	ID          string `json:"id"` // ID of the entity to update
	Name        string `json:"name"`
	Description string `json:"description"`
	UpdatedBy   string `json:"updated_by"`
}

// DeleteCommand holds data needed to delete an entity.
type DeleteCommand struct {
	ID        string `json:"id"` // ID of the entity to delete
	DeletedBy string `json:"deleted_by"`
}

// --- Queries (Implement core.Query) ---
// Queries represent requests to read system state.

// GetQuery holds data needed to retrieve a single entity.
type GetQuery struct {
	ID string // ID of the entity to retrieve
}

// ListQuery holds data needed to retrieve a list of entities, possibly filtered or paginated.
type ListQuery struct {
	Page     int    `json:"page"`      // For pagination
	PageSize int    `json:"page_size"` // For pagination
	Filter   string `json:"filter"`    // Example filter criteria
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
