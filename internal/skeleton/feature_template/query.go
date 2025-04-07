package petrock_example_feature_name

import (
	"context"
	"fmt"
	"log/slog"
	// "time" // Removed unused import

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Querier handles query processing for the feature.
// It typically depends on the feature's state representation.
type Querier struct {
	state *State // Example: Dependency on the feature's state
}

// NewQuerier creates a new Querier instance.
func NewQuerier(state *State) *Querier {
	return &Querier{
		state: state,
	}
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
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
		Version:     item.Version,
	}

	slog.Debug("Successfully processed GetQuery", "feature", "petrock_example_feature_name", "id", getQuery.ID)
	return result, nil
}

// HandleList processes the ListQuery.
// This function signature matches core.QueryHandler.
func (q *Querier) HandleList(ctx context.Context, query core.Query) (core.QueryResult, error) {
	listQuery, ok := query.(ListQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type for HandleList: expected ListQuery, got %T", query)
	}

	slog.Debug("Handling ListQuery", "feature", "petrock_example_feature_name", "page", listQuery.Page, "pageSize", listQuery.PageSize)

	if q.state == nil {
		slog.Error("State is nil in Querier, cannot handle ListQuery")
		return nil, fmt.Errorf("internal state not initialized")
	}

	// 1. Set defaults for pagination
	page := listQuery.Page
	if page < 1 {
		page = 1
	}
	pageSize := listQuery.PageSize
	if pageSize < 1 || pageSize > 100 { // Example max page size
		pageSize = 20 // Default page size
	}

	// 2. Retrieve items from state with filtering and pagination
	items, totalCount := q.state.ListItems(page, pageSize, listQuery.Filter)

	// 3. Map internal state items to QueryResult items
	results := make([]Result, 0, len(items))
	for _, item := range items {
		results = append(results, Result{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			Version:     item.Version,
		})
	}

	// 4. Construct the ListResult
	listResult := &ListResult{
		Items:      results,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	slog.Debug("Successfully processed ListQuery", "feature", "petrock_example_feature_name", "count", len(results), "total", totalCount)
	return listResult, nil
}

// Add more query handlers here...
