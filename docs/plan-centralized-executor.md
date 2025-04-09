# Plan: Implement Centralized core.Executor with Command-Based Validation

**Goal:** Refactor command handling to use a central `core.Executor`. Validation logic needing state access will reside in `Validate(state *State)` methods on the command structs themselves (via a `Validator` interface). The feature executor will bridge the call from the central executor to the command's `Validate` method.

**Definition of Done:**

*   `core.Executor` struct and `Execute` method are implemented in `core/commands.go`, orchestrating validation, logging, and state updates.
*   `core.FeatureExecutor` interface is defined, requiring a `ValidateCommand` method.
*   `core.CommandRegistry` is updated to store state update handlers (`CommandHandler`) and feature executor instances (`FeatureExecutor`), providing methods like `GetHandlerAndFeatureExecutor`. The `Dispatch` method is removed.
*   Feature `messages.go` files define a `Validator` interface (`Validate(state *State) error`) and command structs implement this interface where stateful validation is needed.
*   Feature `execute.go` files define a feature `Executor` struct (implementing `core.FeatureExecutor`) which holds state and contains state update handler methods. Its `ValidateCommand` method checks if a command implements `Validator` and calls its `Validate` method with the state.
*   Feature `register.go` files correctly instantiate the feature `Executor`, register its state update handler methods and the executor instance itself with `core.CommandRegistry`, and pass the central `core.Executor` to the `FeatureServer`.
*   Feature `http.go` files use the central `core.Executor.Execute` method for commands, handling potential validation/logging errors.
*   `cmd/serve.go` initializes the central `core.Executor`, passes it during feature registration, uses it for the core `/commands` endpoint, and correctly replays the message log using only state update handlers, panicking on replay errors.
*   State update errors occurring *after* logging within `core.Executor.Execute` reliably cause a panic.
*   Relevant unit and integration tests cover the new execution flow, including validation logic within command structs.
*   Existing documentation accurately reflects this implementation.

---

## Step-by-Step Implementation Plan

**Step 1: Define `core.FeatureExecutor` Interface**

*   **Summary:** Define the `FeatureExecutor` interface in `core/commands.go` that feature executors must implement.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:** Add `type FeatureExecutor interface { ValidateCommand(ctx context.Context, cmd Command) error }`.
*   **Definition of Done:** The interface exists. Code compiles.

**Step 2: Update `core.CommandRegistry`**

*   **Summary:** Modify `CommandRegistry` to store feature executor instances instead of validators. Update `Register` and related methods. Remove old `Dispatch`.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Replace `validators map[string]CommandValidator` with `featureExecutors map[string]FeatureExecutor`.
    *   Update `NewCommandRegistry` to initialize `featureExecutors`.
    *   Update `Register` signature to `(r *CommandRegistry) Register(cmd Command, handler CommandHandler, featureExecutor FeatureExecutor)`. Update implementation to store handler, feature executor, and type.
    *   Implement `GetHandlerAndFeatureExecutor(name string) (CommandHandler, FeatureExecutor, bool)`.
    *   Implement `GetFeatureExecutor(name string) (FeatureExecutor, bool)`.
    *   Remove `GetValidator` and `GetHandlerAndValidator`.
    *   Remove the old `Dispatch` method.
*   **Definition of Done:** `CommandRegistry` stores and retrieves handlers and feature executors. `Register` has the new signature. `Dispatch` is removed. Code compiles.

**Step 3: Define `core.Executor` Struct and Constructor**

*   **Summary:** Define the central `Executor` struct and its constructor. (No change from previous plan step, but needed for sequence).
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Add the `Executor` struct definition with fields for `*MessageLog` and `*CommandRegistry`.
    *   Add the `NewExecutor` constructor function.
*   **Definition of Done:** The `Executor` struct and `NewExecutor` function exist. Code compiles.

**Step 4: Implement `core.Executor.Execute` Method**

*   **Summary:** Implement the core execution logic: Get Handler/FeatureExecutor -> Call FeatureExecutor.ValidateCommand -> Log -> Execute Handler -> Panic on error.
*   **Files:** `internal/skeleton/core/commands.go` (and corresponding generated project file `core/commands.go`)
*   **Changes:**
    *   Implement `(e *Executor) Execute(ctx context.Context, cmd Command) error`.
    *   Use `registry.GetHandlerAndFeatureExecutor` to get components. Handle not found error.
    *   Call `featureExecutor.ValidateCommand(ctx, cmd)`. Handle validation error.
    *   Call `log.Append(ctx, cmd)`. Handle logging error.
    *   Call `handler(ctx, cmd)`. Panic if it returns an error.
    *   Return `nil` on success.
*   **Definition of Done:** The `Execute` method implements the revised logic using `FeatureExecutor`. Code compiles.

**Step 5: Define `Validator` Interface in Feature Template**

*   **Summary:** Define the `Validator` interface (`Validate(state *State) error`) in the feature message template file.
*   **Files:** `internal/skeleton/feature_template/messages.go`
*   **Changes:** Add the `Validator` interface definition.
*   **Definition of Done:** The `Validator` interface is defined in the template.

**Step 6: Update Command Structs in Feature Template**

*   **Summary:** Update example command structs in the feature message template to optionally implement the `Validator` interface with state-aware logic.
*   **Files:** `internal/skeleton/feature_template/messages.go`
*   **Changes:** Add example `Validate(state *State) error` methods to command structs like `CreateCommand`, `UpdateCommand`, `DeleteCommand`, demonstrating access to `state`.
*   **Definition of Done:** Example commands in the template demonstrate implementing the `Validator` interface.

**Step 7: Refactor Feature `execute.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/execute.go`) to define the feature `Executor` (implementing `core.FeatureExecutor`) with its `ValidateCommand` bridge method and state update handlers.
*   **Files:** `internal/skeleton/feature_template/execute.go`
*   **Changes:**
    *   Define the feature `Executor` struct holding `*State`.
    *   Implement `NewExecutor` constructor.
    *   Implement the `ValidateCommand(ctx context.Context, cmd core.Command) error` method. This method performs the type assertion `cmd.(Validator)` and calls `validator.Validate(e.state)` if the assertion succeeds.
    *   Implement example `Handle<CommandType>` methods (matching `core.CommandHandler`) containing only state update logic using `e.state`.
    *   Update comments.
*   **Definition of Done:** The template `feature_template/execute.go` defines a feature `Executor` implementing `core.FeatureExecutor` with the correct `ValidateCommand` logic and separate state update handlers.

**Step 8: Refactor Feature `register.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/register.go`) to instantiate the feature `Executor` and register it correctly with the `core.CommandRegistry`.
*   **Files:** `internal/skeleton/feature_template/register.go`
*   **Changes:**
    *   Update `RegisterFeature` signature to accept `centralExecutor *core.Executor`.
    *   Instantiate the feature executor: `featureExecutor := NewExecutor(state)`.
    *   Update `commands.Register` calls to pass the command instance, the state update handler method (e.g., `featureExecutor.HandleCreate`), and the feature executor instance itself (e.g., `featureExecutor`).
    *   Update the `NewFeatureServer` call to pass the `centralExecutor`.
    *   Update comments.
*   **Definition of Done:** The template `feature_template/register.go` correctly instantiates the feature `Executor` and registers its components with the updated `Register` signature.

**Step 9: Refactor Feature `http.go` Template**

*   **Summary:** Update the feature template (`internal/skeleton/feature_template/http.go`) to use the central `core.Executor` for write operations and handle its errors. (No change from previous plan step, but needed for sequence).
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

**Step 10: Refactor `cmd/serve.go` Template**

*   **Summary:** Update the template (`internal/skeleton/cmd/petrock_example_project_name/serve.go`) to initialize the central `core.Executor` and ensure log replay uses only state handlers. (No change from previous plan step, but needed for sequence).
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

**Step 11: Apply Template Changes to Existing Projects (Manual or Tooling)**

*   **Summary:** Apply the changes from the updated templates (Steps 1-10) to existing Petrock projects.
*   **Files:** All `core/commands.go`, feature `execute.go`, `register.go`, `http.go`, `messages.go`, and the project's `cmd/<project>/serve.go`.
*   **Changes:** Apply the structural and logical changes.
*   **Definition of Done:** The target project code reflects the new architecture. Project compiles and starts.

**Step 12: Testing**

*   **Summary:** Add/update unit and integration tests.
*   **Files:** Test files (`_test.go`).
*   **Changes:**
    *   Add unit tests for `core.Executor.Execute`.
    *   Add unit tests for command structs' `Validate` methods, mocking state where necessary.
    *   Add unit tests for feature `Executor`'s `ValidateCommand` bridge method.
    *   Update unit tests for feature `Executor`'s `Handle<Type>` methods.
    *   Update integration tests for HTTP endpoints to verify validation errors (400) and success cases originating from command validation.
    *   Update log replay tests.
*   **Definition of Done:** Tests pass.

**Step 13: Review and Refine**

*   **Summary:** Final code review and manual testing.
*   **Files:** All modified files.
*   **Changes:** Address review findings.
*   **Definition of Done:** Code reviewed, consistent, and behaves as expected.
