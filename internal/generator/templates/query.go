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
	EntityName       string // Keep as snake_case for file names
	EntityKebab      string // Add kebab-case version for query names
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

	// Convert kebab-case to snake_case for file matching
	normalizedEntityName := kebabToSnakeCase(entityName)

	// Try to find an exact match first
	knownEntities := []string{"get", "list", "search", "find", "count"}
	for _, knownEntity := range knownEntities {
		if normalizedEntityName == knownEntity {
			skeletonFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/queries/%s.go", knownEntity)
			targetFile := fmt.Sprintf("{{feature}}/queries/%s.go", normalizedEntityName)
			baseFiles[skeletonFile] = targetFile
			return baseFiles
		}
	}
	
	// If no exact match, use the get.go template as a base for new query
	baseFiles["internal/skeleton/petrock_example_feature_name/queries/get.go"] = 
		fmt.Sprintf("{{feature}}/queries/%s.go", normalizedEntityName)

	return baseFiles
}

// GetQueryReplacements returns placeholder replacements for query generation
func GetQueryReplacements(placeholders QueryPlaceholders) map[string]string {
	// Start with more specific replacements first
	replacements := map[string]string{}
	
	// Use kebab-case for query paths in QueryName() method
	replacements["petrock_example_feature_name/get"] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityKebab)
	
	// Then add general replacements
	replacements["petrock_example_feature_name"] = placeholders.FeatureName
	replacements["github.com/petrock/example_module_path"] = placeholders.ModulePath
	replacements["{{feature}}"] = placeholders.FeatureName
	replacements["{{entity}}"] = placeholders.EntityName // snake_case for files
	replacements["{{entity_kebab}}"] = placeholders.EntityKebab // kebab-case for names
	replacements["{{module_path}}"] = placeholders.ModulePath

	// Add query-specific replacements
	replacements[fmt.Sprintf("petrock_example_feature_name/%s", placeholders.EntityName)] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityKebab)

	// Add struct name replacements (e.g., GetQuery -> SearchPublishedQuery)
	if placeholders.QueryStructName != "" {
		replacements["GetQuery"] = placeholders.QueryStructName
		replacements["{{query_struct}}"] = placeholders.QueryStructName
	}

	// Add method name replacements (e.g., HandleGet -> HandleSearchPublished)
	if placeholders.QueryMethodName != "" {
		replacements["HandleGet"] = placeholders.QueryMethodName
		replacements["{{query_method}}"] = placeholders.QueryMethodName
	}

	// Add result name replacements (e.g., GetQueryResult -> SearchPublishedQueryResult)
	if placeholders.QueryResultName != "" {
		replacements["GetQueryResult"] = placeholders.QueryResultName
		replacements["{{query_result}}"] = placeholders.QueryResultName
	}

	// Also replace just the path part for QueryName method (after general replacement)
	replacements[fmt.Sprintf("%s/get", placeholders.FeatureName)] = 
		fmt.Sprintf("%s/%s", placeholders.FeatureName, placeholders.EntityKebab)

	// Add package path replacements
	if placeholders.QueryPackagePath != "" {
		replacements["{{query_package}}"] = placeholders.QueryPackagePath
	}

	return replacements
}

// BuildQueryPlaceholders creates QueryPlaceholders from basic inputs
func BuildQueryPlaceholders(featureName, entityName, modulePath string) QueryPlaceholders {
	// Normalize entity name (convert kebab-case to snake_case)
	normalizedEntityName := kebabToSnakeCase(entityName)
	
	// Generate query struct name (e.g., "get" -> "GetQuery", "search-posts" -> "SearchPostsQuery")
	queryStructName := toTitleCase(entityName) + "Query"
	
	// Generate query method name (e.g., "get" -> "HandleGet", "search-posts" -> "HandleSearchPosts")
	queryMethodName := "Handle" + toTitleCase(entityName)
	
	// Generate query result name (e.g., "get" -> "GetQueryResult", "search-posts" -> "SearchPostsQueryResult")
	queryResultName := toTitleCase(entityName) + "QueryResult"
	
	// Generate query package path
	queryPackagePath := fmt.Sprintf("%s/%s/queries", modulePath, featureName)

	return QueryPlaceholders{
		FeatureName:      featureName,
		EntityName:       normalizedEntityName,    // For file names
		EntityKebab:      entityName,              // Keep original kebab-case
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
