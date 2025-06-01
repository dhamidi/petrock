# Key-Value Store (KVStore)

The KVStore provides persistent key-value storage for Petrock applications. It's used internally by workers for position tracking and is available for application-specific data storage.

## Interface

```go
type KVStore interface {
    // Get retrieves a value by key and unmarshals it into dest
    Get(key string, dest any) error
    
    // Set stores a value by key, marshaling it appropriately
    Set(key string, value any) error
    
    // List returns all keys matching the glob pattern
    List(glob string) ([]string, error)
}
```

## Implementation

Petrock provides `SQLiteKVStore`, which implements the KVStore interface using SQLite for persistence:

```go
type SQLiteKVStore struct {
    db *sql.DB
}
```

### Features

- **JSON Serialization**: Values are automatically marshaled/unmarshaled as JSON
- **SQLite Storage**: Uses a dedicated `kv_store` table in your application's database
- **Glob Patterns**: List method supports SQLite's GLOB syntax for pattern matching

### Glob Pattern Examples

- `*` - All keys
- `user:*` - All keys starting with "user:"
- `*:config` - All keys ending with ":config"
- `worker:posts:*` - All keys starting with "worker:posts:"

## CLI Commands

Generated applications include CLI commands for KVStore operations:

### Get Value

```bash
go run ./cmd/myapp kv get <key>
```

Retrieves and displays a value as formatted JSON.

**Options:**
- `--db-path`: Path to SQLite database (default: "app.db")

**Example:**
```bash
go run ./cmd/myapp kv get user:123:config
```

### Set Value

```bash
go run ./cmd/myapp kv set <key> <value>
go run ./cmd/myapp kv set --json <key> <value>
```

Stores a value in the key-value store.

**Options:**
- `--db-path`: Path to SQLite database (default: "app.db")
- `--json`: Parse value as JSON instead of storing as string

**Examples:**
```bash
# Store a string value
go run ./cmd/myapp kv set user:123:name "John Doe"

# Store a JSON object
go run ./cmd/myapp kv set --json user:123:config '{"theme": "dark", "notifications": true}'

# Store a JSON array
go run ./cmd/myapp kv set --json user:123:tags '["admin", "premium"]'
```

### List Keys

```bash
go run ./cmd/myapp kv list [glob]
```

Lists all keys matching the optional glob pattern. If no pattern is provided, lists all keys.

**Options:**
- `--db-path`: Path to SQLite database (default: "app.db")

**Examples:**
```bash
# List all keys
go run ./cmd/myapp kv list

# List all user-related keys
go run ./cmd/myapp kv list "user:*"

# List all config keys
go run ./cmd/myapp kv list "*:config"
```

## Usage in Application Code

### Accessing KVStore

The KVStore is available through the App struct:

```go
func myHandler(app *core.App) {
    // Store a value
    err := app.KVStore.Set("feature:config", map[string]interface{}{
        "enabled": true,
        "version": "1.0",
    })
    
    // Retrieve a value
    var config map[string]interface{}
    err = app.KVStore.Get("feature:config", &config)
    
    // List keys with pattern
    keys, err := app.KVStore.List("feature:*")
}
```

### Common Patterns

#### Feature Configuration

```go
// Store feature configuration
config := FeatureConfig{
    Enabled: true,
    Settings: map[string]string{
        "api_key": "secret",
        "timeout": "30s",
    },
}
err := app.KVStore.Set("features:posts:config", config)

// Retrieve feature configuration
var config FeatureConfig
err := app.KVStore.Get("features:posts:config", &config)
```

#### User Preferences

```go
// Store user preferences
preferences := UserPreferences{
    Theme: "dark",
    Language: "en",
    Notifications: true,
}
err := app.KVStore.Set(fmt.Sprintf("user:%d:preferences", userID), preferences)

// List all user preference keys
keys, err := app.KVStore.List("user:*:preferences")
```

#### Worker State Persistence

```go
// Workers automatically use KVStore for position tracking
// You can also store worker-specific state
workerState := WorkerState{
    ProcessedCount: 1000,
    LastUpdate: time.Now(),
}
err := app.KVStore.Set("worker:posts:state", workerState)
```

## Database Schema

The KVStore creates the following table:

```sql
CREATE TABLE IF NOT EXISTS kv_store (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
```

Values are stored as JSON strings, allowing for complex data types while maintaining SQLite compatibility.
