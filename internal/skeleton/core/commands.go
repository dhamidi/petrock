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
	NamedMessage
}

// CommandHandler defines the function signature for handling commands.
type CommandHandler func(ctx context.Context, cmd Command) error

// CommandRegistry maps command names (feature/Type) to their handlers and types.
type CommandRegistry struct {
	handlers map[string]CommandHandler // Key: "feature/TypeName"
	types    map[string]reflect.Type   // Key: "feature/TypeName"
	mu       sync.RWMutex
}

// NewCommandRegistry creates a new, initialized CommandRegistry.
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]CommandHandler),
		types:    make(map[string]reflect.Type),
	}
}

// Register associates a command type with its handler using its RegisteredName().
// It panics if a handler for the name is already registered.
func (r *CommandRegistry) Register(cmd Command, handler CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := cmd.RegisteredName()
	if _, exists := r.handlers[name]; exists {
		panic(fmt.Sprintf("handler already registered for command name %q", name))
	}

	cmdType := reflect.TypeOf(cmd)
	// Ensure we store the non-pointer type for consistency if needed
	if cmdType.Kind() == reflect.Ptr {
		cmdType = cmdType.Elem()
	}

	r.handlers[name] = handler
	r.types[name] = cmdType // Store the type for lookup
	slog.Debug("Registered command handler", "name", name, "type", cmdType)
}

// Dispatch finds the handler for the given command's RegisteredName() and executes it.
// It returns an error if no handler is found or if the handler returns an error.
func (r *CommandRegistry) Dispatch(ctx context.Context, cmd Command) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	name := cmd.RegisteredName()
	handler, exists := r.handlers[name]
	if !exists {
		return fmt.Errorf("no command handler registered for name %q (type %T)", name, cmd)
	}

	slog.Debug("Dispatching command", "name", name, "type", reflect.TypeOf(cmd))
	return handler(ctx, cmd)
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
