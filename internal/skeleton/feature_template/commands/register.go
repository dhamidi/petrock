package commands

import (
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// RegisterTypes registers the message types used by this feature.
// This should be called during application initialization where the message log is configured.
func RegisterTypes(log *core.MessageLog) {
	// Register command types as pointers since CommandName has pointer receivers
	log.RegisterType(&CreateCommand{})
	log.RegisterType(&UpdateCommand{})
	log.RegisterType(&DeleteCommand{})
	log.RegisterType(&RequestSummaryGenerationCommand{})
	log.RegisterType(&FailSummaryGenerationCommand{})
	log.RegisterType(&SetGeneratedSummaryCommand{})
}