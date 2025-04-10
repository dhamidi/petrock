# Plan for core/log.go

This file implements the persistent event/message log, backed by SQLite. It handles serialization and storage of commands or events.

## Types

- `Encoder`: An interface for encoding and decoding messages.
    - `Encode(v interface{}) ([]byte, error)`
    - `Decode(data []byte, v interface{}) error`
- `JSONEncoder`: A concrete implementation of `Encoder` using `encoding/json`.
- `Message`: A struct representing a single entry in the log.
    - `ID uint64`: Unique identifier (typically auto-incrementing primary key).
    - `Timestamp time.Time`: Time the message was logged.
    - `Type string`: A string identifier for the type of the message (e.g., "CreatePostCommand"). Used for decoding.
    - `Data []byte`: The serialized message data (e.g., JSON bytes).
- `PersistedMessage`: A struct that combines a raw message with its decoded payload.
    - `Message`: Embedded raw message struct.
    - `DecodedPayload interface{}`: The decoded Go object from the message data.
- `MessageLog`: The main struct for interacting with the log.
    - `db *sql.DB`: The database connection pool.
    - `encoder Encoder`: The encoder used for serializing/deserializing message data.
    - `typeRegistry map[string]reflect.Type`: A map from type name string to `reflect.Type`, used for decoding messages back into concrete Go types.

## Functions

- `NewMessageLog(db *sql.DB, encoder Encoder) (*MessageLog, error)`: Constructor for `MessageLog`. Initializes the type registry.
- `(l *MessageLog) RegisterType(instance interface{})`: Registers a Go type (by passing an instance, e.g., `CreatePostCommand{}`) so the log knows how to decode messages of this type string. It expects the instance to implement `core.Command` or `core.Query` and uses the name returned by `CommandName()` or `QueryName()` as the key. Stores the underlying `reflect.Type`.
- `(l *MessageLog) Append(ctx context.Context, msg interface{}) error`: Encodes the given message using the `encoder`, determines its registered name string (via `CommandName()` or `QueryName()`), and inserts a new row into the `messages` table in the database. Returns an error if the message doesn't implement a known naming interface.
- `(l *MessageLog) Version(ctx context.Context) (uint64, error)`: Returns the highest message ID in the log (the current version). Returns 0 if no messages exist.
- `(l *MessageLog) After(ctx context.Context, startID uint64) iter.Seq[PersistedMessage]`: Returns an iterator over messages after the specified version. Uses Go 1.22's `iter` package for efficient iteration without loading everything into memory.
- `(l *MessageLog) Decode(message Message) (interface{}, error)`: Decodes the `Data` field of a raw `Message` into a concrete Go command/query type. It uses the `message.Type` string to look up the `reflect.Type` in the `typeRegistry`, creates a new instance, and uses the `encoder` to deserialize the `Data` into it.
- `(l *MessageLog) setupSchema(ctx context.Context) error`: Executes SQL `CREATE TABLE IF NOT EXISTS messages (...)` to set up the necessary database table. (Marked as internal as it's called by `NewMessageLog`).

*Note on State Replay:* Application startup logic (e.g., in `cmd/serve.go`) typically involves:
1. Getting the current application version using `messageLog.Version()`.
2. Iterating through messages using `for msg := range iter.Pull(messageLog.After(ctx, lastSeenVersion))` to get new messages.
3. For each PersistedMessage, access the DecodedPayload to get the concrete command instance.
4. Looking up the corresponding *state update handler* using `commandRegistry.GetHandler(decodedCmd.CommandName())`.
5. Executing the state update handler with the decoded command to update the application state.
