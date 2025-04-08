# Plan for <feature>/http.go (Example Feature)

This file contains the HTTP request handlers for the routes defined in `<feature>/routes.go`. It bridges the web layer with the feature's core logic (Executor, Querier, State).

## Types

- `FeatureServer`: A struct designed to hold all dependencies required by the feature's HTTP handlers. This promotes clean dependency injection.
    - `executor *Executor`: Instance of the feature's command executor (from `execute.go`).
    - `querier *Querier`: Instance of the feature's query handler (from `query.go`).
    - `state *State`: Instance of the feature's state (from `state.go`).
    - `log *core.MessageLog`: Shared message log (optional, if handlers need direct logging).
    - `commands *core.CommandRegistry`: Shared command registry (optional, if handlers need to dispatch commands through the core system).
    - `db *sql.DB`: Shared database connection (optional, if handlers need direct DB access).
    - *Other shared dependencies as needed (e.g., config, template renderer).*

## Functions

- `NewFeatureServer(executor *Executor, querier *Querier, state *State, /* other deps */) FeatureServer`: Constructor function to create and initialize the `FeatureServer` struct with its dependencies. This is typically called within `RegisterFeature` in `register.go`.

- **Handler Methods:** These are methods attached to the `FeatureServer` struct, each implementing `http.HandlerFunc` or being compatible with it. They are registered in `routes.go`.
    - `(fs *FeatureServer) HandleGetItem(w http.ResponseWriter, r *http.Request)`: Example handler for retrieving an item.
        - Parses request parameters (e.g., item ID from URL path `r.PathValue("id")`).
        - Calls the appropriate `fs.querier` method (e.g., `fs.querier.HandleGet(ctx, GetQuery{ID: itemID})`).
        - Handles errors returned by the querier (e.g., return 404 Not Found, 500 Internal Server Error).
        - Formats the query result (e.g., marshal to JSON) and writes it to the `http.ResponseWriter`. Sets appropriate headers (e.g., `Content-Type: application/json`).
    - `(fs *FeatureServer) HandleCreateItem(w http.ResponseWriter, r *http.Request)`: Example handler for creating an item.
        - Parses the request body (e.g., decode JSON into a `CreateCommand` struct). Handle decoding errors (return 400 Bad Request).
        - Performs validation if needed.
        - **State Change Strategy:** To ensure changes are logged and state is applied consistently with the event sourcing pattern:
            1. Construct the appropriate `Command` struct (e.g., `CreateCommand`).
            2. Validate the command data.
            3. Append the command to the log: `err := fs.log.Append(r.Context(), cmd)`. Handle errors.
            4. *Optional/Implicit:* The state's `Apply` method will typically be triggered by the log replay mechanism or potentially by a direct call after logging if immediate state update is needed before the next request cycle. Alternatively, use the command registry: `err := fs.commands.Dispatch(r.Context(), cmd)`. This handles logging and state application via the registered command handler.
        - Handles errors (e.g., return 400 Bad Request for validation, 500 Internal Server Error).
        - Writes an appropriate success response (e.g., 201 Created with Location header, 200 OK with created item JSON, 202 Accepted).
    - *Other handlers as needed (List, Update, Delete, custom actions).*
