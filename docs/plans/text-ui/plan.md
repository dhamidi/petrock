# Text UI Refactoring Plan

## FEATURE OBJECTIVE
Refactor command-line I/O in `./cmd/` and `./internal/skeleton/cmd/petrock_example_project_name/` to use a unified UI interface instead of direct logging and print statements. Commands should communicate with users exclusively through this interface, separating user interaction from debug logging.

## ACCEPTANCE CRITERIA
- [ ] All user-facing output in commands goes through UI interface methods
- [ ] Debug/internal logging remains separate using slog for development purposes
- [ ] UI interface supports progress indication, success/error messages, and user prompts
- [ ] Commands can be tested with mock UI implementations
- [ ] No fmt.Printf/fmt.Println statements remain in command logic
- [ ] Consistent user experience across all commands

## IMPLEMENTATION PHASES

### Phase 1: Create UI Interface and Refactor ./cmd/

**Phase 1.1: Define UI Interface**
- Task 1.1.1: [File: internal/ui/interface.go] Create UI interface definition (Effort: 2h, Dependencies: None) - COMPLETED
  - Types: UI interface, MessageType enum, ProgressState struct
  - Functions: Present(), Prompt(), ShowProgress(), ShowError(), ShowSuccess()
  - Packages: context, io

- Task 1.1.2: [File: internal/ui/console.go] Implement console UI (Effort: 3h, Dependencies: 1.1.1) - COMPLETED
  - Types: ConsoleUI struct, OutputWriter interface
  - Functions: NewConsoleUI(), Present(), Prompt(), ShowProgress(), formatMessage()
  - Packages: fmt, os, bufio, strings

- Task 1.1.3: [File: internal/ui/mock.go] Implement mock UI for testing (Effort: 2h, Dependencies: 1.1.1)
  - Types: MockUI struct, CapturedMessage struct
  - Functions: NewMockUI(), Present(), GetMessages(), ClearMessages()
  - Packages: sync, testing

**Phase 1.2: Refactor cmd/petrock/main.go**
- Task 1.2.1: [File: cmd/petrock/main.go] Add UI dependency injection (Effort: 2h, Dependencies: 1.1.2)
  - Types: CommandContext struct
  - Functions: newCommandContext(), injectUI()
  - Packages: github.com/dhamidi/petrock/internal/ui

- Task 1.2.2: [File: cmd/petrock/main.go] Separate debug logging from user output (Effort: 1h, Dependencies: 1.2.1)
  - Types: modify rootCmd configuration
  - Functions: configureLogging(), configureSlogLevel()
  - Packages: log/slog, os

**Phase 1.3: Refactor cmd/petrock/new.go**
- Task 1.3.1: [File: cmd/petrock/new.go] Replace fmt.Printf with UI.Present() (Effort: 2h, Dependencies: 1.2.1)
  - Types: modify newCmd cobra.Command
  - Functions: runNew(), showProjectCreationProgress()
  - Packages: github.com/dhamidi/petrock/internal/ui

- Task 1.3.2: [File: cmd/petrock/new.go] Add progress indication for long operations (Effort: 2h, Dependencies: 1.3.1)
  - Types: ProgressTracker struct
  - Functions: trackProgress(), updateProgress()
  - Packages: time, context

**Phase 1.4: Refactor cmd/petrock/feature.go**
- Task 1.4.1: [File: cmd/petrock/feature.go] Replace fmt.Printf with UI.Present() (Effort: 3h, Dependencies: 1.2.1)
  - Types: modify featureCmd cobra.Command, FeatureCreationSteps enum
  - Functions: runFeature(), showFeatureCreationProgress(), presentNextSteps()
  - Packages: github.com/dhamidi/petrock/internal/ui

- Task 1.4.2: [File: cmd/petrock/feature.go] Remove slog.Info calls for user messages (Effort: 1h, Dependencies: 1.4.1)
  - Types: none
  - Functions: modify runFeature(), insertFeatureRegistration()
  - Packages: remove user-facing slog calls

**Phase 1.5: Refactor cmd/petrock/test.go**
- Task 1.5.1: [File: cmd/petrock/test.go] Replace fmt.Printf with UI.Present() (Effort: 2h, Dependencies: 1.2.1)
  - Types: modify testCmd cobra.Command, TestResultFormatter struct
  - Functions: runTest(), formatTestResults(), showTestProgress()
  - Packages: github.com/dhamidi/petrock/internal/ui

- Task 1.5.2: [File: cmd/petrock/test_types.go] Update test execution to use UI (Effort: 2h, Dependencies: 1.5.1)
  - Types: modify TestRunner struct
  - Functions: ExecuteTests(), RunStep(), handleStepResult()
  - Packages: github.com/dhamidi/petrock/internal/ui

**Phase 1.6: Refactor remaining cmd/petrock/ files**
- Task 1.6.1: [File: cmd/petrock/new_command.go] Replace slog.Info with UI.Present() (Effort: 1h, Dependencies: 1.2.1)
  - Types: modify newCommandCmd cobra.Command
  - Functions: runNewCommand()
  - Packages: github.com/dhamidi/petrock/internal/ui

- Task 1.6.2: [File: cmd/petrock/new_query.go] Replace slog.Info with UI.Present() (Effort: 1h, Dependencies: 1.2.1)
  - Types: modify newQueryCmd cobra.Command
  - Functions: runNewQuery()
  - Packages: github.com/dhamidi/petrock/internal/ui

- Task 1.6.3: [File: cmd/petrock/new_worker.go] Replace slog.Info with UI.Present() (Effort: 1h, Dependencies: 1.2.1)
  - Types: modify newWorkerCmd cobra.Command
  - Functions: runNewWorker()
  - Packages: github.com/dhamidi/petrock/internal/ui

### Phase 2: Refactor ./internal/skeleton/cmd/petrock_example_project_name/

**Phase 2.1: Copy UI Interface to Skeleton**
- Task 2.1.1: [File: internal/skeleton/core/ui/interface.go] Copy UI interface to skeleton (Effort: 1h, Dependencies: Phase 1 complete) - COMPLETED
  - Types: UI interface, MessageType enum, ProgressState struct
  - Functions: Present(), Prompt(), ShowProgress(), ShowError(), ShowSuccess()
  - Packages: context, io

- Task 2.1.2: [File: internal/skeleton/core/ui/console.go] Copy console UI implementation (Effort: 1h, Dependencies: 2.1.1) - COMPLETED
  - Types: ConsoleUI struct, OutputWriter interface
  - Functions: NewConsoleUI(), Present(), Prompt(), ShowProgress()
  - Packages: fmt, os, bufio, strings

**Phase 2.2: Refactor skeleton main.go**
- Task 2.2.1: [File: internal/skeleton/cmd/petrock_example_project_name/main.go] Add UI dependency injection (Effort: 2h, Dependencies: 2.1.2) - COMPLETED
  - Types: CommandContext struct
  - Functions: newCommandContext(), configureUI()
  - Packages: github.com/petrock/example_module_path/core/ui

**Phase 2.3: Refactor skeleton commands**
- Task 2.3.1: [File: internal/skeleton/cmd/petrock_example_project_name/kv.go] Replace fmt.Printf with UI.Present() (Effort: 1h, Dependencies: 2.2.1) - COMPLETED
  - Types: modify kvSetCmd and kvListCmd cobra.Commands
  - Functions: runKvSet(), runKvList()
  - Packages: github.com/petrock/example_module_path/core/ui

- Task 2.3.2: [File: internal/skeleton/cmd/petrock_example_project_name/serve.go] Replace slog.Info user messages with UI.Present() (Effort: 2h, Dependencies: 2.2.1) - COMPLETED
  - Types: modify serveCmd cobra.Command
  - Functions: runServe(), showStartupProgress()
  - Packages: github.com/petrock/example_module_path/core/ui

- Task 2.3.3: [File: internal/skeleton/cmd/petrock_example_project_name/build.go] Replace slog.Info user messages with UI.Present() (Effort: 1h, Dependencies: 2.2.1) - COMPLETED
  - Types: modify buildCmd cobra.Command
  - Functions: runBuild(), showBuildProgress()
  - Packages: github.com/petrock/example_module_path/core/ui

- Task 2.3.4: [File: internal/skeleton/cmd/petrock_example_project_name/deploy.go] Replace slog.Info user messages with UI.Present() (Effort: 2h, Dependencies: 2.2.1) - COMPLETED
  - Types: modify deployCmd cobra.Command
  - Functions: runDeploy(), showDeploymentProgress()
  - Packages: github.com/petrock/example_module_path/core/ui

## CRITICAL PATH
1. Task 1.1.1 (UI interface definition) - blocks all other tasks
2. Task 1.1.2 (console UI implementation) - blocks all cmd refactoring
3. Task 1.2.1 (main.go UI injection) - blocks all command refactoring
4. Phase 1 completion - blocks Phase 2 start

## TECHNICAL RISKS

- **Risk 1**: UI interface too rigid for diverse command needs → **Mitigation**: Design interface with flexible message types and extensible formatting options
- **Risk 2**: Progress indication complexity → **Mitigation**: Start with simple text-based progress, iterate based on usage
- **Risk 3**: Breaking existing command behavior → **Mitigation**: Maintain exact same user-visible output initially, then enhance
- **Risk 4**: Test suite complexity with UI mocking → **Mitigation**: Design MockUI with comprehensive capture capabilities from start
- **Risk 5**: Template placeholders in skeleton UI code → **Mitigation**: Use same placeholder replacement strategy as existing skeleton files

## BUILD REQUIREMENTS

**Go Dependencies**: 
- No new external modules needed
- Uses existing cobra, slog packages

**Build Commands**:
- `./build.sh` - validate template compilation pipeline
- `go test ./cmd/petrock/...` - test command functionality
- `go test ./internal/ui/...` - test UI implementations

**Testing Strategy**:
- Unit tests for UI implementations with MockUI
- Integration tests for commands using MockUI
- Manual testing of console output formatting
- Template generation tests to verify skeleton UI works

## FIRST IMPLEMENTATION STEP

**File**: `internal/ui/interface.go`
**Function**: Define the UI interface with these methods:
```go
type UI interface {
    Present(ctx context.Context, msgType MessageType, message string, args ...interface{}) error
    Prompt(ctx context.Context, question string) (string, error)
    ShowProgress(ctx context.Context, state ProgressState) error
    ShowError(ctx context.Context, err error) error
    ShowSuccess(ctx context.Context, message string, args ...interface{}) error
}
```

This establishes the contract that all subsequent implementations will depend on.
