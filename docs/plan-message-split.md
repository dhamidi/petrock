# Plan for Splitting `messages.go` into `queries.go` and `commands.go`

## Overview

Currently, all message types (commands and queries) are defined in a single `messages.go` file. This plan outlines how to split these into separate files for improved code organization and readability.

## Documentation Updates Required

### 1. Create New Documentation Files

- Create `docs/feature/commands.go.md` - Move command structures from `messages.go.md`
- Create `docs/feature/queries.go.md` - Move query structures from `messages.go.md`

### 2. Update Existing Documentation Files

#### docs/high-level.md

Line 63 currently reads:
```
posts/messages.go # structs for the messages accepted by this feature, and data returned by queries
```

Update to:
```
posts/commands.go # structs for commands that change the state of this feature
posts/queries.go  # structs for queries and their result types
```

#### docs/feature/messages.go.md

The current content should be split with:
- Command-related content moved to `commands.go.md`
- Query-related content moved to `queries.go.md`
- File will be deprecated after migration

#### docs/feature/execute.go.md

Update any references from `messages.go` to `commands.go` for command validation interfaces.

#### docs/feature/query.go.md

Update any references from `messages.go` to `queries.go` for query types.

#### Other files that may need updates

- `docs/feature/register.go.md` - Update references to type registration
- `README.md` - Update file structure reference

## Implementation Steps

1. Create new template files in `internal/skeleton/feature_template/`:
   - `commands.go` - Move command structures from messages.go
   - `queries.go` - Move query structures from messages.go

2. Update `register.go` in the template to import from the new files

3. Update references in other template files (execute.go, query.go, etc.)

4. Remove `messages.go` from the template directory after confirming all content has been moved

5. Test the template generation with a sample feature

6. Update documentation as outlined above

## Verification

To verify this change works correctly:

1. Generate a new feature using the updated template
2. Ensure all imports are working correctly
3. Verify the application builds and runs correctly
4. Verify command and query execution still works as expected