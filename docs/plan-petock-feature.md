# Plan for `petrock feature` Command Implementation

This document outlines the steps required to implement the `petrock feature <feature_name>` command based on the project's high-level description and existing structure.

# High-level overview of `petrock feature`

The `petrock feature <feature_name>` command generates a new Go package within an existing Petrock-generated project. This package provides a basic structure for implementing a new application feature, including files for command/query handling, state management, views, and registration with the core application.

## Mode of operation

Following the principles outlined in `docs/high-level.md`:

1.  **Pre-checks:** The command first verifies it's running within a valid Petrock project directory and that the Git workspace is clean. It also checks if the target feature directory already exists.
2.  **Code Generation:** It copies a predefined feature skeleton (embedded within the Petrock binary) into the project, creating a new directory named after the feature.
3.  **Placeholder Replacement:** It replaces placeholders within the skeleton files (e.g., feature name, module path) with appropriate values.
4.  **Code Modification:** It modifies existing project code (specifically `cmd/<project_name>/features.go`) to import and register the new feature module using predefined markers.
5.  **Dependency Management:** It runs `go mod tidy` to ensure Go module dependencies are consistent.
6.  **Git Commit:** It automatically stages all generated and modified files and creates a new Git commit.

## Step 1: Define and Register the Cobra Command

Define the `feature` subcommand using Cobra and register it with the main `petrock` command.

**Details:**

1.  Create a new file: `cmd/petrock/feature.go`.
2.  In `feature.go`, define a `featureCmd` variable of type `*cobra.Command`.
    *   Set `Use` to `"feature [featureName]"`.
    *   Set `Short` and `Long` descriptions explaining its purpose (generating a new feature module).
    *   Set `Args` to `cobra.ExactArgs(1)` to require the feature name.
    *   Define a `RunE` function (e.g., `runFeature`) to contain the command's logic.
3.  Implement basic validation for the `featureName` argument within `runFeature` (e.g., using a regex similar to `dirNameRegex` in `new.go` to ensure it's a valid Go package name - lowercase, no hyphens).
4.  In `cmd/petrock/main.go`, add `featureCmd` to `rootCmd` within the `init` function.

**Done when:**

-   `cmd/petrock/feature.go` exists with the basic `featureCmd` definition and `runFeature` function.
-   `featureCmd` is registered in `cmd/petrock/main.go`.
-   Running `petrock feature` without arguments shows a Cobra error about missing arguments.
-   Running `petrock feature myfeature` executes the (currently empty) `runFeature` function.
-   Running `petrock feature "Invalid-Name!"` returns a validation error from `runFeature`.

**Files and references:**

-   `cmd/petrock/feature.go` (New file)
-   `cmd/petrock/main.go`
-   `cmd/petrock/new.go` (Example Cobra setup, validation regex)
-   `docs/high-level.md` (Command description)

## Step 2: Implement Pre-run Checks

Add checks at the beginning of `runFeature` to ensure the command runs in a valid environment.

**Details:**

1.  **Git Clean Check:** Call `utils.CheckCleanWorkspace()` and return its error if the workspace is not clean. This enforces the "ruthless override" principle.
2.  **Project Root Check:** Determine if the current directory is a Petrock project root.
    *   Check for the existence of `go.mod`.
    *   Check for the existence of `cmd/<project_name>/main.go` (derive `<project_name>` from the current directory or `go.mod`).
    *   Check for the existence of `core/`.
    *   Return a clear error if any check fails.
3.  **Feature Exists Check:** Check if a directory with the provided `featureName` already exists at the root level. If it does, return an error to prevent overwriting.

**Done when:**

-   `petrock feature myfeature` fails with a specific error if run in a directory with uncommitted Git changes.
-   `petrock feature myfeature` fails with a specific error if run outside a directory that looks like a Petrock project root.
-   `petrock feature myfeature` fails with a specific error if a directory named `myfeature` already exists in the project root.

**Files and references:**

-   `cmd/petrock/feature.go` (`runFeature` function)
-   `internal/utils/git.go` (`CheckCleanWorkspace`)
-   `internal/utils/fs.go` (Potentially needed for directory/file existence checks)
-   `internal/utils/gomod.go` (`GetModuleName` might be useful)
-   `docs/high-level.md` (Requirement for clean workspace)

## Step 3: Create and Embed Feature Skeleton

Define the file structure and content for a new feature as a valid Go package and embed it into the Petrock binary.

**Details:**

1.  Create a new directory `internal/skeleton/feature_template/`. This directory will contain a minimal, compilable Go package representing a feature.
2.  Inside `internal/skeleton/feature_template/`, create the necessary Go files based on the example structure in `docs/high-level.md` and the detailed plans in `docs/feature/*.md`.
    *   Use the placeholder `petrock_example_feature_name` consistently within these files where the actual feature name should appear (e.g., package name, function names, struct names, comments, file content).
    *   Use the placeholder `petrock_example_module_path` where the target project's Go module path is needed (e.g., in import paths within the feature or potentially in generated code/comments).
    *   Example files:
        *   `register.go` (package `petrock_example_feature_name`, defines `RegisterFeature`)
        *   `messages.go` (package `petrock_example_feature_name`, defines command/query structs like `CreatePostCommand`)
        *   `execute.go` (package `petrock_example_feature_name`, defines `Executor` and handlers like `HandleCreatePostCommand`)
        *   `query.go` (package `petrock_example_feature_name`, defines `Querier` and handlers)
        *   `state.go` (package `petrock_example_feature_name`, defines `State`)
        *   `jobs.go` (package `petrock_example_feature_name`)
        *   `view.go` (package `petrock_example_feature_name`)
        *   `assets.go` (package `petrock_example_feature_name`, embeds `assets/`)
    *   Create an empty directory `assets/` within `internal/skeleton/feature_template/`.
    *   Add a minimal `go.mod` file inside `internal/skeleton/feature_template/` declaring its own dummy module path (e.g., `module petrock_internal/feature_template`) and any direct dependencies the template needs (likely just the main project's `core` package placeholder: `require petrock_example_module_path/core v0.0.0`). This helps with validation but will be replaced/ignored during generation.
3.  In `embeds.go` (or `petrock.go`), add an `//go:embed` directive for `internal/skeleton/feature_template`.
    *   `//go:embed all:internal/skeleton/feature_template`
    *   Ensure this embedded FS is accessible, potentially alongside the main `SkeletonFS`.

**Done when:**

-   The `internal/skeleton/feature_template/` directory exists with valid Go files (`.go`), an `assets/` subdirectory, and a minimal `go.mod`.
-   The Go files contain basic structures derived from the `docs/feature/*.md` plans, using placeholders `petrock_example_feature_name` and `petrock_example_module_path` where appropriate.
-   The code within `internal/skeleton/feature_template/` is syntactically valid Go.
-   The feature skeleton is embedded into the binary via `//go:embed` and accessible via a variable (e.g., `petrock.SkeletonFS`).

**Files and references:**

-   `internal/skeleton/feature_template/` (New directory and files)
-   `docs/feature/*.md` (Source for skeleton content)
-   `embeds.go` (Or `petrock.go` - for embedding)
-   `docs/high-level.md` (Feature file structure)
-   `cmd/petrock/new.go` (Example of using placeholders in valid Go files)

## Step 4: Implement Skeleton Copying

Add logic to `runFeature` to copy the embedded feature skeleton into the target project directory.

**Details:**

1.  Get the target project's module path using `utils.GetModuleName(".")`.
2.  Define the source path within the embedded FS (e.g., `internal/skeleton/feature_template`).
3.  Define the destination path (e.g., `./<featureName>`).
4.  Use `utils.CopyDir` (or similar logic) to copy files from the embedded FS source path to the destination path.
    *   Ensure `CopyDir` handles embedded FS correctly.
    *   `CopyDir` will copy `go.mod.skel` from the template.
5.  After copying, rename the `go.mod.skel` file in the destination directory (`./<featureName>/go.mod.skel`) to `go.mod`.

**Done when:**

-   Running `petrock feature myfeature` in a valid project creates a `./myfeature` directory.
-   The `./myfeature` directory contains all the `.go` files, the `assets/` subdirectory, and a `go.mod` file (renamed from `go.mod.skel`).

**Files and references:**

-   `cmd/petrock/feature.go` (`runFeature` function)
-   `internal/utils/fs.go` (`CopyDir`)
-   `internal/utils/gomod.go` (`GetModuleName`)
-   `embeds.go` (Or `petrock.go` - provides the embedded FS)

## Step 5: Implement Placeholder Replacement

Add logic to replace placeholders in the newly copied feature files.

**Details:**

1.  Define the placeholder map:
    *   `"petrock_example_feature_name": featureName`
    *   `"petrock_example_module_path": modulePath` (obtained in Step 4)
2.  Use `utils.ReplaceInFiles` (or similar logic) targeting the newly created feature directory (`./<featureName>`). Pass the placeholder map.
3.  Ensure `ReplaceInFiles` correctly handles file contents and permissions.

**Done when:**

-   After `petrock feature myfeature` runs, files within `./myfeature` (like `register.go`, `messages.go`, etc.) have `petrock_example_feature_name` replaced with `myfeature`.
-   Files within `./myfeature` have `petrock_example_module_path` replaced with the actual module path of the target project.

**Files and references:**

-   `cmd/petrock/feature.go` (`runFeature` function)
-   `internal/utils/fs.go` (`ReplaceInFiles`)

## Step 6: Implement Feature Registration in Project Code

Modify the target project's `cmd/<project_name>/features.go` file to import and register the new feature.

**Details:**

1.  Determine the project name (e.g., from the current directory name or `go.mod`).
2.  Construct the path to the target file: `cmd/<project_name>/features.go`.
3.  Read the content of `features.go`.
4.  **Add Import:**
    *   Locate the `// petrock:import-feature` marker line.
    *   Insert a new import line *before* the marker: `_ "module/path/featureName"` (using the project's module path and the new feature name). Use `_` alias initially if direct usage isn't immediately needed, or determine the correct alias if necessary. Consider adding the feature name as an alias: `featureName "module/path/featureName"`.
5.  **Add Registration Call:**
    *   Locate the `// petrock:register-feature` marker line.
    *   Insert a new registration call *before* the marker: `featureName.RegisterFeature(commands, queries /*, messageLog, state... */)`
    *   Ensure the arguments passed match the expected signature in the skeleton's `register.go.skel`.
6.  Write the modified content back to `cmd/<project_name>/features.go`. This requires careful string manipulation or potentially using Go's AST parser for robustness. Simple string insertion based on markers is likely sufficient initially.

**Done when:**

-   Running `petrock feature myfeature` modifies `cmd/<project_name>/features.go`.
-   The modified file includes a new import line for `module/path/myfeature`.
-   The modified file includes a new line calling `myfeature.RegisterFeature(...)` within the `RegisterAllFeatures` function body.

**Files and references:**

-   `cmd/petrock/feature.go` (`runFeature` function)
-   `internal/skeleton/cmd/petrock_example_project_name/features.go` (Shows the target structure and markers)
-   `internal/skeleton/feature_template/register.go` (Defines the function signature to be called)
-   Go `os` package (for file reading/writing)
-   Go `strings` package (for manipulation)
-   (Optional) Go `go/parser` and `go/ast` packages for more robust code modification.

## Step 7: Run Go Mod Tidy

Execute `go mod tidy` in the project directory to update dependencies.

**Details:**

1.  Call `utils.GoModTidy(".")` from within `runFeature` after modifying the source files.

**Done when:**

-   Running `petrock feature myfeature` executes `go mod tidy` in the project root directory.
-   The project's `go.mod` and `go.sum` files are updated if necessary.

**Files and references:**

-   `cmd/petrock/feature.go` (`runFeature` function)
-   `internal/utils/gomod.go` (`GoModTidy`)

## Step 8: Create Git Commit

Stage all generated and modified files and create a Git commit.

**Details:**

1.  Call `utils.GitAddAll(".")` to stage all changes (new feature directory, modified `features.go`, `go.mod`, `go.sum`).
2.  Create a commit message, e.g., `"feat: Add feature '<featureName>' generated by petrock"`.
3.  Call `utils.GitCommit(".", commitMessage)`.

**Done when:**

-   Running `petrock feature myfeature` creates a new Git commit.
-   `git status` shows a clean working directory after the command finishes.
-   `git log` shows the new commit with the generated message.

**Files and references:**

-   `cmd/petrock/feature.go` (`runFeature` function)
-   `internal/utils/git.go` (`GitAddAll`, `GitCommit`)
-   `docs/high-level.md` (Requirement for automatic commits)

## Step 9: Final Output and Cleanup

Provide informative output to the user upon successful completion.

**Details:**

1.  Print success messages indicating the feature was created and committed.
2.  Consider adding hints about next steps (e.g., "Implement handlers in `./<featureName>/execute.go` and `./<featureName>/query.go`").
3.  Ensure proper error handling throughout `runFeature`, returning informative errors wrapped with context.

**Done when:**

-   `petrock feature myfeature` prints clear success messages upon completion.
-   Errors encountered during any step are reported clearly to the user.

**Files and references:**

-   `cmd/petrock/feature.go` (`runFeature` function)
-   Go `fmt` package
-   Go `log/slog` package
