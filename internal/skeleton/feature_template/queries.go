package petrock_example_feature_name

import (
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

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

// Ensure queries implement the marker interface (optional)
var _ core.Query = (*GetQuery)(nil)
var _ core.Query = (*ListQuery)(nil)