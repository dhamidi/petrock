# Plan for posts/worker.go (Example Feature)

This file defines any long-running background processes or scheduled tasks related to the posts feature. This might be empty for simple features.

This is used to coordinate across multiple events: workers maintain internal state to make decisions.

During their work cycle they react to new events by running queries against the application and issuing commands in return.

By default workers are connected to the `app` instance's message log, but if necessary a worker can have its own message log.

The lifecycle of a worker is managed by `core/app.go`:

- workers are registered with the app,
- the app starts all workers, using one goroutine per worker,
- in that goroutine the worker's `Start()` method is called to initialize the worker,
- when the application shuts down, every worker's `Stop()` method is called.

## Operating Principle

Every worker needs to comply with the `core.Worker` interface:

- `worker.Start(context.Context)` called by `core` to initalize the worker
    - here the worker needs to load its state, usually by iterating over the event log through the provided message log
    - the worker is responsible for maintaining its position in the event log and only requesting newer entries
- `worker.Stop(context.Context)` signals the worker to stop,
    - any necessary cleanup needs to be performed here.
- `worker.Work() error` - performing the actual work:
    - running queries against the `app` to fetch data,
    - dispatching commands as a result of things happenening

## Example: Post Summarization Worker

A good example of a worker is one that summarizes posts using an external LLM service. This worker would:

1. Track posts needing summaries by monitoring events in the message log
2. Make API calls to an external service to generate summaries
3. Dispatch commands to update the post with the generated summary

This would involve the following commands:

- `RequestSummaryGeneration(requestId, postId, content)` (CommandName: `posts/request-summary-generation`)
- `FailSummaryGeneration(requestId, postId, errorMessage)` (CommandName: `posts/fail-summary-generation`)
- `SetGeneratedSummary(requestId, postId, summary)` (CommandName: `posts/set-generated-summary`)

During startup the worker would:
1. Scan the event log for `posts/request-summary-generation` events to build its internal state of posts needing summaries
2. Filter out posts that already have corresponding `posts/fail-summary-generation` or `posts/set-generated-summary` events with matching requestIds

During operation, the worker would:
1. Monitor new `posts/create` events and dispatch `posts/request-summary-generation` commands for new posts
2. Make API calls to the LLM service for pending posts
3. Dispatch either `posts/set-generated-summary` on success or `posts/fail-summary-generation` on failure

## Types

- `PostWorker`: A struct holding dependencies needed for background jobs (e.g., the application object and worker specific state).
    - `state *PostWorkerState`
    - `log *core.MessageLog` 
    - `app *core.App` 
- `PostWorkerState`: A struct holding worker specific state.
    - `pendingSummaries map[string]PendingSummary` // Map of requestId to pending summary requests
    - `lastProcessedID string` // Last message ID processed from the log

## Functions

- `NewPostWorker(app *core.App) *PostWorker`: Constructor for `PostWorker`.
- `(w *PostWorker) Start(ctx context.Context)`: Initializes the worker by scanning the message log and building internal state of pending summaries.
- `(w *PostWorker) Stop(ctx context.Context)`: Cleans up any resources and ensures pending operations are appropriately handled.
- `(w *PostWorker) Work() error`: Processes any pending summaries by calling the external API and dispatches appropriate commands.

*Note: By default, workers poll the event log once per second with a 1s random jitter, as scheduled by core/app.go.*

