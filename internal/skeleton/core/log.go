package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"time"
	"iter"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Encoder defines the interface for encoding and decoding messages.
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// JSONEncoder implements the Encoder interface using encoding/json.
type JSONEncoder struct{}

// Encode marshals the value v into a JSON byte slice.
func (e *JSONEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode unmarshals the JSON data into the value pointed to by v.
func (e *JSONEncoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Message represents a single entry in the persistent log.
type Message struct {
	ID        uint64
	Timestamp time.Time
	Type      string // String identifier for the concrete type of Data
	Data      []byte // Serialized message data
}

// PersistedMessage combines a raw message with its decoded payload.
type PersistedMessage struct {
	Message        // Embedded raw message struct
	DecodedPayload interface{} // The decoded Go object from the message data
}

// MessageLog provides an interface to the persistent message log backed by SQLite.
type MessageLog struct {
	db           *sql.DB
	encoder      Encoder
	typeRegistry map[string]reflect.Type
}

// NewMessageLog creates a new MessageLog instance.
// It requires a database connection and an encoder.
func NewMessageLog(db *sql.DB, encoder Encoder) (*MessageLog, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection cannot be nil")
	}
	if encoder == nil {
		return nil, fmt.Errorf("encoder cannot be nil")
	}
	log := &MessageLog{
		db:           db,
		encoder:      encoder,
		typeRegistry: make(map[string]reflect.Type),
	}
	if err := log.setupSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to setup log schema: %w", err)
	}
	return log, nil
}

// setupSchema creates the necessary 'messages' table if it doesn't exist.
func (l *MessageLog) setupSchema(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		type TEXT NOT NULL,
		data BLOB NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages (timestamp);
	`
	_, err := l.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute schema setup: %w", err)
	}
	slog.Debug("Message log schema setup complete")
	return nil
}

// RegisterType registers a Go type (by passing an instance) so the log
// can decode messages of this type later during Load.
// It uses the name returned by the CommandName() or QueryName() method as the identifier.
func (l *MessageLog) RegisterType(instance interface{}) {
	var typeName string
	var instanceType reflect.Type

	// Check if it's a Command
	if cmd, ok := instance.(Command); ok {
		typeName = cmd.CommandName()
		instanceType = reflect.TypeOf(cmd)
	} else if query, ok := instance.(Query); ok {
		// Check if it's a Query
		typeName = query.QueryName()
		instanceType = reflect.TypeOf(query)
	} else {
		// Fallback or error for types that don't implement Command/Query?
		// For now, let's log a warning and use the reflect name, although
		// ideally only Commands/Queries should be registered if they are the
		// only things expected to be logged and replayed.
		instanceType = reflect.TypeOf(instance)
		typeName = instanceType.Name() // Fallback to struct name
		slog.Warn("Registering type that is not core.Command or core.Query", "type", typeName)
		// Alternatively, return an error:
		// slog.Error("Attempted to register type that is not core.Command or core.Query", "type", reflect.TypeOf(instance).String())
		// return // Or return error
	}


	// If it's a pointer, get the element type for storage
	if instanceType.Kind() == reflect.Ptr {
		instanceType = instanceType.Elem()
	}

	if typeName == "" {
		slog.Error("Attempted to register type with empty name", "type", instanceType.String())
		return // Cannot register with an empty name
	}


	if _, exists := l.typeRegistry[typeName]; exists {
		slog.Warn("Attempted to register already registered type", "name", typeName)
		return
	}
	l.typeRegistry[typeName] = instanceType
	slog.Debug("Registered message type for decoding", "name", typeName, "type", instanceType)
}

// Append encodes the given message, determines its registered name string,
// and inserts it as a new row into the 'messages' table.
func (l *MessageLog) Append(ctx context.Context, msg interface{}) error {
	var typeName string

	// Get the registered name from the message
	if cmd, ok := msg.(Command); ok {
		typeName = cmd.CommandName()
	} else if query, ok := msg.(Query); ok {
		// Should queries be logged? Typically only commands/events are.
		// If queries *can* be logged, use their name.
		typeName = query.QueryName()
		slog.Warn("Appending a Query to the message log", "name", typeName)
	} else {
		// Fallback or error? Only log known types.
		typeName = reflect.TypeOf(msg).Name() // Fallback to struct name
		slog.Error("Attempted to append message without CommandName/QueryName", "type", typeName)
		// Return an error because we likely cannot decode this later if not registered correctly.
		return fmt.Errorf("message type %T does not implement core.Command or core.Query", msg)
	}

	if typeName == "" {
		return fmt.Errorf("cannot append message with empty registered name (type %T)", msg)
	}

	data, err := l.encoder.Encode(msg)
	if err != nil {
		return fmt.Errorf("failed to encode message type %s: %w", typeName, err)
	}

	query := `INSERT INTO messages (timestamp, type, data) VALUES (?, ?, ?)`
	_, err = l.db.ExecContext(ctx, query, time.Now().UTC(), typeName, data)
	if err != nil {
		return fmt.Errorf("failed to insert message type %s into log: %w", typeName, err)
	}

	slog.Debug("Appended message to log", "name", typeName)
	return nil
}

// Version returns the highest message ID in the log (the current version).
// Returns 0 if no messages exist.
func (l *MessageLog) Version(ctx context.Context) (uint64, error) {
	var version uint64
	query := `SELECT COALESCE(MAX(id), 0) FROM messages`
	err := l.db.QueryRowContext(ctx, query).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("failed to query log version: %w", err)
	}
	return version, nil
}

// After returns an iterator over messages after the specified version.
// Uses Go 1.22's iter package for efficient iteration without loading everything into memory.
func (l *MessageLog) After(ctx context.Context, startID uint64) iter.Seq[PersistedMessage] {
	return func(yield func(PersistedMessage) bool) {
		query := `SELECT id, timestamp, type, data FROM messages WHERE id > ? ORDER BY id ASC`
		rows, err := l.db.QueryContext(ctx, query, startID)
		if err != nil {
			slog.Error("Failed to query messages after version", "error", err, "startID", startID)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var m Message
			if err := rows.Scan(&m.ID, &m.Timestamp, &m.Type, &m.Data); err != nil {
				slog.Error("Failed to scan message row", "error", err)
				continue 
			}

			decodedPayload, err := l.Decode(m)
			if err != nil {
				slog.Error("Failed to decode message", "error", err, "type", m.Type, "id", m.ID)
				continue
			}

			pm := PersistedMessage{
				Message:        m,
				DecodedPayload: decodedPayload,
			}

			if !yield(pm) {
				break
			}
		}

		if err := rows.Err(); err != nil {
			slog.Error("Error iterating message rows", "error", err)
		}
	}
}

// Decode decodes the Data field of a raw Message into a concrete Go command/query type.
// It uses the message.Type string to look up the reflect.Type in the typeRegistry,
// creates a new instance, and uses the encoder to deserialize the Data into it.
// Returns a pointer to the decoded type.
func (l *MessageLog) Decode(message Message) (interface{}, error) {
	registeredType, exists := l.typeRegistry[message.Type]
	if !exists {
		return nil, fmt.Errorf("unknown message type: %s", message.Type)
	}

	// Create a new instance of the registered type (must be a pointer for Decode)
	newValue := reflect.New(registeredType).Interface()

	if err := l.encoder.Decode(message.Data, newValue); err != nil {
		return nil, fmt.Errorf("failed to decode message data for type %s: %w", message.Type, err)
	}

	// Return the pointer directly
	return newValue, nil
}

// --- Database Setup Helper ---

// SetupDatabase initializes the SQLite database connection and returns it.
func SetupDatabase(dataSourceName string) (*sql.DB, error) {
	// Append SQLite connection parameters to enable WAL mode and immediate transaction locking
	connString := dataSourceName + "?_journal_mode=WAL&_txlock=IMMEDIATE"
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database %s: %w", dataSourceName, err)
	}

	// Basic connection pool settings
	db.SetMaxOpenConns(1) // SQLite is often best with a single writer
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0) // Connections don't expire

	// Check the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database %s: %w", dataSourceName, err)
	}

	slog.Info("Database connection established", "path", dataSourceName)
	return db, nil
}

// --- Global Log (Optional - consider dependency injection) ---
// var Log *MessageLog

// InitLog initializes the global message log.
// Consider using dependency injection instead.
// func InitLog() {
// 	db, err := SetupDatabase("app.db") // TODO: Make path configurable
// 	if err != nil {
// 		slog.Error("Failed to setup database for global log", "error", err)
// 		panic(err) // Or handle more gracefully
// 	}
// 	// Assume JSONEncoder is suitable
// 	Log, err = NewMessageLog(db, &JSONEncoder{})
// 	if err != nil {
// 		slog.Error("Failed to create global message log", "error", err)
// 		panic(err) // Or handle more gracefully
// 	}
// 	// TODO: Register all known message types here
// 	// Log.RegisterType(SomeCommand{})
// }
