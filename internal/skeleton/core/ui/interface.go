package ui

import (
	"context"
	"io"
)

// MessageType represents the type of message being presented to the user
type MessageType int

const (
	// MessageTypeInfo represents informational messages
	MessageTypeInfo MessageType = iota
	// MessageTypeSuccess represents success messages
	MessageTypeSuccess
	// MessageTypeWarning represents warning messages
	MessageTypeWarning
	// MessageTypeError represents error messages
	MessageTypeError
	// MessageTypeProgress represents progress updates
	MessageTypeProgress
)

// ProgressState represents the current state of a progress operation
type ProgressState struct {
	// Current step being executed
	Step string
	// Current progress (0-100, or -1 for indeterminate)
	Progress int
	// Total number of steps (optional)
	Total int
	// Additional details about the current step
	Details string
}

// UI interface defines the contract for user interaction in petrock commands.
// It separates user-facing output from debug logging, allowing commands to
// be tested with mock implementations while providing consistent user experience.
type UI interface {
	// Present displays a message to the user with the specified type and formatting
	Present(ctx context.Context, msgType MessageType, message string, args ...interface{}) error

	// Prompt asks the user a question and returns their response
	Prompt(ctx context.Context, question string) (string, error)

	// ShowProgress displays progress information for long-running operations
	ShowProgress(ctx context.Context, state ProgressState) error

	// ShowError displays an error message to the user
	ShowError(ctx context.Context, err error) error

	// ShowSuccess displays a success message to the user
	ShowSuccess(ctx context.Context, message string, args ...interface{}) error
}

// OutputWriter interface abstracts the destination for UI output,
// allowing for flexible output destinations (stdout, stderr, files, etc.)
type OutputWriter interface {
	io.Writer
}
