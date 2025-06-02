package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/petrock/example_module_path/core/ui"
	"github.com/spf13/cobra"
)

// CommandContext holds shared dependencies for all commands
type CommandContext struct {
	UI  ui.UI
	Ctx context.Context
}

// Global command context
var cmdCtx *CommandContext

var rootCmd = &cobra.Command{
	Use:   "petrock_example_project_name",
	Short: "The main command for the petrock_example_project_name application.",
	Long:  `petrock_example_project_name application entry point.`,
	// Run: func(cmd *cobra.Command, args []string) { }, // Or remove if subcommands are mandatory
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// Initialization of core components like DB, registries, etc.,
	// is typically done within the specific command that needs them (e.g., 'serve').
	// core.InitLog() // Example - Removed as global init is discouraged
	// core.InitRegistries() // Example - Removed as global init is discouraged

	// Register features
	// TODO: Replace with actual feature registration logic
	// RegisterAllFeatures(core.Commands, core.Queries) // Example call

	return rootCmd.Execute()
}

func init() {
	// Add subcommands here
	rootCmd.AddCommand(NewServeCmd())
	rootCmd.AddCommand(NewBuildCmd())
	rootCmd.AddCommand(NewDeployCmd())
	rootCmd.AddCommand(NewSelfCmd())
	rootCmd.AddCommand(NewKVCmd())

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

// newCommandContext creates a new command context with UI
func newCommandContext() *CommandContext {
	return &CommandContext{
		UI:  ui.NewConsoleUI(),
		Ctx: context.Background(),
	}
}

// configureUI makes the UI available to commands
func configureUI() {
	cmdCtx = newCommandContext()
}

func main() {
	// Initialize UI before executing commands
	configureUI()
	
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
