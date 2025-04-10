package main

import (
	"context"
	"encoding/json" // Added for JSON handling in API endpoints
	"fmt"
	"log/slog"
	"net/http"
	"net/url" // Added for parsing query parameters
	"os"
	"os/signal"
	"reflect" // Added for command/query execution handlers
	"strconv" // Added for converting query parameters
	"strings" // Added for query parameter population helper
	"syscall"
	"time"

	"github.com/petrock/example_module_path/core" // Assuming core package exists

	// Use standard library for routing
	"github.com/spf13/cobra"
)

// AppState is a placeholder for the application's aggregated state.
// In a real application, this might be composed of states from different features.
type AppState struct {
	// In a more complex app, this might hold pointers to feature-specific states:
	// posts *posts.State
	// users *users.State
}

// NewAppState creates a new AppState.
func NewAppState() *AppState {
	return &AppState{}
}

// Apply processes a message (typically a command or event) to update the state.
// This is crucial for rebuilding state from the message log on startup.
// The msgMeta parameter is non-nil during replay, providing access to timestamp and ID.
func (s *AppState) Apply(msg interface{}, msgMeta *core.Message) error {
	// In a real app, this would delegate to the appropriate feature state's Apply method
	// based on the message type.
	slog.Debug("AppState Apply called (placeholder)", "type", fmt.Sprintf("%T", msg), "hasMeta", msgMeta != nil)
	// Example delegation:
	// switch m := msg.(type) {
	// case posts.CreateCommand, posts.UpdateCommand, posts.DeleteCommand:
	//     return s.posts.Apply(m, msgMeta)
	// case users.RegisterCommand:
	//     return s.users.Apply(m, msgMeta)
	// default:
	//     slog.Warn("AppState.Apply received unhandled message type", "type", fmt.Sprintf("%T", msg))
	// }
	return nil
}

// NewServeCmd creates the `serve` subcommand
func NewServeCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Starts the HTTP server",
		Long:  `Starts the HTTP server to handle web requests.`,
		RunE:  runServe,
	}

	// Add flags like --port, --host
	serveCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	serveCmd.Flags().String("host", "localhost", "Host to bind to")
	serveCmd.Flags().String("db-path", "app.db", "Path to the SQLite database file") // Added db-path flag

	return serveCmd
}

func runServe(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	dbPath, _ := cmd.Flags().GetString("db-path") // Get db-path flag value
	addr := fmt.Sprintf("%s:%d", host, port)

	// --- Initialization using core.App ---
	// Initialize the application using the new core.App struct
	app, err := core.NewApp(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer app.Close() // Ensure resources are closed properly

	// Initialize Application State
	slog.Debug("Initializing application state")
	appState := NewAppState() // Using the placeholder defined above

	// Register features BEFORE replaying the log
	// Pass all necessary dependencies through RegisterAllFeatures
	RegisterAllFeatures(nil, app.CommandRegistry, app.QueryRegistry, app.MessageLog, app.Executor, appState, app.DB)

	// Replay the message log to build application state
	if err := app.ReplayLog(); err != nil {
		return fmt.Errorf("failed to replay message log: %w", err)
	}

	// --- HTTP Server Setup ---
	mux := http.NewServeMux()

	// Example: Setup middleware (logging, CSRF)
	// Note: Standard library middleware is often wrapped around specific handlers or the global mux.
	// var handler http.Handler = mux
	// handler = loggingMiddleware(handler)
	// handler = csrfMiddleware(handler) // CSRF often needs session state, which was removed. Re-evaluate CSRF strategy.

	// Example: Setup static file serving (if using embedded assets)
	// coreAssetsFS := core.GetAssetsFS() // Assuming core has embedded assets
	// mux.Handle("/assets/core/", http.StripPrefix("/assets/core/", http.FileServer(http.FS(coreAssetsFS))))
	// TODO: Add handlers for feature assets

	// Setup core HTTP routes
	slog.Info("Setting up HTTP server routes...")
	
	// Setup index route
	mux.HandleFunc("GET /", core.HandleIndex(app.CommandRegistry, app.QueryRegistry))
	
	// Setup core API routes
	mux.HandleFunc("GET /commands", handleListCommands(app.CommandRegistry))
	mux.HandleFunc("POST /commands", handleExecuteCommand(app.Executor, app.CommandRegistry))
	mux.HandleFunc("GET /queries", handleListQueries(app.QueryRegistry))
	mux.HandleFunc("GET /queries/{feature}/{queryName}", handleExecuteQuery(app.QueryRegistry))

	// Register feature-specific handlers and routes
	// The mux is needed here to register HTTP routes
	// Note: Features were already registered above for command/query handling
	RegisterFeatureRoutes(mux, appState, app.Executor) // A new function to register only routes

	// --- Server Start and Shutdown ---
	server := &http.Server{
		Addr:         addr,
		Handler:      mux, // Use the configured mux (potentially wrapped in middleware)
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel to listen for errors starting the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		slog.Info("Starting server", "address", addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt or terminate signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		slog.Info("Shutdown signal received", "signal", sig)

		// Graceful shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Attempt to gracefully shut down the server
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
			// Force close if shutdown fails
			if closeErr := server.Close(); closeErr != nil {
				slog.Error("Failed to close server", "error", closeErr)
			}
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		slog.Info("Server gracefully stopped")
	}

	return nil
}

// handleListCommands creates an http.HandlerFunc that lists registered command types.
func handleListCommands(registry *core.CommandRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		commandNames := registry.RegisteredCommandNames()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(commandNames); err != nil {
			slog.Error("Failed to encode command names list", "error", err)
			// Hard to send error to client if header already written, but log it.
		}
	}
}

// commandRequest is used to decode the incoming JSON payload for command execution.
type commandRequest struct {
	Type    string          `json:"type"`    // The registered name of the command type
	Payload json.RawMessage `json:"payload"` // The command-specific data
}

// handleExecuteCommand creates an http.HandlerFunc that decodes and executes commands
// using the central core.Executor.
func handleExecuteCommand(executor *core.Executor, registry *core.CommandRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Unsupported Media Type: Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		// Decode the request body into the intermediate struct
		var req commandRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Prevent unexpected fields
		err := decoder.Decode(&req)
		if err != nil {
			slog.Error("Failed to decode command request body", "error", err)
			http.Error(w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate type name presence
		if req.Type == "" {
			http.Error(w, "Bad Request: 'type' field is required", http.StatusBadRequest)
			return
		}

		// Look up the command type in the registry using the kebab-case name
		cmdType, found := registry.GetCommandType(req.Type) // req.Type should be "feature/command-name"
		if !found {
			slog.Warn("Received request for unknown command type", "name", req.Type)
			http.Error(w, fmt.Sprintf("Bad Request: unknown command type %q", req.Type), http.StatusBadRequest)
			return
		}

		// Create a new instance of the command struct (must be a pointer for unmarshaling)
		cmdInstancePtr := reflect.New(cmdType).Interface()

		// Unmarshal the payload into the command instance pointer
		if err := json.Unmarshal(req.Payload, cmdInstancePtr); err != nil {
			slog.Error("Failed to unmarshal command payload", "type", req.Type, "error", err)
			http.Error(w, fmt.Sprintf("Bad Request: invalid payload for type %q: %s", req.Type, err.Error()), http.StatusBadRequest)
			return
		}

		// Get the actual command value (dereferenced) to pass to Dispatch
		// Ensure the command instance implements core.Command (which includes CommandName)
		cmdValue, ok := reflect.ValueOf(cmdInstancePtr).Elem().Interface().(core.Command)
		if !ok {
			// Defensive check
			slog.Error("Internal error: command instance does not implement core.Command", "name", req.Type, "type", reflect.TypeOf(cmdInstancePtr).Elem())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the command using the central executor
		slog.Debug("Executing command via API", "name", req.Type)
		execErr := executor.Execute(r.Context(), cmdValue) // Use executor.Execute

		if execErr != nil {
			slog.Error("Error executing command", "name", req.Type, "error", execErr)
			// Handle validation errors vs. other errors
			// Example: Check if the error is a validation error (you might need to define custom error types or check wrapped errors)
			// if errors.As(execErr, &core.ValidationError{}) { // Assuming a ValidationError type
			//     respondJSON(w, http.StatusBadRequest, map[string]string{"error": execErr.Error()})
			// } else {
			// Treat other errors (logging failure, etc.) as internal server errors
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// }
			return
		}

		// Command successful
		slog.Info("Command executed successfully via API", "name", req.Type) // Log the full name
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Or http.StatusAccepted (202) if processing is async
		// Optionally return a success body
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

// handleListQueries creates an http.HandlerFunc that lists registered query types.
func handleListQueries(registry *core.QueryRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		queryNames := registry.RegisteredQueryNames()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(queryNames); err != nil {
			slog.Error("Failed to encode query names list", "error", err)
			// Hard to send error to client if header already written, but log it.
		}
	}
}

// handleExecuteQuery creates an http.HandlerFunc that executes queries based on URL path and parameters.
func handleExecuteQuery(registry *core.QueryRegistry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract feature and kebab-case query name from path (requires Go 1.22+)
		featurePart := r.PathValue("feature")
		queryNamePart := r.PathValue("queryName") // This is now kebab-case
		if featurePart == "" || queryNamePart == "" {
			http.Error(w, "Bad Request: URL path must be /queries/{feature}/{query-name}", http.StatusBadRequest)
			return
		}
		fullQueryName := fmt.Sprintf("%s/%s", featurePart, queryNamePart) // The full kebab-case name
		slog.Debug("Handling query request via API", "name", fullQueryName)

		// Look up the query type in the registry using the full name
		queryType, found := registry.GetQueryType(fullQueryName)
		if !found {
			slog.Warn("Received request for unknown query type", "name", fullQueryName)
			http.Error(w, fmt.Sprintf("Not Found: unknown query type %q", fullQueryName), http.StatusNotFound)
			return
		}

		// Create a new instance of the query struct (must be a pointer for reflection)
		queryInstancePtr := reflect.New(queryType) // Returns a pointer Value
		queryInstance := queryInstancePtr.Elem()   // Get the struct Value

		// Populate the query struct fields from URL query parameters
		urlParams := r.URL.Query()
		if err := populateStructFromURLParams(queryInstance, urlParams); err != nil {
			slog.Error("Failed to populate query struct from URL parameters", "name", fullQueryName, "error", err) // Use fullQueryName here
			http.Error(w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Get the actual query value (non-pointer) to pass to Dispatch
		// Ensure the query instance implements core.Query (which includes QueryName)
		queryValue, ok := queryInstance.Interface().(core.Query)
		if !ok {
			// Defensive check
			slog.Error("Internal error: query instance does not implement core.Query", "name", fullQueryName, "type", queryInstance.Type())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Dispatch the query
		slog.Debug("Dispatching query via API", "name", fullQueryName)
		result, dispatchErr := registry.Dispatch(r.Context(), queryValue)

		if dispatchErr != nil {
			slog.Error("Error dispatching query", "name", fullQueryName, "error", dispatchErr)
			// TODO: Implement more specific error handling (e.g., ErrNotFound -> 404)
			// if errors.Is(dispatchErr, core.ErrNotFound) {
			// 	http.Error(w, "Not Found", http.StatusNotFound)
			// 	return
			// }
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Query successful
		slog.Info("Query executed successfully via API", "name", fullQueryName)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			slog.Error("Failed to encode query result", "name", fullQueryName, "error", err)
			// Hard to send error to client if header already written, but log it.
		}
	}
}

// populateStructFromURLParams uses reflection to set fields of a struct
// based on values found in URL query parameters.
// It supports string, int, and bool field types.
func populateStructFromURLParams(structVal reflect.Value, params url.Values) error {
	if structVal.Kind() != reflect.Struct {
		return fmt.Errorf("internal error: expected a struct value, got %s", structVal.Kind())
	}
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		fieldVal := structVal.Field(i)
		fieldType := structType.Field(i)
		fieldName := fieldType.Name // Use struct field name directly

		// Check if field is settable (exported)
		if !fieldVal.CanSet() {
			continue
		}

		// Get parameter value (case-sensitive match with field name)
		paramValueStr, exists := params[fieldName]
		if !exists || len(paramValueStr) == 0 {
			// Also check lowercase version for convenience? Optional.
			lowerFieldName := strings.ToLower(fieldName)
			paramValueStr, exists = params[lowerFieldName]
			if !exists || len(paramValueStr) == 0 {
				continue // No parameter found for this field
			}
		}
		valueStr := paramValueStr[0] // Use the first value if multiple are provided

		// Set field value based on its type
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(valueStr)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer value %q for field %q: %w", valueStr, fieldName, err)
			}
			if fieldVal.OverflowInt(intValue) {
				return fmt.Errorf("integer value %q overflows field %q", valueStr, fieldName)
			}
			fieldVal.SetInt(intValue)
		case reflect.Bool:
			// Handle common boolean representations: "true", "false", "1", "0"
			boolValue, err := strconv.ParseBool(strings.ToLower(valueStr))
			if err != nil && valueStr != "1" && valueStr != "0" { // Allow 1/0 as bool
				return fmt.Errorf("invalid boolean value %q for field %q", valueStr, fieldName)
			}
			if valueStr == "1" { // Handle "1" explicitly if ParseBool fails
				boolValue = true
			}
			fieldVal.SetBool(boolValue)
		// Add cases for other supported types (float, etc.) if needed
		default:
			slog.Warn("Unsupported field type for URL parameter population", "field", fieldName, "type", fieldVal.Kind())
		}
	}
	return nil
}

// --- Placeholder Middleware/Handlers (Replace with actual implementations) ---

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		slog.Info("Request received", "method", r.Method, "path", r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func csrfMiddleware(next http.Handler) http.Handler {
// 	// TODO: Implement CSRF protection (e.g., using standard library techniques or other allowed libraries)
// 	// Re-evaluate CSRF strategy as session state (often used by CSRF libraries) was removed.
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Add CSRF token check logic here
// 		next.ServeHTTP(w, r)
// 	})
// }

// Session middleware removed as gorilla/sessions is not allowed.
// Consider alternative session management if needed (e.g., client-side tokens, other libraries).

// func HandleIndex(/* queryRegistry *core.QueryRegistry */) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Example: Render index page
// 		// component := core.IndexPage() // Get component from core/page_index.go
// 		// layout := core.Layout("Home", component) // Wrap in layout
// 		// layout.Render(w) // Render component
// 		fmt.Fprintln(w, "Welcome to petrock_example_project_name!") // Placeholder
// 	}
// }
