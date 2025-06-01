# Petrock - Go Application Generator

<p align="center">
  <img src="static/petrock-transparent.png" alt="Petrock Mascot" width="256" height="256">
</p>

Petrock is a command-line tool designed to bootstrap and develop Go web applications quickly. Think of it like Ruby on Rails, but for Go, focusing on generating boilerplate code so you can focus on your application's logic.

It generates applications built with an event-sourcing-inspired architecture using Go, SQLite for persistence, Gomponents for server-side HTML rendering, Tailwind CSS for styling, and optionally Stimulus/Hotwire for front-end interactivity.

## Philosophy

Petrock acts purely as a **code generator**. It creates files based on templates and relies heavily on Git for version control.

- **No Runtime Dependency:** Applications built with Petrock *do not* depend on Petrock itself at runtime. This favors code generation over adding external framework dependencies to your project.
- **Opinionated Structure:** It generates a specific project structure and technology stack based on these core ideas:
    - *Event Sourcing / Intent Capture:* The system logs user intent (commands) first and derives application state later, providing auditability and flexibility.
    - *Simplicity:* It aims to keep the number of moving pieces low, using standard Go libraries and well-established tools like SQLite and Gomponents.
- **Git Integration:** Petrock refuses to run if your Git working directory is not clean (no uncommitted changes). It automatically creates Git commits after generating code, ensuring changes are tracked. You can always revert unwanted changes using Git.
- **AI Ready:** Petrock is designed with AI-assisted development in mind:
    - *Modifiable Code:* It generates straightforward Go code that is relatively easy for AI tools to understand and modify.
    - *Contextual Documentation:* It generates detailed documentation (`docs/`) describing the plan for each generated file, which can be fed to coding agents as context.
    - *Cohesive Structure:* The enforced project and feature structure limits the scope an AI needs to inspect when working on specific parts of the application.

## Usage

### Prerequisites

- Go (latest stable version recommended)
- Git

### Installation

```bash
go install github.com/dhamidi/petrock@latest
```

### Creating a New Project

To start a new web application project:

```bash
petrock new <project-name> <go-module-path>
```

- `<project-name>`: The name of the directory to create for your project (e.g., `myblog`).
- `<go-module-path>`: The Go module path for your project (e.g., `github.com/me/myblog`).

**Example:**

```bash
petrock new myblog github.com/me/myblog
```

This command will:
1. Create a directory named `myblog`.
2. Initialize a Git repository within `myblog`.
3. Set up a Go module (`go.mod`) with the specified path.
4. Generate the initial project structure:
    - `cmd/myblog/`: Contains the main application entrypoint and subcommands (`serve`, `build`, `deploy`).
    - `core/`: Contains shared core components (command/query registries, logging, persistence, base views).
5. Install necessary Go dependencies (`go mod tidy`).
6. Create an initial Git commit with the generated files.

### Adding a Feature

Once you have a project, you can add features (domain-specific modules):

```bash
cd <project-name>
petrock feature <feature-name>
```

- `<feature-name>`: The name of the feature package (e.g., `posts`, `users`). Must be a valid Go package name (lowercase).

**Example:**

```bash
cd myblog
petrock feature posts
```

This command will:
1. Check if the Git workspace is clean.
2. Generate a new Go package directory named `posts`.
3. Populate the `posts/` directory with a modular, organized structure:
    - `main.go`: Main feature package exports and initialization.
    - `assets.go`: Asset registration and embedding.
    - `assets/`: Directory for static assets.
    - `commands/`: Command definitions and handlers:
      - Command interfaces, creation, update, deletion commands
      - Registration logic and command execution
    - `handlers/`: HTTP and command handling:
      - API endpoints and form handlers
      - Middleware and core handler functionality
    - `queries/`: Query definitions and handlers:
      - Query interfaces, single item and list queries
    - `state/`: In-memory state management:
      - Core item state and related metadata
    - `ui/`: UI components and templates:
      - Form, table, and layout components
      - Page views for listing, detail, and forms
    - `routes/`: Feature-specific HTTP route definitions:
      - API and web UI routes
    - `workers/`: Background job definitions and workers.
4. Automatically update `cmd/myblog/features.go` to import and register the new `posts` feature.
5. Run `go mod tidy` to update dependencies.
6. Create a Git commit with the newly added feature files and modifications.

### Testing Petrock Itself

To run a self-test that verifies the `new` and `feature` commands and basic build process:

```bash
petrock test
```

This creates a temporary project, runs `petrock new`, adds a test feature, attempts to build it, and tests basic HTTP functionality.

## Worker System

Petrock includes a powerful worker abstraction for background processing, event handling, and asynchronous operations. Workers eliminate boilerplate code while maintaining flexibility for complex business logic.

### Key Features

- **Command-Based Processing**: Workers respond to specific commands with focused handlers
- **Automatic Message Processing**: Core infrastructure handles message iteration and routing
- **State Management**: Built-in support for worker-specific state with persistence
- **Periodic Work**: Background processing that runs during each work cycle
- **Graceful Lifecycle**: Proper startup, shutdown, and error handling

### Worker Benefits

- **Reduced Boilerplate**: ~80% reduction in worker code compared to manual implementation
- **Consistent Infrastructure**: All workers use the same message processing and error handling
- **Easy Testing**: Command handlers are pure functions that are simple to test
- **Scalable**: Projects can easily support 50-100 workers without code duplication

### Example Worker

```go
func NewWorker(app *core.App, state *State, log *core.MessageLog, executor *core.Executor) core.Worker {
    worker := core.NewWorker(
        "Feature Worker",
        "Handles background processing for the feature",
        &WorkerState{pendingTasks: make(map[string]Task)},
    )
    
    worker.SetDependencies(log, executor)
    
    // Register command handlers
    worker.OnCommand("feature/process", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
        return handleProcessCommand(ctx, cmd, msg, worker.State().(*WorkerState))
    })
    
    // Set periodic work
    worker.SetPeriodicWork(func(ctx context.Context) error {
        return processPendingTasks(ctx, worker.State().(*WorkerState))
    })
    
    return worker
}
```

For complete documentation on workers, see [`docs/workers.md`](docs/workers.md).

## Key-Value Store

Petrock applications include a built-in key-value store for persistent data storage. The KVStore provides a simple interface for storing JSON-serializable data with SQLite backing.

### Features

- **JSON Serialization**: Automatic marshaling/unmarshaling of Go types
- **Pattern Matching**: List keys using SQLite GLOB patterns
- **CLI Access**: Built-in commands for get, set, and list operations
- **Worker Integration**: Used internally for worker position tracking

### Usage Examples

```bash
# Store configuration
go run ./cmd/myapp kv set --json "app:config" '{"theme": "dark", "debug": true}'

# Retrieve configuration
go run ./cmd/myapp kv get "app:config"

# List all configuration keys
go run ./cmd/myapp kv list "*:config"
```

For complete documentation on the KVStore, see [`docs/core/kv.md`](docs/core/kv.md).

## Generated Application

A Petrock-generated application includes its own command-line interface:

```bash
# Navigate into your generated project directory
cd <project-name>

# Run the development server (default: http://localhost:8080)
go run ./cmd/<project-name> serve

# Build a distributable binary
go run ./cmd/<project-name> build

# Deploy the binary (requires configuration)
go run ./cmd/<project-name> deploy --target-host user@hostname

# Inspect the application (view registered commands, queries, routes, etc.)
go run ./cmd/<project-name> self inspect

# Key-value store operations
go run ./cmd/<project-name> kv get <key>
go run ./cmd/<project-name> kv set <key> <value>
go run ./cmd/<project-name> kv set --json <key> <json-value>
go run ./cmd/<project-name> kv list [glob-pattern]
```

Refer to the generated code and the `docs/` directory within Petrock's repository for more in-depth details on the architecture and specific components.
