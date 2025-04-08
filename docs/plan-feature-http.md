# Plan: Implement Feature-Specific HTTP Routes

**Goal:** Extend the Petrock framework and generated application skeleton to allow features to define and register their own HTTP routes and handlers, coexisting with the core `/commands` and `/queries` API endpoints.

**Context:** Currently, all interactions primarily go through the centralized `/commands` and `/queries` API. This plan introduces the ability for features like `posts` to expose more conventional RESTful or custom endpoints (e.g., `GET /posts/{id}`, `POST /posts`) directly.

## Implementation Steps

### Chunk 1: Create Feature HTTP Template Files - DONE

- **Goal:** Add the necessary template files (`routes.go`, `http.go`) to the feature skeleton.
- **Tasks:**
  1. Create `internal/skeleton/feature_template/routes.go`:
      - Define the `RegisterRoutes(mux *http.ServeMux, deps FeatureServer)` function signature.
      - Include example `mux.HandleFunc` calls mapping prefixed paths (e.g., `/petrock_example_feature_name/{id}`) to methods on the `deps` object.
      - Add comments explaining prefixing conventions and overriding potential.
  2. Create `internal/skeleton/feature_template/http.go`:
      - Define the `FeatureServer` struct holding dependencies (`*Executor`, `*Querier`, `*State`, `*core.MessageLog`, `*core.CommandRegistry`, etc.).
      - Implement the `NewFeatureServer(...)` constructor.
      - Add example handler methods (e.g., `HandleGetItem`, `HandleCreateItem`) attached to `FeatureServer`.
      - Implement basic request parsing, dependency usage (calling querier/executor), error handling, and response writing (e.g., JSON).
      - Implement the recommended state change strategy in write handlers (using `MessageLog.Append` or `CommandRegistry.Dispatch`).
- **References:**
  - `docs/feature/routes.go.md`
  - `docs/feature/http.go.md`
- **Definition of Done:** The two new template files exist in `internal/skeleton/feature_template/` with the basic structure and example implementations.

### Chunk 2: Update `petrock feature` Command - DONE

- **Goal:** Modify the `petrock feature` command to copy the new templates and update the feature registration call site.
- **Tasks:**
  1. **Update Copy Logic:** Modify the file copying logic within `cmd/petrock/feature.go` (likely in `runFeature` or a helper) to include copying `routes.go` and `http.go` from the skeleton to the new feature directory. Ensure placeholder replacement (`petrock_example_feature_name`, etc.) works for these new files.
  2. **Update `insertFeatureRegistration`:** Modify the `insertFeatureRegistration` function in `cmd/petrock/feature.go`:
      - Ensure it correctly identifies the `RegisterAllFeatures` function call within `cmd/<project>/features.go`.
      - Update the generated feature registration call to match the new signature expected by the updated `feature/register.go` template (including the `mux` argument and potentially others like `db`, `log`). Example: `featureName.RegisterFeature(mux, commands, queries, messageLog, appState, db)`
- **References:**
  - `cmd/petrock/feature.go`
- **Definition of Done:** Running `petrock feature <name>` correctly creates the `<name>/routes.go` and `<name>/http.go` files and inserts the updated `<name>.RegisterFeature(...)` call into `cmd/<project>/features.go`.

### Chunk 3: Update Feature Registration Template (`register.go`) - DONE

- **Goal:** Modify the feature's main registration function template to accommodate HTTP route registration.
- **Tasks:**
  1. **Update Signature:** In `internal/skeleton/feature_template/register.go`, change the `RegisterFeature` function signature to accept `mux *http.ServeMux` and other potential shared dependencies needed by `FeatureServer` (e.g., `*sql.DB`, `*core.MessageLog`, `*core.CommandRegistry`).
  2. **Initialize `FeatureServer`:** Add code to call `NewFeatureServer` (from `http.go`), passing the required dependencies (executor, querier, state, log, db, commands, etc.).
  3. **Call `RegisterRoutes`:** Add a call to `RegisterRoutes(mux, server)` (from `routes.go`), passing the router and the initialized `FeatureServer`.
  4. **Maintain Core Registrations:** Ensure the existing calls to `commands.Register`, `queries.Register`, and `RegisterTypes` remain.
- **References:**
  - `internal/skeleton/feature_template/register.go`
  - `internal/skeleton/feature_template/http.go` (for `NewFeatureServer` signature)
  - `docs/feature/register.go.md`
- **Definition of Done:** The `feature_template/register.go` file reflects the new signature and includes the logic to initialize `FeatureServer` and call `RegisterRoutes`.

### Chunk 4: Update Project `features.go` Template - DONE

- **Goal:** Update the template for the project-level feature registration function to pass the necessary dependencies.
- **Tasks:**
  1. **Update Signature:** In `internal/skeleton/cmd/petrock_example_project_name/features.go`, modify the `RegisterAllFeatures` function signature to accept `mux *http.ServeMux` and any other shared dependencies that features might need via `RegisterFeature` (e.g., `*sql.DB`).
  2. **Update Call Site Placeholder:** Ensure the placeholder comment `// petrock:register-feature` is positioned correctly within `RegisterAllFeatures` so that inserted `feature.RegisterFeature` calls receive the correct arguments (mux, commands, queries, log, state, db, etc.).
- **References:**
  - `internal/skeleton/cmd/petrock_example_project_name/features.go`
- **Definition of Done:** The `features.go` template has the updated `RegisterAllFeatures` signature and the registration placeholder is correctly positioned.

### Chunk 5: Update Project `serve.go` Template - DONE

- **Goal:** Ensure the main server setup passes the HTTP router to the feature registration process at the correct time.
- **Tasks:**
  1. **Pass Dependencies to `RegisterAllFeatures`:** In `internal/skeleton/cmd/petrock_example_project_name/serve.go`, locate the call to `RegisterAllFeatures`. Update it to pass the `mux *http.ServeMux` instance and any other newly required shared dependencies (like `db`).
  2. **Ensure Registration Order:** Verify that `RegisterAllFeatures` is called _after_ the core application routes (e.g., `/`, `/commands`, `/queries`) have been registered on the `mux`. This preserves the ability for features to override core routes.
- **References:**
  - `internal/skeleton/cmd/petrock_example_project_name/serve.go`
  - `internal/skeleton/cmd/petrock_example_project_name/features.go` (for `RegisterAllFeatures` signature)
  - `docs/cmd/serve.go.md`
- **Definition of Done:** The `serve.go` template correctly passes the `mux` and other dependencies to `RegisterAllFeatures` after core routes are defined.

### Chunk 6: Testing and Refinement - DONE

- **Goal:** Verify the entire workflow functions correctly.
- **Tasks:**
  1. **Run `petrock new ...`:** Create a new test project.
  2. **Run `petrock feature testfeature`:** Add a feature.
  3. **Inspect Generated Code:** Check that `testfeature/routes.go`, `testfeature/http.go` are created correctly. Verify `cmd/testproject/features.go` has the correct import and `testfeature.RegisterFeature` call with the right arguments.
  4. **Implement Basic Handler:** Add a simple `http.HandlerFunc` in `testfeature/http.go` (e.g., returning "Hello from feature") and register it in `testfeature/routes.go` (e.g., at `/testfeature/hello`).
  5. **Run `go run ./cmd/testproject serve`:** Start the server.
  6. **Test Endpoints:**
      - Verify core endpoints (`/`, `/commands`, `/queries`) still work.
      - Verify the new feature endpoint (`/testfeature/hello`) works.
  7. **Test Overriding (Optional):** Implement a handler for `/` in the feature and verify it overrides the core index page.
  8. **Refine Templates:** Adjust template code (`routes.go`, `http.go`) based on testing for clarity and correctness.
- **Definition of Done:** A project generated with `petrock new` and extended with `petrock feature` successfully runs, serving both core API endpoints and the new feature-specific HTTP endpoints.

## Documentation Updates

- All relevant documentation files (`docs/`) have already been updated or created in the previous step. This plan focuses solely on the code implementation.
