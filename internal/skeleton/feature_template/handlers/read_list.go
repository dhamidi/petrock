package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

// HandleListItems handles requests to list items.
// Example route: GET /petrock_example_feature_name/
func (fs *FeatureServer) HandleListItems(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleListItems called", "feature", "petrock_example_feature_name")

	// Parse query parameters for filtering/pagination (example)
	page := ParseIntParam(r.URL.Query().Get("page"), 1)
	pageSize := ParseIntParam(r.URL.Query().Get("pageSize"), 20)
	filter := r.URL.Query().Get("filter")

	// Check for success message in query parameters
	successAction := r.URL.Query().Get("success")
	var successMsg string
	switch successAction {
	case "created":
		successMsg = "Item created successfully"
	case "updated":
		successMsg = "Item updated successfully"
	case "deleted":
		successMsg = "Item deleted successfully"
	}

	// Construct the query message
	query := ListQuery{
		Page:     page,
		PageSize: pageSize,
		Filter:   filter,
	}

	// Execute the query
	result, err := fs.querier.HandleList(r.Context(), query)
	if err != nil {
		slog.Error("Error handling ListQuery", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the list view HTML
	// Type assert the result to the appropriate type
	listResult, ok := result.(*ListResult)
	if !ok {
		slog.Error("Invalid result type for list view", "type", fmt.Sprintf("%T", result))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create page title
	pageTitle := "All Items"

	// Render the page with our helper and success message if present
	if err := RenderPageWithSuccess(w, pageTitle, ItemsListView(*listResult), successMsg); err != nil {
		slog.Error("Error rendering list view", "error", err)
		http.Error(w, "Error rendering view", http.StatusInternalServerError)
	}
}