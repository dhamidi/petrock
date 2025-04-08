package petrock_example_feature_name

import (
	"database/sql" // Added for db dependency
	"log/slog"
	"net/http" // Added for mux dependency

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// RegisterFeature initializes and registers the feature's handlers with the core registries
// and registers feature-specific HTTP routes.
// It connects the command/query messages to their respective handlers and registers
// message types needed for log replay.
func RegisterFeature(
	mux *http.ServeMux, // The main HTTP router
	commands *core.CommandRegistry,
	queries *core.QueryRegistry,
	messageLog *core.MessageLog, // For registering message types for decoding
	state *State, // The feature's specific state instance
	db *sql.DB, // Shared database connection pool
	coreExecutor core.Executor, // Centralized executor for standardized command handling
	// Add other core dependencies if needed (e.g., config, external clients)
) {
	slog.Debug("Registering feature", "feature", "petrock_example_feature_name")

	// --- 1. Initialize Core Logic Handlers (FeatureExecutor, Querier) ---
	// These components encapsulate the logic for handling commands and queries.
	// They typically depend on the feature's state and potentially other core services.

	// Assumes execute.go defines NewFeatureExecutor and its handler methods.
	// FeatureExecutor only depends on state, no longer needs MessageLog (handled by core.Executor)
	featureExecutor := NewFeatureExecutor(state)

	// Assumes query.go defines NewQuerier and its handler methods.
	querier := NewQuerier(state)

	// --- 2. Initialize HTTP Handler Dependencies ---
	// Create the FeatureServer which holds dependencies needed by HTTP handlers.
	// Pass all necessary components (featureExecutor, querier, state, executor, commands, db, etc.).
	server := NewFeatureServer(featureExecutor, querier, state, coreExecutor, commands, db)

	// --- 3. Register Feature-Specific HTTP Routes ---
	// Call the function in routes.go to define routes on the main router.
	slog.Debug("Registering feature HTTP routes", "feature", "petrock_example_feature_name")
	RegisterRoutes(mux, server)

	// --- 4. Register Core Command Handlers ---
	// Map command message types (from messages.go) to their handler functions (from execute.go).
	// These are used by the core /commands API endpoint.
	slog.Debug("Registering command handlers", "feature", "petrock_example_feature_name")
	commands.Register(CreateCommand{}, featureExecutor.HandleCreate) // Map CreateCommand to featureExecutor.HandleCreate
	commands.Register(UpdateCommand{}, featureExecutor.HandleUpdate) // Map UpdateCommand to featureExecutor.HandleUpdate
	commands.Register(DeleteCommand{}, featureExecutor.HandleDelete) // Map DeleteCommand to featureExecutor.HandleDelete
	// Add registrations for other commands specific to this feature...

	// --- 5. Register Core Query Handlers ---
	// Map query message types (from messages.go) to their handler functions (from query.go).
	// These are used by the core /queries API endpoint.
	slog.Debug("Registering query handlers", "feature", "petrock_example_feature_name")
	queries.Register(GetQuery{}, querier.HandleGet)     // Map GetQuery to querier.HandleGet
	queries.Register(ListQuery{}, querier.HandleList)   // Map ListQuery to querier.HandleList
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
