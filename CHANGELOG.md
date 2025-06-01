# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **New Unified Form Validation System**: Tag-based validation system supporting multiple input sources
  - `ParseFromURLValues()` for HTTP form parsing with validation
  - `ParseFromMap()` for JSON API parsing with validation  
  - `ParseFromArgs()` for CLI argument parsing with validation
  - Built-in validators: `required`, `minlen`, `maxlen`, `min`, `max`, `email`
  - Advanced validators: `confirm_field`, `required_if`, custom messages
  - Structured error handling with `ParseError` and `ParseErrors`
  - Extensible converter and validator registries
- **CLI Argument Parsing**: Support for `--key=value`, `--key value`, and boolean flag formats
- **Enhanced API Error Responses**: Structured JSON error responses for validation failures
- **Comprehensive Documentation**: 
  - Form validation guide with examples and best practices
  - Migration guide from old form system
  - Updated README with validation system overview

### Changed

- **BREAKING**: Command and query structs now use validation tags instead of manual validation
- **BREAKING**: HTTP handlers use new parsing system instead of `core.NewForm()` validation methods
- **BREAKING**: API endpoints return structured validation errors instead of simple error messages
- Form validation methods (`ValidateRequired`, `ValidateEmail`, etc.) removed from `core.Form`
- Generated feature templates include validation tags by default

### Deprecated

- Manual form validation methods (removed from `core.Form`)
- Direct usage of `form.ValidateRequired()` and similar methods (use validation tags instead)

### Migration Guide

#### For HTTP Form Handlers

**Before:**
```go
form := core.NewForm(r.PostForm)
form.ValidateRequired("name", "email")
if !form.IsValid() {
    // handle errors
}
cmd := CreateCommand{
    Name:  form.Get("name"),
    Email: form.Get("email"),
}
```

**After:**
```go
var cmd CreateCommand
if err := core.ParseFromURLValues(r.PostForm, &cmd); err != nil {
    if parseErrors, ok := err.(*core.ParseErrors); ok {
        // Convert to form errors for template rendering
        form := core.NewForm(r.PostForm)
        for _, parseErr := range parseErrors.Errors {
            form.AddError(parseErr.Field, parseErr.Message)
        }
        // render with errors
    }
    return
}
```

#### For Command/Query Structs

**Before:**
```go
type CreateCommand struct {
    Name        string    `json:"name"`
    Description string    `json:"description"`
}

func (c *CreateCommand) Validate(state *State) error {
    if strings.TrimSpace(c.Name) == "" {
        return errors.New("name cannot be empty")
    }
    // ... more validation
}
```

**After:**
```go
type CreateCommand struct {
    Name        string    `json:"name" validate:"required,minlen=2,maxlen=100"`
    Description string    `json:"description" validate:"required,minlen=5,maxlen=500"`
}

func (c *CreateCommand) Validate(state *State) error {
    // Only business logic validation remains
    // Basic field validation handled by tags
    return nil
}
```

#### For API Error Handling

**Before:**
```json
{
  "error": "validation failed: name is required"
}
```

**After:**
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "name",
      "message": "This field is required",
      "code": "required"
    }
  ]
}
```

### Notes

- Template rendering compatibility maintained through `core.Form` error handling
- Existing projects can migrate incrementally
- Old validation methods still available in git history if needed for reference
