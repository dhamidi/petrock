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

	// GET /petrock_example_feature_name/{id} - Get a specific item
	// Note: Go 1.22+ required for path parameters in ServeMux patterns
	mux.HandleFunc("GET "+featurePrefix+"/{id}", deps.HandleGetItem)

	// GET /petrock_example_feature_name/ - List items (example)
	mux.HandleFunc("GET "+featurePrefix+"/", deps.HandleListItems)

	// POST /petrock_example_feature_name/ - Create a new item
	mux.HandleFunc("POST "+featurePrefix+"/", deps.HandleCreateItem)

	// PUT /petrock_example_feature_name/{id} - Update an item (example)
	mux.HandleFunc("PUT "+featurePrefix+"/{id}", deps.HandleUpdateItem)

	// DELETE /petrock_example_feature_name/{id} - Delete an item (example)
	mux.HandleFunc("DELETE "+featurePrefix+"/{id}", deps.HandleDeleteItem)

	// Add more feature-specific routes here...

	// --- Overriding Example (Use with caution!) ---
	// If you uncomment the line below, this feature's HandleCustomIndex
	// will handle requests to the root path "/", overriding the core index handler.
	// mux.HandleFunc("GET /", deps.HandleCustomIndex)

	slog.Info("Registered feature HTTP routes", "feature", "petrock_example_feature_name")
}
