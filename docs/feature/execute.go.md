# Plan for posts/execute.go (Example Feature)

This file defines the command handlers responsible for executing domain-specific logic within the feature. With the centralized executor pattern, these handlers focus on business logic rather than common execution flow concerns like validation and logging.

## Types

- `PostExecutor`: A struct that holds dependencies needed for command execution, primarily the feature's state.
    - `state *PostState`

## Functions

- `NewPostExecutor(state *PostState) *PostExecutor`: Constructor for the `PostExecutor`.

- `(e *PostExecutor) HandleCreatePost(ctx context.Context, cmd core.Command) error`: Handles the `CreatePostCommand`.
    - Type-assert `cmd` to `CreatePostCommand`.
    - Create a new `Post` object within the `e.state`.
    - Return `nil` on success, or an error on state update issues.
    *Note: This function focuses solely on state manipulation. Validation and persistence are handled by the core.Executor.*

- `(e *PostExecutor) HandleUpdatePost(ctx context.Context, cmd core.Command) error`: Handles the `UpdatePostCommand`.
    - Type-assert `cmd` to `UpdatePostCommand`.
    - Find the existing post in `e.state` using `PostID`.
    - Update the post's fields in `e.state`.
    - Return error if post not found.
    *Note: This function focuses solely on state manipulation.*

- `(e *PostExecutor) HandleDeletePost(ctx context.Context, cmd core.Command) error`: Handles the `DeletePostCommand`.
    - Type-assert `cmd` to `DeletePostCommand`.
    - Remove the post from `e.state` using `PostID`.
    - Return error if post not found.
    *Note: This function focuses solely on state manipulation.*

## Command Validation

Validation logic is now implemented in the command types rather than in the handlers:

- `(c CreatePostCommand) Validate() error`: Validates command fields (e.g., title not empty).
- `(c UpdatePostCommand) Validate() error`: Validates command fields and ensures ID is provided.
- `(c DeletePostCommand) Validate() error`: Ensures ID is provided.

These validation methods are automatically called by the core.Executor before the command is logged or dispatched to handlers.
