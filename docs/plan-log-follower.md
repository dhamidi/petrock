# Plan: Log Follower Abstraction - DONE

## Overview

Introduce a new abstraction for workers to maintain and persist their position in the message log, enabling reliable replay and preventing duplicate processing across restarts.

## New Interfaces

### 1. LogFollower Interface

```go
// LogFollower maintains a position in the message log
type LogFollower interface {
    // LogPosition returns the current position in the log
    LogPosition() uint64
    
    // LogSeek sets the position to a new value
    LogSeek(newPosition uint64)
}
```

**Purpose**: Track where a worker has processed up to in the message log.

### 2. KVStore Interface

```go
// KVStore provides persistent key-value storage
type KVStore interface {
    // Get retrieves a value by key and unmarshals it into dest
    Get(key string, dest any) error
    
    // Set stores a value by key, marshaling it appropriately
    Set(key string, value any) error
}
```

**Location**: `core/kv.go`

**Default Implementation**: SQLite-backed table in the main database
- Table: `kv_store` with columns `key TEXT PRIMARY KEY, value TEXT`
- Values stored as JSON for simplicity
- Uses same database connection as message log

## Application Changes

### 1. App Initialization

```go
type App struct {
    // existing fields...
    log     *MessageLog
    kvStore *KVStore  // new
}

func NewApp() *App {
    db := openDatabase()
    return &App{
        log:     NewMessageLog(db),
        kvStore: NewSQLiteKVStore(db),  // new
        // ...
    }
}
```

### 2. Worker Interface Updates

```go
type Worker interface {
    // existing methods...
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Work() error
    WorkerInfo() *WorkerInfo
    
    // new method for initial replay
    Replay(ctx context.Context) error
}
```

### 3. App Startup Sequence

```go
func (app *App) Start() error {
    // 1. Initialize core components
    if err := app.initializeDatabase(); err != nil {
        return err
    }
    
    // 2. Register features (creates workers)
    for _, feature := range app.features {
        feature.Register(app)
    }
    
    // 3. Start workers (creates LogFollowers)
    for _, worker := range app.workers {
        if err := worker.Start(ctx); err != nil {
            return err
        }
    }
    
    // 4. Replay messages to catch up workers
    for _, worker := range app.workers {
        if err := worker.Replay(ctx); err != nil {
            return err
        }
    }
    
    // 5. Start periodic worker execution
    app.startWorkerScheduler()
    
    // 6. Start HTTP server (now safe to accept requests)
    return app.startHTTPServer()
}
```

## Worker Implementation Changes

### 1. CommandWorker Updates

```go
type ProcessingContext struct {
    IsReplay bool  // true during replay, false during normal processing
}

// Updated handler signature to support replay mode
type CommandHandler func(ctx context.Context, cmd Command, msg *Message, pctx *ProcessingContext) error

type CommandWorker struct {
    // existing fields...
    follower LogFollower  // new
    kvStore  *KVStore     // new
    
    // remove lastProcessedID - now handled by LogFollower
}

func NewWorker(name, description string, initialState interface{}) *CommandWorker {
    return &CommandWorker{
        name:        name,
        description: description,
        state:       initialState,
        handlers:    make(map[string]CommandHandler),
        follower:    NewLogFollower(),  // new
    }
}

func (w *CommandWorker) SetDependencies(log *MessageLog, executor *Executor, kvStore *KVStore) {
    w.log = log
    w.executor = executor
    w.kvStore = kvStore  // new
    w.follower.SetKVStore(kvStore, w.positionKey())  // new
}

func (w *CommandWorker) positionKey() string {
    return fmt.Sprintf("worker:%s:position", w.name)
}
```

### 2. New Replay Method

```go
func (w *CommandWorker) Replay(ctx context.Context) error {
    // Load last known position from KV store
    if err := w.follower.LoadPosition(w.kvStore, w.positionKey()); err != nil {
        // If no position found, start from beginning
        w.follower.LogSeek(0)
    }
    
    startPosition := w.follower.LogPosition()
    slog.Info("Starting message replay", "worker", w.name, "fromPosition", startPosition)
    
    // Replay ALL messages from beginning to reconstruct state
    messageCount := 0
    for msg := range w.log.After(ctx, 0) {  // Always start from beginning for state reconstruction
        if err := w.replayMessage(msg); err != nil {  // Use separate replay method
            slog.Error("Failed to replay message", "worker", w.name, "id", msg.ID, "error", err)
            continue
        }
        
        messageCount++
    }
    
    // Set position to latest message ID to avoid reprocessing
    if messageCount > 0 {
        // Get the latest message ID
        latest := w.log.Latest(ctx)
        if latest != nil {
            w.follower.LogSeek(latest.ID)
        }
    }
    
    // Save final position
    if err := w.follower.SavePosition(w.kvStore, w.positionKey()); err != nil {
        return fmt.Errorf("failed to save final position: %w", err)
    }
    
    slog.Info("Message replay completed", "worker", w.name, "messagesProcessed", messageCount, "finalPosition", w.follower.LogPosition())
    return nil
}

// replayMessage processes a message for state reconstruction only (no side effects)
func (w *CommandWorker) replayMessage(msg PersistedMessage) error {
    cmd, ok := msg.DecodedPayload.(Command)
    if !ok {
        return nil
    }

    handler, found := w.handlers[cmd.CommandName()]
    if !found {
        return nil
    }

    // Create replay context to indicate this is state-only processing
    replayCtx := &ProcessingContext{IsReplay: true}
    
    return handler(w.ctx, cmd, &msg.Message, replayCtx)
}
```

### 3. Updated Work Method

```go
func (w *CommandWorker) Work() error {
    if !w.started {
        return ErrWorkerStopped
    }

    // Process only new messages (after current position)
    currentPosition := w.follower.LogPosition()
    messagesProcessed := 0
    
    for msg := range w.log.After(w.ctx, currentPosition) {
        if err := w.processMessage(msg); err != nil {
            slog.Error("Failed to process message", "worker", w.name, "id", msg.ID, "error", err)
            continue
        }
        
        w.follower.LogSeek(msg.ID)
        messagesProcessed++
    }

// processMessage handles normal message processing with side effects
func (w *CommandWorker) processMessage(msg PersistedMessage) error {
    cmd, ok := msg.DecodedPayload.(Command)
    if !ok {
        return nil
    }

    handler, found := w.handlers[cmd.CommandName()]
    if !found {
        return nil
    }

    // Create normal processing context (allows side effects)
    normalCtx := &ProcessingContext{IsReplay: false}
    
    return handler(w.ctx, cmd, &msg.Message, normalCtx)
}
    
    // Persist position if we processed any messages
    if messagesProcessed > 0 {
        if err := w.follower.SavePosition(w.kvStore, w.positionKey()); err != nil {
            slog.Error("Failed to save worker position", "worker", w.name, "error", err)
        }
    }

    // Execute periodic work
    if w.periodicWork != nil {
        if err := w.periodicWork(w.ctx); err != nil {
            return fmt.Errorf("periodic work failed: %w", err)
        }
    }

    return nil
}
```

## Implementation Files

### 1. `core/log_follower.go`
- `LogFollower` interface
- `SimpleLogFollower` implementation with position tracking
- Methods for loading/saving position via KVStore

### 2. `core/kv.go`
- `KVStore` interface
- `SQLiteKVStore` implementation
- Database table creation and JSON marshaling

### 3. Updated `core/worker.go`
- Add `Replay` method to `Worker` interface
- Update `CommandWorker` to use `LogFollower` and `KVStore`
- Remove manual `lastProcessedID` tracking

### 4. Updated `core/app.go`
- Add KVStore to App struct
- Update startup sequence to replay before accepting requests
- Pass KVStore to workers via `SetDependencies`

## Handler Implementation Pattern

Workers must implement handlers that distinguish between replay and normal processing:

```go
// Example: Posts worker handles summary requests
func (w *PostsWorker) handleRequestSummary(ctx context.Context, cmd Command, msg *Message, pctx *ProcessingContext) error {
    reqCmd := cmd.(*RequestSummaryCommand)
    
    // ALWAYS update internal state (both replay and normal)
    w.state.AddPendingSummary(reqCmd.ItemID, reqCmd.RequestID)
    
    // Skip side effects during replay
    if pctx.IsReplay {
        return nil
    }
    
    // Execute side effects only during normal processing
    slog.Info("Added content to pending summarization queue", 
        "feature", "posts", "itemID", reqCmd.ItemID, "requestID", reqCmd.RequestID)
    
    return w.generateSummary(ctx, reqCmd)
}

func (w *PostsWorker) handleSetSummary(ctx context.Context, cmd Command, msg *Message, pctx *ProcessingContext) error {
    setCmd := cmd.(*SetSummaryCommand)
    
    // ALWAYS update internal state
    w.state.SetSummary(setCmd.ItemID, setCmd.Summary)
    
    // No side effects for this command - it's a state update
    return nil
}
```

## Benefits

1. **Reliable Position Tracking**: Workers remember where they left off
2. **Clean Separation**: LogFollower handles position, KVStore handles persistence  
3. **Side-Effect Free Replay**: Replay only reconstructs state, no duplicate work
4. **Ordered Startup**: Replay completes before serving requests
5. **Extensible**: KVStore can be used for other worker state
6. **Testable**: Interfaces make mocking easy

## Migration Path

1. Implement new interfaces and SQLite KVStore
2. Add LogFollower to CommandWorker (keeping old logic)
3. Update App startup sequence
4. Test with existing workers
5. Remove old `lastProcessedID` tracking

This design ensures workers maintain consistent state across restarts while providing a clean, testable abstraction for position tracking.
