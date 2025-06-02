package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dhamidi/petrock/internal/generator"
	"github.com/dhamidi/petrock/internal/generator/templates"
	"github.com/spf13/cobra"
)

// NewWorkerOptions holds options for worker generation
type NewWorkerOptions struct {
	FeatureName string
	EntityName  string
	TargetDir   string
	ModulePath  string
}

// NewWorkerSubcommand creates the worker-specific subcommand
func NewWorkerSubcommand() *cobra.Command {
	return &cobra.Command{
		Use:   "worker <feature>/<name-of-thing>",
		Short: "Generate a worker component",
		Long: `Generate worker files for a specific feature and entity from skeleton templates.

Workers handle background processing, side effects, and asynchronous operations
in your petrock application. They respond to commands and events to perform
tasks that don't need to be done synchronously.

Workers are ideal for:
- External API calls and integrations
- File processing and data transformation
- Email and notification sending
- Backup and synchronization tasks
- Analytics and reporting
- Cleanup and maintenance operations

Examples:
  petrock new worker posts/summary           - Generate SummaryWorker for content summarization
  petrock new worker posts/email-digest      - Generate EmailDigestWorker for sending email digests
  petrock new worker users/notification      - Generate NotificationWorker for user alerts
  petrock new worker orders/backup           - Generate BackupWorker for order data backup
  petrock new worker analytics/process       - Generate ProcessWorker for data analysis

Generated files:
  - <feature>/workers/main.go           - Worker registry and startup logic
  - <feature>/workers/types.go          - Worker state and type definitions
  - <feature>/workers/<name-of-thing>_worker.go - Entity-specific worker implementation

Common worker patterns:` + formatWorkerPatterns(),
		Args: cobra.ExactArgs(1),
		RunE: runWorkerGeneration,
	}
}

// formatWorkerPatterns formats worker patterns for help text
func formatWorkerPatterns() string {
	patterns := templates.GetWorkerPatterns()
	var formatted []string
	
	for pattern, description := range patterns {
		formatted = append(formatted, fmt.Sprintf("  - %s: %s", pattern, description))
	}
	
	return "\n" + strings.Join(formatted, "\n")
}

// runWorkerGeneration handles the worker generation process
func runWorkerGeneration(cmd *cobra.Command, args []string) error {
	// Parse feature/entity from argument
	featureName, entityName, err := parseFeatureEntityName(args[0])
	if err != nil {
		return err
	}

	// Create worker options
	options := NewWorkerOptions{
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
	if err := validateWorkerOptions(options); err != nil {
		return fmt.Errorf("invalid worker options: %w", err)
	}

	// Generate worker component
	if err := generateWorker(options); err != nil {
		return fmt.Errorf("failed to generate worker: %w", err)
	}

	slog.Info("Worker component generated successfully",
		"feature", options.FeatureName,
		"entity", options.EntityName)

	return nil
}

// validateWorkerOptions validates the worker generation options
func validateWorkerOptions(options NewWorkerOptions) error {
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

// generateWorker performs the actual worker generation
func generateWorker(options NewWorkerOptions) error {
	slog.Debug("Generating worker component",
		"feature", options.FeatureName,
		"entity", options.EntityName,
		"target", options.TargetDir,
		"module", options.ModulePath)

	// Create worker generator
	workerGen := generator.NewWorkerGenerator(".")

	// Generate worker component
	return workerGen.GenerateWorkerComponent(
		options.FeatureName,
		options.EntityName,
		options.TargetDir,
		options.ModulePath,
	)
}
