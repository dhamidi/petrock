# Plan for posts/view.go (Example Feature)

This file defines Gomponents specific to rendering HTML for the posts feature. These components use the feature's state/query results and potentially shared core components.

## Types

- None specific to this file. Component logic is encapsulated in functions.

## Functions

*These functions return `gomponents.Node`.*

- `PostView(post PostQueryResult) gomponents.Node`: Renders the HTML representation of a single post, perhaps showing title, content, author, and timestamps.
- `PostForm(form *core.Form, post *PostQueryResult, csrfToken string) gomponents.Node`: Renders an HTML `<form>` for creating or editing a post.
    - Uses `core.Input`, `core.TextArea`, `core.Button`, `core.FormError` components.
    - Populates fields with data from `post` if provided (for editing).
    - Uses `form.Values` to repopulate fields on validation error.
    - Includes CSRF token input (`core.CSRFTokenInput`).
    - Sets appropriate `action` and `method` attributes. *Note: Form submission can be handled in multiple ways: via JavaScript making calls to the core API (`POST /commands`), via JavaScript calling feature-specific routes (e.g., `POST /posts`), or potentially using libraries like HTMX targeting either core or feature routes.*
- `PostsListView(result PostsListQueryResult) gomponents.Node`: Renders a list of posts, often as a table or a series of divs.
    - Iterates over `result.Posts`.
    - Calls `PostView` for each post (or renders a summary row).
    - May include pagination controls based on `result.TotalCount`, `result.Page`, `result.PageSize`.
- `NewItemButton() gomponents.Node`: Renders a button or link that might navigate to a page containing the `ItemForm` or trigger client-side logic (JavaScript, HTMX) to display a form or initiate an action.

*Note: Feature-specific views primarily focus on rendering data retrieved via queries (either from the core API or feature-specific endpoints). User interactions (creating, updating, deleting) can be handled by client-side logic targeting either the core API endpoints (`/commands`, `/queries`) or custom endpoints defined by the feature in `routes.go`.*
