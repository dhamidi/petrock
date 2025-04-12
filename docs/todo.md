# Deficiencies

## feature_template/job.go

- the skeleton is useless and idempotent execution is difficult to achieve with that as a base.
- I need to think more about a useful pattern here

## feature_template/messages.go

Currently all messages are defined in `messages.go`.

In the first iteration, we should split `messages.go` into `queries.go` and `commands.go` to make the directory structure more legible.

Which designd documents in `docs/` need to change for this?  Once you have identified them, list what changes need to be made in each document.

## Rules

- working in a petrock generated project is still rough, as the AI is lacking context about important rules (e.g. no direct mutation of the state)
- we should automatically generate rules 
