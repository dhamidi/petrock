# Description

Petrock is a Go command line tool which bootstraps new Go projects and generates new files in a Go project.

The target application is a web application that is using event sourcing as its basic idea.

The technology stack includes Tailwind CSS, Gomponents, SQLite and Stimulus and Hotwire for the front-end.

## Mode of operation

Petrock is a code generator. It ruthlessly overrides files relying on Git to be able to restore any unwanted changes.

Basically, Petrock refuses to run if there are any uncommitted or changed files in Git.

It then spits out a bunch of files based on the command that was invoked and automatically commits them.

The basic idea is that applications that are built with Petrock do not depend on Petrock at runtime, only at development time.

The Petrock binary has all of these templates compiled in and just extracts them when invoked.

It automatically creates Git commits.

## Example use

### Creating a new project

```sh
petrock new blog
```

This creates a new directory called blog, initializes a git repository in there, sets up a go module with all necessary dependencies installed.

The basic structure is this:

```
cmd                # contains the application's command line interface
cmd/blog/main.go   # the main entrypoint for the `blog` command
cmd/blog/serve.go  # the `blog serve` subcommand which starts the http server, by default on port 8080
cmd/blog/build.go  # the `blog build` subcommand which builds a single binary including all assets that can be shipped to the target host
cmd/blog/deploy.go # the `blog deploy` subcommand which copies the binary via SSH to the target host

core                 # package core takes care of all infrastructure concerns
core/commands.go     # a registry for commands and their associated handlers
core/queries.go      # a registry for queries and their associated handlers
core/form.go         # a flexible data structure for capturing data with error states
core/log.go          # a persistent event log, backed by sqlite3
core/view.go         # provides shared components
core/view_layout.go  # Layout gomponent + views
core/page_index.go   # The body gomponent for the index page
```

### Adding a new feature

```sh
petrock feature posts
```

This generates a new Go package called "posts", which contains all of the functionality related to authoring and editing posts:

```
posts/            # the package for this feature
posts/register.go # the entrypoint for the module which registers it with the core
posts/messages.go # structs for the messages accepted by this feature, and data returned by queries
posts/execute.go  # functions for accepting messages that change the state of posts
posts/query.go    # functions for accepting messages that returns parts of the state of posts
posts/state.go    # application state that needs to be kept
posts/jobs.go     # long-running processes that
posts/view.go     # components for rendering
posts/assets.go   # a file that builds an in-memory FS using go:embed for the assets directory
posts/routes.go   # defines feature-specific HTTP routes
posts/http.go     # contains feature-specific HTTP handlers
posts/assets.go   # a file that builds an in-memory FS using go:embed for the assets directory
posts/assets/     # a directory containing binary assets that should get included in the final binary
```

# Inside Petrock generated code

Petrock is about generating the simplest code that could possibly work to achieve the goal at hand.

The high-level overview of a Petrock application is this:

1.  **Startup:** The application initializes core components (database connection, message log, command/query registries, application state). It then replays all messages from the persistent log (`messages` table in SQLite) to rebuild the in-memory application state. Finally, it registers all features defined in the project.
2.  **Feature Registration:** When `petrock feature <name>` is run, it automatically adds an import and a registration call (e.g., `posts.RegisterFeature(...)`) to `cmd/<project>/features.go`. During startup, `RegisterAllFeatures` calls each feature's `RegisterFeature` function. This function registers the feature's command handlers, query handlers, and message types (for decoding) with the core registries and message log.
3.  **API Interaction:** The application exposes a core API for interacting with commands and queries:
    *   `GET /`: Displays an HTML index page listing available commands and queries.
    *   `GET /commands`: Returns a JSON list of registered command names (e.g., `["posts/create", "posts/update"]`).
    *   `POST /commands`: Executes a command. Expects a JSON body like `{"type": "feature/create", "payload": {...}}`. The handler decodes the payload into the appropriate command struct, dispatches it via the `CommandRegistry` using its `CommandName()`, logs the command, and applies it to the state. Returns `200 OK` or `202 Accepted` on success, or `400/500` on error.
    *   `GET /queries`: Returns a JSON list of registered query names (e.g., `["posts/get", "posts/list"]`).
    *   `GET /queries/{feature}/{query-name}`: Executes a query. The path contains the full kebab-case name (e.g., `/queries/posts/list`). Query parameters (e.g., `?ID=123&page=1`) are automatically parsed and mapped to the fields of the query struct. The handler dispatches the query via the `QueryRegistry` using its `QueryName()` and returns the JSON result on success (`200 OK`), or `400/404/500` on error.
4.  **Feature-Specific HTTP Routes:** In addition to the core API, features can define their own HTTP routes and handlers:
    *   Routes are defined in `<feature>/routes.go` using the standard `net/http.ServeMux`.
    *   Handlers are implemented in `<feature>/http.go`, typically as methods on a `FeatureServer` struct holding dependencies (Executor, Querier, State, Log, etc.).
    *   These routes are registered in `cmd/<project>/serve.go` *after* the core routes. This means features can add new endpoints (e.g., `GET /posts/{id}`) or even override core endpoints (like `/`) if needed. Conventionally, feature routes should be prefixed (e.g., `/posts/...`) to avoid accidental overrides.
5.  **Command Handling:** Commands are processed through a centralized Executor pattern:
     * All commands are executed via the `core.Executor.Execute(cmd)` method, which provides a standardized flow.
     * Each command implements a `Validate()` method which is automatically called by the Executor before logging.
     * If validation passes, the command is appended to the message log (`core.MessageLog`).
     * The Executor then dispatches the command to its registered handler via the CommandRegistry.
     * Handlers focus exclusively on domain-specific logic rather than repeating the validation/logging pattern.
     * State changes are applied to the relevant feature's in-memory state (`feature.State.Apply`).
     * Application code and feature-specific HTTP handlers use the core Executor to process commands rather than directly accessing the message log or command registry.
     * This pattern centralizes validation, logging, and dispatching logic, eliminating duplication across features.
6.  **Query Handling:** When a query is dispatched (either via `GET /queries/...` or a feature-specific handler), the registered handler reads directly from the feature's in-memory state (`feature.State`) to produce the result.
7.  **State Management:** Each feature manages its own slice of the application state within its `state.go` file. The `Apply` method is crucial for both rebuilding state from the log at startup and applying live updates after a command is logged. Handlers in `<feature>/http.go` access this state, usually via the `Querier` or direct `State` access.

All commands are serialized with a configurable encoder (JSON by default) and persisted in a SQLite database in a single table called `messages`.
