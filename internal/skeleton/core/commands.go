package core

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
)

// Command is a marker interface for command messages.
type Command interface{}

// CommandHandler defines the function signature for handling commands.
type CommandHandler func(ctx context.Context, cmd Command) error

// CommandRegistry maps command types to their handlers.
type CommandRegistry struct {
	handlers map[reflect.Type]CommandHandler
	mu       sync.RWMutex
}

// NewCommandRegistry creates a new, initialized CommandRegistry.
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[reflect.Type]CommandHandler),
	}
}

// Register associates a command type with its handler.
// It panics if a handler for the command type is already registered.
func (r *CommandRegistry) Register(cmd Command, handler CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cmdType := reflect.TypeOf(cmd)
	if _, exists := r.handlers[cmdType]; exists {
		panic(fmt.Sprintf("handler already registered for command type %v", cmdType))
	}

	r.handlers[cmdType] = handler
	slog.Debug("Registered command handler", "type", cmdType)
}

// Dispatch finds the handler for the given command's type and executes it.
// It returns an error if no handler is found or if the handler returns an error.
func (r *CommandRegistry) Dispatch(ctx context.Context, cmd Command) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmdType := reflect.TypeOf(cmd)
	handler, exists := r.handlers[cmdType]
	if !exists {
		return fmt.Errorf("no command handler registered for type %v", cmdType)
	}

	slog.Debug("Dispatching command", "type", cmdType)
	// TODO: Add instrumentation/tracing here if needed
	return handler(ctx, cmd)
}

// RegisteredCommandNames returns a slice of strings containing the names
// of all registered command types.
func (r *CommandRegistry) RegisteredCommandNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.handlers))
	for cmdType := range r.handlers {
		names = append(names, cmdType.Name()) // Use the simple type name
	}
	// Sort for predictable output order
	// sort.Strings(names) // Optional: uncomment if consistent order is desired
	return names
}

// GetCommandType retrieves the reflect.Type for a registered command by its name.
// This is useful for decoding commands from external sources like API requests.
func (r *CommandRegistry) GetCommandType(name string) (reflect.Type, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for cmdType := range r.handlers {
		// Compare the simple name of the registered type
		if cmdType.Name() == name {
			return cmdType, true
		}
	}
	return nil, false
}

// --- Global Registry (Optional - consider dependency injection instead) ---
// var Commands = NewCommandRegistry()

// InitRegistries initializes global registries (if used).
// Consider using dependency injection instead of globals.
// func InitRegistries() {
// 	Commands = NewCommandRegistry()
// 	Queries = NewQueryRegistry() // Assuming Queries registry exists
// }
