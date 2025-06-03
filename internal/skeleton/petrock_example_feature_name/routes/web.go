package routes

import (
	"log/slog"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/petrock_example_feature_name/handlers"
)

// registerWebRoutes registers the web UI routes for the feature
func registerWebRoutes(app *core.App, deps *handlers.FeatureServer) {
	featurePrefix := "/petrock_example_feature_name" // Base path for this feature's routes
	slog.Debug("Registering web UI routes", "feature", "petrock_example_feature_name", "prefix", featurePrefix)

	// GET /petrock_example_feature_name/ - List items
	app.RegisterRoute("GET "+featurePrefix+"/", deps.HandleListItems)

	// Register specific routes before parameterized routes to avoid conflicts
	// GET /petrock_example_feature_name/new - Show form for creating a new item
	app.RegisterRoute("GET "+featurePrefix+"/new", deps.HandleNewForm)

	// POST /petrock_example_feature_name/new - Process the creation form
	app.RegisterRoute("POST "+featurePrefix+"/new", deps.HandleCreateForm)

	// GET /petrock_example_feature_name/{id}/edit - Show form for editing an item
	app.RegisterRoute("GET "+featurePrefix+"/{id}/edit", deps.HandleEditForm)

	// POST /petrock_example_feature_name/{id}/edit - Process the edit form
	app.RegisterRoute("POST "+featurePrefix+"/{id}/edit", deps.HandleUpdateForm)

	// GET /petrock_example_feature_name/{id}/delete - Show delete confirmation
	app.RegisterRoute("GET "+featurePrefix+"/{id}/delete", deps.HandleDeleteForm)

	// POST /petrock_example_feature_name/{id}/delete - Process the deletion
	app.RegisterRoute("POST "+featurePrefix+"/"+"{id}/delete", deps.HandleDeleteConfirm)

	// GET /petrock_example_feature_name/{id} - View a specific item
	// Note: Go 1.22+ required for path parameters in ServeMux patterns
	// This must be registered last as it's the most general pattern
	app.RegisterRoute("GET "+featurePrefix+"/{id}", deps.HandleGetItem)
}
