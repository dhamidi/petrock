# Overview

Read through docs/high-level.md and internal/skeleton/feature_template.

We'll be mostly working with these files.

The goal for this project is to make the default feature template more enticing.

Ultimately the following should be possible after running `petrock feature posts` (`posts` corresponds to `petrock_example_feature_name` in internal/skeleton/feature_template)

These routes will be defined:

- `GET /posts/new` renders a form using (core/form.go) for adding a new post
- `POST /posts/new` accepts the form, converts it into a `CreateCommand`, and dispatches it to the core.Executor
  - during validation, this command trims all incoming strings and makes sure they are not empty
  - it also checks that the provided post string ID is unique
  - in case of errors the corresponding http handler will render the form with errors
  - in case of success, it redirects to `GET /posts` with HTTP Status SeeOther
- `GET /posts` lists all posts by issuing a posts.ListQuery
  - it includes a button which links to the new posts form
- `GET /posts/{id}` renders a given post as a series of labels + pre-formatted text fields for each attribute on a post
- `GET /posts/{id}/edit` renders a form for editing the `content` field of a post
- `POST /posts/{id}/edit` handles the form for editing and dispatches a corresponding `UpdateCommand`
- `GET /posts/{id}/delete` renders a form with a button that will trigger a `DeleteCommand` for the given post
- `POST /posts/{id}/delete` will accept the form from the previous route and actually dispatch the `DeleteCommand`

## Implementation Plan

### 1. Route & HTTP Handler Implementation

1. **Update `routes.go`**
   - [ ] Remove routes using PUT/DELETE methods
   - [ ] Add GET/POST routes as specified in the requirements
   - [ ] Define the correct feature prefix path (`/petrock_example_feature_name` → `/posts`)

2. **Create New Handlers in `http.go`**
   - [ ] `HandleNewForm` (GET /posts/new) - renders the form for adding new post
   - [ ] `HandleCreateForm` (POST /posts/new) - processes form submission → CreateCommand
   - [ ] `HandleEditForm` (GET /posts/{id}/edit) - renders edit form for a post
   - [ ] `HandleUpdateForm` (POST /posts/{id}/edit) - processes edit form → UpdateCommand
   - [ ] `HandleDeleteForm` (GET /posts/{id}/delete) - renders delete confirmation
   - [ ] `HandleDeleteConfirm` (POST /posts/{id}/delete) - processes delete → DeleteCommand

3. **Update Existing Handlers in `http.go`**
   - [ ] Modify `HandleGetItem` to render HTML using `ItemView` instead of JSON
   - [ ] Modify `HandleListItems` to render HTML using `ItemsListView` instead of JSON
   - [ ] Replace JSON responses with HTML responses using gomponents

### 2. View Components Implementation in `view.go`

1. **Form Components**
   - [ ] Update `ItemForm` to work as both create and edit form
   - [ ] Implement `DeleteConfirmForm` for deletion confirmation
   - [ ] Ensure forms use `core.Form` methods like `HasError`, `GetError`, `ValidateRequired`
   - [ ] Set proper form actions and methods (POST)

2. **Item View Components**
   - [ ] Update `ItemView` to render all post attributes with labels
   - [ ] Add proper links to edit and delete actions
   - [ ] Update `ItemsListView` to include "New Post" button linking to /posts/new
   - [ ] Add navigation between views (back to list, edit links, etc.)

### 3. Command/Query Functionality

1. **Update Command Validation in `messages.go`**
   - [ ] Ensure `CreateCommand.Validate()` trims strings using `strings.TrimSpace()`
   - [ ] Add empty string validation in `CreateCommand.Validate()`
   - [ ] Add ID uniqueness check in state map
   - [ ] Update `UpdateCommand.Validate()` for proper field validation
   - [ ] Update `DeleteCommand.Validate()` to verify post exists

2. **Update HTTP Handler Logic for Form Handling**
   - [ ] Add `parseForm` helper to extract form data
   - [ ] Implement form → command conversion in handlers
   - [ ] Add error handling to render forms with validation errors
   - [ ] Implement proper HTTP 303 See Other redirects on success

### 4. Form Integration Helpers

1. **HTML Form Utilities**
   - [ ] Create helper function to render form errors
   - [ ] Create helper to generate CSRF tokens for forms
   - [ ] Add utility to preserve form values on validation failure
