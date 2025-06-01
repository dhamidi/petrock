package core

import (
	"testing"
)

func TestArgsSource_KeyValueFormat(t *testing.T) {
	args := []string{"--name=john", "--age=25"}
	source := NewArgsSource(args)

	if source.Get("name") != "john" {
		t.Errorf("Expected name 'john', got '%s'", source.Get("name"))
	}

	if source.Get("age") != "25" {
		t.Errorf("Expected age '25', got '%s'", source.Get("age"))
	}
}

func TestArgsSource_SpaceSeparatedFormat(t *testing.T) {
	args := []string{"--name", "john", "--age", "25"}
	source := NewArgsSource(args)

	if source.Get("name") != "john" {
		t.Errorf("Expected name 'john', got '%s'", source.Get("name"))
	}

	if source.Get("age") != "25" {
		t.Errorf("Expected age '25', got '%s'", source.Get("age"))
	}
}

func TestArgsSource_BooleanFlags(t *testing.T) {
	args := []string{"--active", "--no-debug"}
	source := NewArgsSource(args)

	if source.Get("active") != "true" {
		t.Errorf("Expected active 'true', got '%s'", source.Get("active"))
	}

	if source.Get("debug") != "false" {
		t.Errorf("Expected debug 'false', got '%s'", source.Get("debug"))
	}
}

func TestArgsSource_MultipleValues(t *testing.T) {
	args := []string{"--tags=go", "--tags=web", "--tags=test"}
	source := NewArgsSource(args)

	tags := source.GetAll("tags")
	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}

	expected := []string{"go", "web", "test"}
	for i, tag := range tags {
		if tag != expected[i] {
			t.Errorf("Expected tag '%s', got '%s'", expected[i], tag)
		}
	}
}

func TestArgsSource_MixedFormats(t *testing.T) {
	args := []string{"--name=john", "--age", "25", "--active", "--tags=go", "--tags=web"}
	source := NewArgsSource(args)

	if source.Get("name") != "john" {
		t.Errorf("Expected name 'john', got '%s'", source.Get("name"))
	}

	if source.Get("age") != "25" {
		t.Errorf("Expected age '25', got '%s'", source.Get("age"))
	}

	if source.Get("active") != "true" {
		t.Errorf("Expected active 'true', got '%s'", source.Get("active"))
	}

	tags := source.GetAll("tags")
	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}
}

func TestParseFromArgs(t *testing.T) {
	type CLIConfig struct {
		Name   string   `json:"name" validate:"required"`
		Age    int      `json:"age" validate:"min=0"`
		Active bool     `json:"active"`
		Tags   []string `json:"tags"`
	}

	args := []string{"--name=john", "--age=25", "--active", "--tags=go", "--tags=web"}

	var config CLIConfig
	err := ParseFromArgs(args, &config)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if config.Name != "john" {
		t.Errorf("Expected name 'john', got '%s'", config.Name)
	}

	if config.Age != 25 {
		t.Errorf("Expected age 25, got %d", config.Age)
	}

	if !config.Active {
		t.Errorf("Expected active true, got false")
	}

	if len(config.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(config.Tags))
	}
}

func TestParseFromArgs_ValidationError(t *testing.T) {
	type CLIConfig struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"min=0"`
	}

	args := []string{"--age=-5"} // Missing required name, invalid age

	var config CLIConfig
	err := ParseFromArgs(args, &config)

	if err == nil {
		t.Errorf("Expected validation error")
	}

	parseErrors, ok := err.(*ParseErrors)
	if !ok {
		t.Errorf("Expected ParseErrors, got %T", err)
	}

	if len(parseErrors.Errors) == 0 {
		t.Errorf("Expected validation errors, got none")
	}

	// Should have both required and min_value errors
	hasRequiredError := false
	hasMinValueError := false

	for _, e := range parseErrors.Errors {
		if e.Code == "required" && e.Field == "name" {
			hasRequiredError = true
		}
		if e.Code == "min_value" && e.Field == "age" {
			hasMinValueError = true
		}
	}

	if !hasRequiredError {
		t.Errorf("Expected required validation error for 'name' field")
	}

	if !hasMinValueError {
		t.Errorf("Expected min_value validation error for 'age' field")
	}
}
