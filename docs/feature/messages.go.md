# Plan for posts/messages.go (Example Feature)

This file defines the specific command and query message structures, as well as query result structures, for the feature. These structs represent the data transferred for operations within this feature.

## Types

### Commands (Implement `core.Command`)
- `CreatePostCommand`:
    - `Title string`
    - `Content string`
    - `AuthorID string` // Or appropriate user identifier type
- `UpdatePostCommand`:
    - `PostID string` // Identifier for the post to update
    - `Title string`
    - `Content string`
- `DeletePostCommand`:
    - `PostID string` // Identifier for the post to delete

### Queries (Implement `core.Query`)
- `GetPostQuery`:
    - `PostID string` // Identifier for the post to retrieve
- `ListPostsQuery`:
    - `Page int` // For pagination
    - `PageSize int` // For pagination
    - `AuthorIDFilter string` // Optional filter criteria

### Query Results (Implement `core.QueryResult`)
- `PostQueryResult`: Represents a single post's data.
    - `ID string`
    - `Title string`
    - `Content string`
    - `AuthorID string`
    - `CreatedAt time.Time`
    - `UpdatedAt time.Time`
- `PostsListQueryResult`: Represents a list of posts, potentially with pagination info.
    - `Posts []PostQueryResult`
    - `TotalCount int`
    - `Page int`
    - `PageSize int`

## Functions

- None, this file primarily defines data structures (structs).
