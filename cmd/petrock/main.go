package main

import (
	"log/slog" // Import slog
	"os"
	"strings" // Import strings

	"github.com/dhamidi/petrock/internal/utils" // Import the utils package
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "petrock",
	Short: "Petrock is a tool for bootstrapping and managing Go web projects.",
	Long: `Petrock helps create new Go projects based on event sourcing principles
and generate feature modules within existing Petrock projects.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Only check git workspace for major operations that warrant such caution
		requireCleanGit := cmd.Name() == "feature"
		if !requireCleanGit {
			return nil
		}
		// Check if the Git workspace is clean before running major operations
		if err := utils.CheckCleanWorkspace(); err != nil {
			// Return the error directly; CheckCleanWorkspace provides context.
			// Adding more context here might be redundant unless clarifying *why* it's checked.
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands defined in other files
	// Note: newCmd and testCmd are registered in their respective files
	rootCmd.AddCommand(featureCmd) // From feature.go

	// Use Cobra's default error handling

	// Configure logging level based on environment variable
	logLevel := slog.LevelInfo // Default level
	levelStr := strings.ToLower(os.Getenv("PETROCK_LOG_LEVEL"))
	switch levelStr {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel <= slog.LevelDebug, // Add source only for debug or lower
	}
	handler := slog.NewTextHandler(os.Stderr, opts)
	slog.SetDefault(slog.New(handler))
}

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
