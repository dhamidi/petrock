package handlers

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
	"github.com/petrock/example_module_path/petrock_example_feature_name/commands"
)

// HandleNewForm handles requests to display a form for creating a new item.
// Example route: GET /petrock_example_feature_name/new
func (fs *FeatureServer) HandleNewForm(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleNewForm called", "feature", "petrock_example_feature_name")

	// Create empty form data
	formData := ui.NewFormData(nil, nil)

	// Get CSRF token
	csrfToken := "token" // Replace with actual CSRF token generation

	// Render the form
	// Create page title
	pageTitle := "Create New Item"

	// Render the page with our helper
	if err := RenderPage(w, pageTitle, ItemForm(formData, nil, csrfToken)); err != nil {
		slog.Error("Error rendering new item form", "error", err)
		http.Error(w, "Error rendering form", http.StatusInternalServerError)
	}
}

// HandleCreateForm handles requests to create a new item from a form submission.
// Example route: POST /petrock_example_feature_name/new
func (fs *FeatureServer) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleCreateForm called", "feature", "petrock_example_feature_name", "method", r.Method, "url", r.URL.String())

	// Parse the form
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	slog.Debug("Form parsed successfully", "postForm", r.PostForm, "form", r.Form)

	// Create the command and parse from form data
	var cmd commands.CreateCommand
	slog.Debug("Attempting to parse command from form data", "postForm", r.PostForm)
	if err := core.ParseFromURLValues(r.PostForm, &cmd); err != nil {
		slog.Error("Failed to parse command from form data", "error", err, "postForm", r.PostForm)
		// Handle validation errors
		if parseErrors, ok := err.(*core.ParseErrors); ok {
			// Convert ParseErrors to ui.ParseError format
			var uiErrors []ui.ParseError
			for _, parseErr := range parseErrors.Errors {
				uiErrors = append(uiErrors, ui.ParseError{
					Field:   parseErr.Field,
					Message: parseErr.Message,
					Code:    parseErr.Code,
					Meta:    parseErr.Meta,
				})
			}

			// Create FormData with values and errors
			formData := ui.NewFormData(r.PostForm, uiErrors)

			// Render the form with validation errors
			pageTitle := "Create New Item"
			csrfToken := "token" // Replace with actual CSRF token

			if err := RenderPage(w, pageTitle, ItemForm(formData, nil, csrfToken)); err != nil {
				slog.Error("Error rendering form with validation errors", "error", err)
				http.Error(w, "Error rendering form", http.StatusInternalServerError)
			}
			return
		}

		// Handle other parsing errors
		slog.Error("Failed to parse form data", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Set additional fields not from form
	cmd.CreatedBy = "user" // Replace with actual user ID if authentication is implemented
	cmd.CreatedAt = time.Now().UTC()

	slog.Debug("Command parsed successfully", "command", cmd)

	// Execute the command
	slog.Debug("Executing command", "command", cmd)
	err := fs.app.Executor.Execute(r.Context(), &cmd)
	if err != nil {
		slog.Error("Command execution failed", "error", err, "command", cmd)
		// Check if it's a validation error
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "already exists") {
			// Create FormData with validation error
			uiErrors := []ui.ParseError{
				{
					Field:   "name",
					Message: err.Error(),
					Code:    "validation_error",
				},
			}
			formData := ui.NewFormData(r.PostForm, uiErrors)

			// Create page title for validation error
			pageTitle := "Create New Item"
			csrfToken := "token" // Replace with actual CSRF token

			// Render the page with validation errors
			if err := RenderPage(w, pageTitle, ItemForm(formData, nil, csrfToken)); err != nil {
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

	slog.Debug("Command executed successfully", "command", cmd)

	// Set a success message in session (this would be implemented with a real session mechanism)
	// For now, we'll use a direct redirect, but in a real implementation you would:
	// session.SetFlash("success", "Item created successfully")

	// Redirect to the list view on success
	redirectURL := "/petrock_example_feature_name?success=created"
	slog.Debug("Redirecting to list view", "url", redirectURL)
	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusSeeOther) // 303 See Other
}
