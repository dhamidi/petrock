# Plan for posts/register.go (Example Feature)

This file acts as the entry point for the feature module. Its primary role is to register the feature's command and query handlers with the core registries.

## Types

- None specific to this file.

## Functions

- `RegisterFeature(commands *core.CommandRegistry, queries *core.QueryRegistry, state *PostState)`: This function initializes the feature's handlers and registers them.
    - It creates instances of the feature's executor (e.g., `NewPostExecutor(state)`) and querier (e.g., `NewPostQuerier(state)`).
    - It calls `commands.Register` for each command type defined in `posts/messages.go`, passing the corresponding handler method (e.g., `executor.HandleCreatePost`).
    - It calls `queries.Register` for each query type defined in `posts/messages.go`, passing the corresponding handler method (e.g., `querier.HandleGetPost`).
    - It calls the feature's `RegisterTypes` function (defined in `state.go`) to register command/event types with the `core.MessageLog` for decoding during replay.
    - It might initialize and register background jobs/workers if defined in `jobs.go`.

*Note: The `petrock feature <name>` command automatically adds the necessary import and the call to this `RegisterFeature` function within the project's `cmd/<project>/features.go` file.*
