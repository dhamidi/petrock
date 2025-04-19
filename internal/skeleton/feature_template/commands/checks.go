package commands

import (
	"github.com/petrock/example_module_path/core" // Placeholder for target project's core package
)

// Ensure commands implement the marker interface and Validator where applicable
var _ core.Command = (*CreateCommand)(nil)
var _ Validator = (*CreateCommand)(nil)
var _ core.Command = (*UpdateCommand)(nil)
var _ Validator = (*UpdateCommand)(nil)
var _ core.Command = (*DeleteCommand)(nil)
var _ Validator = (*DeleteCommand)(nil)
var _ core.Command = (*RequestSummaryGenerationCommand)(nil)
var _ Validator = (*RequestSummaryGenerationCommand)(nil)
var _ core.Command = (*FailSummaryGenerationCommand)(nil)
var _ core.Command = (*SetGeneratedSummaryCommand)(nil)
