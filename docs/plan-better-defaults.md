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

### 1. Routes and Handlers Updates

1. Update `routes.go` to:
   - Remove existing routes that use PUT/DELETE methods
   - Add all the GET/POST routes as specified
   - Ensure route naming follows the convention (prefix + path)

2. Update `http.go` to:
   - Implement new handlers that work with HTML forms instead of JSON
   - Add handlers for the new routes (new post form, edit form, delete confirmation)
   - Modify existing handlers to return HTML responses
   - Add form handling logic (validation, error display, redirection)

### 2. View Component Updates

1. Modify `view.go` to:
   - Implement form components that use `core.Form` for validation
   - Create view components for each page (list, detail, new, edit, delete)
   - Ensure components follow Tailwind CSS styling
   - Add navigation elements between views

### 3. Command and Query Modifications

1. Update `messages.go` to:
   - Ensure all commands have proper validation
   - Add string trimming to command validation
   - Add uniqueness checks for IDs

2. Update `query.go` to:
   - Ensure query handlers return appropriate data for HTML views
   - Add any needed query methods for the new routes

### 4. Testing and Verification

1. Create test cases for each route
2. Verify form validation works correctly
3. Test navigation between pages
4. Ensure all error cases are handled appropriately

## Detailed Tasks

### Routes and Handlers

1. **Add Route Definitions**
   - [ ] Update `/routes.go` to define all required GET/POST routes
   - [ ] Remove PUT/DELETE routes
   - [ ] Update route prefix to use the feature name

2. **Create New Form Handlers**
   - [ ] Implement `HandleNewForm` (GET /posts/new)
   - [ ] Implement `HandleCreateForm` (POST /posts/new)
   - [ ] Implement `HandleEditForm` (GET /posts/{id}/edit)
   - [ ] Implement `HandleUpdateForm` (POST /posts/{id}/edit)
   - [ ] Implement `HandleDeleteForm` (GET /posts/{id}/delete)
   - [ ] Implement `HandleDeleteConfirm` (POST /posts/{id}/delete)

3. **Update Existing Handlers**
   - [ ] Modify `HandleGetItem` to render HTML instead of JSON
   - [ ] Modify `HandleListItems` to render HTML instead of JSON
   - [ ] Remove or adapt JSON-specific handlers

### View Components

1. **Create Form Components**
   - [ ] Implement `NewItemForm` using core.Form
   - [ ] Implement `EditItemForm` using core.Form
   - [ ] Implement `DeleteConfirmForm` 

2. **Update View Components**
   - [ ] Enhance `ItemView` to display all attributes
   - [ ] Update `ItemsListView` to include new item button
   - [ ] Add navigation links between views

3. **Styling and UX**
   - [ ] Apply consistent Tailwind CSS styling
   - [ ] Add basic form validation feedback
   - [ ] Ensure responsive design

### Command and Query Updates

1. **Command Validation**
   - [ ] Update `CreateCommand.Validate()` to trim strings
   - [ ] Add empty string validation
   - [ ] Add ID uniqueness check

2. **Form Handling**
   - [ ] Implement logic to convert form data to commands
   - [ ] Add error handling for validation failures
   - [ ] Implement proper redirects after successful operations

### Common Utilities

1. **Core Form Integration**
   - [ ] Ensure proper use of the core.Form package
   - [ ] Add helper functions for form validation
   - [ ] Add utilities for form rendering

2. **HTTP Utilities**
   - [ ] Add redirect helpers
   - [ ] Add form parsing helpers
   - [ ] Add CSRF token management
