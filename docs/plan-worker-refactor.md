# Worker Abstraction Specification and Implementation Plan

## Specification

### Overview

The current worker implementation in Petrock requires significant boilerplate code for common operations like message iteration, state tracking, and command routing. With projects potentially having 50-100 workers, we need to abstract this infrastructure into the core while keeping business logic simple and clear.

### Core Abstraction: Pattern-Based Workers

#### BaseWorker Structure
```go
type BaseWorker struct {
    name            string
    description     string
    lastProcessedID uint64
    state           interface{}
    patterns        map[string]PatternHandler
    periodicWork    func(context.Context) error
    log             *MessageLog
    executor        *Executor
    ctx             context.Context
    cancel          context.CancelFunc
}

type PatternHandler func(ctx context.Context, cmd Command, state interface{}) error
```

#### Core-Provided Infrastructure
1. **Message Processing Loop**: Automatically iterates through new messages since last processed ID
2. **Command Routing**: Dispatches commands to registered pattern handlers based on command name
3. **State Management**: Tracks `lastProcessedID` and provides access to worker-specific state
4. **Lifecycle Management**: Handles Start/Stop/Work cycle with proper context management
5. **Error Handling**: Provides consistent error handling and logging
6. **Periodic Execution**: Calls user-defined periodic work function

#### Feature-Specific Implementation
Features only need to provide:
1. **Pattern Handlers**: Simple functions that handle specific command types
2. **Periodic Work**: Business logic for background processing
3. **State Structure**: Worker-specific state (if needed)

#### Example Usage
```go
// In feature/workers/main.go
func NewWorker(app *core.App, state *State) core.Worker {
    worker := core.NewPatternWorker(
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

func handleCreate(ctx context.Context, cmd core.Command, state interface{}) error {
    createCmd := cmd.(*CreateCommand)
    wState := state.(*WorkerState)
    // Only business logic - no infrastructure code
    return requestSummarization(ctx, createCmd.Name)
}
```

### Benefits
- **Reduced Boilerplate**: ~80% reduction in worker code
- **Consistent Infrastructure**: All workers use the same message processing, error handling, and lifecycle management
- **Easier Testing**: Business logic handlers are pure functions
- **Scalable**: Easy to create 50-100 workers without duplication
- **Incremental Migration**: Can migrate existing workers gradually

---

## Plan

### Step 1: Create BaseWorker Infrastructure
**Files Modified:**
- `core/worker.go` - Add PatternWorker implementation and supporting types
- `core/app.go` - Update worker management to support new pattern workers

**Modifications:**
- Add `PatternWorker` struct with message processing loop
- Add `PatternHandler` function type and command routing logic
- Add `NewPatternWorker()` constructor and fluent API methods (`OnCommand`, `SetPeriodicWork`)
- Update `App.StartWorkers()` to work with both old and new worker types
- Add helper methods for state access and context management

**Acceptance Criteria:**
1. `PatternWorker` can be instantiated with name, description, and initial state
2. `OnCommand()` method successfully registers handlers for specific command names
3. `SetPeriodicWork()` method accepts and stores a periodic work function that gets called during Work() cycles
4. Message processing loop correctly iterates through messages after `lastProcessedID` and routes commands to registered handlers
5. Worker state is properly encapsulated and accessible to handlers through the state parameter
6. Context management works correctly with timeout handling and cancellation support

### Step 2: Migrate Example Feature Worker
**Files Modified:**
- `petrock_example_feature_name/workers/main.go` - Replace current worker with PatternWorker usage
- `petrock_example_feature_name/workers/summary_worker.go` - Extract business logic into pattern handlers
- `petrock_example_feature_name/workers/types.go` - Update type definitions as needed

**Modifications:**
- Replace `Worker` struct with `NewWorker()` function returning `core.Worker`
- Convert `handleCreateCommand`, `handleSummaryRequestCommand`, etc. into pattern handlers
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

### Step 3: Add Advanced Features and Utilities
**Files Modified:**
- `core/worker.go` - Add retry logic, metrics, and advanced configuration options
- `core/app.go` - Add worker monitoring and health check capabilities
- `petrock_example_feature_name/workers/summary_worker.go` - Demonstrate usage of advanced features

**Modifications:**
- Add retry logic with exponential backoff for failed pattern handlers
- Add metrics collection (message processing rate, error counts, handler execution time)
- Add configurable timeouts for pattern handlers and periodic work
- Add health check interface for workers to report their status
- Add structured logging with consistent format across all workers
- Add graceful degradation when external services are unavailable

**Acceptance Criteria:**
1. Pattern handlers that return errors are automatically retried with exponential backoff (1s, 2s, 4s, max 30s)
2. Worker metrics are collected and accessible through app inspection (messages/sec, errors/minute, avg handler time)
3. Workers can report health status (healthy, degraded, unhealthy) based on recent error rates and external service availability
4. Worker timeouts are enforced: pattern handlers timeout after 5s, periodic work after 30s, with configurable overrides
5. All worker operations produce structured logs with consistent fields (worker_name, command_type, execution_time, error_details)
6. When external services fail, workers gracefully degrade (skip processing, log warnings) rather than crashing or blocking

### Step 4: Create Worker Testing Framework
**Files Modified:**
- `core/worker_test.go` - Add comprehensive test framework for PatternWorker
- `petrock_example_feature_name/workers/main_test.go` - Create example worker tests using new framework
- `core/testing.go` - Add test utilities for worker behavior verification

**Modifications:**
- Create test framework that allows mocking message log, executor, and external services
- Add test utilities for simulating message sequences and verifying handler calls
- Create helpers for testing periodic work functions and error scenarios
- Add performance testing utilities for worker throughput and latency
- Create integration test helpers that work with real message log and state

**Acceptance Criteria:**
1. Test framework can simulate arbitrary message sequences and verify correct handler invocation and state changes
2. Mock utilities allow testing worker behavior without external dependencies (database, HTTP services)
3. Test helpers can verify timing behavior (periodic work intervals, timeout handling, retry delays)
4. Performance tests can measure worker throughput (messages/second) and verify it meets minimum thresholds (>100 msg/sec)
5. Integration tests can verify worker behavior with real message log replay and state persistence
6. Error injection tests verify workers handle failures gracefully (network errors, timeout errors, invalid commands)

### Step 5: Documentation and Migration Guide
**Files Modified:**
- `docs/workers.md` - Create comprehensive worker development guide
- `README.md` - Update with worker abstraction information
- `petrock_example_feature_name/workers/README.md` - Add feature-specific worker documentation

**Modifications:**
- Write developer guide explaining PatternWorker usage and best practices
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
