package petrock_example_feature_name

import (
	"database/sql" // Example: If handlers need direct DB access
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv" // Added for parseIntParam helper
	"strings"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// FeatureServer holds dependencies required by the feature's HTTP handlers.
// This struct is initialized in register.go and passed to RegisterRoutes.
type FeatureServer struct {
	executor *core.Executor   // Central command executor
	querier  *Querier         // Query execution logic
	state    *State           // Direct state access (use querier/executor preferably)
	log      *core.MessageLog // For logging commands/events directly (less common now)
	db       *sql.DB          // Example: Shared DB connection pool
	// Add other dependencies like config, template renderers, etc.
}

// NewFeatureServer creates and initializes the FeatureServer with its dependencies.
// Note: It now receives the central core.Executor.
func NewFeatureServer(
	executor *core.Executor, // Changed from feature executor to core executor
	querier *Querier,
	state *State,
	log *core.MessageLog, // Keep log if needed for other purposes, but not primary command path
	db *sql.DB, // Add other dependencies here
) *FeatureServer {
	// Basic validation
	if executor == nil || querier == nil || state == nil || log == nil {
		// Depending on requirements, some dependencies might be optional (e.g., db, log)
		panic("missing required dependencies for FeatureServer")
	}
	return &FeatureServer{
		executor: executor, // Store the central executor
		querier:  querier,
		state:    state,
		log:      log,
		db:       db,
	}
}

// --- Handler Methods ---
// These methods are attached to FeatureServer and registered in routes.go.

// HandleGetItem handles requests to retrieve a single item.
// Example route: GET /petrock_example_feature_name/{id}
func (fs *FeatureServer) HandleGetItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("id") // Requires Go 1.22+
	if itemID == "" {
		http.Error(w, "Bad Request: Missing item ID in path", http.StatusBadRequest)
		return
	}
	slog.Debug("HandleGetItem called", "feature", "petrock_example_feature_name", "id", itemID)

	// Construct the query message
	query := GetQuery{ID: itemID}

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ItemView(result).Render(w); err != nil {
		slog.Error("Error rendering item view", "error", err)
		http.Error(w, "Error rendering view", http.StatusInternalServerError)
	}
}

// HandleListItems handles requests to list items.
// Example route: GET /petrock_example_feature_name/
func (fs *FeatureServer) HandleListItems(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleListItems called", "feature", "petrock_example_feature_name")

	// Parse query parameters for filtering/pagination (example)
	page := parseIntParam(r.URL.Query().Get("page"), 1)
	pageSize := parseIntParam(r.URL.Query().Get("pageSize"), 20)
	filter := r.URL.Query().Get("filter")

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ItemsListView(result.(ListResult)).Render(w); err != nil {
		slog.Error("Error rendering list view", "error", err)
		http.Error(w, "Error rendering view", http.StatusInternalServerError)
	}
}

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
	err := fs.executor.Execute(r.Context(), cmd)
	if err != nil {
		slog.Error("Failed to execute CreateCommand", "error", err)
		// Distinguish between validation errors (client-side, 400) and other errors (server-side, 500)
		// This requires core.Executor.Execute to wrap validation errors or use specific types.
		// Assuming validation errors are returned directly or wrapped:
		// TODO: Define specific validation error types or use error wrapping checks
		// Example check (adjust based on actual error handling strategy):
		if strings.Contains(err.Error(), "validation failed") { // Simple string check, improve this
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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
	respondJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

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
	err := fs.executor.Execute(r.Context(), cmd)
	if err != nil {
		slog.Error("Failed to execute UpdateCommand", "error", err, "id", itemID)
		// Distinguish validation (e.g., not found, invalid name) from internal errors
		// TODO: Define specific validation error types or use error wrapping checks
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "not found") { // Simple check
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) // Or 404 if specifically "not found"
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	respondJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

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
	cmd := DeleteCommand{ID: itemID /* DeletedBy: "user_from_context" */}

	// Execute the command using the central executor
	err := fs.executor.Execute(r.Context(), cmd)
	if err != nil {
		slog.Error("Failed to execute DeleteCommand", "error", err, "id", itemID)
		// Distinguish validation (e.g., not found) from internal errors
		// TODO: Define specific validation error types or use error wrapping checks
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "not found") { // Simple check
			// If the validation allows deleting non-existent items idempotently, Execute might return nil.
			// If validation returns "not found", return 404 or 400.
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) // Or 404
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"}) // Or 204 No Content
}

// HandleCustomIndex is an example of a feature overriding a core route.
// Example route: GET /
// func (fs *FeatureServer) HandleCustomIndex(w http.ResponseWriter, r *http.Request) {
// 	slog.Debug("HandleCustomIndex called", "feature", "petrock_example_feature_name")
// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
// 	fmt.Fprintf(w, "Hello from the %s feature's custom index page!", "petrock_example_feature_name")
// }

// --- Helper Functions ---

// respondJSON is a utility to send JSON responses.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			// Log error, but can't change status code now
			slog.Error("Failed to encode JSON response", "error", err)
		}
	}
}

// parseIntParam is a helper to parse integer query parameters with a default value.
func parseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}
	return val
}

// HandleNewForm handles requests to display a form for creating a new item.
// Example route: GET /petrock_example_feature_name/new
func (fs *FeatureServer) HandleNewForm(w http.ResponseWriter, r *http.Request) {
	slog.Debug("HandleNewForm called", "feature", "petrock_example_feature_name")

	// Create an empty form
	form := core.NewForm(nil)

	// Get CSRF token
	csrfToken := "token" // Replace with actual CSRF token generation

	// Render the form
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ItemForm(form, nil, csrfToken).Render(w); err != nil {
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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		csrfToken := "token" // Replace with actual CSRF token
		if err := ItemForm(form, nil, csrfToken).Render(w); err != nil {
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
	err := fs.executor.Execute(r.Context(), cmd)
	if err != nil {
		// Check if it's a validation error
		if strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "already exists") {
			// Add the error to the form and re-render
			form.AddError("name", err.Error())
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			csrfToken := "token" // Replace with actual CSRF token
			if err := ItemForm(form, nil, csrfToken).Render(w); err != nil {
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

	// Redirect to the list view on success
	w.Header().Set("Location", "/petrock_example_feature_name")
	w.WriteHeader(http.StatusSeeOther) // 303 See Other
}

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := ItemForm(form, item, csrfToken).Render(w); err != nil {
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

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		csrfToken := "token" // Replace with actual CSRF token
		if err := ItemForm(form, item, csrfToken).Render(w); err != nil {
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
	err := fs.executor.Execute(r.Context(), cmd)
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

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			csrfToken := "token" // Replace with actual CSRF token
			if err := ItemForm(form, item, csrfToken).Render(w); err != nil {
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

	// Redirect to the item view on success
	w.Header().Set("Location", "/petrock_example_feature_name/"+itemID)
	w.WriteHeader(http.StatusSeeOther) // 303 See Other
}

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
	query := GetQuery{ID: itemID}
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := DeleteConfirmForm(item, csrfToken).Render(w); err != nil {
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
	cmd := DeleteCommand{
		ID:        itemID,
		DeletedBy: "user", // Replace with actual user ID if authentication is implemented
		DeletedAt: time.Now().UTC(),
	}

	// Execute the command
	err := fs.executor.Execute(r.Context(), cmd)
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

	// Redirect to the list view on success
	w.Header().Set("Location", "/petrock_example_feature_name")
	w.WriteHeader(http.StatusSeeOther) // 303 See Other
}

// Add more handlers as needed...
