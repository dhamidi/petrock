# Plan for core/form.go

This file provides a flexible data structure for capturing and validating HTML form data, including error handling.

## Types

- `Form`: A struct holding form data and validation errors.
    - `Values url.Values`: The raw form values, typically parsed from an HTTP request.
    - `Errors map[string][]string`: A map where keys are field names and values are slices of error messages for that field.

## Functions

- `NewForm(data url.Values) *Form`: Creates a new `Form` instance initialized with the provided `url.Values`.
- `(f *Form) HasError(field string) bool`: Checks if a specific field has any associated validation errors.
- `(f *Form) GetError(field string) string`: Returns the *first* error message for a given field. Returns an empty string if there are no errors for that field.
- `(f *Form) AddError(field, message string)`: Adds a new error message to a specific field.
- `(f *Form) ValidateRequired(fields ...string)`: Checks if the specified fields are present and non-empty in `f.Values`. Adds errors using `AddError` if validation fails.
- `(f *Form) ValidateMinLength(field string, length int)`: Checks if the value of the specified field has at least `length` characters. Adds an error if validation fails.
- `(f *Form) ValidateEmail(field string)`: Checks if the value of the specified field looks like a valid email address. Adds an error if validation fails.
- `(f *Form) IsValid() bool`: Returns `true` if the `Errors` map is empty, `false` otherwise.
