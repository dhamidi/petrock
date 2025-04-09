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
