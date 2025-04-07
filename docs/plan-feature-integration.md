# Plan: Feature Integration and Core API

**Goal:** Ensure that after running `petrock new <project> ...` followed by `petrock feature <feature>`, the command `go run ./cmd/<project> serve` starts a functional application where the new feature is fully integrated with core services (logging, state, registries) and accessible via defined API endpoints.

**Target Workflow:**

```sh
# create a project
petrock new blog github.com/dhamidi/blog

# go into the project and add a feature
cd blog
petrock feature posts

# run the application
go run ./cmd/blog serve
```

**Expected Outcome:**

1.  The `posts` feature's commands, queries, and message types are registered with the core registries and message log upon startup.
2.  The application serves HTTP requests on `localhost:8080` (by default).
3.  The following API endpoints are functional:
    *   `GET /commands`: Returns a JSON list of registered command type names.
    *   `POST /commands`: Accepts a JSON payload (`{"type": "CmdName", "payload": {...}}`), executes the command via the registry, logs it, updates state, and returns success/error.
    *   `GET /queries`: Returns a JSON list of registered query names (e.g., `posts/ListQuery`).
    *   `GET /queries/{feature}/{QueryName}`: Accepts query parameters, executes the named query via the registry using the parameters, and returns the JSON result.
4.  Executing commands via `POST /commands` correctly updates the application's in-memory state via the event log replay mechanism.
5.  Visiting the root `/` displays a simple HTML page listing the registered command and query names (e.g., `posts/CreateCommand`).

---

## Implementation Steps

This plan breaks down the work into iterative chunks. Each step builds upon the previous ones.

### Chunk 1: Core Initialization in `serve`

*   **Goal:** Ensure `MessageLog`, `CommandRegistry`, `QueryRegistry`, and a basic `State` representation are correctly initialized and wired together within the `runServe` function of the generated application.
*   **Tasks:**
    1.  **Initialize Registries:** In `internal/skeleton/cmd/petrock_example_project_name/serve.go`, instantiate `core.NewCommandRegistry()` and `core.NewQueryRegistry()` at the start of `runServe`. Store them in local variables.
    2.  **Initialize Encoder:** Instantiate `core.JSONEncoder{}`. Store it in a local variable.
    3.  **Configure DB Path:** Modify `runServe` to accept the database path via a command-line flag (e.g., `--db-path`, defaulting to `app.db`). Update `core.SetupDatabase` call to use this path. Remove the hardcoded path and the corresponding TODO.
    4.  **Initialize MessageLog:** Instantiate `core.NewMessageLog()` using the `*sql.DB` connection and the `JSONEncoder` instance. Handle potential errors during initialization. Store the log in a local variable.
    5.  **Define Placeholder AppState:** Create a simple, temporary `AppState` struct within `serve.go` (it can be moved later). Give it an `Apply(msg interface{}) error` method that currently does nothing but log the received message type. Instantiate this `AppState`.
    6.  **Implement Log Replay:** Add the log replay loop in `runServe` after initializing the `AppState` and `MessageLog`:
        *   Call `messageLog.Load(context.Background())`.
        *   Iterate through the loaded messages.
        *   For each message, call `appState.Apply(msg)`.
        *   Log the number of replayed messages and handle potential errors during load or apply.
    7.  **Pass Dependencies:** Ensure the initialized components (`commandRegistry`, `queryRegistry`, `messageLog`, `appState`, `db`) are available as local variables within `runServe` for subsequent steps.
*   **References:**
    *   `docs/core/log.go.md`
    *   `docs/core/commands.go.md`
    *   `docs/core/queries.go.md`
    *   `internal/skeleton/cmd/petrock_example_project_name/serve.go`
    *   `internal/skeleton/core/log.go`
*   **Definition of Done:**
    *   The `serve` command has a `--db-path` flag.
    *   `runServe` successfully initializes registries, encoder, message log (using the flag), and a placeholder `AppState`.
    *   The log replay loop executes without crashing (it will replay zero messages initially).
    *   All initialized core components are held in local variables within `runServe`.
    *   Relevant TODOs in `serve.go` regarding initialization are removed.

### Chunk 2: Feature Registration Plumbing

*   **Goal:** Connect the `RegisterAllFeatures` function to the `runServe` lifecycle and ensure the `petrock feature` command correctly modifies `features.go` to include new features.
*   **Tasks:**
    1.  **Update `RegisterAllFeatures` Signature:** Modify `internal/skeleton/cmd/petrock_example_project_name/features.go`:
        *   Change the signature of `RegisterAllFeatures` to accept pointers to the initialized `core.CommandRegistry`, `core.QueryRegistry`, `core.MessageLog`, and the placeholder `AppState` (or a more specific state interface later).
        *   Remove the commented-out `messageLog` parameter.
    2.  **Call `RegisterAllFeatures`:** In `runServe` (`internal/skeleton/cmd/petrock_example_project_name/serve.go`), *after* initializing core components and replaying the log, call `RegisterAllFeatures`, passing the initialized instances.
    3.  **Verify `petrock feature` Tooling:** Review the `insertFeatureRegistration` function in `cmd/petrock/feature.go`. Ensure it correctly:
        *   Finds the `// petrock:import-feature` and `// petrock:register-feature` markers.
        *   Inserts the feature package import (e.g., `_ "github.com/dhamidi/blog/posts"`). The underscore import might be sufficient if `RegisterFeature` is called via an init function in the feature package later, but for now, let's assume an explicit call. *Correction:* The current template calls `feature.RegisterFeature`, so a named import is needed.
        *   Inserts the feature registration call (e.g., `posts.RegisterFeature(commands, queries, messageLog, appState)`) using the variable names passed to `RegisterAllFeatures`.
*   **References:**
    *   `internal/skeleton/cmd/petrock_example_project_name/features.go`
    *   `internal/skeleton/cmd/petrock_example_project_name/serve.go`
    *   `cmd/petrock/feature.go`
    *   `docs/feature/register.go.md`
*   **Definition of Done:**
    *   After running `petrock new blog ...` and `petrock feature posts`, the generated `cmd/blog/features.go` contains:
        *   `import "github.com/dhamidi/blog/posts"` (or similar based on module path).
        *   `posts.RegisterFeature(commands, queries, messageLog, appState)` inside `RegisterAllFeatures`.
    *   `go run ./cmd/blog serve` executes without errors related to calling `RegisterAllFeatures` or the feature's `RegisterFeature` function.

### Chunk 3: Basic State Management in Feature Template

*   **Goal:** Implement a minimal functional `State` and `Apply` method in the feature template (`internal/skeleton/feature_template`) so that features can manage their own state based on logged commands.
*   **Tasks:**
    1.  **Refine `feature_template/state.go`:**
        *   Ensure `NewState()` initializes `Items map[string]*Item` correctly.
        *   Implement the `Apply(msg interface{}) error` method:
            *   Use a type switch (`switch cmd := msg.(type)`) to handle `CreateCommand`, `UpdateCommand`, `DeleteCommand`.
            *   For each command, perform the corresponding action on the `s.Items` map (add, update, delete). Remember to use the mutex (`s.mu.Lock()/Unlock()`).
            *   Use `slog` for logging state changes within `Apply`.
            *   Handle potential errors gracefully (e.g., updating/deleting a non-existent item - return an error or log a warning based on desired semantics).
        *   Implement `GetItem(id string) (*Item, bool)` and `ListItems(page, pageSize int, filter string) ([]*Item, int)` to read from the `Items` map safely using the read mutex (`s.mu.RLock()/RUnlock()`).
    2.  **Refine `feature_template/register.go`:**
        *   Ensure the call `RegisterTypes(messageLog)` correctly registers `CreateCommand`, `UpdateCommand`, `DeleteCommand` with the `core.MessageLog` instance passed into `RegisterFeature`. This function `RegisterTypes` should exist in `state.go`.
        *   Ensure `RegisterFeature` accepts the `*State` instance and passes it to `NewExecutor` and `NewQuerier`.
*   **References:**
    *   `internal/skeleton/feature_template/state.go`
    *   `internal/skeleton/feature_template/register.go`
    *   `internal/skeleton/feature_template/execute.go` (ensure `NewExecutor` takes state)
    *   `internal/skeleton/feature_template/query.go` (ensure `NewQuerier` takes state)
    *   `docs/feature/state.go.md`
    *   `docs/core/log.go.md`
*   **Definition of Done:**
    *   The generated feature package (e.g., `posts`) has a `State` struct capable of modifying its internal `Items` map based on its own command types passed to its `Apply` method.
    *   The feature correctly registers its command types via `RegisterTypes` during the `RegisterFeature` call.
    *   The feature's `Executor` and `Querier` are initialized with the feature's `State`.

### Chunk 4: API Endpoint - List Commands (`GET /commands`)

*   **Goal:** Create an HTTP handler that returns a list of registered command type names.
*   **Tasks:**
    1.  **Expose Registered Commands:** In `internal/skeleton/core/commands.go`:
        *   Add a method to `CommandRegistry`, e.g., `RegisteredCommandNames() []string`.
        *   This method should acquire the read lock (`r.mu.RLock()`), iterate over the `r.handlers` map, extract the type name (`reflect.Type.Name()`) for each key, collect them into a slice of strings, release the lock, and return the slice.
    2.  **Create Handler:** In `internal/skeleton/cmd/petrock_example_project_name/serve.go`:
        *   Create a handler function `handleListCommands(registry *core.CommandRegistry) http.HandlerFunc`.
    3.  **Implement Handler Logic:**
        *   Inside the handler, call `registry.RegisteredCommandNames()` to get the list.
        *   Set the `Content-Type` header to `application/json`.
        *   Use `encoding/json` to marshal the slice of names into the response body (`http.ResponseWriter`). Handle potential marshaling errors (return HTTP 500).
    4.  **Register Route:** In `runServe`, register the handler with the `ServeMux`: `mux.HandleFunc("GET /commands", handleListCommands(commandRegistry))`.
*   **References:**
    *   `internal/skeleton/core/commands.go`
    *   `internal/skeleton/cmd/petrock_example_project_name/serve.go`
*   **Definition of Done:**
    *   Running `go run ./cmd/blog serve` starts the server.
    *   Sending a `GET /commands` request returns a `200 OK` response with `Content-Type: application/json` and a body like `["CreateCommand", "UpdateCommand", "DeleteCommand"]` (assuming the `posts` feature was added).

### Chunk 5: API Endpoint - Execute Command (`POST /commands`)

*   **Goal:** Create an HTTP handler that accepts a command payload, decodes it, dispatches it through the registry and state mechanism, and returns a result.
*   **Tasks:**
    1.  **Define Command Request Structure:** Decide on a standard JSON structure for requests, e.g.:
        ```json
        {
          "type": "CreateCommand", // The registered name of the command type
          "payload": {
            "name": "My Item",
            "description": "Details here"
            // ... other fields specific to CreateCommand
          }
        }
        ```
    2.  **Enhance CommandRegistry:** In `internal/skeleton/core/commands.go`:
        *   Add a method `GetCommandType(name string) (reflect.Type, bool)` that looks up the registered `reflect.Type` by its name string. Requires iterating through the `handlers` map keys. Thread-safe (read lock).
    3.  **Create Handler:** In `internal/skeleton/cmd/petrock_example_project_name/serve.go`:
        *   Create `handleExecuteCommand(registry *core.CommandRegistry) http.HandlerFunc`.
    4.  **Implement Handler Logic:**
        *   Define an intermediate struct to decode the request body, e.g., `type commandRequest struct { Type string `json:"type"`; Payload json.RawMessage `json:"payload"` }`.
        *   Decode the request body into this intermediate struct. Handle JSON decoding errors (return HTTP 400).
        *   Use `registry.GetCommandType(req.Type)` to find the `reflect.Type` for the command. If not found, return HTTP 404 or 400.
        *   Create a new zero-value instance of the command struct using `reflect.New(cmdType).Interface()`. This returns a pointer.
        *   Unmarshal the `req.Payload` (which is `json.RawMessage`) into the command instance pointer. Handle JSON unmarshaling errors (return HTTP 400).
        *   Dispatch the command: `err := registry.Dispatch(r.Context(), reflect.ValueOf(cmdInstance).Elem().Interface())`. We pass the actual struct value, not the pointer, assuming handlers expect the value type. *Correction:* Handlers likely expect the value type, but `Dispatch` takes `interface{}`. Let's stick to passing the value: `cmdValue := reflect.ValueOf(cmdInstance).Elem().Interface()`. Dispatch `err := registry.Dispatch(r.Context(), cmdValue)`.
        *   Handle errors returned by `Dispatch`:
            *   If it's a known validation error type (needs defining, maybe in `core`), return HTTP 400 with the error message.
            *   For other errors (e.g., persistence errors from the handler, state update errors), return HTTP 500. Log the error server-side.
        *   If `Dispatch` succeeds, return HTTP 200 OK or HTTP 202 Accepted. Optionally return a simple JSON success message like `{"status": "success"}`.
    5.  **Register Route:** In `runServe`, register `mux.HandleFunc("POST /commands", handleExecuteCommand(commandRegistry))`.
*   **References:**
    *   `internal/skeleton/core/commands.go`
    *   `internal/skeleton/cmd/petrock_example_project_name/serve.go`
    *   `internal/skeleton/feature_template/execute.go` (Command Handlers)
    *   `docs/core/commands.go.md`
*   **Definition of Done:**
    *   Running `go run ./cmd/blog serve` starts the server.
    *   Sending a `POST /commands` with a valid JSON payload (e.g., for `CreateCommand`) results in:
        *   The corresponding command handler in the `posts` feature being executed.
        *   A message being appended to the `app.db` log file.
        *   The feature's in-memory state being updated via its `Apply` method.
        *   A `200 OK` (or `202 Accepted`) response.
    *   Sending requests with invalid JSON, unknown command types, or payloads causing validation errors returns appropriate HTTP error codes (400, 404, 500).

### Chunk 6: API Endpoint - List Queries (`GET /queries`)

*   **Goal:** Create an HTTP handler that returns a list of registered query type names.
*   **Tasks:**
    1.  **Expose Registered Queries:** In `internal/skeleton/core/queries.go`:
        *   Add `RegisteredQueryNames() []string` method to `QueryRegistry`, similar to the command registry. Thread-safe (read lock).
    2.  **Create Handler:** In `internal/skeleton/cmd/petrock_example_project_name/serve.go`:
        *   Create `handleListQueries(registry *core.QueryRegistry) http.HandlerFunc`.
    3.  **Implement Handler Logic:**
        *   Call `registry.RegisteredQueryNames()`.
        *   Set `Content-Type: application/json`.
        *   Marshal the list of names to the response. Handle errors (HTTP 500).
    4.  **Register Route:** In `runServe`, register `mux.HandleFunc("GET /queries", handleListQueries(queryRegistry))`.
*   **References:**
    *   `internal/skeleton/core/queries.go`
    *   `internal/skeleton/cmd/petrock_example_project_name/serve.go`
*   **Definition of Done:**
    *   Running `go run ./cmd/blog serve` starts the server.
    *   Sending a `GET /queries` request returns a `200 OK` response with `Content-Type: application/json` and a body like `["GetQuery", "ListQuery"]` (assuming the `posts` feature was added).

### Chunk 7: API Endpoint - Execute Query (`GET /queries/{name}`)

*   **Goal:** Create an HTTP handler that executes a named query, populating its fields from URL query parameters, and returns the result.
*   **Tasks:**
    1.  **Query Naming Convention:** Use the full `feature/QueryName` string. The URL path will be `/queries/{feature}/{QueryName}`.
    2.  **Enhance QueryRegistry:** In `internal/skeleton/core/queries.go`:
        *   Ensure `GetQueryType(name string) (reflect.Type, bool)` accepts the full name (e.g., "posts/ListQuery") and looks up the type correctly from the registry's internal storage.
    3.  **Create Handler:** In `internal/skeleton/cmd/petrock_example_project_name/serve.go`:
        *   Create `handleExecuteQuery(registry *core.QueryRegistry) http.HandlerFunc`.
    4.  **Implement Handler Logic:**
        *   Extract the `feature` and `queryName` parts from the URL path using `r.PathValue()`. Construct the full name `feature + "/" + queryName`. Handle errors if parts are missing.
        *   Use `registry.GetQueryType(fullQueryName)` to find the `reflect.Type`. If not found, return HTTP 404.
        *   Create a new zero-value instance (pointer) of the query struct using `reflect.New(queryType).Interface()`.
        *   **Populate Query Struct from URL Params:** (Logic remains the same)
            *   Get URL query parameters using `r.URL.Query()`.
            *   Iterate through the fields of the query struct instance (using reflection on the pointer: `reflect.ValueOf(queryInstance).Elem()`).
            *   For each field, check if a corresponding key exists in the URL parameters.
            *   If found, get the string value(s) from the `url.Values`.
            *   Convert the string value to the field's type (e.g., `string`, `int`, `bool`). Use `strconv` for conversions. Handle potential conversion errors (return HTTP 400).
            *   Set the field's value using reflection (`field.SetString()`, `field.SetInt()`, etc.). Handle potential errors during setting (return HTTP 400).
        *   Get the query value and ensure it implements `core.Query`: `queryValue, ok := queryInstance.Interface().(core.Query)`. Handle `!ok`.
        *   Dispatch the populated query: `result, err := registry.Dispatch(r.Context(), queryValue)`.
        *   Handle errors from `Dispatch`:
            *   If it's a "not found" error from the handler (needs defining, maybe a standard `ErrNotFound` in `core`), return HTTP 404.
            *   For other errors, return HTTP 500. Log the error.
        *   If successful:
            *   Set `Content-Type: application/json`.
            *   Marshal the `result` (which implements `core.QueryResult`) to the response body. Handle marshaling errors (HTTP 500).
            *   Return HTTP 200 OK.
    5.  **Register Route:** In `runServe`, register `mux.HandleFunc("GET /queries/{feature}/{queryName}", handleExecuteQuery(queryRegistry))`.
*   **References:**
    *   `internal/skeleton/core/queries.go`
    *   `internal/skeleton/cmd/petrock_example_project_name/serve.go`
    *   `internal/skeleton/feature_template/query.go` (Query Handlers)
    *   `docs/core/queries.go.md`
    *   Go `reflect` package documentation
    *   Go `strconv` package documentation
*   **Definition of Done:**
    *   Running `go run ./cmd/blog serve` starts the server.
    *   Sending `GET /queries/ListQuery?page=1&pageSize=5` executes the `ListQuery` handler in the `posts` feature, populating the `ListQuery` struct from the URL parameters. It returns a `200 OK` with a JSON body representing the `ListResult` based on the current state.
    *   Sending `GET /queries/GetQuery?ID=some-id` works similarly for retrieving single items.
    *   Requests with invalid query names, missing/invalid parameters, or queries resulting in "not found" return appropriate HTTP error codes (404, 400, 500).

### Chunk 8: Discoverability - Update Index Page

*   **Goal:** Modify the application's root page (`/`) to display lists of registered command and query names.
*   **Tasks:**
    1.  **Update `HandleIndex` Signature:** In `internal/skeleton/core/page_index.go`, modify `HandleIndex` to accept `*core.CommandRegistry` and `*core.QueryRegistry` as parameters.
    2.  **Fetch Lists in Handler:** Inside `HandleIndex`, call `commandRegistry.RegisteredCommandNames()` and `queryRegistry.RegisteredQueryNames()` to get the lists.
    3.  **Update `IndexPage` Signature:** Modify the `IndexPage` component function to accept `[]string` for command names and `[]string` for query names.
    4.  **Pass Lists to Component:** In `HandleIndex`, pass the fetched lists when calling `IndexPage`.
    5.  **Render Lists in Component:** Update the Gomponents structure within `IndexPage` to:
        *   Render a heading like "Available Commands".
        *   Render an unordered list (`html.Ul`) where each list item (`html.Li`) displays a command name.
        *   Render a heading like "Available Queries".
        *   Render an unordered list for query names.
        *   Remove the previous static welcome text.
    6.  **Update `serve.go` Call:** In `runServe`, update the call to `core.HandleIndex` in the route registration to pass the initialized registries: `mux.HandleFunc("GET /", core.HandleIndex(commandRegistry, queryRegistry))`.
*   **References:**
    *   `internal/skeleton/core/page_index.go`
    *   `internal/skeleton/cmd/petrock_example_project_name/serve.go`
*   **Definition of Done:**
    *   Visiting `/` on the running application displays an HTML page containing two lists: one for registered command names and one for registered query names.

### Chunk 9: Template Cleanup and Refinement

*   **Goal:** Remove addressed TODO comments and refine the skeleton code based on the implementations from previous chunks.
*   **Tasks:**
    1.  **Review `serve.go`:** Remove TODOs related to registry/log/state initialization, feature registration, basic API route setup, DB path configuration, and encoder instantiation.
    2.  **Review `commands.go`/`queries.go`:** Remove TODOs for instrumentation if not added. Ensure registry methods for listing/lookup are implemented as described.
    3.  **Review `log.go`:** Remove TODO for DB path config and type registration.
    4.  **Review `feature_template/...`:**
        *   Remove the TODO regarding `ErrNotFound` in `query.go` (or implement a shared `core.ErrNotFound`).
        *   Review `view.go`: The API endpoints are now the primary interaction method. Decide whether to keep the HTMX examples in the template (perhaps commented out or updated to use the API) or remove them to avoid confusion. For now, let's keep them but ensure placeholders like `/feature-path` are replaced with something indicative like `/api/commands` or `/api/queries/...` if they were to be adapted. *Decision:* Remove the HTMX attributes from the default view templates (`ItemView`, `ItemForm`, `NewItemButton`) as the primary interaction is now the API. Keep the basic structure rendering data.
    5.  **Review `build.go`/`deploy.go`:** Leave TODOs for version flags and remote commands, as they are outside the scope of this integration plan.
*   **References:** All `internal/skeleton/...` files.
*   **Definition of Done:** Skeleton code is cleaner, reflects the implemented API patterns, and has fewer TODO comments related to the core feature integration logic. Feature view templates render data but do not contain HTMX attributes pointing to non-existent feature-specific handlers.

### Chunk 10: Documentation Update

*   **Goal:** Update project documentation (`docs/`) to reflect the new API endpoints and the implemented feature integration process.
*   **Tasks:**
    1.  **Update `docs/high-level.md`:**
        *   Clearly describe the `/commands` and `/queries` API endpoints, including request/response formats.
        *   Explain that features are integrated by calling `RegisterFeature` from `features.go`, which is automatically updated by `petrock feature`.
    2.  **Update Feature Docs:**
        *   In `docs/feature/register.go.md`, confirm that the registration call is added automatically to `features.go`.
        *   In `docs/feature/view.go.md`, clarify that views are primarily for rendering state, and interaction logic now primarily uses the core API endpoints.
    3.  **Update Core Docs:**
        *   Update `docs/core/commands.go.md` and `docs/core/queries.go.md` to document the registry methods for listing names and looking up types (`RegisteredCommandNames`, `GetCommandType`, etc.).
*   **References:** `docs/high-level.md`, `docs/feature/*.md`, `docs/core/*.md`
*   **Definition of Done:** Project documentation accurately reflects the implemented feature integration mechanism, the core API endpoints, and the updated role of feature-specific views.

---
This plan provides a structured approach to achieving the desired feature integration and API functionality. Each chunk represents a testable increment.
