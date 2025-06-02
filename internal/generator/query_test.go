package generator

import (
	"strings"
	"testing"
)

func TestModifyQueryStructWithFields(t *testing.T) {
	// Sample query template content (simplified)
	templateContent := `package queries

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/petrock/example_module_path/core"
)

type GetQuery struct {
	ID string ` + "`json:\"id\" validate:\"required\"`" + ` // ID of the entity to retrieve
}

func (q GetQuery) QueryName() string {
	return "petrock_example_feature_name/get"
}

func (q *Querier) HandleGet(ctx context.Context, query core.Query) (core.QueryResult, error) {
	// Lots of handler logic here...
	return result, nil
}
`

	fields := []QueryField{
		{Name: "postID", Type: "string"},
		{Name: "status", Type: "string"},
	}

	options := ExtractionOptions{
		FeatureName: "posts",
		EntityName:  "search",
		QueryFields: fields,
	}

	qg := &QueryGenerator{}
	result, err := qg.modifyQueryStructWithFields(templateContent, options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the struct has the right fields
	if !strings.Contains(result, "PostID string") {
		t.Error("expected PostID string field in result")
	}
	if !strings.Contains(result, "Status string") {
		t.Error("expected Status string field in result")
	}

	// Check that original fields are gone
	if strings.Contains(result, "ID string") {
		t.Error("expected original ID field to be replaced")
	}

	// Check that handler method is simplified
	handlerLines := strings.Split(result, "\n")
	handlerFound := false
	handlerSimple := false

	for i, line := range handlerLines {
		if strings.Contains(line, "func (q *Querier) Handle") {
			handlerFound = true
			// Check if the method body is simple
			for j := i + 1; j < len(handlerLines) && j < i+5; j++ {
				if strings.TrimSpace(handlerLines[j]) == "return nil, nil" {
					handlerSimple = true
					break
				}
			}
			break
		}
	}

	if !handlerFound {
		t.Error("expected handler method to be present")
	}
	if !handlerSimple {
		t.Error("expected handler method to be simplified to just return nil, nil")
	}

	t.Logf("Modified content:\n%s", result)
}

func TestModifyItemResultWithFields(t *testing.T) {
	// Sample base.go template content (simplified)
	templateContent := `package queries

import (
	"time"

	"github.com/petrock/example_module_path/petrock_example_feature_name/state"
)

type ItemResult struct {
	ID          string    ` + "`json:\"id\"`" + `
	Name        string    ` + "`json:\"name\"`" + `
	Description string    ` + "`json:\"description\"`" + `
	CreatedAt   time.Time ` + "`json:\"created_at\"`" + `
}
`

	fields := []QueryField{
		{Name: "postID", Type: "string"},
		{Name: "status", Type: "string"},
		{Name: "publishedAt", Type: "time.Time"},
	}

	options := ExtractionOptions{
		FeatureName: "posts",
		EntityName:  "search",
		QueryFields: fields,
	}

	qg := &QueryGenerator{}
	result, err := qg.modifyItemResultWithFields(templateContent, options)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the struct has the right fields
	if !strings.Contains(result, "PostID string") {
		t.Error("expected PostID string field in result")
	}
	if !strings.Contains(result, "Status string") {
		t.Error("expected Status string field in result")
	}
	if !strings.Contains(result, "PublishedAt time.Time") {
		t.Error("expected PublishedAt time.Time field in result")
	}

	// Check that original fields are gone
	if strings.Contains(result, "Name        string") {
		t.Error("expected original Name field to be replaced")
	}

	t.Logf("Modified ItemResult content:\n%s", result)
}
