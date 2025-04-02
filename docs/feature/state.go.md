# Plan for posts/state.go (Example Feature)

This file defines the in-memory application state specific to the posts feature. This state is built by replaying logged commands/events at startup and updated by command handlers during runtime.

## Types

- `Post`: The internal struct representing a single post within the application state.
    - `ID string`
    - `Title string`
    - `Content string`
    - `AuthorID string`
    - `CreatedAt time.Time`
    - `UpdatedAt time.Time`
    - `IsDeleted bool` // Or handle deletion by removing from the map
- `PostState`: The main struct holding the collective state for the posts feature.
    - `Posts map[string]*Post`: A map from Post ID to the Post object pointer. Using pointers allows in-place updates.
    - `mu sync.RWMutex`: A mutex to protect concurrent access to the `Posts` map.

## Functions

- `NewPostState() *PostState`: Constructor to create an initialized (empty) `PostState`.
- `(s *PostState) Apply(msg core.Message)`: Updates the state based on a logged message (replayed event/command). This function contains the core logic for state reconstruction.
    - It needs to decode `msg.Data` based on `msg.Type` (which should correspond to command types like `CreatePostCommand`, `UpdatePostCommand`, `DeletePostCommand`).
    - Based on the decoded command, it modifies the `s.Posts` map accordingly (adds, updates, or marks/removes posts). Requires locking/unlocking `s.mu`.
- `(s *PostState) GetPost(id string) (*Post, bool)`: Retrieves a post by its ID. Returns the post pointer and `true` if found, `nil` and `false` otherwise. Requires read-locking `s.mu`.
- `(s *PostState) ListPosts(page, pageSize int, authorIDFilter string) ([]*Post, int)`: Retrieves a slice of posts, applying filtering and pagination. Returns the slice of posts for the current page and the total count of matching posts. Requires read-locking `s.mu`.
- `(s *PostState) AddPost(post *Post)`: Adds a new post to the state map. Requires write-locking `s.mu`.
- `(s *PostState) UpdatePost(post *Post)`: Updates an existing post in the state map. Requires write-locking `s.mu`.
- `(s *PostState) DeletePost(id string)`: Removes (or marks as deleted) a post from the state map. Requires write-locking `s.mu`.
