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

During Sta
## Types

- `PostWorker`: A struct holding dependencies needed for background jobs (e.g., the application object and worker specific state).
    - `state *PostWorkerState`
    - `log *core.MessageLog` 
    - `app *core.App` 
- `PostWorkerState`: A struct holding worker specific state.
    - `posts []*Post` // Example state
    - `lastProcessedID string` // Example state

## Functions

- `NewPostWorker(app *core.App, state *PostWorkerState) *PostWorker`: Constructor for `PostWorker`.
- `(w *PostWorker) Start(ctx context.Context)`: Example background worker function. Might poll a queue, check new posts in `w.state`, interact with an external moderation service, and potentially dispatch new commands via `w.app.Log`. Should respect the `ctx` for cancellation.

*Note: The actual implementation of starting/managing these workers/schedulers often resides in `cmd/serve.go` or a dedicated worker process entry point.*

