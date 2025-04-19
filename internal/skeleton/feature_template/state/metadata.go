package state

import (
	"fmt"
)

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