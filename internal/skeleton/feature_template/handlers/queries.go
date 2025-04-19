package handlers

import (
	"time"
)

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

// Result holds the data for a single entity returned by a query.
type Result struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`    // Content field that gets summarized
	Summary     string    `json:"summary"`    // Generated summary
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