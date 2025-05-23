package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	// Placeholder for SSH library: "golang.org/x/crypto/ssh"
)

// NewDeployCmd creates the `deploy` subcommand
func NewDeployCmd() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploys the application binary via SSH",
		Long:  `Builds the application (if necessary) and copies the binary to a target host using SSH.`,
		RunE:  runDeploy,
	}

	// Flags for deployment target
	deployCmd.Flags().String("target-host", "", "Target host (e.g., user@hostname)")
	deployCmd.Flags().String("target-path", "", "Target path on the remote host for the binary")
	deployCmd.Flags().String("ssh-key", "", "Path to the SSH private key")
	deployCmd.Flags().Int("ssh-port", 22, "SSH port on the target host")
	deployCmd.Flags().String("binary-path", "", "Path to the pre-built binary (optional, builds if not provided)")
	// TODO: Add flags for remote commands (e.g., restart service)

	deployCmd.MarkFlagRequired("target-host")
	deployCmd.MarkFlagRequired("target-path")

	return deployCmd
}

func runDeploy(cmd *cobra.Command, args []string) error {
	targetHost, _ := cmd.Flags().GetString("target-host")
	targetPath, _ := cmd.Flags().GetString("target-path")
	sshKeyPath, _ := cmd.Flags().GetString("ssh-key")
	sshPort, _ := cmd.Flags().GetInt("ssh-port")
	binaryPath, _ := cmd.Flags().GetString("binary-path")

	slog.Info("Starting deployment", "target", targetHost, "path", targetPath)

	// 1. Ensure binary exists (build if not provided)
	if binaryPath == "" {
		slog.Info("Binary path not provided, running build step first...")
		// Determine default binary name based on OS
		outputName := "petrock_example_project_name"
		if runtime.GOOS == "windows" {
			outputName += ".exe" // Although deploying windows binary via ssh is less common
		}
		// Use the build command's logic (could be refactored into a shared function)
		buildCmd := NewBuildCmd()
		// Set flags for build suitable for deployment (e.g., target OS/Arch if needed)
		// For simplicity, assume local build is sufficient or configure flags appropriately.
		buildCmd.SetArgs([]string{"-o", outputName}) // Pass necessary build flags
		if err := buildCmd.Execute(); err != nil {
			return fmt.Errorf("build step failed during deploy: %w", err)
		}
		binaryPath = outputName // Use the newly built binary
		slog.Info("Build completed", "binary", binaryPath)
	}

	// Ensure the local binary path exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary file not found at %s", binaryPath)
	}
	absBinaryPath, err := filepath.Abs(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for binary %s: %w", binaryPath, err)
	}

	// 2. Connect via SSH (Placeholder - requires SSH library)
	slog.Info("Connecting via SSH...", "host", targetHost, "port", sshPort)
	// --- SSH Connection Logic ---
	// Use golang.org/x/crypto/ssh or os/exec with scp/ssh commands
	// Example using os/exec (simpler, less secure key handling):
	// scp -P <port> -i <key_path> <local_binary_path> <user@host>:<target_path>
	scpArgs := []string{}
	if sshPort != 22 {
		scpArgs = append(scpArgs, "-P", fmt.Sprintf("%d", sshPort))
	}
	if sshKeyPath != "" {
		scpArgs = append(scpArgs, "-i", sshKeyPath)
	}
	scpArgs = append(scpArgs, absBinaryPath, fmt.Sprintf("%s:%s", targetHost, targetPath))

	scpCmd := exec.Command("scp", scpArgs...)
	scpCmd.Stdout = os.Stdout
	scpCmd.Stderr = os.Stderr
	slog.Info("Executing scp", "command", scpCmd.String())
	if err := scpCmd.Run(); err != nil {
		return fmt.Errorf("scp failed: %w", err)
	}
	slog.Info("Binary copied successfully")

	// 3. Execute remote commands (Placeholder)
	// Example: ssh -p <port> -i <key_path> <user@host> "sudo systemctl restart myapp.service"
	// sshArgs := []string{}
	// ... add port, key, host ...
	// sshArgs = append(sshArgs, "your remote command here") // e.g., "sudo systemctl restart petrock_example_project_name"
	// sshCmd := exec.Command("ssh", sshArgs...)
	// ... run command ...
	slog.Warn("Remote command execution (e.g., service restart) not implemented yet.")

	slog.Info("Deployment finished successfully (manual service restart might be required).")
	return nil
}
