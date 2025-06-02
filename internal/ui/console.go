package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// ConsoleUI implements the UI interface for console-based interaction
type ConsoleUI struct {
	stdout OutputWriter
	stderr OutputWriter
	stdin  *bufio.Scanner
}

// NewConsoleUI creates a new console UI implementation
func NewConsoleUI() *ConsoleUI {
	return &ConsoleUI{
		stdout: os.Stdout,
		stderr: os.Stderr,
		stdin:  bufio.NewScanner(os.Stdin),
	}
}

// NewConsoleUIWithWriters creates a console UI with custom writers for testing
func NewConsoleUIWithWriters(stdout, stderr OutputWriter, stdin *bufio.Scanner) *ConsoleUI {
	return &ConsoleUI{
		stdout: stdout,
		stderr: stderr,
		stdin:  stdin,
	}
}

// Present displays a message to the user with the specified type and formatting
func (c *ConsoleUI) Present(ctx context.Context, msgType MessageType, message string, args ...interface{}) error {
	formattedMessage := c.formatMessage(msgType, message, args...)

	// Write to appropriate output stream based on message type
	var writer OutputWriter
	switch msgType {
	case MessageTypeError:
		writer = c.stderr
	default:
		writer = c.stdout
	}

	_, err := fmt.Fprint(writer, formattedMessage)
	return err
}

// Prompt asks the user a question and returns their response
func (c *ConsoleUI) Prompt(ctx context.Context, question string) (string, error) {
	// Display the question
	_, err := fmt.Fprint(c.stdout, question+" ")
	if err != nil {
		return "", err
	}

	// Read user input
	if !c.stdin.Scan() {
		if err := c.stdin.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("no input received")
	}

	return strings.TrimSpace(c.stdin.Text()), nil
}

// ShowProgress displays progress information for long-running operations
func (c *ConsoleUI) ShowProgress(ctx context.Context, state ProgressState) error {
	var progressStr string

	if state.Progress >= 0 {
		// Determinate progress
		if state.Total > 0 {
			progressStr = fmt.Sprintf("[%d/%d] ", state.Progress, state.Total)
		} else {
			progressStr = fmt.Sprintf("[%d%%] ", state.Progress)
		}
	} else {
		// Indeterminate progress
		progressStr = "→ "
	}

	message := progressStr + state.Step
	if state.Details != "" {
		message += ": " + state.Details
	}

	_, err := fmt.Fprint(c.stdout, message+"\n")
	return err
}

// ShowError displays an error message to the user
func (c *ConsoleUI) ShowError(ctx context.Context, err error) error {
	return c.Present(ctx, MessageTypeError, "Error: %s\n", err.Error())
}

// ShowSuccess displays a success message to the user
func (c *ConsoleUI) ShowSuccess(ctx context.Context, message string, args ...interface{}) error {
	return c.Present(ctx, MessageTypeSuccess, message, args...)
}

// ShowHeader displays a section header to group related output
func (c *ConsoleUI) ShowHeader(ctx context.Context, title string) error {
	_, err := fmt.Fprintf(c.stdout, "\n%s\n", title)
	return err
}

// ShowFileOperation displays a file operation (create, update, etc.) to the user
func (c *ConsoleUI) ShowFileOperation(ctx context.Context, operation, filePath string) error {
	_, err := fmt.Fprintf(c.stdout, "%12s  %s\n", operation, filePath)
	return err
}

// formatMessage formats a message based on its type
func (c *ConsoleUI) formatMessage(msgType MessageType, message string, args ...interface{}) string {
	formatted := fmt.Sprintf(message, args...)

	switch msgType {
	case MessageTypeSuccess:
		return "✓ " + formatted
	case MessageTypeWarning:
		return "⚠ " + formatted
	case MessageTypeError:
		return "✗ " + formatted
	case MessageTypeProgress:
		return "→ " + formatted
	default:
		return formatted
	}
}
