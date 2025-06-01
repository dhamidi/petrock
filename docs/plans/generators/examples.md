# Component Generator Usage Examples

This document provides comprehensive examples of using the new component generators in petrock projects.

## Command Generator Examples

### Basic Command Generation
```bash
# Generate a create command for posts feature
petrock new command posts/create

# Generate a update command with validation
petrock new command posts/update

# Generate a delete command 
petrock new command posts/delete
```

### Generated File Structure
```
posts/
├── commands/
│   ├── base.go           # Shared command infrastructure
│   ├── create.go         # CreateCommand struct and handler
│   └── register.go       # Command registration logic
```

### Command Template Example
```go
// Generated: posts/commands/create.go
package commands

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "strings"
    "time"

    "github.com/myproject/core"
    "github.com/myproject/posts/state"
)

type CreateCommand struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedBy   string    `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
}

func (c *CreateCommand) CommandName() string {
    return "posts/create"
}

func (c *CreateCommand) Validate(state *state.State) error {
    if strings.TrimSpace(c.Name) == "" {
        return errors.New("item name cannot be empty")
    }
    return nil
}
```

## Query Generator Examples

### Basic Query Generation
```bash
# Generate a get query for posts feature
petrock new query posts/get

# Generate a list query
petrock new query posts/list

# Generate a search query
petrock new query posts/search
```

### Generated File Structure
```
posts/
├── queries/
│   ├── base.go           # Shared query infrastructure  
│   ├── get.go            # GetQuery and GetQueryResult
│   └── list.go           # ListQuery and ListQueryResult
```

### Query Template Example
```go
// Generated: posts/queries/get.go
package queries

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/myproject/core"
)

type GetQuery struct {
    ID string `json:"id"`
}

func (q GetQuery) QueryName() string {
    return "posts/get"
}

type GetQueryResult struct {
    Item ItemResult `json:"item"`
}

func (q *Querier) HandleGet(ctx context.Context, query core.Query) (core.QueryResult, error) {
    getQuery, ok := query.(GetQuery)
    if !ok {
        return nil, fmt.Errorf("invalid query type for HandleGet: expected GetQuery, got %T", query)
    }

    item, found := q.state.GetItem(getQuery.ID)
    if !found {
        return nil, fmt.Errorf("item with ID %s not found", getQuery.ID)
    }

    // Map to result...
    return &GetQueryResult{Item: itemResult}, nil
}
```

## Worker Generator Examples

### Basic Worker Generation
```bash
# Generate a summary worker for posts feature
petrock new worker posts/summary

# Generate an indexing worker
petrock new worker posts/indexer

# Generate a notification worker
petrock new worker posts/notifier
```

### Generated File Structure
```
posts/
├── workers/
│   ├── main.go           # Worker registration and setup
│   ├── summary_worker.go # SummaryWorker implementation
│   └── types.go          # Worker-specific types
```

### Worker Template Example
```go
// Generated: posts/workers/summary_worker.go
package workers

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    "github.com/myproject/core"
)

func handleCreateCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState, pctx *core.ProcessingContext) error {
    createCmd, ok := cmd.(*CreateCommand)
    if !ok {
        return fmt.Errorf("unexpected command type: %T", cmd)
    }

    // Skip side effects during replay
    if pctx.IsReplay {
        return nil
    }

    // Process the command...
    return nil
}
```

## Collision Detection Examples

### Detecting Existing Components
```bash
# Try to generate existing command - should fail
$ petrock new command posts/create
Error: Command 'posts/create' already exists in this project.
Use --force to overwrite or choose a different name.

# Check what exists
$ go run . self inspect --json | jq '.commands[].name'
"posts/create"
"posts/update"
"posts/delete"
```

### Force Overwrite (Future Feature)
```bash
# Overwrite existing component (planned feature)
petrock new command posts/create --force
Warning: Overwriting existing command 'posts/create'
```

## Integration Examples

### Adding to Existing Feature
```bash
# Add new query to existing posts feature
cd myproject/
petrock new query posts/search

# Generated files integrate with existing posts/ structure
ls posts/
commands/  handlers/  queries/  routes/  state/  ui/  workers/

# New query appears in inspection
go run . self inspect --json | jq '.queries[] | select(.name | contains("search"))'
{
  "name": "posts/search",
  "type": "github.com/myproject/posts/queries.SearchQuery"
}
```

### Cross-Feature Dependencies
```bash
# Generate worker that processes multiple features
petrock new worker analytics/aggregator

# Worker can import from multiple features
cat analytics/workers/aggregator_worker.go
import (
    "github.com/myproject/posts/state"
    "github.com/myproject/users/state" 
    "github.com/myproject/analytics/state"
)
```

## Advanced Usage

### Custom Entity Names
```bash
# Use descriptive entity names
petrock new command posts/publish
petrock new query posts/by-author  
petrock new worker posts/auto-tagger
```

### Component Validation
```bash
# Generators validate feature/entity format
$ petrock new command posts-create
Error: Invalid format. Use 'feature/entity' (e.g., 'posts/create')

$ petrock new command posts/
Error: Entity name cannot be empty
```

### Build Integration
```bash
# Generated components are validated during build
./build.sh
✓ Template compilation successful
✓ Generated project compilation successful  
✓ All tests pass
```

## Template Customization Examples

### Extending Base Templates
```go
// Custom command extending generated base
package commands

import "github.com/myproject/posts/commands"

type AdvancedCreateCommand struct {
    *commands.CreateCommand
    Tags     []string `json:"tags"`
    Category string   `json:"category"`
}

func (c *AdvancedCreateCommand) CommandName() string {
    return "posts/advanced-create"
}
```

### Feature-Specific Variations
```bash
# Different features can have specialized components
petrock new command users/activate    # User-specific command
petrock new query   orders/by-date    # Order-specific query  
petrock new worker  images/processor  # Image-specific worker
```

These examples demonstrate the flexibility and power of the component generator system while maintaining the simplicity and consistency that makes petrock effective for rapid application development.
