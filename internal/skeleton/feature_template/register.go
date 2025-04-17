package petrock_example_feature_name

import (
	// Added for db dependency
	"log/slog"
	// Added for mux dependency
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// RegisterFeature initializes and registers the feature's handlers with the core registries
// and registers feature-specific HTTP routes.
// It connects the command/query messages to their respective handlers, registers
// message types needed for log replay, and sets up feature-specific HTTP routes.
func RegisterFeature(app *core.App, state *State) {
	slog.Debug("Registering feature", "feature", "petrock_example_feature_name")

	// Validate required dependencies
	if app == nil {
		slog.Error("Cannot register feature: App is nil", "feature", "petrock_example_feature_name")
		return
	}
	if app.CommandRegistry == nil {
		slog.Error("Cannot register feature: App.CommandRegistry is nil", "feature", "petrock_example_feature_name")
		return
	}
	if app.QueryRegistry == nil {
		slog.Error("Cannot register feature: App.QueryRegistry is nil", "feature", "petrock_example_feature_name")
		return
	}
	if app.MessageLog == nil {
		slog.Error("Cannot register feature: App.MessageLog is nil", "feature", "petrock_example_feature_name")
		return
	}
	if app.Executor == nil {
		slog.Error("Cannot register feature: App.Executor is nil", "feature", "petrock_example_feature_name")
		return
	}
	if state == nil {
		slog.Error("Cannot register feature: State is nil", "feature", "petrock_example_feature_name")
		return
	}

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
	// Pass the app instance, querier, and state
	server := NewFeatureServer(app, querier, state)

	// --- 3. Register Feature-Specific HTTP Routes ---
	// Call the function in routes.go to define routes on the main router.
	slog.Debug("Registering feature HTTP routes", "feature", "petrock_example_feature_name")
	RegisterRoutes(app, server)

	// --- 4. Register Core Command Handlers ---
	// Map command message types (from commands.go) to their handler functions (from execute.go).
	// These are used by the central core.Executor.
	// Register the command type, the state update handler method, and the feature executor instance.
	slog.Debug("Registering command handlers and feature executor", "feature", "petrock_example_feature_name")
	app.CommandRegistry.Register(CreateCommand{}, featureExecutor.HandleCreate, featureExecutor) // Pass handler AND executor instance
	app.CommandRegistry.Register(UpdateCommand{}, featureExecutor.HandleUpdate, featureExecutor) // Pass handler AND executor instance
	app.CommandRegistry.Register(DeleteCommand{}, featureExecutor.HandleDelete, featureExecutor) // Pass handler AND executor instance
	
	// Register summary-related commands
	app.CommandRegistry.Register(RequestSummaryGenerationCommand{}, featureExecutor.HandleRequestSummaryGeneration, featureExecutor)
	app.CommandRegistry.Register(FailSummaryGenerationCommand{}, featureExecutor.HandleFailSummaryGeneration, featureExecutor)
	app.CommandRegistry.Register(SetGeneratedSummaryCommand{}, featureExecutor.HandleSetGeneratedSummary, featureExecutor)

	// --- 5. Register Core Query Handlers ---
	// Map query message types (from queries.go) to their handler functions (from query.go).
	// These are used by the core /queries API endpoint.
	slog.Debug("Registering query handlers", "feature", "petrock_example_feature_name")
	app.QueryRegistry.Register(GetQuery{}, querier.HandleGet)   // Map GetQuery to querier.HandleGet
	app.QueryRegistry.Register(ListQuery{}, querier.HandleList) // Map ListQuery to querier.HandleList
	// Add registrations for other queries specific to this feature...

	// --- 6. Register Message Types for Decoding ---
	// Register message types (commands, events) with the MessageLog so it can
	// decode them correctly during replay.
	slog.Debug("Registering message types with MessageLog", "feature", "petrock_example_feature_name")
	RegisterTypes(app.MessageLog) // state.go provides RegisterTypes(*core.MessageLog)

	// --- 7. Register Worker (replacing jobs registration) ---
	// Initialize and register the worker with the app
	// The worker lifecycle (starting/stopping) is managed by the App
	slog.Debug("Registering worker", "feature", "petrock_example_feature_name")
	worker := NewWorker(app, state, app.MessageLog, app.Executor)
	app.RegisterWorker(worker)

	slog.Info("Feature registered successfully", "feature", "petrock_example_feature_name")
}
