# Plan for posts/query.go (Example Feature)

This file defines the query handlers responsible for retrieving data (reading state) from the feature.

## Types

- `PostQuerier`: A struct that holds dependencies needed for query execution, primarily the feature's state.
    - `state *PostState`

## Functions

- `NewPostQuerier(state *PostState) *PostQuerier`: Constructor for the `PostQuerier`.
- `(q *PostQuerier) HandleGetPost(ctx context.Context, query core.Query) (core.QueryResult, error)`: Handles the `GetPostQuery`.
    - Type-assert `query` to `GetPostQuery`.
    - Retrieve the post from `q.state` using `PostID`.
    - If found, map the internal `Post` state struct to a `PostQueryResult` struct.
    - Return the `PostQueryResult` and `nil` error.
    - If not found, return `nil` and an appropriate error (e.g., `ErrPostNotFound`).
    *Note: This function signature matches `core.QueryHandler`.*
- `(q *PostQuerier) HandleListPosts(ctx context.Context, query core.Query) (core.QueryResult, error)`: Handles the `ListPostsQuery`.
    - Type-assert `query` to `ListPostsQuery`.
    - Retrieve the list of posts from `q.state`, applying filtering (e.g., by `AuthorIDFilter`) and pagination (`Page`, `PageSize`).
    - Map the internal `Post` state structs to `PostQueryResult` structs.
    - Construct a `PostsListQueryResult` containing the list and pagination details (total count, current page, page size).
    - Return the `PostsListQueryResult` and `nil` error.
    *Note: This function signature matches `core.QueryHandler`.*
