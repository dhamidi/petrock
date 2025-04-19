# State

The `state` directory contains the state management for the feature, including data models and state handling logic.

## Structure

- `main.go` - Main state container and interfaces
- `item.go` - Core item state definition and management
- `metadata.go` - Related metadata state

## Responsibilities

- Define the data model for the feature
- Manage state transitions
- Apply commands to modify state
- Provide access to current state for queries

## State Management Pattern

The state follows an event-sourced pattern where:

1. Commands generate events
2. Events are applied to the state
3. The current state is derived from the sequence of events
4. The state is the single source of truth for queries

## Usage

State is typically not accessed directly from handlers but instead through commands (for modifications) and queries (for reads).