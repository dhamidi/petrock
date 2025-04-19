# Queries

The `queries` directory contains query definitions and implementations for the feature, following the Command-Query Separation pattern for retrieving data without modifying state.

## Structure

- `base.go` - Common query interfaces and types
- `get.go` - Query for retrieving single items
- `list.go` - Query for retrieving lists of items

## Query Pattern

Queries follow a pattern where:

1. Each query is a distinct type with its own parameters
2. Queries are executed by a query handler
3. Queries return results without modifying state
4. Queries may be optimized for specific use cases

## Key Components

- Query types that define the request parameters
- Result types that structure the returned data
- Query handlers that process the queries

## Usage

Queries are typically called from handlers to retrieve data needed for rendering views or responding to API requests.