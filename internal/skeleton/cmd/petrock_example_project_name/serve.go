package main

import (
	"context"
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
	// TODO: Add flags for database path, etc.

	return serveCmd
}

func runServe(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")
	host, _ := cmd.Flags().GetString("host")
	addr := fmt.Sprintf("%s:%d", host, port)

	// --- Initialization ---
	slog.Info("Initializing application...")

	// Example: Initialize database connection (replace with actual logic)
	dbPath := "app.db" // TODO: Make configurable
	db, err := core.SetupDatabase(dbPath)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}
	defer db.Close()

	// Example: Initialize message log
	// TODO: Instantiate a real encoder (e.g., JSONEncoder)
	// messageLog, err := core.NewMessageLog(db, &core.JSONEncoder{})
	// if err != nil {
	// 	return fmt.Errorf("failed to initialize message log: %w", err)
	// }
	// TODO: Register message types with messageLog.RegisterType(...)

	// Example: Initialize command/query registries (assuming global or passed instances)
	// commandRegistry := core.NewCommandRegistry()
	// queryRegistry := core.NewQueryRegistry()

	// Example: Initialize application state by replaying log
	// TODO: Implement state struct(s) and Apply method
	// appState := core.NewAppState() // Replace with actual state struct
	// messages, err := messageLog.Load(context.Background())
	// if err != nil {
	// 	return fmt.Errorf("failed to load messages for replay: %w", err)
	// }
	// for _, msg := range messages {
	// 	if err := appState.Apply(msg); err != nil { // Assuming Apply method exists
	// 		slog.Error("Failed to apply message during replay", "error", err, "message", msg)
	// 		// Decide whether to continue or fail startup
	// 	}
	// }
	// slog.Info("State replay completed", "message_count", len(messages))

	// Example: Register feature handlers
	// RegisterAllFeatures(commandRegistry, queryRegistry /*, appState, messageLog */) // Pass necessary dependencies

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
	// TODO: Add routes for features (e.g., mux.HandleFunc("GET /posts", handleListPosts), mux.HandleFunc("GET /posts/{id}", handleGetPost))

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
