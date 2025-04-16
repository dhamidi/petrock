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
posts/commands.go # structs for commands that change the state of this feature
posts/queries.go  # structs for queries and their result types
posts/execute.go  # functions for accepting messages that change the state of posts
posts/query.go    # functions for accepting messages that returns parts of the state of posts
posts/state.go    # application state that needs to be kept
posts/worker.go   # long-running background processes that react to events in the message log
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

1.  **Startup:** The application initializes core components through the central `App` struct in `core/app.go` which manages database connection, message log, command/query registries, and application state. `App.RegisterFeatures()` registers all features defined in the project first, which is crucial for message deserialization. Then `App.ReplayLog()` replays messages from the persistent log (`messages` table in SQLite) by iterating through them with `messageLog.After(ctx, 0)` to rebuild the in-memory application state. The `serve.go` file handles all HTTP concerns: creating the App instance, setting up HTTP routes, and managing server lifecycle.
2.  **Feature Registration:** When `petrock feature <name>` is run, it automatically adds an import and a registration call (e.g., `posts.RegisterFeature(...)`) to `cmd/<project>/features.go`. During startup, `RegisterAllFeatures` calls each feature's `RegisterFeature` function. This function registers the feature's command handlers, query handlers, and message types (for decoding) with the core registries and message log.
3.  **API Interaction:** The application exposes a core API for interacting with commands and queries:
    *   `GET /`: Displays an HTML index page listing available commands and queries.
    *   `GET /commands`: Returns a JSON list of registered command names (e.g., `["posts/create", "posts/update"]`).
    *   `POST /commands`: Executes a command. Expects JSON like `{"type": "feature/create", "payload": {...}}`. The core handler decodes this into the appropriate command struct and passes it to the `core.Executor.Execute`. The Executor retrieves the feature's executor instance, calls its `ValidateCommand` method (which in turn calls the command's `Validate(state)` method if implemented), logs the command if valid, and then calls the registered feature-specific *state update handler*. Returns `200 OK`/`202 Accepted` on success, `400` on validation/decoding errors, `500` on logging errors. (State update errors cause panic).
    *   `GET /queries`: Returns JSON list of registered query names.
    *   `GET /queries/{feature}/{query-name}`: Executes a query. Path gives the name (e.g., `/queries/posts/list`). Query params (e.g., `?ID=123`) map to query struct fields. The core handler decodes, dispatches via `QueryRegistry` to the feature's query handler, and returns JSON result (`200 OK`) or `400/404/500` error.
4.  **Feature-Specific HTTP Routes:** Features define routes in `<feature>/routes.go` and handlers in `<feature>/http.go`.
    *   Handlers are methods on a `FeatureServer` struct holding dependencies like the central `core.Executor`, the feature's `Querier`, and `State`.
    *   Routes are registered *after* core routes, allowing overrides. Conventionally prefixed (e.g., `/posts/...`).
    *   Feature handlers needing to perform writes **must** go through the central `core.Executor.Execute(ctx, cmd)` method to ensure validation (via the command's `Validate` method if present), logging, and consistent state updates.
    *   Feature handlers performing reads use the feature's `Querier`.
5.  **Command Handling (via `core.Executor`):**
    *   A command (`core.Command`) is constructed (e.g., from an HTTP request).
    *   It's passed to `core.Executor.Execute(ctx, cmd)`.
    *   The Executor looks up the feature's executor instance (`core.FeatureExecutor`) and state update handler (`core.CommandHandler`) in `core.CommandRegistry` using `cmd.CommandName()`. Not found -> return error.
    *   The Executor calls `featureExecutor.ValidateCommand(ctx, cmd)`. This checks if `cmd` implements `Validator` and calls `cmd.Validate(state)` if so. Fails -> return validation error.
    *   The Executor appends the command to `core.MessageLog`. Fails -> return logging error.
    *   The Executor calls the state update handler (defined in `<feature>/execute.go`). Fails -> **panic**.
6.  **Query Handling:** Queries (`core.Query`) are dispatched via the `core.QueryRegistry` to the appropriate handler defined in `<feature>/query.go`. These handlers read directly from the feature's in-memory state (`<feature>/state.go`).
7.  **State Management:** Each feature manages its state in `<feature>/state.go`. State is rebuilt at startup by replaying the `core.MessageLog` using the iterator pattern: `for msg := range messageLog.After(ctx, 0)`. For each message, the corresponding *state update handler* (retrieved from `core.CommandRegistry`) is executed with both the decoded payload and a pointer to the message metadata (`handler(ctx, msg.DecodedPayload, &msg.Message)`). Live updates also happen via these same state update handlers, called by the `core.Executor` after successful validation and logging, but without message metadata (`handler(ctx, cmd, nil)`).
8.  **Workers:** Workers, defined in `<feature>/worker.go`, are long-running background processes that react to events in the message log. Each worker implements the `core.Worker` interface with `Start()`, `Stop()`, and `Work()` methods. Workers maintain their own internal state by tracking events from the message log and perform operations that span multiple events, often interacting with external systems. During application startup, workers are registered with the central `App` instance which manages their lifecycle: starting them in separate goroutines, scheduling their `Work()` method periodically (by default every 1-2 seconds with jitter), and stopping them during shutdown. Workers typically dispatch commands through the central `core.Executor` when they need to update application state.

All commands are serialized with a configurable encoder (JSON by default) and persisted in a SQLite database in a single table called `messages`.
