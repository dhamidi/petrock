Currently, the entire Go code is encoded as string templates inside the repository. Instead, I would like the Go code to be actual Go code using placeholder identifiers in a separate directory in the codebase.

Please create a plan that details the migration steps. We would need to take the current templates and convert them into actual Go code using placeholder identifiers.
# Migration Plan: Example Project Skeleton

This document outlines the steps to migrate the `petrock new` command from using `text/template` to an example project skeleton approach with search-and-replace.

## Goal

Replace the Go template rendering (`text/template`) with a fully functional, compilable example Go project located within the `petrock` repository (e.g., in `internal/skeleton/`). This example project will use hardcoded placeholder strings for the project name and module path. The `petrock new` command will copy this example project and perform search-and-replace operations on the copied files.

## Placeholders

The following placeholder strings will be used consistently throughout the example project skeleton files and directory names:

*   **Project Name Placeholder:** `petrock_example_project_name`
    *   Used within file content where the project name is needed (e.g., binary name, output messages).
    *   Used as the directory name for the command entry point (e.g., `cmd/petrock_example_project_name/`).
*   **Module Path Placeholder:** `github.com/petrock/example_module_path`
    *   Used in `go.mod` and import paths within `.go` files.

## Migration Steps

1.  **Create Example Project Skeleton Directory:**
    *   Create a new directory `internal/skeleton/` to hold the example project structure.

2.  **Populate Example Project Skeleton:**
    *   **Move & Rename Files:** Move all `.tmpl` files from `internal/template/templates/new/` (and its subdirectories `cmd`, `core`) into the corresponding location within `internal/skeleton/`. Rename them to remove the `.tmpl` extension (e.g., `internal/template/templates/new/core/view.go.tmpl` becomes `internal/skeleton/core/view.go`, `internal/template/templates/new/go.mod.tmpl` becomes `internal/skeleton/go.mod`).
    *   **Rename Placeholder Directory:** Rename the command directory from `internal/skeleton/cmd/{{ .ProjectName }}` (or similar if already moved) to `internal/skeleton/cmd/petrock_example_project_name`.
    *   **Replace Placeholders in Content:** Edit *all* files within `internal/skeleton/` (including `.go`, `go.mod`, `.gitignore`, etc.).
        *   Replace all occurrences of `{{ .ProjectName }}` with `petrock_example_project_name`.
        *   Replace all occurrences of `{{ .ModuleName }}` with `github.com/petrock/example_module_path`.

3.  **Verify Example Project Compiles:**
    *   Navigate into the `internal/skeleton/` directory.
    *   Run `go mod tidy`.
    *   Run `go build ./...`.
    *   Iteratively fix any compilation errors within the skeleton code until it builds successfully. This ensures the base copied by `petrock new` is valid Go code.

4.  **Update `cmd/petrock/new.go` (`runNew` function):**
    *   **Remove Template Rendering:** Remove all logic related to `text/template` parsing and rendering (`template.RenderTemplate`, `templatesToRender` map, `templateData`).
    *   **Implement Skeleton Copy:** Add logic to copy the entire `internal/skeleton/` directory structure to the new project directory (`projectName`). This copy *must* handle renaming the placeholder directory `internal/skeleton/cmd/petrock_example_project_name` to `cmd/<actual_project_name>` during the copy process.
    *   **Implement Search & Replace:** After copying, add logic to:
        *   Walk through *all* files (not just `.go`) within the newly created project directory.
        *   Read each file's content.
        *   Perform string replacement for `petrock_example_project_name` with the actual `projectName`.
        *   Perform string replacement for `github.com/petrock/example_module_path` with the actual `modulePath`.
        *   Write the modified content back to the file.

5.  **Refactor/Add Utility Functions:**
    *   **Directory Copy:** Ensure/Create a utility function (e.g., in `internal/utils/fs.go`) to recursively copy a directory structure, handling the specific placeholder directory rename (`cmd/petrock_example_project_name`).
    *   **Search & Replace:** Ensure/Create a utility function (e.g., in a new `internal/utils/replace.go` or within `fs.go`) that takes a file path and the placeholder mappings, reads the file, performs replacements efficiently, and writes it back.

6.  **Clean Up:**
    *   Remove the original `internal/template/templates/` directory and the `internal/template/template.go` file after confirming the new mechanism works reliably.

7.  **Testing:**
    *   Thoroughly test the `petrock new` command with various project/module names, including edge cases if any.
    *   Ensure the `petrock test` command still passes, as it relies on `petrock new`.
