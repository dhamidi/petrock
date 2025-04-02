# Plan for core/commands.go

This file defines the registry for commands and their associated handlers, forming the core of the command dispatch system.

## Types

- `Command`: An interface representing a command message. All command structs should implicitly satisfy this (e.g., by being `interface{}`). It's primarily a marker.
- `CommandHandler func(ctx context.Context, cmd Command) error`: A function type for handlers that process commands. Takes context and the command message, returns an error if processing fails.
- `CommandRegistry`: A struct responsible for mapping command types to their handlers.
    - `handlers map[reflect.Type]CommandHandler`: The internal map storing the registrations.
    - `mu sync.RWMutex`: For thread-safe access to the handlers map.

## Functions

- `NewCommandRegistry() *CommandRegistry`: Constructor function to create and initialize a new `CommandRegistry`.
- `(r *CommandRegistry) Register(cmd Command, handler CommandHandler)`: Registers a command handler for a specific command type. It uses reflection (`reflect.TypeOf(cmd)`) to get the type key.
- `(r *CommandRegistry) Dispatch(ctx context.Context, cmd Command) error`: Looks up the handler for the given command's type and executes it. Returns an error if no handler is registered or if the handler itself returns an error.
