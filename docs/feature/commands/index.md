# Commands

The `commands` directory contains the command definitions and implementations for the feature, following the command pattern for handling operations that change state.

## Structure

- `base.go` - Common command interfaces and types
- `create.go` - Commands for creating new items
- `update.go` - Commands for updating existing items
- `delete.go` - Commands for deleting items
- `request_summary.go` - Commands for requesting summary generation
- `set_summary.go` - Commands for setting generated summaries
- `fail_summary.go` - Commands for handling failed summary generation
- `register.go` - Command registration with the system

## Command Pattern

Commands follow a pattern where:

1. Each command is a distinct type with its own fields
2. Commands are executed by a command handler
3. Commands represent intent to change state rather than directly changing it
4. Commands may be validated before execution

## Registering Commands

Commands must be registered with the system to be available for use. The registration happens in `register.go` and is typically called during feature initialization.