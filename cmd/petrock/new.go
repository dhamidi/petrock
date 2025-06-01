package main

import (
	// "embed" // Removed embed import here
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath" // Added import for filepath.Join
	"regexp"
	"strings"

	// "github.com/dhamidi/petrock/internal/template" // Removed template import
	// "github.com/dhamidi/petrock/internal/skeletonfs" // Removed import for skeletonfs
	petrock "github.com/dhamidi/petrock" // Import root package for embedded FS
	"github.com/dhamidi/petrock/internal/utils"

	"github.com/spf13/cobra"
)

var (
	// Simple regex for basic validation of project/module names.
	// Allows letters, numbers, underscores, hyphens, dots, and forward slashes (for module path).
	// Does not allow starting/ending with separators or consecutive separators.
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._\-/]*[a-zA-Z0-9])?$`)
	// Regex to validate a simple directory name (no slashes).
	dirNameRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._-]*[a-zA-Z0-9])?$`)
)

// //go:embed all:../../internal/skeleton // Removed embed directive here
// var skeletonFS embed.FS // Removed local embed FS variable

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [projectName] [modulePath]",
	Short: "Creates a new Petrock project structure",
	Long: `Creates a new directory with the specified project name, initializes a Go module
with the given module path, sets up a git repository, and generates the
initial project files based on Petrock templates.

Example:
  petrock new myblog github.com/youruser/myblog`,
	Args: cobra.ExactArgs(2), // Require both project name and module path
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	// Add flags here if needed in the future (e.g., --template-set)
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]
	modulePath := args[1]

	slog.Debug("Starting new project creation", "project", projectName, "module", modulePath) // Changed to Debug

	// Validate inputs
	if !dirNameRegex.MatchString(projectName) {
		return fmt.Errorf("invalid project name %q: must contain only letters, numbers, '.', '_', '-' and cannot start/end with separators", projectName)
	}
	if !nameRegex.MatchString(modulePath) || strings.Contains(modulePath, "//") || strings.Contains(modulePath, "..") {
		return fmt.Errorf("invalid module path %q: must be a valid Go module path", modulePath)
	}

	// Check if project directory already exists
	if _, err := os.Stat(projectName); !errors.Is(err, os.ErrNotExist) {
		if err == nil {
			return fmt.Errorf("directory %q already exists", projectName)
		}
		// Other error (e.g., permission denied)
		return fmt.Errorf("failed to check directory status for %q: %w", projectName, err)
	}

	// Create project directory
	slog.Debug("Creating project directory", "path", projectName)
	if err := utils.EnsureDir(projectName); err != nil {
		return fmt.Errorf("failed to create project directory %q: %w", projectName, err)
	}

	// Initialize Git repository
	slog.Debug("Initializing Git repository", "path", projectName)
	if err := utils.GitInit(projectName); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// --- Copy Skeleton and Replace Placeholders ---
	projectNamePlaceholder := "petrock_example_project_name"
	modulePathPlaceholder := "github.com/petrock/example_module_path"

	// Copy skeleton directory structure from embedded FS, excluding the feature template
	slog.Debug("Copying skeleton project structure from embedded FS", "to", projectName)
	// Pass the embedded FS from the root petrock package
	// Start copying from the 'internal/skeleton' directory within the embed FS
	exclude := []string{"internal/skeleton/petrock_example_feature_name"}
	err := utils.CopyDir(petrock.SkeletonFS, "internal/skeleton", projectName, projectNamePlaceholder, projectName, exclude)
	if err != nil {
		return fmt.Errorf("failed to copy skeleton directory from embedded FS: %w", err)
	}

	// Rename go.mod.skel to go.mod after copying
	slog.Debug("Renaming go.mod.skel to go.mod", "path", projectName)
	skelModPath := filepath.Join(projectName, "go.mod.skel")
	targetModPath := filepath.Join(projectName, "go.mod")
	if err := os.Rename(skelModPath, targetModPath); err != nil {
		// Check if the source file exists, maybe CopyDir failed silently?
		if _, statErr := os.Stat(skelModPath); os.IsNotExist(statErr) {
			return fmt.Errorf("failed to rename go.mod.skel: source file %s not found after copy", skelModPath)
		}
		return fmt.Errorf("failed to rename %s to %s: %w", skelModPath, targetModPath, err)
	}


	// Define replacements
	replacements := map[string]string{
		projectNamePlaceholder: projectName,
		modulePathPlaceholder:  modulePath,
	}

	// Replace placeholders in copied files
	slog.Debug("Replacing placeholders in project files", "path", projectName)
	if err := utils.ReplaceInFiles(projectName, replacements); err != nil {
		return fmt.Errorf("failed to replace placeholders in project files: %w", err)
	}
	// --- End Copy & Replace ---

	// Tidy Go module dependencies (after go.mod and source files are created)
	slog.Debug("Running go mod tidy", "path", projectName)
	if err := utils.GoModTidy(projectName); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	// Initial Git commit
	slog.Debug("Creating initial Git commit", "path", projectName)
	if err := utils.GitAddAll(projectName); err != nil {
		return fmt.Errorf("failed to stage files in git: %w", err)
	}
	commitMsg := fmt.Sprintf("Initial project structure for %s generated by petrock", projectName)
	if err := utils.GitCommit(projectName, commitMsg); err != nil {
		return fmt.Errorf("failed to create initial git commit: %w", err)
	}

	slog.Debug("Project created successfully!", "path", projectName) // Changed to Debug
	fmt.Printf("\nSuccess! Created project %s at ./%s\n", projectName, projectName)
	fmt.Printf("Module path: %s\n", modulePath)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd ./%s\n", projectName)
	fmt.Printf("  go run ./cmd/%s serve\n", projectName)

	return nil
}
