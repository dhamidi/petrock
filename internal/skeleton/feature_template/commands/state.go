package commands

import (
	"sync"
	"time"
)

// State holds the collective in-memory state for the feature.
// This is a temporary reference to allow for command validation
// It will be updated properly when migrating the state code
type State struct {
	Items map[string]*Item // Map from Item ID to the Item object pointer
	mu    sync.RWMutex     // Protects concurrent access to Items map
}

// Item represents the internal state of a single entity managed by this feature.
type Item struct {
	ID          string
	Name        string
	Description string
	Content     string // The main content that will be summarized
	Summary     string // The generated summary
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

// GetItem retrieves an item by its ID.
func (s *State) GetItem(id string) (*Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, found := s.Items[id]
	return item, found
}

// Apply is a stub method to make the commands work with the State
func (s *State) Apply(cmd interface{}, msg interface{}) error {
	return nil
}