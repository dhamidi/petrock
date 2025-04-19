package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/petrock/example_module_path/petrock_example_feature_name/queries"
)

// HandleGetItem handles requests to retrieve a single item.
// Example route: GET /petrock_example_feature_name/{id}
func (fs *FeatureServer) HandleGetItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleGetItem called", "feature", "petrock_example_feature_name", "id", itemID)

	// Check for success message in query parameters
	successAction := r.URL.Query().Get("success")
	var successMsg string
	if successAction == "updated" {
		successMsg = "Item updated successfully"
	}

	// Construct the query message
	query := queries.GetQuery{ID: itemID}

	// Execute the query using the feature's querier
	result, err := fs.querier.HandleGet(r.Context(), query)
	if err != nil {
		// Handle not found error
		if strings.Contains(err.Error(), "not found") {
			slog.Warn("Item not found", "feature", "petrock_example_feature_name", "id", itemID)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		// Generic error handling for other errors
		slog.Error("Error handling GetQuery", "error", err, "id", itemID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the item view HTML
	// Type assert the result to the appropriate type
	getResult, ok := result.(*queries.GetQueryResult)
	if !ok {
		slog.Error("Invalid result type for item view", "type", fmt.Sprintf("%T", result))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create page title
	pageTitle := fmt.Sprintf("%s - Detail", getResult.Item.Name)

	// Render the page with our helper and success message if present
	if err := RenderPageWithSuccess(w, pageTitle, ItemView(getResult.Item), successMsg); err != nil {
		slog.Error("Error rendering item view", "error", err)
		http.Error(w, "Error rendering view", http.StatusInternalServerError)
	}
}
