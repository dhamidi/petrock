package core

import (
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

// RequiredValidator checks for required fields
type RequiredValidator struct{}

func (v RequiredValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.GetTagBool("required")
}

func (v RequiredValidator) Validate(ctx *FieldContext) []ParseError {
	switch val := ctx.Value.(type) {
	case string:
		if strings.TrimSpace(val) == "" {
			return []ParseError{{
				Field: ctx.Name, Message: "This field is required", Code: "required",
			}}
		}
	case nil:
		return []ParseError{{
			Field: ctx.Name, Message: "This field is required", Code: "required",
		}}
	default:
		if reflect.ValueOf(val).IsZero() {
			return []ParseError{{
				Field: ctx.Name, Message: "This field is required", Code: "required",
			}}
		}
	}
	return nil
}

// LengthValidator checks string length constraints
type LengthValidator struct{}

func (v LengthValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.FieldType.Kind() == reflect.String &&
		(ctx.GetTag("minlen", "") != "" || ctx.GetTag("maxlen", "") != "")
}

func (v LengthValidator) Validate(ctx *FieldContext) []ParseError {
	str, ok := ctx.Value.(string)
	if !ok || str == "" {
		return nil // Skip validation for non-strings or empty strings
	}

	var errors []ParseError
	length := utf8.RuneCountInString(str)

	if minLen := ctx.GetTagInt("minlen", 0); minLen > 0 && length < minLen {
		errors = append(errors, ParseError{
			Field:   ctx.Name,
			Message: fmt.Sprintf("Must be at least %d characters long", minLen),
			Code:    "min_length",
			Meta:    map[string]interface{}{"min_length": minLen, "actual_length": length},
		})
	}

	if maxLen := ctx.GetTagInt("maxlen", 0); maxLen > 0 && length > maxLen {
		errors = append(errors, ParseError{
			Field:   ctx.Name,
			Message: fmt.Sprintf("Must be no more than %d characters long", maxLen),
			Code:    "max_length",
			Meta:    map[string]interface{}{"max_length": maxLen, "actual_length": length},
		})
	}

	return errors
}

// RangeValidator checks numeric range constraints
type RangeValidator struct{}

func (v RangeValidator) CanValidate(ctx *FieldContext) bool {
	kind := ctx.FieldType.Kind()
	isNumeric := kind >= reflect.Int && kind <= reflect.Float64
	return isNumeric && (ctx.GetTag("min", "") != "" || ctx.GetTag("max", "") != "")
}

func (v RangeValidator) Validate(ctx *FieldContext) []ParseError {
	var numVal int64
	var isValid bool

	switch val := ctx.Value.(type) {
	case int:
		numVal, isValid = int64(val), true
	case int8:
		numVal, isValid = int64(val), true
	case int16:
		numVal, isValid = int64(val), true
	case int32:
		numVal, isValid = int64(val), true
	case int64:
		numVal, isValid = val, true
	case uint:
		numVal, isValid = int64(val), true
	case uint8:
		numVal, isValid = int64(val), true
	case uint16:
		numVal, isValid = int64(val), true
	case uint32:
		numVal, isValid = int64(val), true
	case uint64:
		numVal, isValid = int64(val), true
	}

	if !isValid {
		return nil
	}

	var errors []ParseError

	if minTag := ctx.GetTag("min", ""); minTag != "" {
		if min, err := strconv.ParseInt(minTag, 10, 64); err == nil && numVal < min {
			errors = append(errors, ParseError{
				Field:   ctx.Name,
				Message: fmt.Sprintf("Must be at least %d", min),
				Code:    "min_value",
				Meta:    map[string]interface{}{"min_value": min, "actual_value": numVal},
			})
		}
	}

	if maxTag := ctx.GetTag("max", ""); maxTag != "" {
		if max, err := strconv.ParseInt(maxTag, 10, 64); err == nil && numVal > max {
			errors = append(errors, ParseError{
				Field:   ctx.Name,
				Message: fmt.Sprintf("Must be no more than %d", max),
				Code:    "max_value",
				Meta:    map[string]interface{}{"max_value": max, "actual_value": numVal},
			})
		}
	}

	return errors
}

// EmailValidator validates email format
type EmailValidator struct{}

func (v EmailValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.FieldType.Kind() == reflect.String && ctx.GetTagBool("email")
}

func (v EmailValidator) Validate(ctx *FieldContext) []ParseError {
	str, ok := ctx.Value.(string)
	if !ok || str == "" {
		return nil
	}

	if _, err := mail.ParseAddress(str); err != nil {
		return []ParseError{{
			Field:   ctx.Name,
			Message: "Must be a valid email address",
			Code:    "invalid_email",
		}}
	}

	return nil
}

// CrossFieldValidator validates fields against other fields in the same struct
type CrossFieldValidator struct{}

func (v CrossFieldValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.GetTag("confirm_field", "") != ""
}

func (v CrossFieldValidator) Validate(ctx *FieldContext) []ParseError {
	confirmField := ctx.GetTag("confirm_field", "")
	if confirmField == "" {
		return nil
	}

	// This is a simplified implementation - in a full implementation,
	// you'd need access to the entire struct being validated
	// For now, this serves as a demonstration of the pattern
	return []ParseError{{
		Field:   ctx.Name,
		Message: fmt.Sprintf("Must match %s field", confirmField),
		Code:    "field_mismatch",
		Meta: map[string]interface{}{
			"confirm_field": confirmField,
		},
	}}
}

// ConditionalValidator validates fields only when conditions are met
type ConditionalValidator struct{}

func (v ConditionalValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.GetTag("required_if", "") != ""
}

func (v ConditionalValidator) Validate(ctx *FieldContext) []ParseError {
	requiredIf := ctx.GetTag("required_if", "")
	if requiredIf == "" {
		return nil
	}

	// Parse the condition: "field:value"
	parts := strings.SplitN(requiredIf, ":", 2)
	if len(parts) != 2 {
		return []ParseError{{
			Field:   ctx.Name,
			Message: "Invalid required_if condition format",
			Code:    "invalid_condition",
		}}
	}

	// This is a demonstration - full implementation would check the actual field values
	return nil // Conditional logic would go here
}

// CustomMessageValidator allows custom error messages
type CustomMessageValidator struct{}

func (v CustomMessageValidator) CanValidate(ctx *FieldContext) bool {
	return ctx.GetTag("message", "") != ""
}

func (v CustomMessageValidator) Validate(ctx *FieldContext) []ParseError {
	// This validator doesn't actually validate - it's used by other validators
	// to override their default messages. This is a pattern demonstration.
	return nil
}

// StandardTagParser parses common validation tags
type StandardTagParser struct{}

func (p StandardTagParser) ParseTags(field reflect.StructField) map[string]string {
	tags := make(map[string]string)

	// Parse validate tag: "required,minlen=2,maxlen=100,email"
	if validateTag := field.Tag.Get("validate"); validateTag != "" {
		parts := strings.Split(validateTag, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				tags[kv[0]] = kv[1]
			} else {
				tags[part] = "true"
			}
		}
	}

	return tags
}
