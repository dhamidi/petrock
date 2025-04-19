package workers

import (
	"time"
)

// PendingSummary tracks a content item waiting for summarization
type PendingSummary struct {
	RequestID string
	ItemID    string
	Content   string
	CreatedAt time.Time
}

// State is a reference to the feature's state
type State struct {
	Items map[string]*Item
}

// Item is a representation of an item in the state
type Item struct {
	ID          string
	Name        string
	Description string
	Content     string
	Summary     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

// GetItem retrieves an item from the state
func (s *State) GetItem(id string) (*Item, bool) {
	item, found := s.Items[id]
	return item, found
}

// Command types used by the worker

// CreateCommand holds data needed to create a new entity.
type CreateCommand struct {
	Name        string
	Description string
}

// CommandName returns the unique kebab-case name for this command type.
func (c *CreateCommand) CommandName() string {
	return "petrock_example_feature_name/create"
}

// RequestSummaryGenerationCommand requests a summary be generated for an item
type RequestSummaryGenerationCommand struct {
	ID        string // ID of the item to summarize
	RequestID string // Unique ID for this summary request
}

// CommandName returns the unique kebab-case name for this command type
func (c *RequestSummaryGenerationCommand) CommandName() string {
	return "petrock_example_feature_name/request-summary-generation"
}

// FailSummaryGenerationCommand indicates a summary generation request failed
type FailSummaryGenerationCommand struct {
	ID        string // ID of the item
	RequestID string // References the original request
	Reason    string // Reason for failure
}

// CommandName returns the unique kebab-case name for this command type
func (c *FailSummaryGenerationCommand) CommandName() string {
	return "petrock_example_feature_name/fail-summary-generation"
}

// SetGeneratedSummaryCommand sets the generated summary for an item
type SetGeneratedSummaryCommand struct {
	ID        string // ID of the item
	RequestID string // References the original request
	Summary   string // The generated summary text
}

// CommandName returns the unique kebab-case name for this command type
func (c *SetGeneratedSummaryCommand) CommandName() string {
	return "petrock_example_feature_name/set-generated-summary"
}