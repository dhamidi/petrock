package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CheckCleanWorkspace checks if the git working directory is clean.
// It returns an error if there are uncommitted changes or untracked files.
func CheckCleanWorkspace() error {
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // Capture stderr as well, in case of git errors

	err := cmd.Run()
	if err != nil {
		// If git status itself fails, return that error
		return fmt.Errorf("failed to run git status: %w\nOutput:\n%s", err, out.String())
	}

	if strings.TrimSpace(out.String()) != "" {
		return fmt.Errorf("git workspace is not clean:\n%s", out.String())
	}

	return nil
}

// GitInit initializes a new git repository in the specified directory.
func GitInit(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run git init in %s: %w\nOutput:\n%s", dir, err, string(output))
	}
	return nil
}

// GitAddAll stages all changes in the specified directory.
func GitAddAll(dir string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run git add . in %s: %w\nOutput:\n%s", dir, err, string(output))
	}
	return nil
}

// GitCommit creates a commit with the given message in the specified directory.
func GitCommit(dir string, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run git commit in %s: %w\nOutput:\n%s", dir, err, string(output))
	}
	return nil
}
