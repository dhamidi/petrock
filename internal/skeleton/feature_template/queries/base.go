package queries

import (
	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
)

// State is an alias to the state package's State type
type State = state.State

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
