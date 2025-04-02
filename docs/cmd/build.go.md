# Plan for cmd/blog/build.go

This file defines the `build` subcommand, responsible for compiling the application into a single distributable binary, potentially including embedded assets.

## Types

- None specific to this file.

## Functions

- `NewBuildCmd() *cobra.Command`: Creates and configures the `build` subcommand, including flags (e.g., `--output`, `--os`, `--arch`). Returns the Cobra command object.
- `runBuild(cmd *cobra.Command, args []string) error`: The function executed when the `build` command is invoked. It runs the `go build` command with appropriate flags (like `-ldflags="-s -w"`) and embeds assets if necessary.
