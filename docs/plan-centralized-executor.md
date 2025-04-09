# Plan: Implement Centralized core.Executor

**Goal:** Refactor the command handling logic to use a central `core.Executor` that manages validation, logging, and state updates, ensuring consistency and simplifying feature code.

**Definition of Done:**

*   `core.Executor` struct and `Execute` method are implemented in `core/commands.go`, orchestrating validation, logging, and state updates.
*   `core.CommandValidator` interface is defined.
*   `core.CommandRegistry` is updated to store both state update handlers (`CommandHandler`) and feature validators (`CommandValidator`), providing methods like `GetHandlerAndValidator`. The `Dispatch` method is removed.
*   Feature `execute.go` files define a feature `Executor` struct (implementing `core.CommandValidator`) containing both validation methods (accessing state) and state update handler methods (matching `core.CommandHandler`).
*   Feature `register.go` files correctly instantiate the feature `Executor`, register both its state update handler methods and the executor instance itself (as the validator) with `core.CommandRegistry`, and pass the central `core.Executor` to the `FeatureServer`.
*   Feature `http.go` files use the central `core.Executor.Execute` method for all command-driven state changes, handling potential validation/logging errors returned by it.
*   `cmd/serve.go` initializes the central `core.Executor`, passes it during feature registration, uses it for the core `/commands` endpoint, and correctly replays the message log using only the state update handlers (retrieved via `GetHandler`), panicking on replay errors.
*   The `Validate()` method is removed from the `core.Command` interface definition in the documentation.
*   State update errors occurring *after* logging within `core.Executor.Execute` reliably cause a panic.
*   Relevant unit and integration tests are added or updated to cover the new execution flow, including validation logic within feature executors.
*   Existing documentation accurately reflects the final implementation (as per these updates).

---

## Step-by-Step Implementation Plan

**Step 1: Define `core.CommandValidator` Interface**

*   **Summary:** Define the `CommandValidator` interface in `core/commands.go`.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:** Add the `CommandValidator` interface definition with the `ValidateCommand(ctx context.Context, cmd Command) error` method.
*   **Definition of Done:** The interface exists. Code compiles.

**Step 2: Update `core.CommandRegistry`**

*   **Summary:** Modify `CommandRegistry` to store validators alongside handlers, update `Register`, and add `GetHandlerAndValidator`. Remove old `Dispatch`.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Add `validators map[string]CommandValidator` field.
    *   Update `NewCommandRegistry` to initialize the new map.
    *   Update `Register` signature to `(r *CommandRegistry) Register(cmd Command, handler CommandHandler, validator CommandValidator)`. Update implementation to store handler, validator, and type.
    *   Implement `GetHandlerAndValidator(name string) (CommandHandler, CommandValidator, bool)`.
    *   Implement `GetValidator(name string) (CommandValidator, bool)`.
    *   Remove the old `Dispatch` method.
*   **Definition of Done:** `CommandRegistry` stores and retrieves handlers and validators. `Register` has the new signature. `Dispatch` is removed. Code compiles.

**Step 3: Define `core.Executor` Struct and Constructor**

*   **Summary:** Define the central `Executor` struct and its constructor.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Add the `Executor` struct definition with fields for `*MessageLog` and `*CommandRegistry`.
    *   Add the `NewExecutor` constructor function.
*   **Definition of Done:** The `Executor` struct and `NewExecutor` function exist. Code compiles.

**Step 4: Implement `core.Executor.Execute` Method**

*   **Summary:** Implement the core execution logic (Get Handler/Validator -> Validate -> Log -> Execute Handler -> Panic on error) in `core.Executor.Execute`.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Implement `(e *Executor) Execute(ctx context.Context, cmd Command) error`.
    *   Use `registry.GetHandlerAndValidator` to get components. Handle not found error.
    *   Call `validator.ValidateCommand(ctx, cmd)`. Handle validation error.
    *   Call `log.Append(ctx, cmd)`. Handle logging error.
    *   Call `handler(ctx, cmd)`. Panic if it returns an error.
    *   Return `nil` on success.
*   **Definition of Done:** The `Execute` method implements the revised logic. Code compiles.

**Step 5: Refactor Feature `execute.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/execute.go`) to define the feature `Executor` (implementing `core.CommandValidator`) with both validation and state update methods.
*   **Files:** `internal/skeleton/feature_template/execute.go`
*   **Changes:**
    *   Define the feature `Executor` struct holding `*State`.
    *   Implement `NewExecutor` constructor.
    *   Implement `ValidateCommand` method with a type switch delegating to specific `validate<CommandType>` methods.
    *   Implement example `validate<CommandType>` methods that access `e.state`.
    *   Implement example `Handle<CommandType>` methods (matching `core.CommandHandler`) containing only state update logic using `e.state`.
    *   Update comments.
*   **Definition of Done:** The template `feature_template/execute.go` defines a feature `Executor` with separate validation and state update methods.

**Step 6: Refactor Feature `register.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/register.go`) to instantiate the feature `Executor` and register it correctly with the `core.CommandRegistry`.
*   **Files:** `internal/skeleton/feature_template/register.go`
*   **Changes:**
    *   Update `RegisterFeature` signature to accept `centralExecutor *core.Executor`.
    *   Instantiate the feature executor: `featureExecutor := NewExecutor(state)`.
    *   Update `commands.Register` calls to pass the command instance, the state update handler method (e.g., `featureExecutor.HandleCreate`), and the feature executor instance itself as the validator (e.g., `featureExecutor`).
    *   Update the `NewFeatureServer` call to pass the `centralExecutor`.
    *   Update comments.
*   **Definition of Done:** The template `feature_template/register.go` correctly instantiates the feature `Executor` and registers its components.

**Step 7: Refactor Feature `http.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/http.go`) to use the central `core.Executor` for write operations and handle its errors.
*   **Files:** `internal/skeleton/feature_template/http.go`
*   **Changes:**
    *   Ensure `FeatureServer` struct has `executor *core.Executor` field (representing the central executor).
    *   Update `NewFeatureServer` to accept and store the central `executor`.
    *   Modify write handlers (e.g., `HandleCreate`) to:
        *   Construct the command.
        *   Call `fs.executor.Execute(r.Context(), cmd)`.
        *   Handle potential validation or logging errors returned by `Execute` (e.g., return 400, 500).
    *   Remove direct logging or old validation calls.
*   **Definition of Done:** The template `feature_template/http.go` uses the central `core.Executor.Execute` for commands.

**Step 8: Refactor `cmd/serve.go` Template**

*   **Summary:** Update the template (`internal/skeleton/cmd/petrock_example_project_name/serve.go`) to initialize the central `core.Executor` and adjust the log replay logic.
*   **Files:** `internal/skeleton/cmd/petrock_example_project_name/serve.go`
*   **Changes:**
    *   In `runServe`:
        *   Instantiate the central `core.Executor`: `executor := core.NewExecutor(messageLog, commandRegistry)`.
        *   Update `RegisterAllFeatures` call to pass the `executor`.
        *   Modify the log replay loop:
            *   Load and decode messages.
            *   If it's a command, look up *only the handler* using `commandRegistry.GetHandler`. **Do not validate during replay.**
            *   Execute the handler: `err := handler(ctx, cmd)`.
            *   If `handler` returns an error, `panic`.
        *   Ensure the core `POST /commands` HTTP handler uses `executor.Execute`.
*   **Definition of Done:** The template `serve.go` initializes the central `executor`, passes it correctly, and replays the log using only state update handlers.

**Step 9: Remove `Validate()` Method from Command Templates**

*   **Summary:** Remove the now-unused `Validate()` method from example command structs in the feature template.
*   **Files:** `internal/skeleton/feature_template/messages.go`
*   **Changes:** Delete the `Validate() error` methods from command structs.
*   **Definition of Done:** Example command structs no longer have a `Validate()` method.

**Step 10: Apply Template Changes to Existing Projects (Manual or Tooling)**

*   **Summary:** Apply the changes from the updated templates (Steps 1-9) to existing Petrock projects.
*   **Files:** All `core/commands.go`, feature `execute.go`, `register.go`, `http.go`, `messages.go`, and the project's `cmd/<project>/serve.go`.
*   **Changes:** Apply the structural and logical changes.
*   **Definition of Done:** The target project code reflects the new architecture. Project compiles and starts.

**Step 11: Testing**

*   **Summary:** Add/update unit and integration tests.
*   **Files:** Test files (`_test.go`).
*   **Changes:**
    *   Add unit tests for `core.Executor.Execute`.
    *   Add unit tests for feature `Executor`'s `ValidateCommand` and specific `validate<Type>` methods, mocking state where necessary.
    *   Update unit tests for feature `Executor`'s `Handle<Type>` methods to ensure they only contain state logic.
    *   Update integration tests for HTTP endpoints to verify validation errors (400) and success cases.
    *   Update log replay tests.
*   **Definition of Done:** Tests pass.

**Step 12: Review and Refine**

*   **Summary:** Final code review and manual testing.
*   **Files:** All modified files.
*   **Changes:** Address review findings.
*   **Definition of Done:** Code reviewed, consistent, and behaves as expected.
