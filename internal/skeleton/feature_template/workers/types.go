package workers

import (
	"time"

	"github.com/petrock/example_module_path/petrock_example_feature_name/commands"
	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
)

// PendingSummary tracks a content item waiting for summarization
type PendingSummary struct {
	RequestID string
	ItemID    string
	Content   string
	CreatedAt time.Time
}

// State is an alias to the state package's State type
type State = state.State

// Item is an alias to the state package's Item type
type Item = state.Item

// Command types used by the worker

// Use commands from the commands package
type CreateCommand = commands.CreateCommand
type RequestSummaryGenerationCommand = commands.RequestSummaryGenerationCommand
type FailSummaryGenerationCommand = commands.FailSummaryGenerationCommand
type SetGeneratedSummaryCommand = commands.SetGeneratedSummaryCommand
