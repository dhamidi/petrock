# Plan for core/log.go

This file implements the persistent event/message log, backed by SQLite. It handles serialization and storage of commands or events.

## Types

- `Encoder`: An interface for encoding and decoding messages.
    - `Encode(v interface{}) ([]byte, error)`
    - `Decode(data []byte, v interface{}) error`
- `JSONEncoder`: A concrete implementation of `Encoder` using `encoding/json`.
- `Message`: A struct representing a single entry in the log.
    - `ID int64`: Unique identifier (typically auto-incrementing primary key).
    - `Timestamp time.Time`: Time the message was logged.
    - `Type string`: A string identifier for the type of the message (e.g., "CreatePostCommand"). Used for decoding.
    - `Data []byte`: The serialized message data (e.g., JSON bytes).
- `MessageLog`: The main struct for interacting with the log.
    - `db *sql.DB`: The database connection pool.
    - `encoder Encoder`: The encoder used for serializing/deserializing message data.
    - `typeRegistry map[string]reflect.Type`: A map from type name string to `reflect.Type`, used for decoding messages back into concrete Go types.

## Functions

- `NewMessageLog(db *sql.DB, encoder Encoder) (*MessageLog, error)`: Constructor for `MessageLog`. Initializes the type registry.
- `(l *MessageLog) RegisterType(instance interface{})`: Registers a Go type (by passing an instance, e.g., `CreatePostCommand{}`) so the log knows how to decode messages of this type string. It expects the instance to implement `core.Command` or `core.Query` and uses the name returned by `CommandName()` or `QueryName()` as the key. Stores the underlying `reflect.Type`.
- `(l *MessageLog) Append(ctx context.Context, msg interface{}) error`: Encodes the given message using the `encoder`, determines its registered name string (via `CommandName()` or `QueryName()`), and inserts a new row into the `messages` table in the database. Returns an error if the message doesn't implement a known naming interface.
- `(l *MessageLog) Load(ctx context.Context) ([]interface{}, error)`: Loads all messages from the database table, ordered by ID. For each row, it uses the `Type` string (which should be the kebab-case name) to look up the `reflect.Type` in `typeRegistry`, creates a new instance of that type, decodes the `Data` into it using the `encoder`, and returns a slice of these decoded messages (`[]interface{}`).
- `(l *MessageLog) setupSchema(ctx context.Context) error`: Executes SQL `CREATE TABLE IF NOT EXISTS messages (...)` to set up the necessary database table. (Marked as internal as it's called by `NewMessageLog`).
