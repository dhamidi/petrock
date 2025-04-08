# Plan for core/commands.go

This file defines the registry for commands and their associated handlers, forming the core of the command dispatch system. It works in conjunction with the Executor system to provide a standardized workflow for command processing.

## Types

- `Command`: An interface representing a command message. All command structs should implicitly satisfy this (e.g., by being `interface{}`). It requires `CommandName() string`.

- `CommandHandler func(ctx context.Context, cmd Command) error`: A function type for handlers that process commands. Takes context and the command message, returns an error if processing fails. With the centralized executor pattern, handlers focus on business logic rather than execution flow.

- `CommandRegistry`: A struct responsible for mapping command types to their handlers.
    - `handlers map[reflect.Type]CommandHandler`: The internal map storing the registrations.
    - `mu sync.RWMutex`: For thread-safe access to the handlers map.

## Functions

- `NewCommandRegistry() *CommandRegistry`: Constructor function to create and initialize a new `CommandRegistry`.

- `(r *CommandRegistry) Register(cmd Command, handler CommandHandler)`: Registers a handler using the name returned by `cmd.CommandName()`. Stores the handler and `reflect.Type`. Panics if the name is already registered.

- `(r *CommandRegistry) Dispatch(ctx context.Context, cmd Command) error`: Internal method that looks up the handler using `cmd.CommandName()` and executes it. This is typically called by the Executor after validation and logging, not directly by application code.

- `(r *CommandRegistry) RegisteredCommandNames() []string`: Returns a slice containing the full registered kebab-case names (e.g., "posts/create") of all commands.

- `(r *CommandRegistry) GetCommandType(name string) (reflect.Type, bool)`: Looks up and returns the `reflect.Type` for a command based on its full registered kebab-case name (e.g., "posts/create").

## Integration with Executor

The CommandRegistry works closely with the Executor system:

1. Features register their command handlers with the CommandRegistry during setup.
2. Application code executes commands via the Executor's Execute method.
3. The Executor handles validation, logging, and then uses the CommandRegistry to dispatch to the appropriate handler.
4. Handlers now focus on domain-specific logic rather than repeating the validation/logging pattern.
