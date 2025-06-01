package core

import (
	"fmt"
)

// LogFollower maintains a position in the message log
type LogFollower interface {
	// LogPosition returns the current position in the log
	LogPosition() uint64
	
	// LogSeek sets the position to a new value
	LogSeek(newPosition uint64)
}

// SimpleLogFollower implements LogFollower with basic position tracking
type SimpleLogFollower struct {
	position uint64
}

// NewLogFollower creates a new LogFollower implementation
func NewLogFollower() LogFollower {
	return &SimpleLogFollower{position: 0}
}

// LogPosition returns the current position in the log
func (f *SimpleLogFollower) LogPosition() uint64 {
	return f.position
}

// LogSeek sets the position to a new value
func (f *SimpleLogFollower) LogSeek(newPosition uint64) {
	f.position = newPosition
}

// LoadPosition loads the position from KVStore using the given key
func (f *SimpleLogFollower) LoadPosition(kvStore KVStore, key string) error {
	var position uint64
	err := kvStore.Get(key, &position)
	if err != nil {
		return fmt.Errorf("failed to load position from key %s: %w", key, err)
	}
	f.position = position
	return nil
}

// SavePosition saves the current position to KVStore using the given key
func (f *SimpleLogFollower) SavePosition(kvStore KVStore, key string) error {
	err := kvStore.Set(key, f.position)
	if err != nil {
		return fmt.Errorf("failed to save position to key %s: %w", key, err)
	}
	return nil
}
