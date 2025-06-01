package templates

import (
	"fmt"
	"strings"
)

// WorkerFileMap defines the mapping from skeleton files to target files for workers
type WorkerFileMap map[string]string

// WorkerPlaceholders holds worker-specific placeholder definitions
type WorkerPlaceholders struct {
	FeatureName        string
	EntityName         string
	ModulePath         string
	WorkerStructName   string
	WorkerMethodName   string
	WorkerFileName     string
	WorkerPackagePath  string
}

// GetWorkerTemplateFiles returns the skeleton files needed for worker generation
func GetWorkerTemplateFiles(entityName string) WorkerFileMap {
	baseFiles := WorkerFileMap{
		"internal/skeleton/petrock_example_feature_name/workers/main.go":  "{{feature}}/workers/main.go",
		"internal/skeleton/petrock_example_feature_name/workers/types.go": "{{feature}}/workers/types.go",
	}

	// Add entity-specific file if it matches known patterns
	knownEntities := []string{"summary", "notification", "backup", "sync", "process", "analyze"}
	for _, knownEntity := range knownEntities {
		if entityName == knownEntity {
			skeletonFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/workers/%s_worker.go", entityName)
			targetFile := fmt.Sprintf("{{feature}}/workers/%s_worker.go", entityName)
			baseFiles[skeletonFile] = targetFile
			break
		}
	}

	return baseFiles
}

// GetWorkerReplacements returns placeholder replacements for worker generation
func GetWorkerReplacements(placeholders WorkerPlaceholders) map[string]string {
	replacements := map[string]string{
		"petrock_example_feature_name":              placeholders.FeatureName,
		"github.com/petrock/example_module_path":    placeholders.ModulePath,
		"{{feature}}":                               placeholders.FeatureName,
		"{{entity}}":                                placeholders.EntityName,
		"{{module_path}}":                           placeholders.ModulePath,
	}

	// Add worker-specific replacements
	replacements[fmt.Sprintf("petrock_example_feature_name/%s", placeholders.EntityName)] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityName)

	// Add struct name replacements (e.g., SummaryWorker)
	if placeholders.WorkerStructName != "" {
		replacements["{{worker_struct}}"] = placeholders.WorkerStructName
	}

	// Add method name replacements (e.g., ProcessSummary)
	if placeholders.WorkerMethodName != "" {
		replacements["{{worker_method}}"] = placeholders.WorkerMethodName
	}

	// Add file name replacements (e.g., summary_worker.go)
	if placeholders.WorkerFileName != "" {
		replacements["{{worker_file}}"] = placeholders.WorkerFileName
	}

	// Add package path replacements
	if placeholders.WorkerPackagePath != "" {
		replacements["{{worker_package}}"] = placeholders.WorkerPackagePath
	}

	return replacements
}

// BuildWorkerPlaceholders creates WorkerPlaceholders from basic inputs
func BuildWorkerPlaceholders(featureName, entityName, modulePath string) WorkerPlaceholders {
	// Generate worker struct name (e.g., "summary" -> "SummaryWorker")
	workerStructName := toTitleCase(entityName) + "Worker"
	
	// Generate worker method name (e.g., "summary" -> "ProcessSummary")
	workerMethodName := "Process" + toTitleCase(entityName)
	
	// Generate worker file name (e.g., "summary" -> "summary_worker.go")
	workerFileName := entityName + "_worker.go"
	
	// Generate worker package path
	workerPackagePath := fmt.Sprintf("%s/%s/workers", modulePath, featureName)

	return WorkerPlaceholders{
		FeatureName:        featureName,
		EntityName:         entityName,
		ModulePath:         modulePath,
		WorkerStructName:   workerStructName,
		WorkerMethodName:   workerMethodName,
		WorkerFileName:     workerFileName,
		WorkerPackagePath:  workerPackagePath,
	}
}

// ResolveWorkerTargetPath resolves template placeholders in a target file path
func ResolveWorkerTargetPath(templatePath string, placeholders WorkerPlaceholders) string {
	resolved := templatePath
	resolved = strings.ReplaceAll(resolved, "{{feature}}", placeholders.FeatureName)
	resolved = strings.ReplaceAll(resolved, "{{entity}}", placeholders.EntityName)
	resolved = strings.ReplaceAll(resolved, "{{module_path}}", placeholders.ModulePath)
	return resolved
}

// GetWorkerDependencies returns additional files that worker generation might need
func GetWorkerDependencies(featureName string) []string {
	// Workers typically depend on state and commands packages
	return []string{
		fmt.Sprintf("%s/state", featureName),
		fmt.Sprintf("%s/commands", featureName),
	}
}

// ValidateWorkerEntity checks if an entity name is valid for worker generation
func ValidateWorkerEntity(entityName string) error {
	return ValidateEntityName(entityName)
}

// GetWorkerPatterns returns common worker patterns with descriptions
func GetWorkerPatterns() map[string]string {
	return map[string]string{
		"summary":      "Generate summaries or reports from data",
		"notification": "Send notifications to users or external systems",
		"backup":       "Backup data to external storage",
		"sync":         "Synchronize data with external systems",
		"process":      "Process and transform data",
		"analyze":      "Analyze data and generate insights",
		"cleanup":      "Clean up old or unused data",
		"export":       "Export data to external formats",
		"import":       "Import data from external sources",
		"monitor":      "Monitor system health and performance",
	}
}
