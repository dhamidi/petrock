package handlers

import (
	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
	"github.com/petrock/example_module_path/petrock_example_feature_name/queries"
)

// State is an alias for the feature's state type
type State = state.State

// Querier is the query handler for the feature
type Querier = queries.Querier