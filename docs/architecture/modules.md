# Modules

Modules are the building blocks of a Petrock application. Each module encapsulates a specific feature or functionality, such as authentication, file storage, or background jobs.

## Module Interface

```go
// Module interface
type Module interface {
    Name() string
    Commands() []Message
    Init(core *Core) error
    RegisterRoutes(router *Router)
    RegisterAdminSections(admin *Admin)
    Shutdown() error
}

// CommandExecutor is implemented by modules that can execute commands
type CommandExecutor interface {
    Execute(cmd interface{}) (interface{}, error)
}

// QueryExecutor is implemented by modules that can execute queries
type QueryExecutor interface {
    Query(query interface{}) (interface{}, error)
}
```

## Module Registry

Modules are registered with the core system:

```go
var moduleRegistry = make(map[string]Module)

// Register a module
func RegisterModule(module Module) {
    moduleRegistry[module.Name()] = module
}

// Get a module by name
func GetModule(name string) (Module, bool) {
    module, ok := moduleRegistry[name]
    return module, ok
}

// Initialize all modules
func InitModules(core *Core) error {
    for _, module := range moduleRegistry {
        if err := module.Init(core); err != nil {
            return err
        }
    }
    return nil
}

// Shutdown all modules
func ShutdownModules() error {
    for _, module := range moduleRegistry {
        if err := module.Shutdown(); err != nil {
            return err
        }
    }
    return nil
}
```

## Module Communication

Modules can communicate through commands and queries:

```go
// Execute a command on a module
func ExecuteCommand(module string, cmd interface{}) (interface{}, error) {
    mod, ok := GetModule(module)
    if !ok {
        return nil, fmt.Errorf("module %s not found", module)
    }
    
    executor, ok := mod.(CommandExecutor)
    if !ok {
        return nil, fmt.Errorf("module %s does not support command execution", module)
    }
    
    return executor.Execute(cmd)
}

// Execute a query on a module
func ExecuteQuery(module string, query interface{}) (interface{}, error) {
    mod, ok := GetModule(module)
    if !ok {
        return nil, fmt.Errorf("module %s not found", module)
    }
    
    querier, ok := mod.(QueryExecutor)
    if !ok {
        return nil, fmt.Errorf("module %s does not support queries", module)
    }
    
    return querier.Query(query)
}
```

## Module Lifecycle

A module's lifecycle includes:

1. **Registration**: Module is registered with the core
2. **Initialization**: Module's `Init` method is called
3. **Route Registration**: Module registers its HTTP routes
4. **Command Registration**: Module registers its command types
5. **Operation**: Module processes commands and queries
6. **Shutdown**: Module's `Shutdown` method is called when the application stops

## Auth Module Example

Here's an example of the auth module structure:

```go
// Auth module
type AuthModule struct {
    core *Core
    
    // State
    users  map[string]User
    tokens map[string]Token
    
    // Version tracking
    version uint64
    mu      sync.RWMutex
}

func NewAuthModule() *AuthModule {
    return &AuthModule{
        users:  make(map[string]User),
        tokens: make(map[string]Token),
    }
}

func (m *AuthModule) Name() string {
    return "auth"
}

func (m *AuthModule) Commands() []Message {
    return []Message{
        RegisterUser{},
        Login{},
        Logout{},
    }
}

func (m *AuthModule) Init(core *Core) error {
    // Implementation...
    return nil
}

func (m *AuthModule) RegisterRoutes(router *Router) {
    // Implementation...
}

func (m *AuthModule) RegisterAdminSections(admin *Admin) {
    // Implementation...
}

func (m *AuthModule) Shutdown() error {
    // Implementation...
    return nil
}

func (m *AuthModule) Execute(cmd interface{}) (interface{}, error) {
    // Implementation...
    return nil, nil
}

func (m *AuthModule) Query(query interface{}) (interface{}, error) {
    // Implementation...
    return nil, nil
}
```

## Auth Message Types

```go
// RegisterUser command
type RegisterUser struct {
    Username     string
    Email        string
    PasswordHash string
    CreatedAt    time.Time
}

func (c RegisterUser) Type() string {
    return "auth.register_user"
}

func (c RegisterUser) EntityID() string {
    return c.Username
}

// Login command
type Login struct {
    Username     string
    PasswordHash string
    SessionID    string
    LoginAt      time.Time
}

func (c Login) Type() string {
    return "auth.login"
}

func (c Login) EntityID() string {
    return c.Username
}

// Logout command
type Logout struct {
    SessionID string
    LogoutAt  time.Time
}

func (c Logout) Type() string {
    return "auth.logout"
}

func (c Logout) EntityID() string {
    return c.SessionID
}
```

## Storage Module Example

```go
// Storage module
type StorageModule struct {
    core *Core
    
    // State
    files map[string]FileInfo
    
    // Version tracking
    version uint64
    mu      sync.RWMutex
}

func NewStorageModule() *StorageModule {
    return &StorageModule{
        files: make(map[string]FileInfo),
    }
}

func (m *StorageModule) Name() string {
    return "storage"
}

func (m *StorageModule) Commands() []Message {
    return []Message{
        UploadFile{},
        DeleteFile{},
    }
}

func (m *StorageModule) Init(core *Core) error {
    // Implementation...
    return nil
}

func (m *StorageModule) RegisterRoutes(router *Router) {
    // Implementation...
}

func (m *StorageModule) RegisterAdminSections(admin *Admin) {
    // Implementation...
}

func (m *StorageModule) Shutdown() error {
    // Implementation...
    return nil
}

func (m *StorageModule) Execute(cmd interface{}) (interface{}, error) {
    // Implementation...
    return nil, nil
}

func (m *StorageModule) Query(query interface{}) (interface{}, error) {
    // Implementation...
    return nil, nil
}
```

## Storage Message Types

```go
// UploadFile command
type UploadFile struct {
    ID          string
    Name        string
    ContentType string
    Size        int64
    Path        string
    UploadedAt  time.Time
}

func (c UploadFile) Type() string {
    return "storage.upload_file"
}

func (c UploadFile) EntityID() string {
    return c.ID
}

// DeleteFile command
type DeleteFile struct {
    ID        string
    DeletedAt time.Time
}

func (c DeleteFile) Type() string {
    return "storage.delete_file"
}

func (c DeleteFile) EntityID() string {
    return c.ID
}
```

## Application Module Organization

In a Petrock application, modules are organized in feature-specific directories:

```
app/
├── auth/
│   ├── messages.go    # Command definitions
│   ├── actions.go     # Command handlers
│   ├── state.go       # State management
│   ├── ui.go          # UI components
│   ├── routes.go      # HTTP routes
│   └── main.go        # Module initialization
├── storage/
│   ├── messages.go
│   ├── actions.go
│   ├── state.go
│   ├── ui.go
│   ├── routes.go
│   └── main.go
└── other_modules/
    ├── ...
```

Each feature module follows the same structure:

1. **messages.go**: Defines command types and their serialization
2. **actions.go**: Implements command handlers
3. **state.go**: Manages in-memory state and rebuilding from event log
4. **ui.go**: UI components specific to the module
5. **routes.go**: HTTP route registration
6. **main.go**: Module initialization and registration

## Module Dependencies

Modules can only depend on the core package:

```go
import (
    "github.com/yourusername/petrock/pkg/core"
)
```

This ensures a clean dependency graph and prevents circular dependencies.

## Module Communication Pattern

Modules communicate with each other using the command pattern:

```go
// In auth module
func (m *AuthModule) IsAuthorized(userID string, resource string) bool {
    query := AuthorizationQuery{
        UserID:   userID,
        Resource: resource,
    }
    
    result, err := core.ExecuteQuery("auth", query)
    if err != nil {
        return false
    }
    
    authorized, ok := result.(bool)
    if !ok {
        return false
    }
    
    return authorized
}

// In another module that needs auth
func (m *SomeModule) ProtectedAction(ctx *Context) error {
    // Get current user from context
    user := ctx.CurrentUser
    
    // Check authorization
    authorized, err := core.ExecuteQuery("auth", AuthorizationQuery{
        UserID:   user.ID,
        Resource: "some_resource",
    })
    
    if err != nil || !authorized.(bool) {
        return errors.New("unauthorized")
    }
    
    // Perform protected action
    return nil
}
```
