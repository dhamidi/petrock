# Worker Abstraction Specification and Implementation Plan

## Specification

### Overview

The current worker implementation in Petrock requires significant boilerplate code for common operations like message iteration, state tracking, and command routing. With projects potentially having 50-100 workers, we need to abstract this infrastructure into the core while keeping business logic simple and clear.

### Core Abstraction: Command-Based Workers

#### Worker Structure
```go
type Worker struct {
    name            string
    description     string
    lastProcessedID uint64
    state           interface{}
    handlers        map[string]CommandHandler
    periodicWork    func(context.Context) error
    log             *MessageLog
    executor        *Executor
    ctx             context.Context
    cancel          context.CancelFunc
}

// Use existing CommandHandler from core/commands.go
type CommandHandler func(ctx context.Context, cmd Command, msg Message) error
```

#### Core-Provided Infrastructure
1. **Message Processing Loop**: Automatically iterates through new messages since last processed ID
2. **Command Routing**: Dispatches commands to registered command handlers based on command name
3. **State Management**: Tracks `lastProcessedID` and provides access to worker-specific state
4. **Lifecycle Management**: Handles Start/Stop/Work cycle with proper context management
5. **Error Handling**: Provides consistent error handling and logging
6. **Periodic Execution**: Calls user-defined periodic work function

#### Feature-Specific Implementation
Features only need to provide:
1. **Command Handlers**: Simple functions that handle specific command types
2. **Periodic Work**: Business logic for background processing
3. **State Structure**: Worker-specific state (if needed)

#### Example Usage
```go
// In feature/workers/main.go
func NewWorker(app *core.App, state *State) core.Worker {
    worker := core.NewWorker(
        "feature_name Worker",
        "Handles background processing for feature",
        &WorkerState{pendingSummaries: make(map[string]PendingSummary)},
    )
    
    worker.OnCommand("feature/create", handleCreate)
    worker.OnCommand("feature/request-summary", handleSummaryRequest)
    
    worker.SetPeriodicWork(func(ctx context.Context) error {
        return processPendingSummaries(ctx, worker.State().(*WorkerState))
    })
    
    return worker
}

func handleCreate(ctx context.Context, cmd core.Command, msg core.Message) error {
    createCmd := cmd.(*CreateCommand)
    wState := worker.State().(*WorkerState)
    // Only business logic - no infrastructure code
    return requestSummarization(ctx, createCmd.Name)
}
```

### Benefits
- **Reduced Boilerplate**: ~80% reduction in worker code
- **Consistent Infrastructure**: All workers use the same message processing, error handling, and lifecycle management
- **Easier Testing**: Business logic command handlers are pure functions
- **Scalable**: Easy to create 50-100 workers without duplication
- **Incremental Migration**: Can migrate existing workers gradually

---

## Plan

### Step 1: Create Worker Infrastructure - DONE
**Files Modified:**
- `core/worker.go` - Add new Worker implementation and supporting types
- `core/app.go` - Update worker management to support new workers

**Modifications:**
- Add new `Worker` struct with message processing loop
- Use existing `CommandHandler` function type and add command routing logic
- Add `NewWorker()` constructor and fluent API methods (`OnCommand`, `SetPeriodicWork`)
- Update `App.StartWorkers()` to work with both old and new worker types
- Add helper methods for state access and context management

**Acceptance Criteria:**
1. `Worker` can be instantiated with name, description, and initial state
2. `OnCommand()` method successfully registers command handlers for specific command names
3. `SetPeriodicWork()` method accepts and stores a periodic work function that gets called during Work() cycles
4. Message processing loop correctly iterates through messages after `lastProcessedID` and routes commands to registered handlers
5. Worker state is properly encapsulated and accessible to handlers
6. Context management works correctly with timeout handling and cancellation support

### Step 2: Migrate Example Feature Worker - DONE
**Files Modified:**
- `petrock_example_feature_name/workers/main.go` - Replace current worker with PatternWorker usage
- `petrock_example_feature_name/workers/summary_worker.go` - Extract business logic into pattern handlers
- `petrock_example_feature_name/workers/types.go` - Update type definitions as needed

**Modifications:**
- Replace existing `Worker` struct with `NewWorker()` function returning `core.Worker`
- Convert `handleCreateCommand`, `handleSummaryRequestCommand`, etc. into command handlers
- Extract `processPendingSummaries` into standalone periodic work function
- Remove infrastructure code (message iteration, ID tracking, context management)
- Simplify error handling to focus on business logic

**Acceptance Criteria:**
1. New worker implementation produces identical behavior to the original worker for all command types
2. Summary generation still works end-to-end: create item → request summary → process summary → store result
3. Worker successfully processes existing messages from the log during startup replay
4. Periodic summary processing continues to work with the same timing and retry logic
5. All existing tests pass without modification (demonstrates behavioral compatibility)
6. Worker code is reduced by at least 60% while maintaining all functionality

### Step 3: Documentation and Migration Guide
**Files Modified:**
- `docs/workers.md` - Create comprehensive worker development guide
- `README.md` - Update with worker abstraction information
- `petrock_example_feature_name/workers/README.md` - Add feature-specific worker documentation

**Modifications:**
- Write developer guide explaining Worker usage and best practices
- Create migration guide for converting existing workers to new pattern
- Document performance characteristics and scaling considerations
- Add troubleshooting guide for common worker issues
- Create examples for different worker patterns (event processing, batch jobs, external integrations)

**Acceptance Criteria:**
1. Documentation includes complete working examples that can be copy-pasted and modified for new features
2. Migration guide provides step-by-step instructions that allow converting existing workers in <30 minutes
3. Performance guide includes specific benchmarks and optimization recommendations for different workload types
4. Troubleshooting guide covers at least 10 common issues with diagnostic steps and solutions
5. Best practices section includes guidance on state management, error handling, and testing strategies
6. All code examples in documentation are tested and verified to work with the current implementation
