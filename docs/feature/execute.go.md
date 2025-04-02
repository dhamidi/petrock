# Plan for posts/execute.go (Example Feature)

This file defines the command handlers responsible for executing state changes within the feature. These handlers interact with the feature's state.

## Types

- `PostExecutor`: A struct that holds dependencies needed for command execution, primarily the feature's state.
    - `state *PostState`

## Functions

- `NewPostExecutor(state *PostState) *PostExecutor`: Constructor for the `PostExecutor`.
- `(e *PostExecutor) HandleCreatePost(ctx context.Context, cmd core.Command) error`: Handles the `CreatePostCommand`.
    - Type-assert `cmd` to `CreatePostCommand`.
    - Perform validation (e.g., non-empty title/content).
    - Create a new `Post` object within the `e.state`.
    - Return `nil` on success, or an error on validation failure or state update issues.
    *Note: This function signature matches `core.CommandHandler`.*
- `(e *PostExecutor) HandleUpdatePost(ctx context.Context, cmd core.Command) error`: Handles the `UpdatePostCommand`.
    - Type-assert `cmd` to `UpdatePostCommand`.
    - Validate input.
    - Find the existing post in `e.state` using `PostID`.
    - Update the post's fields in `e.state`.
    - Return error if post not found or validation fails.
    *Note: This function signature matches `core.CommandHandler`.*
- `(e *PostExecutor) HandleDeletePost(ctx context.Context, cmd core.Command) error`: Handles the `DeletePostCommand`.
    - Type-assert `cmd` to `DeletePostCommand`.
    - Remove the post from `e.state` using `PostID`.
    - Return error if post not found.
    *Note: This function signature matches `core.CommandHandler`.*
