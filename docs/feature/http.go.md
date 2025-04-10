# Plan for <feature>/http.go (Example Feature)

This file contains the HTTP request handlers for the routes defined in `<feature>/routes.go`. It bridges the web layer with the feature's core logic (Executor, Querier, State).

## Types

- `FeatureServer`: A struct holding dependencies for the feature's HTTP handlers.
    - `executor *core.Executor`: The central command executor.
    - `querier *Querier`: Instance of the feature's query handler (from `query.go`).
    - `state *State`: Instance of the feature's state (from `state.go`).
    - `db *sql.DB`: Shared database connection (optional).
    - *Other shared dependencies (e.g., config, template engine).*

## Functions

- `NewFeatureServer(executor *core.Executor, querier *Querier, state *State, /* other deps */) *FeatureServer`: Constructor for `FeatureServer`. Called in `RegisterFeature`.

- **Handler Methods:** Methods on `FeatureServer` implementing `http.HandlerFunc`.
    - `(fs *FeatureServer) HandleGetItem(w http.ResponseWriter, r *http.Request)`: Handles item retrieval.
        - Parses request (e.g., ID from path).
        - Creates the appropriate `GetQuery` struct.
        - Calls `fs.querier.HandleGetPost(r.Context(), query)`.
        - Handles querier errors (404, 500).
        - Writes JSON response.
    - `(fs *FeatureServer) HandleCreateItem(w http.ResponseWriter, r *http.Request)`: Handles item creation.
        - Parses request body into a `CreateCommand` struct. Handle decoding errors (400 Bad Request).
        - **State Change Strategy:** Uses the central executor:
            ```go
            cmd := CreateCommand{ /* ... data from request ... */ }
            err := fs.executor.Execute(r.Context(), cmd)
            if err != nil {
                // Handle validation errors (e.g., return 400 Bad Request)
                // Handle logging errors (e.g., return 500 Internal Server Error)
                // Note: State update errors cause panic within Execute, so won't be handled here.
                http.Error(w, err.Error(), http.StatusBadRequest) // Example
                return
            }
            ```
        - Writes success response (e.g., 201 Created, 200 OK, 202 Accepted).
    - *Other handlers (List, Update, Delete) follow similar patterns, using `fs.querier` for reads and `fs.executor.Execute` for writes.*
