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
posts/assets/     # a directory containing binary assets that should get included in the final binary
```

# Inside Petrock generated code

Petrock is about generating the simplest code that could possibly work to achieve the goal at hand.

The high-level overview of a Petrock application is this:

Petrock accepts a command, usually via HTTP,
and then hands this command over to the domain logic layer,
which either accepts or rejects the command.

If the command is accepted, it gets persisted in the command log, and the application state is updated based on the now-accepted command.

At application startup, Petrock iterates over all commands to build the in-memory application state before it starts accepting requests.

All commands are serialized with a configurable encoder, JSON by default, and then persisted in a sqlite3 database in a single table called messages.
