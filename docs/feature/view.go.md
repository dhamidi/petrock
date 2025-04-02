# Plan for posts/view.go (Example Feature)

This file defines Gomponents specific to rendering HTML for the posts feature. These components use the feature's state/query results and potentially shared core components.

## Types

- None specific to this file. Component logic is encapsulated in functions.

## Functions

*These functions return `gomponents.Node`.*

- `PostView(post PostQueryResult) gomponents.Node`: Renders the HTML representation of a single post, perhaps showing title, content, author, and timestamps.
- `PostForm(form *core.Form, post *PostQueryResult) gomponents.Node`: Renders an HTML `<form>` for creating or editing a post.
    - Uses `core.Input`, `core.TextArea`, `core.Button`, `core.FormError` components.
    - Populates fields with data from `post` if provided (for editing).
    - Uses `form.Values` to repopulate fields on validation error.
    - Includes CSRF token input (`core.CSRFTokenInput`).
    - Sets appropriate `action` and `method` attributes. Can include HTMX attributes (`hx-post`, `hx-put`, `hx-target`).
- `PostsListView(result PostsListQueryResult) gomponents.Node`: Renders a list of posts, often as a table or a series of divs.
    - Iterates over `result.Posts`.
    - Calls `PostView` for each post (or renders a summary row).
    - May include pagination controls based on `result.TotalCount`, `result.Page`, `result.PageSize`.
- `NewPostButton() gomponents.Node`: Renders a button or link that navigates or triggers loading the `PostForm` for creating a new post (e.g., using HTMX `hx-get`).
