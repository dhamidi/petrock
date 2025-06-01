package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/petrock/example_module_path/core"
	"github.com/spf13/cobra"
)

// NewMCPCmd creates the mcp command
func NewMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start MCP (Model Context Protocol) server on stdio",
		Long: `Start an MCP server that communicates over stdio using JSON-RPC 2.0.

This allows AI applications to connect to your petrock application and:
- Access application data through resources
- Execute queries and operations through tools  
- Use templated prompts for common workflows

The server runs on stdin/stdout for integration with MCP clients.`,
		RunE: runMCPServer,
	}

	// Add flags for configuration
	cmd.Flags().String("db-path", "app.db", "Path to the SQLite database")
	cmd.Flags().String("log-level", "info", "Log level (debug, info, warn, error)")

	return cmd
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	// Get flags
	dbPath, _ := cmd.Flags().GetString("db-path")
	logLevel, _ := cmd.Flags().GetString("log-level")

	// Setup logger - log to stderr so it doesn't interfere with stdio JSON-RPC
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))

	// Initialize the core application
	app, err := core.NewApp(dbPath, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}
	defer app.Close()

	// Register features - this would normally be done by the main application
	// but we need to do it here for MCP server as well
	registerFeatures(app)

	// Replay event log to get current state
	if err := app.ReplayEventLog(); err != nil {
		logger.Warn("Failed to replay event log", "error", err)
	}

	// Create MCP server
	mcpServer := core.NewMCPServer(app, logger)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping MCP server")
		cancel()
	}()

	// Start the MCP server
	logger.Info("Starting MCP server on stdio", 
		"db_path", dbPath,
		"log_level", logLevel,
	)

	if err := mcpServer.Serve(ctx); err != nil && err != context.Canceled {
		return fmt.Errorf("MCP server error: %w", err)
	}

	logger.Info("MCP server stopped")
	return nil
}

// registerFeatures registers all application features
// This is a placeholder - in real applications this would import and register
// actual feature modules
func registerFeatures(app *core.App) {
	// Example feature registration - replace with actual features
	
	// Register command handlers
	app.Commands.Register("ping", func(ctx context.Context, input map[string]interface{}) (interface{}, error) {
		return map[string]interface{}{
			"message": "pong",
			"timestamp": ctx.Value("timestamp"),
		}, nil
	})

	app.Commands.Register("version", func(ctx context.Context, input map[string]interface{}) (interface{}, error) {
		return map[string]interface{}{
			"version": "1.0.0",
			"name":    "petrock_example_project_name",
		}, nil
	})

	// Register query handlers  
	app.Queries.Register("health", func(ctx context.Context, input map[string]interface{}) (interface{}, error) {
		return map[string]interface{}{
			"status": "healthy",
			"database": "connected",
		}, nil
	})

	// Add more feature registrations here as needed
	// For example:
	// posts.RegisterFeature(app)
	// users.RegisterFeature(app)
	// etc.
}
