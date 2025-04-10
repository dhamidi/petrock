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

1. **Update `internal/skeleton/feature_template/routes.go`**
   - [ ] Remove routes with PUT method (`mux.HandleFunc("PUT "+featurePrefix+"/{id}", deps.HandleUpdateItem)`) 
   - [ ] Remove routes with DELETE method (`mux.HandleFunc("DELETE "+featurePrefix+"/{id}", deps.HandleDeleteItem)`) 
   - [ ] Add route for new form: `mux.HandleFunc("GET "+featurePrefix+"/new", deps.HandleNewForm)`
   - [ ] Add route for create form submission: `mux.HandleFunc("POST "+featurePrefix+"/new", deps.HandleCreateForm)`
   - [ ] Add route for edit form: `mux.HandleFunc("GET "+featurePrefix+"/{id}/edit", deps.HandleEditForm)`
   - [ ] Add route for update form submission: `mux.HandleFunc("POST "+featurePrefix+"/{id}/edit", deps.HandleUpdateForm)`
   - [ ] Add route for delete confirmation: `mux.HandleFunc("GET "+featurePrefix+"/{id}/delete", deps.HandleDeleteForm)`
   - [ ] Add route for delete confirmation submission: `mux.HandleFunc("POST "+featurePrefix+"/{id}/delete", deps.HandleDeleteConfirm)`

2. **Create New Handler Methods in `internal/skeleton/feature_template/http.go`**
   - [ ] Add `HandleNewForm(w http.ResponseWriter, r *http.Request)` method to render form (reference ItemForm component)
   - [ ] Add `HandleCreateForm(w http.ResponseWriter, r *http.Request)` method that parses form using `r.ParseForm()` and creates CreateCommand
   - [ ] Add `HandleEditForm(w http.ResponseWriter, r *http.Request)` method that loads item and renders edit form
   - [ ] Add `HandleUpdateForm(w http.ResponseWriter, r *http.Request)` method that parses form and creates UpdateCommand
   - [ ] Add `HandleDeleteForm(w http.ResponseWriter, r *http.Request)` method that loads item and renders delete confirmation
   - [ ] Add `HandleDeleteConfirm(w http.ResponseWriter, r *http.Request)` method that creates DeleteCommand

3. **Update Existing Handler Methods in `internal/skeleton/feature_template/http.go`**
   - [ ] Modify `HandleGetItem` to use `core.RenderView` with the `ItemView` component instead of `respondJSON`
   - [ ] Modify `HandleListItems` to use `core.RenderView` with the `ItemsListView` component instead of `respondJSON`
   - [ ] Update content-type headers to `text/html` instead of `application/json`

### 2. View Components Implementation in `internal/skeleton/feature_template/view.go`

1. **Form Components**
   - [ ] Update `ItemForm` to work with both create and edit by handling nil items in `func ItemForm(form *core.Form, item *Result, csrfToken string)`
   - [ ] Create new function `func DeleteConfirmForm(item *Result, csrfToken string) g.Node` for delete confirmation
   - [ ] Add explicit validation messages using `core.Form.GetError()` for each form field
   - [ ] Set form action URLs to new routes (`/posts/new` and `/posts/{id}/edit`) using html.Action attribute
   - [ ] Set form method to POST using html.Method attribute

2. **Item View Components**
   - [ ] Update `ItemView` to display each field with labels using html.Div and html.Label elements
   - [ ] Add edit link with href=`/posts/{id}/edit` using g.Attr("href", ...)
   - [ ] Add delete link with href=`/posts/{id}/delete` using g.Attr("href", ...)
   - [ ] Update `ItemsListView` to include a "New Post" button linking to `/posts/new`
   - [ ] Add "Back to list" links on item pages using g.Attr("href", "/posts")

### 3. Command/Query Functionality

1. **Update Command Structure in `internal/skeleton/feature_template/messages.go`**
   - [ ] Add `CreatedAt time.Time` field to `CreateCommand` struct
   - [ ] Add `UpdatedAt time.Time` field to `UpdateCommand` struct
   - [ ] Add `DeletedAt time.Time` field to `DeleteCommand` struct
   - [ ] Set these timestamp fields in handlers using `time.Now().UTC()` before passing to `executor.Execute`

2. **Update Command Validation Methods in `internal/skeleton/feature_template/messages.go`**
   - [ ] Modify `CreateCommand.Validate()` to trim all string fields with `strings.TrimSpace(c.Name)` etc.
   - [ ] Add validation in `CreateCommand.Validate()` to check if fields are empty after trimming
   - [ ] Add uniqueness check in `CreateCommand.Validate()` using the state map: `_, exists := state.Items[c.Name]`
   - [ ] Update `UpdateCommand.Validate()` to trim strings and check for empty values after trimming
   - [ ] Update `DeleteCommand.Validate()` to check if item exists: `_, found := state.GetItem(c.ID)`

3. **Update HTTP Handler Logic for Form Handling in `internal/skeleton/feature_template/http.go`**
   - [ ] Add helper function `parseItemForm(r *http.Request) (*core.Form, error)` to parse form data
   - [ ] Add code to convert form values to commands: `cmd := CreateCommand{Name: form.Get("name"), ...}`
   - [ ] Add error handling pattern to render forms with errors: `if !form.IsValid() { /* render form with errors */ }`
   - [ ] Use `w.Header().Set("Location", "/posts")` and `w.WriteHeader(http.StatusSeeOther)` for redirects

### 4. Form Integration Helpers

1. **HTML Form Utilities in `internal/skeleton/feature_template/view.go`**
   - [ ] Create helper function `func formErrorMessage(form *core.Form, field string) g.Node` that returns error message HTML
   - [ ] Create helper function `func csrfField(token string) g.Node` for CSRF token input field
   - [ ] Add helper function `func preserveFormValues(form *core.Form, fields []string) map[string]string` for validation failures
