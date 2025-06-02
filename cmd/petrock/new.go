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
	Short: "Creates a new Petrock project structure or generates components",
	Long: `Creates a new directory with the specified project name, initializes a Go module
with the given module path, sets up a git repository, and generates the
initial project files based on Petrock templates.

For component generation, use the subcommands:
  petrock new command <feature>/<name-of-thing>   - Generate command component
  petrock new query <feature>/<name-of-thing>     - Generate query component  
  petrock new worker <feature>/<name-of-thing>    - Generate worker component

Examples:
  petrock new myblog github.com/youruser/myblog
  petrock new command posts/create
  petrock new query posts/get
  petrock new worker posts/summary`,
	Args: cobra.ExactArgs(2), // Require both project name and module path
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	
	// Add component subcommands
	newCmd.AddCommand(newCommandCmd())
	newCmd.AddCommand(newQueryCmd())
	newCmd.AddCommand(newWorkerCmd())
	
	// Add flags here if needed in the future (e.g., --template-set)
}



// parseFeatureEntityName parses feature/entity format
func parseFeatureEntityName(input string) (string, string, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format %q: expected <feature>/<name-of-thing>", input)
	}
	
	featureName := strings.TrimSpace(parts[0])
	entityName := strings.TrimSpace(parts[1])
	
	if featureName == "" || entityName == "" {
		return "", "", fmt.Errorf("invalid format %q: feature and entity names cannot be empty", input)
	}
	
	// Basic validation of names (similar to existing validation)
	if !dirNameRegex.MatchString(featureName) {
		return "", "", fmt.Errorf("invalid feature name %q: must contain only letters, numbers, '.', '_', '-'", featureName)
	}
	if !dirNameRegex.MatchString(entityName) {
		return "", "", fmt.Errorf("invalid entity name %q: must contain only letters, numbers, '.', '_', '-'", entityName)
	}
	
	return featureName, entityName, nil
}

// newCommandCmd creates the command subcommand
func newCommandCmd() *cobra.Command {
	return NewCommandSubcommand()
}

// newQueryCmd creates the query subcommand  
func newQueryCmd() *cobra.Command {
	return NewQuerySubcommand()
}

// newWorkerCmd creates the worker subcommand
func newWorkerCmd() *cobra.Command {
	return NewWorkerSubcommand()
}



// detectModulePath reads go.mod to determine the current module path
func detectModulePath(projectPath string) (string, error) {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}
	
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}
	
	return "", fmt.Errorf("module directive not found in go.mod")
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
