package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GoModInit runs `go mod init <modulePath>` in the specified directory.
func GoModInit(dir, modulePath string) error {
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run go mod init %s in %s: %w\nOutput:\n%s", modulePath, dir, err, string(output))
	}
	return nil
}

// GoModTidy runs `go mod tidy` in the specified directory.
func GoModTidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run go mod tidy in %s: %w\nOutput:\n%s", dir, err, string(output))
	}
	return nil
}

// GetModuleName reads the go.mod file in the specified directory and returns the module path.
func GetModuleName(dir string) (string, error) {
	goModPath := filepath.Join(dir, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("could not open %s: %w", goModPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				return parts[1], nil
			}
			return "", fmt.Errorf("malformed module line in %s: %s", goModPath, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning %s: %w", goModPath, err)
	}

	return "", fmt.Errorf("module line not found in %s", goModPath)
}
