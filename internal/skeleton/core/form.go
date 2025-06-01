package core

import (
	"net/url"
)

// Form holds form data and validation errors.
type Form struct {
	Values url.Values
	Errors map[string][]string
}

// NewForm creates a new Form instance initialized with the provided url.Values.
func NewForm(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(map[string][]string),
	}
}

// HasError checks if a specific field has any associated validation errors.
func (f *Form) HasError(field string) bool {
	_, exists := f.Errors[field]
	return exists
}

// GetError returns the first error message for a given field.
// Returns an empty string if there are no errors for that field.
func (f *Form) GetError(field string) string {
	if errs, exists := f.Errors[field]; exists && len(errs) > 0 {
		return errs[0]
	}
	return ""
}

// AddError adds a new error message to a specific field.
// If the field already has errors, the message is appended.
func (f *Form) AddError(field, message string) {
	f.Errors[field] = append(f.Errors[field], message)
}

// IsValid returns true if the Errors map is empty, false otherwise.
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

// Get returns the first value for the given form field key.
// Returns an empty string if the key is not present.
func (f *Form) Get(key string) string {
	return f.Values.Get(key)
}

// Note: Validation methods have been removed in favor of the new tag-based validation system.
// See docs/form-validation-guide.md for migration guidance.
// This form is now used only for template rendering and error display.
