package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// NewBuildCmd creates the `build` subcommand
func NewBuildCmd() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Builds the application binary",
		Long:  `Compiles the application into a single executable binary, optionally embedding assets.`,
		RunE:  runBuild,
	}

	outputName := "petrock_example_project_name"
	if runtime.GOOS == "windows" {
		outputName += ".exe"
	}

	buildCmd.Flags().StringP("output", "o", outputName, "Output binary name")
	buildCmd.Flags().String("goos", runtime.GOOS, "Target operating system (GOOS)")
	buildCmd.Flags().String("goarch", runtime.GOARCH, "Target architecture (GOARCH)")
	buildCmd.Flags().String("ldflags", "-s -w", "Linker flags (e.g., '-s -w' to strip symbols)")
	// TODO: Add flags for version injection via ldflags

	return buildCmd
}

func runBuild(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	goos, _ := cmd.Flags().GetString("goos")
	goarch, _ := cmd.Flags().GetString("goarch")
	ldflags, _ := cmd.Flags().GetString("ldflags")

	slog.Info("Starting build process", "output", output, "goos", goos, "goarch", goarch, "ldflags", ldflags)

	// Ensure the output path is absolute or relative to the current dir
	if !filepath.IsAbs(output) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}
		output = filepath.Join(cwd, output)
	}

	// Prepare the build command
	// Assumes the main package is in the current directory or specified correctly.
	// If the main package is under cmd/{{.ProjectName}}, adjust the target path.
	buildArgs := []string{
		"build",
		"-ldflags=" + ldflags,
		"-o", output,
		"./cmd/petrock_example_project_name", // Target the main package
	}

	buildCmd := exec.Command("go", buildArgs...)
	buildCmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", goos), fmt.Sprintf("GOARCH=%s", goarch))
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	slog.Info("Executing go build", "command", buildCmd.String())

	err := buildCmd.Run()
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	slog.Info("Build successful", "output", output)
	return nil
}
