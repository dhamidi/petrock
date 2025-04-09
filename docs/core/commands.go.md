# Plan for core/commands.go

This file defines the core components for command handling: the `Executor` which orchestrates command processing, the `CommandRegistry` which maps command names to state update handlers, and associated interfaces/types.

## Types

- `Command`: An interface representing a command message.
    - `CommandName() string`: Returns the unique kebab-case name (e.g., "feature/create-item").
- `CommandHandler func(ctx context.Context, cmd Command) error`: A function type for handlers that **apply the state changes** for a command. This handler is executed *after* the command has been validated (by the feature executor) and appended to the log. It should only contain the logic to modify the application state based on the command data. Returning an error from this handler will cause the `core.Executor` to panic, as it indicates an inconsistency between the log and the state logic.
- `CommandValidator`: An interface that feature-specific executors can implement to handle command validation.
    - `ValidateCommand(ctx context.Context, cmd Command) error`: Validates the command, potentially using the feature's current state. Returns `nil` if valid, or an error describing the validation failure.
- `CommandRegistry`: A struct responsible for mapping command names to their state update handlers, validators (feature executors), and types.
    - `handlers map[string]CommandHandler`: Map from command name to state update handler function.
    - `validators map[string]CommandValidator`: Map from command name to the feature executor instance responsible for validating this command.
    - `types map[string]reflect.Type`: Map from command name to `reflect.Type`.
    - `mu sync.RWMutex`: For thread-safe access.
- `Executor`: The central orchestrator for command execution.
    - `log *MessageLog`: Dependency for appending commands.
    - `registry *CommandRegistry`: Dependency for finding state update handlers and validators.

## Functions

- `NewCommandRegistry() *CommandRegistry`: Constructor for `CommandRegistry`.
- `(r *CommandRegistry) Register(cmd Command, handler CommandHandler, validator CommandValidator)`: Registers a state update `handler` and a `validator` (typically the feature executor instance) for the given `cmd` instance. Uses `cmd.CommandName()` as the key. Stores the handler, validator, and `reflect.Type`. Panics if the name is already registered.
- `(r *CommandRegistry) GetHandler(name string) (CommandHandler, bool)`: Retrieves the registered state update handler for a given command name.
- `(r *CommandRegistry) GetValidator(name string) (CommandValidator, bool)`: Retrieves the registered validator instance for a given command name.
- `(r *CommandRegistry) GetHandlerAndValidator(name string) (CommandHandler, CommandValidator, bool)`: Retrieves both the handler and validator for a given command name.
- `(r *CommandRegistry) GetCommandType(name string) (reflect.Type, bool)`: Looks up and returns the `reflect.Type` for a command based on its registered name.
- `(r *CommandRegistry) RegisteredCommandNames() []string`: Returns a slice containing the registered command names.
- `NewExecutor(log *MessageLog, registry *CommandRegistry) *Executor`: Constructor for `Executor`.
- `(e *Executor) Execute(ctx context.Context, cmd Command) error`: Orchestrates command execution:
    1. Retrieves the state update handler and validator from `e.registry.GetHandlerAndValidator(cmd.CommandName())`. Returns error if not found.
    2. Calls the validator's validation method: `err := validator.ValidateCommand(ctx, cmd)`. Returns validation error if it fails.
    3. Appends the command to the message log via `e.log.Append(ctx, cmd)`. Returns logging error if it fails.
    4. Executes the state update handler: `handlerErr := handler(ctx, cmd)`.
    5. If the handler returns an error (`handlerErr != nil`), `panic` immediately. This indicates an unrecoverable state inconsistency requiring a restart.
    6. Returns `nil` on successful execution (validation, logging, and state update).
