package main

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/spf13/cobra"
)

// featureNameRegex validates that a feature name is a valid Go package name
// (lowercase letters, numbers, underscore, starting with a letter).
// It intentionally disallows hyphens which are valid in directory names but not package names.
var featureNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// featureCmd represents the feature command
var featureCmd = &cobra.Command{
	Use:   "feature [featureName]",
	Short: "Generates a new feature module in the current Petrock project",
	Long: `Generates a new Go package directory containing skeleton code for a
new application feature (e.g., command/query handlers, state, views).

The [featureName] must be a valid Go package name (lowercase letters, numbers, underscores, starting with a letter).
This command must be run from the root directory of an existing Petrock project.`,
	Args: cobra.ExactArgs(1), // Requires exactly one argument: the feature name
	RunE: runFeature,
}

func init() {
	// This function will be called by Cobra's initialization process
	// No need to explicitly add to rootCmd here, main.go's init will do it.
}

// runFeature executes the logic for the 'petrock feature' command.
func runFeature(cmd *cobra.Command, args []string) error {
	featureName := args[0]
	slog.Debug("Starting feature creation", "featureName", featureName)

	// Validate feature name format
	if !featureNameRegex.MatchString(featureName) {
		return fmt.Errorf("invalid feature name %q: must be a valid Go package name (lowercase letters, numbers, underscores, starting with a letter)", featureName)
	}

	slog.Info("Validated feature name", "name", featureName)

	// --- Placeholder for subsequent steps ---
	// 1. Pre-run checks (Git clean, project root, feature exists)
	// 2. Copy skeleton
	// 3. Replace placeholders
	// 4. Modify features.go
	// 5. Run go mod tidy
	// 6. Git commit
	// 7. Output success message
	// --- End Placeholder ---

	fmt.Printf("Feature command executed for: %s (Implementation pending)\n", featureName) // Placeholder output
	slog.Debug("runFeature completed (placeholder)", "featureName", featureName)
	return nil // Return nil on success
}
