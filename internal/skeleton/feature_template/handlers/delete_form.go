package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/petrock/example_module_path/petrock_example_feature_name/commands"
	"github.com/petrock/example_module_path/petrock_example_feature_name/queries"
)

// HandleDeleteForm handles requests to display a confirmation form for deleting an item.
// Example route: GET /petrock_example_feature_name/{id}/delete
func (fs *FeatureServer) HandleDeleteForm(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleDeleteForm called", "feature", "petrock_example_feature_name", "id", itemID)

	// Retrieve the item to confirm deletion
	query := queries.GetQuery{ID: itemID}
	result, err := fs.querier.HandleGet(r.Context(), query)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		slog.Error("Error retrieving item for delete confirmation", "error", err, "id", itemID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Cast the result to the correct type
	item, ok := result.(*Result)
	if !ok {
		slog.Error("Invalid result type for delete confirmation", "type", fmt.Sprintf("%T", result))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get CSRF token
	csrfToken := "token" // Replace with actual CSRF token generation

	// Render the delete confirmation view
	// Create page title
	pageTitle := fmt.Sprintf("Delete %s", item.Name)

	// Render the page with our helper
	if err := RenderPage(w, pageTitle, DeleteConfirmForm(item, csrfToken)); err != nil {
		slog.Error("Error rendering delete confirmation", "error", err)
		http.Error(w, "Error rendering confirmation", http.StatusInternalServerError)
	}
}

// HandleDeleteConfirm handles requests to delete an item after confirmation.
// Example route: POST /petrock_example_feature_name/{id}/delete
func (fs *FeatureServer) HandleDeleteConfirm(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleDeleteConfirm called", "feature", "petrock_example_feature_name", "id", itemID)

	// Parse form to get CSRF token (if needed)
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	// Verify CSRF token (if implemented)
	// ...

	// Create the delete command
	cmd := &commands.DeleteCommand{
		ID:        itemID,
		DeletedBy: "user", // Replace with actual user ID if authentication is implemented
		DeletedAt: time.Now().UTC(),
	}

	// Execute the command
	err := fs.app.Executor.Execute(r.Context(), cmd)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// Could redirect with a message that the item was already deleted
			w.Header().Set("Location", "/petrock_example_feature_name")
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		slog.Error("Failed to execute DeleteCommand", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set a success message in session (this would be implemented with a real session mechanism)
	// For now, we'll use a direct redirect, but in a real implementation you would:
	// session.SetFlash("success", "Item deleted successfully")

	// Redirect to the list view on success
	w.Header().Set("Location", "/petrock_example_feature_name?success=deleted")
	w.WriteHeader(http.StatusSeeOther) // 303 See Other
}
