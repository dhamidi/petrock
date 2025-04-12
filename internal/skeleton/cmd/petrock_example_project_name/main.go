package main

import (
	"fmt"
	"os"

	// "github.com/petrock/example_module_path/core" // Removed as it's not directly used in main.go anymore
	// Core components are initialized/used within specific commands (e.g., serve)

	"github.com/spf13/cobra"
)

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
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
