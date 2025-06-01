package main

import (
	"fmt"
	"log/slog"

	"github.com/dhamidi/petrock/internal/generator"
	"github.com/spf13/cobra"
)

// NewCommandOptions holds options for command generation
type NewCommandOptions struct {
	FeatureName string
	EntityName  string
	TargetDir   string
	ModulePath  string
}

// NewCommandSubcommand creates the command-specific subcommand
func NewCommandSubcommand() *cobra.Command {
	return &cobra.Command{
		Use:   "command <feature>/<entity>",
		Short: "Generate a command component",
		Long: `Generate command files for a specific feature and entity from skeleton templates.

Commands handle business logic and state changes in your petrock application.
They are part of the CQRS (Command Query Responsibility Segregation) pattern.

Examples:
  petrock new command posts/create     - Generate CreateCommand for posts feature
  petrock new command users/register   - Generate RegisterCommand for users feature  
  petrock new command orders/cancel    - Generate CancelCommand for orders feature

Generated files:
  - <feature>/commands/base.go       - Base command interfaces and types
  - <feature>/commands/register.go   - Command registration logic
  - <feature>/commands/<entity>.go   - Entity-specific command implementation`,
		Args: cobra.ExactArgs(1),
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

	// Create command options
	options := NewCommandOptions{
		FeatureName: featureName,
		EntityName:  entityName,
		TargetDir:   ".",
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

	slog.Info("Command component generated successfully",
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

// generateCommand performs the actual command generation
func generateCommand(options NewCommandOptions) error {
	slog.Debug("Generating command component",
		"feature", options.FeatureName,
		"entity", options.EntityName,
		"target", options.TargetDir,
		"module", options.ModulePath)

	// Create command generator
	cmdGen := generator.NewCommandGenerator(".")

	// Generate command component
	return cmdGen.GenerateCommandComponent(
		options.FeatureName,
		options.EntityName,
		options.TargetDir,
		options.ModulePath,
	)
}
