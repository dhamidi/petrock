# Commands for petrock Go project

## Build Commands

- Full build: `./build.sh`
- Single component: `./build.sh build_skeleton` or `./build.sh build_petrock`
- Run tests: `./build.sh`
- During builds, you might see the following error: `Error: git workspace is not clean`. When you encounter this error, write a conventional commit message, committing all uncommitted files and then run the build command again.

## Lint Commands

- Lint specific files: `./lint.sh file1.go file2.go`
- Lint all Go files: `find . -name "*.go" | xargs ./lint.sh`

## Code Style Guidelines

- Follow Go's standard conventions for naming and formatting
- Use goimports for imports organization
- Shell scripts must follow rules in rules/shell.md (use main function, decompose logic, quote variables)
- Error handling: Always check errors and provide contextual messages
- Logging: Use log/slog package for structured logging
- Use cobra for CLI commands (see cmd/petrock/\*.go for examples)

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
- **IMPORTANT**: Use `ui.CSSClass("class1", "class2", ...)` pattern, NOT `Classes{...}` maps - the old Classes syntax causes compilation errors
- Never create custom `Classes` functions that conflict with gomponents' built-in `Classes` type
- CSS should be embedded as string constants in gomponents files, not external files
- Only TailwindCSS is allowed, no raw CSS
- `html.Style()` takes raw CSS strings for inline styles
- Use `g.Text()` for text content, not bare `Text()` function calls
- Use `html.THead` and `html.TBody` (capitalized) for table elements

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
- Use `BuildSidebar()` function in gallery.go to create consistent sidebar navigation across all component pages
- Individual component pages must implement their own sidebar layout using flex pattern: `ui.CSSClass("flex", "min-h-screen", "-mx-4", "-mt-4")`
- Component pages need full sidebar navigation, not just back links - users expect to navigate directly between components

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

## Grid Component Insights

- Grid component supports flexible CSS Grid layouts with customizable columns, gaps, and grid template areas
- Use CSS Grid properties: `grid-template-columns`, `gap`, `grid-template-areas` for layout control
- GridProps pattern: `Columns`, `Gap`, and `Areas` properties provide comprehensive grid configuration
- GridItem helper function enables named grid areas for complex layouts
- Default values work well: `Columns: "1fr"` (single column), `Gap: "1rem"` (consistent spacing)
- Common patterns: `"repeat(3, 1fr)"` for equal columns, `"200px 1fr 100px"` for mixed sizing
- Grid areas enable semantic layout: `"header header header" "sidebar main aside" "footer footer footer"`

## Form Component Insights

- Form components use `CSSClass()` directly as g.Node attributes, NOT wrapped in `html.Class()` - `html.Class()` expects strings
- HTML element naming: use `html.Textarea` not `html.TextArea`, `html.Input` for text inputs, `html.Select` for dropdowns
- Type conversions required: `html.Rows(strconv.Itoa(rows))` - HTML attributes expect strings, not integers
- Validation state pattern: use `ValidationState` string field with values "valid", "invalid", "pending" for consistent styling
- Import management: avoid unused imports like `. "maragu.dev/gomponents/components"` in form components
- Gallery styling: use `ui.CSSClass()` directly as first argument to `html.Div()`, not wrapped in `html.Class()`
- Component structure: TextInput/TextArea/Select props follow consistent pattern with Type, Placeholder, Value, ValidationState, Required, Disabled fields

