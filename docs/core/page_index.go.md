# Plan for core/page_index.go

This file defines the Gomponent for rendering the body content of the main index page (`/`). It uses `maragu.dev/gomponents`.

## Types

- None specific to this file.

## Functions

- `IndexPage(commandNames, queryNames []string) g.Node`: Returns a `maragu.dev/gomponents.Node` representing the main content block for the application's home page (`/`). By default, it lists registered command and query names. It's intended to be passed as the `body` argument to `core.Layout`. *Note: The handler for the `/` route, which calls this function, can be overridden by a feature if it registers its own handler for `/` in its `routes.go` file.*
