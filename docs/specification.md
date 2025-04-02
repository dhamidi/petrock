# Petrock Framework Specification

## Overview

Petrock is an opinionated Go framework for building web applications, inspired by Ruby on Rails but with a specific focus on command sourcing, deterministic behavior, and a predefined technology stack. This document outlines the architecture, components, and design principles of the Petrock framework.

## Core Principles

1. **Command Sourcing**: All state changes in a Petrock application are captured as commands (user intents) that are persisted in a log.
2. **Modularity**: Functionality is provided through modules that can be added to a project.
3. **Determinism**: Given the same sequence of commands, a Petrock application will always reach the same state.
4. **Simple Deployment**: Applications are compiled into single static binaries that include all dependencies.
5. **Developer Experience**: A comprehensive CLI facilitates all aspects of development and deployment.

## Technology Stack

Petrock uses a carefully selected technology stack:

- **Backend**: Go
- **Frontend**: TailwindCSS + gomponents
- **Interactivity**: Hotwired stack (Turbo + Stimulus.js)
- **Database**: SQLite3
- **Persistence**: Command sourcing with deterministic command handlers

## Architecture

### Module Structure

The Petrock framework is organized as a collection of Go modules in a single repository:

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

### Module Dependencies

- All extra modules may only depend on `petrock/core`
- The meta-package at the root serves as an easy way to import common packages

### Command Messaging System

The foundation of Petrock is a command messaging system that:

1. Validates commands against a schema
2. Stores commands in a persistent log
3. Processes commands to update application state
4. Provides a consistent interface for querying and manipulating data

## Core Components

### pkg/core

The core package provides foundational functionality for the entire framework:

#### Message Log

```go
// Core message interface
type Message interface {
    Type() string        // Returns the message type for deserialization
    EntityID() string    // Entity this message relates to
}

// Persistent log store
type LogStore interface {
    Append(messages []Message) (newVersion uint64, error)
    Version() (uint64, error)
    ByType(typeGlob string, version uint64) (<-chan Message, error)
    ByEntity(entityGlob string, version uint64) (<-chan Message, error)
}
```

#### Web Server

```go
// Context for request handling
type Context struct {
    Request *http.Request
    ResponseWriter http.ResponseWriter
    CurrentUser *User
    Flash map[string]string
}

// Form registration
func RegisterForm(path string, form Form, handler func(form Form, ctx *Context))

// Response helpers
func RenderError(w http.ResponseWriter, err error)
func Redirect(w http.ResponseWriter, r *http.Request, path string)
func Render(w http.ResponseWriter, component gomponents.Node)
```

#### Forms

```go
// Form interface
type Form interface {
    Validate() []ValidationError
    ToCommand(ctx *Context) (interface{}, error)
}

// Validation utilities
type ValidationError struct {
    Field string
    Message string
}

func ParseRequest(r *http.Request, form interface{}) error
```

#### Admin Interface

```go
// Admin registration functions
func RegisterAdminSection(name string, component func(ctx *Context) gomponents.Node)
func RegisterAdminAction(name string, handler func(ctx *Context) error)
```

### pkg/ui

The UI package provides components for building web interfaces:

```go
// Basic UI components
func Form(action string, method string, turboFrame string, children ...gomponents.Node) gomponents.Node
func Button(text string, attrs ...Attribute) gomponents.Node
func Input(name string, type_ string, value string, attrs ...Attribute) gomponents.Node
func ErrorMessage(message string) gomponents.Node

// Turbo-specific components
func TurboFrame(id string, children ...gomponents.Node) gomponents.Node
func TurboStream(action string, target string, content gomponents.Node) gomponents.Node
```

### pkg/auth

The authentication module manages users and sessions:

```go
// Command definitions
type RegisterUser struct {
    Username string
    Email string
    PasswordHash string
}

type Login struct {
    Username string
    PasswordHash string
    SessionID string
}

// Command handler
func Do(cmd interface{}) (interface{}, error)

// Utilities
func NewPasswordHash(password string) (string, error)
func VerifyPasswordHash(password, hash string) bool
func CurrentUser(r *http.Request) (*User, error)
```

### pkg/jobs

The jobs system handles background processing:

```go
// Job handler definition
type JobHandler func(ctx context.Context, params interface{}) error

// Job registration and execution
func RegisterJob(name string, handler JobHandler)
func EnqueueJob(name string, params interface{}, runAt time.Time) (string, error)
func CancelJob(id string) error
```

### pkg/storage

The storage package provides utilities for data persistence:

```go
// SQLite database management
func InitDatabase(path string) (*sql.DB, error)
func Backup(db *sql.DB, backupPath string) error
```

### pkg/cli

The CLI package provides commands for project management:

```go
// Project generation
func NewProject(name string) error
func NewFeature(projectPath string, featureName string) error

// Development utilities
func StartServer(projectPath string) error
func GenerateMigration(projectPath string, name string) error

// Deployment utilities
func BuildBinary(projectPath string, outputPath string) error
func DeployToServer(binaryPath string, serverConfig map[string]string) error
```

## Application Structure

A Petrock application has the following structure:

```
myapp/
├── main.go                 # Application entrypoint
├── app/
│   ├── shared/
│   │   └── ui.go           # App-specific shared components
│   └── [feature]/          # Feature modules
│       ├── messages.go     # Command/event definitions
│       ├── actions.go      # Command handlers
│       ├── state.go        # State management
│       ├── ui.go           # UI components
│       ├── routes.go       # HTTP routes
│       └── main.go         # Feature initialization
├── config/
│   └── app.go              # Application configuration
└── tmp/
    └── .gitkeep            # Temporary files
```

## Command Processing

Petrock uses a direct processing approach with transaction log:

1. Command is validated against its schema
2. Command is serialized and stored in the message log
3. Command is immediately processed by the appropriate handler
4. State changes are directly applied to the application state
5. Confirmation/response is returned to the user

```go
// Example command processing flow
func CreatePost(cmd CreatePostCommand) error {
    // Validate command
    if err := core.Validate(cmd); err != nil {
        return err
    }
    
    // Store command in message log
    messageID, err := core.Store(cmd)
    if err != nil {
        return err
    }
    
    // Process command immediately
    post := posts.NewPost(cmd.Title, cmd.Content, cmd.AuthorID)
    err = postRepository.Save(post)
    
    return err
}
```

## Form Processing

Forms are a key part of the web interface:

1. Form is rendered with the initial request
2. User submits the form
3. Form data is validated on the server
4. If validation fails, form is re-rendered with errors
5. If validation succeeds, form data is converted to a command
6. Command is processed and stored
7. User is redirected or shown a success page

```go
// Example form definition
type CreatePostForm struct {
    Title   string `form:"title"`
    Content string `form:"content"`
}

func (f CreatePostForm) Validate() []ValidationError {
    var errors []ValidationError
    
    if f.Title == "" {
        errors = append(errors, ValidationError{Field: "title", Message: "Title is required"})
    }
    
    if len(f.Content) < 10 {
        errors = append(errors, ValidationError{Field: "content", Message: "Content must be at least 10 characters"})
    }
    
    return errors
}

func (f CreatePostForm) ToCommand(ctx *core.Context) (interface{}, error) {
    return CreatePostCommand{
        Title:    f.Title,
        Content:  f.Content,
        AuthorID: ctx.CurrentUser.ID,
        PostedAt: time.Now(),
    }, nil
}
```

## Background Jobs

Jobs in Petrock follow a specific pattern:

1. Job builds state from the message log
2. Job runs in a separate goroutine
3. Before performing any action, the job catches up with the event log
4. Job performs actions based on its state
5. Job persists the fact that it performed an action in the message log

```go
// Example job definition
func ProcessImageJob(ctx context.Context, params interface{}) error {
    imageParams := params.(ImageProcessingParams)
    
    // Build state from message log
    store := core.GetLogStore()
    lastVersion, _ := store.Version()
    
    // Process image
    processed, err := processImage(imageParams.ImagePath)
    if err != nil {
        return err
    }
    
    // Record that processing was done
    cmd := ImageProcessedCommand{
        ImageID: imageParams.ImageID,
        ProcessedPath: processed,
        ProcessedAt: time.Now(),
    }
    
    _, err = store.Append([]core.Message{cmd})
    return err
}
```

## Admin Interface

The admin interface provides access to:

1. All actions in the system
2. In-memory and persisted state of feature modules
3. Event log browsing
4. Operational statistics (CPU, memory, disk usage)
5. Feature-module specific functionality

## CLI Commands

The Petrock CLI supports:

1. `petrock new [name]` - Create a new project
2. `petrock feature [name]` - Generate a new feature module
3. `petrock start` - Start the development server
4. `petrock build` - Build a production binary
5. `petrock deploy` - Deploy the application
6. `petrock db` - Database management commands
7. `petrock generate` - Generate various components

## Deployment

Petrock applications are deployed as static binaries:

1. The application is compiled into a single binary
2. The binary contains all dependencies, including SQLite
3. Cross-compilation is supported for different target platforms
4. Simple deployment via rsync or similar tools
5. SystemD service files can be generated for Linux deployment

## Extension Points

Petrock can be extended in several ways:

1. Custom feature modules
2. Custom UI components
3. Integration with external services via job handlers
4. Custom admin sections

## Implementation Guidelines

1. Use Go's standard library whenever possible
2. Minimize external dependencies
3. Follow Go best practices for package structure
4. Prioritize simplicity and readability over premature optimization
5. Ensure thorough documentation and testing

## Future Considerations

1. Plugin system for third-party extensions
2. Support for alternative databases
3. API-first development capabilities
4. Integration with cloud services
5. Monitoring and observability tools
