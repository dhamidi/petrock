package generator

import (
	"strings"
	"testing"

	"github.com/dhamidi/petrock/internal/ed"
)

func TestCommandModificationWithEditor(t *testing.T) {
	// Test just the editor logic for command modification
	templateContent := `type CreateCommand struct {
	// Example fields - replace with actual data needed for creation
	Name        string    
	Description string    
	CreatedBy   string    
	CreatedAt   time.Time 
}`

	// Test field replacement
	fields := []CommandField{
		{Name: "postID", Type: "string"},
		{Name: "publishAt", Type: "time.Time"},
	}

	// Build the field definitions string
	var fieldDefs []string
	for _, field := range fields {
		// Capitalize first letter for exported fields
		capitalizedName := strings.ToUpper(field.Name[:1]) + field.Name[1:]
		fieldDef := "\t" + capitalizedName + " " + field.Type
		fieldDefs = append(fieldDefs, fieldDef)
	}
	fieldDefsStr := strings.Join(fieldDefs, "\n") + "\n"

	editor := ed.New(templateContent)
	
	// Find the command struct and replace its fields
	err := editor.Do(
		ed.BeginningOfBuffer(),
		ed.Search("type"),
		ed.Search("Command struct {"),
		ed.Search("{"),
		ed.ForwardChar(1), // Move past the opening brace
		ed.SetMark(),
		ed.Search("}"),     // Find closing brace
		ed.ReplaceRegion("\n"+fieldDefsStr),
	)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := editor.String()
	
	// Check that the struct has the right fields
	if !strings.Contains(result, "PostID string") {
		t.Error("expected PostID string field in result")
	}
	if !strings.Contains(result, "PublishAt time.Time") {
		t.Error("expected PublishAt time.Time field in result")
	}

	t.Logf("Modified content:\n%s", result)
}

func TestValidateMethodSimplification(t *testing.T) {
	content := `func (c *CreateCommand) Validate(state *state.State) error {
	// Lots of validation logic here...
	if c.Name == "" {
		return errors.New("name required")
	}
	return nil
}`

	editor := ed.New(content)
	
	err := editor.Do(
		ed.BeginningOfBuffer(),
		ed.Search("func (c *"),
		ed.Search(") Validate("),
		ed.Search("{"),
		ed.ForwardChar(1), // Move past the opening brace
		ed.SetMark(),
		ed.Search("return"),
		ed.Search("nil"),
		ed.ForwardChar(3), // Move past "nil"
		ed.ReplaceRegion("\n\treturn nil\n"),
	)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := editor.String()
	expected := `func (c *CreateCommand) Validate(state *state.State) error {
	return nil
}`

	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	t.Logf("Simplified method:\n%s", result)
}
