package main

import (
	"net/http"
	
	"github.com/petrock/example_module_path/core"
)

// RegisterFeatureRoutes registers feature-specific HTTP routes
// This is called after features are already registered for command/query handling
func RegisterFeatureRoutes(app *core.App) {
	// This function should call each feature's route registration function
	// Example: posts.RegisterRoutes(app)
	// 
	// Features typically need:
	// - App to register routes
	// - App's Executor for command execution
	// - Access to feature-specific state for queries
}