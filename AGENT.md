# Commands for petrock Go project

## Build Commands
- Full build: `./build.sh`
- Single component: `./build.sh build_skeleton` or `./build.sh build_petrock`
- Run tests: `./build.sh`
- During builds, you might see the following error: `Error: git workspace is not clean`.  When you encounter this error, write a conventional commit message, committing all uncommitted files and then run the build command again.

## Lint Commands
- Lint specific files: `./lint.sh file1.go file2.go`
- Lint all Go files: `find . -name "*.go" | xargs ./lint.sh`

## Code Style Guidelines
- Follow Go's standard conventions for naming and formatting
- Use goimports for imports organization
- Shell scripts must follow rules in rules/shell.md (use main function, decompose logic, quote variables)
- Error handling: Always check errors and provide contextual messages
- Logging: Use log/slog package for structured logging
- Use cobra for CLI commands (see cmd/petrock/*.go for examples)

## Important project idiosyncracies

The project uses template code in internal/skeleton.

This project is a template that needs to be **valid Go code**.

Various placeholders in the code will be replaced via simple string substitution when petrock initializes a new project.

The template should **never** use string-based templating to generate code.

The following replacements are supported:

- `petrock_example_feature_name` – for the name of newly generated features (e.g. `posts`)
- `github.com/petrock/example_module_path` – the path to the finally generated module when using `petrock feature`.

## Project status

This project has no users yet – we are in the prototyping phase.

When making changes, ignore backwards compatibility.

# Lessons

## Template System
- Template code in `internal/skeleton/` must be valid, compilable Go code that gets string-substituted during generation
- Template system uses simple string replacement, never string-based templating engines
- All placeholder imports like `"github.com/petrock/example_module_path/core"` must be valid during skeleton compilation
- Build process (`./build.sh`) validates the entire template→generation→compilation pipeline

## Gomponents Patterns
- Gomponents use `g "maragu.dev/gomponents"`, dot-import components with `. "maragu.dev/gomponents/components"`
- Use `Classes{...}` maps from components package, not custom utilities - the map syntax `Classes{"class": true}` works directly
- Never create custom `Classes` functions that conflict with gomponents' built-in `Classes` type
- CSS should be embedded as string constants in gomponents files, not external files
- `html.Style()` takes raw CSS strings for inline styles
- Use `g.Text()` for text content, not bare `Text()` function calls

## UI Component Architecture
- UI components live in `internal/skeleton/core/ui/` with each component in its own file
- Gallery handlers live in `internal/skeleton/core/ui/gallery/` with routing in `component.go`
- Component demo pages use dedicated handlers (e.g., `HandleContainerDetail`) for rich examples
- Gallery imports both `core` and `core/ui` packages using template placeholders
- Use `core.Layout()` and `core.Page()` for consistent page structure in gallery
- Component props follow pattern: `type ComponentProps struct { Variant string; ... }`

## Design System Gallery
- Gallery accessible at `/_/ui` with component routes at `/_/ui/{component}`
- Gallery components import core package using template placeholder `"github.com/petrock/example_module_path/core"`
- Use color-coded examples in demo pages to visually distinguish component variants
- Include comprehensive documentation: description, usage examples, and properties table
- Gallery navigation uses sidebar with components grouped by category

## Build & Testing
- Integration tests include generating projects, adding features, and running live servers
- The user maintains the test server process - assume it runs or ask the user to run it when testing web functionality
- Use `./build.sh` to validate the full template compilation pipeline
- Gallery routes must be registered in `internal/skeleton/cmd/petrock_example_project_name/serve.go`

## Container Component Insights
- Container variants work well: default (896px), narrow (672px), wide (1280px), full (no limit)
- Use Tailwind utility classes: `mx-auto`, `px-4`, and responsive max-width classes
- Custom max-width override provides flexibility beyond predefined variants
- Visual examples with colored backgrounds effectively demonstrate width differences