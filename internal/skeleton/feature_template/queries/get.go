package queries

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure query and result implement the marker interfaces
var _ core.Query = (*GetQuery)(nil)
var _ core.QueryResult = (*GetQueryResult)(nil)

// GetQuery holds data needed to retrieve a single entity.
type GetQuery struct {
	ID string // ID of the entity to retrieve
}

// QueryName returns the unique kebab-case name for this query type.
func (q GetQuery) QueryName() string {
	return "petrock_example_feature_name/get" // Removed suffix
}

// GetQueryResult wraps an ItemResult as a specific result type for GetQuery
type GetQueryResult struct {
	Item ItemResult `json:"item"`
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
	itemResult := ItemResult{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Content:     item.Content,
		Summary:     item.Summary,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		Version:     item.Version,
	}

	result := &GetQueryResult{
		Item: itemResult,
	}

	slog.Debug("Successfully processed GetQuery", "feature", "petrock_example_feature_name", "id", getQuery.ID)
	return result, nil
}
