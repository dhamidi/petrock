package generator

import (
	"fmt"
	"log/slog"

	petrock "github.com/dhamidi/petrock"
	"github.com/dhamidi/petrock/internal/generator/templates"
)

// QueryGenerator implements ComponentGenerator for query-specific generation
type QueryGenerator struct {
	inspector ComponentInspector
}

// NewQueryGenerator creates a new query-specific generator
func NewQueryGenerator(projectPath string) *QueryGenerator {
	return &QueryGenerator{
		inspector: NewComponentInspector(projectPath),
	}
}

// ExtractQueryFiles extracts query-specific files from skeleton
func (qg *QueryGenerator) ExtractQueryFiles(featureName, entityName string, options ExtractionOptions) error {
	slog.Debug("Extracting query files",
	"feature", featureName,
	"name", entityName)

	// Get query-specific file list
	queryFiles, err := qg.getQueryFileList(featureName, entityName)
	if err != nil {
		return fmt.Errorf("failed to get query file list: %w", err)
	}

	// Update extraction options with query files and file mapping
	options.SkeletonFiles = queryFiles
	options.FileMapping = templates.GetQueryTemplateFiles(entityName)
	
	// Use base ComponentGenerator extraction logic
	baseGen := NewComponentGenerator(".")
	return baseGen.ExtractComponent(options)
}

// GenerateQueryComponent generates a complete query component
func (qg *QueryGenerator) GenerateQueryComponent(featureName, entityName, targetDir, modulePath string) error {
	slog.Debug("Generating query component", 
		"feature", featureName, 
		"name", entityName,
		"target", targetDir)

	// Validate entity name  
	if err := templates.ValidateQueryEntity(entityName); err != nil {
		return fmt.Errorf("invalid query entity: %w", err)
	}

	// Check for collisions
	exists, err := qg.inspector.ComponentExists(ComponentTypeQuery, featureName, entityName)
	if err != nil {
		slog.Warn("Could not check for existing queries", "error", err.Error())
	} else if exists {
		return fmt.Errorf("query %s/%s already exists", featureName, entityName)
	}

	// Build query placeholders
	placeholders := templates.BuildQueryPlaceholders(featureName, entityName, modulePath)
	
	// Prepare extraction options
	extractOptions := ExtractionOptions{
		ComponentType: ComponentTypeQuery,
		FeatureName:   featureName,
		EntityName:    entityName,
		TargetDir:     targetDir,
		Replacements:  templates.GetQueryReplacements(placeholders),
	}

	// Extract query files
	return qg.ExtractQueryFiles(featureName, entityName, extractOptions)
}

// ValidateQueryStructure validates the generated query structure
func (qg *QueryGenerator) ValidateQueryStructure(featureName, entityName, targetDir string) error {
	slog.Debug("Validating query structure", 
		"feature", featureName, 
		"entity", entityName)

	// TODO: Implement query-specific validation
	// - Check if query files compile
	// - Check if query is properly registered
	// - Check if imports are correct
	// - Check if query follows naming conventions
	// - Check if query implements required interfaces

	return nil
}

// getQueryFileList returns the list of skeleton files needed for query generation
func (qg *QueryGenerator) getQueryFileList(featureName, entityName string) ([]string, error) {
	// Get query file mapping from templates
	fileMap := templates.GetQueryTemplateFiles(entityName)
	
	// Extract source files and verify they exist
	var queryFiles []string
	for skeletonFile := range fileMap {
		if qg.skeletonFileExists(skeletonFile) {
			queryFiles = append(queryFiles, skeletonFile)
			slog.Debug("Found query skeleton file", "file", skeletonFile)
		} else {
			slog.Debug("Query skeleton file not found, skipping", 
				"file", skeletonFile, "entity", entityName)
		}
	}

	if len(queryFiles) == 0 {
		return nil, fmt.Errorf("no query skeleton files found for entity %s", entityName)
	}

	return queryFiles, nil
}

// skeletonFileExists checks if a file exists in the embedded skeleton
func (qg *QueryGenerator) skeletonFileExists(filePath string) bool {
	_, err := petrock.SkeletonFS.ReadFile(filePath)
	return err == nil
}
