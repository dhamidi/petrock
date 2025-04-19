package state

import (
	"sort"
	"strings"
	"time"
)

// Item represents the internal state of a single entity managed by this feature.
// Adapt fields based on the specific feature's needs.
type Item struct {
	ID          string
	Name        string
	Description string
	Content     string // The main content that will be summarized
	Summary     string // The generated summary
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
	// Add other feature-specific fields here
	// IsDeleted bool // Example: for soft deletes
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
