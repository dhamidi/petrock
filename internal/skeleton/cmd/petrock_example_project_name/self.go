package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/petrock/example_module_path/core"
	"github.com/spf13/cobra"
)

// NewSelfCmd creates the 'self' parent command for introspection commands
func NewSelfCmd() *cobra.Command {
	selfCmd := &cobra.Command{
		Use:   "self",
		Short: "Commands for application self-inspection",
		Long:  `Commands that provide information about the application itself.`,
	}

	// Add subcommands
	selfCmd.AddCommand(NewSelfInspectCmd())

	return selfCmd
}

// NewSelfInspectCmd creates the 'self inspect' command
func NewSelfInspectCmd() *cobra.Command {
	inspectCmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect the application structure",
		Long:  `Initializes the application and dumps information about its structure in the specified format.`,
		RunE:  runSelfInspect,
	}

	// Add flags
	inspectCmd.Flags().String("format", "json", "Output format: json")
	inspectCmd.Flags().String("db-path", "app.db", "Path to the SQLite database file")

	return inspectCmd
}

func runSelfInspect(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")
	dbPath, _ := cmd.Flags().GetString("db-path")

	// Only support JSON for now
	if format != "json" {
		return fmt.Errorf("unsupported format: %s (only 'json' is currently supported)", format)
	}

	// Initialize the application
	app, err := core.NewApp(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer app.Close()

	// Initialize Application State
	appState := NewAppState()
	app.AppState = appState

	// Create HTTP mux for capturing routes
	app.Mux = http.NewServeMux()

	// Register features
	RegisterAllFeatures(app)

	// We don't need to replay the log since we're only inspecting structure

	// Register core HTTP routes to ensure they're captured
	app.RegisterRoute("GET /", core.HandleIndex(app.CommandRegistry, app.QueryRegistry))
	app.RegisterRoute("GET /commands", handleListCommands(app.CommandRegistry))
	app.RegisterRoute("POST /commands", handleExecuteCommand(app.Executor, app.CommandRegistry))
	app.RegisterRoute("GET /queries", handleListQueries(app.QueryRegistry))
	app.RegisterRoute("GET /queries/{feature}/{queryName}", handleExecuteQuery(app.QueryRegistry))

	// Gather application metadata
	result := app.GetInspectResult()

	// Output as JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("failed to encode result as JSON: %w", err)
	}

	return nil
}