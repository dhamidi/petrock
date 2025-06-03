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
	
	// If we have custom fields, use enhanced extraction with template modification
	if len(options.QueryFields) > 0 {
		return qg.ExtractQueryFilesWithFields(options)
	}

	// Use base ComponentGenerator extraction logic
	baseGen := NewComponentGenerator(".")
	return baseGen.ExtractComponent(options)
}

// GenerateQueryComponent generates a complete query component
func (qg *QueryGenerator) GenerateQueryComponent(featureName, entityName, targetDir, modulePath string) error {
	return qg.GenerateQueryComponentWithFields(featureName, entityName, targetDir, modulePath, nil)
}

// GenerateQueryComponentWithFields generates a complete query component with custom fields
func (qg *QueryGenerator) GenerateQueryComponentWithFields(featureName, entityName, targetDir, modulePath string, fields []QueryField) error {
	slog.Debug("Generating query component", 
		"feature", featureName, 
		"name", entityName,
		"target", targetDir,
		"fields", len(fields))

	// Validate entity name  
	if err := templates.ValidateQueryEntity(entityName); err != nil {
		return fmt.Errorf("invalid query entity: %w", err)
	}

	// Check for collisions
	exists, err := qg.inspector.ComponentExists(ComponentTypeQuery, featureName, entityName)
	if err != nil {
		slog.Warn("Could not check for existing queries", "error", err.Error())
	}
	if exists {
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
		QueryFields:   fields,
	}

	// Extract query files
	if err := qg.ExtractQueryFiles(featureName, entityName, extractOptions); err != nil {
		return err
	}

	// Register the query in main.go
	return qg.registerQueryInMainFile(featureName, entityName, targetDir, modulePath)
}

// ExtractQueryFilesWithFields extracts query files and modifies them with custom fields using the editor
func (qg *QueryGenerator) ExtractQueryFilesWithFields(options ExtractionOptions) error {
	slog.Debug("Extracting query files with custom fields",
		"feature", options.FeatureName,
		"entity", options.EntityName,
		"fields", len(options.QueryFields))

	// First extract files normally
	baseGen := NewComponentGenerator(".")
	if err := baseGen.ExtractComponent(options); err != nil {
		return fmt.Errorf("failed to extract base query files: %w", err)
	}

	// Then modify the entity-specific query file to include custom fields
	entityFile := fmt.Sprintf("%s/queries/%s.go", options.FeatureName, options.EntityName)
	entityFilePath := filepath.Join(options.TargetDir, entityFile)

	// Read the generated file
	content, err := os.ReadFile(entityFilePath)
	if err != nil {
		return fmt.Errorf("failed to read generated query file %s: %w", entityFilePath, err)
	}

	// Modify the content using the editor
	modifiedContent, err := qg.modifyQueryStructWithFields(string(content), options)
	if err != nil {
		return fmt.Errorf("failed to modify query struct: %w", err)
	}

	// Write the modified content back
	if err := os.WriteFile(entityFilePath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("failed to write modified query file %s: %w", entityFilePath, err)
	}

	// Also modify the base.go file to update ItemResult with custom fields
	baseFilePath := filepath.Join(options.TargetDir, options.FeatureName, "queries", "base.go")
	baseContent, err := os.ReadFile(baseFilePath)
	if err != nil {
		return fmt.Errorf("failed to read generated base file %s: %w", baseFilePath, err)
	}

	// Modify the ItemResult struct
	modifiedBaseContent, err := qg.modifyItemResultWithFields(string(baseContent), options)
	if err != nil {
		return fmt.Errorf("failed to modify ItemResult struct: %w", err)
	}

	// Write the modified base content back
	if err := os.WriteFile(baseFilePath, []byte(modifiedBaseContent), 0644); err != nil {
		return fmt.Errorf("failed to write modified base file %s: %w", baseFilePath, err)
	}

	slog.Debug("Successfully modified query files with custom fields", "entityFile", entityFilePath, "baseFile", baseFilePath)
	return nil
}

// modifyQueryStructWithFields uses the editor to modify the query struct with custom fields
func (qg *QueryGenerator) modifyQueryStructWithFields(content string, options ExtractionOptions) (string, error) {
	// Build the field definitions string
	var fieldDefs []string
	for _, field := range options.QueryFields {
		// Capitalize first letter for exported fields
		capitalizedName := strings.ToUpper(field.Name[:1]) + field.Name[1:]
		fieldDef := fmt.Sprintf("\t%s %s `json:\"%s\" validate:\"required\"`", capitalizedName, field.Type, field.Name)
		fieldDefs = append(fieldDefs, fieldDef)
	}
	fieldDefsStr := strings.Join(fieldDefs, "\n") + "\n"

	editor := ed.New(content)

	// Find the query struct and replace its fields
	err := editor.Do(
		ed.BeginningOfBuffer(),
		ed.Search("type"),
		ed.Search("Query struct {"),
		ed.Search("{"),
		ed.ForwardChar(1), // Move past the opening brace
		ed.SetMark(),
		ed.Search("}"),     // Find closing brace
		ed.ReplaceRegion("\n"+fieldDefsStr),
	)

	if err != nil {
		return "", fmt.Errorf("failed to modify query struct: %w", err)
	}

	// Simplify the handler method to just return nil
	content = editor.String()
	editor = ed.New(content)

	err = editor.Do(
		ed.BeginningOfBuffer(),
		ed.Search("func (q *Querier) Handle"),
		ed.Search("{"),
		ed.ForwardChar(1), // Move past the opening brace
		ed.SetMark(),
		ed.Search("return result, nil"),
		ed.ForwardChar(17), // Move past "return result, nil"
		ed.ReplaceRegion("\n\treturn nil, nil\n"),
	)

	if err != nil {
		// If we can't find the handler method, that's okay
		slog.Debug("Could not simplify handler method, this is expected for some templates")
	}

	return editor.String(), nil
}

// modifyItemResultWithFields uses the editor to modify the ItemResult struct with custom fields
func (qg *QueryGenerator) modifyItemResultWithFields(content string, options ExtractionOptions) (string, error) {
	// Build the field definitions string
	var fieldDefs []string
	for _, field := range options.QueryFields {
		// Capitalize first letter for exported fields
		capitalizedName := strings.ToUpper(field.Name[:1]) + field.Name[1:]
		fieldDef := fmt.Sprintf("\t%s %s `json:\"%s\"`", capitalizedName, field.Type, field.Name)
		fieldDefs = append(fieldDefs, fieldDef)
	}
	fieldDefsStr := strings.Join(fieldDefs, "\n") + "\n"

	editor := ed.New(content)

	// Find the ItemResult struct and replace its fields
	err := editor.Do(
		ed.BeginningOfBuffer(),
		ed.Search("type ItemResult struct {"),
		ed.Search("{"),
		ed.ForwardChar(1), // Move past the opening brace
		ed.SetMark(),
		ed.Search("}"),     // Find closing brace
		ed.ReplaceRegion("\n"+fieldDefsStr),
	)

	if err != nil {
		return "", fmt.Errorf("failed to modify ItemResult struct: %w", err)
	}

	return editor.String(), nil
}

// registerQueryInMainFile adds the registration line for a new query to the feature's main.go file
func (qg *QueryGenerator) registerQueryInMainFile(featureName, entityName, targetDir, modulePath string) error {
	slog.Debug("Registering query in main.go", 
		"feature", featureName, 
		"entity", entityName)

	// Build the query placeholders
	placeholders := templates.BuildQueryPlaceholders(featureName, entityName, modulePath)
	
	// Path to the main.go file
	mainFilePath := filepath.Join(targetDir, featureName, "main.go")
	
	// Read the main.go file
	content, err := os.ReadFile(mainFilePath)
	if err != nil {
		return fmt.Errorf("failed to read main.go file %s: %w", mainFilePath, err)
	}

	// Generate the registration line
	registrationLine := fmt.Sprintf("\tapp.QueryRegistry.Register(queries.%s{}, featureQuerier.%s)", 
		placeholders.QueryStructName, placeholders.QueryMethodName)

	editor := ed.New(string(content))
	
	// Find the end of the query registration block by looking for the comment about message types
	err = editor.Do(
		ed.BeginningOfBuffer(),
		ed.Search("// --- 6. Register Message Types for Decoding ---"),
		ed.SetMark(),
		ed.ForwardChar(0), // Stay at current position
		ed.ReplaceRegion(registrationLine+"\n\n\t"),
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert query registration: %w", err)
	}

	// Write the modified content back
	if err := os.WriteFile(mainFilePath, []byte(editor.String()), 0644); err != nil {
		return fmt.Errorf("failed to write modified main.go file %s: %w", mainFilePath, err)
	}

	slog.Debug("Successfully registered query in main.go", "file", mainFilePath, "query", placeholders.QueryStructName)
	return nil
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
