package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dhamidi/petrock/internal/generator"
	"github.com/spf13/cobra"
)



// NewCommandOptions holds options for command generation
type NewCommandOptions struct {
	FeatureName string
	EntityName  string
	TargetDir   string
	ModulePath  string
	Fields      []generator.CommandField
}

// NewCommandSubcommand creates the command-specific subcommand
func NewCommandSubcommand() *cobra.Command {
	return &cobra.Command{
		Use:   "command <feature>/<name-of-thing> [field1:type1] [field2:type2] ...",
		Short: "Generate a command component",
		Long: `Generate command files for a specific feature and entity from skeleton templates.

Commands handle business logic and state changes in your petrock application.
They are part of the CQRS (Command Query Responsibility Segregation) pattern.

Examples:
  petrock new command posts/create              - Generate CreateCommand for posts feature
  petrock new command posts/schedule-publication - Generate SchedulePublicationCommand for posts feature  
  petrock new command users/register            - Generate RegisterCommand for users feature  
  petrock new command orders/cancel             - Generate CancelCommand for orders feature
  petrock new command posts/publish postID:string publishAt:time.Time - Generate PublishCommand with typed fields

Generated files:
  - <feature>/commands/base.go       - Base command interfaces and types
  - <feature>/commands/register.go   - Command registration logic
  - <feature>/commands/<name-of-thing>.go   - Entity-specific command implementation`,
		Args: cobra.MinimumNArgs(1),
		RunE: runCommandGeneration,
	}
}

// runCommandGeneration handles the command generation process
func runCommandGeneration(cmd *cobra.Command, args []string) error {
	// Parse feature/entity from argument
	featureName, entityName, err := parseFeatureEntityName(args[0])
	if err != nil {
		return err
	}

	// Parse field definitions from remaining arguments
	var fields []generator.CommandField
	for _, fieldArg := range args[1:] {
		field, err := parseFieldDefinition(fieldArg)
		if err != nil {
			return fmt.Errorf("invalid field definition %q: %w", fieldArg, err)
		}
		fields = append(fields, field)
	}

	// Create command options
	options := NewCommandOptions{
		FeatureName: featureName,
		EntityName:  entityName,
		TargetDir:   ".",
		Fields:      fields,
	}

	// Detect module path
	options.ModulePath, err = detectModulePath(".")
	if err != nil {
		return fmt.Errorf("failed to detect module path: %w", err)
	}

	// Validate options
	if err := validateCommandOptions(options); err != nil {
		return fmt.Errorf("invalid command options: %w", err)
	}

	// Generate command component
	if err := generateCommand(options); err != nil {
		return fmt.Errorf("failed to generate command: %w", err)
	}

	slog.Debug("Command component generated successfully",
		"feature", options.FeatureName,
		"entity", options.EntityName)

	return nil
}

// validateCommandOptions validates the command generation options
func validateCommandOptions(options NewCommandOptions) error {
	if options.FeatureName == "" {
		return fmt.Errorf("feature name cannot be empty")
	}
	if options.EntityName == "" {
		return fmt.Errorf("entity name cannot be empty")
	}
	if options.ModulePath == "" {
		return fmt.Errorf("module path cannot be empty")
	}
	if options.TargetDir == "" {
		return fmt.Errorf("target directory cannot be empty")
	}

	return nil
}

// parseFieldDefinition parses a field definition in format "name:type"
func parseFieldDefinition(input string) (generator.CommandField, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return generator.CommandField{}, fmt.Errorf("expected format 'name:type', got %q", input)
	}
	
	name := strings.TrimSpace(parts[0])
	fieldType := strings.TrimSpace(parts[1])
	
	if name == "" {
		return generator.CommandField{}, fmt.Errorf("field name cannot be empty")
	}
	if fieldType == "" {
		return generator.CommandField{}, fmt.Errorf("field type cannot be empty")
	}
	
	// Basic validation for field names (should be valid Go identifiers)
	if !isValidGoIdentifier(name) {
		return generator.CommandField{}, fmt.Errorf("field name %q is not a valid Go identifier", name)
	}
	
	return generator.CommandField{
		Name: name,
		Type: fieldType,
	}, nil
}

// isValidGoIdentifier checks if a string is a valid Go identifier
func isValidGoIdentifier(s string) bool {
	if s == "" {
		return false
	}
	
	// First character must be letter or underscore
	if !((s[0] >= 'A' && s[0] <= 'Z') || (s[0] >= 'a' && s[0] <= 'z') || s[0] == '_') {
		return false
	}
	
	// Subsequent characters must be letters, digits, or underscores
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	
	return true
}

// generateCommand performs the actual command generation
func generateCommand(options NewCommandOptions) error {
	slog.Debug("Generating command component",
		"feature", options.FeatureName,
		"entity", options.EntityName,
		"target", options.TargetDir,
		"module", options.ModulePath,
		"fields", len(options.Fields))

	// Create command generator
	cmdGen := generator.NewCommandGenerator(".")

	// Generate command component
	return cmdGen.GenerateCommandComponentWithFields(
		options.FeatureName,
		options.EntityName,
		options.TargetDir,
		options.ModulePath,
		options.Fields,
	)
}
