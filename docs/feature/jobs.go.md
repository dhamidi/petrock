# Plan for posts/jobs.go (Example Feature)

This file defines any long-running background processes or scheduled tasks related to the posts feature. This might be empty for simple features.

## Types

- `PostJobs`: A struct holding dependencies needed for background jobs (e.g., database connection, external service clients, state access).
    - `state *PostState` // Example dependency
    - `log *core.MessageLog` // Example dependency

## Functions

- `NewPostJobs(state *PostState, log *core.MessageLog) *PostJobs`: Constructor for `PostJobs`.
- `(j *PostJobs) StartContentModerationWorker(ctx context.Context)`: Example background worker function. Might poll a queue, check new posts in `j.state`, interact with an external moderation service, and potentially dispatch new commands via `j.log.Append`. Should respect the `ctx` for cancellation.
- `(j *PostJobs) StartScheduledDigestEmailer(ctx context.Context, schedule string)`: Example scheduled task. Might use a cron library to run periodically, query recent posts using `j.state`, and send emails. Should respect the `ctx` for cancellation.

*Note: The actual implementation of starting/managing these workers/schedulers often resides in `cmd/serve.go` or a dedicated worker process entry point.*
