package generator

import (
	"fmt"
	"log/slog"
	"strings"

	petrock "github.com/dhamidi/petrock"
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

	// Check for collisions
	exists, err := cg.inspector.ComponentExists(ComponentTypeCommand, featureName, entityName)
	if err != nil {
		slog.Warn("Could not check for existing commands", "error", err.Error())
	} else if exists {
		return fmt.Errorf("command %s/%s already exists", featureName, entityName)
	}

	// Prepare extraction options
	extractOptions := ExtractionOptions{
		ComponentType: ComponentTypeCommand,
		FeatureName:   featureName,
		EntityName:    entityName,
		TargetDir:     targetDir,
		Replacements:  cg.buildCommandReplacements(featureName, entityName, modulePath),
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
	commandFiles := []string{
		"internal/skeleton/petrock_example_feature_name/commands/base.go",
		"internal/skeleton/petrock_example_feature_name/commands/register.go",
	}

	// Check if entity-specific command file exists in skeleton
	entityFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/commands/%s.go", entityName)
	if cg.skeletonFileExists(entityFile) {
		commandFiles = append(commandFiles, entityFile)
		slog.Debug("Found entity-specific command file", "file", entityFile)
	} else {
		// Generate a generic command file based on naming patterns
		slog.Debug("Entity-specific command file not found, will use base patterns", 
			"entity", entityName, "expectedFile", entityFile)
		
		// For now, if the specific entity file doesn't exist, we don't generate it
		// In a more sophisticated implementation, we could generate from templates
	}

	return commandFiles, nil
}

// buildCommandReplacements creates command-specific placeholder replacements
func (cg *CommandGenerator) buildCommandReplacements(featureName, entityName, modulePath string) map[string]string {
	replacements := map[string]string{
		"petrock_example_feature_name": featureName,
		"github.com/petrock/example_module_path": modulePath,
	}

	// Add command-specific entity replacements
	// These handle entity name variations in different contexts
	replacements[fmt.Sprintf("petrock_example_feature_name/%s", entityName)] = fmt.Sprintf("%s/%s", featureName, entityName)
	
	// Handle command registration patterns
	if strings.Contains(entityName, "_") {
		// Convert snake_case to camelCase for struct names
		camelEntity := toCamelCase(entityName)
		replacements[fmt.Sprintf("%sCommand", strings.Title(entityName))] = fmt.Sprintf("%sCommand", camelEntity)
	}

	return replacements
}

// skeletonFileExists checks if a file exists in the embedded skeleton
func (cg *CommandGenerator) skeletonFileExists(filePath string) bool {
	_, err := petrock.SkeletonFS.ReadFile(filePath)
	return err == nil
}

// toCamelCase converts snake_case to CamelCase
func toCamelCase(input string) string {
	if input == "" {
		return ""
	}
	
	parts := strings.Split(input, "_")
	result := ""
	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}
	return result
}
