# Plan for posts/queries.go (Example Feature)

This file defines the query structures and query result types for the feature.

## Types

*Query structs in the template implement `QueryName()`. These methods return the kebab-case name (e.g., `petrock_example_feature_name/get`) which is updated with the actual feature name during placeholder replacement.*

### Queries (Implement `core.Query`)
- `GetPostQuery`: (Implements `QueryName() string`)
    - `PostID string`
- `ListPostsQuery`: (Implements `QueryName() string`)
    - `Page int`
    - `PageSize int`
    - `AuthorIDFilter string`

### Query Results (Implement `core.QueryResult`)
- `PostQueryResult`:
    - `ID string`
    - `Title string`
    - `Content string`
    - `AuthorID string`
    - `CreatedAt time.Time`
    - `UpdatedAt time.Time`
- `PostsListQueryResult`:
    - `Posts []PostQueryResult`
    - `TotalCount int`
    - `Page int`
    - `PageSize int`

## Functions

- None, this file primarily defines data structures for queries and their results.