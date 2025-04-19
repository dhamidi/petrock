package petrock_example_feature_name

import (
	"log/slog"

	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
	"github.com/petrock/example_module_path/petrock_example_feature_name/commands"
	"github.com/petrock/example_module_path/petrock_example_feature_name/handlers"
	"github.com/petrock/example_module_path/petrock_example_feature_name/queries"
	"github.com/petrock/example_module_path/petrock_example_feature_name/routes"
	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
	"github.com/petrock/example_module_path/petrock_example_feature_name/workers"
)

// RegisterFeature initializes and registers the feature's handlers with the core registries
// and registers feature-specific HTTP routes.
// It connects the command/query messages to their respective handlers, registers
// message types needed for log replay, and sets up feature-specific HTTP routes.
func RegisterFeature(app *core.App, featureState *state.State) {
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
	if featureState == nil {
		slog.Error("Cannot register feature: State is nil", "feature", "petrock_example_feature_name")
		return
	}

	// --- 1. Initialize Feature-Specific Logic Handlers ---
	// These components encapsulate the logic for handling commands (validation + state updates) and queries.
	// They typically depend on the feature's state and potentially other core services.

	// Create the command executor for handling commands
	featureExecutor := commands.NewExecutor(featureState)

	// Create the querier for handling queries
	featureQuerier := queries.NewQuerier(featureState)

	// --- 2. Initialize HTTP Handler Dependencies ---
	// Create the FeatureServer which holds dependencies needed by HTTP handlers.
	server := handlers.NewFeatureServer(app, featureQuerier, featureState)

	// --- 3. Register Feature-Specific HTTP Routes ---
	// Call the function in routes.go to define routes on the main router.
	slog.Debug("Registering feature HTTP routes", "feature", "petrock_example_feature_name")
	routes.RegisterRoutes(app, server)

	// --- 4. Register Core Command Handlers ---
	// Map command message types to their handler functions
	slog.Debug("Registering command handlers and feature executor", "feature", "petrock_example_feature_name")
	app.CommandRegistry.Register(&commands.CreateCommand{}, featureExecutor.HandleCreate, featureExecutor)
	app.CommandRegistry.Register(&commands.UpdateCommand{}, featureExecutor.HandleUpdate, featureExecutor)
	app.CommandRegistry.Register(&commands.DeleteCommand{}, featureExecutor.HandleDelete, featureExecutor)
	
	// Register summary-related commands
	app.CommandRegistry.Register(&commands.RequestSummaryGenerationCommand{}, featureExecutor.HandleRequestSummaryGeneration, featureExecutor)
	app.CommandRegistry.Register(&commands.FailSummaryGenerationCommand{}, featureExecutor.HandleFailSummaryGeneration, featureExecutor)
	app.CommandRegistry.Register(&commands.SetGeneratedSummaryCommand{}, featureExecutor.HandleSetGeneratedSummary, featureExecutor)

	// --- 5. Register Core Query Handlers ---
	// Map query message types to their handler functions
	slog.Debug("Registering query handlers", "feature", "petrock_example_feature_name")
	app.QueryRegistry.Register(queries.GetQuery{}, featureQuerier.HandleGet)
	app.QueryRegistry.Register(queries.ListQuery{}, featureQuerier.HandleList)

	// --- 6. Register Message Types for Decoding ---
	// Register message types (commands, events) with the MessageLog
	slog.Debug("Registering message types with MessageLog", "feature", "petrock_example_feature_name")
	state.RegisterTypes(app.MessageLog)

	// --- 7. Register Worker ---
	// Initialize and register the worker with the app
	slog.Debug("Registering worker", "feature", "petrock_example_feature_name")
	worker := workers.NewWorker(app, featureState, app.MessageLog, app.Executor)
	app.RegisterWorker(worker)

	slog.Info("Feature registered successfully", "feature", "petrock_example_feature_name")
}

// NewState creates a new instance of the feature's state
func NewState() *state.State {
	return state.NewState()
}