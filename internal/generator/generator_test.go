package generator

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dhamidi/petrock/internal/generator/templates"
)

// GeneratorTestSuite represents the test suite for component generators
type GeneratorTestSuite struct {
	t               *testing.T
	tempDir         string
	testModulePath  string
	testFeatureName string
}

// TestCase represents a single test case
type TestCase struct {
	Name           string
	ComponentType  ComponentType
	FeatureName    string
	EntityName     string
	ExpectedFiles  []string
	ShouldFail     bool
	ExpectedError  string
}

// NewGeneratorTestSuite creates a new test suite
func NewGeneratorTestSuite(t *testing.T) *GeneratorTestSuite {
	tempDir, err := os.MkdirTemp("", "petrock-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	return &GeneratorTestSuite{
		t:               t,
		tempDir:         tempDir,
		testModulePath:  "github.com/test/project",
		testFeatureName: "testfeature",
	}
}

// Cleanup removes the test directory
func (suite *GeneratorTestSuite) Cleanup() {
	if err := os.RemoveAll(suite.tempDir); err != nil {
		suite.t.Errorf("Failed to cleanup test directory: %v", err)
	}
}

// CreateTestGoMod creates a test go.mod file
func (suite *GeneratorTestSuite) CreateTestGoMod() error {
	goModContent := fmt.Sprintf("module %s\n\ngo 1.21\n", suite.testModulePath)
	goModPath := filepath.Join(suite.tempDir, "go.mod")
	return os.WriteFile(goModPath, []byte(goModContent), 0644)
}

// TestCommandGeneration tests command component generation
func TestCommandGeneration(t *testing.T) {
	suite := NewGeneratorTestSuite(t)
	defer suite.Cleanup()

	if err := suite.CreateTestGoMod(); err != nil {
		t.Fatalf("Failed to create test go.mod: %v", err)
	}

	testCases := []TestCase{
		{
			Name:          "Valid create command",
			ComponentType: ComponentTypeCommand,
			FeatureName:   "posts",
			EntityName:    "create",
			ExpectedFiles: []string{"posts/commands/base.go", "posts/commands/register.go", "posts/commands/create.go"},
			ShouldFail:    false,
		},
		{
			Name:          "Valid update command",
			ComponentType: ComponentTypeCommand,
			FeatureName:   "users",
			EntityName:    "update",
			ExpectedFiles: []string{"users/commands/base.go", "users/commands/register.go", "users/commands/update.go"},
			ShouldFail:    false,
		},
		{
			Name:          "Invalid entity name",
			ComponentType: ComponentTypeCommand,
			FeatureName:   "posts",
			EntityName:    "123invalid",
			ShouldFail:    true,
			ExpectedError: "invalid entity name",
		},
		{
			Name:          "Empty feature name",
			ComponentType: ComponentTypeCommand,
			FeatureName:   "",
			EntityName:    "create",
			ShouldFail:    true,
			ExpectedError: "feature name cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			suite.runTestCase(tc)
		})
	}
}

// TestQueryGeneration tests query component generation
func TestQueryGeneration(t *testing.T) {
	suite := NewGeneratorTestSuite(t)
	defer suite.Cleanup()

	if err := suite.CreateTestGoMod(); err != nil {
		t.Fatalf("Failed to create test go.mod: %v", err)
	}

	testCases := []TestCase{
		{
			Name:          "Valid get query",
			ComponentType: ComponentTypeQuery,
			FeatureName:   "posts",
			EntityName:    "get",
			ExpectedFiles: []string{"posts/queries/base.go", "posts/queries/get.go"},
			ShouldFail:    false,
		},
		{
			Name:          "Valid list query",
			ComponentType: ComponentTypeQuery,
			FeatureName:   "users",
			EntityName:    "list",
			ExpectedFiles: []string{"users/queries/base.go", "users/queries/list.go"},
			ShouldFail:    false,
		},
		{
			Name:          "Unknown entity query",
			ComponentType: ComponentTypeQuery,
			FeatureName:   "posts",
			EntityName:    "unknown",
			ExpectedFiles: []string{"posts/queries/base.go"}, // Only base file expected
			ShouldFail:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			suite.runTestCase(tc)
		})
	}
}

// TestWorkerGeneration tests worker component generation
func TestWorkerGeneration(t *testing.T) {
	suite := NewGeneratorTestSuite(t)
	defer suite.Cleanup()

	if err := suite.CreateTestGoMod(); err != nil {
		t.Fatalf("Failed to create test go.mod: %v", err)
	}

	testCases := []TestCase{
		{
			Name:          "Valid summary worker",
			ComponentType: ComponentTypeWorker,
			FeatureName:   "posts",
			EntityName:    "summary",
			ExpectedFiles: []string{"posts/workers/main.go", "posts/workers/types.go", "posts/workers/summary_worker.go"},
			ShouldFail:    false,
		},
		{
			Name:          "Valid notification worker",
			ComponentType: ComponentTypeWorker,
			FeatureName:   "users",
			EntityName:    "notification",
			ExpectedFiles: []string{"users/workers/main.go", "users/workers/types.go"},
			ShouldFail:    false,
		},
		{
			Name:          "Unknown entity worker",
			ComponentType: ComponentTypeWorker,
			FeatureName:   "orders",
			EntityName:    "unknown",
			ExpectedFiles: []string{"orders/workers/main.go", "orders/workers/types.go"},
			ShouldFail:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			suite.runTestCase(tc)
		})
	}
}

// TestCollisionDetection tests collision detection functionality
func TestCollisionDetection(t *testing.T) {
	suite := NewGeneratorTestSuite(t)
	defer suite.Cleanup()

	if err := suite.CreateTestGoMod(); err != nil {
		t.Fatalf("Failed to create test go.mod: %v", err)
	}

	// Create a mock inspector that simulates existing components
	inspector := &MockComponentInspector{
		existingComponents: map[string]bool{
			"command:posts:create": true,
			"query:posts:get":      true,
			"worker:posts:summary": true,
		},
	}

	// Test collision detection
	exists, err := inspector.ComponentExists(ComponentTypeCommand, "posts", "create")
	if err != nil {
		t.Fatalf("Collision detection failed: %v", err)
	}
	if !exists {
		t.Error("Expected collision detection to find existing command")
	}

	exists, err = inspector.ComponentExists(ComponentTypeCommand, "posts", "delete")
	if err != nil {
		t.Fatalf("Collision detection failed: %v", err)
	}
	if exists {
		t.Error("Expected collision detection to not find non-existing command")
	}
}

// TestPlaceholderReplacements tests placeholder replacement functionality
func TestPlaceholderReplacements(t *testing.T) {
	tests := []struct {
		name           string
		componentType  ComponentType
		featureName    string
		entityName     string
		modulePath     string
		expectedKeys   []string
	}{
		{
			name:          "Command placeholders",
			componentType: ComponentTypeCommand,
			featureName:   "posts",
			entityName:    "create",
			modulePath:    "github.com/test/project",
			expectedKeys: []string{
				"petrock_example_feature_name",
				"github.com/petrock/example_module_path",
				"posts/create",
			},
		},
		{
			name:          "Query placeholders",
			componentType: ComponentTypeQuery,
			featureName:   "users",
			entityName:    "get",
			modulePath:    "github.com/test/project",
			expectedKeys: []string{
				"petrock_example_feature_name",
				"github.com/petrock/example_module_path",
				"users/get",
			},
		},
		{
			name:          "Worker placeholders",
			componentType: ComponentTypeWorker,
			featureName:   "analytics",
			entityName:    "process",
			modulePath:    "github.com/test/project",
			expectedKeys: []string{
				"petrock_example_feature_name",
				"github.com/petrock/example_module_path",
				"analytics/process",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var replacements map[string]string

			switch tt.componentType {
			case ComponentTypeCommand:
				placeholders := templates.BuildCommandPlaceholders(tt.featureName, tt.entityName, tt.modulePath)
				replacements = templates.GetCommandReplacements(placeholders)
			case ComponentTypeQuery:
				placeholders := templates.BuildQueryPlaceholders(tt.featureName, tt.entityName, tt.modulePath)
				replacements = templates.GetQueryReplacements(placeholders)
			case ComponentTypeWorker:
				placeholders := templates.BuildWorkerPlaceholders(tt.featureName, tt.entityName, tt.modulePath)
				replacements = templates.GetWorkerReplacements(placeholders)
			}

			for _, expectedKey := range tt.expectedKeys {
				if _, exists := replacements[expectedKey]; !exists {
					t.Errorf("Expected replacement key %q not found", expectedKey)
				}
			}

			// Verify that feature name replacement is correct
			if replacements["petrock_example_feature_name"] != tt.featureName {
				t.Errorf("Expected feature name %q, got %q", tt.featureName, replacements["petrock_example_feature_name"])
			}

			// Verify that module path replacement is correct
			if replacements["github.com/petrock/example_module_path"] != tt.modulePath {
				t.Errorf("Expected module path %q, got %q", tt.modulePath, replacements["github.com/petrock/example_module_path"])
			}
		})
	}
}

// runTestCase runs a single test case
func (suite *GeneratorTestSuite) runTestCase(tc TestCase) {
	// Change to test directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(suite.tempDir)

	var err error
	switch tc.ComponentType {
	case ComponentTypeCommand:
		cmdGen := NewCommandGenerator(".")
		err = cmdGen.GenerateCommandComponent(tc.FeatureName, tc.EntityName, ".", suite.testModulePath)
	case ComponentTypeQuery:
		queryGen := NewQueryGenerator(".")
		err = queryGen.GenerateQueryComponent(tc.FeatureName, tc.EntityName, ".", suite.testModulePath)
	case ComponentTypeWorker:
		workerGen := NewWorkerGenerator(".")
		err = workerGen.GenerateWorkerComponent(tc.FeatureName, tc.EntityName, ".", suite.testModulePath)
	}

	if tc.ShouldFail {
		if err == nil {
			suite.t.Errorf("Test case %q should have failed but didn't", tc.Name)
			return
		}
		if tc.ExpectedError != "" && !strings.Contains(err.Error(), tc.ExpectedError) {
			suite.t.Errorf("Test case %q expected error containing %q, got %q", tc.Name, tc.ExpectedError, err.Error())
		}
		return
	}

	if err != nil {
		suite.t.Errorf("Test case %q failed: %v", tc.Name, err)
		return
	}

	// Check that expected files were created
	for _, expectedFile := range tc.ExpectedFiles {
		fullPath := filepath.Join(suite.tempDir, expectedFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			suite.t.Errorf("Test case %q: expected file %q was not created", tc.Name, expectedFile)
		}
	}

	// Check that generated files contain correct replacements
	suite.verifyPlaceholderReplacements(tc)
}

// verifyPlaceholderReplacements checks that placeholders were correctly replaced
func (suite *GeneratorTestSuite) verifyPlaceholderReplacements(tc TestCase) {
	err := filepath.WalkDir(suite.tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := string(content)

		// Check that old placeholders don't exist
		if strings.Contains(contentStr, "petrock_example_feature_name") {
			suite.t.Errorf("File %q still contains unreplaced placeholder 'petrock_example_feature_name'", path)
		}
		if strings.Contains(contentStr, "github.com/petrock/example_module_path") {
			suite.t.Errorf("File %q still contains unreplaced placeholder 'github.com/petrock/example_module_path'", path)
		}

		// Check that new values exist
		if !strings.Contains(contentStr, tc.FeatureName) {
			suite.t.Errorf("File %q doesn't contain feature name %q", path, tc.FeatureName)
		}
		if !strings.Contains(contentStr, suite.testModulePath) {
			suite.t.Errorf("File %q doesn't contain module path %q", path, suite.testModulePath)
		}

		return nil
	})

	if err != nil {
		suite.t.Errorf("Error verifying placeholder replacements: %v", err)
	}
}

// MockComponentInspector mocks the ComponentInspector for testing
type MockComponentInspector struct {
	existingComponents map[string]bool
}

// InspectExistingComponents mocks the self inspect functionality
func (m *MockComponentInspector) InspectExistingComponents() (*InspectResult, error) {
	result := &InspectResult{
		Commands: make(map[string][]string),
		Queries:  make(map[string][]string),
		Workers:  make(map[string][]string),
		Routes:   []string{},
	}

	for key := range m.existingComponents {
		parts := strings.Split(key, ":")
		if len(parts) != 3 {
			continue
		}
		componentType, feature, entity := parts[0], parts[1], parts[2]

		switch componentType {
		case "command":
			if result.Commands[feature] == nil {
				result.Commands[feature] = []string{}
			}
			result.Commands[feature] = append(result.Commands[feature], entity)
		case "query":
			if result.Queries[feature] == nil {
				result.Queries[feature] = []string{}
			}
			result.Queries[feature] = append(result.Queries[feature], entity)
		case "worker":
			if result.Workers[feature] == nil {
				result.Workers[feature] = []string{}
			}
			result.Workers[feature] = append(result.Workers[feature], entity)
		}
	}

	return result, nil
}

// ComponentExists checks if a component exists in the mock
func (m *MockComponentInspector) ComponentExists(componentType ComponentType, featureName, entityName string) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", componentType, featureName, entityName)
	return m.existingComponents[key], nil
}

// TestInspectResultParsing tests JSON parsing of inspect results
func TestInspectResultParsing(t *testing.T) {
	jsonData := `{
		"commands": {
			"posts": ["create", "update", "delete"],
			"users": ["register", "login"]
		},
		"queries": {
			"posts": ["get", "list"],
			"users": ["get"]
		},
		"workers": {
			"posts": ["summary"],
			"notifications": ["email"]
		},
		"routes": ["/", "/posts", "/users"]
	}`

	var result InspectResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		t.Fatalf("Failed to parse inspect result JSON: %v", err)
	}

	// Verify commands
	if len(result.Commands["posts"]) != 3 {
		t.Errorf("Expected 3 commands for posts, got %d", len(result.Commands["posts"]))
	}

	// Verify queries
	if len(result.Queries["posts"]) != 2 {
		t.Errorf("Expected 2 queries for posts, got %d", len(result.Queries["posts"]))
	}

	// Verify workers
	if len(result.Workers["posts"]) != 1 {
		t.Errorf("Expected 1 worker for posts, got %d", len(result.Workers["posts"]))
	}

	// Verify routes
	if len(result.Routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(result.Routes))
	}
}
