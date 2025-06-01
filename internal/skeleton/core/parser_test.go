package core

import (
	"net/url"
	"testing"
	"time"
)

// Test structs for validation
type TestStruct struct {
	Name        string    `json:"name" validate:"required,minlen=2,maxlen=50"`
	Email       string    `json:"email" validate:"email"`
	Age         int       `json:"age" validate:"min=0,max=120"`
	Active      bool      `json:"active"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description" validate:"maxlen=500"`
}

func TestParseFromURLValues(t *testing.T) {
	values := url.Values{
		"name":        []string{"John Doe"},
		"email":       []string{"john@example.com"},
		"age":         []string{"25"},
		"active":      []string{"true"},
		"tags":        []string{"go", "web", "test"},
		"created_at":  []string{"2023-01-01T00:00:00Z"},
		"description": []string{"A test user"},
	}

	var result TestStruct
	err := ParseFromURLValues(values, &result)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", result.Name)
	}

	if result.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", result.Email)
	}

	if result.Age != 25 {
		t.Errorf("Expected age 25, got %d", result.Age)
	}

	if !result.Active {
		t.Errorf("Expected active true, got false")
	}

	if len(result.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(result.Tags))
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	if !result.CreatedAt.Equal(expectedTime) {
		t.Errorf("Expected created_at %v, got %v", expectedTime, result.CreatedAt)
	}
}

func TestValidationRequired(t *testing.T) {
	values := url.Values{
		"email": []string{"john@example.com"},
		"age":   []string{"25"},
	}

	var result TestStruct
	err := ParseFromURLValues(values, &result)

	if err == nil {
		t.Errorf("Expected validation error for missing required field")
	}

	parseErrors, ok := err.(*ParseErrors)
	if !ok {
		t.Errorf("Expected ParseErrors, got %T", err)
	}

	if len(parseErrors.Errors) == 0 {
		t.Errorf("Expected validation errors, got none")
	}

	// Check that we have error for required field
	foundRequiredError := false
	for _, e := range parseErrors.Errors {
		if e.Code == "required" && e.Field == "name" {
			foundRequiredError = true
			break
		}
	}

	if !foundRequiredError {
		t.Errorf("Expected required validation error for 'name' field")
	}
}

func TestValidationMinLength(t *testing.T) {
	values := url.Values{
		"name":  []string{"A"}, // Too short
		"email": []string{"john@example.com"},
		"age":   []string{"25"},
	}

	var result TestStruct
	err := ParseFromURLValues(values, &result)

	if err == nil {
		t.Errorf("Expected validation error for short name")
	}

	parseErrors, ok := err.(*ParseErrors)
	if !ok {
		t.Errorf("Expected ParseErrors, got %T", err)
	}

	// Check for min_length error
	foundMinLengthError := false
	for _, e := range parseErrors.Errors {
		if e.Code == "min_length" && e.Field == "name" {
			foundMinLengthError = true
			break
		}
	}

	if !foundMinLengthError {
		t.Errorf("Expected min_length validation error for 'name' field")
	}
}

func TestValidationEmail(t *testing.T) {
	values := url.Values{
		"name":  []string{"John Doe"},
		"email": []string{"invalid-email"},
		"age":   []string{"25"},
	}

	var result TestStruct
	err := ParseFromURLValues(values, &result)

	if err == nil {
		t.Errorf("Expected validation error for invalid email")
	}

	parseErrors, ok := err.(*ParseErrors)
	if !ok {
		t.Errorf("Expected ParseErrors, got %T", err)
	}

	// Check for email validation error
	foundEmailError := false
	for _, e := range parseErrors.Errors {
		if e.Code == "invalid_email" && e.Field == "email" {
			foundEmailError = true
			break
		}
	}

	if !foundEmailError {
		t.Errorf("Expected invalid_email validation error for 'email' field")
	}
}

func TestValidationRange(t *testing.T) {
	values := url.Values{
		"name":  []string{"John Doe"},
		"email": []string{"john@example.com"},
		"age":   []string{"150"}, // Too high
	}

	var result TestStruct
	err := ParseFromURLValues(values, &result)

	if err == nil {
		t.Errorf("Expected validation error for age out of range")
	}

	parseErrors, ok := err.(*ParseErrors)
	if !ok {
		t.Errorf("Expected ParseErrors, got %T", err)
	}

	// Check for max_value error
	foundMaxValueError := false
	for _, e := range parseErrors.Errors {
		if e.Code == "max_value" && e.Field == "age" {
			foundMaxValueError = true
			break
		}
	}

	if !foundMaxValueError {
		t.Errorf("Expected max_value validation error for 'age' field")
	}
}

func TestParseFromMap(t *testing.T) {
	data := map[string]interface{}{
		"name":        "Jane Doe",
		"email":       "jane@example.com",
		"age":         30,
		"active":      true,
		"tags":        []string{"api", "json"},
		"created_at":  "2023-01-01T00:00:00Z",
		"description": "A test user from JSON",
	}

	var result TestStruct
	err := ParseFromMap(data, &result)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result.Name != "Jane Doe" {
		t.Errorf("Expected name 'Jane Doe', got '%s'", result.Name)
	}

	if result.Age != 30 {
		t.Errorf("Expected age 30, got %d", result.Age)
	}
}
