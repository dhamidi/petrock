# Technical Implementation Plan for Feature Template Structure Reorganization

## Overview

Currently, the feature template uses a flat directory structure with all major files in the root directory. This results in large, monolithic files as features grow in complexity, making the codebase harder to navigate and maintain.

After this change, the feature template will use a hierarchical directory structure that organizes code by responsibility. This will distribute functionality across smaller, more focused files, improving maintainability, discoverability, and organization as features scale.

## Workflow

1. Pick the first task that is not marked DONE
2. For each subtask
3. Implement the subtask
4. Commit all changes including the task identifier in the commit message
5. Mark the subtask as DONE
6. Evaluate the definition of done for the task
7. If it is met, mark the task as DONE
8. Otherwise ask for feedback from the user.

## Detailed Task Breakdown

### T1: Create Directory Structure - DONE

**T1.1:** Create the base directories for the new structure - DONE

- Create directories in internal/skeleton/feature_template/: commands/, handlers/, queries/, state/, ui/, routes/, workers/
- Add .keep files in each directory to ensure they're tracked by git

**T1.2:** Create subdirectory structure - DONE

- Create ui/components/, ui/layouts/, ui/pages/ directories
- No subdirectories for handlers - use descriptive file names instead

**Definition of Done for T1:**

- All directories and subdirectories exist in the feature template
- Each directory has a .keep file to ensure it's tracked by git
- Directory structure matches the specification in docs/new-feature-spec.md

### T2: Migrate Command-Related Code - DONE

**T2.1:** Create base files in commands/ directory - DONE

- Create commands/base.go with common interfaces and types
- Create commands/create.go, commands/update.go, commands/delete.go files

**T2.2:** Extract command definitions from commands.go - DONE

- Move CreateCommand to commands/create.go
- Move UpdateCommand to commands/update.go
- Move DeleteCommand to commands/delete.go
- Move RequestSummaryGenerationCommand to commands/request_summary.go
- Move SetGeneratedSummaryCommand to commands/set_summary.go
- Move FailSummaryGenerationCommand to commands/fail_summary.go

**T2.3:** Extract command handlers from execute.go - DONE

- Move handler functions to their respective command files
- Update imports and references

**Definition of Done for T2:**

- All command definitions are moved to appropriate files in commands/ directory
- All command handler functions are moved to their respective command files
- Original commands.go and execute.go are removed
- Code compiles successfully

### T3: Migrate Query-Related Code - DONE

**T3.1:** Create base files in queries/ directory - DONE

- Create queries/base.go with common interfaces and types
- Create queries/get.go and queries/list.go files

**T3.2:** Extract query definitions and result types from queries.go - DONE

- Move list query and list result types to queries/list.go
- Move get query and item result type to queries/get.go
- Include result type definitions in the same file as their corresponding queries

**T3.3:** Extract query handlers from query.go - DONE

- Move handler functions to their respective query files
- Update imports and references

**Definition of Done for T3:**

- All query definitions are moved to appropriate files in queries/ directory
- All query handler functions are moved to their respective query files
- Original queries.go and query.go are removed
- Code compiles successfully

### T4: Migrate UI Components - DONE

**T4.1:** Create base files in ui/ directory - DONE

- Create files in ui/components/, ui/layouts/, ui/pages/ directories

**T4.2:** Extract UI components from view.go - DONE

- Move table components to ui/components/tables.go
- Move form components to ui/components/forms.go
- Move layout components to ui/layouts/
- Move page components to ui/pages/

**Definition of Done for T4:**

- All UI components are moved to appropriate files in ui/ directory
- Original view.go is removed
- Code compiles successfully

### T5: Migrate HTTP Handlers - DONE

**T5.1:** Create base files in handlers/ directory - DONE

- Create handlers/base.go for common utilities and types
- Create handlers/middleware.go for common middleware
- Create individual handler files with descriptive names (e.g., create_item.go, read_list.go)

**T5.2:** Extract HTTP handlers from http.go - DONE

- Move create handlers to create_item.go and create_form.go
- Move read/list handlers to read_item.go and read_list.go
- Move update handlers to update_item.go and update_form.go
- Move delete handlers to delete_item.go and delete_form.go
- Move middleware functions to middleware.go

**T5.3:** Update imports and references - DONE

- Ensure all imports and references are updated to reflect new file locations

**Definition of Done for T5:**

- All HTTP handlers are moved to appropriate files in handlers/ directory
- Original http.go is removed
- Code compiles successfully

### T6: Migrate State Management - DONE

**T6.1:** Create base files in state/ directory - DONE

- Create state/main.go with main state container and interfaces
- Create state/item.go and state/metadata.go files

**T6.2:** Extract state code from state.go - DONE

- Move State struct to state/main.go
- Move item-related functions to state/item.go
- Move metadata-related functions to state/metadata.go

**T6.3:** Update imports and references - DONE

- Ensure all imports and references are updated to reflect new file locations

**Definition of Done for T6:**

- All state management code is moved to appropriate files in state/ directory
- Original state.go is removed
- Code compiles successfully

### T7: Migrate Worker Code - DONE

**T7.1:** Create base files in workers/ directory - DONE

- Create workers/main.go with common worker interfaces and building blocks
- Create workers/summary_worker.go for the complete summary generation worker
- Create workers/types.go for shared worker type definitions

**T7.2:** Extract worker code from worker.go - DONE

- Move common interfaces and building blocks to workers/main.go
- Move complete summary worker implementation to workers/summary_worker.go
- Move shared type definitions to workers/types.go

**T7.3:** Update imports and references - DONE

- Ensure all imports and references are updated to reflect new file locations

**Definition of Done for T7:**

- All worker code is moved to appropriate files in workers/ directory
- Original worker.go is removed
- Code compiles successfully

### T8: Create Main Package File - DONE

**T8.1:** Create main.go in feature_template/ root - DONE

- Create main.go file with imports for all subpackages

**T8.2:** Migrate registration logic from register.go - DONE

- Move feature registration logic to main.go
- Update imports and references

**Definition of Done for T8:**

- main.go contains the core feature initialization and registration logic
- register.go is removed
- Code compiles successfully

### T9: Migrate Routes

**T9.1:** Create base files in routes/ directory

- Create routes/main.go, routes/api.go, routes/web.go files

**T9.2:** Extract route definitions from routes.go

- Move API routes to routes/api.go
- Move web UI routes to routes/web.go
- Move route registration to routes/main.go

**T9.3:** Update imports and references

- Ensure all imports and references are updated to reflect new file locations

**Definition of Done for T9:**

- All route definitions are moved to appropriate files in routes/ directory
- Original routes.go is removed
- Code compiles successfully

### T10: Verification and Testing

**T10.1:** Run build to verify new structure

- Run ./build.sh to verify the new structure compiles
- Fix any compile errors or issues

**T10.2:** Generate a test feature with the new template

- Use petrock to generate a test feature with the new template
- Verify that the generated feature works as expected

**T10.3:** Update documentation

- Update any documentation that references the old structure
- Add documentation for the new structure

**Definition of Done for T10:**

- Build passes successfully
- Test feature generates and works correctly
- Documentation is updated to reflect the new structure

## Implementation Details

When implementing the new structure, use these guidelines:

1. Each file should have a clear, focused responsibility
2. Use consistent package naming across the new structure
3. Use proper imports to reference code in other packages
4. Follow Go conventions for package organization

Example file structure for a command file:

```go
package commands

import (
	// Add necessary imports
)

// CreateCommand holds data needed to create a new entity
type CreateCommand struct {
	// Fields...
}

// CommandName returns the unique name for this command type
func (c *CreateCommand) CommandName() string {
	return "petrock_example_feature_name/create"
}

// HandleCreate handles the create command
func HandleCreate(cmd *CreateCommand, state *State) error {
	// Implementation...
}
```