# Plan for posts/execute.go (Example Feature)

This file defines the feature-specific **Executor**. This component is responsible for both **validating** commands (using the feature's state) and **applying state changes** for those commands.

## Types

- `Executor`: A struct holding dependencies needed for command validation and state updates, primarily the feature's state. It implements `core.CommandValidator`.
    - `state *PostState`
    - *Other dependencies if needed (e.g., external service clients for validation).*

## Functions

- `NewExecutor(state *PostState) *Executor`: Constructor for the feature `Executor`.

### Validation Methods (Implement `core.CommandValidator` implicitly via `ValidateCommand`)

- `(e *Executor) ValidateCommand(ctx context.Context, cmd core.Command) error`: The central validation method called by `core.Executor`. It type-switches on the command and delegates to specific validation methods.
    ```go
    switch c := cmd.(type) {
    case CreatePostCommand:
        return e.validateCreatePost(ctx, c)
    case UpdatePostCommand:
        return e.validateUpdatePost(ctx, c)
    // ... other commands
    default:
        return fmt.Errorf("unknown command type for validation: %T", cmd)
    }
    ```
- `(e *Executor) validateCreatePost(ctx context.Context, cmd CreatePostCommand) error`: Performs validation specific to `CreatePostCommand`.
    - Checks basic field constraints (e.g., non-empty title).
    - Checks state-dependent constraints using `e.state` (e.g., uniqueness of a slug derived from the title, if applicable).
    - Returns `nil` if valid, or an error otherwise.
- `(e *Executor) validateUpdatePost(ctx context.Context, cmd UpdatePostCommand) error`: Performs validation specific to `UpdatePostCommand`.
    - Checks basic field constraints.
    - Checks if the post with `cmd.PostID` exists in `e.state`.
    - Returns `nil` if valid, or an error otherwise.
- `(e *Executor) validateDeletePost(ctx context.Context, cmd DeletePostCommand) error`: Performs validation specific to `DeletePostCommand`.
    - Checks if the post with `cmd.PostID` exists in `e.state`.
    - Returns `nil` if valid, or an error otherwise.

### State Update Handlers (Match `core.CommandHandler` signature)

- `(e *Executor) HandleCreatePost(ctx context.Context, cmd core.Command) error`: Applies state changes for `CreatePostCommand`.
    - Type-assert `cmd` to `CreatePostCommand`. **No validation here.**
    - Create and add the new `Post` object to `e.state`.
    - Return `nil` on success. An error returned here will cause `core.Executor` to panic.
    *Note: This function is registered with `core.CommandRegistry` as the state update handler and called by `core.Executor` after validation and logging, and also during log replay.*
- `(e *Executor) HandleUpdatePost(ctx context.Context, cmd core.Command) error`: Applies state changes for `UpdatePostCommand`.
    - Type-assert `cmd` to `UpdatePostCommand`. **No validation here.**
    - Find the existing post in `e.state` using `PostID`. If not found, return an error (will cause panic).
    - Update the post's fields in `e.state`.
    - Return `nil` on success. An error returned here will cause `core.Executor` to panic.
    *Note: This function signature matches `core.CommandHandler`.*
- `(e *Executor) HandleDeletePost(ctx context.Context, cmd core.Command) error`: Applies state changes for `DeletePostCommand`.
    - Type-assert `cmd` to `DeletePostCommand`. **No validation here.**
    - Remove (or mark as deleted) the post from `e.state` using `PostID`. If not found, return an error (will cause panic).
    - Return `nil` on success. An error returned here will cause `core.Executor` to panic.
    *Note: This function signature matches `core.CommandHandler`.*
