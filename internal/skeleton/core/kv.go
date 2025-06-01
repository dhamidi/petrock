package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// KVStore provides persistent key-value storage
type KVStore interface {
	// Get retrieves a value by key and unmarshals it into dest
	Get(key string, dest any) error
	
	// Set stores a value by key, marshaling it appropriately
	Set(key string, value any) error
	
	// List returns all keys matching the glob pattern
	List(glob string) ([]string, error)
}

// SQLiteKVStore implements KVStore using SQLite storage
type SQLiteKVStore struct {
	db *sql.DB
}

// NewSQLiteKVStore creates a new SQLite-backed KV store
func NewSQLiteKVStore(db *sql.DB) (*SQLiteKVStore, error) {
	store := &SQLiteKVStore{db: db}
	if err := store.createTable(); err != nil {
		return nil, fmt.Errorf("failed to create kv_store table: %w", err)
	}
	return store, nil
}

// createTable creates the kv_store table if it doesn't exist
func (s *SQLiteKVStore) createTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS kv_store (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)
	`
	_, err := s.db.Exec(query)
	return err
}

// Get retrieves a value by key and unmarshals it into dest
func (s *SQLiteKVStore) Get(key string, dest any) error {
	var valueJSON string
	err := s.db.QueryRow("SELECT value FROM kv_store WHERE key = ?", key).Scan(&valueJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get value for key %s: %w", key, err)
	}
	
	if err := json.Unmarshal([]byte(valueJSON), dest); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}
	
	return nil
}

// Set stores a value by key, marshaling it appropriately  
func (s *SQLiteKVStore) Set(key string, value any) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}
	
	_, err = s.db.Exec("INSERT OR REPLACE INTO kv_store (key, value) VALUES (?, ?)", key, string(valueJSON))
	if err != nil {
		return fmt.Errorf("failed to set value for key %s: %w", key, err)
	}
	
	return nil
}

// List returns all keys matching the glob pattern
func (s *SQLiteKVStore) List(glob string) ([]string, error) {
	rows, err := s.db.Query("SELECT key FROM kv_store WHERE key GLOB ?", glob)
	if err != nil {
		return nil, fmt.Errorf("failed to query keys with glob %s: %w", glob, err)
	}
	defer rows.Close()
	
	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("failed to scan key: %w", err)
		}
		keys = append(keys, key)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during key iteration: %w", err)
	}
	
	return keys, nil
}
