package petrock_example_feature_name

import (
	"context"
	"errors" // Import errors package
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"                       // For generating IDs
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Executor handles command processing for the feature.
// It validates commands, logs them, and applies changes to the state.
type Executor struct {
	state *State           // Dependency on the feature's state
	log   *core.MessageLog // Dependency on the message log for persistence
}

// NewExecutor creates a new Executor instance.
func NewExecutor(state *State, log *core.MessageLog) *Executor {
	if state == nil {
		panic("state cannot be nil for Executor") // Or return an error
	}
	if log == nil {
		panic("log cannot be nil for Executor") // Or return an error
	}
	return &Executor{
		state: state,
		log:   log,
	}
}

// HandleCreate processes the CreateCommand.
// This function signature matches core.CommandHandler.
func (e *Executor) HandleCreate(ctx context.Context, command core.Command) error {
	cmd, ok := command.(CreateCommand)
	if !ok {
		return fmt.Errorf("invalid command type for HandleCreate: expected CreateCommand, got %T", command)
	}

	slog.Debug("Handling CreateCommand", "feature", "petrock_example_feature_name", "name", cmd.Name)

	// --- 1. Validation ---
	if cmd.Name == "" {
		// Consider returning a more specific validation error type
		return errors.New("item name cannot be empty")
	}
	// Add other validation rules as needed (e.g., length, format)

	// --- 2. Generate ID (if not provided by client or command) ---
	// Assuming ID is generated server-side for new items.
	// If the command should specify the ID, use cmd.ID and validate uniqueness.
	// newItemID := uuid.NewString() // Removed unused variable
	// We need to update the command *if* the state's Apply method relies on the ID being present
	// in the logged command. Alternatively, the state's Apply could generate the ID if missing.
	// For simplicity here, let's assume the ID is primarily managed by the state upon creation.
	// The logged command might not *need* the final ID if Apply handles it.
	// Let's refine this: Log the command as received, Apply generates/assigns ID.

	// --- 3. Log the Command ---
	// Log the original command before applying it to the state.
	if err := e.log.Append(ctx, cmd); err != nil {
		slog.Error("Failed to append CreateCommand to log", "error", err, "name", cmd.Name)
		return fmt.Errorf("failed to persist command: %w", err)
	}
	slog.Debug("CreateCommand appended to log", "name", cmd.Name)

	// --- 4. Apply to State ---
	// Create the item struct to be added to the state.
	// The Apply method in state.go currently uses cmd.Name as ID, let's adapt.
	// It's better if Apply uses a consistent ID source. Let's modify Apply logic conceptually.
	// For now, we'll stick to the current state.go Apply logic which uses Name as ID.
	// If state.Apply fails, the command is logged but state might be inconsistent until restart/replay.
	// This highlights a potential design choice: transactional log+state update or eventual consistency.
	// Let's assume Apply *can* fail and we should report it.
	// Note: state.Apply uses Name as ID in the current template.
	if err := e.state.Apply(cmd); err != nil {
		// Log the inconsistency. Manual intervention or a compensating command might be needed.
		slog.Error("Failed to apply CreateCommand to state after logging", "error", err, "name", cmd.Name)
		// Return the error to the caller.
		return fmt.Errorf("failed to update state after persisting command: %w", err)
	}

	slog.Info("Successfully processed CreateCommand", "feature", "petrock_example_feature_name", "id", cmd.Name) // Using Name as ID per state.go
	return nil
}

// HandleUpdate processes the UpdateCommand.
// This function signature matches core.CommandHandler.
func (e *Executor) HandleUpdate(ctx context.Context, command core.Command) error {
	cmd, ok := command.(UpdateCommand)
	if !ok {
		return fmt.Errorf("invalid command type for HandleUpdate: expected UpdateCommand, got %T", command)
	}

	slog.Debug("Handling UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// --- 1. Validation ---
	if cmd.ID == "" {
		return errors.New("item ID cannot be empty for update")
	}
	if cmd.Name == "" {
		// Allow empty description? Depends on requirements.
		return errors.New("item name cannot be empty")
	}
	// Add other validation rules...

	// --- 2. Check Existence (Optional but recommended) ---
	// Check if the item exists before logging the update command.
	// This prevents logging commands for non-existent items.
	_, found := e.state.GetItem(cmd.ID)
	if !found {
		slog.Warn("Attempted to update non-existent item", "id", cmd.ID)
		return fmt.Errorf("item with ID %s not found", cmd.ID) // Return a "not found" error
	}

	// --- 3. Log the Command ---
	if err := e.log.Append(ctx, cmd); err != nil {
		slog.Error("Failed to append UpdateCommand to log", "error", err, "id", cmd.ID)
		return fmt.Errorf("failed to persist command: %w", err)
	}
	slog.Debug("UpdateCommand appended to log", "id", cmd.ID)

	// --- 4. Apply to State ---
	if err := e.state.Apply(cmd); err != nil {
		// Log inconsistency
		slog.Error("Failed to apply UpdateCommand to state after logging", "error", err, "id", cmd.ID)
		// Return the error
		return fmt.Errorf("failed to update state after persisting command: %w", err)
	}

	slog.Info("Successfully processed UpdateCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}

// HandleDelete processes the DeleteCommand.
// This function signature matches core.CommandHandler.
func (e *Executor) HandleDelete(ctx context.Context, command core.Command) error {
	cmd, ok := command.(DeleteCommand)
	if !ok {
		return fmt.Errorf("invalid command type for HandleDelete: expected DeleteCommand, got %T", command)
	}

	slog.Debug("Handling DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)

	// --- 1. Validation ---
	if cmd.ID == "" {
		return errors.New("item ID cannot be empty for deletion")
	}

	// --- 2. Check Existence (Optional but recommended) ---
	_, found := e.state.GetItem(cmd.ID)
	if !found {
		slog.Warn("Attempted to delete non-existent item", "id", cmd.ID)
		// Decide if this is an error or idempotent success
		return fmt.Errorf("item with ID %s not found", cmd.ID) // Return error
		// return nil // Alternative: Treat deletion of non-existent item as success
	}

	// --- 3. Log the Command ---
	if err := e.log.Append(ctx, cmd); err != nil {
		slog.Error("Failed to append DeleteCommand to log", "error", err, "id", cmd.ID)
		return fmt.Errorf("failed to persist command: %w", err)
	}
	slog.Debug("DeleteCommand appended to log", "id", cmd.ID)

	// --- 4. Apply to State ---
	if err := e.state.Apply(cmd); err != nil {
		// Log inconsistency
		slog.Error("Failed to apply DeleteCommand to state after logging", "error", err, "id", cmd.ID)
		// Return the error
		return fmt.Errorf("failed to update state after persisting command: %w", err)
	}

	slog.Info("Successfully processed DeleteCommand", "feature", "petrock_example_feature_name", "id", cmd.ID)
	return nil
}

// Add more command handlers here...
