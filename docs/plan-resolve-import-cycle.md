# Technical Implementation Plan for Resolving Import Cycle Between State and Commands

## Overview

Currently, we have an import cycle between the state and commands packages:
- The state package imports the commands package to use command types
- The commands package imports the state package to access State for validation

This circular dependency makes the code impossible to compile.

After the change, we'll have a clean separation of concerns:
- Commands package will define all commands and handlers
- State package will provide only an interface/API surface for manipulating state
- Command handlers will use this API to directly manipulate state

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

### T1: Remove state/commands.go and Move State Updates to Command Handlers

**T1.1:** Remove state/commands.go file entirely - DONE

- Delete the whole file state/commands.go
- Confirm deletion with git status

**T1.2:** Update main.go to remove command imports - DONE

- Remove imports for commands package in state/main.go
- Remove all command type references in state/main.go

**T1.3:** Move RegisterTypes function to commands package - DONE

- Create a new function in commands package to register command types
- Remove RegisterTypes from state package

**T1.4:** Remove Apply method from state package

- Delete the Apply method from state/main.go
- Create specific state manipulation functions in state package if needed

**T1.5:** Update command handlers to manipulate state directly

- Modify each command handler to call the appropriate state function
- Replace the state.Apply call with direct state manipulation

**Definition of Done for T1:**

- No import of commands package in any state package file
- No import cycle detected when running build
- All command handlers correctly update state
- Build succeeds without errors

## Implementation Details

Current pattern (in command handlers):
```go
// Apply the change using the state's Apply method
if err := e.state.Apply(cmd, msg); err != nil {
    slog.Error("State Apply failed for CreateCommand", "error", err, "id", cmd.ID)
    return fmt.Errorf("state.Apply failed for CreateCommand: %w", err)
}
```

New pattern (in command handlers):
```go
// Directly manipulate state using state package functions
existingItem, err := e.state.AddItem(&state.Item{
    ID:          cmd.Name,
    Name:        cmd.Name,
    Description: cmd.Description,
    CreatedAt:   getTimestamp(msg),
    UpdatedAt:   getTimestamp(msg),
    Version:     1,
})
if err != nil {
    slog.Error("Failed to add item", "error", err, "name", cmd.Name)
    return fmt.Errorf("failed to add item: %w", err)
}
```