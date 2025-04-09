# Plan for posts/register.go (Example Feature)

This file acts as the entry point for the feature module. Its primary role is to register the feature's command and query handlers with the core registries.

## Types

- None specific to this file.

## Functions

- `RegisterFeature(mux *http.ServeMux, commands *core.CommandRegistry, queries *core.QueryRegistry, messageLog *core.MessageLog, state *State, db *sql.DB, executor core.Executor)`: This function initializes the feature's handlers, registers them with core registries, and registers any feature-specific HTTP routes.
  - **Dependencies:** It receives shared core components like the main HTTP router (`mux`), command/query registries, message log, database connection (`db`), the feature's specific state, and the core executor. The executor is a critical dependency that centralizes command validation, logging, and dispatch.
  - **Initialization:**
    - Creates instances of the feature's command handler (e.g., `NewPostExecutor(state)`) and querier (e.g., `NewQuerier(state)`).
    - Uses the provided core `executor` parameter for command execution rather than accessing the message log directly.
    - Creates an instance of the feature's HTTP handler container (e.g., `server := NewFeatureServer(coreExecutor, featureExecutor, querier, state, commands, db)` from `http.go`), passing necessary dependencies.
  - **Core Registration:**
    - Calls `commands.Register` for each command type (e.g., `commands.Register(CreateCommand{}, executor.HandleCreate)`).
    - Calls `queries.Register` for each query type (e.g., `queries.Register(GetQuery{}, querier.HandleGet)`).
    - Calls `RegisterTypes(messageLog)` (typically defined in `state.go` or `messages.go`) to register command/query types with the `core.MessageLog` for decoding during replay.
  - **HTTP Route Registration:**
    - Calls `RegisterRoutes(mux, server)` (defined in `routes.go`) to register the feature's specific HTTP routes with the main application router.
  - **Background Jobs:** It might initialize and register background jobs/workers if defined in `jobs.go`.

_Note: The `petrock feature <name>` command automatically adds the necessary import and the call to this `RegisterFeature` function (with the updated signature) within the project's `cmd/<project>/features.go` file. The registration happens *after* core routes are registered, allowing features to override core routes._
