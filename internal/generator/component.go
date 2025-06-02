package generator

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	petrock "github.com/dhamidi/petrock"
	"github.com/dhamidi/petrock/internal/utils"
)

// ComponentGenerator defines the interface for component generation
type ComponentGenerator interface {
	ExtractComponent(options ExtractionOptions) error
	GenerateComponent(options GenerateOptions) error
	ValidateComponent(options ValidateOptions) error
}

// ExtractionOptions holds options for template extraction
type ExtractionOptions struct {
	ComponentType   ComponentType
	FeatureName     string
	EntityName      string
	TargetDir       string
	SkeletonFiles   []string
	Replacements    map[string]string
	FileMapping     map[string]string // Optional mapping from skeleton file to target file path
}

// GenerateOptions holds options for component generation
type GenerateOptions struct {
	ComponentType ComponentType
	FeatureName   string
	EntityName    string
	TargetDir     string
	ModulePath    string
}

// ValidateOptions holds options for component validation
type ValidateOptions struct {
	ComponentType ComponentType
	FeatureName   string
	EntityName    string
	TargetDir     string
}

// ComponentTemplate represents a template component structure
type ComponentTemplate struct {
	Type         ComponentType
	SourceFiles  []string
	TargetFiles  []string
	Replacements map[string]string
}

// ComponentGeneratorImpl implements ComponentGenerator
type ComponentGeneratorImpl struct {
	inspector ComponentInspector
}

// NewComponentGenerator creates a new component generator
func NewComponentGenerator(projectPath string) ComponentGenerator {
	return &ComponentGeneratorImpl{
		inspector: NewComponentInspector(projectPath),
	}
}

// ExtractComponent extracts component files from skeleton templates
func (cg *ComponentGeneratorImpl) ExtractComponent(options ExtractionOptions) error {
	slog.Debug("Extracting component", 
		"type", options.ComponentType,
		"feature", options.FeatureName,
		"entity", options.EntityName,
		"target", options.TargetDir,
		"skeletonFiles", options.SkeletonFiles)

	// Ensure target directory exists
	if err := utils.EnsureDir(options.TargetDir); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", options.TargetDir, err)
	}

	// Extract each skeleton file
	for i, sourceFile := range options.SkeletonFiles {
		slog.Debug("Extracting skeleton file", "index", i, "file", sourceFile)
		if err := cg.extractSkeletonFile(sourceFile, options); err != nil {
			return fmt.Errorf("failed to extract %s: %w", sourceFile, err)
		}
	}

	slog.Debug("Component extraction completed", 
		"type", options.ComponentType,
		"files", len(options.SkeletonFiles))

	return nil
}

// GenerateComponent generates a complete component with all required files
func (cg *ComponentGeneratorImpl) GenerateComponent(options GenerateOptions) error {
	slog.Debug("Generating component", 
		"type", options.ComponentType,
		"feature", options.FeatureName,
		"entity", options.EntityName)

	// Use specialized generators for each component type
	switch options.ComponentType {
	case ComponentTypeCommand:
		cmdGen := NewCommandGenerator(".")
		return cmdGen.GenerateCommandComponent(options.FeatureName, options.EntityName, options.TargetDir, options.ModulePath)
	case ComponentTypeQuery:
		queryGen := NewQueryGenerator(".")
		return queryGen.GenerateQueryComponent(options.FeatureName, options.EntityName, options.TargetDir, options.ModulePath)
	case ComponentTypeWorker:
		workerGen := NewWorkerGenerator(".")
		return workerGen.GenerateWorkerComponent(options.FeatureName, options.EntityName, options.TargetDir, options.ModulePath)
	default:
		return fmt.Errorf("unknown component type: %s", options.ComponentType)
	}
}

// generateGenericComponent provides fallback generic component generation
func (cg *ComponentGeneratorImpl) generateGenericComponent(options GenerateOptions) error {
	// Check for collisions
	exists, err := cg.inspector.ComponentExists(options.ComponentType, options.FeatureName, options.EntityName)
	if err != nil {
		slog.Warn("Could not check for existing components", "error", err.Error())
	} else if exists {
		return fmt.Errorf("component %s %s/%s already exists", 
			options.ComponentType, options.FeatureName, options.EntityName)
	}

	// Get component template
	template, err := cg.getComponentTemplate(options.ComponentType, options.FeatureName, options.EntityName)
	if err != nil {
		return fmt.Errorf("failed to get component template: %w", err)
	}

	// Prepare extraction options
	extractOptions := ExtractionOptions{
		ComponentType: options.ComponentType,
		FeatureName:   options.FeatureName,
		EntityName:    options.EntityName,
		TargetDir:     options.TargetDir,
		SkeletonFiles: template.SourceFiles,
		Replacements:  cg.buildReplacements(options.FeatureName, options.EntityName, options.ModulePath),
	}

	// Extract component files
	if err := cg.ExtractComponent(extractOptions); err != nil {
		return fmt.Errorf("failed to extract component: %w", err)
	}

	return nil
}

// ValidateComponent validates a generated component
func (cg *ComponentGeneratorImpl) ValidateComponent(options ValidateOptions) error {
	slog.Debug("Validating component", "type", options.ComponentType)
	
	// TODO: Implement validation logic
	// - Check if generated files compile
	// - Check if imports are correct
	// - Check if component registers correctly
	
	return nil
}

// extractSkeletonFile extracts a single skeleton file to target location
func (cg *ComponentGeneratorImpl) extractSkeletonFile(sourceFile string, options ExtractionOptions) error {
	// Read from embedded skeleton
	content, err := petrock.SkeletonFS.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read skeleton file %s: %w", sourceFile, err)
	}

	// Apply replacements in order (most specific first)
	contentStr := string(content)
	
	// First apply specific path replacements that might be affected by general replacements
	specificReplacements := []string{
		"petrock_example_feature_name/create",
		"petrock_example_feature_name/get", 
		"petrock_example_feature_name/summary",
	}
	
	for _, placeholder := range specificReplacements {
		if replacement, exists := options.Replacements[placeholder]; exists {
			contentStr = strings.ReplaceAll(contentStr, placeholder, replacement)
		}
	}
	
	// Then apply all other replacements
	for placeholder, replacement := range options.Replacements {
		// Skip the specific ones we already did
		isSpecific := false
		for _, specific := range specificReplacements {
			if placeholder == specific {
				isSpecific = true
				break
			}
		}
		if !isSpecific {
			contentStr = strings.ReplaceAll(contentStr, placeholder, replacement)
		}
	}

	// Determine target file path
	targetFile := cg.buildTargetPath(sourceFile, options)
	
	slog.Debug("Extracting skeleton file", "source", sourceFile, "target", targetFile)

	// Ensure target directory exists
	targetDir := filepath.Dir(targetFile)
	if err := utils.EnsureDir(targetDir); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	// Write target file
	if err := os.WriteFile(targetFile, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("failed to write target file %s: %w", targetFile, err)
	}

	return nil
}

// buildTargetPath builds the target file path from skeleton source path
func (cg *ComponentGeneratorImpl) buildTargetPath(sourceFile string, options ExtractionOptions) string {
	// Check if there's a specific mapping for this file
	if options.FileMapping != nil {
		if targetPath, exists := options.FileMapping[sourceFile]; exists {
			// Apply placeholder replacements to the target path
			resolved := targetPath
			for placeholder, replacement := range options.Replacements {
				resolved = strings.ReplaceAll(resolved, placeholder, replacement)
			}
			return filepath.Join(options.TargetDir, resolved)
		}
	}
	
	// Fall back to default path building
	relativePath := strings.TrimPrefix(sourceFile, "internal/skeleton/petrock_example_feature_name/")
	
	// Replace feature name placeholder with actual feature name
	relativePath = strings.ReplaceAll(relativePath, "petrock_example_feature_name", options.FeatureName)
	
	// Apply entity name replacements (if file contains entity references)
	// This will be refined in component-specific implementations
	
	return filepath.Join(options.TargetDir, relativePath)
}

// buildReplacements creates the replacement map for placeholders
func (cg *ComponentGeneratorImpl) buildReplacements(featureName, entityName, modulePath string) map[string]string {
	replacements := map[string]string{
		"petrock_example_feature_name": featureName,
		"github.com/petrock/example_module_path": modulePath,
	}
	
	// Add entity-specific replacements based on naming patterns
	// These will be refined in component-specific implementations
	
	return replacements
}

// getComponentTemplate returns the template for a specific component type
func (cg *ComponentGeneratorImpl) getComponentTemplate(componentType ComponentType, featureName, entityName string) (*ComponentTemplate, error) {
	switch componentType {
	case ComponentTypeCommand:
		return cg.getCommandTemplate(featureName, entityName)
	case ComponentTypeQuery:
		return cg.getQueryTemplate(featureName, entityName)
	case ComponentTypeWorker:
		return cg.getWorkerTemplate(featureName, entityName)
	default:
		return nil, fmt.Errorf("unknown component type: %s", componentType)
	}
}

// getCommandTemplate returns template for command components
func (cg *ComponentGeneratorImpl) getCommandTemplate(featureName, entityName string) (*ComponentTemplate, error) {
	sourceFiles := []string{
		"internal/skeleton/petrock_example_feature_name/commands/base.go",
		"internal/skeleton/petrock_example_feature_name/commands/register.go",
	}
	
	// Check if entity-specific command file exists in skeleton
	entityFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/commands/%s.go", entityName)
	if cg.skeletonFileExists(entityFile) {
		sourceFiles = append(sourceFiles, entityFile)
	} else {
		slog.Debug("Entity-specific command file not found in skeleton, will generate from base pattern", 
			"entity", entityName, "expectedFile", entityFile)
	}
	
	return &ComponentTemplate{
		Type:        ComponentTypeCommand,
		SourceFiles: sourceFiles,
		TargetFiles: nil, // Will be computed from source files
	}, nil
}

// getQueryTemplate returns template for query components  
func (cg *ComponentGeneratorImpl) getQueryTemplate(featureName, entityName string) (*ComponentTemplate, error) {
	sourceFiles := []string{
		"internal/skeleton/petrock_example_feature_name/queries/base.go",
	}
	
	// Check if entity-specific query file exists in skeleton
	entityFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/queries/%s.go", entityName)
	if cg.skeletonFileExists(entityFile) {
		sourceFiles = append(sourceFiles, entityFile)
	} else {
		slog.Debug("Entity-specific query file not found in skeleton, will generate from base pattern", 
			"entity", entityName, "expectedFile", entityFile)
	}
	
	return &ComponentTemplate{
		Type:        ComponentTypeQuery,
		SourceFiles: sourceFiles,
		TargetFiles: nil,
	}, nil
}

// getWorkerTemplate returns template for worker components
func (cg *ComponentGeneratorImpl) getWorkerTemplate(featureName, entityName string) (*ComponentTemplate, error) {
	sourceFiles := []string{
		"internal/skeleton/petrock_example_feature_name/workers/main.go",
		"internal/skeleton/petrock_example_feature_name/workers/types.go",
	}
	
	// Check if entity-specific worker file exists in skeleton  
	entityFile := fmt.Sprintf("internal/skeleton/petrock_example_feature_name/workers/%s_worker.go", entityName)
	if cg.skeletonFileExists(entityFile) {
		sourceFiles = append(sourceFiles, entityFile)
	} else {
		slog.Debug("Entity-specific worker file not found in skeleton, will generate from base pattern", 
			"entity", entityName, "expectedFile", entityFile)
	}
	
	return &ComponentTemplate{
		Type:        ComponentTypeWorker,
		SourceFiles: sourceFiles,
		TargetFiles: nil,
	}, nil
}

// skeletonFileExists checks if a file exists in the embedded skeleton
func (cg *ComponentGeneratorImpl) skeletonFileExists(filePath string) bool {
	_, err := petrock.SkeletonFS.ReadFile(filePath)
	return err == nil
}
