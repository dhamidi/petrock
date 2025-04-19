package handlers

import (
	"sync"
	"time"
)

// State holds the collective in-memory state for the feature.
// This is a temporary reference to allow for handler operations
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