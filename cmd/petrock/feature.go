package main

import (
	"errors" // Import errors package
	"fmt"
	"log/slog"
	"os" // Import os package
	"path/filepath" // Import filepath package
	"regexp"
	"strings" // Import strings package

	"github.com/dhamidi/petrock/internal/utils" // Import utils
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
	slog.Debug("Validated feature name", "name", featureName)

	// --- Step 2: Pre-run Checks ---
	slog.Debug("Performing pre-run checks...")

	// 1. Git Clean Check (Handled by PersistentPreRunE in main.go, but double-checking doesn't hurt)
	// Note: PersistentPreRunE in main.go already performs this check.
	// If we want feature-specific pre-run logic beyond the global check, add it here.
	// For now, we rely on the global check.
	// if err := utils.CheckCleanWorkspace(); err != nil {
	// 	return fmt.Errorf("git workspace check failed: %w", err)
	// }
	// slog.Debug("Git workspace is clean.")

	// 2. Project Root Check
	if err := checkIsPetrockProjectRoot("."); err != nil {
		return fmt.Errorf("failed project root validation: %w", err)
	}
	slog.Debug("Current directory appears to be a Petrock project root.")

	// 3. Feature Exists Check
	if _, err := os.Stat(featureName); !errors.Is(err, os.ErrNotExist) {
		if err == nil {
			// Directory exists
			return fmt.Errorf("feature directory %q already exists", featureName)
		}
		// Other error (e.g., permission denied)
		return fmt.Errorf("failed to check status of potential feature directory %q: %w", featureName, err)
	}
	slog.Debug("Feature directory does not already exist", "name", featureName)

	slog.Info("Pre-run checks passed.")

	// --- Placeholder for subsequent steps ---
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

// checkIsPetrockProjectRoot verifies if the given directory looks like a Petrock project root.
func checkIsPetrockProjectRoot(dir string) error {
	// Check for go.mod
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in %s", dir)
	} else if err != nil {
		return fmt.Errorf("failed to check for go.mod in %s: %w", dir, err)
	}

	// Check for core directory
	corePath := filepath.Join(dir, "core")
	if info, err := os.Stat(corePath); os.IsNotExist(err) {
		return fmt.Errorf("core directory not found in %s", dir)
	} else if err != nil {
		return fmt.Errorf("failed to check for core directory in %s: %w", dir, err)
	} else if !info.IsDir() {
		return fmt.Errorf("core path exists but is not a directory in %s", dir)
	}

	// Check for cmd/<project_name>/main.go
	modulePath, err := utils.GetModuleName(dir)
	if err != nil {
		// If GetModuleName fails (e.g., go.mod is invalid), we can't reliably find cmd dir
		slog.Warn("Could not determine module path from go.mod, skipping cmd/<project>/main.go check", "error", err)
		return nil // Or return the error if this check is critical
	}
	parts := strings.Split(modulePath, "/")
	projectName := parts[len(parts)-1] // Assume last part of module path is project name

	cmdMainPath := filepath.Join(dir, "cmd", projectName, "main.go")
	if _, err := os.Stat(cmdMainPath); os.IsNotExist(err) {
		return fmt.Errorf("cmd/%s/main.go not found in %s", projectName, dir)
	} else if err != nil {
		return fmt.Errorf("failed to check for cmd/%s/main.go in %s: %w", projectName, dir, err)
	}

	return nil
}
