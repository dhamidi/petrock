# Handlers

The `handlers` directory contains HTTP handlers and middleware for processing web requests related to the feature.

## Structure

- `base.go` - Common handler types and utilities
- `commands.go` - Handlers for processing commands
- `core.go` - Core handler functionality
- `middleware.go` - Common middleware functions
- `create_item.go` - Handlers for item creation (API)
- `create_form.go` - Form handlers for item creation (UI)
- `read_item.go` - Handlers for retrieving single items
- `read_list.go` - Handlers for list views
- `update_item.go` - Handlers for item updates (API)
- `update_form.go` - Form handlers for item updates (UI)
- `delete_item.go` - Handlers for item deletion (API)
- `delete_form.go` - Confirmation forms for deletion (UI)
- `views.go` - View-related handler functionality

## Responsibilities

- Process incoming HTTP requests
- Parse and validate input data using the tag-based validation system (see [Form Validation Guide](../../form-validation-guide.md))
- Execute commands or queries based on request
- Render appropriate responses
- Apply middleware for cross-cutting concerns

## Usage

Handlers are registered with the application's router in the `routes` package, which maps URLs to specific handler functions.