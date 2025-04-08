package petrock_example_feature_name

import (
	"database/sql" // Example: If handlers need direct DB access
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// FeatureServer holds dependencies required by the feature's HTTP handlers.
// This struct is initialized in register.go and passed to RegisterRoutes.
type FeatureServer struct {
	executor *Executor             // Command execution logic
	querier  *Querier              // Query execution logic
	state    *State                // Direct state access (use querier/executor preferably)
	log      *core.MessageLog      // For logging commands/events directly
	commands *core.CommandRegistry // For dispatching commands via the core system
	db       *sql.DB               // Example: Shared DB connection pool
	// Add other dependencies like config, template renderers, etc.
}

// NewFeatureServer creates and initializes the FeatureServer with its dependencies.
func NewFeatureServer(
	executor *Executor,
	querier *Querier,
	state *State,
	log *core.MessageLog,
	commands *core.CommandRegistry,
	db *sql.DB, // Add other dependencies here
) *FeatureServer {
	// Basic validation
	if executor == nil || querier == nil || state == nil || log == nil || commands == nil {
		// Depending on requirements, some dependencies might be optional (e.g., db)
		panic("missing required dependencies for FeatureServer")
	}
	return &FeatureServer{
		executor: executor,
		querier:  querier,
		state:    state,
		log:      log,
		commands: commands,
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
		// Handle specific errors, e.g., not found
		// if errors.Is(err, core.ErrNotFound) { // Assuming a standard error type
		// 	http.Error(w, "Not Found", http.StatusNotFound)
		// 	return
		// }
		// Generic error handling
		slog.Error("Error handling GetQuery", "error", err, "id", itemID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with the result (e.g., as JSON)
	respondJSON(w, http.StatusOK, result)
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

	// Respond with the result
	respondJSON(w, http.StatusOK, result)
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

	// --- State Change Strategy ---
	// Option A (Recommended): Log the command directly using MessageLog.
	// The state's Apply method will handle the actual state update, triggered by log replay or explicitly.
	err := fs.log.Append(r.Context(), cmd)
	if err != nil {
		// Handle potential errors (validation errors should ideally be caught before logging)
		slog.Error("Failed to append CreateCommand via feature handler", "error", err)
		// Check if the error is a validation error type and return 400 if so
		// if validationErr, ok := err.(*core.ValidationError); ok {
		//     respondJSON(w, http.StatusBadRequest, map[string]string{"error": validationErr.Error()})
		//     return
		// }
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Option B: Dispatch via CommandRegistry (handles logging and state application via registered handler)
	// err := fs.commands.Dispatch(r.Context(), cmd)
	// if err != nil {
	//     slog.Error("Failed to dispatch CreateCommand via feature handler", "error", err)
	//     // Handle validation vs internal errors appropriately
	//     http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//     return
	// }

	// Respond with success (e.g., 201 Created or 202 Accepted)
	// Optionally return the created resource or its ID
	respondJSON(w, http.StatusCreated, map[string]string{"status": "created"}) // Simple status response
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

	// Log the command (Option A)
	err := fs.log.Append(r.Context(), cmd)
	if err != nil {
		slog.Error("Failed to append UpdateCommand via feature handler", "error", err, "id", itemID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	// Log the command (Option A)
	err := fs.log.Append(r.Context(), cmd)
	if err != nil {
		// Handle not found errors if Append checks existence first, or let Apply handle it
		slog.Error("Failed to append DeleteCommand via feature handler", "error", err, "id", itemID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

// Add more handlers as needed...
