# Plan for cmd/blog/deploy.go

This file defines the `deploy` subcommand, responsible for copying the built binary to a target host via SSH and potentially restarting the service.

## Types

- None specific to this file.

## Functions

- `NewDeployCmd() *cobra.Command`: Creates and configures the `deploy` subcommand, including flags (e.g., `--target-host`, `--target-path`, `--ssh-user`, `--ssh-key`). Returns the Cobra command object.
- `runDeploy(cmd *cobra.Command, args []string) error`: The function executed when the `deploy` command is invoked. It likely first runs the build process (or ensures a binary exists), then uses SSH (e.g., via `golang.org/x/crypto/ssh`) to connect to the target host, copy the binary, and execute remote commands (like restarting a systemd service).
