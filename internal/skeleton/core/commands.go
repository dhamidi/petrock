package core

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
)

// NamedMessage defines an interface for messages that know their registered name.
type NamedMessage interface {
	RegisteredName() string // e.g., "feature/TypeName"
}

// Command is an interface combining NamedMessage for command messages.
type Command interface {
	CommandName() string // e.g., "feature/create-command"
}

// CommandHandler defines the function signature for handling commands.
// These handlers are responsible *only* for applying state changes after validation and logging.
// Returning an error from a CommandHandler will cause the core.Executor to panic.
type CommandHandler func(ctx context.Context, cmd Command) error

// FeatureExecutor defines an interface that feature-specific executors must implement.
// This allows the central core.Executor to delegate command validation to the appropriate feature.
type FeatureExecutor interface {
	// ValidateCommand checks if a command is valid according to the feature's rules and state.
	// The implementation should typically check if the command implements a feature-specific
	// Validator interface and call its Validate(state) method if it does.
	ValidateCommand(ctx context.Context, cmd Command) error
}

// CommandRegistry maps command names (feature/Type) to their state update handlers,
// feature executors, and types.
type CommandRegistry struct {
	handlers         map[string]CommandHandler  // Key: "feature/TypeName" -> State update handler
	featureExecutors map[string]FeatureExecutor // Key: "feature/TypeName" -> Feature executor instance
	types            map[string]reflect.Type    // Key: "feature/TypeName" -> Command type
	mu               sync.RWMutex
}

// NewCommandRegistry creates a new, initialized CommandRegistry.
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers:         make(map[string]CommandHandler),
		featureExecutors: make(map[string]FeatureExecutor),
		types:            make(map[string]reflect.Type),
	}
}

// Register associates a command type with its state update handler and the responsible
// feature executor instance using its CommandName().
// It panics if a handler or executor for the name is already registered.
func (r *CommandRegistry) Register(cmd Command, handler CommandHandler, featureExecutor FeatureExecutor) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := cmd.CommandName() // Use CommandName()
	if _, exists := r.handlers[name]; exists {
		panic(fmt.Sprintf("handler already registered for command name %q", name))
	}
	if _, exists := r.featureExecutors[name]; exists {
		// This shouldn't happen if Register is called correctly once per command type
		panic(fmt.Sprintf("feature executor already registered for command name %q", name))
	}
	if handler == nil {
		panic(fmt.Sprintf("attempted to register nil handler for command name %q", name))
	}
	if featureExecutor == nil {
		panic(fmt.Sprintf("attempted to register nil feature executor for command name %q", name))
	}

	cmdType := reflect.TypeOf(cmd)
	// Ensure we store the non-pointer type for consistency if needed
	if cmdType.Kind() == reflect.Ptr {
		cmdType = cmdType.Elem()
	}

	r.handlers[name] = handler
	r.featureExecutors[name] = featureExecutor
	r.types[name] = cmdType // Store the type for lookup
	slog.Debug("Registered command handler and feature executor", "name", name, "type", cmdType)
}

// GetHandler retrieves the registered state update handler for a given command name.
func (r *CommandRegistry) GetHandler(name string) (CommandHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, found := r.handlers[name]
	return handler, found
}

// GetFeatureExecutor retrieves the registered feature executor instance for a given command name.
func (r *CommandRegistry) GetFeatureExecutor(name string) (FeatureExecutor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	executor, found := r.featureExecutors[name]
	return executor, found
}

// GetHandlerAndFeatureExecutor retrieves both the handler and feature executor for a given command name.
func (r *CommandRegistry) GetHandlerAndFeatureExecutor(name string) (CommandHandler, FeatureExecutor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, handlerFound := r.handlers[name]
	executor, executorFound := r.featureExecutors[name]
	// Both must be found for the registration to be considered complete
	found := handlerFound && executorFound
	if found {
		return handler, executor, true
	}
	// Log if one is found but not the other, indicating an incomplete registration
	if handlerFound != executorFound {
		slog.Error("Inconsistent registration found for command", "name", name, "handlerFound", handlerFound, "executorFound", executorFound)
	}
	return nil, nil, false
}

// RegisteredCommandNames returns a slice of strings containing the full names
// (e.g., "feature/TypeName") of all registered command types.
func (r *CommandRegistry) RegisteredCommandNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		names = append(names, name)
	}
	// Sort for predictable output order
	// sort.Strings(names) // Optional: uncomment if consistent order is desired
	return names
}

// GetCommandType retrieves the reflect.Type for a registered command by its full name.
// This is useful for decoding commands from external sources like API requests.
func (r *CommandRegistry) GetCommandType(name string) (reflect.Type, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Direct lookup using the stored type map
	cmdType, found := r.types[name]
	return cmdType, found
}

// --- Global Registry (Optional - consider dependency injection instead) ---
// var Commands = NewCommandRegistry()

// InitRegistries initializes global registries (if used).
// Consider using dependency injection instead of globals.
// func InitRegistries() {
// 	Commands = NewCommandRegistry()
// 	Queries = NewQueryRegistry() // Assuming Queries registry exists
// }

// --- Central Command Executor ---

// Executor orchestrates the validation, logging, and execution of commands.
type Executor struct {
	log      *MessageLog      // Dependency for appending commands
	registry *CommandRegistry // Dependency for finding handlers and feature executors
}

// NewExecutor creates a new central command executor.
func NewExecutor(log *MessageLog, registry *CommandRegistry) *Executor {
	if log == nil {
		panic("MessageLog cannot be nil for Executor")
	}
	if registry == nil {
		panic("CommandRegistry cannot be nil for Executor")
	}
	return &Executor{
		log:      log,
		registry: registry,
	}
}

// Execute orchestrates the full lifecycle of a command:
// 1. Retrieves the state update handler and the responsible feature executor.
// 2. Calls the feature executor's ValidateCommand method.
// 3. Appends the command to the message log.
// 4. Executes the state update handler.
// It returns an error if validation or logging fails.
// It panics if the state update handler returns an error after the command has been logged,
// indicating an unrecoverable inconsistency.
func (e *Executor) Execute(ctx context.Context, cmd Command) error {
	name := cmd.CommandName()
	slog.Debug("Executing command", "name", name)

	// 1. Get Handler and Feature Executor
	handler, featureExecutor, found := e.registry.GetHandlerAndFeatureExecutor(name)
	if !found {
		slog.Error("No handler and/or feature executor registered for command", "name", name)
		return fmt.Errorf("command %q not registered", name)
	}

	// 2. Validate Command using Feature Executor
	slog.Debug("Validating command", "name", name)
	if err := featureExecutor.ValidateCommand(ctx, cmd); err != nil {
		slog.Warn("Command validation failed", "name", name, "error", err)
		// TODO: Consider defining specific validation error types to return distinct HTTP status codes (e.g., 400 Bad Request)
		return fmt.Errorf("validation failed for command %q: %w", name, err)
	}
	slog.Debug("Command validation successful", "name", name)

	// 3. Append Command to Log
	slog.Debug("Appending command to log", "name", name)
	if err := e.log.Append(ctx, cmd); err != nil {
		slog.Error("Failed to append command to log", "name", name, "error", err)
		// This is a critical error, as the action wasn't persisted.
		return fmt.Errorf("failed to persist command %q: %w", name, err)
	}
	slog.Debug("Command appended to log successfully", "name", name)

	// 4. Execute State Update Handler
	slog.Debug("Executing state update handler", "name", name)
	handlerErr := handler(ctx, cmd)
	if handlerErr != nil {
		// PANIC! If the handler fails *after* the command was logged,
		// the state is inconsistent with the log. This is unrecoverable
		// without manual intervention or complex compensation logic.
		// A panic forces a restart, allowing state to be rebuilt from the log.
		slog.Error("State update handler failed after command was logged! PANICKING.", "name", name, "error", handlerErr)
		panic(fmt.Sprintf("unrecoverable state inconsistency: handler for %q failed after logging: %v", name, handlerErr))
	}
	slog.Debug("State update handler executed successfully", "name", name)

	slog.Info("Command executed successfully", "name", name)
	return nil
}
