# Plan for <feature>/routes.go (Example Feature)

This file is responsible for defining the specific HTTP routes exposed by the feature.

## Types

- None specific to this file.

## Functions

- `RegisterRoutes(mux *http.ServeMux, deps FeatureServer)`: This function takes the main application's router (`*http.ServeMux`) and the feature's dependency container (`FeatureServer` struct, typically defined in `http.go`) as input.
    - It defines the feature's HTTP routes using `mux.HandleFunc` or `mux.Handle`.
    - Example: `mux.HandleFunc("GET /feature-prefix/{id}", deps.HandleGetItem)`
    - Example: `mux.HandleFunc("POST /feature-prefix/", deps.HandleCreateItem)`
    - **Convention:** It's strongly recommended to prefix feature-specific routes (e.g., `/feature-prefix/`) to avoid accidental collision with core routes or routes from other features.
    - **Overriding:** Since features are registered *after* core routes in `serve.go`, defining a route here with the same pattern as a core route (e.g., `"GET /"`) will effectively override the core handler for that route. This should be done intentionally and with caution.
    - Handlers referenced here (e.g., `deps.HandleGetItem`) are typically methods on the `FeatureServer` struct defined in `http.go`.
