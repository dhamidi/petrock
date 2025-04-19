package queries

import (
	"time"
	
	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
)

// State is an alias to the state package's State type
type State = state.State

// ItemResult holds the data for a single entity returned by a query.
type ItemResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"` // Content field that gets summarized
	Summary     string    `json:"summary"` // Generated summary
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version"` // Example version field
}

// Querier handles query processing for the feature.
// It typically depends on the feature's state representation.
type Querier struct {
	state *State // Dependency on the feature's state
}

// NewQuerier creates a new Querier instance.
func NewQuerier(state *State) *Querier {
	return &Querier{
		state: state,
	}
}
