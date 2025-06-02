package queries

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure query and result implement the marker interfaces
var _ core.Query = (*ListQuery)(nil)
var _ core.QueryResult = (*ListQueryResult)(nil)

// ListQuery holds data needed to retrieve a list of entities, possibly filtered or paginated.
type ListQuery struct {
	Page     int    `json:"page" validate:"min=1"`           // For pagination
	PageSize int    `json:"page_size" validate:"min=1,max=100"` // For pagination
	Filter   string `json:"filter" validate:"maxlen=100"`    // Example filter criteria
}

// QueryName returns the unique kebab-case name for this query type.
func (q ListQuery) QueryName() string {
	return "petrock_example_feature_name/list" // Removed suffix
}

// ListQueryResult holds a list of entities and pagination details.
type ListQueryResult struct {
	Items      []ItemResult `json:"items"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
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
	results := make([]ItemResult, 0, len(items))
	for _, item := range items {
		results = append(results, ItemResult{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Content:     item.Content,
			Summary:     item.Summary,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			Version:     item.Version,
		})
	}

	// 4. Construct the ListQueryResult
	listResult := &ListQueryResult{
		Items:      results,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	slog.Debug("Successfully processed ListQuery", "feature", "petrock_example_feature_name", "count", len(results), "total", totalCount)
	return listResult, nil
}
