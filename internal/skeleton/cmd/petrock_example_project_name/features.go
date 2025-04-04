package main

import (
	"github.com/petrock/example_module_path/core"
	// petrock:import-feature - Do not remove or modify this line
)

// RegisterAllFeatures registers handlers and types for all compiled-in features.
func RegisterAllFeatures(
	commands *core.CommandRegistry,
	queries *core.QueryRegistry,
	messageLog *core.MessageLog, // Uncommented: Needed by generated feature registration
	// Add other dependencies like state if needed
) {
	// petrock:register-feature - Do not remove or modify this line
}
