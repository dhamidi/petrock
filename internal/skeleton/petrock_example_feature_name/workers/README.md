# petrock_example_feature_name Worker

This worker demonstrates the command-based worker pattern for handling background processing in Petrock applications.

## Overview

The petrock_example_feature_name worker handles:
- **Content Summarization**: Processes new items and generates summaries via external API
- **Background Processing**: Manages pending summarization requests with retry logic
- **State Management**: Tracks summarization status and handles failures gracefully

## Architecture

### Worker State

```go
type WorkerState struct {
    pendingSummaries map[string]PendingSummary // Tracks ongoing summarization requests
    state            *State                     // Reference to application state
    executor         *core.Executor             // Command execution
    apiURL           string                     // External API configuration
    apiKey           string
    client           *http.Client               // HTTP client for API calls
}
```

### Command Handlers

The worker responds to these commands:

1. **`petrock_example_feature_name/create`**: When new items are created, automatically requests summarization
2. **`petrock_example_feature_name/request-summary-generation`**: Adds content to the pending summarization queue
3. **`petrock_example_feature_name/fail-summary-generation`**: Removes failed requests from the queue
4. **`petrock_example_feature_name/set-generated-summary`**: Removes completed requests from the queue

### Periodic Work

The worker's periodic function:
- Processes pending summarization requests
- Calls external API for content summarization
- Handles timeouts and failures
- Cleans up old requests (24-hour timeout)

## Usage Example

### Creating Items

When you create a new item via the web interface or API:

```bash
curl -X POST http://localhost:8080/petrock_example_feature_name/new \
  -d "name=My Article" \
  -d "content=This is the content to be summarized..."
```

The flow is:
1. `CreateCommand` is executed
2. Worker's `handleCreateCommand` receives the command
3. Worker automatically requests summarization via `RequestSummaryGenerationCommand`
4. Content is added to the pending summarization queue
5. Periodic work processes the queue and calls external API
6. Summary is stored via `SetGeneratedSummaryCommand`

### Monitoring Worker Activity

Check logs for worker activity:

```bash
go run ./cmd/yourapp serve --log-level=debug
```

You'll see logs like:
```
INFO  Worker started name=petrock_example_feature_name Worker
DEBUG Processing pending summaries count=3
INFO  Added content to pending summarization queue itemID=article-1 requestID=req-123
INFO  Successfully generated summary itemID=article-1 requestID=req-123
```

## Configuration

### External API Setup

In a real application, configure the external API:

```go
workerState := &WorkerState{
    // Configure from environment variables
    apiURL: os.Getenv("SUMMARIZATION_API_URL"),
    apiKey: os.Getenv("SUMMARIZATION_API_KEY"),
    client: &http.Client{
        Timeout: 30 * time.Second,
    },
}
```

### Environment Variables

```bash
export SUMMARIZATION_API_URL="https://api.openai.com/v1/completions"
export SUMMARIZATION_API_KEY="your-api-key-here"
```

## Implementation Details

### Mock API Implementation

The current implementation includes a mock API call that:
- Simulates network latency (500ms-1.5s delay)
- Generates fake summaries for demonstration
- Includes commented code showing real HTTP API integration

### Real API Integration

To use a real API service, uncomment and modify the code in `callSummarizationAPI()`:

```go
// Example for OpenAI API
type SummarizeRequest struct {
    Prompt      string `json:"prompt"`
    MaxTokens   int    `json:"max_tokens"`
    Temperature float32 `json:"temperature"`
}

type SummarizeResponse struct {
    Choices []struct {
        Text string `json:"text"`
    } `json:"choices"`
}

func callSummarizationAPI(ctx context.Context, workerState *WorkerState, summary PendingSummary) error {
    prompt := fmt.Sprintf("Summarize the following content: %s", summary.Content)
    
    payload := SummarizeRequest{
        Prompt:      prompt,
        MaxTokens:   150,
        Temperature: 0.7,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal request: %w", err)
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
    
    var result SummarizeResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }
    
    if len(result.Choices) == 0 {
        return fmt.Errorf("no summary generated")
    }
    
    // Use the generated summary
    generatedSummary := strings.TrimSpace(result.Choices[0].Text)
    
    // Store the summary
    setCmd := &SetGeneratedSummaryCommand{
        ID:        summary.ItemID,
        RequestID: summary.RequestID,
        Summary:   generatedSummary,
    }
    
    if err := workerState.executor.Execute(ctx, setCmd); err != nil {
        return fmt.Errorf("failed to store summary: %w", err)
    }
    
    return nil
}
```

## Error Handling

### Timeout Handling

Requests older than 24 hours are automatically failed:

```go
if time.Since(summary.CreatedAt) > 24*time.Hour {
    failCmd := &FailSummaryGenerationCommand{
        ID:        summary.ItemID,
        RequestID: summary.RequestID,
        Reason:    "timeout",
    }
    // Execute failure command...
}
```

### API Failure Handling

When API calls fail, the worker:
1. Logs the error with structured logging
2. Keeps the request in the pending queue for retry on the next cycle
3. Eventually times out after 24 hours

### Context Cancellation

All operations respect context cancellation:
- API calls use `context.WithTimeout`
- Command execution uses separate contexts
- Graceful shutdown is supported

## Testing

### Unit Testing Command Handlers

```go
func TestHandleCreateCommand(t *testing.T) {
    workerState := &WorkerState{
        pendingSummaries: make(map[string]PendingSummary),
        executor:         &mockExecutor{},
    }
    
    cmd := &CreateCommand{
        Name:    "test-item",
        Content: "test content",
    }
    
    ctx := context.Background()
    msg := &core.Message{ID: 1}
    
    err := handleCreateCommand(ctx, cmd, msg, workerState)
    assert.NoError(t, err)
    
    // Verify that RequestSummaryGenerationCommand was executed
    assert.True(t, workerState.executor.commandExecuted)
}
```

### Integration Testing

```go
func TestWorkerIntegration(t *testing.T) {
    // Set up test environment
    app := core.NewApp()
    state := NewState()
    log := core.NewMessageLog()
    executor := core.NewExecutor()
    
    // Create worker
    worker := NewWorker(app, state, log, executor)
    
    // Start worker
    ctx := context.Background()
    err := worker.Start(ctx)
    require.NoError(t, err)
    
    // Create test item
    createCmd := &CreateCommand{
        Name:    "test-item",
        Content: "This is test content that should be summarized.",
    }
    
    err = executor.Execute(ctx, createCmd)
    require.NoError(t, err)
    
    // Run worker cycle
    err = worker.Work()
    require.NoError(t, err)
    
    // Verify summarization was requested
    workerState := worker.State().(*WorkerState)
    assert.NotEmpty(t, workerState.pendingSummaries)
    
    // Run another cycle to process pending summaries
    err = worker.Work()
    require.NoError(t, err)
    
    // Verify summary was generated and stored
    item, found := state.GetItem("test-item")
    require.True(t, found)
    assert.NotEmpty(t, item.Summary)
}
```

## Performance Considerations

### Batch Processing

For high-volume scenarios, consider batching API calls:

```go
func processPendingSummariesBatched(ctx context.Context, workerState *WorkerState) error {
    const batchSize = 10
    
    summaries := getSummariesAsList(workerState.pendingSummaries)
    
    for i := 0; i < len(summaries); i += batchSize {
        end := i + batchSize
        if end > len(summaries) {
            end = len(summaries)
        }
        
        batch := summaries[i:end]
        if err := processBatch(ctx, workerState, batch); err != nil {
            slog.Error("Batch processing failed", "error", err)
            continue
        }
        
        // Rate limiting between batches
        time.Sleep(100 * time.Millisecond)
    }
    
    return nil
}
```

### Rate Limiting

For APIs with rate limits:

```go
type RateLimitedWorkerState struct {
    *WorkerState
    lastAPICall time.Time
    minInterval time.Duration
}

func callAPIWithRateLimit(ctx context.Context, workerState *RateLimitedWorkerState, summary PendingSummary) error {
    // Enforce rate limiting
    if time.Since(workerState.lastAPICall) < workerState.minInterval {
        return nil // Skip this cycle
    }
    
    err := callSummarizationAPI(ctx, workerState.WorkerState, summary)
    if err == nil {
        workerState.lastAPICall = time.Now()
    }
    
    return err
}
```

## Monitoring and Observability

### Metrics

Track key metrics:

```go
func (w *WorkerState) getMetrics() WorkerMetrics {
    return WorkerMetrics{
        PendingSummaries:    len(w.pendingSummaries),
        OldestPendingAge:    w.getOldestPendingAge(),
        SuccessfulSummaries: w.successCount,
        FailedSummaries:     w.failureCount,
    }
}
```

### Health Checks

Implement health checks:

```go
func (w *WorkerState) isHealthy() bool {
    // Check if we have too many pending requests
    if len(w.pendingSummaries) > 1000 {
        return false
    }
    
    // Check if oldest request is too old
    if w.getOldestPendingAge() > 2*time.Hour {
        return false
    }
    
    return true
}
```

This worker serves as a comprehensive example of the command-based worker pattern and can be adapted for various background processing scenarios in your Petrock applications.
