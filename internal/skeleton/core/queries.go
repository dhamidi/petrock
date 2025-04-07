package core

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
)

// Query is an interface combining NamedMessage for query messages.
type Query interface {
	NamedMessage
}

// QueryResult is a marker interface for the data returned by a query handler.
type QueryResult interface{}

// QueryHandler defines the function signature for handling queries.
type QueryHandler func(ctx context.Context, query Query) (QueryResult, error)

// QueryRegistry maps query names (feature/Type) to their handlers and types.
type QueryRegistry struct {
	handlers map[string]QueryHandler // Key: "feature/TypeName"
	types    map[string]reflect.Type // Key: "feature/TypeName"
	mu       sync.RWMutex
}

// NewQueryRegistry creates a new, initialized QueryRegistry.
func NewQueryRegistry() *QueryRegistry {
	return &QueryRegistry{
		handlers: make(map[string]QueryHandler),
		types:    make(map[string]reflect.Type),
	}
}

// Register associates a query type with its handler using its RegisteredName().
// It panics if a handler for the name is already registered.
func (r *QueryRegistry) Register(query Query, handler QueryHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := query.RegisteredName()
	if _, exists := r.handlers[name]; exists {
		panic(fmt.Sprintf("handler already registered for query name %q", name))
	}

	queryType := reflect.TypeOf(query)
	// Ensure we store the non-pointer type for consistency if needed
	if queryType.Kind() == reflect.Ptr {
		queryType = queryType.Elem()
	}

	r.handlers[name] = handler
	r.types[name] = queryType // Store the type for lookup
	slog.Debug("Registered query handler", "name", name, "type", queryType)
}

// Dispatch finds the handler for the given query's RegisteredName() and executes it.
// It returns the result and an error if no handler is found or if the handler returns an error.
func (r *QueryRegistry) Dispatch(ctx context.Context, query Query) (QueryResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	name := query.RegisteredName()
	handler, exists := r.handlers[name]
	if !exists {
		return nil, fmt.Errorf("no query handler registered for name %q (type %T)", name, query)
	}

	slog.Debug("Dispatching query", "name", name, "type", reflect.TypeOf(query))
	return handler(ctx, query)
}

// RegisteredQueryNames returns a slice of strings containing the full names
// (e.g., "feature/TypeName") of all registered query types.
func (r *QueryRegistry) RegisteredQueryNames() []string {
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

// GetQueryType retrieves the reflect.Type for a registered query by its full name.
// This is useful for decoding/constructing queries from external sources like API requests.
func (r *QueryRegistry) GetQueryType(name string) (reflect.Type, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Direct lookup using the stored type map
	queryType, found := r.types[name]
	return queryType, found
}

// --- Global Registry (Optional - consider dependency injection instead) ---
// var Queries = NewQueryRegistry()
