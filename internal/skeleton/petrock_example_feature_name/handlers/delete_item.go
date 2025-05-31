package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/petrock/example_module_path/petrock_example_feature_name/commands"
)

// HandleDeleteItem handles requests to delete an item.
// Example route: DELETE /petrock_example_feature_name/{id}
func (fs *FeatureServer) HandleDeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleDeleteItem called", "feature", "petrock_example_feature_name", "id", itemID)

	// Construct the command
	cmd := &commands.DeleteCommand{ID: itemID /* DeletedBy: "user_from_context" */}

	// Execute the command using the central executor
	err := fs.app.Executor.Execute(r.Context(), cmd)
	if err != nil {
		slog.Error("Failed to execute DeleteCommand", "error", err, "id", itemID)
		// Distinguish validation (e.g., not found) from internal errors
		// TODO: Define specific validation error types or use error wrapping checks
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "not found") { // Simple check
			// If the validation allows deleting non-existent items idempotently, Execute might return nil.
			// If validation returns "not found", return 404 or 400.
			RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) // Or 404
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	RespondJSON(w, http.StatusOK, map[string]string{"status": "deleted"}) // Or 204 No Content
}
