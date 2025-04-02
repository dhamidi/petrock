# LogStore

The LogStore is the central repository for commands in Petrock. It provides a persistent, append-only log of all commands executed in the system, serving as the source of truth for application state.

## Interface

```go
type Message interface {
    Type() string        // Returns the message type for deserialization
    EntityID() string    // Entity this message relates to
}

type LogStore interface {
    Append(messages []Message) (newVersion uint64, error)
    Version() (uint64, error)
    ByType(typeGlob string, version uint64) (<-chan Message, error)
    ByEntity(entityGlob string, version uint64) (<-chan Message, error)
}
```

## Implementation

The default implementation uses SQLite as the storage backend:

```go
// Simplified implementation
type SQLiteLogStore struct {
    db *sql.DB
    mu sync.RWMutex
}

func NewSQLiteLogStore(dbPath string) (*SQLiteLogStore, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Create tables if they don't exist
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY,
            version INTEGER NOT NULL,
            type TEXT NOT NULL,
            entity_id TEXT NOT NULL,
            data BLOB NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_messages_type ON messages(type);
        CREATE INDEX IF NOT EXISTS idx_messages_entity_id ON messages(entity_id);
        CREATE INDEX IF NOT EXISTS idx_messages_version ON messages(version);
    `)
    
    return &SQLiteLogStore{db: db}, err
}

func (s *SQLiteLogStore) Append(messages []Message) (uint64, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    tx, err := s.db.Begin()
    if err != nil {
        return 0, err
    }
    defer tx.Rollback()
    
    // Get current version
    var currentVersion uint64
    err = tx.QueryRow("SELECT COALESCE(MAX(version), 0) FROM messages").Scan(&currentVersion)
    if err != nil {
        return 0, err
    }
    
    stmt, err := tx.Prepare("INSERT INTO messages (version, type, entity_id, data) VALUES (?, ?, ?, ?)")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()
    
    // Insert all messages with incrementing versions
    for i, msg := range messages {
        newVersion := currentVersion + uint64(i) + 1
        
        data, err := json.Marshal(msg)
        if err != nil {
            return 0, err
        }
        
        _, err = stmt.Exec(newVersion, msg.Type(), msg.EntityID(), data)
        if err != nil {
            return 0, err
        }
    }
    
    err = tx.Commit()
    if err != nil {
        return 0, err
    }
    
    return currentVersion + uint64(len(messages)), nil
}

func (s *SQLiteLogStore) Version() (uint64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    var version uint64
    err := s.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM messages").Scan(&version)
    return version, err
}

func (s *SQLiteLogStore) ByType(typeGlob string, startVersion uint64) (<-chan Message, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    rows, err := s.db.Query(
        "SELECT id, version, type, entity_id, data FROM messages WHERE type GLOB ? AND version > ? ORDER BY version ASC",
        typeGlob, startVersion,
    )
    if err != nil {
        return nil, err
    }
    
    ch := make(chan Message)
    
    go func() {
        defer rows.Close()
        defer close(ch)
        
        for rows.Next() {
            var id, version uint64
            var typ, entityID string
            var data []byte
            
            if err := rows.Scan(&id, &version, &typ, &entityID, &data); err != nil {
                // Log error but continue
                continue
            }
            
            // Deserialize message based on type
            msg := deserializeMessage(typ, data)
            if msg != nil {
                ch <- msg
            }
        }
    }()
    
    return ch, nil
}

func (s *SQLiteLogStore) ByEntity(entityGlob string, startVersion uint64) (<-chan Message, error) {
    // Similar implementation to ByType but filtering on entity_id instead
    // ...
}
```

## Command Registration

Commands need to be registered with the system for proper deserialization:

```go
var commandRegistry = make(map[string]reflect.Type)

// Register a command type
func RegisterCommand(cmd Message) {
    t := reflect.TypeOf(cmd)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    commandRegistry[cmd.Type()] = t
}

// Deserialize a message based on its type
func deserializeMessage(typ string, data []byte) Message {
    t, ok := commandRegistry[typ]
    if !ok {
        return nil
    }
    
    // Create a new instance of the command type
    v := reflect.New(t).Interface().(Message)
    
    // Unmarshal JSON data into the command
    if err := json.Unmarshal(data, v); err != nil {
        return nil
    }
    
    return v
}
```

## Usage

Example of using the LogStore to record and retrieve commands:

```go
// Record a command
cmd := CreatePost{
    ID:        "post-123",
    Title:     "Hello World",
    Content:   "My first post",
    AuthorID:  "user-456",
    PostedAt:  time.Now(),
}

// Store the command in the log
store := core.GetLogStore()
version, err := store.Append([]Message{cmd})
if err != nil {
    // Handle error
}

// Later, retrieve all post-related commands
messages, err := store.ByType("post.*", 0)
if err != nil {
    // Handle error
}

// Process messages
for msg := range messages {
    switch m := msg.(type) {
    case CreatePost:
        // Process create post command
    case UpdatePost:
        // Process update post command
    case DeletePost:
        // Process delete post command
    }
}
```
