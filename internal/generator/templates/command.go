package templates

import (
	"fmt"
	"strings"
)

// CommandFileMap defines the mapping from skeleton files to target files for commands
type CommandFileMap map[string]string

// CommandPlaceholders holds command-specific placeholder definitions
type CommandPlaceholders struct {
	FeatureName        string
	EntityName         string
	ModulePath         string
	CommandStructName  string
	CommandMethodName  string
	CommandPackagePath string
}

// GetCommandTemplateFiles returns the skeleton files needed for command generation
func GetCommandTemplateFiles(entityName string) CommandFileMap {
	baseFiles := CommandFileMap{
		"internal/skeleton/petrock_example_feature_name/commands/base.go":     "{{feature}}/commands/base.go",
		"internal/skeleton/petrock_example_feature_name/commands/register.go": "{{feature}}/commands/register.go",
	}

	// Convert kebab-case to snake_case for file matching
	normalizedEntityName := kebabToSnakeCase(entityName)
	
	// Try to find an exact match first
	knownEntities := []string{"create", "update", "delete", "get", "list", "request_summary", "set_summary", "fail_summary"}
	for _, knownEntity := range knownEntities {
		if normalizedEntityName == knownEntity {
			skeletonFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/commands/%s.go", knownEntity)
			targetFile := fmt.Sprintf("{{feature}}/commands/%s.go", normalizedEntityName)
			baseFiles[skeletonFile] = targetFile
			return baseFiles
		}
	}
	
	// If no exact match, use the create.go template as a base for new command
	baseFiles["internal/skeleton/petrock_example_feature_name/commands/create.go"] = 
		fmt.Sprintf("{{feature}}/commands/%s.go", normalizedEntityName)

	return baseFiles
}

// GetCommandReplacements returns placeholder replacements for command generation
func GetCommandReplacements(placeholders CommandPlaceholders) map[string]string {
	// Start with more specific replacements first
	replacements := map[string]string{}
	
	// Replace command path first (more specific)
	replacements["petrock_example_feature_name/create"] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityName)
	
	// Then add general replacements
	replacements["petrock_example_feature_name"] = placeholders.FeatureName
	replacements["github.com/petrock/example_module_path"] = placeholders.ModulePath
	replacements["{{feature}}"] = placeholders.FeatureName
	replacements["{{entity}}"] = placeholders.EntityName
	replacements["{{module_path}}"] = placeholders.ModulePath

	// Add command-specific replacements
	replacements[fmt.Sprintf("petrock_example_feature_name/%s", placeholders.EntityName)] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityName)

	// Add struct name replacements (e.g., CreateCommand -> SchedulePublicationCommand)
	if placeholders.CommandStructName != "" {
		replacements["CreateCommand"] = placeholders.CommandStructName
		replacements["{{command_struct}}"] = placeholders.CommandStructName
	}

	// Add method name replacements (e.g., HandleCreate -> HandleSchedulePublication)
	if placeholders.CommandMethodName != "" {
		replacements["HandleCreate"] = placeholders.CommandMethodName
		replacements["{{command_method}}"] = placeholders.CommandMethodName
	}

	// Also replace just the path part for CommandName method (after general replacement)
	replacements[fmt.Sprintf("%s/create", placeholders.FeatureName)] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityName)

	// Add package path replacements
	if placeholders.CommandPackagePath != "" {
		replacements["{{command_package}}"] = placeholders.CommandPackagePath
	}

	return replacements
}

// BuildCommandPlaceholders creates CommandPlaceholders from basic inputs
func BuildCommandPlaceholders(featureName, entityName, modulePath string) CommandPlaceholders {
	// Normalize entity name (convert kebab-case to snake_case)
	normalizedEntityName := kebabToSnakeCase(entityName)
	
	// Generate command struct name (e.g., "create" -> "CreateCommand", "schedule-publication" -> "SchedulePublicationCommand")
	commandStructName := toTitleCase(entityName) + "Command"
	
	// Generate command method name (e.g., "create" -> "HandleCreate", "schedule-publication" -> "HandleSchedulePublication")
	commandMethodName := "Handle" + toTitleCase(entityName)
	
	// Generate command package path
	commandPackagePath := fmt.Sprintf("%s/%s/commands", modulePath, featureName)

	return CommandPlaceholders{
		FeatureName:        featureName,
		EntityName:         normalizedEntityName,
		ModulePath:         modulePath,
		CommandStructName:  commandStructName,
		CommandMethodName:  commandMethodName,
		CommandPackagePath: commandPackagePath,
	}
}

// ResolveTargetPath resolves template placeholders in a target file path
func ResolveTargetPath(templatePath string, placeholders CommandPlaceholders) string {
	resolved := templatePath
	resolved = strings.ReplaceAll(resolved, "{{feature}}", placeholders.FeatureName)
	resolved = strings.ReplaceAll(resolved, "{{entity}}", placeholders.EntityName)
	resolved = strings.ReplaceAll(resolved, "{{module_path}}", placeholders.ModulePath)
	return resolved
}

// GetCommandDependencies returns additional files that command generation might need
func GetCommandDependencies(featureName string) []string {
	// Commands typically depend on state package
	return []string{
		fmt.Sprintf("%s/state", featureName),
	}
}

// ValidateCommandEntity checks if an entity name is valid for command generation
func ValidateCommandEntity(entityName string) error {
	return ValidateEntityName(entityName)
}
