package petrock_example_feature_name

import (
	"database/sql" // Added for db dependency
	"log/slog"
	"net/http" // Added for mux dependency

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// RegisterFeature initializes and registers the feature's handlers with the core registries
// and registers feature-specific HTTP routes.
// It connects the command/query messages to their respective handlers, registers
// message types needed for log replay, and sets up feature-specific HTTP routes.
func RegisterFeature(
	mux *http.ServeMux, // The main HTTP router
	commands *core.CommandRegistry,
	queries *core.QueryRegistry,
	messageLog *core.MessageLog, // For registering message types for decoding
	centralExecutor *core.Executor, // The central command executor
	state *State, // The feature's specific state instance
	db *sql.DB, // Shared database connection pool
	// Add other core dependencies if needed (e.g., config, external clients)
) {
	slog.Debug("Registering feature", "feature", "petrock_example_feature_name")
	
	// Validate required dependencies
	if commands == nil {
		slog.Error("Cannot register feature: CommandRegistry is nil", "feature", "petrock_example_feature_name")
		return
	}
	if queries == nil {
		slog.Error("Cannot register feature: QueryRegistry is nil", "feature", "petrock_example_feature_name")
		return
	}
	if messageLog == nil {
		slog.Error("Cannot register feature: MessageLog is nil", "feature", "petrock_example_feature_name")
		return
	}
	if centralExecutor == nil {
		slog.Error("Cannot register feature: Executor is nil", "feature", "petrock_example_feature_name")
		return
	}
	if state == nil {
		slog.Error("Cannot register feature: State is nil", "feature", "petrock_example_feature_name")
		return
	}
	// db can be optional depending on the feature

	// --- 1. Initialize Feature-Specific Logic Handlers ---
	// These components encapsulate the logic for handling commands (validation + state updates) and queries.
	// They typically depend on the feature's state and potentially other core services.

	// Assumes execute.go defines NewExecutor (the feature executor) and its handler methods.
	// Pass dependencies like state.
	featureExecutor := NewExecutor(state) // Feature executor holds state for validation

	// Assumes query.go defines NewQuerier and its handler methods.
	querier := NewQuerier(state)

	// --- 2. Initialize HTTP Handler Dependencies ---
	// Create the FeatureServer which holds dependencies needed by HTTP handlers.
	// Pass the *central* executor, querier, state, log, db, etc.
	server := NewFeatureServer(centralExecutor, querier, state, messageLog, db) // Pass centralExecutor

	// --- 3. Register Feature-Specific HTTP Routes ---
	// Call the function in routes.go to define routes on the main router.
	slog.Debug("Registering feature HTTP routes", "feature", "petrock_example_feature_name")
	RegisterRoutes(mux, server)

	// --- 4. Register Core Command Handlers ---
	// Map command message types (from messages.go) to their handler functions (from execute.go).
	// These are used by the central core.Executor.
	// Register the command type, the state update handler method, and the feature executor instance.
	slog.Debug("Registering command handlers and feature executor", "feature", "petrock_example_feature_name")
	commands.Register(CreateCommand{}, featureExecutor.HandleCreate, featureExecutor) // Pass handler AND executor instance
	commands.Register(UpdateCommand{}, featureExecutor.HandleUpdate, featureExecutor) // Pass handler AND executor instance
	commands.Register(DeleteCommand{}, featureExecutor.HandleDelete, featureExecutor) // Pass handler AND executor instance
	// Add registrations for other commands specific to this feature...

	// --- 5. Register Core Query Handlers ---
	// Map query message types (from messages.go) to their handler functions (from query.go).
	// These are used by the core /queries API endpoint.
	slog.Debug("Registering query handlers", "feature", "petrock_example_feature_name")
	queries.Register(GetQuery{}, querier.HandleGet)   // Map GetQuery to querier.HandleGet
	queries.Register(ListQuery{}, querier.HandleList) // Map ListQuery to querier.HandleList
	// Add registrations for other queries specific to this feature...

	// --- 6. Register Message Types for Decoding ---
	// Register message types (commands, events) with the MessageLog so it can
	// decode them correctly during replay.
	slog.Debug("Registering message types with MessageLog", "feature", "petrock_example_feature_name")
	RegisterTypes(messageLog) // Assumes messages.go or state.go provides RegisterTypes(*core.MessageLog)

	// --- 7. Register Background Jobs/Workers (Optional) ---
	// If the feature includes background processes (defined in jobs.go),
	// initialize them here. The actual launching (e.g., starting goroutines)
	// launching (e.g., starting goroutines) is typically done in the main application
	// entry point (e.g., cmd/serve.go) to manage their lifecycle.
	// Example:
	// jobs := NewJobs(state, messageLog)
	// // Register jobs with a scheduler or worker pool if applicable
	// // scheduler.Register(jobs.SomeScheduledTask, "*/5 * * * *")

	slog.Info("Feature registered successfully", "feature", "petrock_example_feature_name")
}
