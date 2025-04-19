package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// HandleCreateItem handles requests to create a new item.
// Example route: POST /petrock_example_feature_name/
func (fs *FeatureServer) HandleCreateItem(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleCreateItem called", "feature", "petrock_example_feature_name")

	// Decode request body into the command struct
	var cmd CreateCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		slog.Error("Failed to decode CreateCommand request body", "error", err)
		http.Error(w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Execute the command using the central executor
	err := fs.app.Executor.Execute(r.Context(), &cmd)
	if err != nil {
		slog.Error("Failed to execute CreateCommand", "error", err)
		// Distinguish between validation errors (client-side, 400) and other errors (server-side, 500)
		// This requires core.Executor.Execute to wrap validation errors or use specific types.
		// Assuming validation errors are returned directly or wrapped:
		// TODO: Define specific validation error types or use error wrapping checks
		// Example check (adjust based on actual error handling strategy):
		if strings.Contains(err.Error(), "validation failed") { // Simple string check, improve this
			RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with success (e.g., 201 Created or 202 Accepted)
	// Optionally return the created resource or its ID
	// If the command generated an ID, it might be available in the cmd struct *after* Apply,
	// but Apply runs *after* Execute returns nil. Need a way to get the ID if needed.
	// For now, just return status.
	RespondJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}