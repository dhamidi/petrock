package petrock_example_feature_name

import (
	"log/slog"
	"net/http"

	// No direct dependency on core needed here usually, unless for shared types used in handlers
)

// RegisterRoutes defines the feature-specific HTTP routes and registers them
// with the main application's router.
// It receives the router and the feature's dependency container (FeatureServer).
func RegisterRoutes(mux *http.ServeMux, deps *FeatureServer) {
	featurePrefix := "/petrock_example_feature_name" // Base path for this feature's routes
	slog.Debug("Registering feature HTTP routes", "feature", "petrock_example_feature_name", "prefix", featurePrefix)

	// Example Routes:
	// It's strongly recommended to prefix feature routes to avoid collisions.

	// GET /petrock_example_feature_name/ - List items 
	mux.HandleFunc("GET "+featurePrefix+"/", deps.HandleListItems)

	// GET /petrock_example_feature_name/{id} - View a specific item
	// Note: Go 1.22+ required for path parameters in ServeMux patterns
	mux.HandleFunc("GET "+featurePrefix+"/{id}", deps.HandleGetItem)

	// GET /petrock_example_feature_name/new - Show form for creating a new item
	mux.HandleFunc("GET "+featurePrefix+"/new", deps.HandleNewForm)

	// POST /petrock_example_feature_name/new - Process the creation form
	mux.HandleFunc("POST "+featurePrefix+"/new", deps.HandleCreateForm)

	// GET /petrock_example_feature_name/{id}/edit - Show form for editing an item
	mux.HandleFunc("GET "+featurePrefix+"/{id}/edit", deps.HandleEditForm)

	// POST /petrock_example_feature_name/{id}/edit - Process the edit form
	mux.HandleFunc("POST "+featurePrefix+"/{id}/edit", deps.HandleUpdateForm)

	// GET /petrock_example_feature_name/{id}/delete - Show delete confirmation
	mux.HandleFunc("GET "+featurePrefix+"/{id}/delete", deps.HandleDeleteForm)

	// POST /petrock_example_feature_name/{id}/delete - Process the deletion
	mux.HandleFunc("POST "+featurePrefix+"/"+"{id}/delete", deps.HandleDeleteConfirm)

	// Add more feature-specific routes here...

	// --- Overriding Example (Use with caution!) ---
	// If you uncomment the line below, this feature's HandleCustomIndex
	// will handle requests to the root path "/", overriding the core index handler.
	// mux.HandleFunc("GET /", deps.HandleCustomIndex)

	slog.Info("Registered feature HTTP routes", "feature", "petrock_example_feature_name")
}
