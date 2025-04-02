package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"petrock/internal/utils" // Import the utils package
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
	// Add subcommands here
	// newCmd is defined in cmd/petrock/new.go
	// featureCmd will be defined in cmd/petrock/feature.go
	// rootCmd.AddCommand(featureCmd)
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
