# Plan for posts/execute.go (Example Feature)

This file defines the feature-specific **Executor**. This component holds the feature's state and provides the state update handlers. It also implements the `core.FeatureExecutor` interface to bridge validation calls from the central `core.Executor` to the command structs themselves.

## Types

- `Executor`: A struct holding dependencies needed for state updates, primarily the feature's state. It implements `core.FeatureExecutor`.
    - `state *PostState`
    - *Other dependencies if needed.*
- `Validator` (Interface, potentially defined in `messages.go`): An interface implemented by command structs that require stateful validation.
    - `Validate(state *PostState) error`

## Functions

- `NewExecutor(state *PostState) *Executor`: Constructor for the feature `Executor`.

### Validation Bridge Method (Implements `core.FeatureExecutor`)

- `(e *Executor) ValidateCommand(ctx context.Context, cmd core.Command) error`: This method is called by the central `core.Executor`. It checks if the received command `cmd` implements the feature's `Validator` interface. If it does, it calls the command's `Validate` method, passing its own state (`e.state`).
    ```go
    // Check if the command implements the stateful validator interface
    if validator, ok := cmd.(Validator); ok {
        // If yes, call the command's Validate method with the feature state
        return validator.Validate(e.state)
    }
    // If the command doesn't implement Validator, assume no stateful validation needed
    return nil
    ```
    *Note: Basic, stateless validation (e.g., checking required fields) can still be done within the command's `Validate` method or potentially before calling `core.Executor.Execute` in the HTTP handler.*

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
