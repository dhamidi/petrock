Currently, the entire Go code is encoded as string templates inside the repository. Instead, I would like the Go code to be actual Go code using placeholder identifiers in a separate directory in the codebase.

Please create a plan that details the migration steps. We would need to take the current templates and convert them into actual Go code using placeholder identifiers.
# Migration Plan: Go Templates to Skeleton Code

This document outlines the steps to migrate the `petrock new` command from using `text/template` for Go files to a skeleton-based approach with search-and-replace.

## Goal

Instead of rendering Go code using `text/template`, create a "skeleton" project structure containing actual Go files with placeholder identifiers. The `petrock new` command will then copy this skeleton and perform search-and-replace operations on the copied files.

## Placeholders

The following placeholder strings will be used in the skeleton Go files:

*   `PETROCK_PROJECT_NAME`: Replaced with the actual project name provided by the user.
*   `PETROCK_MODULE_NAME`: Replaced with the actual module path provided by the user.

## Migration Steps

1.  **Create Skeleton Directory:**
    *   Create a new directory `internal/skeleton/` to hold the template project structure.

2.  **Convert Go Templates to Skeleton Files:**
    *   **Move & Rename:** Move each `.go.tmpl` file from `internal/template/templates/new/` (and its subdirectories `cmd`, `core`) into the corresponding location within `internal/skeleton/`. Rename them to remove the `.tmpl` extension (e.g., `internal/template/templates/new/core/view.go.tmpl` becomes `internal/skeleton/core/view.go`).
    *   **Replace Placeholders:** Edit each of the newly moved `.go` files. Replace `{{ .ProjectName }}` with `PETROCK_PROJECT_NAME` and `{{ .ModuleName }}` with `PETROCK_MODULE_NAME`.
    *   **Directory Placeholder:** Rename the `internal/skeleton/cmd/{{ .ProjectName }}` directory (once moved) to `internal/skeleton/cmd/PETROCK_PROJECT_NAME`. The copy mechanism will need to handle renaming this directory based on the actual project name.

3.  **Handle Non-Go Templates (`.gitignore`, `go.mod`):**
    *   Keep `.gitignore.tmpl` and `go.mod.tmpl` in their current location (`internal/template/templates/new/`).
    *   Continue using the existing `text/template` rendering mechanism for these specific files.

4.  **Update `cmd/petrock/new.go` (`runNew` function):**
    *   **Remove Go Template Rendering Loop:** Remove the loop that iterates over `templatesToRender` for the `.go.tmpl` files.
    *   **Implement Skeleton Copy:** Add logic to copy the entire `internal/skeleton/` directory structure to the new project directory (`projectName`). This copy must handle renaming the `internal/skeleton/cmd/PETROCK_PROJECT_NAME` directory to `cmd/<actual_project_name>`.
    *   **Implement Search & Replace:** After copying, add logic to:
        *   Walk through all `.go` files within the newly created project directory.
        *   Read each file's content.
        *   Perform string replacement for `PETROCK_PROJECT_NAME` and `PETROCK_MODULE_NAME` with the actual `projectName` and `modulePath` values.
        *   Write the modified content back to the file.
    *   **Keep Non-Go Template Rendering:** Keep the logic that uses `template.RenderTemplate` for `.gitignore.tmpl` and `go.mod.tmpl`. Update the `templatesToRender` map to only include these two files.

5.  **Refactor/Add Utility Functions:**
    *   **Directory Copy:** Create a utility function (e.g., in `internal/utils/fs.go`) to recursively copy a directory structure, handling the placeholder directory rename.
    *   **Search & Replace:** Create a utility function (e.g., in a new `internal/utils/replace.go` or within `fs.go`) that takes a file path and the placeholder mappings, reads the file, performs replacements, and writes it back.
    *   **Template Rendering:** Review `internal/template/template.go`. It might still be used for the non-Go files, or its logic could be simplified/moved if only used in `runNew`.

6.  **Clean Up:**
    *   Remove the original `.go.tmpl` files from `internal/template/templates/new/` and its subdirectories after confirming the new mechanism works.

7.  **Testing:**
    *   Thoroughly test the `petrock new` command with various project/module names.
    *   Ensure the `petrock test` command still passes.
