package main

import (
	"context"
	"encoding/json" // Added for JSON handling in API endpoints
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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
func (s *AppState) Apply(msg interface{}) error {
	// In a real app, this would delegate to the appropriate feature state's Apply method
	// based on the message type.
	slog.Debug("AppState Apply called (placeholder)", "type", fmt.Sprintf("%T", msg))
	// Example delegation:
	// switch m := msg.(type) {
	// case posts.CreateCommand, posts.UpdateCommand, posts.DeleteCommand:
	//     return s.posts.Apply(m)
	// case users.RegisterCommand:
	//     return s.users.Apply(m)
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

	// --- Initialization ---
	slog.Info("Initializing application...")

	// 1. Initialize Core Registries
	commandRegistry := core.NewCommandRegistry()
	queryRegistry := core.NewQueryRegistry()
	slog.Debug("Initialized command and query registries")

	// 2. Initialize Encoder
	encoder := &core.JSONEncoder{} // Using JSON encoder
	slog.Debug("Initialized JSON encoder")

	// 3. Initialize Database Connection
	slog.Debug("Setting up database connection", "path", dbPath)
	db, err := core.SetupDatabase(dbPath)
	if err != nil {
		return fmt.Errorf("failed to setup database at %s: %w", dbPath, err)
	}
	defer func() {
		slog.Debug("Closing database connection", "path", dbPath)
		if err := db.Close(); err != nil {
			slog.Error("Error closing database", "path", dbPath, "error", err)
		}
	}()

	// 4. Initialize Message Log
	slog.Debug("Initializing message log")
	messageLog, err := core.NewMessageLog(db, encoder)
	if err != nil {
		// This also runs setupSchema, so errors are possible here
		return fmt.Errorf("failed to initialize message log: %w", err)
	}
	// Note: Message type registration (messageLog.RegisterType) will happen
	// within feature registration later.

	// 5. Initialize Application State
	slog.Debug("Initializing application state")
	appState := NewAppState() // Using the placeholder defined above

	// 6. Replay Message Log to Build State
	slog.Info("Replaying message log to build application state...")
	messages, err := messageLog.Load(context.Background())
	if err != nil {
		// Log loading errors can be critical, might indicate corruption
		return fmt.Errorf("failed to load messages from log for replay: %w", err)
	}
	slog.Debug("Loaded messages from log", "count", len(messages))
	replayErrors := 0
	for i, msg := range messages {
		// Apply each message to the state
		if err := appState.Apply(msg); err != nil {
			// Log errors during replay but continue if possible, depending on Apply logic
			slog.Error("Failed to apply message during replay", "error", err, "message_index", i, "message_type", fmt.Sprintf("%T", msg))
			replayErrors++
			// Decide whether to fail startup on replay errors. For now, just log.
		}
	}
	slog.Info("State replay completed", "message_count", len(messages), "replay_errors", replayErrors)
	if replayErrors > 0 {
		slog.Warn("Some messages failed to apply during state replay. State might be incomplete.")
	}

	// 7. Register Feature Handlers
	slog.Debug("Registering features...")
	RegisterAllFeatures(commandRegistry, queryRegistry, messageLog, appState) // Pass initialized components
	slog.Info("Features registered")

	// --- HTTP Server Setup ---
	slog.Info("Setting up HTTP server...")
	mux := http.NewServeMux()

	// Example: Setup middleware (logging, CSRF)
	// Note: Standard library middleware is often wrapped around specific handlers or the global mux.
	// var handler http.Handler = mux
	// handler = loggingMiddleware(handler)
	// handler = csrfMiddleware(handler) // CSRF often needs session state, which was removed. Re-evaluate CSRF strategy.

	// Example: Setup static file serving (if using embedded assets)
	// coreAssetsFS := core.GetAssetsFS() // Assuming core has embedded assets
	// mux.Handle("/assets/core/", http.StripPrefix("/assets/core/", http.FileServer(http.FS(coreAssetsFS))))
	// TODO: Add similar handlers for feature assets

	// Example: Define HTTP routes/handlers
	// Note: net/http mux uses pattern-based routing. For path parameters like /posts/{id},
	// you'd typically check r.URL.Path inside the handler or use a small helper/library.
	mux.HandleFunc("GET /", core.HandleIndex( /* queryRegistry */ )) // Pass dependencies - Use core.HandleIndex
	mux.HandleFunc("GET /commands", handleListCommands(commandRegistry))
	mux.HandleFunc("POST /commands", handleExecuteCommand(commandRegistry)) // Added route for executing commands
	// TODO: Add routes for features (e.g., mux.HandleFunc("GET /posts", handleListPosts), mux.HandleFunc("GET /posts/{id}", handleGetPost))
	// TODO: Add routes for other API endpoints (/queries, GET /queries/{name})

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

// handleExecuteCommand creates an http.HandlerFunc that decodes and dispatches commands.
func handleExecuteCommand(registry *core.CommandRegistry) http.HandlerFunc {
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

		// Look up the command type in the registry
		cmdType, found := registry.GetCommandType(req.Type)
		if !found {
			slog.Warn("Received request for unknown command type", "type", req.Type)
			http.Error(w, fmt.Sprintf("Bad Request: unknown command type %q", req.Type), http.StatusBadRequest) // Or 404? 400 seems better for unknown type name.
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
		cmdValue := reflect.ValueOf(cmdInstancePtr).Elem().Interface()

		// Dispatch the command
		slog.Debug("Dispatching command via API", "type", req.Type)
		dispatchErr := registry.Dispatch(r.Context(), cmdValue)

		if dispatchErr != nil {
			slog.Error("Error dispatching command", "type", req.Type, "error", dispatchErr)
			// TODO: Implement more specific error handling (e.g., validation errors -> 400)
			// For now, treat all dispatch errors as internal server errors.
			// If dispatchErr is a validation error type: http.Error(w, dispatchErr.Error(), http.StatusBadRequest); return
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Command successful
		slog.Info("Command executed successfully via API", "type", req.Type)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Or http.StatusAccepted (202) if processing is async
		// Optionally return a success body
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
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
