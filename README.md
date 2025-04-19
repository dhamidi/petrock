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
*(Replace `github.com/dhamidi/petrock` with the actual repository path)*

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
5. Create a Git commit with the newly added feature files and modifications.

### Testing Petrock Itself

To run a self-test that verifies the `new` command and basic build process:

```bash
petrock test
```

This creates a temporary project, runs `petrock new`, and attempts to build it.

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
```

Refer to the generated code and the `docs/` directory within Petrock's repository for more in-depth details on the architecture and specific components.
