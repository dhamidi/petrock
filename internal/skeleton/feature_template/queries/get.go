package queries

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// GetQuery holds data needed to retrieve a single entity.
type GetQuery struct {
	ID string // ID of the entity to retrieve
}

// QueryName returns the unique kebab-case name for this query type.
func (q GetQuery) QueryName() string {
	return "petrock_example_feature_name/get" // Removed suffix
}

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

// HandleGet processes the GetQuery.
// This function signature matches core.QueryHandler.
func (q *Querier) HandleGet(ctx context.Context, query core.Query) (core.QueryResult, error) {
	getQuery, ok := query.(GetQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for HandleGet: expected GetQuery, got %T", query)
	}

	slog.Debug("Handling GetQuery", "feature", "petrock_example_feature_name", "id", getQuery.ID)

	if q.state == nil {
		slog.Error("State is nil in Querier, cannot handle GetQuery")
		return nil, fmt.Errorf("internal state not initialized")
	}

	// 1. Retrieve item from state
	item, found := q.state.GetItem(getQuery.ID)
	if !found {
		// Return a specific "not found" error if defined, otherwise a generic one
		return nil, fmt.Errorf("item with ID %s not found", getQuery.ID)
	}

	// 2. Map internal state representation to the QueryResult struct
	result := &Result{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Content:     item.Content,
		Summary:     item.Summary,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		Version:     item.Version,
	}

	slog.Debug("Successfully processed GetQuery", "feature", "petrock_example_feature_name", "id", getQuery.ID)
	return result, nil
}
