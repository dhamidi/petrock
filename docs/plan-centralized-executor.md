# Plan: Implement Centralized core.Executor

**Goal:** Refactor the command handling logic to use a central `core.Executor` that manages validation, logging, and state updates, ensuring consistency and simplifying feature code.

**Definition of Done:**

*   `core.Executor` struct and `Execute` method are implemented in `core/commands.go`.
*   `core.CommandRegistry` is updated to store state update handlers and provide `GetHandler`. The `Dispatch` method is removed.
*   Feature `execute.go` files contain only state update logic in handlers matching the `core.CommandHandler` signature. The `Executor` struct is renamed to `StateUpdater`.
*   Feature `register.go` files correctly instantiate `StateUpdater`, register state update handlers with `core.CommandRegistry`, and pass the central `core.Executor` to the `FeatureServer`.
*   Feature `http.go` files use the central `core.Executor.Execute` method for all command-driven state changes, handling potential validation/logging errors.
*   `cmd/serve.go` initializes the `core.Executor`, passes it during feature registration, uses it for the core `/commands` endpoint, and correctly replays the message log using state update handlers, panicking on replay errors.
*   Optional `Validate()` methods are added to command structs where needed and are called by `core.Executor`.
*   State update errors occurring *after* logging within `core.Executor.Execute` reliably cause a panic.
*   Relevant unit and integration tests are added or updated to cover the new execution flow.
*   Existing documentation accurately reflects the final implementation (as per commit 760b63e).

---

## Step-by-Step Implementation Plan

**Step 1: Define `core.Executor` Struct and Constructor**

*   **Summary:** Create the basic `Executor` struct definition and its constructor function in `core/commands.go`.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Add the `Executor` struct definition with fields for `*MessageLog` and `*CommandRegistry`.
    *   Add the `NewExecutor` constructor function.
*   **Definition of Done:** The `Executor` struct and `NewExecutor` function exist in `core/commands.go`. The code compiles.

**Step 2: Implement `core.Executor.Execute` Method**

*   **Summary:** Add the core execution logic (Validate -> Log -> Get Handler -> Execute Handler -> Panic on error) to the `Execute` method of the `core.Executor`.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Implement the `(e *Executor) Execute(ctx context.Context, cmd Command) error` method according to the documented logic.
    *   Include checks for the optional `Validate()` method on the command.
    *   Include the call to `log.Append`.
    *   Include the lookup of the handler via `registry.GetHandler`.
    *   Include the execution of the handler.
    *   Include the panic logic if the handler returns an error.
*   **Definition of Done:** The `Execute` method is implemented in `core/commands.go` with the specified logic. The code compiles.

**Step 3: Update `core.CommandRegistry`**

*   **Summary:** Modify the `CommandRegistry` to store state update handlers, add `GetHandler`, and remove the old `Dispatch` method.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Ensure the `handlers` map in `CommandRegistry` stores `CommandHandler`.
    *   Implement the `GetHandler(name string) (CommandHandler, bool)` method.
    *   Remove the `Dispatch(ctx context.Context, cmd Command) error` method.
    *   Verify the `Register` method correctly stores the `CommandHandler`.
*   **Definition of Done:** `CommandRegistry` correctly stores and retrieves `CommandHandler` functions, and the `Dispatch` method is removed. The code compiles.

**Step 4: Refactor Feature `execute.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/execute.go`) to reflect the new pattern: rename `Executor` to `StateUpdater` and simplify handlers to only perform state updates.
*   **Files:** `internal/skeleton/feature_template/execute.go`
*   **Changes:**
    *   Rename the `Executor` struct to `StateUpdater`.
    *   Rename the constructor `NewExecutor` to `NewStateUpdater`.
    *   Remove any validation or logging logic from the example command handlers (`HandleCreate`, `HandleUpdate`, `HandleDelete`).
    *   Ensure handler signatures match `core.CommandHandler`.
    *   Ensure handlers interact correctly with the `state *State` dependency for updates.
    *   Update comments to reflect the new role (state updates only, panic on error).
*   **Definition of Done:** The template `feature_template/execute.go` uses `StateUpdater` and contains only state update logic in its handlers.

**Step 5: Refactor Feature `register.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/register.go`) to accept `*core.Executor`, instantiate `StateUpdater`, register the correct handlers, and pass `core.Executor` to the `FeatureServer`.
*   **Files:** `internal/skeleton/feature_template/register.go`
*   **Changes:**
    *   Add `executor *core.Executor` to the `RegisterFeature` function signature.
    *   Instantiate `StateUpdater` using `NewStateUpdater(state)`.
    *   Update `commands.Register` calls to pass the methods from the `StateUpdater` instance (e.g., `updater.HandleCreate`).
    *   Update the `NewFeatureServer` call to pass the `executor` dependency.
    *   Update comments/documentation within the function.
*   **Definition of Done:** The template `feature_template/register.go` correctly wires the `core.Executor` and registers state update handlers.

**Step 6: Refactor Feature `http.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/http.go`) to use the central `core.Executor` for write operations.
*   **Files:** `internal/skeleton/feature_template/http.go`
*   **Changes:**
    *   Ensure the `FeatureServer` struct has an `executor *core.Executor` field.
    *   Update the `NewFeatureServer` constructor signature and logic to accept and store the `executor`.
    *   Modify example HTTP handlers that perform writes (e.g., `HandleCreate`, `HandleUpdate`, `HandleDelete`) to:
        *   Construct the appropriate command struct.
        *   Call `fs.executor.Execute(r.Context(), cmd)`.
        *   Handle potential errors returned by `Execute` (validation, logging errors) by returning appropriate HTTP status codes (e.g., 400, 500). State update errors will cause a panic before returning.
    *   Remove any direct calls to `log.Append` or old feature executor methods from these handlers.
*   **Definition of Done:** The template `feature_template/http.go` uses `core.Executor.Execute` for commands initiated via HTTP handlers.

**Step 7: Refactor `cmd/serve.go` Template**

*   **Summary:** Update the application startup logic in the template (`internal/skeleton/cmd/petrock_example_project_name/serve.go`) to initialize and use the `core.Executor`, and correctly replay the log.
*   **Files:** `internal/skeleton/cmd/petrock_example_project_name/serve.go`
*   **Changes:**
    *   In `runServe`:
        *   Instantiate `core.Executor` after initializing the log and command registry: `executor := core.NewExecutor(messageLog, commandRegistry)`.
        *   Update the call to `RegisterAllFeatures` (or equivalent feature registration logic) to pass the `executor` instance.
        *   Modify the log replay loop:
            *   Load raw messages using `messageLog.Load`.
            *   Decode each message using `messageLog.Decode`.
            *   Type-assert the decoded message to `core.Command`.
            *   If it's a command, look up the handler using `commandRegistry.GetHandler`.
            *   Execute the handler: `err := handler(ctx, cmd)`.
            *   If `handler` returns an error, `panic`.
        *   Update the core `POST /commands` HTTP handler (if defined directly in `serve.go` or a related core http file) to use `executor.Execute`.
*   **Definition of Done:** The template `serve.go` correctly initializes `core.Executor`, uses it for the core command endpoint, passes it to features, and implements the correct log replay logic with panics on handler errors.

**Step 8: Add `Validate()` Method to Command Templates**

*   **Summary:** Add the optional `Validate()` method to the example command structs in the feature template (`internal/skeleton/feature_template/messages.go`).
*   **Files:** `internal/skeleton/feature_template/messages.go`
*   **Changes:**
    *   Add a `Validate() error` method to `CreateCommand`, `UpdateCommand`, etc.
    *   Implement basic example validation logic within these methods (e.g., checking for empty strings).
    *   Add comments explaining that this method is optional but will be called by `core.Executor`.
*   **Definition of Done:** Example command structs in the template have a `Validate()` method.

**Step 9: Apply Template Changes to Existing Projects (Manual or Tooling)**

*   **Summary:** Regenerate or manually update existing Petrock-based projects to incorporate the changes made to the templates in the previous steps. This involves running `petrock feature <name>` again (if safe) or carefully applying the diffs.
*   **Files:** All `execute.go`, `register.go`, `http.go`, `messages.go` in existing features, and the project's `cmd/<project>/serve.go`. Also `core/commands.go`.
*   **Changes:** Apply the structural and logical changes from Steps 1-8 to the actual project code.
*   **Definition of Done:** The target project code reflects the new `core.Executor` architecture. The project compiles and the basic server starts.

**Step 10: Testing**

*   **Summary:** Add or update unit and integration tests to verify the new execution flow.
*   **Files:** Test files (`_test.go`) for `core`, features, and `cmd/serve`.
*   **Changes:**
    *   Add unit tests for `core.Executor.Execute`, covering validation success/failure, logging success/failure, handler success, and handler failure (panic).
    *   Update unit tests for feature state updaters (`StateUpdater`) to ensure they only contain state logic.
    *   Update integration tests for HTTP endpoints (feature-specific and core `/commands`) to ensure they correctly trigger the executor and handle responses/errors.
    *   Add tests for the log replay mechanism in `cmd/serve.go` to ensure state is rebuilt correctly and errors during replay cause panics.
*   **Definition of Done:** New and updated tests pass, providing confidence in the refactored implementation.

**Step 11: Review and Refine**

*   **Summary:** Perform a final code review, check for TODOs, ensure consistency, and manually test the application flow.
*   **Files:** All modified files.
*   **Changes:** Address any findings from the review. Ensure logging provides adequate information about the execution flow.
*   **Definition of Done:** Code is reviewed, consistent, and behaves as expected according to the new design. The implementation is complete.
