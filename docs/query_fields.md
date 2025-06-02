# Query Fields Support

This document describes the enhanced query generator that supports custom typed fields, similar to the command field support.

## Overview

The query generator now supports custom field definitions for both:
1. Query struct fields (for query input parameters)
2. ItemResult struct fields (for query output data)

## Usage

### Programmatic Usage

```go
// Define custom fields
fields := []QueryField{
    {Name: "postID", Type: "string"},
    {Name: "status", Type: "string"},
    {Name: "publishedAfter", Type: "time.Time"},
}

// Generate query with custom fields
queryGen := NewQueryGenerator(targetDir)
err := queryGen.GenerateQueryComponentWithFields(
    "posts",     // featureName
    "search",    // entityName
    targetDir,   // targetDir
    modulePath,  // modulePath
    fields,      // custom fields
)
```

### Generated Structure

#### Query Struct
When custom fields are provided, the query struct is modified to include only the specified fields:

```go
type SearchQuery struct {
    PostID        string    `json:"postID" validate:"required"`
    Status        string    `json:"status" validate:"required"`
    PublishedAfter time.Time `json:"publishedAfter" validate:"required"`
}
```

#### ItemResult Struct
The ItemResult struct in base.go is also updated to match the custom fields:

```go
type ItemResult struct {
    PostID        string    `json:"postID"`
    Status        string    `json:"status"`
    PublishedAfter time.Time `json:"publishedAfter"`
}
```

#### Handler Method
Handler methods are simplified to stub implementations when custom fields are used:

```go
func (q *Querier) HandleSearch(ctx context.Context, query core.Query) (core.QueryResult, error) {
    return nil, nil
}
```

## Implementation Details

### Editor-Based Modification
The implementation uses the `internal/ed` editor package to:
1. Extract standard query templates
2. Replace struct field definitions with custom fields
3. Simplify handler methods to stub implementations

### Field Naming
- Input field names are automatically capitalized for Go export (e.g., `postID` → `PostID`)
- JSON tags maintain the original field names for API compatibility
- Validation tags are added automatically for query fields

### Template Modifications
Two files are modified when custom fields are provided:
1. `{feature}/queries/{entity}.go` - Query struct and handler
2. `{feature}/queries/base.go` - ItemResult struct

## Equivalence with Command Fields
This implementation provides equivalent functionality to the command field support added in commit 8751560:

- **Field Definition**: Same `FieldType` struct pattern with Name and Type
- **Editor Usage**: Same editor-based template modification approach
- **Method Simplification**: Similar stub implementation generation
- **Struct Replacement**: Complete field replacement rather than addition

## Testing

The query field functionality can be tested using:
- Unit tests for field modification logic
- Integration tests for complete generation workflow
- Build verification to ensure generated code compiles correctly
