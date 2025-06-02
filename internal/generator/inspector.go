package generator

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// ComponentInspector provides collision detection via self inspect integration
type ComponentInspector interface {
	InspectExistingComponents() (*InspectResult, error)
	ComponentExists(componentType ComponentType, featureName, entityName string) (bool, error)
}

// InspectResult represents the JSON structure returned by self inspect
type InspectResult struct {
	Commands []CommandInfo `json:"commands"`
	Queries  []QueryInfo   `json:"queries"`
	Workers  []WorkerInfo  `json:"workers"`
	Routes   []string      `json:"routes"`
	Features []string      `json:"features"`
}

// CommandInfo represents a command in the inspect output
type CommandInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// QueryInfo represents a query in the inspect output
type QueryInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// WorkerInfo represents a worker in the inspect output
type WorkerInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// ComponentInspectorImpl implements ComponentInspector
type ComponentInspectorImpl struct {
	projectPath string
}

// NewComponentInspector creates a new component inspector
func NewComponentInspector(projectPath string) ComponentInspector {
	if projectPath == "" {
		projectPath = "." // Current directory by default
	}
	return &ComponentInspectorImpl{
		projectPath: projectPath,
	}
}

// InspectExistingComponents runs self inspect command and returns parsed result
func (ci *ComponentInspectorImpl) InspectExistingComponents() (*InspectResult, error) {
	slog.Debug("Running self inspect command", "path", ci.projectPath)
	
	// Prepare command: go run ./cmd/{project} self inspect --format=json
	// We need to determine the project binary name from the directory structure
	projectBinary, err := ci.detectProjectBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to detect project binary: %w", err)
	}
	
	cmd := exec.Command("go", "run", fmt.Sprintf("./cmd/%s", projectBinary), "self", "inspect", "--format=json")
	cmd.Dir = ci.projectPath
	
	// Capture output
	output, err := cmd.Output()
	if err != nil {
		// Check if this is because we're not in a petrock project
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			slog.Debug("Self inspect command failed", "stderr", stderr, "stdout", string(output))
			return nil, fmt.Errorf("self inspect failed (are you in a petrock project?): %w", err)
		}
		return nil, fmt.Errorf("failed to execute self inspect: %w", err)
	}
	
	// Parse JSON result
	var result InspectResult
	if err := json.Unmarshal(output, &result); err != nil {
		slog.Debug("Failed to parse self inspect output", "output", string(output))
		return nil, fmt.Errorf("failed to parse self inspect JSON output: %w", err)
	}
	
	slog.Debug("Self inspect completed successfully", 
		"commands", len(result.Commands),
		"queries", len(result.Queries), 
		"workers", len(result.Workers),
		"features", len(result.Features))
	
	return &result, nil
}

// ComponentExists checks if a specific component already exists
func (ci *ComponentInspectorImpl) ComponentExists(componentType ComponentType, featureName, entityName string) (bool, error) {
	result, err := ci.InspectExistingComponents()
	if err != nil {
		return false, err
	}
	
	expectedName := fmt.Sprintf("%s/%s", featureName, entityName)
	
	switch componentType {
	case ComponentTypeCommand:
		for _, cmd := range result.Commands {
			if cmd.Name == expectedName {
				return true, nil
			}
		}
	case ComponentTypeQuery:
		for _, query := range result.Queries {
			if query.Name == expectedName {
				return true, nil
			}
		}
	case ComponentTypeWorker:
		// For workers, check if any worker exists for this feature
		// since each feature typically has one worker infrastructure that handles multiple commands
		for _, worker := range result.Workers {
			// Workers are reported as individual command handlers (e.g., "posts/create")
			// If any worker exists for this feature, the worker infrastructure already exists
			if strings.HasPrefix(worker.Name, featureName+"/") {
				return true, nil
			}
		}
	default:
		return false, fmt.Errorf("unknown component type: %s", componentType)
	}
	
	return false, nil
}

// detectProjectBinary tries to determine the project binary name from cmd directory
func (ci *ComponentInspectorImpl) detectProjectBinary() (string, error) {
	cmdDir := fmt.Sprintf("%s/cmd", ci.projectPath)
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return "", fmt.Errorf("failed to read cmd directory: %w", err)
	}
	
	// Look for directories in cmd/ (should be the project binary)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			return entry.Name(), nil
		}
	}
	
	return "", fmt.Errorf("no project binary found in cmd directory")
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
