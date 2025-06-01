# Component Generator File Mappings

This document defines the exact file mappings and extraction rules for each component type.

## Command Component Mapping

### Source Template Files
From `internal/skeleton/petrock_example_feature_name/commands/`:

| Template File | Target File | Required | Description |
|---------------|-------------|----------|-------------|
| `base.go` | `{feature}/commands/base.go` | Yes | Common command interfaces and validation |
| `create.go` | `{feature}/commands/{entity}.go` | Yes | Main command implementation |
| `register.go` | `{feature}/commands/register.go` | Yes | Command registration logic |

### Placeholder Replacements
| Placeholder | Replacement | Example |
|-------------|-------------|---------|
| `petrock_example_feature_name` | `{feature}` | `posts` |
| `CreateCommand` | `{Entity}Command` | `PublishCommand` |
| `HandleCreate` | `Handle{Entity}` | `HandlePublish` |
| `create` (in command name) | `{entity}` | `publish` |

### Dependencies
Commands require:
- Feature state package (`{feature}/state/`)
- Core package import path update

## Query Component Mapping

### Source Template Files  
From `internal/skeleton/petrock_example_feature_name/queries/`:

| Template File | Target File | Required | Description |
|---------------|-------------|----------|-------------|
| `base.go` | `{feature}/queries/base.go` | Yes | Common query interfaces and result types |
| `get.go` | `{feature}/queries/{entity}.go` | Yes | Main query implementation |
| `list.go` | `{feature}/queries/list.go` | Optional | List query (if entity is 'list') |

### Placeholder Replacements
| Placeholder | Replacement | Example |
|-------------|-------------|---------|
| `petrock_example_feature_name` | `{feature}` | `posts` |
| `GetQuery` | `{Entity}Query` | `SearchQuery` |
| `GetQueryResult` | `{Entity}QueryResult` | `SearchQueryResult` |
| `HandleGet` | `Handle{Entity}` | `HandleSearch` |
| `get` (in query name) | `{entity}` | `search` |

### Dependencies
Queries require:
- Feature state package (`{feature}/state/`)
- Core package import path update
- ItemResult type from base.go

## Worker Component Mapping

### Source Template Files
From `internal/skeleton/petrock_example_feature_name/workers/`:

| Template File | Target File | Required | Description |
|---------------|-------------|----------|-------------|
| `main.go` | `{feature}/workers/main.go` | Yes | Worker registration and setup |
| `summary_worker.go` | `{feature}/workers/{entity}_worker.go` | Yes | Main worker implementation |
| `types.go` | `{feature}/workers/types.go` | Yes | Worker-specific types and state |

### Placeholder Replacements
| Placeholder | Replacement | Example |
|-------------|-------------|---------|
| `petrock_example_feature_name` | `{feature}` | `posts` |
| `summary_worker` | `{entity}_worker` | `indexer_worker` |
| `SummaryWorker` | `{Entity}Worker` | `IndexerWorker` |
| `summary` (in function names) | `{entity}` | `indexer` |

### Dependencies
Workers require:
- Feature commands package (`{feature}/commands/`)
- Feature state package (`{feature}/state/`)
- Core package import path update

## Supporting File Generation

### State Package Requirements
If target feature doesn't exist, generators may need to create minimal state files:

```
{feature}/
├── state/
│   ├── main.go          # State struct and basic methods
│   ├── item.go          # Item type definition
│   └── metadata.go      # Metadata helpers
```

### Registration Integration
Generated components must register themselves:

#### Commands
```go
// In {feature}/commands/register.go
func RegisterCommands(registry *core.CommandRegistry) {
    registry.RegisterCommandType(&CreateCommand{})
    registry.RegisterCommandType(&UpdateCommand{})
    // ... other commands
}
```

#### Queries  
```go
// In {feature}/queries/base.go
func RegisterQueries(registry *core.QueryRegistry, querier *Querier) {
    registry.RegisterQuery("posts/get", querier.HandleGet)
    registry.RegisterQuery("posts/list", querier.HandleList)
    // ... other queries
}
```

#### Workers
```go
// In {feature}/workers/main.go  
func RegisterWorkers(app *core.App) error {
    worker := NewSummaryWorker()
    return app.RegisterWorker(worker)
}
```

## Extraction Algorithm

### File Selection
1. **Identify template source**: Map component type to skeleton subdirectory
2. **Filter template files**: Select only files relevant to component type
3. **Apply entity mapping**: Replace generic entity names with specific entity
4. **Validate dependencies**: Ensure required supporting files exist

### Placeholder Processing
1. **Module path replacement**: Update import paths to target project
2. **Feature name replacement**: Replace template feature with actual feature  
3. **Entity name replacement**: Replace template entity with specific entity
4. **Type name replacement**: Update struct/interface names consistently

### Target Structure Creation
1. **Create directories**: Ensure target directory structure exists
2. **Copy and transform**: Apply replacements while copying files
3. **Validate imports**: Ensure all import paths are correct
4. **Check compilation**: Verify generated code compiles

## File Interdependencies

### Command Dependencies
- **base.go**: Provides Validator interface and common command utilities
- **{entity}.go**: Implements specific command logic
- **register.go**: Registers all commands with the application

### Query Dependencies  
- **base.go**: Provides Querier struct and common result types
- **{entity}.go**: Implements specific query handlers
- **Integration**: Queries must be registered in main application

### Worker Dependencies
- **main.go**: Provides worker registration and lifecycle management
- **{entity}_worker.go**: Implements specific worker logic
- **types.go**: Defines worker-specific data structures

## Error Handling

### Missing Dependencies
If required files don't exist:
```
Error: Cannot generate command 'posts/create' - feature 'posts' does not exist.
Run 'petrock feature posts' first to create the base feature structure.
```

### Import Path Resolution
If module path cannot be determined:
```
Error: Cannot determine module path. Ensure you're in a valid Go module directory.
Run 'go mod init <module-path>' if this is a new project.
```

### Compilation Validation
If generated code doesn't compile:
```
Error: Generated component failed compilation check.
This may indicate missing dependencies or template issues.
Run './build.sh' for detailed error information.
```
