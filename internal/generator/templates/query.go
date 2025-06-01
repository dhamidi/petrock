package templates

import (
	"fmt"
	"strings"
)

// QueryFileMap defines the mapping from skeleton files to target files for queries
type QueryFileMap map[string]string

// QueryPlaceholders holds query-specific placeholder definitions
type QueryPlaceholders struct {
	FeatureName      string
	EntityName       string
	ModulePath       string
	QueryStructName  string
	QueryMethodName  string
	QueryResultName  string
	QueryPackagePath string
}

// GetQueryTemplateFiles returns the skeleton files needed for query generation
func GetQueryTemplateFiles(entityName string) QueryFileMap {
	baseFiles := QueryFileMap{
		"internal/skeleton/petrock_example_feature_name/queries/base.go": "{{feature}}/queries/base.go",
	}

	// Add entity-specific file if it matches known patterns
	knownEntities := []string{"get", "list", "search", "find", "count"}
	for _, knownEntity := range knownEntities {
		if entityName == knownEntity {
			skeletonFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/queries/%s.go", entityName)
			targetFile := fmt.Sprintf("{{feature}}/queries/%s.go", entityName)
			baseFiles[skeletonFile] = targetFile
			break
		}
	}

	return baseFiles
}

// GetQueryReplacements returns placeholder replacements for query generation
func GetQueryReplacements(placeholders QueryPlaceholders) map[string]string {
	replacements := map[string]string{
		"petrock_example_feature_name":              placeholders.FeatureName,
		"github.com/petrock/example_module_path":    placeholders.ModulePath,
		"{{feature}}":                               placeholders.FeatureName,
		"{{entity}}":                                placeholders.EntityName,
		"{{module_path}}":                           placeholders.ModulePath,
	}

	// Add query-specific replacements
	replacements[fmt.Sprintf("petrock_example_feature_name/%s", placeholders.EntityName)] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityName)

	// Add struct name replacements (e.g., GetQuery)
	if placeholders.QueryStructName != "" {
		replacements["{{query_struct}}"] = placeholders.QueryStructName
	}

	// Add method name replacements (e.g., HandleGet)
	if placeholders.QueryMethodName != "" {
		replacements["{{query_method}}"] = placeholders.QueryMethodName
	}

	// Add result name replacements (e.g., GetQueryResult)
	if placeholders.QueryResultName != "" {
		replacements["{{query_result}}"] = placeholders.QueryResultName
	}

	// Add package path replacements
	if placeholders.QueryPackagePath != "" {
		replacements["{{query_package}}"] = placeholders.QueryPackagePath
	}

	return replacements
}

// BuildQueryPlaceholders creates QueryPlaceholders from basic inputs
func BuildQueryPlaceholders(featureName, entityName, modulePath string) QueryPlaceholders {
	// Generate query struct name (e.g., "get" -> "GetQuery")
	queryStructName := toTitleCase(entityName) + "Query"
	
	// Generate query method name (e.g., "get" -> "HandleGet")
	queryMethodName := "Handle" + toTitleCase(entityName)
	
	// Generate query result name (e.g., "get" -> "GetQueryResult")
	queryResultName := toTitleCase(entityName) + "QueryResult"
	
	// Generate query package path
	queryPackagePath := fmt.Sprintf("%s/%s/queries", modulePath, featureName)

	return QueryPlaceholders{
		FeatureName:      featureName,
		EntityName:       entityName,
		ModulePath:       modulePath,
		QueryStructName:  queryStructName,
		QueryMethodName:  queryMethodName,
		QueryResultName:  queryResultName,
		QueryPackagePath: queryPackagePath,
	}
}

// ResolveQueryTargetPath resolves template placeholders in a target file path
func ResolveQueryTargetPath(templatePath string, placeholders QueryPlaceholders) string {
	resolved := templatePath
	resolved = strings.ReplaceAll(resolved, "{{feature}}", placeholders.FeatureName)
	resolved = strings.ReplaceAll(resolved, "{{entity}}", placeholders.EntityName)
	resolved = strings.ReplaceAll(resolved, "{{module_path}}", placeholders.ModulePath)
	return resolved
}

// GetQueryDependencies returns additional files that query generation might need
func GetQueryDependencies(featureName string) []string {
	// Queries typically depend on state package
	return []string{
		fmt.Sprintf("%s/state", featureName),
	}
}

// ValidateQueryEntity checks if an entity name is valid for query generation
func ValidateQueryEntity(entityName string) error {
	return ValidateEntityName(entityName)
}
