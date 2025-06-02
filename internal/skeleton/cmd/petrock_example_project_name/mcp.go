package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/petrock/example_module_path/core"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP (Model Context Protocol) server",
	Long: `Start an MCP server that exposes petrock application functionality
to AI assistants via JSON-RPC 2.0 over stdio.

The MCP server provides tools for:
- generate_command: Generate command components with optional custom fields
- generate_query: Generate query components with optional custom fields
- generate_worker: Generate worker components for background processing
- generate_component: Universal component generator supporting all types

This allows AI assistants like Claude Desktop to generate petrock
components programmatically, accelerating development workflows.

Example usage with Claude Desktop:
1. Add this server to your Claude Desktop MCP configuration
2. Ask Claude to "generate a command component for posts/create"
3. Claude will use the MCP tools to generate the component files`,
	Run: runMCPServer,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

func runMCPServer(cmd *cobra.Command, args []string) {
	// Set up logging to stderr so it doesn't interfere with stdio communication
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting MCP server", "version", "1.0.0")

	// Create and configure the MCP server
	server := core.NewMCPServer()

	// Start the stdio transport
	if err := core.StartStdioServer(server); err != nil {
		slog.Error("MCP server failed", "error", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	slog.Info("MCP server shutdown")
}
