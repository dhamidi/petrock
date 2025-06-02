# Refactor cmd/petrock/test.go Plan

**FEATURE OBJECTIVE**: Refactor the monolithic `runTest` function in `cmd/petrock/test.go` into a modular, step-based testing framework where each test operation is isolated, trackable, and has clear success/failure reporting.

**ACCEPTANCE CRITERIA**:
- Test execution consists of discrete, named steps
- Each step has isolated setup, execution, and cleanup phases
- Step results include success/failure status and associated logs
- Main test function reads as setup followed by step execution
- Failed steps provide clear error context without affecting other steps
- Test output shows progress through each step with clear success/failure indicators

**IMPLEMENTATION PHASES**:

PHASE 1: Define Core Testing Types
- Task 1.1: [File: cmd/petrock/test_types.go] Create step and result types (Effort: Small, Dependencies: None) - COMPLETED
  - Types: `TestStep` interface, `StepResult` struct, `TestContext` struct, `StepFunc` type
  - Functions: `NewStepResult()`, `(*StepResult).Success()`, `(*StepResult).Failure()`
  - Packages: `fmt`, `time`, `log/slog`

- Task 1.2: [File: cmd/petrock/test_types.go] Add step execution framework (Effort: Small, Dependencies: Task 1.1) - COMPLETED
  - Types: `TestRunner` struct
  - Functions: `NewTestRunner()`, `(*TestRunner).RunStep()`, `(*TestRunner).RunAllSteps()`
  - Packages: `context`

Phase 2: Extract Individual Test Steps
- Task 2.1: [File: cmd/petrock/test_steps.go] Create environment setup steps (Effort: Medium, Dependencies: Phase 1)
  - Types: `SetupTempDirStep`, `CreateProjectStep`, `AddFeatureStep`
  - Functions: `NewSetupTempDirStep()`, `NewCreateProjectStep()`, `NewAddFeatureStep()`
  - Packages: `os`, `path/filepath`

- Task 2.2: [File: cmd/petrock/test_steps.go] Create build and server steps (Effort: Medium, Dependencies: Phase 1)
  - Types: `BuildProjectStep`, `StartServerStep`, `StopServerStep`
  - Functions: `NewBuildProjectStep()`, `NewStartServerStep()`, `NewStopServerStep()`
  - Packages: `os/exec`, `time`

- Task 2.3: [File: cmd/petrock/test_steps.go] Create HTTP testing steps (Effort: Medium, Dependencies: Phase 1)
  - Types: `HTTPGetStep`, `HTTPPostStep`, `CommandAPIStep`
  - Functions: `NewHTTPGetStep()`, `NewHTTPPostStep()`, `NewCommandAPIStep()`
  - Packages: `net/http`, `encoding/json`, `bytes`, `net/url`

- Task 2.4: [File: cmd/petrock/test_steps.go] Create inspection testing steps (Effort: Small, Dependencies: Phase 1)
  - Types: `SelfInspectStep`
  - Functions: `NewSelfInspectStep()`, `validateInspectJSON()`
  - Packages: `encoding/json`, `os/exec`

Phase 3: Implement Step Interface Methods
- Task 3.1: [File: cmd/petrock/test_steps.go] Implement Execute() for setup steps (Effort: Medium, Dependencies: Phase 2)
  - Types: Methods on `SetupTempDirStep`, `CreateProjectStep`, `AddFeatureStep`
  - Functions: `(*SetupTempDirStep).Execute()`, `(*CreateProjectStep).Execute()`, `(*AddFeatureStep).Execute()`
  - Packages: None (using existing functionality)

- Task 3.2: [File: cmd/petrock/test_steps.go] Implement Execute() for build/server steps (Effort: Medium, Dependencies: Phase 2)
  - Types: Methods on `BuildProjectStep`, `StartServerStep`, `StopServerStep`
  - Functions: `(*BuildProjectStep).Execute()`, `(*StartServerStep).Execute()`, `(*StopServerStep).Execute()`
  - Packages: None (using existing functionality)

- Task 3.3: [File: cmd/petrock/test_steps.go] Implement Execute() for HTTP and inspection steps (Effort: Medium, Dependencies: Phase 2)
  - Types: Methods on `HTTPGetStep`, `HTTPPostStep`, `CommandAPIStep`, `SelfInspectStep`
  - Functions: `(*HTTPGetStep).Execute()`, `(*HTTPPostStep).Execute()`, `(*CommandAPIStep).Execute()`, `(*SelfInspectStep).Execute()`
  - Packages: None (using existing functionality)

Phase 4: Refactor Main Test Function
- Task 4.1: [File: cmd/petrock/test.go] Replace runTest with step-based implementation (Effort: Large, Dependencies: Phase 3)
  - Types: Modify `runTest` function signature and implementation
  - Functions: `runTest()`, `setupTestSteps()`, `executeTestSteps()`
  - Packages: None (reorganizing existing code)

- Task 4.2: [File: cmd/petrock/test.go] Add step progress reporting (Effort: Small, Dependencies: Task 4.1)
  - Types: None
  - Functions: `reportStepProgress()`, `reportFinalResults()`
  - Packages: `fmt`

Phase 5: Add Context and Cleanup Management
- Task 5.1: [File: cmd/petrock/test_context.go] Implement shared test context (Effort: Medium, Dependencies: Phase 4)
  - Types: `TestContext` struct with fields for temp dirs, server processes, etc.
  - Functions: `NewTestContext()`, `(*TestContext).Cleanup()`, `(*TestContext).SetTempDir()`
  - Packages: `os`, `os/exec`

- Task 5.2: [File: cmd/petrock/test.go] Integrate context management (Effort: Small, Dependencies: Task 5.1)
  - Types: Update step implementations to use context
  - Functions: Modify step Execute methods to accept and use TestContext
  - Packages: None

**CRITICAL PATH**: 
1. Phase 1 (core types) must be completed before any step implementation
2. Task 4.1 (main refactor) blocks final integration and testing
3. Context management (Phase 5) should be completed before production use

**TECHNICAL RISKS**:
- Risk 1: Step isolation complexity → Mitigation: Use shared TestContext for state that must persist between steps
- Risk 2: Error handling across step boundaries → Mitigation: Standardize error wrapping and context preservation in StepResult
- Risk 3: Server process lifecycle management → Mitigation: Implement robust cleanup in TestContext with defer patterns
- Risk 4: Parallel step execution complexity → Mitigation: Start with sequential execution, add parallelism in future iteration

**BUILD REQUIREMENTS**:
- Go Dependencies: No new external dependencies (uses existing cobra, slog, etc.)
- Build Commands: `./build.sh` should continue to work, `go test ./cmd/petrock` for unit testing
- Testing Strategy: Maintain existing integration test behavior while adding unit tests for individual steps

**FIRST IMPLEMENTATION STEP**: Create `cmd/petrock/test_types.go` with the basic `TestStep` interface and `StepResult` struct:

```go
type TestStep interface {
    Name() string
    Execute(ctx *TestContext) *StepResult
}

type StepResult struct {
    StepName  string
    Success   bool
    Error     error
    Duration  time.Duration
    Logs      []string
    StartTime time.Time
}
```

**IMPLEMENTATION DETAILS**:

**TestStep Interface Design**:
```go
type TestStep interface {
    Name() string                    // Human-readable step name
    Execute(ctx *TestContext) *StepResult // Execute the step
}

type StepResult struct {
    StepName  string        // Step identifier
    Success   bool          // Whether step passed
    Error     error         // Error if failed
    Duration  time.Duration // Execution time
    Logs      []string      // Step-specific log messages
    StartTime time.Time     // When step started
}

type TestContext struct {
    TempDir     string          // Base temporary directory
    ProjectDir  string          // Generated project directory
    ProjectName string          // Name of generated project
    ModulePath  string          // Go module path
    ServerCmd   *exec.Cmd       // Running server process
    ServerPort  string          // Server port number
    Cleanup     []func() error  // Cleanup functions
}
```

**Step Categories**:
1. **Setup Steps**: `SetupTempDirStep`, `CreateProjectStep`, `AddFeatureStep`
2. **Build Steps**: `BuildProjectStep`, `BuildServerStep`
3. **Server Steps**: `StartServerStep`, `StopServerStep`
4. **Test Steps**: `HTTPGetStep`, `HTTPPostStep`, `CommandAPIStep`, `SelfInspectStep`

**Error Handling Strategy**:
- Each step captures its own errors without affecting others
- TestContext maintains cleanup functions for proper resource management
- StepResult provides structured error information with context
- Main test function reports overall success based on all step results

**Logging Strategy**:
- Each step maintains its own log buffer
- slog continues to provide structured logging
- Step results include relevant log messages for debugging
- Progress reporting shows step completion status

This refactoring transforms the monolithic test function into a modular, maintainable testing framework while preserving all existing test functionality and improving error reporting and debugging capabilities.
