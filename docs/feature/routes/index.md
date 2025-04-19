# Routes

The `routes` directory contains route definitions and registration for the feature's HTTP endpoints.

## Structure

- `main.go` - Central route registration
- `api.go` - API routes (REST endpoints)
- `web.go` - Web UI routes (HTML pages)

## Responsibilities

- Define URL patterns for the feature
- Map URLs to specific handlers
- Group related routes together
- Configure route-specific middleware
- Handle route parameters

## Route Registration Pattern

Routes follow a registration pattern where:

1. Routes are defined in their respective files (api.go, web.go)
2. The registration functions are called during feature initialization
3. Routes are organized by their purpose (API vs. web UI)

## Usage

Routes connect incoming HTTP requests to the appropriate handlers, allowing the feature to respond to different URLs and HTTP methods.