package routes

import (
	"log/slog"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/petrock_example_feature_name/handlers"
)

// registerAPIRoutes registers the REST API routes for the feature
func registerAPIRoutes(app *core.App, deps *handlers.FeatureServer) {
	apiPrefix := "/api/petrock_example_feature_name" // Base path for this feature's API routes
	slog.Debug("Registering API routes", "feature", "petrock_example_feature_name", "prefix", apiPrefix)

	// Example API routes - not implemented in the original code but showing as an example
	// POST /api/petrock_example_feature_name - Create a new item via API
	// app.RegisterRoute("POST "+apiPrefix, deps.HandleCreateItem)

	// GET /api/petrock_example_feature_name/{id} - Get a specific item via API
	// app.RegisterRoute("GET "+apiPrefix+"/{id}", deps.HandleGetItemAPI)

	// PUT /api/petrock_example_feature_name/{id} - Update a specific item via API
	// app.RegisterRoute("PUT "+apiPrefix+"/{id}", deps.HandleUpdateItem)

	// DELETE /api/petrock_example_feature_name/{id} - Delete a specific item via API
	// app.RegisterRoute("DELETE "+apiPrefix+"/{id}", deps.HandleDeleteItem)
}