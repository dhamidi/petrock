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