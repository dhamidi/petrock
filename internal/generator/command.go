package generator

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	petrock "github.com/dhamidi/petrock"
	"github.com/dhamidi/petrock/internal/ed"
	"github.com/dhamidi/petrock/internal/generator/templates"
)

// CommandGenerator implements ComponentGenerator for command-specific generation
type CommandGenerator struct {
	inspector ComponentInspector
}

// NewCommandGenerator creates a new command-specific generator
func NewCommandGenerator(projectPath string) *CommandGenerator {
	return &CommandGenerator{
		inspector: NewComponentInspector(projectPath),
	}
}

// ExtractCommandFiles extracts command-specific files from skeleton
func (cg *CommandGenerator) ExtractCommandFiles(featureName, entityName string, options ExtractionOptions) error {
	slog.Debug("Extracting command files",
	"feature", featureName,
	"name", entityName)

	// Get command-specific file list
	commandFiles, err := cg.getCommandFileList(featureName, entityName)
	if err != nil {
		return fmt.Errorf("failed to get command file list: %w", err)
	}

	// Update extraction options with command files and file mapping
	options.SkeletonFiles = commandFiles
	options.FileMapping = templates.GetCommandTemplateFiles(entityName)
	
	// Use base ComponentGenerator extraction logic
	baseGen := NewComponentGenerator(".")
	
	// If we have custom fields, use enhanced extraction with template modification
	if len(options.Fields) > 0 {
		return cg.ExtractCommandFilesWithFields(options)
	}
	
	return baseGen.ExtractComponent(options)
}

// CommandField represents a field definition for a command struct
type CommandField struct {
	Name string
	Type string
}

// GenerateCommandComponent generates a complete command component
func (cg *CommandGenerator) GenerateCommandComponent(featureName, entityName, targetDir, modulePath string) error {
	return cg.GenerateCommandComponentWithFields(featureName, entityName, targetDir, modulePath, nil)
}

// GenerateCommandComponentWithFields generates a complete command component with custom fields
func (cg *CommandGenerator) GenerateCommandComponentWithFields(featureName, entityName, targetDir, modulePath string, fields []CommandField) error {
	slog.Debug("Generating command component", 
		"feature", featureName, 
		"name", entityName,
		"target", targetDir,
		"fields", len(fields))

	// Validate entity name
	if err := templates.ValidateCommandEntity(entityName); err != nil {
		return fmt.Errorf("invalid command entity: %w", err)
	}

	// Check for collisions
	exists, err := cg.inspector.ComponentExists(ComponentTypeCommand, featureName, entityName)
	if err != nil {
		slog.Warn("Could not check for existing commands", "error", err.Error())
	}
	if exists {
		return fmt.Errorf("command %s/%s already exists", featureName, entityName)
	}

	// Build command placeholders
	placeholders := templates.BuildCommandPlaceholders(featureName, entityName, modulePath)
	
	// Prepare extraction options
	extractOptions := ExtractionOptions{
		ComponentType: ComponentTypeCommand,
		FeatureName:   featureName,
		EntityName:    entityName,
		TargetDir:     targetDir,
		Replacements:  templates.GetCommandReplacements(placeholders),
		Fields:        fields,
	}

	// Extract command files
	return cg.ExtractCommandFiles(featureName, entityName, extractOptions)
}

// ExtractCommandFilesWithFields extracts command files and modifies them with custom fields using the editor
func (cg *CommandGenerator) ExtractCommandFilesWithFields(options ExtractionOptions) error {
	slog.Debug("Extracting command files with custom fields",
		"feature", options.FeatureName,
		"entity", options.EntityName,
		"fields", len(options.Fields))

	// First extract files normally
	baseGen := NewComponentGenerator(".")
	if err := baseGen.ExtractComponent(options); err != nil {
		return fmt.Errorf("failed to extract base command files: %w", err)
	}

	// Then modify the entity-specific command file to include custom fields
	entityFile := fmt.Sprintf("%s/commands/%s.go", options.FeatureName, options.EntityName)
	entityFilePath := filepath.Join(options.TargetDir, entityFile)

	// Read the generated file
	content, err := os.ReadFile(entityFilePath)
	if err != nil {
		return fmt.Errorf("failed to read generated command file %s: %w", entityFilePath, err)
	}

	// Modify the content using the editor
	modifiedContent, err := cg.modifyCommandStructWithFields(string(content), options)
	if err != nil {
		return fmt.Errorf("failed to modify command struct: %w", err)
	}

	// Write the modified content back
	if err := os.WriteFile(entityFilePath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("failed to write modified command file %s: %w", entityFilePath, err)
	}

	slog.Debug("Successfully modified command file with custom fields", "file", entityFilePath)
	return nil
}

// modifyCommandStructWithFields uses the editor to modify the command struct with custom fields
func (cg *CommandGenerator) modifyCommandStructWithFields(content string, options ExtractionOptions) (string, error) {
	// Build the field definitions string
	var fieldDefs []string
	for _, field := range options.Fields {
		// Capitalize first letter for exported fields
		capitalizedName := strings.ToUpper(field.Name[:1]) + field.Name[1:]
		fieldDef := fmt.Sprintf("\t%s %s", capitalizedName, field.Type)
		fieldDefs = append(fieldDefs, fieldDef)
	}
	fieldDefsStr := strings.Join(fieldDefs, "\n") + "\n"

	editor := ed.New(content)
	
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
		return "", fmt.Errorf("failed to modify command struct: %w", err)
	}

	// Also simplify the Validate method to just return nil
	content = editor.String()
	editor = ed.New(content)
	
	err = editor.Do(
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
		// If we can't find the Validate method, that's okay - not all templates have it
		slog.Debug("Could not simplify Validate method, this is expected for some templates")
	}

	// Also simplify the handler method to just return nil
	content = editor.String()
	editor = ed.New(content)
	
	err = editor.Do(
		ed.BeginningOfBuffer(),
		ed.Search("func (e *Executor) Handle"),
		ed.Search("{"),
		ed.ForwardChar(1), // Move past the opening brace
		ed.SetMark(),
		ed.Search("return nil"),
		ed.ForwardChar(10), // Move past "return nil"
		ed.ReplaceRegion("\n\treturn nil\n"),
	)
	
	if err != nil {
		// If we can't find the handler method, that's okay
		slog.Debug("Could not simplify handler method, this is expected for some templates")
	}

	return editor.String(), nil
}

// ValidateCommandStructure validates the generated command structure
func (cg *CommandGenerator) ValidateCommandStructure(featureName, entityName, targetDir string) error {
	slog.Debug("Validating command structure", 
		"feature", featureName, 
		"entity", entityName)

	// TODO: Implement command-specific validation
	// - Check if command files compile
	// - Check if command is properly registered
	// - Check if imports are correct
	// - Check if command follows naming conventions

	return nil
}

// getCommandFileList returns the list of skeleton files needed for command generation
func (cg *CommandGenerator) getCommandFileList(featureName, entityName string) ([]string, error) {
	// Get command file mapping from templates
	fileMap := templates.GetCommandTemplateFiles(entityName)
	
	// Extract source files and verify they exist
	var commandFiles []string
	for skeletonFile := range fileMap {
		if cg.skeletonFileExists(skeletonFile) {
			commandFiles = append(commandFiles, skeletonFile)
			slog.Debug("Found command skeleton file", "file", skeletonFile)
		} else {
			slog.Debug("Command skeleton file not found, skipping", 
				"file", skeletonFile, "entity", entityName)
		}
	}

	if len(commandFiles) == 0 {
		return nil, fmt.Errorf("no command skeleton files found for entity %s", entityName)
	}

	return commandFiles, nil
}



// skeletonFileExists checks if a file exists in the embedded skeleton
func (cg *CommandGenerator) skeletonFileExists(filePath string) bool {
	_, err := petrock.SkeletonFS.ReadFile(filePath)
	return err == nil
}


