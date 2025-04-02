package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// Placeholder for utils when created: "petrock/internal/utils"
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
		// Placeholder: Call utils.CheckCleanWorkspace() here
		// if err := utils.CheckCleanWorkspace(); err != nil {
		// 	return fmt.Errorf("git workspace is not clean, please commit or stash changes: %w", err)
		// }
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Placeholder: Add subcommands here later
	// rootCmd.AddCommand(newCmd)
	// rootCmd.AddCommand(featureCmd)
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
