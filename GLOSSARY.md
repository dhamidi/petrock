# Petrock Glossary

## Architectural Patterns

**Command Query Responsibility Segregation (CQRS)** - An architectural pattern that separates read and write operations into distinct models, where commands modify state and queries retrieve data without side effects. In petrock, this enables optimized read models while maintaining a clear audit trail of all state changes through the command log.

**Command Sourcing** - A data persistence pattern where all user intents and operations are stored as a sequence of immutable commands in a message log, rather than storing current state directly. Petrock logs commands before execution and replays them to rebuild application state, enabling time-travel debugging, complete audit trails, and auditability of all user actions.

**Domain-Driven Design (DDD)** - A software development approach that structures code around business domains and uses domain-specific language throughout the codebase. Petrock organizes features as self-contained business domains with clear boundaries and consistent terminology.

**Registry Pattern** - A design pattern that provides a centralized location for storing and retrieving objects by name or type. Petrock uses registries to map command names to handlers, query names to handlers, and manage the lifecycle of workers and other components.

**Template System** - Petrock's code generation mechanism that uses simple string substitution to transform skeleton templates into working Go code, avoiding complex templating engines while maintaining valid, compilable code throughout the generation process.

## Core Components

**App** - The central dependency injection container that holds all application components including the database, message log, command and query registries, workers, and other shared services. It serves as the main orchestration point for the entire application.

**Command** - An interface representing state-changing operations that must implement validation logic and provide a unique command name. Commands are logged to the message log before execution, ensuring all state changes are auditable and reproducible.

**CommandHandler** - Functions that process specific command types, performing business logic and state modifications. Handlers receive validated commands and execute the necessary operations to fulfill the command's intent.

**CommandRegistry** - A centralized mapping system that associates command names with their corresponding handler functions, enabling dynamic command dispatch and registration of new command types at runtime.

**CommandWorker** - A concrete implementation of the Worker interface that processes commands from the message log in the background, enabling asynchronous command execution and event-driven architectures.

**Converter** - A pluggable component in the form processing system that transforms raw input data into specific Go types, supporting custom type conversion logic for complex data structures.

**Entity** - Domain objects that represent core business concepts within a feature, typically having identity and lifecycle management. In petrock, entities are stored within feature-specific state containers.

**Command Sourcing** - The foundational persistence strategy where all user commands are captured as immutable entries in a sequential log before execution, enabling complete state reconstruction through command replay and full auditability of user actions.

**Executor** - The command orchestration component that coordinates command validation, logging, and state updates, ensuring consistent execution flow across all command types.

**Feature** - A self-contained business domain module that encapsulates related commands, queries, workers, and state management for a specific area of functionality, enabling clean separation of concerns and independent development.

**FeatureExecutor** - An interface that provides feature-specific validation and execution logic, allowing each feature to define custom behavior while maintaining consistency with the overall command execution framework.

**FormSource** - An abstraction that normalizes different input sources (HTTP requests, JSON payloads, CLI arguments) into a consistent interface for form data processing and validation.

**Item** - Domain entities managed within a feature's state container, representing the core business objects that the application manipulates and queries.

**JSONRPCServer** - The communication protocol implementation that enables external systems to interact with petrock applications through standardized JSON-RPC 2.0 messaging.

**KVStore** - A key-value storage interface with SQLite implementation that provides persistent storage for worker state, application configuration, and other auxiliary data that doesn't belong in the event stream.

**LogFollower** - A component that tracks reading position within the message log, enabling workers and other consumers to resume processing from their last known position after restarts.

**MCPServer** - A Model Context Protocol server implementation that exposes petrock applications as AI-assistant-compatible tools, enabling intelligent automation and interaction with the application's command and query interfaces.

**MessageLog** - The central command store that persists all commands as immutable log entries, supporting type-safe serialization and deserialization through a type registry system for command sourcing.

**Parser** - The form processing component that converts raw input data into validated Go structs, handling type conversion and validation rules to ensure data integrity.

**Processing Context** - The execution environment provided to commands and workers, containing necessary dependencies and state information required for processing business logic.

**Props** - The base interface for UI component properties in petrock's gomponents-based frontend system, ensuring consistent component configuration and styling patterns.

**Query** - An interface for read-only operations that retrieve data without side effects, implementing a unique query name for registry-based dispatch and maintaining clear separation from state-modifying commands.

**QueryHandler** - Functions that process specific query types, retrieving and formatting data for client consumption without modifying application state.

**QueryRegistry** - A centralized mapping system that associates query names with their corresponding handler functions, enabling dynamic query dispatch and registration.

**QueryResult** - The interface that query responses must implement, providing a consistent contract for query return values and enabling type-safe query result handling.

**Skeleton** - The template code structure used by petrock's code generation system, consisting of valid Go code with placeholder strings that get replaced during project and feature generation.

**State** - Feature-specific in-memory data structures that are built by replaying commands from the message log, providing fast read access to current application state while maintaining command sourcing benefits.

**Validator** - A pluggable component in the form processing system that enforces business rules and data constraints, ensuring input data meets application requirements before processing.

**Worker** - Background processing components that implement the Worker interface to provide asynchronous event handling, message processing, and background task execution capabilities.

## CLI Commands

**petrock new** - The primary command for bootstrapping new Go web applications with command sourcing architecture, creating a complete project structure with dependency injection, message logging, and HTTP server setup.

**petrock feature** - Generates new feature modules within existing petrock projects, including commands, queries, workers, and UI components organized around specific business domains.

**petrock test** - Runs comprehensive integration tests that validate the entire template generation and compilation pipeline by creating sample projects and testing their functionality.

**petrock new command** - Generates command components within features, creating the necessary files and registration code for new state-changing operations.

**petrock new query** - Generates query components within features, creating the necessary files and registration code for new read-only operations.

**petrock new worker** - Generates worker components within features, creating background processing capabilities for asynchronous event handling and message processing.
