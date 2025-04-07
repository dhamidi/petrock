# Plan for posts/messages.go (Example Feature)

This file defines the specific command and query message structures, as well as query result structures, for the feature. These structs represent the data transferred for operations within this feature.

## Types

*The `petrock feature` command automatically generates `CommandName()` and `QueryName()` methods for command and query structs respectively, returning the kebab-case name (e.g., `posts/create-post-command`).*

### Commands (Implement `core.Command`)
- `CreatePostCommand`: (Method `CommandName() string` generated)
    - `Title string`
    - `Content string`
    - `AuthorID string` // Or appropriate user identifier type
- `UpdatePostCommand`:
    - `PostID string` // Identifier for the post to update
    - `Title string`
    - `Content string`
- `DeletePostCommand`: (Method `CommandName() string` generated)
    - `PostID string` // Identifier for the post to delete

### Queries (Implement `core.Query`)
- `GetPostQuery`: (Method `QueryName() string` generated)
    - `PostID string` // Identifier for the post to retrieve
- `ListPostsQuery`: (Method `QueryName() string` generated)
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
