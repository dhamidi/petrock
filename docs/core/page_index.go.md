# Plan for core/page_index.go

This file defines the Gomponent for rendering the body content of the main index page (`/`).

## Types

- None specific to this file.

## Functions

- `IndexPage() gomponents.Node`: Returns a `gomponents.Node` representing the main content block for the application's home page. This might include welcome text, links to features, etc. It's intended to be passed as the `body` argument to `core.Layout`.
