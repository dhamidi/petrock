# Plan for posts/register.go (Example Feature)

This file acts as the entry point for the feature module. Its primary role is to register the feature's command and query handlers with the core registries.

## Types

- None specific to this file.

## Functions

- `RegisterFeature(mux *http.ServeMux, commands *core.CommandRegistry, queries *core.QueryRegistry, messageLog *core.MessageLog, executor *core.Executor, state *State, /* other shared deps */)`: This function initializes the feature's components, registers its state update handlers and query handlers with the core registries, registers message types, and sets up feature-specific HTTP routes.
  - **Dependencies:** Receives core components: HTTP router (`mux`), command/query registries, message log, the central `core.Executor`, and the feature's specific state (`*State`). May receive others like `*sql.DB`.
  - **Initialization:**
    - Creates the feature's state updater (e.g., `updater := NewStateUpdater(state)` from `execute.go`).
    - Creates the feature's querier (e.g., `querier := NewQuerier(state)` from `query.go`).
    - Creates the feature's HTTP server (e.g., `server := NewFeatureServer(executor, querier, state, /* other deps */)` from `http.go`), passing the *central* `core.Executor` and other needed dependencies.
  - **Core Registration:**
    - Calls `commands.Register` for each command type, passing the corresponding *state update handler* from the `updater` (e.g., `commands.Register(CreateCommand{}, updater.HandleCreatePost)`).
    - Calls `queries.Register` for each query type, passing the handler from the `querier` (e.g., `queries.Register(GetQuery{}, querier.HandleGetPost)`).
    - Calls `RegisterTypes(messageLog)` (typically defined in `messages.go` or `state.go`) to register command/query types with the `core.MessageLog` for decoding.
  - **HTTP Route Registration:**
    - Calls `RegisterRoutes(mux, server)` (defined in `routes.go`) to register the feature's specific HTTP routes.
  - **Background Jobs:** May initialize and start background jobs (from `jobs.go`).

_Note: The `petrock feature <name>` command automatically adds the necessary import and the call to this `RegisterFeature` function (with the updated signature) within the project's `cmd/<project>/features.go` file. The registration happens *after* core routes are registered, allowing features to override core routes._
