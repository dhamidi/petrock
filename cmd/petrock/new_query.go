package main

import (
	"fmt"
	"log/slog"

	"github.com/dhamidi/petrock/internal/generator"
	"github.com/spf13/cobra"
)

// NewQueryOptions holds options for query generation
type NewQueryOptions struct {
	FeatureName string
	EntityName  string
	TargetDir   string
	ModulePath  string
}

// NewQuerySubcommand creates the query-specific subcommand
func NewQuerySubcommand() *cobra.Command {
	return &cobra.Command{
		Use:   "query <feature>/<name-of-thing>",
		Short: "Generate a query component",
		Long: `Generate query files for a specific feature and entity from skeleton templates.

Queries handle data retrieval and read operations in your petrock application.
They are part of the CQRS (Command Query Responsibility Segregation) pattern.

Queries are read-only operations that return data without modifying application state.
They can include filtering, pagination, sorting, and aggregation logic.

Examples:
  petrock new query posts/get               - Generate GetQuery for retrieving a single post
  petrock new query posts/search-published  - Generate SearchPublishedQuery for searching published posts
  petrock new query users/list              - Generate ListQuery for retrieving user lists
  petrock new query orders/search           - Generate SearchQuery for searching orders
  petrock new query analytics/count         - Generate CountQuery for counting entities

Generated files:
  - <feature>/queries/base.go     - Base query interfaces and types
  - <feature>/queries/<name-of-thing>.go - Entity-specific query implementation

Common query patterns:
  - get: Retrieve a single entity by ID
  - list: Retrieve multiple entities with pagination
  - search: Find entities matching criteria
  - find: Locate entities with complex filters
  - count: Count entities matching conditions`,
		Args: cobra.ExactArgs(1),
		RunE: runQueryGeneration,
	}
}

// runQueryGeneration handles the query generation process
func runQueryGeneration(cmd *cobra.Command, args []string) error {
	// Parse feature/entity from argument
	featureName, entityName, err := parseFeatureEntityName(args[0])
	if err != nil {
		return err
	}

	// Create query options
	options := NewQueryOptions{
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
	if err := validateQueryOptions(options); err != nil {
		return fmt.Errorf("invalid query options: %w", err)
	}

	// Generate query component
	if err := generateQuery(options); err != nil {
		return fmt.Errorf("failed to generate query: %w", err)
	}

	slog.Info("Query component generated successfully",
		"feature", options.FeatureName,
		"entity", options.EntityName)

	return nil
}

// validateQueryOptions validates the query generation options
func validateQueryOptions(options NewQueryOptions) error {
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

// generateQuery performs the actual query generation
func generateQuery(options NewQueryOptions) error {
	slog.Debug("Generating query component",
		"feature", options.FeatureName,
		"entity", options.EntityName,
		"target", options.TargetDir,
		"module", options.ModulePath)

	// Create query generator
	queryGen := generator.NewQueryGenerator(".")

	// Generate query component
	return queryGen.GenerateQueryComponent(
		options.FeatureName,
		options.EntityName,
		options.TargetDir,
		options.ModulePath,
	)
}
