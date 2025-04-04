package petrock_example_feature_name

import (
	"log/slog"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// RegisterFeature initializes and registers the feature's handlers with the core registries.
// It connects the command/query messages to their respective handlers and registers
// message types needed for log replay.
func RegisterFeature(
	commands *core.CommandRegistry,
	queries *core.QueryRegistry,
	messageLog *core.MessageLog, // For registering message types for decoding
	state *State, // The feature's specific state instance
	// Add other core dependencies if needed (e.g., config, external clients)
) {
	slog.Debug("Registering feature", "feature", "petrock_example_feature_name")

	// --- 1. Initialize handlers/executors/queriers ---
	// These components encapsulate the logic for handling commands and queries.
	// They typically depend on the feature's state and potentially other core services.

	// Assumes execute.go defines NewExecutor and its handler methods.
	// Pass dependencies like state and messageLog (if executor needs to append events/commands).
	executor := NewExecutor(state, messageLog)

	// Assumes query.go defines NewQuerier and its handler methods.
	querier := NewQuerier(state)

	// --- 2. Register Command Handlers ---
	// Map command message types (from messages.go) to their handler functions (from execute.go).
	// The core.CommandRegistry ensures type safety during dispatch.
	slog.Debug("Registering command handlers", "feature", "petrock_example_feature_name")
	commands.Register(CreateCommand{}, executor.HandleCreate) // Map CreateCommand to executor.HandleCreate
	commands.Register(UpdateCommand{}, executor.HandleUpdate) // Map UpdateCommand to executor.HandleUpdate
	commands.Register(DeleteCommand{}, executor.HandleDelete) // Map DeleteCommand to executor.HandleDelete
	// Add registrations for other commands specific to this feature...

	// --- 3. Register Query Handlers ---
	// Map query message types (from messages.go) to their handler functions (from query.go).
	// The core.QueryRegistry ensures type safety during dispatch.
	slog.Debug("Registering query handlers", "feature", "petrock_example_feature_name")
	queries.Register(GetQuery{}, querier.HandleGet)     // Map GetQuery to querier.HandleGet
	queries.Register(ListQuery{}, querier.HandleList)   // Map ListQuery to querier.HandleList
	// Add registrations for other queries specific to this feature...

	// --- 4. Register Message Types for Decoding ---
	// If the application replays the message log to rebuild state (common in event sourcing/CQRS),
	// the core.MessageLog needs to know how to decode the stored message data back into
	// concrete Go types. Call the type registration function (conventionally in state.go).
	slog.Debug("Registering message types with MessageLog", "feature", "petrock_example_feature_name")
	RegisterTypes(messageLog) // Assumes state.go provides RegisterTypes(*core.MessageLog)

	// --- 5. Register Background Jobs/Workers (Optional) ---
	// If the feature includes background processes (defined in jobs.go),
	// they might need initialization or registration here, although the actual
	// launching (e.g., starting goroutines) is typically done in the main application
	// entry point (e.g., cmd/serve.go) to manage their lifecycle.
	// Example:
	// jobs := NewJobs(state, messageLog)
	// // Register jobs with a scheduler or worker pool if applicable
	// // scheduler.Register(jobs.SomeScheduledTask, "*/5 * * * *")

	slog.Info("Feature registered successfully", "feature", "petrock_example_feature_name")
}
