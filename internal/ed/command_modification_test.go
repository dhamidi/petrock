package ed

import (
	"strings"
	"testing"
)

type CommandField struct {
	Name string
	Type string
}

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

	editor := New(templateContent)
	
	// Find the command struct and replace its fields
	err := editor.Do(
		BeginningOfBuffer(),
		Search("type"),
		Search("Command struct {"),
		Search("{"),
		ForwardChar(1), // Move past the opening brace
		SetMark(),
		Search("}"),     // Find closing brace
		ReplaceRegion("\n"+fieldDefsStr),
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

	// Check that original fields are gone
	if strings.Contains(result, "Name        string") {
		t.Error("expected original Name field to be replaced")
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

	editor := New(content)
	
	err := editor.Do(
		BeginningOfBuffer(),
		Search("func (c *"),
		Search(") Validate("),
		Search("{"),
		ForwardChar(1), // Move past the opening brace
		SetMark(),
		Search("return"),
		Search("nil"),
		ForwardChar(3), // Move past "nil"
		ReplaceRegion("\n\treturn nil\n"),
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

func TestHandlerMethodSimplification(t *testing.T) {
	content := `func (e *Executor) HandleCreate(ctx context.Context, command core.Command, msg *core.Message, pctx *core.ProcessingContext) error {
	// Type assertion for pointer type
	cmd, ok := command.(*CreateCommand)
	if !ok {
		err := fmt.Errorf("internal error: incorrect command type (%T) passed to HandleCreate, expected *CreateCommand", command)
		slog.Error("Type assertion failed in HandleCreate", "error", err)
		return err // Returning error causes panic in core.Executor
	}

	slog.Debug("Applying state change for CreateCommand", "feature", "posts", "name", cmd.Name)

	// Lots of logic here...
	return nil
}`

	editor := New(content)
	
	err := editor.Do(
		BeginningOfBuffer(),
		Search("func (e *Executor) Handle"),
		Search("{"),
		ForwardChar(1), // Move past the opening brace
		SetMark(),
		Search("return nil"),
		ForwardChar(10), // Move past "return nil"
		ReplaceRegion("\n\treturn nil\n"),
	)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := editor.String()
	
	// Check that it's simplified
	if strings.Contains(result, "Type assertion") {
		t.Error("expected handler to be simplified, but still contains type assertion logic")
	}

	expectedStart := "func (e *Executor) HandleCreate(ctx context.Context, command core.Command, msg *core.Message, pctx *core.ProcessingContext) error {\n\treturn nil\n\n}"
	if result != expectedStart {
		t.Errorf("expected simplified handler, got:\n%s", result)
	}

	t.Logf("Simplified handler:\n%s", result)
}
