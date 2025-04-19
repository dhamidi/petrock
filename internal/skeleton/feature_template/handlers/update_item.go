package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// HandleUpdateItem handles requests to update an existing item.
// Example route: PUT /petrock_example_feature_name/{id}
func (fs *FeatureServer) HandleUpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleUpdateItem called", "feature", "petrock_example_feature_name", "id", itemID)

	var cmd UpdateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		slog.Error("Failed to decode UpdateCommand request body", "error", err, "id", itemID)
		http.Error(w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Ensure the ID in the path matches the ID in the body (important for PUT)
	if cmd.ID != itemID {
		http.Error(w, "Bad Request: Item ID in path does not match ID in request body", http.StatusBadRequest)
		return
	}

	// Execute the command using the central executor
	err := fs.app.Executor.Execute(r.Context(), &cmd)
	if err != nil {
		slog.Error("Failed to execute UpdateCommand", "error", err, "id", itemID)
		// Distinguish validation (e.g., not found, invalid name) from internal errors
		// TODO: Define specific validation error types or use error wrapping checks
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "not found") { // Simple check
			RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) // Or 404 if specifically "not found"
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	RespondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
