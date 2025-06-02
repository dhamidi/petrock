# Component Generator Usage Examples

This document provides comprehensive examples of using the petrock component generators to create commands, queries, and workers.

## Overview

The petrock component generators allow you to quickly create individual components without generating entire features. This is useful for:

- Adding new operations to existing features
- Creating specialized components for specific use cases  
- Iterative development and testing
- Learning the petrock architecture patterns

## Basic Usage

All component generators follow the same pattern:

```bash
petrock new <component-type> <feature>/<name-of-thing>
```

Where:
- `<component-type>` is one of: `command`, `query`, `worker`
- `<feature>` is the feature name (e.g., `posts`, `users`, `orders`)
- `<name-of-thing>` is the specific operation name (e.g., `create`, `get`, `summary`)

## Command Generation Examples

Commands handle business logic and state changes in your application.

### Basic CRUD Operations

```bash
# Generate basic CRUD commands for a posts feature
petrock new command posts/create
petrock new command posts/update  
petrock new command posts/delete

# Generate user management commands
petrock new command users/register
petrock new command users/login
petrock new command users/logout
petrock new command users/activate

# Generate order processing commands
petrock new command orders/place
petrock new command orders/cancel
petrock new command orders/fulfill
petrock new command orders/refund
```

### Specialized Commands

```bash
# Content management commands
petrock new command posts/publish
petrock new command posts/archive
petrock new command posts/feature

# User administration commands  
petrock new command users/ban
petrock new command users/promote
petrock new command users/reset_password

# Business logic commands
petrock new command inventory/restock
petrock new command payments/process
petrock new command notifications/send
```

### Generated Command Files

When you run `petrock new command posts/create`, it generates:

```
posts/
└── commands/
    ├── base.go          # Base command interfaces and types
    ├── register.go      # Command registration logic
    └── create.go        # CreateCommand implementation
```

Example generated `create.go`:

```go
package commands

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "strings"
    "time"

    "github.com/yourmodule/yourproject/core"
    "github.com/yourmodule/yourproject/posts/state"
)

// CreateCommand holds data needed to create a new entity.
type CreateCommand struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedBy   string    `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
}

// CommandName returns the unique kebab-case name for this command type.
func (c *CreateCommand) CommandName() string {
    return "posts/create"
}

// Validate implements the Validator interface for CreateCommand.
func (c *CreateCommand) Validate(state *state.State) error {
    // Validation logic here
    return nil
}

// HandleCreate applies state changes for CreateCommand.
func (e *Executor) HandleCreate(ctx context.Context, command core.Command, msg *core.Message, pctx *core.ProcessingContext) error {
    // Implementation here
    return nil
}
```

## Query Generation Examples

Queries handle data retrieval and read operations in your application.

### Basic Query Operations

```bash
# Generate basic query operations for posts
petrock new query posts/get     # Get single post by ID
petrock new query posts/list    # List posts with pagination

# Generate user query operations  
petrock new query users/get     # Get user by ID
petrock new query users/search  # Search users by criteria

# Generate order query operations
petrock new query orders/get    # Get order details
petrock new query orders/list   # List orders for user
petrock new query orders/search # Search orders by criteria
```

### Specialized Queries

```bash
# Analytics and reporting queries
petrock new query analytics/count
petrock new query analytics/summary
petrock new query reports/daily
petrock new query reports/monthly

# Complex business queries
petrock new query products/featured
petrock new query products/recommendations
petrock new query users/active
petrock new query orders/pending
```

### Generated Query Files

When you run `petrock new query posts/get`, it generates:

```
posts/
└── queries/
    ├── base.go    # Base query interfaces and types
    └── get.go     # GetQuery implementation
```

Example generated `get.go`:

```go
package queries

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/yourmodule/yourproject/core"
)

// GetQuery holds data needed to retrieve a single entity.
type GetQuery struct {
    ID string // ID of the entity to retrieve
}

// QueryName returns the unique kebab-case name for this query type.
func (q GetQuery) QueryName() string {
    return "posts/get"
}

// GetQueryResult wraps an ItemResult as a specific result type for GetQuery
type GetQueryResult struct {
    Item ItemResult `json:"item"`
}

// HandleGet processes the GetQuery.
func (e *Executor) HandleGet(ctx context.Context, query core.Query) (core.QueryResult, error) {
    // Implementation here
    return &GetQueryResult{}, nil
}
```

## Worker Generation Examples

Workers handle background processing, side effects, and asynchronous operations.

### Background Processing Workers

```bash
# Generate content processing workers
petrock new worker posts/summary       # Summarize post content
petrock new worker posts/analysis      # Analyze post metrics
petrock new worker posts/cleanup       # Clean up old posts

# Generate user-related workers
petrock new worker users/notification  # Send user notifications
petrock new worker users/backup        # Backup user data
petrock new worker users/sync          # Sync with external systems
```

### Integration Workers

```bash
# External service integration workers
petrock new worker payments/process    # Process payment transactions
petrock new worker email/send          # Send email notifications
petrock new worker sms/send            # Send SMS notifications
petrock new worker storage/backup      # Backup to cloud storage

# Data processing workers
petrock new worker analytics/process   # Process analytics data
petrock new worker reports/generate    # Generate reports
petrock new worker images/resize       # Resize uploaded images
petrock new worker search/index        # Index content for search
```

### Generated Worker Files

When you run `petrock new worker posts/summary`, it generates:

```
posts/
└── workers/
    ├── main.go           # Worker registry and startup logic
    ├── types.go          # Worker state and type definitions
    └── summary_worker.go # SummaryWorker implementation
```

Example generated `summary_worker.go`:

```go
package workers

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    "github.com/yourmodule/yourproject/core"
)

// handleCreateCommand processes new item creation commands
func handleCreateCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState, pctx *core.ProcessingContext) error {
    // Worker implementation here
    return nil
}

// processPendingSummaries calls external API for pending summaries
func processPendingSummaries(ctx context.Context, workerState *WorkerState) error {
    // Background processing logic here
    return nil
}
```

## Common Patterns and Best Practices

### Naming Conventions

**Commands:**
- Use action verbs: `create`, `update`, `delete`, `publish`, `archive`
- Use business terms: `place` (order), `fulfill` (order), `ban` (user)
- Keep names concise but descriptive

**Queries:**
- Use retrieval verbs: `get`, `list`, `search`, `find`, `count`
- Use descriptive qualifiers: `active`, `pending`, `featured`
- Use time-based qualifiers: `recent`, `daily`, `monthly`

**Workers:**
- Use processing verbs: `process`, `analyze`, `generate`, `sync`
- Use domain-specific terms: `summary`, `notification`, `backup`
- Use integration terms: `import`, `export`, `upload`

### Entity Name Guidelines

Entity names should:
- Start with a letter (a-z, A-Z)
- Contain only letters, numbers, and underscores
- Use snake_case for multi-word names: `reset_password`, `send_notification`
- Be descriptive and unambiguous

Valid examples:
```bash
petrock new command users/reset_password
petrock new query products/search_featured
petrock new worker analytics/generate_report
```

Invalid examples:
```bash
petrock new command users/123invalid    # Cannot start with number
petrock new query posts/get-data        # Hyphens not allowed
petrock new worker orders/process data  # Spaces not allowed
```

### Feature Organization

Organize components by business domains:

```
# User management feature
petrock new command users/register
petrock new command users/login
petrock new query users/get
petrock new query users/search
petrock new worker users/notification

# Blog/content feature  
petrock new command posts/create
petrock new command posts/publish
petrock new query posts/get
petrock new query posts/list
petrock new worker posts/summary

# E-commerce feature
petrock new command orders/place
petrock new command orders/cancel
petrock new query orders/get
petrock new query orders/search
petrock new worker orders/process
```

## Integration with Existing Features

### Adding Components to Existing Features

If you already have a `posts` feature, you can add new components:

```bash
# Add new command to existing posts feature
petrock new command posts/feature    # Feature a post
petrock new command posts/moderate   # Moderate post content

# Add new queries to existing posts feature
petrock new query posts/trending     # Get trending posts
petrock new query posts/archived     # Get archived posts

# Add new workers to existing posts feature
petrock new worker posts/analytics   # Analyze post performance
petrock new worker posts/backup      # Backup post data
```

### Collision Detection

The generator automatically detects existing components and prevents overwrites:

```bash
$ petrock new command posts/create
Error: command posts/create already exists
```

This prevents accidental overwrites of existing code.

## Advanced Usage

### Generating Multiple Related Components

Create a complete CRUD interface for a new entity:

```bash
# Create all CRUD commands
petrock new command products/create
petrock new command products/update
petrock new command products/delete

# Create corresponding queries
petrock new query products/get
petrock new query products/list
petrock new query products/search

# Create supporting workers
petrock new worker products/index      # Search indexing
petrock new worker products/analyze    # Performance analytics
```

### Building Microservice Components

Create components for microservice communication:

```bash
# Payment service integration
petrock new command payments/process
petrock new command payments/refund
petrock new query payments/status
petrock new worker payments/webhook

# Inventory service integration
petrock new command inventory/reserve
petrock new command inventory/release
petrock new query inventory/check
petrock new worker inventory/sync
```

### Event-Driven Architecture

Create components for event handling:

```bash
# Event producers (commands)
petrock new command events/publish
petrock new command events/schedule

# Event consumers (workers)
petrock new worker events/process
petrock new worker events/replay
petrock new worker events/archive

# Event queries
petrock new query events/get
petrock new query events/search
```

## Troubleshooting

### Common Issues

**Issue: "git workspace is not clean"**
```bash
$ petrock new command posts/create
Error: git workspace is not clean
```

**Solution:** Commit or stash your changes before running generators:
```bash
git add .
git commit -m "Work in progress"
petrock new command posts/create
```

**Issue: "failed to detect module path"**
```bash
$ petrock new command posts/create  
Error: failed to detect module path: failed to read go.mod
```

**Solution:** Ensure you're in a Go module directory with a `go.mod` file:
```bash
go mod init github.com/youruser/yourproject
petrock new command posts/create
```

**Issue: "invalid entity name"**
```bash
$ petrock new command posts/123invalid
Error: invalid entity name "123invalid": must contain only letters, numbers, and underscores
```

**Solution:** Use valid entity names that start with a letter:
```bash
petrock new command posts/create_post
```

### Best Practices for Success

1. **Plan your component structure** before generating
2. **Use consistent naming conventions** across your project
3. **Generate related components together** (e.g., command + query + worker)
4. **Test generated components** to ensure they integrate properly
5. **Customize generated code** to fit your specific business logic
6. **Document your component architecture** for team members

## Next Steps

After generating components:

1. **Implement business logic** in the generated files
2. **Add validation rules** to commands
3. **Implement query logic** for data retrieval
4. **Configure worker processing** for background tasks
5. **Write tests** for your components
6. **Update documentation** with your specific use cases

For more information, see:
- [Petrock Architecture Guide](../architecture.md)
- [CQRS Patterns](../patterns/cqrs.md)
- [Testing Guide](../testing.md)
