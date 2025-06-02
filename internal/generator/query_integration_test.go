package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQueryGenerationWithFields(t *testing.T) {
	// Skip if running with regular tests to avoid filesystem conflicts
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "petrock_query_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up test fields
	fields := []QueryField{
		{Name: "postID", Type: "string"},
		{Name: "status", Type: "string"},
		{Name: "publishedAfter", Type: "time.Time"},
	}

	// Create query generator
	queryGen := NewQueryGenerator(tempDir)

	// Generate query component with fields
	err = queryGen.GenerateQueryComponentWithFields(
		"posts",     // featureName
		"search",    // entityName
		tempDir,     // targetDir
		"github.com/test/app", // modulePath
		fields,
	)

	if err != nil {
		t.Fatalf("failed to generate query: %v", err)
	}

	// Read the generated query file
	queryFile := filepath.Join(tempDir, "posts", "queries", "search.go")
	content, err := os.ReadFile(queryFile)
	if err != nil {
		t.Fatalf("failed to read generated query file: %v", err)
	}

	contentStr := string(content)
	t.Logf("Generated query content:\n%s", contentStr)

	// Verify the struct has the correct fields
	if !strings.Contains(contentStr, "PostID string") {
		t.Error("expected PostID string field in generated query")
	}
	if !strings.Contains(contentStr, "Status string") {
		t.Error("expected Status string field in generated query")
	}
	if !strings.Contains(contentStr, "PublishedAfter time.Time") {
		t.Error("expected PublishedAfter time.Time field in generated query")
	}

	// Verify the struct is named correctly
	if !strings.Contains(contentStr, "type SearchQuery struct") {
		t.Error("expected SearchQuery struct name")
	}

	// Verify the handler method is named correctly
	if !strings.Contains(contentStr, "func (q *Querier) HandleSearch(") {
		t.Error("expected HandleSearch method name")
	}

	// Read the generated base file
	baseFile := filepath.Join(tempDir, "posts", "queries", "base.go")
	baseContent, err := os.ReadFile(baseFile)
	if err != nil {
		t.Fatalf("failed to read generated base file: %v", err)
	}

	baseContentStr := string(baseContent)
	t.Logf("Generated base content:\n%s", baseContentStr)

	// Verify ItemResult has the custom fields
	if !strings.Contains(baseContentStr, "PostID string") {
		t.Error("expected PostID string field in ItemResult")
	}
	if !strings.Contains(baseContentStr, "Status string") {
		t.Error("expected Status string field in ItemResult")
	}
	if !strings.Contains(baseContentStr, "PublishedAfter time.Time") {
		t.Error("expected PublishedAfter time.Time field in ItemResult")
	}

	// Verify that original fields are replaced
	if strings.Contains(baseContentStr, "Name        string") {
		t.Error("expected original Name field to be replaced in ItemResult")
	}
}
