# Deficiencies

## core/log.go
- is missing a LoadAfter method
- LoadAfter should return a Go iterator
- The iterator should return PersistedMessage objects, which correspond to database rows but with a decoded payload.

## cmd/serve.go

- the application intialization logic should move to core/app.go
- serve.go just uses the logic from core/app.go
- as for the logic: features need to be registered before replay, as deserializing log entries depends on this

## feature_template/job.go

- the skeleton is useless and idempotent execution is difficult to achieve with that as a base.
- I need to think more about a useful pattern here

## Rules

- working in a petrock generated project is still rough, as the AI is lacking context about important rules (e.g. no direct mutation of the state)
- we should automatically generate rules 