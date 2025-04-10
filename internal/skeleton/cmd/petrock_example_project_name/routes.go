package main

import (
	"net/http"
	
	"github.com/petrock/example_module_path/core"
)

// RegisterFeatureRoutes registers feature-specific HTTP routes
// This is called after features are already registered for command/query handling
func RegisterFeatureRoutes(mux *http.ServeMux, appState *AppState, executor *core.Executor) {
	// This function should call each feature's route registration function
	// Example: posts.RegisterRoutes(mux, executor, posts.NewQuerier(appState.posts))
	// 
	// Features typically need:
	// - HTTP mux to register routes
	// - Executor for command execution
	// - Feature-specific querier for data access
}