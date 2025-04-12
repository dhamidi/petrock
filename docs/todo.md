# Deficiencies

## feature_template/job.go

- the skeleton is useless and idempotent execution is difficult to achieve with that as a base.
- I need to think more about a useful pattern here

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
