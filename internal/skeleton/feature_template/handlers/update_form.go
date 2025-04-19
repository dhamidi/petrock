package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/petrock/example_module_path/core"
)

// HandleEditForm handles requests to display a form for editing an existing item.
// Example route: GET /petrock_example_feature_name/{id}/edit
func (fs *FeatureServer) HandleEditForm(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleEditForm called", "feature", "petrock_example_feature_name", "id", itemID)

	// Retrieve the item to edit
	query := GetQuery{ID: itemID}
	result, err := fs.querier.HandleGet(r.Context(), query)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		slog.Error("Error retrieving item for edit form", "error", err, "id", itemID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create an empty form
	form := core.NewForm(nil)

	// Get CSRF token
	csrfToken := "token" // Replace with actual CSRF token generation

	// Cast the result to the correct type
	item, ok := result.(*Result)
	if !ok {
		slog.Error("Invalid result type for edit form", "type", fmt.Sprintf("%T", result))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the edit form
	// Create page title
	pageTitle := fmt.Sprintf("Edit %s", item.Name)

	// Render the page with our helper
	if err := RenderPage(w, pageTitle, ItemForm(form, item, csrfToken)); err != nil {
		slog.Error("Error rendering edit form", "error", err)
		http.Error(w, "Error rendering form", http.StatusInternalServerError)
	}
}

// HandleUpdateForm handles requests to update an item from an edit form submission.
// Example route: POST /petrock_example_feature_name/{id}/edit
func (fs *FeatureServer) HandleUpdateForm(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleUpdateForm called", "feature", "petrock_example_feature_name", "id", itemID)

	// Parse the form
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	// Create a form instance with the parsed data
	form := core.NewForm(r.PostForm)

	// Validate required fields
	form.ValidateRequired("name", "description")

	// If the form has errors, re-render it with validation messages
	if !form.IsValid() {
		// Retrieve the original item to re-render the form with both the item and errors
		query := GetQuery{ID: itemID}
		result, err := fs.querier.HandleGet(r.Context(), query)
		if err != nil {
			slog.Error("Error retrieving item for form re-render", "error", err, "id", itemID)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Cast the result and render the form with errors
		item, ok := result.(*Result)
		if !ok {
			slog.Error("Invalid result type for edit form", "type", fmt.Sprintf("%T", result))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create page title for validation error
		pageTitle := fmt.Sprintf("Edit %s", item.Name)
		csrfToken := "token" // Replace with actual CSRF token

		// Render the page with validation errors
		if err := RenderPage(w, pageTitle, ItemForm(form, item, csrfToken)); err != nil {
			slog.Error("Error rendering form with validation errors", "error", err)
			http.Error(w, "Error rendering form", http.StatusInternalServerError)
		}
		return
	}

	// Create the update command from form data
	cmd := UpdateCommand{
		ID:          itemID,
		Name:        form.Get("name"),
		Description: form.Get("description"),
		UpdatedBy:   "user", // Replace with actual user ID if authentication is implemented
		UpdatedAt:   time.Now().UTC(),
	}

	// Execute the command
	err := fs.app.Executor.Execute(r.Context(), &cmd)
	if err != nil {
		// Check if it's a validation error
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "not found") {
			// Add the error to the form and re-render
			form.AddError("name", err.Error())

			// Retrieve the original item to re-render the form
			query := GetQuery{ID: itemID}
			result, getErr := fs.querier.HandleGet(r.Context(), query)
			if getErr != nil {
				slog.Error("Error retrieving item for form re-render", "error", getErr, "id", itemID)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Cast the result and render
			item, ok := result.(*Result)
			if !ok {
				slog.Error("Invalid result type for edit form", "type", fmt.Sprintf("%T", result))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Create page title for validation error
			pageTitle := fmt.Sprintf("Edit %s", item.Name)
			csrfToken := "token" // Replace with actual CSRF token

			// Render the page with validation errors
			if err := RenderPage(w, pageTitle, ItemForm(form, item, csrfToken)); err != nil {
				slog.Error("Error rendering form with validation errors", "error", err)
				http.Error(w, "Error rendering form", http.StatusInternalServerError)
			}
			return
		}

		// Handle other errors
		slog.Error("Failed to execute UpdateCommand", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set a success message in session (this would be implemented with a real session mechanism)
	// For now, we'll use a direct redirect, but in a real implementation you would:
	// session.SetFlash("success", "Item updated successfully")

	// Redirect to the item view on success
	w.Header().Set("Location", "/petrock_example_feature_name/"+itemID+"?success=updated")
	w.WriteHeader(http.StatusSeeOther) // 303 See Other
}
