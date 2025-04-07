package core

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
)

// Query is a marker interface for query messages.
type Query interface{}

// QueryResult is a marker interface for the data returned by a query handler.
type QueryResult interface{}

// QueryHandler defines the function signature for handling queries.
type QueryHandler func(ctx context.Context, query Query) (QueryResult, error)

// QueryRegistry maps query types to their handlers.
type QueryRegistry struct {
	handlers map[reflect.Type]QueryHandler
	mu       sync.RWMutex
}

// NewQueryRegistry creates a new, initialized QueryRegistry.
func NewQueryRegistry() *QueryRegistry {
	return &QueryRegistry{
		handlers: make(map[reflect.Type]QueryHandler),
	}
}

// Register associates a query type with its handler.
// It panics if a handler for the query type is already registered.
func (r *QueryRegistry) Register(query Query, handler QueryHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	queryType := reflect.TypeOf(query)
	if _, exists := r.handlers[queryType]; exists {
		panic(fmt.Sprintf("handler already registered for query type %v", queryType))
	}

	r.handlers[queryType] = handler
	slog.Debug("Registered query handler", "type", queryType)
}

// Dispatch finds the handler for the given query's type and executes it.
// It returns the result and an error if no handler is found or if the handler returns an error.
func (r *QueryRegistry) Dispatch(ctx context.Context, query Query) (QueryResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	queryType := reflect.TypeOf(query)
	handler, exists := r.handlers[queryType]
	if !exists {
		return nil, fmt.Errorf("no query handler registered for type %v", queryType)
	}

	slog.Debug("Dispatching query", "type", queryType)
	return handler(ctx, query)
}

// RegisteredQueryNames returns a slice of strings containing the names
// of all registered query types.
func (r *QueryRegistry) RegisteredQueryNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.handlers))
	for queryType := range r.handlers {
		names = append(names, queryType.Name()) // Use the simple type name
	}
	// Sort for predictable output order
	// sort.Strings(names) // Optional: uncomment if consistent order is desired
	return names
}

// GetQueryType retrieves the reflect.Type for a registered query by its name.
// This is useful for decoding/constructing queries from external sources like API requests.
func (r *QueryRegistry) GetQueryType(name string) (reflect.Type, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for queryType := range r.handlers {
		// Compare the simple name of the registered type
		if queryType.Name() == name {
			return queryType, true
		}
	}
	return nil, false
}

// --- Global Registry (Optional - consider dependency injection instead) ---
// var Queries = NewQueryRegistry()
