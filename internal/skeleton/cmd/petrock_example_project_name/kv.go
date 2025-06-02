package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/petrock/example_module_path/core"
	"github.com/petrock/example_module_path/core/ui"
	"github.com/spf13/cobra"
)

// NewKVCmd creates the 'kv' parent command for key-value store operations
func NewKVCmd() *cobra.Command {
	kvCmd := &cobra.Command{
		Use:   "kv",
		Short: "Commands for key-value store operations",
		Long:  `Commands for interacting with the application's key-value store.`,
	}

	// Add subcommands
	kvCmd.AddCommand(NewKVGetCmd())
	kvCmd.AddCommand(NewKVSetCmd())
	kvCmd.AddCommand(NewKVListCmd())

	return kvCmd
}

// NewKVGetCmd creates the 'kv get' command
func NewKVGetCmd() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a value from the key-value store",
		Long:  `Retrieves and displays a value from the key-value store by its key.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runKVGet,
	}

	getCmd.Flags().String("db-path", "app.db", "Path to the SQLite database file")

	return getCmd
}

// NewKVSetCmd creates the 'kv set' command
func NewKVSetCmd() *cobra.Command {
	setCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a value in the key-value store",
		Long:  `Stores a value in the key-value store with the specified key.`,
		Args:  cobra.ExactArgs(2),
		RunE:  runKVSet,
	}

	setCmd.Flags().String("db-path", "app.db", "Path to the SQLite database file")
	setCmd.Flags().Bool("json", false, "Parse value as JSON")

	return setCmd
}

// NewKVListCmd creates the 'kv list' command
func NewKVListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list [glob]",
		Short: "List keys in the key-value store",
		Long:  `Lists all keys in the key-value store, optionally filtered by a glob pattern.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runKVList,
	}

	listCmd.Flags().String("db-path", "app.db", "Path to the SQLite database file")

	return listCmd
}

func runKVGet(cmd *cobra.Command, args []string) error {
	dbPath, _ := cmd.Flags().GetString("db-path")
	key := args[0]

	// Initialize the application
	app, err := core.NewApp(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer app.Close()

	// Get the value as interface{}
	var value interface{}
	err = app.KVStore.Get(key, &value)
	if err != nil {
		return fmt.Errorf("failed to get value for key '%s': %w", key, err)
	}

	// Pretty print the value as JSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("failed to encode value as JSON: %w", err)
	}

	return cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, buf.String())
}

func runKVSet(cmd *cobra.Command, args []string) error {
	dbPath, _ := cmd.Flags().GetString("db-path")
	isJSON, _ := cmd.Flags().GetBool("json")
	key := args[0]
	valueStr := args[1]

	// Initialize the application
	app, err := core.NewApp(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer app.Close()

	var value interface{}
	if isJSON {
		// Parse as JSON
		if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
			return fmt.Errorf("failed to parse value as JSON: %w", err)
		}
	} else {
		// Store as string
		value = valueStr
	}

	// Set the value
	if err := app.KVStore.Set(key, value); err != nil {
		return fmt.Errorf("failed to set value for key '%s': %w", key, err)
	}

	return cmdCtx.UI.ShowSuccess(cmdCtx.Ctx, "Successfully set key '%s'\n", key)
}

func runKVList(cmd *cobra.Command, args []string) error {
	dbPath, _ := cmd.Flags().GetString("db-path")
	
	// Default to listing all keys if no glob provided
	glob := "*"
	if len(args) > 0 {
		glob = args[0]
	}

	// Initialize the application
	app, err := core.NewApp(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer app.Close()

	// List keys matching the glob
	keys, err := app.KVStore.List(glob)
	if err != nil {
		return fmt.Errorf("failed to list keys with glob '%s': %w", glob, err)
	}

	// Print each key on a separate line
	for _, key := range keys {
		if err := cmdCtx.UI.Present(cmdCtx.Ctx, ui.MessageTypeInfo, "%s\n", key); err != nil {
			return err
		}
	}

	return nil
}
