package routes

import (
	"log/slog"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/petrock_example_feature_name/handlers"
)

// RegisterRoutes defines the feature-specific HTTP routes and registers them
// with the main application's router.
// It receives the router and the feature's dependency container (FeatureServer).
func RegisterRoutes(app *core.App, deps *handlers.FeatureServer) {
	// Safety check to prevent nil pointer dereference
	if app == nil || app.Mux == nil {
		slog.Error("Error registering routes: nil app or nil HTTP mux provided", "feature", "petrock_example_feature_name")
		return
	}

	// Safety check for dependencies
	if deps == nil {
		slog.Error("Error registering routes: nil FeatureServer provided", "feature", "petrock_example_feature_name")
		return
	}

	// Register Web UI routes
	registerWebRoutes(app, deps)

	// Register API routes
	registerAPIRoutes(app, deps)

	// --- Overriding Example (Use with caution!) ---
	// If you uncomment the line below, this feature's HandleCustomIndex
	// will handle requests to the root path "/", overriding the core index handler.
	// app.RegisterRoute("GET /", deps.HandleCustomIndex)

	slog.Info("Registered feature HTTP routes", "feature", "petrock_example_feature_name")
}