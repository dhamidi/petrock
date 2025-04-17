# Technical Implementation Plan for Workers

## Overview

Workers provide a mechanism for features to run background processes that react to events in the message log. Workers are ideal for handling cross-cutting concerns such as interacting with external services, processing data asynchronously, or implementing business logic that spans multiple events.

This plan outlines how to replace the current `jobs.go` pattern with a more robust worker system that better integrates with the event sourcing architecture of Petrock applications.

## Detailed Task Breakdown

### T1: Core Worker Interface - DONE

**T1.1:** Define the Worker interface in core/worker.go - DONE
- Create new file core/worker.go with Context import - DONE
- Define Worker interface with Start, Stop, and Work methods - DONE
- Add documentation comments for each method - DONE

**T1.2:** Create error types for worker operations - DONE
- Define appropriate error types for worker initialization and processing failures - DONE
- Add documentation for error handling - DONE

**Definition of Done for T1:**
- File core/worker.go exists with properly documented Worker interface - DONE
- All methods have proper documentation with usage examples - DONE
- Interface is in line with the event sourcing architecture in core/ - DONE

### T2: App Worker Management

**T2.1:** Add worker tracking to App struct - DONE
- Add workers slice to App struct in core/app.go - DONE
- Add worker control fields (workerCtx, workerCancel, workerWg) - DONE

**T2.2:** Implement RegisterWorker method - DONE
- Add RegisterWorker method to App struct - DONE
- Method should accept Worker interface and add to the workers slice - DONE
- Add appropriate logging - DONE

**T2.3:** Implement StartWorkers method - DONE
- Create StartWorkers method that initializes workers with context - DONE
- Start each worker in its own goroutine - DONE
- Implement ticker with jitter for periodic Work() calls - DONE
- Handle initialization errors properly - DONE

**T2.4:** Implement StopWorkers method - DONE
- Create StopWorkers method that signals workers to stop - DONE
- Wait for workers to finish with timeout - DONE
- Handle cleanup failures properly - DONE

**T2.5:** Update Close method - DONE
- Modify existing Close method to call StopWorkers - DONE
- Ensure proper error handling during shutdown - DONE

**Definition of Done for T2: - DONE**
- App struct has worker management fields - DONE
- RegisterWorker correctly adds workers to the app - DONE
- StartWorkers successfully starts workers in goroutines with randomized intervals - DONE
- StopWorkers gracefully shuts down workers and handles timeouts - DONE
- Close method properly includes worker shutdown in its sequence - DONE
- Logging statements exist for major state changes - DONE

### T3: Feature Template Worker Implementation

**T3.1:** Create worker.go template file - DONE
- Create internal/skeleton/feature_template/worker.go - DONE
- Use correct package declaration: package petrock_example_feature_name - DONE
- Use correct import path: github.com/petrock/example_module_path/core - DONE

**T3.2:** Define worker state struct - DONE
- Create struct for tracking pending operations (like summarization) - DONE
- Include lastProcessedID field for tracking message log position - DONE
- Include mutex for thread safety - DONE

**T3.3:** Define worker struct - DONE
- Create worker struct with app, executor, state, log dependencies - DONE
- Include internal state reference - DONE

**T3.4:** Implement Start method - DONE
- Add Start method that scans message log from beginning - DONE
- Process messages to rebuild internal state - DONE
- Update lastProcessedID as messages are processed - DONE
- Return appropriate errors - DONE

**T3.5:** Implement Stop method - DONE
- Add Stop method for graceful shutdown - DONE
- Clean up any resources used by the worker - DONE
- Return appropriate errors - DONE

**T3.6:** Implement Work method - DONE
- Add Work method that processes new messages since lastProcessedID - DONE
- Update worker state based on messages - DONE
- Perform any pending tasks - DONE
- Return appropriate errors - DONE

**T3.7:** Add helper methods - DONE
- Create processMessage method for handling different command types - DONE
- Create utility methods for worker-specific tasks - DONE

**Definition of Done for T3: - DONE**
- worker.go exists in feature template directory - DONE
- All worker methods are properly implemented with correct error handling - DONE
- Worker correctly tracks its position in the message log - DONE
- Code follows project's conventions and uses appropriate placeholders - DONE
- File has proper documentation for all methods and types - DONE

### T4: Feature Template Registration

**T4.1:** Update RegisterFeature in register.go
- Modify internal/skeleton/feature_template/register.go
- Add worker initialization code in the RegisterFeature function
- Ensure worker is registered with the app

**Definition of Done for T4:**
- RegisterFeature in register.go initializes and registers the worker
- Worker is created with the correct dependencies
- Code follows existing conventions and style

### T5: Serve Command Integration - DONE

**T5.1:** Update serve.go to start workers - DONE
- Modify cmd/template/serve.go to call app.StartWorkers after feature registration - DONE
- Add proper error handling - DONE

**T5.2:** Update shutdown sequence - DONE
- Add app.StopWorkers call during server shutdown - DONE
- Ensure proper cancellation context is passed - DONE

**Definition of Done for T5: - DONE**
- serve.go calls StartWorkers at the appropriate point in startup - DONE
- StopWorkers is called during shutdown before server resources are released - DONE
- Both calls have proper error handling - DONE

### T6: Migration Strategy


**T6.1:** Update petrock feature command
- Update feature generation to use worker.go instead of jobs.go

**Definition of Done for T6:**
- feature command generates worker.go instead of jobs.go
- Existing applications have a clear path to adopt the new pattern

## Feature Template Worker Implementation Details

### Worker Interface

The core Worker interface will be defined in `core/worker.go`:

- `Worker` interface with methods:
  - `Start(context.Context) error` - Initializes worker state from the message log
  - `Stop(context.Context) error` - Cleans up resources when shutting down
  - `Work() error` - Performs processing on each work cycle

### App Worker Management

Extend `core/app.go` to include worker management:

- Add a slice of workers to the App struct
- Add worker control fields (context, cancel func, wait group)
- Implement methods:
  - `RegisterWorker(Worker)` - Adds a worker to the app
  - `StartWorkers(context.Context) error` - Initializes and starts workers
  - `StopWorkers(context.Context) error` - Gracefully stops workers
  - Update `Close()` to stop workers

### Post Summarization Worker Template

The feature template in `internal/skeleton/feature_template/worker.go` will implement a post summarization worker pattern. This worker will:

1. Track content that needs to be summarized using an external service
2. Process relevant commands:
   - Respond to `petrock_example_feature_name/create` by requesting summarization
   - Track `petrock_example_feature_name/request-summary-generation` requests
   - Handle `petrock_example_feature_name/fail-summary-generation` responses
   - Process `petrock_example_feature_name/set-generated-summary` responses

3. Make API calls to an external service when content needs summarization
4. Dispatch commands to update the application state with results

The implementation will include:

```go
// In worker.go:

// PendingSummary tracks a content item waiting for summarization
type PendingSummary struct {
	RequestID string
	ItemID    string
	Content   string
	CreatedAt time.Time
}

// WorkerState holds worker-specific state
type WorkerState struct {
	mu               sync.Mutex
	lastProcessedID  string
	pendingSummaries map[string]PendingSummary // keyed by RequestID
}

// Worker implements background processing for the feature
type Worker struct {
	app      *core.App
	executor *core.Executor
	state    *State
	log      *core.MessageLog
	
	// Worker's internal state
	wState   *WorkerState
	
	// Configuration for external service
	apiURL   string
	apiKey   string
	client   *http.Client
}
```

### Feature Template Registration

Update `internal/skeleton/feature_template/register.go` to include worker registration in the `RegisterFeature` function:

```go
// --- 7. Register Worker (replacing jobs registration) ---
worker := NewWorker(app, state, app.MessageLog, app.Executor)
app.RegisterWorker(worker)
```

### Serve Command Integration

Update `cmd/project/serve.go` template to:

1. Call `app.StartWorkers(ctx)` after registering features
2. Call `app.StopWorkers(ctx)` during graceful shutdown

## Benefits

1. Better integration with the event sourcing architecture
2. More structured lifecycle management
3. Improved separation of concerns
4. Consistent tracking of message log position