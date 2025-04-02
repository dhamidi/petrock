package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"petrock/internal/template"
	"petrock/internal/utils"

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

	slog.Info("Starting new project creation", "project", projectName, "module", modulePath)

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

	// Prepare template data
	templateData := map[string]string{
		"ProjectName": projectName,
		"ModuleName":  modulePath,
	}

	// List of templates to render (source path relative to embed FS root -> target path relative to project dir)
	// Using filepath.Join for target paths ensures OS compatibility.
	// Using forward slashes for template names ensures embed.FS compatibility.
	templatesToRender := map[string]string{
		"new/.gitignore.tmpl":                                ".gitignore",
		"new/go.mod.tmpl":                                    "go.mod",
		"new/cmd/main.go.tmpl":                               filepath.Join("cmd", projectName, "main.go"),
		"new/cmd/serve.go.tmpl":                              filepath.Join("cmd", projectName, "serve.go"),
		"new/cmd/build.go.tmpl":                              filepath.Join("cmd", projectName, "build.go"),
		"new/cmd/deploy.go.tmpl":                             filepath.Join("cmd", projectName, "deploy.go"),
		"new/cmd/features.go.tmpl":                           filepath.Join("cmd", projectName, "features.go"),
		"new/core/commands.go.tmpl":                          filepath.Join("core", "commands.go"),
		"new/core/queries.go.tmpl":                           filepath.Join("core", "queries.go"),
		"new/core/form.go.tmpl":                              filepath.Join("core", "form.go"),
		"new/core/log.go.tmpl":                               filepath.Join("core", "log.go"),
		"new/core/view.go.tmpl":                              filepath.Join("core", "view.go"),
		"new/core/view_layout.go.tmpl":                       filepath.Join("core", "view_layout.go"),
		"new/core/page_index.go.tmpl":                        filepath.Join("core", "page_index.go"),
		// Add other core files here if needed
	}

	// Render templates
	slog.Debug("Rendering project templates...")
	for tmplName, targetRelPath := range templatesToRender {
		targetAbsPath := filepath.Join(projectName, targetRelPath)
		slog.Debug("Rendering template", "template", tmplName, "target", targetAbsPath)
		err := template.RenderTemplate(template.Templates, targetAbsPath, tmplName, templateData)
		if err != nil {
			return fmt.Errorf("failed to render template %s to %s: %w", tmplName, targetAbsPath, err)
		}
	}

	// Initialize Go module (after go.mod is created)
	slog.Debug("Initializing Go module", "path", projectName, "module", modulePath)
	if err := utils.GoModInit(projectName, modulePath); err != nil {
		// GoModInit already includes modulePath in its error message
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	// Tidy Go module dependencies
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

	slog.Info("Project created successfully!", "path", projectName)
	fmt.Printf("\nSuccess! Created project %s at ./%s\n", projectName, projectName)
	fmt.Printf("Module path: %s\n", modulePath)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd ./%s\n", projectName)
	fmt.Printf("  go run ./cmd/%s serve\n", projectName)

	return nil
}
