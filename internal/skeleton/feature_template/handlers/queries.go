package handlers

import (
	"time"
)

// Moving query types to queries package

// Result holds the data for a single entity returned by a query.
type Result struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"` // Content field that gets summarized
	Summary     string    `json:"summary"` // Generated summary
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
