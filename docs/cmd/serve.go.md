# Plan for cmd/blog/serve.go

This file defines the `serve` subcommand, responsible for starting the HTTP server and managing background workers.

## Types

- None specific to this file.

## Functions

- `NewServeCmd() *cobra.Command`: Creates and configures the `serve` subcommand, including flags (e.g., `--port`, `--host`). Returns the Cobra command object.
- `runServe(cmd *cobra.Command, args []string) error`: The function executed when the `serve` command is invoked.
    1. Parses flags (`--port`, `--host`).
    2. Initializes core components:
        - Database connection (`*sql.DB`).
        - Message Log (`*core.MessageLog`).
        - Command Registry (`*core.CommandRegistry`).
        - Query Registry (`*core.QueryRegistry`).
        - Central Command Executor (`*core.Executor`), passing it the log and command registry.
        - Application State (potentially a map or struct holding each feature's state, e.g., `map[string]interface{}`).
    3. **Register Message Types:** Call `RegisterTypes` for all known command/query types with the `messageLog`. This is crucial *before* replay.
    4. **Initialize Feature States:** Create initial instances of each feature's state (e.g., `posts.NewPostState()`).
    5. **Register Features:** Call `RegisterAllFeatures(...)` (from `cmd/<project>/features.go`), passing registries, log, executor, feature states, etc. This populates the command registry with state update handlers.
    6. **Replay Log & Rebuild State:**
        - Get starting version: `startVersion := uint64(0)` (start from beginning).
        - Iterate through messages using the iterator: 
            `for msg := range messageLog.After(ctx, startVersion) {`
            - Access decoded message directly: `decodedMsg := msg.DecodedPayload`.
            - If it's a command (`core.Command`):
                - Look up state update handler: `handler, ok := commandRegistry.GetHandler(cmd.CommandName())`.
                - If handler found, execute it with both the payload and message metadata: `err := handler(ctx, decodedMsg, &msg.Message)`. **Panic on error here.**
            - If it's a query or unknown type, ignore for replay.
    7. Creates the main HTTP router (`mux := http.NewServeMux()`).
    8. **Registers core HTTP handlers** (e.g., `/`, `/commands`, `/queries`) using the `executor`, `queryRegistry`, etc.
    9. **Feature HTTP Routes:** Feature routes were already registered inside `RegisterAllFeatures` by calling each feature's `RegisterRoutes`.
    10. **Start Workers:** Call `app.StartWorkers(ctx)` to initialize and start all registered workers.
    11. Sets up and starts the HTTP server.
    12. On shutdown signal, calls `app.StopWorkers(ctx)` to gracefully stop all workers.
