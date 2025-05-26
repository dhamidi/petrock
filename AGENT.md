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

- Template code in `internal/skeleton/` must be valid, compilable Go code that gets string-substituted during generation
- Gomponents use `g "maragu.dev/gomponents"`, dot-import components, `Classes{...}` maps, and `html.Style()` takes raw CSS strings
- CSS should be embedded as string constants in gomponents files, not external files
- Build process (`./build.sh`) validates the entire template→generation→compilation pipeline
- UI gallery components import core package using template placeholder `"github.com/petrock/example_module_path/core"`
- Template system uses simple string replacement, never string-based templating engines
- Integration tests include generating projects, adding features, and running live servers
- The user maintains the test server process - assume it runs or ask the user to run it when testing web functionality