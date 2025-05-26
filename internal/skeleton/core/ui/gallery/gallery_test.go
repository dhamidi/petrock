package gallery

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/petrock/example_module_path/core"
)

func TestHandleGallery(t *testing.T) {
	// Create a test app (can be nil for this test since handler doesn't use app)
	var app *core.App

	// Create the handler
	handler := HandleGallery(app)

	// Create a test request
	req, err := http.NewRequest("GET", "/_/ui", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	expected := "text/html; charset=utf-8"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("handler returned wrong content type: got %v want %v", ct, expected)
	}

	// Check that the response contains expected HTML structure
	body := rr.Body.String()

	// Check for essential HTML structure
	if !strings.Contains(body, "<html") {
		t.Error("response should contain HTML tag")
	}

	if !strings.Contains(body, "<title>UI Component Gallery</title>") {
		t.Error("response should contain page title")
	}

	if !strings.Contains(body, "gallery-container") {
		t.Error("response should contain gallery container")
	}

	if !strings.Contains(body, "sidebar") {
		t.Error("response should contain sidebar")
	}

	if !strings.Contains(body, "content") {
		t.Error("response should contain content area")
	}

	// Check for welcome content since no components exist initially
	if !strings.Contains(body, "Welcome to the UI component gallery") {
		t.Error("response should contain welcome message")
	}

	if !strings.Contains(body, "No components available yet") {
		t.Error("response should show empty state message when no components exist")
	}
}

func TestHandleComponentDetail(t *testing.T) {
	// Create a test app (can be nil for this test since handler doesn't use app)
	var app *core.App

	// Create the handler
	handler := HandleComponentDetail(app)

	t.Run("ValidComponentName", func(t *testing.T) {
		// Test with a component name (even though no components exist initially)
		req, err := http.NewRequest("GET", "/_/ui/button", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler(rr, req)

		// Should return 200 even if component doesn't exist (shows not found page)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		body := rr.Body.String()

		// Should contain "Component Not Found" since no components exist initially
		if !strings.Contains(body, "Component Not Found") {
			t.Error("response should contain 'Component Not Found' for non-existent component")
		}

		if !strings.Contains(body, "button") {
			t.Error("response should mention the requested component name")
		}
	})

	t.Run("EmptyComponentName", func(t *testing.T) {
		// Test with empty component name
		req, err := http.NewRequest("GET", "/_/ui/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler(rr, req)

		// Should return 400 for empty component name
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("HTMLStructure", func(t *testing.T) {
		// Test HTML structure for valid request
		req, err := http.NewRequest("GET", "/_/ui/test-component", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler(rr, req)

		body := rr.Body.String()

		// Check for essential HTML structure
		if !strings.Contains(body, "<html") {
			t.Error("response should contain HTML tag")
		}

		if !strings.Contains(body, "component-container") {
			t.Error("response should contain component container")
		}

		if !strings.Contains(body, "sidebar") {
			t.Error("response should contain sidebar")
		}

		if !strings.Contains(body, "← Gallery") {
			t.Error("response should contain back link to gallery")
		}
	})
}

func TestGetAllComponents(t *testing.T) {
	components := GetAllComponents()

	// Initially should return empty slice
	if len(components) != 0 {
		t.Errorf("GetAllComponents should initially return empty slice, got %d components", len(components))
	}

	// Verify the return type is correct
	if components == nil {
		t.Error("GetAllComponents should return empty slice, not nil")
	}
}

func TestComponentInfo(t *testing.T) {
	// Test ComponentInfo struct
	component := ComponentInfo{
		Name:        "Button",
		Description: "A clickable button component",
		Category:    "Interactive",
	}

	if component.Name != "Button" {
		t.Error("ComponentInfo.Name should be set correctly")
	}

	if component.Description != "A clickable button component" {
		t.Error("ComponentInfo.Description should be set correctly")
	}

	if component.Category != "Interactive" {
		t.Error("ComponentInfo.Category should be set correctly")
	}
}

// Benchmark tests for performance
func BenchmarkHandleGallery(b *testing.B) {
	var app *core.App
	handler := HandleGallery(app)
	req, _ := http.NewRequest("GET", "/_/ui", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler(rr, req)
	}
}

func BenchmarkHandleComponentDetail(b *testing.B) {
	var app *core.App
	handler := HandleComponentDetail(app)
	req, _ := http.NewRequest("GET", "/_/ui/button", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler(rr, req)
	}
}