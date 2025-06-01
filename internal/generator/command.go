package generator

import (
	"fmt"
	"log/slog"

	petrock "github.com/dhamidi/petrock"
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
		"entity", entityName)

	// Get command-specific file list
	commandFiles, err := cg.getCommandFileList(featureName, entityName)
	if err != nil {
		return fmt.Errorf("failed to get command file list: %w", err)
	}

	// Update extraction options with command files
	options.SkeletonFiles = commandFiles
	
	// Use base ComponentGenerator extraction logic
	baseGen := NewComponentGenerator(".")
	return baseGen.ExtractComponent(options)
}

// GenerateCommandComponent generates a complete command component
func (cg *CommandGenerator) GenerateCommandComponent(featureName, entityName, targetDir, modulePath string) error {
	slog.Debug("Generating command component", 
		"feature", featureName, 
		"entity", entityName,
		"target", targetDir)

	// Validate entity name
	if err := templates.ValidateCommandEntity(entityName); err != nil {
		return fmt.Errorf("invalid command entity: %w", err)
	}

	// Check for collisions
	exists, err := cg.inspector.ComponentExists(ComponentTypeCommand, featureName, entityName)
	if err != nil {
		slog.Warn("Could not check for existing commands", "error", err.Error())
	} else if exists {
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
	}

	// Extract command files
	return cg.ExtractCommandFiles(featureName, entityName, extractOptions)
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


