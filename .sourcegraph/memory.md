# Commands for petrock Go project

## Development workflow:

1. Make changes
2. Commit changes with a auto-generated commit messages following conventional commits: `git add <changed-or-added-files> && git commit -m $'conventional commit message with multiple lines'
3. Run tests to verify everything works

## Build Commands
- Full build: `./build.sh`
- Single component: `./build.sh build_skeleton` or `./build.sh build_petrock`
- Run tests: `./build.sh --debug`
- Single test: `go test -v ./path/to/package -run TestName`

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