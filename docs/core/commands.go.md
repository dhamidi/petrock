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
- `Command`: Interface that command structs must implement. Embeds `NamedMessage`.
- `NamedMessage`: Interface with `RegisteredName() string` method.
- `CommandHandler`: Function type for command handlers.
- `CommandRegistry`: Maps command names (`feature/Type`) to handlers and `reflect.Type`.
- `NewCommandRegistry()`: Constructor.
- `(r *CommandRegistry) Register(cmd Command, handler CommandHandler)`: Registers a handler using the name returned by `cmd.RegisteredName()`. Stores the handler and `reflect.Type`. Panics if the name is already registered.
- `(r *CommandRegistry) Dispatch(ctx context.Context, cmd Command) error`: Looks up the handler using `cmd.RegisteredName()` and executes it.
- `(r *CommandRegistry) RegisteredCommandNames() []string`: Returns a slice containing the full registered names (e.g., "posts/CreateCommand") of all commands.
- `(r *CommandRegistry) GetCommandType(name string) (reflect.Type, bool)`: Looks up and returns the `reflect.Type` for a command based on its full registered name (e.g., "posts/CreateCommand").
