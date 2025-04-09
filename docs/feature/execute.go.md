# Plan for posts/execute.go (Example Feature)

This file defines the **state update handlers** for the feature's commands. These handlers are responsible *only* for applying the changes described by a command to the feature's state (`posts.State`). They assume the command has already been validated and logged by the `core.Executor`.

## Types

- `StateUpdater`: A struct holding dependencies needed for state updates, primarily the feature's state.
    - `state *PostState`
    - *Other dependencies if state updates require them (rare).*

## Functions

- `NewStateUpdater(state *PostState) *StateUpdater`: Constructor for the `StateUpdater`.
- `(su *StateUpdater) HandleCreatePost(ctx context.Context, cmd core.Command) error`: Applies state changes for `CreatePostCommand`.
    - Type-assert `cmd` to `CreatePostCommand`. **No validation here.**
    - Create and add the new `Post` object to `su.state`.
    - Return `nil` on success. An error returned here will cause `core.Executor` to panic.
    *Note: This function signature matches `core.CommandHandler`. It's registered with the `core.CommandRegistry` and called by `core.Executor` after logging, and also during log replay.*
- `(su *StateUpdater) HandleUpdatePost(ctx context.Context, cmd core.Command) error`: Applies state changes for `UpdatePostCommand`.
    - Type-assert `cmd` to `UpdatePostCommand`. **No validation here.**
    - Find the existing post in `su.state` using `PostID`. If not found, return an error (will cause panic).
    - Update the post's fields in `su.state`.
    - Return `nil` on success. An error returned here will cause `core.Executor` to panic.
    *Note: This function signature matches `core.CommandHandler`.*
- `(su *StateUpdater) HandleDeletePost(ctx context.Context, cmd core.Command) error`: Applies state changes for `DeletePostCommand`.
    - Type-assert `cmd` to `DeletePostCommand`. **No validation here.**
    - Remove (or mark as deleted) the post from `su.state` using `PostID`. If not found, return an error (will cause panic).
    - Return `nil` on success. An error returned here will cause `core.Executor` to panic.
    *Note: This function signature matches `core.CommandHandler`.*
