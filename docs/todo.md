# Deficiencies

## feature_template/worker.go

- Ensure workers can properly track their position in the event log
- Define clear patterns for error handling and retry strategies
- Consider standardizing common worker operations (like making HTTP calls to external services)

## Design system

* Part of the core 
* probably start with DaisyUI

## Tools and MCP support 

* need to be added to introspection
* subcommand `serve` should accept protocol (http, MCP)

## Command/Query generators

## Asset pipeline

* plain JS with importmap support
* compression + hashing


## Rules

- working in a petrock generated project is still rough, as the AI is lacking context about important rules (e.g. no direct mutation of the state)
- we should automatically generate rules 
