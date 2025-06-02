package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandGenerationWithFields(t *testing.T) {
	// Skip if running with regular tests to avoid filesystem conflicts
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "petrock_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up test fields
	fields := []CommandField{
		{Name: "postID", Type: "string"},
		{Name: "publishAt", Type: "time.Time"},
	}

	// Create command generator
	cmdGen := NewCommandGenerator(tempDir)

	// Generate command component with fields
	err = cmdGen.GenerateCommandComponentWithFields(
		"posts",     // featureName
		"publish",   // entityName
		tempDir,     // targetDir
		"github.com/test/app", // modulePath
		fields,
	)

	if err != nil {
		t.Fatalf("failed to generate command: %v", err)
	}

	// Read the generated command file
	commandFile := filepath.Join(tempDir, "posts", "commands", "publish.go")
	content, err := os.ReadFile(commandFile)
	if err != nil {
		t.Fatalf("failed to read generated command file: %v", err)
	}

	contentStr := string(content)
	t.Logf("Generated command content:\n%s", contentStr)

	// Verify the struct has the correct fields
	if !strings.Contains(contentStr, "PostID string") {
		t.Error("expected PostID string field in generated command")
	}
	if !strings.Contains(contentStr, "PublishAt time.Time") {
		t.Error("expected PublishAt time.Time field in generated command")
	}

	// Verify the struct is named correctly
	if !strings.Contains(contentStr, "type PublishCommand struct") {
		t.Error("expected PublishCommand struct name")
	}

	// Verify the handler method is named correctly
	if !strings.Contains(contentStr, "func (e *Executor) HandlePublish(") {
		t.Error("expected HandlePublish method name")
	}

	// Verify the methods are simplified
	// The Validate method should be simple
	validateLines := strings.Split(contentStr, "\n")
	validateFound := false
	validateSimple := false
	
	for i, line := range validateLines {
		if strings.Contains(line, "func (c *PublishCommand) Validate(") {
			validateFound = true
			// Check if the next few lines contain only "return nil"
			for j := i + 1; j < len(validateLines) && j < i+5; j++ {
				if strings.TrimSpace(validateLines[j]) == "return nil" {
					validateSimple = true
					break
				}
			}
			break
		}
	}
	
	if !validateFound {
		t.Error("expected Validate method to be present")
	}
	if !validateSimple {
		t.Error("expected Validate method to be simplified to just return nil")
	}

	// Similar check for handler method
	handlerLines := strings.Split(contentStr, "\n")
	handlerFound := false
	handlerSimple := false
	
	for i, line := range handlerLines {
		if strings.Contains(line, "func (e *Executor) HandlePublish(") {
			handlerFound = true
			// Check if the method body is simple
			for j := i + 1; j < len(handlerLines) && j < i+5; j++ {
				if strings.TrimSpace(handlerLines[j]) == "return nil" {
					handlerSimple = true
					break
				}
			}
			break
		}
	}
	
	if !handlerFound {
		t.Error("expected HandlePublish method to be present")
	}
	if !handlerSimple {
		t.Error("expected HandlePublish method to be simplified to just return nil")
	}
}
