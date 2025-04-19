package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// HandleNewForm handles requests to display a form for creating a new item.
// Example route: GET /petrock_example_feature_name/new
func (fs *FeatureServer) HandleNewForm(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleNewForm called", "feature", "petrock_example_feature_name")

	// Create an empty form
	form := core.NewForm(nil)

	// Get CSRF token
	csrfToken := "token" // Replace with actual CSRF token generation

	// Render the form
	// Create page title
	pageTitle := "Create New Item"

	// Render the page with our helper
	if err := RenderPage(w, pageTitle, ItemForm(form, nil, csrfToken)); err != nil {
		slog.Error("Error rendering new item form", "error", err)
		http.Error(w, "Error rendering form", http.StatusInternalServerError)
	}
}

// HandleCreateForm handles requests to create a new item from a form submission.
// Example route: POST /petrock_example_feature_name/new
func (fs *FeatureServer) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleCreateForm called", "feature", "petrock_example_feature_name")

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
		// Create page title for validation error
		pageTitle := "Create New Item"
		csrfToken := "token" // Replace with actual CSRF token

		// Render the page with validation errors
		if err := RenderPage(w, pageTitle, ItemForm(form, nil, csrfToken)); err != nil {
			slog.Error("Error rendering form with validation errors", "error", err)
			http.Error(w, "Error rendering form", http.StatusInternalServerError)
		}
		return
	}

	// Create the command from form data
	cmd := CreateCommand{
		Name:        form.Get("name"),
		Description: form.Get("description"),
		CreatedBy:   "user", // Replace with actual user ID if authentication is implemented
		CreatedAt:   time.Now().UTC(),
	}

	// Execute the command
	err := fs.app.Executor.Execute(r.Context(), &cmd)
	if err != nil {
		// Check if it's a validation error
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "already exists") {
			// Add the error to the form and re-render
			form.AddError("name", err.Error())

			// Create page title for validation error
			pageTitle := "Create New Item"
			csrfToken := "token" // Replace with actual CSRF token

			// Render the page with validation errors
			if err := RenderPage(w, pageTitle, ItemForm(form, nil, csrfToken)); err != nil {
				slog.Error("Error rendering form with validation errors", "error", err)
				http.Error(w, "Error rendering form", http.StatusInternalServerError)
			}
			return
		}

		// Handle other errors
		slog.Error("Failed to execute CreateCommand", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set a success message in session (this would be implemented with a real session mechanism)
	// For now, we'll use a direct redirect, but in a real implementation you would:
	// session.SetFlash("success", "Item created successfully")

	// Redirect to the list view on success
	w.Header().Set("Location", "/petrock_example_feature_name?success=created")
	w.WriteHeader(http.StatusSeeOther) // 303 See Other
}