# Petrock Architecture Overview

## System Overview

Petrock is a web application framework built around the concept of command sourcing. Unlike traditional MVC frameworks, Petrock treats user interactions as commands that represent intent. These commands are validated, stored, and then processed to update application state.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    HTTP     │     │    Form     │     │   Command   │     │    State    │
│   Request   │────▶│  Processing │────▶│  Processing │────▶│   Update    │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                                                                   │
┌─────────────┐     ┌─────────────┐     ┌─────────────┐           │
│  HTML/Turbo │     │     UI      │     │  Command    │◀───────────┘
│  Response   │◀────│  Rendering  │◀────│    Log      │
└─────────────┘     └─────────────┘     └─────────────┘
```

## Command Lifecycle

1. **Creation**: Commands are created from user input, typically via form submissions
2. **Validation**: Commands are validated against their schema
3. **Storage**: Commands are stored in the persistent log
4. **Processing**: Commands are processed to update application state
5. **Effect**: The result of command processing is used to render a response

## Application Flow

A typical request in a Petrock application follows this flow:

1. HTTP request is received by the web server
2. Router matches the request to a handler
3. Form middleware processes the request if applicable
4. Form data is validated
5. Form is converted to a command
6. Command is stored in the log
7. Command handler processes the command
8. State is updated based on the command
9. Response is rendered using gomponents
10. HTML is sent back to the browser

## Module Organization

Petrock is organized into a collection of Go modules in a single repository:

```
petrock/
├── go.work                # Workspace file for development
├── go.mod                 # Root module file
├── petrock.go             # Meta-package for importing all modules
├── cmd/                   # Command-line applications
│   └── petrock/           # Main CLI entry point
└── pkg/                   # Package directory
    ├── ui/                # UI components
    ├── core/              # Core functionality
    ├── auth/              # Authentication and user management
    ├── jobs/              # Background job system
    ├── storage/           # Storage utilities
    └── cli/               # CLI framework
```

## Module Dependencies

- All extra modules may only depend on `petrock/core`
- The meta-package at the root serves as an easy way to import common packages
