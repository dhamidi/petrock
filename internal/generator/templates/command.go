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

	// Add entity-specific file if it matches known patterns
	knownEntities := []string{"create", "update", "delete", "get", "list", "request_summary", "set_summary", "fail_summary"}
	for _, knownEntity := range knownEntities {
		if entityName == knownEntity {
			skeletonFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/commands/%s.go", entityName)
			targetFile := fmt.Sprintf("{{feature}}/commands/%s.go", entityName)
			baseFiles[skeletonFile] = targetFile
			break
		}
	}

	return baseFiles
}

// GetCommandReplacements returns placeholder replacements for command generation
func GetCommandReplacements(placeholders CommandPlaceholders) map[string]string {
	replacements := map[string]string{
		"petrock_example_feature_name":              placeholders.FeatureName,
		"github.com/petrock/example_module_path":    placeholders.ModulePath,
		"{{feature}}":                               placeholders.FeatureName,
		"{{entity}}":                                placeholders.EntityName,
		"{{module_path}}":                           placeholders.ModulePath,
	}

	// Add command-specific replacements
	replacements[fmt.Sprintf("petrock_example_feature_name/%s", placeholders.EntityName)] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityName)

	// Add struct name replacements (e.g., CreateCommand)
	if placeholders.CommandStructName != "" {
		replacements["{{command_struct}}"] = placeholders.CommandStructName
	}

	// Add method name replacements (e.g., HandleCreate)
	if placeholders.CommandMethodName != "" {
		replacements["{{command_method}}"] = placeholders.CommandMethodName
	}

	// Add package path replacements
	if placeholders.CommandPackagePath != "" {
		replacements["{{command_package}}"] = placeholders.CommandPackagePath
	}

	return replacements
}

// BuildCommandPlaceholders creates CommandPlaceholders from basic inputs
func BuildCommandPlaceholders(featureName, entityName, modulePath string) CommandPlaceholders {
	// Generate command struct name (e.g., "create" -> "CreateCommand")
	commandStructName := toTitleCase(entityName) + "Command"
	
	// Generate command method name (e.g., "create" -> "HandleCreate")
	commandMethodName := "Handle" + toTitleCase(entityName)
	
	// Generate command package path
	commandPackagePath := fmt.Sprintf("%s/%s/commands", modulePath, featureName)

	return CommandPlaceholders{
		FeatureName:        featureName,
		EntityName:         entityName,
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
