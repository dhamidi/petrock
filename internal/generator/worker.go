package generator

import (
	"fmt"
	"log/slog"

	petrock "github.com/dhamidi/petrock"
	"github.com/dhamidi/petrock/internal/generator/templates"
)

// WorkerGenerator implements ComponentGenerator for worker-specific generation
type WorkerGenerator struct {
	inspector ComponentInspector
}

// NewWorkerGenerator creates a new worker-specific generator
func NewWorkerGenerator(projectPath string) *WorkerGenerator {
	return &WorkerGenerator{
		inspector: NewComponentInspector(projectPath),
	}
}

// ExtractWorkerFiles extracts worker-specific files from skeleton
func (wg *WorkerGenerator) ExtractWorkerFiles(featureName, entityName string, options ExtractionOptions) error {
	slog.Debug("Extracting worker files", 
		"feature", featureName, 
		"entity", entityName)

	// Get worker-specific file list
	workerFiles, err := wg.getWorkerFileList(featureName, entityName)
	if err != nil {
		return fmt.Errorf("failed to get worker file list: %w", err)
	}

	// Update extraction options with worker files and file mapping
	options.SkeletonFiles = workerFiles
	options.FileMapping = templates.GetWorkerTemplateFiles(entityName)
	
	// Use base ComponentGenerator extraction logic
	baseGen := NewComponentGenerator(".")
	return baseGen.ExtractComponent(options)
}

// GenerateWorkerComponent generates a complete worker component
func (wg *WorkerGenerator) GenerateWorkerComponent(featureName, entityName, targetDir, modulePath string) error {
	slog.Debug("Generating worker component", 
		"feature", featureName, 
		"entity", entityName,
		"target", targetDir)

	// Validate entity name  
	if err := templates.ValidateWorkerEntity(entityName); err != nil {
		return fmt.Errorf("invalid worker entity: %w", err)
	}

	// Check for collisions
	exists, err := wg.inspector.ComponentExists(ComponentTypeWorker, featureName, entityName)
	if err != nil {
		slog.Warn("Could not check for existing workers", "error", err.Error())
	} else if exists {
		return fmt.Errorf("worker %s/%s already exists", featureName, entityName)
	}

	// Build worker placeholders
	placeholders := templates.BuildWorkerPlaceholders(featureName, entityName, modulePath)
	
	// Prepare extraction options
	extractOptions := ExtractionOptions{
		ComponentType: ComponentTypeWorker,
		FeatureName:   featureName,
		EntityName:    entityName,
		TargetDir:     targetDir,
		Replacements:  templates.GetWorkerReplacements(placeholders),
	}

	// Extract worker files
	return wg.ExtractWorkerFiles(featureName, entityName, extractOptions)
}

// ValidateWorkerStructure validates the generated worker structure
func (wg *WorkerGenerator) ValidateWorkerStructure(featureName, entityName, targetDir string) error {
	slog.Debug("Validating worker structure", 
		"feature", featureName, 
		"entity", entityName)

	// TODO: Implement worker-specific validation
	// - Check if worker files compile
	// - Check if worker is properly registered  
	// - Check if imports are correct
	// - Check if worker follows naming conventions
	// - Check if worker implements required interfaces
	// - Check if worker handles commands correctly

	return nil
}

// getWorkerFileList returns the list of skeleton files needed for worker generation
func (wg *WorkerGenerator) getWorkerFileList(featureName, entityName string) ([]string, error) {
	// Get worker file mapping from templates
	fileMap := templates.GetWorkerTemplateFiles(entityName)
	
	// Extract source files and verify they exist
	var workerFiles []string
	for skeletonFile := range fileMap {
		if wg.skeletonFileExists(skeletonFile) {
			workerFiles = append(workerFiles, skeletonFile)
			slog.Debug("Found worker skeleton file", "file", skeletonFile)
		} else {
			slog.Debug("Worker skeleton file not found, skipping", 
				"file", skeletonFile, "entity", entityName)
		}
	}

	if len(workerFiles) == 0 {
		return nil, fmt.Errorf("no worker skeleton files found for entity %s", entityName)
	}

	return workerFiles, nil
}

// skeletonFileExists checks if a file exists in the embedded skeleton
func (wg *WorkerGenerator) skeletonFileExists(filePath string) bool {
	_, err := petrock.SkeletonFS.ReadFile(filePath)
	return err == nil
}
