package core

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"
	"unicode/utf8"
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

// --- Validation Methods ---

// ValidateRequired checks if the specified fields are present and non-empty in f.Values.
// Adds errors using AddError if validation fails.
func (f *Form) ValidateRequired(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.AddError(field, "This field cannot be blank")
		}
	}
}

// ValidateMinLength checks if the value of the specified field has at least `length` characters.
// Adds an error if validation fails. Does not add error if field is empty (use ValidateRequired first).
func (f *Form) ValidateMinLength(field string, length int) {
	value := f.Get(field)
	if value == "" {
		return // Don't check length if field is empty
	}
	if utf8.RuneCountInString(value) < length {
		f.AddError(field, fmt.Sprintf("This field must be at least %d characters long", length))
	}
}

// ValidateMaxLength checks if the value of the specified field has at most `length` characters.
// Adds an error if validation fails.
func (f *Form) ValidateMaxLength(field string, length int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > length {
		f.AddError(field, fmt.Sprintf("This field must be no more than %d characters long", length))
	}
}

// ValidateEmail checks if the value of the specified field looks like a valid email address.
// Adds an error if validation fails. Does not add error if field is empty.
func (f *Form) ValidateEmail(field string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		f.AddError(field, "This field must be a valid email address")
	}
}

// ValidateAllowedValues checks if the value of the specified field is one of the allowed values.
// Adds an error if validation fails.
func (f *Form) ValidateAllowedValues(field string, allowed ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, v := range allowed {
		if value == v {
			return // Found a valid value
		}
	}
	f.AddError(field, "This field has an invalid value")
}
