# Centralized Command Execution Plan

## Problem Statement

Currently, command execution logic is duplicated across feature command handlers. Each handler must:
1. Validate the command
2. Persist the command to the message log
3. Update the in-memory state

This duplication leads to boilerplate code, potential inconsistencies in how commands are processed, and difficulty in changing the execution flow globally.

## Proposed Solution

Implement a centralized executor pattern with the following components:

### 1. Core Executor Interface and Implementation

```go
// core/executor.go

// Executor defines the contract for command execution
type Executor interface {
    Execute(ctx context.Context, cmd Command) error
}

// BaseExecutor provides a standard implementation of command execution flow
type BaseExecutor struct {
    log         *MessageLog
    cmdRegistry *CommandRegistry
}

// NewBaseExecutor creates a new executor with dependencies
func NewBaseExecutor(log *MessageLog, cmdRegistry *CommandRegistry) *BaseExecutor {
    // Implementation
}

// Execute implements the standard command execution flow:
// 1. Validate the command
// 2. Log the command
// 3. Dispatch to appropriate handler
func (e *BaseExecutor) Execute(ctx context.Context, cmd Command) error {
    // Implementation of standard execution flow
}
```

### 2. Validator Interface for Commands

```go
// core/commands.go

// Validator interface for self-validating commands
type Validator interface {
    Validate() error
}

// ValidationError provides structured validation errors
type ValidationError struct {
    CommandName string
    Fields      map[string]string
    Err         error
}
```

### 3. Updated Feature Command Handlers

Feature command handlers will be simplified to focus on domain logic:

```go
// posts/execute.go

// HandleCreatePost focuses only on business logic, not execution flow
func (e *PostExecutor) HandleCreatePost(ctx context.Context, cmd core.Command) error {
    createCmd := cmd.(CreatePostCommand)
    
    // Focus on domain logic only
    // No validation or logging here
    return e.state.Apply(createCmd)
}
```

### 4. Self-Validating Commands

Commands will implement the Validator interface:

```go
// posts/messages.go

// Validate checks if the command is valid
func (c CreatePostCommand) Validate() error {
    if c.Title == "" {
        return errors.New("title cannot be empty")
    }
    if c.Content == "" {
        return errors.New("content cannot be empty")
    }
    return nil
}
```

## Benefits

1. **Reduced duplication**: The standard execution flow is defined once in the core executor
2. **Separation of concerns**: Command handlers focus on domain logic, not infrastructure concerns
3. **Consistent handling**: All commands are processed through the same execution flow
4. **Centralized error handling**: Validation errors are handled consistently
5. **Easier to modify**: Global changes to execution flow can be made in one place

## Implementation Steps

### Step 1: Define Executor Interface and Implementation
**Files to modify:**
- Create new file: `core/executor.go`

**Changes:**
- Define `Executor` interface with `Execute(ctx context.Context, cmd Command) error` method
- Create `BaseExecutor` struct with dependencies on `MessageLog` and `CommandRegistry`
- Implement constructor `NewBaseExecutor(log *MessageLog, cmdRegistry *CommandRegistry) *BaseExecutor`
- Implement `Execute()` method with standard validation → logging → dispatch flow
- Implement error handling for validation failures vs. execution failures

**Definition of Done:**
- New `core/executor.go` file created with all required functionality
- Unit tests for `BaseExecutor` verify validation checking, logging, and dispatching
- Error types and handling patterns defined and tested

**Goal:** Create a centralized command execution component that standardizes how all commands are processed, eliminating duplicated logic across feature handlers.

### Step 2: Add Validation Support to Core Commands
**Files to modify:**
- `core/commands.go`

**Changes:**
- Add `Validator` interface with `Validate() error` method
- Create `ValidationError` type with fields for command name, error message, and field-specific errors
- Add helper functions `IsValidationError(err error) bool` and `NewValidationError(...) *ValidationError`
- Update documentation for `Command` interface to reference validation pattern

**Definition of Done:**
- Validator interface and error types added to `core/commands.go`
- Error creation helpers and check functions implemented and tested
- Unit tests verify error creation and type checking functions

**Goal:** Enable commands to self-validate and provide a standardized way to report validation errors.

### Step 3: Update Feature Command Types
**Files to modify:**
- Feature message files (e.g., `posts/messages.go`)
- Templates for code generation (`internal/skeleton/feature_template/messages.go`)

**Changes:**
- Add `Validate()` methods to all command types in each feature
- Move validation logic from handlers to the command types
- Ensure commands implement the `Validator` interface
- Update code templates to include `Validate()` methods for new features

**Definition of Done:**
- All existing commands implement `Validate()` methods with proper validation logic
- Template file updated to include validation in generated commands
- Unit tests for command validations verify error reporting

**Goal:** Move validation logic from handlers to command objects, enabling the executor to standardize validation checking.

### Step 4: Refactor Feature Command Handlers
**Files to modify:**
- Feature execution files (e.g., `posts/execute.go`)
- Templates for code generation (`internal/skeleton/feature_template/execute.go`)

**Changes:**
- Remove validation logic from handler methods
- Remove message logging logic from handler methods
- Focus handlers solely on domain-specific state manipulation
- Update constructor to no longer require direct MessageLog access
- Update code templates for new features

**Definition of Done:**
- Handlers in all features simplified to focus only on domain logic
- No direct MessageLog interactions in handlers
- Template file updated to generate simplified handlers
- Unit tests for handlers focus on state manipulation, not validation

**Goal:** Simplify command handlers to focus exclusively on domain logic without infrastructure concerns.

### Step 5: Update HTTP and API Handlers
**Files to modify:**
- Core API handlers (`core/http.go` or similar)
- Feature-specific HTTP handlers (`posts/http.go`, etc.)
- `cmd/<app>/serve.go` for initialization

**Changes:**
- Add Executor as a dependency to HTTP handlers
- Replace direct CommandRegistry and MessageLog calls with Executor.Execute()
- Update initialization code to create and inject the BaseExecutor
- Update error handling to process ValidationErrors appropriately

**Definition of Done:**
- All command execution flows use the Executor interface
- HTTP handlers properly format validation errors for API responses
- No direct MessageLog.Append() or CommandRegistry.Dispatch() calls remain
- Integration tests verify end-to-end command execution via HTTP

**Goal:** Ensure all application entry points use the centralized executor pattern for consistent command processing.

### Step 6: Update Documentation and Examples – DONE
**Files to modify:**
- `docs/high-level.md`
- `docs/core/*.md` files
- `docs/feature/*.md` files
- Code examples and inline documentation

**Changes:**
- Update architecture documentation to describe the executor pattern
- Add examples showing how to use the Executor
- Update feature documentation to reflect simplified handler approach
- Ensure inline code comments reflect the new pattern

**Definition of Done:**
- All documentation consistently describes the executor pattern
- Examples show proper usage of the executor
- No documentation references the old direct dispatch approach

**Goal:** Ensure documentation and examples consistently reflect the new architecture for both users and developers.

## Migration Strategy

When generating new features with petrock, the executor pattern will be used by default. For existing features, a migration path will be provided to adapt them to the new pattern.