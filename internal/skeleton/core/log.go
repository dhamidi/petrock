package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"time"

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
	ID        int64
	Timestamp time.Time
	Type      string // String identifier for the concrete type of Data
	Data      []byte // Serialized message data
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

// Load retrieves all messages from the database, ordered by ID.
// It attempts to decode each message's data into the registered Go type.
// Returns a slice of decoded messages (as interface{}) or an error.
func (l *MessageLog) Load(ctx context.Context) ([]interface{}, error) {
	query := `SELECT id, timestamp, type, data FROM messages ORDER BY id ASC`
	rows, err := l.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var loadedMessages []interface{}
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Timestamp, &m.Type, &m.Data); err != nil {
			slog.Error("Failed to scan message row", "error", err)
			continue // Or return error, depending on desired strictness
		}

		registeredType, exists := l.typeRegistry[m.Type]
		if !exists {
			slog.Warn("Skipping message: unknown type found in log", "type", m.Type, "id", m.ID)
			continue // Skip messages we don't know how to decode
		}

		// Create a new instance of the registered type (must be a pointer for Decode)
		newValue := reflect.New(registeredType).Interface()

		if err := l.encoder.Decode(m.Data, newValue); err != nil {
			slog.Error("Failed to decode message data", "error", err, "type", m.Type, "id", m.ID)
			continue // Or return error
		}

		// Append the decoded message (dereferenced from the pointer)
		loadedMessages = append(loadedMessages, reflect.ValueOf(newValue).Elem().Interface())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	slog.Debug("Loaded messages from log", "count", len(loadedMessages))
	return loadedMessages, nil
}

// --- Database Setup Helper ---

// SetupDatabase initializes the SQLite database connection and returns it.
func SetupDatabase(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
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
