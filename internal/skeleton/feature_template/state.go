package petrock_example_feature_name

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Item represents the internal state of a single entity managed by this feature.
// Adapt fields based on the specific feature's needs.
type Item struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
	// Add other feature-specific fields here
	// IsDeleted bool // Example: for soft deletes
}

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

// Apply updates the state based on a logged message (typically a command).
// This is the core logic for state reconstruction during replay and updates during runtime.
// Note: This example assumes commands are logged directly. If events are logged,
// the logic here would react to event types instead.
// The msg parameter is non-nil during replay (providing timestamp and ID) and nil during direct execution.
func (s *State) Apply(payload interface{}, msg *core.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch cmd := payload.(type) {
	case CreateCommand:
		// Check if item already exists (optional, depends on desired idempotency)
		if _, exists := s.Items[cmd.Name]; exists { // Assuming Name is used as ID for creation simplicity
			slog.Warn("Attempted to create already existing item", "id", cmd.Name)
			return fmt.Errorf("item with name %s already exists", cmd.Name) // Or return nil for idempotency
		}
		newItem := &Item{
			ID:          cmd.Name, // Use Name as ID for simplicity, replace with generated ID if needed
			Name:        cmd.Name,
			Description: cmd.Description,
			CreatedAt:   getTimestamp(msg), // Use message timestamp if available, otherwise current time
			UpdatedAt:   getTimestamp(msg),
			Version:     1,
		}
		s.Items[newItem.ID] = newItem
		slog.Debug("Applied CreateCommand to state", "id", newItem.ID)

	case UpdateCommand:
		existingItem, found := s.Items[cmd.ID]
		if !found {
			slog.Warn("Attempted to update non-existent item", "id", cmd.ID)
			return fmt.Errorf("item with ID %s not found for update", cmd.ID) // Or handle as upsert
		}
		existingItem.Name = cmd.Name
		existingItem.Description = cmd.Description
		existingItem.UpdatedAt = getTimestamp(msg) // Use message timestamp if available, otherwise current time
		existingItem.Version++
		slog.Debug("Applied UpdateCommand to state", "id", existingItem.ID, "version", existingItem.Version)

	case DeleteCommand:
		if _, found := s.Items[cmd.ID]; !found {
			slog.Warn("Attempted to delete non-existent item", "id", cmd.ID)
			return fmt.Errorf("item with ID %s not found for deletion", cmd.ID) // Or return nil for idempotency
		}
		delete(s.Items, cmd.ID) // Hard delete for this example
		slog.Debug("Applied DeleteCommand to state", "id", cmd.ID)

	default:
		// This might happen if the message log contains message types not handled here
		// (e.g., events if event sourcing is used, or other command types).
		slog.Warn("Apply received unhandled message type", "type", fmt.Sprintf("%T", msg))
		// Decide whether to return an error or just ignore
		// return fmt.Errorf("unhandled message type in State.Apply: %T", msg)
	}

	return nil
}

// GetItem retrieves an item by its ID.
// Returns the item pointer and true if found, nil and false otherwise.
func (s *State) GetItem(id string) (*Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, found := s.Items[id]
	// Return a copy to prevent external modification? Depends on usage patterns.
	// For now, returning the pointer for efficiency. Be mindful of callers modifying it.
	return item, found
}

// ListItems retrieves a slice of items, applying filtering and pagination.
// Returns the slice of items for the current page and the total count of matching items.
// Note: Basic filtering and pagination implemented here. More complex queries might need optimization.
func (s *State) ListItems(page, pageSize int, filter string) ([]*Item, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filteredItems []*Item
	lowerFilter := strings.ToLower(filter)

	// Filter items
	for _, item := range s.Items {
		// Example filter: checks Name and Description (case-insensitive)
		if filter == "" || strings.Contains(strings.ToLower(item.Name), lowerFilter) || strings.Contains(strings.ToLower(item.Description), lowerFilter) {
			filteredItems = append(filteredItems, item)
		}
	}

	totalCount := len(filteredItems)

	// Sort items (e.g., by CreatedAt descending for recent items first)
	sort.Slice(filteredItems, func(i, j int) bool {
		return filteredItems[i].CreatedAt.After(filteredItems[j].CreatedAt)
	})

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start < 0 {
		start = 0
	}
	if start >= totalCount {
		return []*Item{}, totalCount // Page out of bounds
	}
	if end > totalCount {
		end = totalCount
	}

	// Return copies? See GetItem comment.
	return filteredItems[start:end], totalCount
}

// --- Direct State Modifiers (used internally by Apply or potentially command handlers if not logging) ---

// AddItem adds a new item directly to the state map. USE WITH CAUTION outside Apply.
// Primarily useful if command handlers modify state directly *before* logging.
func (s *State) AddItem(item *Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Items[item.ID]; exists {
		return fmt.Errorf("item with ID %s already exists", item.ID)
	}
	s.Items[item.ID] = item
	return nil
}

// UpdateItem updates an existing item directly in the state map. USE WITH CAUTION outside Apply.
func (s *State) UpdateItem(item *Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Items[item.ID]; !exists {
		return fmt.Errorf("item with ID %s not found for update", item.ID)
	}
	// Consider version checking here if needed
	s.Items[item.ID] = item // Replace existing pointer
	return nil
}

// DeleteItem removes an item directly from the state map. USE WITH CAUTION outside Apply.
func (s *State) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Items[id]; !exists {
		return fmt.Errorf("item with ID %s not found for deletion", id)
	}
	delete(s.Items, id)
	return nil
}

// getTimestamp returns the timestamp from the message metadata if available, otherwise current time
func getTimestamp(msg *core.Message) time.Time {
	if msg != nil {
		return msg.Timestamp
	}
	return time.Now().UTC()
}

// --- Helper to register types with the message log ---

// RegisterTypes registers the message types used by this feature's state Apply method.
// This should be called during application initialization where the message log is configured.
func RegisterTypes(log *core.MessageLog) {
	log.RegisterType(CreateCommand{})
	log.RegisterType(UpdateCommand{})
	log.RegisterType(DeleteCommand{})
	// Register any event types here if using event sourcing and Apply reacts to events
}
