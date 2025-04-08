# Plan for cmd/blog/serve.go

This file defines the `serve` subcommand, responsible for starting the HTTP server.

## Types

- None specific to this file.

## Functions

- `NewServeCmd() *cobra.Command`: Creates and configures the `serve` subcommand, including flags (e.g., `--port`, `--host`). Returns the Cobra command object.
- `runServe(cmd *cobra.Command, args []string) error`: The function executed when the `serve` command is invoked. It performs the following steps:
    1. Parses flags (`--port`, `--host`, etc.).
    2. Initializes core components: database connection (`*sql.DB`), message log (`*core.MessageLog`), command registry (`*core.CommandRegistry`), query registry (`*core.QueryRegistry`), application state (`AppState`).
    3. Replays messages from the log to rebuild the application state.
    4. Creates the main HTTP router (`mux := http.NewServeMux()`).
    5. **Registers core HTTP handlers** on the `mux` (e.g., for `/`, `/commands`, `/queries`).
    6. **Calls `RegisterAllFeatures(mux, commands, queries, messageLog, appState, ...)`** (defined in `cmd/<project>/features.go`), passing the router and other initialized components. This allows features to register their core handlers *and* their specific HTTP routes. Crucially, this happens *after* core routes are registered, enabling features to override them if necessary.
    7. Sets up and starts the HTTP server using the configured `mux`.
