# Plan for posts/messages.go (Example Feature)

This file defines the specific command and query message structures, query result structures, and the `Validator` interface for the feature.

## Types

*Command and Query structs in the template implement `CommandName()` and `QueryName()` respectively. These methods return the kebab-case name (e.g., `petrock_example_feature_name/create`) which is updated with the actual feature name during placeholder replacement.*

### Interfaces
- `Validator`: Implemented by command structs requiring stateful validation.
    - `Validate(state *PostState) error`: Performs validation using the provided feature state.

### Commands (Implement `core.Command`)
- `CreatePostCommand`: (Implements `CommandName() string`, optionally `Validator`)
    - `Title string`
    - `Content string`
    - `AuthorID string`
    ```go
    // Example implementation of Validator
    func (cmd CreatePostCommand) Validate(state *PostState) error {
        if cmd.Title == "" {
            return errors.New("title cannot be empty")
        }
        // Example stateful validation: Check if a post with this title already exists
        // (Note: This requires iterating or having an index in PostState)
        // if state.PostTitleExists(cmd.Title) {
        //     return errors.New("a post with this title already exists")
        // }
        return nil
    }
    ```
- `UpdatePostCommand`: (Implements `CommandName() string`, optionally `Validator`)
    - `PostID string` // Identifier for the post to update
    - `Title string`
    - `Content string`
    ```go
    // Example implementation of Validator
    func (cmd UpdatePostCommand) Validate(state *PostState) error {
        if cmd.PostID == "" {
            return errors.New("post ID cannot be empty")
        }
        if cmd.Title == "" {
            return errors.New("title cannot be empty")
        }
        // Example stateful validation: Check if the post exists
        if _, exists := state.GetPost(cmd.PostID); !exists {
             return fmt.Errorf("post with ID %s not found", cmd.PostID)
        }
        // Example: Check if *another* post already has the new title
        // if state.PostTitleExistsForOtherID(cmd.Title, cmd.PostID) {
        //     return errors.New("another post with this title already exists")
        // }
        return nil
    }
    ```
- `DeletePostCommand`: (Implements `CommandName() string`, optionally `Validator`)
    - `PostID string` // Identifier for the post to delete
    ```go
    // Example implementation of Validator
    func (cmd DeletePostCommand) Validate(state *PostState) error {
        if cmd.PostID == "" {
            return errors.New("post ID cannot be empty")
        }
        // Example stateful validation: Check if the post exists
        if _, exists := state.GetPost(cmd.PostID); !exists {
             return fmt.Errorf("post with ID %s not found", cmd.PostID)
        }
        return nil
    }
    ```

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

- None, this file primarily defines data structures and interfaces.
