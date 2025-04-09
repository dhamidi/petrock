# Plan for core/commands.go

This file defines the core components for command handling: the `Executor` which orchestrates command processing, the `CommandRegistry` which maps command names to state update handlers, and associated interfaces/types.

## Types

- `Command`: An interface representing a command message.
    - `CommandName() string`: Returns the unique kebab-case name (e.g., "feature/create-item").
    - `Validate() error` (Optional): If implemented, this method is called by the `core.Executor` before logging and execution. It should return `nil` if the command data is valid, or an error otherwise.
- `CommandHandler func(ctx context.Context, cmd Command) error`: A function type for handlers that **apply the state changes** for a command. This handler is executed *after* the command has been validated and appended to the log. It should only contain the logic to modify the application state based on the command data. Returning an error from this handler will cause the `core.Executor` to panic, as it indicates an inconsistency between the log and the state logic.
- `CommandRegistry`: A struct responsible for mapping command names to their state update handlers and types.
    - `handlers map[string]CommandHandler`: Map from command name to state update handler.
    - `types map[string]reflect.Type`: Map from command name to `reflect.Type`.
    - `mu sync.RWMutex`: For thread-safe access.
- `Executor`: The central orchestrator for command execution.
    - `log *MessageLog`: Dependency for appending commands.
    - `registry *CommandRegistry`: Dependency for finding state update handlers.

## Functions

- `NewCommandRegistry() *CommandRegistry`: Constructor for `CommandRegistry`.
- `(r *CommandRegistry) Register(cmd Command, handler CommandHandler)`: Registers a state update `handler` for the given `cmd` instance. Uses `cmd.CommandName()` as the key. Stores the handler and `reflect.Type`. Panics if the name is already registered.
- `(r *CommandRegistry) GetHandler(name string) (CommandHandler, bool)`: Retrieves the registered state update handler for a given command name.
- `(r *CommandRegistry) GetCommandType(name string) (reflect.Type, bool)`: Looks up and returns the `reflect.Type` for a command based on its registered name.
- `(r *CommandRegistry) RegisteredCommandNames() []string`: Returns a slice containing the registered command names.
- `NewExecutor(log *MessageLog, registry *CommandRegistry) *Executor`: Constructor for `Executor`.
- `(e *Executor) Execute(ctx context.Context, cmd Command) error`: Orchestrates command execution:
    1. Calls `cmd.Validate()` if implemented. Returns validation error if it fails.
    2. Appends the command to the message log via `e.log.Append(ctx, cmd)`. Returns logging error if it fails.
    3. Retrieves the state update handler from `e.registry.GetHandler(cmd.CommandName())`. Returns error if handler not found.
    4. Executes the state update handler: `handler(ctx, cmd)`.
    5. If the handler returns an error, `panic` immediately. This indicates an unrecoverable state inconsistency requiring a restart.
    6. Returns `nil` on successful execution (validation, logging, and state update).
