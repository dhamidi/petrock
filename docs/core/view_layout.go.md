# Plan for core/view_layout.go

This file defines the main HTML page layout using Gomponents (`maragu.dev/gomponents`).

## Types

- None specific to this file.

## Functions

- `Layout(title string, body ...gomponents.Node) gomponents.Node`: Renders the full HTML document structure (`<!DOCTYPE html>`, `<html>`, `<head>`, `<body>`) using `maragu.dev/gomponents`.
    - The `<head>` section includes the page `title`, links to CSS (Tailwind), potentially JS (Stimulus, Hotwire), meta tags.
    - The `<body>` section includes the content passed in the `body` parameter(s). It might also include common elements like headers, footers, or navigation bars. It should include necessary attributes for Hotwire/Stimulus if used (e.g. `data-controller`).
