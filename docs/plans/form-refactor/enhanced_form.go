package core

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// FormSource represents different input sources for form data
type FormSource interface {
	// Get returns the first value for the given key
	Get(key string) string
	// GetAll returns all values for the given key
	GetAll(key string) []string
	// Keys returns all available keys
	Keys() []string
}

// URLValuesSource wraps url.Values to implement FormSource
type URLValuesSource struct {
	Values url.Values
}

func (u URLValuesSource) Get(key string) string {
	return u.Values.Get(key)
}

func (u URLValuesSource) GetAll(key string) []string {
	return u.Values[key]
}

func (u URLValuesSource) Keys() []string {
	keys := make([]string, 0, len(u.Values))
	for k := range u.Values {
		keys = append(keys, k)
	}
	return keys
}

// MapSource wraps a map[string]interface{} to implement FormSource
type MapSource struct {
	Data map[string]interface{}
}

func (m MapSource) Get(key string) string {
	if val, exists := m.Data[key]; exists {
		switch v := val.(type) {
		case string:
			return v
		case []string:
			if len(v) > 0 {
				return v[0]
			}
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

func (m MapSource) GetAll(key string) []string {
	if val, exists := m.Data[key]; exists {
		switch v := val.(type) {
		case string:
			return []string{v}
		case []string:
			return v
		case []interface{}:
			result := make([]string, len(v))
			for i, item := range v {
				result[i] = fmt.Sprintf("%v", item)
			}
			return result
		default:
			return []string{fmt.Sprintf("%v", v)}
		}
	}
	return nil
}

func (m MapSource) Keys() []string {
	keys := make([]string, 0, len(m.Data))
	for k := range m.Data {
		keys = append(keys, k)
	}
	return keys
}

// ArgsSource wraps command line arguments to implement FormSource
type ArgsSource struct {
	Args map[string][]string
}

func NewArgsSource(args []string) *ArgsSource {
	result := &ArgsSource{
		Args: make(map[string][]string),
	}
	
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(arg, "--")
			if strings.Contains(key, "=") {
				// Handle --key=value format
				parts := strings.SplitN(key, "=", 2)
				result.Args[parts[0]] = append(result.Args[parts[0]], parts[1])
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				// Handle --key value format
				result.Args[key] = append(result.Args[key], args[i+1])
				i++ // Skip the value
			} else {
				// Handle --key (boolean flag)
				result.Args[key] = append(result.Args[key], "true")
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			key := strings.TrimPrefix(arg, "-")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				result.Args[key] = append(result.Args[key], args[i+1])
				i++ // Skip the value
			} else {
				result.Args[key] = append(result.Args[key], "true")
			}
		}
	}
	
	return result
}

func (a ArgsSource) Get(key string) string {
	if vals := a.Args[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (a ArgsSource) GetAll(key string) []string {
	return a.Args[key]
}

func (a ArgsSource) Keys() []string {
	keys := make([]string, 0, len(a.Args))
	for k := range a.Args {
		keys = append(keys, k)
	}
	return keys
}

// FieldError represents an error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"` // Optional error code for programmatic handling
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// FormErrors holds validation errors
type FormErrors struct {
	FieldErrors []FieldError `json:"field_errors"`
	FormErrors  []string     `json:"form_errors"` // Form-wide errors
}

func (e *FormErrors) AddFieldError(field, message string) {
	e.FieldErrors = append(e.FieldErrors, FieldError{
		Field:   field,
		Message: message,
	})
}

func (e *FormErrors) AddFieldErrorWithCode(field, message, code string) {
	e.FieldErrors = append(e.FieldErrors, FieldError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

func (e *FormErrors) AddFormError(message string) {
	e.FormErrors = append(e.FormErrors, message)
}

func (e *FormErrors) HasErrors() bool {
	return len(e.FieldErrors) > 0 || len(e.FormErrors) > 0
}

func (e *FormErrors) GetFieldErrors(field string) []FieldError {
	var errors []FieldError
	for _, err := range e.FieldErrors {
		if err.Field == field {
			errors = append(errors, err)
		}
	}
	return errors
}

func (e *FormErrors) GetFirstFieldError(field string) string {
	for _, err := range e.FieldErrors {
		if err.Field == field {
			return err.Message
		}
	}
	return ""
}

// ValidationRule represents a validation rule that can be applied to fields
type ValidationRule interface {
	Validate(value interface{}, fieldName string) *FieldError
}

// RequiredRule validates that a field is not empty
type RequiredRule struct {
	Message string
}

func (r RequiredRule) Validate(value interface{}, fieldName string) *FieldError {
	message := r.Message
	if message == "" {
		message = "This field is required"
	}
	
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return &FieldError{Field: fieldName, Message: message, Code: "required"}
		}
	case nil:
		return &FieldError{Field: fieldName, Message: message, Code: "required"}
	default:
		// For other types, check if it's the zero value
		if reflect.ValueOf(v).IsZero() {
			return &FieldError{Field: fieldName, Message: message, Code: "required"}
		}
	}
	return nil
}

// MinLengthRule validates minimum string length
type MinLengthRule struct {
	MinLength int
	Message   string
}

func (r MinLengthRule) Validate(value interface{}, fieldName string) *FieldError {
	str, ok := value.(string)
	if !ok {
		return nil // Skip validation for non-strings
	}
	
	if str == "" {
		return nil // Don't validate empty strings (use RequiredRule for that)
	}
	
	if utf8.RuneCountInString(str) < r.MinLength {
		message := r.Message
		if message == "" {
			message = fmt.Sprintf("Must be at least %d characters long", r.MinLength)
		}
		return &FieldError{Field: fieldName, Message: message, Code: "min_length"}
	}
	return nil
}

// MaxLengthRule validates maximum string length
type MaxLengthRule struct {
	MaxLength int
	Message   string
}

func (r MaxLengthRule) Validate(value interface{}, fieldName string) *FieldError {
	str, ok := value.(string)
	if !ok {
		return nil
	}
	
	if utf8.RuneCountInString(str) > r.MaxLength {
		message := r.Message
		if message == "" {
			message = fmt.Sprintf("Must be no more than %d characters long", r.MaxLength)
		}
		return &FieldError{Field: fieldName, Message: message, Code: "max_length"}
	}
	return nil
}

// EmailRule validates email format
type EmailRule struct {
	Message string
}

func (r EmailRule) Validate(value interface{}, fieldName string) *FieldError {
	str, ok := value.(string)
	if !ok || str == "" {
		return nil
	}
	
	_, err := mail.ParseAddress(str)
	if err != nil {
		message := r.Message
		if message == "" {
			message = "Must be a valid email address"
		}
		return &FieldError{Field: fieldName, Message: message, Code: "invalid_email"}
	}
	return nil
}

// CustomRule allows for custom validation logic
type CustomRule struct {
	ValidateFunc func(value interface{}, fieldName string) *FieldError
}

func (r CustomRule) Validate(value interface{}, fieldName string) *FieldError {
	return r.ValidateFunc(value, fieldName)
}

// FieldConverter handles type conversion from string to target type
type FieldConverter interface {
	Convert(value string, targetType reflect.Type) (interface{}, error)
}

// DefaultConverter provides basic type conversion
type DefaultConverter struct{}

func (c DefaultConverter) Convert(value string, targetType reflect.Type) (interface{}, error) {
	switch targetType.Kind() {
	case reflect.String:
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			return reflect.Zero(targetType).Interface(), nil
		}
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer: %s", value)
		}
		// Convert to the specific int type
		switch targetType.Kind() {
		case reflect.Int:
			return int(intVal), nil
		case reflect.Int8:
			return int8(intVal), nil
		case reflect.Int16:
			return int16(intVal), nil
		case reflect.Int32:
			return int32(intVal), nil
		case reflect.Int64:
			return intVal, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			return reflect.Zero(targetType).Interface(), nil
		}
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid unsigned integer: %s", value)
		}
		switch targetType.Kind() {
		case reflect.Uint:
			return uint(uintVal), nil
		case reflect.Uint8:
			return uint8(uintVal), nil
		case reflect.Uint16:
			return uint16(uintVal), nil
		case reflect.Uint32:
			return uint32(uintVal), nil
		case reflect.Uint64:
			return uintVal, nil
		}
	case reflect.Float32, reflect.Float64:
		if value == "" {
			return reflect.Zero(targetType).Interface(), nil
		}
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float: %s", value)
		}
		if targetType.Kind() == reflect.Float32 {
			return float32(floatVal), nil
		}
		return floatVal, nil
	case reflect.Bool:
		if value == "" {
			return false, nil
		}
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean: %s", value)
		}
		return boolVal, nil
	case reflect.Slice:
		// For slices, we expect the source to provide multiple values
		return nil, fmt.Errorf("slice conversion requires multiple values")
	default:
		// Handle time.Time specifically
		if targetType == reflect.TypeOf(time.Time{}) {
			if value == "" {
				return time.Time{}, nil
			}
			// Try common time formats
			formats := []string{
				time.RFC3339,
				"2006-01-02",
				"2006-01-02 15:04:05",
				"01/02/2006",
				"01/02/2006 15:04:05",
			}
			for _, format := range formats {
				if t, err := time.Parse(format, value); err == nil {
					return t, nil
				}
			}
			return nil, fmt.Errorf("invalid time format: %s", value)
		}
		return nil, fmt.Errorf("unsupported type: %s", targetType.Kind())
	}
	return nil, fmt.Errorf("conversion failed")
}

// Form represents a form with validation and conversion capabilities
type Form struct {
	source    FormSource
	converter FieldConverter
	errors    FormErrors
	rules     map[string][]ValidationRule
}

// NewForm creates a new form with the given source
func NewForm(source FormSource) *Form {
	return &Form{
		source:    source,
		converter: DefaultConverter{},
		errors:    FormErrors{},
		rules:     make(map[string][]ValidationRule),
	}
}

// NewFormFromURLValues creates a form from url.Values
func NewFormFromURLValues(values url.Values) *Form {
	return NewForm(URLValuesSource{Values: values})
}

// NewFormFromMap creates a form from a map
func NewFormFromMap(data map[string]interface{}) *Form {
	return NewForm(MapSource{Data: data})
}

// NewFormFromArgs creates a form from command line arguments
func NewFormFromArgs(args []string) *Form {
	return NewForm(NewArgsSource(args))
}

// SetConverter sets a custom field converter
func (f *Form) SetConverter(converter FieldConverter) {
	f.converter = converter
}

// AddRule adds a validation rule for a field
func (f *Form) AddRule(fieldName string, rule ValidationRule) {
	f.rules[fieldName] = append(f.rules[fieldName], rule)
}

// AddRules adds multiple validation rules for a field
func (f *Form) AddRules(fieldName string, rules ...ValidationRule) {
	f.rules[fieldName] = append(f.rules[fieldName], rules...)
}

// Get returns the first value for a field
func (f *Form) Get(key string) string {
	return f.source.Get(key)
}

// GetAll returns all values for a field
func (f *Form) GetAll(key string) []string {
	return f.source.GetAll(key)
}

// Errors returns the form errors
func (f *Form) Errors() *FormErrors {
	return &f.errors
}

// HasErrors returns true if the form has any errors
func (f *Form) HasErrors() bool {
	return f.errors.HasErrors()
}

// Validate runs all validation rules and returns whether the form is valid
func (f *Form) Validate() bool {
	f.errors = FormErrors{} // Reset errors
	
	// Run field-specific validations
	for fieldName, rules := range f.rules {
		value := f.source.Get(fieldName)
		for _, rule := range rules {
			if err := rule.Validate(value, fieldName); err != nil {
				f.errors.FieldErrors = append(f.errors.FieldErrors, *err)
			}
		}
	}
	
	return !f.errors.HasErrors()
}

// BindTo populates a struct with form data, performing type conversion and validation
func (f *Form) BindTo(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	
	targetValue = targetValue.Elem()
	targetType := targetValue.Type()
	
	f.errors = FormErrors{} // Reset errors
	
	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Field(i)
		fieldType := targetType.Field(i)
		
		// Skip unexported fields
		if !field.CanSet() {
			continue
		}
		
		// Get field name from JSON tag or struct field name
		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "-" && parts[0] != "" {
				fieldName = parts[0]
			}
		}
		
		// Handle different field types
		if err := f.bindField(field, fieldType.Type, fieldName); err != nil {
			f.errors.AddFieldError(fieldName, err.Error())
		}
	}
	
	// Run validation rules
	for fieldName, rules := range f.rules {
		// Get the converted value from the struct
		fieldValue := f.getFieldValue(targetValue, targetType, fieldName)
		for _, rule := range rules {
			if err := rule.Validate(fieldValue, fieldName); err != nil {
				f.errors.FieldErrors = append(f.errors.FieldErrors, *err)
			}
		}
	}
	
	if f.errors.HasErrors() {
		return &f.errors
	}
	
	return nil
}

// bindField binds a single field
func (f *Form) bindField(field reflect.Value, fieldType reflect.Type, fieldName string) error {
	if fieldType.Kind() == reflect.Slice {
		// Handle slice fields
		values := f.source.GetAll(fieldName)
		if len(values) == 0 {
			return nil // Leave as zero value
		}
		
		elemType := fieldType.Elem()
		slice := reflect.MakeSlice(fieldType, len(values), len(values))
		
		for i, value := range values {
			convertedValue, err := f.converter.Convert(value, elemType)
			if err != nil {
				return fmt.Errorf("error converting slice element %d: %w", i, err)
			}
			slice.Index(i).Set(reflect.ValueOf(convertedValue))
		}
		
		field.Set(slice)
		return nil
	}
	
	// Handle single value fields
	value := f.source.Get(fieldName)
	if value == "" && fieldType.Kind() != reflect.String {
		return nil // Leave as zero value for non-strings
	}
	
	convertedValue, err := f.converter.Convert(value, fieldType)
	if err != nil {
		return err
	}
	
	field.Set(reflect.ValueOf(convertedValue))
	return nil
}

// getFieldValue gets a field value from a struct for validation
func (f *Form) getFieldValue(structValue reflect.Value, structType reflect.Type, fieldName string) interface{} {
	for i := 0; i < structValue.NumField(); i++ {
		fieldType := structType.Field(i)
		
		// Check if this is the field we're looking for
		checkName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "-" && parts[0] != "" {
				checkName = parts[0]
			}
		}
		
		if checkName == fieldName {
			return structValue.Field(i).Interface()
		}
	}
	return nil
}

// ToCommand creates and populates a command struct from form data
func (f *Form) ToCommand(commandType reflect.Type) (Command, error) {
	if commandType.Kind() != reflect.Ptr || commandType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("command type must be a pointer to a struct")
	}
	
	// Create new instance
	cmd := reflect.New(commandType.Elem()).Interface()
	
	// Bind form data to the command
	if err := f.BindTo(cmd); err != nil {
		return nil, err
	}
	
	// Type assert to Command interface
	command, ok := cmd.(Command)
	if !ok {
		return nil, fmt.Errorf("target does not implement Command interface")
	}
	
	return command, nil
}

// ToQuery creates and populates a query struct from form data
func (f *Form) ToQuery(queryType reflect.Type) (Query, error) {
	if queryType.Kind() != reflect.Ptr || queryType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("query type must be a pointer to a struct")
	}
	
	// Create new instance
	q := reflect.New(queryType.Elem()).Interface()
	
	// Bind form data to the query
	if err := f.BindTo(q); err != nil {
		return nil, err
	}
	
	// Type assert to Query interface
	query, ok := q.(Query)
	if !ok {
		return nil, fmt.Errorf("target does not implement Query interface")
	}
	
	return query, nil
}

// Helper functions for common validation rules
func Required(message ...string) ValidationRule {
	msg := "This field is required"
	if len(message) > 0 {
		msg = message[0]
	}
	return RequiredRule{Message: msg}
}

func MinLength(length int, message ...string) ValidationRule {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	return MinLengthRule{MinLength: length, Message: msg}
}

func MaxLength(length int, message ...string) ValidationRule {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	return MaxLengthRule{MaxLength: length, Message: msg}
}

func Email(message ...string) ValidationRule {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	return EmailRule{Message: msg}
}

func Custom(validateFunc func(value interface{}, fieldName string) *FieldError) ValidationRule {
	return CustomRule{ValidateFunc: validateFunc}
}
