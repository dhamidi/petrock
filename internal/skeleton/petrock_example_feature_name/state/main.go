package state

import (
	"sync"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// State holds the collective in-memory state for the feature.
// It's built by replaying logged messages and updated by command handlers.
type State struct {
	Items map[string]*Item // Map from Item ID to the Item object pointer
	mu    sync.RWMutex     // Protects concurrent access to Items map
}

// NewState creates an initialized (empty) State.
func NewState() *State {
	return &State{
		Items: make(map[string]*Item),
	}
}

// getTimestamp returns the timestamp from the message metadata if available, otherwise current time
func getTimestamp(msg *core.Message) time.Time {
	if msg != nil {
		return msg.Timestamp
	}
	return time.Now().UTC()
}
