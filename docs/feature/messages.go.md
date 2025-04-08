# Plan for posts/messages.go (Example Feature)

This file defines the commands, queries, and query results for the posts feature. Commands now include validation methods to support the centralized executor pattern.

## Types

### Commands (Must implement core.Command and core.Validator)

- `CreatePostCommand`: Represents a command to create a new post.
    - `Title string`: The title of the post. (Required)
    - `Content string`: The content of the post. (Required)
    - `Tags []string`: Optional tags for the post.
    - `AuthorID string`: The ID of the author creating the post. (Required)
    - `func (c CreatePostCommand) CommandName() string`: Returns "posts/create".
    - `func (c CreatePostCommand) Validate() error`: Validates required fields (Title, Content, AuthorID).

- `UpdatePostCommand`: Represents a command to update an existing post.
    - `PostID string`: The ID of the post to update. (Required)
    - `Title string`: The updated title. (Required)
    - `Content string`: The updated content. (Required)
    - `Tags []string`: Updated tags.
    - `func (c UpdatePostCommand) CommandName() string`: Returns "posts/update".
    - `func (c UpdatePostCommand) Validate() error`: Validates required fields (PostID, Title, Content).

- `DeletePostCommand`: Represents a command to delete an existing post.
    - `PostID string`: The ID of the post to delete. (Required)
    - `func (c DeletePostCommand) CommandName() string`: Returns "posts/delete".
    - `func (c DeletePostCommand) Validate() error`: Validates that PostID is provided.

### Queries (Must implement core.Query)

- `GetPostQuery`: Represents a query to fetch a single post by ID.
    - `PostID string`: The ID of the post to retrieve. (Required)
    - `func (q GetPostQuery) QueryName() string`: Returns "posts/get".
    - `func (q GetPostQuery) Validate() error`: Validates that PostID is provided.

- `ListPostsQuery`: Represents a query to fetch multiple posts with optional filtering and pagination.
    - `Page int`: The page number for pagination. Defaults to 1.
    - `PageSize int`: The number of items per page. Defaults to 10.
    - `Tags []string`: Optional filter by tags.
    - `AuthorID string`: Optional filter by author.
    - `func (q ListPostsQuery) QueryName() string`: Returns "posts/list".
    - `func (q ListPostsQuery) Validate() error`: Validates that Page and PageSize are positive.

### Query Results (Must implement core.QueryResult)

- `Post`: Represents a single post in query results.
    - `ID string`: The unique identifier of the post.
    - `Title string`: The title of the post.
    - `Content string`: The content of the post.
    - `Tags []string`: Tags associated with the post.
    - `AuthorID string`: The ID of the post's author.
    - `CreatedAt time.Time`: When the post was created.
    - `UpdatedAt time.Time`: When the post was last updated.

- `PostListResult`: Represents the result of a list query.
    - `Posts []Post`: The slice of posts matching the query criteria.
    - `TotalCount int`: The total number of matching posts (for pagination).
    - `Page int`: The current page number.
    - `PageSize int`: The number of items per page.