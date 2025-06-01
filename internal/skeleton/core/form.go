package core

// This file previously contained the legacy Form struct and validation methods.
// All form functionality has been moved to the new tag-based validation system.
// 
// For form handling, use:
// - core.ParseFromURLValues() for parsing and validation
// - ui.FormData for template rendering
// - ui.FormGroupWithValidation(), ui.TextInputWithValidation(), etc. for UI components
//
// See docs/form-validation-guide.md for complete documentation.
