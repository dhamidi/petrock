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
- `(r *QueryRegistry) Register(query Query, handler QueryHandler)`: Registers a query handler for a specific query type. It uses reflection (`reflect.TypeOf(query)`) to get the type key.
- `(r *QueryRegistry) Dispatch(ctx context.Context, query Query) (QueryResult, error)`: Looks up the handler for the given query's type and executes it. Returns the result and an error if no handler is registered or if the handler itself returns an error.
