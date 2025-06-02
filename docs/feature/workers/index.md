# Workers

The `workers` directory contains background worker definitions and implementations for the feature.

## Structure

- `main.go` - Common worker interfaces and building blocks
- `summary_worker.go` - Worker for handling summary generation
- `types.go` - Shared worker type definitions

## Responsibilities

- Define background processing jobs
- Handle asynchronous operations
- Process work items from queues
- Update state based on background processing results

## Worker Pattern

Workers follow a background processing pattern where:

1. Jobs are submitted to workers through commands
2. Workers process jobs asynchronously
3. Workers report results through commands
4. Workers may retry failed operations

## Usage

Workers are typically registered during feature initialization and handle operations that are too time-consuming for synchronous processing in HTTP handlers.

