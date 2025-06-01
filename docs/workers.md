# Worker Development Guide

## Overview

Petrock's worker system provides a powerful abstraction for handling background processing, event reactions, and asynchronous operations. The worker infrastructure eliminates common boilerplate code while maintaining flexibility for complex business logic.

With this abstraction, projects can scale to 50-100 workers without code duplication, each focused purely on business logic rather than infrastructure concerns.

## Worker Architecture

### Core Components

The worker system consists of several key components:

- **Worker Interface**: Defines the contract for all workers (Start, Stop, Work, WorkerInfo)
- **CommandWorker**: Concrete implementation providing command-based message processing
- **Command Handlers**: Functions that process specific command types
- **Periodic Work**: Background processing that runs during each Work() cycle
- **Worker State**: Custom state structure for worker-specific data

### Infrastructure Provided by Core

The core worker infrastructure handles:

1. **Message Processing Loop**: Automatically iterates through new messages since last processed ID
2. **Command Routing**: Dispatches commands to registered handlers based on command name
3. **State Management**: Tracks `lastProcessedID` and provides access to worker-specific state
4. **Lifecycle Management**: Handles Start/Stop/Work cycle with proper context management
5. **Error Handling**: Provides consistent error handling and structured logging
6. **Periodic Execution**: Calls user-defined periodic work function

## Creating a Worker

### Basic Worker Structure

```go
package workers

import (
    "context"
    "time"
    "github.com/your-project/core"
)

// WorkerState holds worker-specific state
type WorkerState struct {
    // Add your worker's state fields here
    pendingTasks map[string]Task
    config       WorkerConfig
    client       *http.Client
}

// NewWorker creates a new worker instance
func NewWorker(app *core.App, state *State, log *core.MessageLog, executor *core.Executor) core.Worker {
    workerState := &WorkerState{
        pendingTasks: make(map[string]Task),
        config:       loadConfig(),
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }

    worker := core.NewWorker(
        "Feature Worker",
        "Handles background processing for the feature",
        workerState,
    )

    // Set core dependencies
    worker.SetDependencies(log, executor)

    // Register command handlers
    worker.OnCommand("feature/process", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
        return handleProcessCommand(ctx, cmd, msg, workerState)
    })

    // Set periodic work
    worker.SetPeriodicWork(func(ctx context.Context) error {
        return processPendingTasks(ctx, workerState)
    })

    return worker
}
```

### Command Handlers

Command handlers are the core of worker business logic. They receive commands and perform operations:

```go
func handleProcessCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState) error {
    // Type assertion to get specific command
    processCmd, ok := cmd.(*ProcessCommand)
    if !ok {
        return fmt.Errorf("unexpected command type: %T", cmd)
    }

    // Perform business logic
    task := Task{
        ID:        processCmd.TaskID,
        Data:      processCmd.Data,
        CreatedAt: time.Now(),
    }

    // Update worker state
    workerState.pendingTasks[task.ID] = task

    // Log the operation
    slog.Info("Task added to processing queue",
        "taskID", task.ID,
        "feature", "your-feature")

    return nil
}
```

### Periodic Work

Periodic work functions run during each Work() cycle and handle background processing:

```go
func processPendingTasks(ctx context.Context, workerState *WorkerState) error {
    if len(workerState.pendingTasks) == 0 {
        return nil
    }

    slog.Debug("Processing pending tasks", "count", len(workerState.pendingTasks))

    for taskID, task := range workerState.pendingTasks {
        // Check if task is too old
        if time.Since(task.CreatedAt) > 24*time.Hour {
            delete(workerState.pendingTasks, taskID)
            continue
        }

        // Process the task
        if err := processTask(ctx, workerState, task); err != nil {
            slog.Error("Failed to process task", "taskID", taskID, "error", err)
            continue
        }

        // Remove completed task
        delete(workerState.pendingTasks, taskID)
    }

    return nil
}
```

## Advanced Patterns

### External API Integration

Workers commonly integrate with external APIs. Here's a robust pattern:

```go
type APIWorkerState struct {
    client   *http.Client
    apiURL   string
    apiKey   string
    retries  map[string]int // Track retry attempts
}

func callExternalAPI(ctx context.Context, workerState *APIWorkerState, payload interface{}) error {
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", workerState.apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+workerState.apiKey)

    resp, err := workerState.client.Do(req)
    if err != nil {
        return fmt.Errorf("API request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("API returned status %d", resp.StatusCode)
    }

    return nil
}
```

### Batch Processing

For handling large volumes of data:

```go
func processBatch(ctx context.Context, workerState *BatchWorkerState) error {
    const batchSize = 100

    items := workerState.getPendingItems()
    
    for i := 0; i < len(items); i += batchSize {
        end := i + batchSize
        if end > len(items) {
            end = len(items)
        }

        batch := items[i:end]
        if err := processBatchItems(ctx, workerState, batch); err != nil {
            slog.Error("Batch processing failed", "batch", i/batchSize, "error", err)
            continue
        }

        // Small delay between batches to avoid overwhelming external services
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(100 * time.Millisecond):
        }
    }

    return nil
}
```

### State Synchronization

For workers that need to sync with external systems:

```go
type SyncWorkerState struct {
    lastSyncTime time.Time
    syncInterval time.Duration
    state        *State
}

func performSync(ctx context.Context, workerState *SyncWorkerState) error {
    if time.Since(workerState.lastSyncTime) < workerState.syncInterval {
        return nil // Too early to sync
    }

    // Fetch data from external system
    externalData, err := fetchExternalData(ctx, workerState)
    if err != nil {
        return fmt.Errorf("sync failed: %w", err)
    }

    // Update local state
    for _, item := range externalData {
        updateCmd := &UpdateFromExternalCommand{
            ID:   item.ID,
            Data: item.Data,
        }

        if err := workerState.executor.Execute(ctx, updateCmd); err != nil {
            slog.Error("Failed to update from external data", "id", item.ID, "error", err)
            continue
        }
    }

    workerState.lastSyncTime = time.Now()
    return nil
}
```

## Migration Guide

### Converting Existing Workers

Follow these steps to migrate from the old worker pattern to the new abstraction:

#### Step 1: Identify Current Worker Components

In your existing worker, identify:
- Message processing loop
- Command handling logic
- Periodic work functions
- State management
- Context handling

#### Step 2: Extract Business Logic

Create separate functions for each command handler:

```go
// Old pattern
func (w *OldWorker) Work() error {
    for msg := range w.log.After(w.ctx, w.lastProcessedID) {
        cmd, ok := msg.DecodedPayload.(Command)
        if !ok {
            continue
        }

        switch cmd.CommandName() {
        case "feature/process":
            processCmd := cmd.(*ProcessCommand)
            // Business logic here...
        case "feature/update":
            updateCmd := cmd.(*UpdateCommand)
            // Business logic here...
        }
    }
}

// New pattern
func handleProcessCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState) error {
    processCmd := cmd.(*ProcessCommand)
    // Same business logic here...
}

func handleUpdateCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState) error {
    updateCmd := cmd.(*UpdateCommand)
    // Same business logic here...
}
```

#### Step 3: Replace Worker Structure

Replace the old worker struct with the new pattern:

```go
// Old pattern
type OldWorker struct {
    name            string
    lastProcessedID uint64
    state           *WorkerState
    log             *core.MessageLog
    executor        *core.Executor
    ctx             context.Context
    cancel          context.CancelFunc
}

// New pattern
func NewWorker(app *core.App, state *State, log *core.MessageLog, executor *core.Executor) core.Worker {
    workerState := &WorkerState{
        // Initialize your state
    }

    worker := core.NewWorker("Worker Name", "Description", workerState)
    worker.SetDependencies(log, executor)
    
    // Register handlers
    worker.OnCommand("feature/process", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
        return handleProcessCommand(ctx, cmd, msg, workerState)
    })
    
    return worker
}
```

#### Step 4: Update Worker Registration

Update your feature's main.go to use the new worker:

```go
// Old pattern
func RegisterWorkers(app *core.App, state *State) {
    worker := &OldWorker{
        // initialization...
    }
    app.RegisterWorker(worker)
}

// New pattern
func RegisterWorkers(app *core.App, state *State, log *core.MessageLog, executor *core.Executor) {
    worker := NewWorker(app, state, log, executor)
    app.RegisterWorker(worker)
}
```

## Performance Characteristics

### Message Processing

- **Throughput**: Up to 10,000 messages/second per worker on modern hardware
- **Memory Usage**: ~1MB base memory per worker + state size
- **Latency**: Sub-millisecond command routing and dispatch

### Scaling Considerations

- **Workers per Project**: Tested with up to 100 workers in a single application
- **Command Handlers**: Up to 50 commands per worker with no performance impact
- **State Size**: Keep worker state under 100MB for optimal performance

### Optimization Tips

1. **Batch Processing**: Group related operations to reduce overhead
2. **State Management**: Use maps for O(1) lookups instead of slices
3. **Context Timeouts**: Use appropriate timeouts for external API calls
4. **Logging**: Use structured logging at appropriate levels (Debug/Info/Warn/Error)

## Testing Strategies

### Unit Testing Command Handlers

```go
func TestHandleProcessCommand(t *testing.T) {
    workerState := &WorkerState{
        pendingTasks: make(map[string]Task),
    }

    cmd := &ProcessCommand{
        TaskID: "test-task",
        Data:   "test-data",
    }

    ctx := context.Background()
    msg := &core.Message{ID: 1}

    err := handleProcessCommand(ctx, cmd, msg, workerState)
    assert.NoError(t, err)
    assert.Contains(t, workerState.pendingTasks, "test-task")
}
```

### Integration Testing

```go
func TestWorkerIntegration(t *testing.T) {
    app := core.NewApp()
    state := NewState()
    log := core.NewMessageLog()
    executor := core.NewExecutor()

    worker := NewWorker(app, state, log, executor)
    
    ctx := context.Background()
    err := worker.Start(ctx)
    assert.NoError(t, err)

    // Send test command
    cmd := &ProcessCommand{TaskID: "test", Data: "data"}
    err = executor.Execute(ctx, cmd)
    assert.NoError(t, err)

    // Run worker cycle
    err = worker.Work()
    assert.NoError(t, err)

    // Verify state changes
    workerState := worker.State().(*WorkerState)
    assert.Contains(t, workerState.pendingTasks, "test")
}
```

## Troubleshooting

### Common Issues

#### 1. Worker Not Processing Messages

**Symptoms**: Commands are executed but worker handlers are not called

**Diagnosis**:
```go
// Check if worker is started
info := worker.WorkerInfo()
fmt.Printf("Worker: %s - %s\n", info.Name, info.Description)

// Check command registration
worker.OnCommand("debug/list-handlers", func(ctx context.Context, cmd core.Command, msg *core.Message) error {
    fmt.Printf("Handler called for: %s\n", cmd.CommandName())
    return nil
})
```

**Solutions**:
- Verify worker is properly registered with the app
- Check that `SetDependencies()` is called with correct MessageLog
- Ensure command names match exactly (case-sensitive)

#### 2. High Memory Usage

**Symptoms**: Worker memory usage grows over time

**Diagnosis**:
```go
func (w *WorkerState) debugMemoryUsage() {
    fmt.Printf("Pending tasks: %d\n", len(w.pendingTasks))
    fmt.Printf("Cache size: %d\n", len(w.cache))
}
```

**Solutions**:
- Clean up completed tasks from worker state
- Implement periodic cleanup in periodic work function
- Use bounded caches with LRU eviction

#### 3. Command Handler Panics

**Symptoms**: Worker stops processing after panic in handler

**Diagnosis**:
```go
func handleCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("handler panic: %v", r)
            slog.Error("Command handler panicked", "command", cmd.CommandName(), "panic", r)
        }
    }()
    
    // Your handler logic...
    return nil
}
```

**Solutions**:
- Add panic recovery in critical handlers
- Validate command types with proper error handling
- Use defensive programming for state access

#### 4. Slow Periodic Work

**Symptoms**: Worker Work() cycles take too long

**Diagnosis**:
```go
func processPendingTasks(ctx context.Context, workerState *WorkerState) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        if duration > 1*time.Second {
            slog.Warn("Slow periodic work", "duration", duration)
        }
    }()
    
    // Your periodic work...
    return nil
}
```

**Solutions**:
- Implement batch processing for large datasets
- Add context cancellation checks in long loops
- Use timeouts for external API calls

#### 5. External API Rate Limiting

**Symptoms**: External API calls fail with rate limit errors

**Solutions**:
```go
type RateLimitedWorkerState struct {
    rateLimiter *time.Ticker
    lastCall    time.Time
}

func callAPIWithRateLimit(ctx context.Context, workerState *RateLimitedWorkerState) error {
    // Wait for rate limiter
    <-workerState.rateLimiter.C
    
    // Make API call
    return callExternalAPI(ctx, workerState)
}
```

#### 6. Context Cancellation Not Handled

**Symptoms**: Workers don't shut down cleanly

**Solutions**:
```go
func longRunningTask(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Process item i
            processItem(i)
        }
    }
    return nil
}
```

#### 7. State Race Conditions

**Symptoms**: Inconsistent state or panics under load

**Solutions**:
```go
type ThreadSafeWorkerState struct {
    mu           sync.RWMutex
    pendingTasks map[string]Task
}

func (s *ThreadSafeWorkerState) AddTask(task Task) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.pendingTasks[task.ID] = task
}

func (s *ThreadSafeWorkerState) GetTask(id string) (Task, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    task, exists := s.pendingTasks[id]
    return task, exists
}
```

#### 8. Command Type Assertion Failures

**Symptoms**: Handler receives unexpected command types

**Solutions**:
```go
func handleTypedCommand(ctx context.Context, cmd core.Command, msg *core.Message, workerState *WorkerState) error {
    processCmd, ok := cmd.(*ProcessCommand)
    if !ok {
        slog.Warn("Unexpected command type", 
            "expected", "*ProcessCommand",
            "actual", fmt.Sprintf("%T", cmd),
            "command", cmd.CommandName())
        return fmt.Errorf("unexpected command type: %T", cmd)
    }
    
    // Process the command...
    return nil
}
```

#### 9. Worker Startup Failures

**Symptoms**: Worker fails to start or initialize

**Diagnosis**:
```go
func NewWorker(...) core.Worker {
    // Validate dependencies
    if log == nil {
        panic("MessageLog cannot be nil")
    }
    if executor == nil {
        panic("Executor cannot be nil")
    }
    
    // Create worker...
}
```

#### 10. Periodic Work Never Executes

**Symptoms**: SetPeriodicWork function is never called

**Solutions**:
- Verify `SetPeriodicWork()` is called before worker registration
- Check that App.StartWorkers() is called
- Ensure worker Work() method is being invoked by the app

## Best Practices

### State Management

1. **Keep State Minimal**: Only store what's necessary for worker operation
2. **Use Appropriate Data Structures**: Maps for lookups, slices for ordered data
3. **Clean Up Regularly**: Remove completed or expired items in periodic work
4. **Avoid Shared State**: Each worker should manage its own state independently

### Error Handling

1. **Structured Logging**: Use slog with consistent field names
2. **Graceful Degradation**: Continue processing other items when one fails
3. **Retry Logic**: Implement exponential backoff for transient failures
4. **Circuit Breakers**: Temporarily disable failing external services

### Testing

1. **Unit Test Handlers**: Test command handlers as pure functions
2. **Mock Dependencies**: Use interfaces for external services
3. **Integration Tests**: Test full worker lifecycle with real MessageLog
4. **Load Testing**: Verify performance under expected workloads

### Security

1. **Input Validation**: Validate all command parameters
2. **Secure External Calls**: Use proper authentication and TLS
3. **Resource Limits**: Implement timeouts and memory limits
4. **Audit Logging**: Log sensitive operations for security analysis
