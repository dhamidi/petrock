package main

import (
	"fmt"
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
		// Skip git check for the 'new' command as it runs before the repo exists.
		if cmd.Name() == "new" {
			return nil
		}
		// Check if the Git workspace is clean before running commands other than 'new'
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
	// Subcommands are added in their respective files' init() functions (e.g., new.go, test.go).
	// Placeholder for adding featureCmd later:
	// rootCmd.AddCommand(featureCmd)

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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
