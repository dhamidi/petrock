package generator

import (
	"strings"
	"testing"
)

func TestModifyCommandStructWithFields(t *testing.T) {
	// Sample command template content (simplified)
	templateContent := `package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/petrock/example_module_path/core"
)

type CreateCommand struct {
	// Example fields - replace with actual data needed for creation
	Name        string    ` + "`json:\"name\" validate:\"required\"`" + `
	Description string    ` + "`json:\"description\" validate:\"required\"`" + `
	CreatedBy   string    ` + "`json:\"created_by\"`" + `
	CreatedAt   time.Time ` + "`json:\"created_at\"`" + `
}

func (c *CreateCommand) CommandName() string {
	return "posts/create"
}

func (c *CreateCommand) Validate(state *state.State) error {
	// Lots of validation logic here...
	if c.Name == "" {
		return errors.New("name required")
	}
	return nil
}

func (e *Executor) HandleCreate(ctx context.Context, command core.Command, msg *core.Message, pctx *core.ProcessingContext) error {
	// Lots of handler logic here...
	return nil
}
`

	fields := []CommandField{
		{Name: "postID", Type: "string"},
		{Name: "publishAt", Type: "time.Time"},
	}

	options := ExtractionOptions{
		FeatureName: "posts",
		EntityName:  "publish",
		Fields:      fields,
	}

	cg := &CommandGenerator{}
	result, err := cg.modifyCommandStructWithFields(templateContent, options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the struct has the right fields
	if !strings.Contains(result, "PostID string") {
		t.Error("expected PostID string field in result")
	}
	if !strings.Contains(result, "PublishAt time.Time") {
		t.Error("expected PublishAt time.Time field in result")
	}

	// Check that validate method is simplified
	validateLines := strings.Split(result, "\n")
	validateStarted := false
	validateSimplified := false
	
	for _, line := range validateLines {
		if strings.Contains(line, "func (c *") && strings.Contains(line, ") Validate(") {
			validateStarted = true
			continue
		}
		if validateStarted && strings.Contains(line, "return nil") && strings.TrimSpace(line) == "return nil" {
			validateSimplified = true
			break
		}
	}
	
	if !validateSimplified {
		t.Error("expected Validate method to be simplified to just return nil")
	}

	t.Logf("Modified content:\n%s", result)
}
