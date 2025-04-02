# Plan for cmd/blog/serve.go

This file defines the `serve` subcommand, responsible for starting the HTTP server.

## Types

- None specific to this file.

## Functions

- `NewServeCmd() *cobra.Command`: Creates and configures the `serve` subcommand, including flags (e.g., `--port`, `--host`). Returns the Cobra command object.
- `runServe(cmd *cobra.Command, args []string) error`: The function executed when the `serve` command is invoked. It parses flags, initializes the core components (like registries, state, log), sets up the HTTP router (`net/http.ServeMux`) and handlers, and starts the HTTP server.
