# Plan for core/queries.go

This file defines the registry for queries and their associated handlers, enabling read operations on the application state.

## Types

- `Query`: An interface representing a query message. All query structs should implicitly satisfy this (e.g., by being `interface{}`). It's primarily a marker.
- `QueryResult`: An interface representing the data returned by a query handler. All query result structs should implicitly satisfy this (e.g., by being `interface{}`).
- `QueryHandler func(ctx context.Context, query Query) (QueryResult, error)`: A function type for handlers that process queries. Takes context and the query message, returns a result and an error.
- `QueryRegistry`: A struct responsible for mapping query types to their handlers.
    - `handlers map[reflect.Type]QueryHandler`: The internal map storing the registrations.
    - `mu sync.RWMutex`: For thread-safe access to the handlers map.

## Functions

- `NewQueryRegistry() *QueryRegistry`: Constructor function to create and initialize a new `QueryRegistry`.
- `Query`: Interface that query structs must implement. Requires `QueryName() string`.
- `QueryResult`: Marker interface for query results.
- `QueryHandler`: Function type for query handlers.
- `QueryRegistry`: Maps query names (`feature/kebab-case-name`) to handlers and `reflect.Type`.
- `NewQueryRegistry()`: Constructor.
- `(r *QueryRegistry) Register(query Query, handler QueryHandler)`: Registers a handler using the name returned by `query.QueryName()`. Stores the handler and `reflect.Type`. Panics if the name is already registered.
- `(r *QueryRegistry) Dispatch(ctx context.Context, query Query) (QueryResult, error)`: Looks up the handler using `query.QueryName()` and executes it.
- `(r *QueryRegistry) RegisteredQueryNames() []string`: Returns a slice containing the full registered kebab-case names (e.g., "posts/list-query") of all queries.
- `(r *QueryRegistry) GetQueryType(name string) (reflect.Type, bool)`: Looks up and returns the `reflect.Type` for a query based on its full registered kebab-case name (e.g., "posts/list-query").
