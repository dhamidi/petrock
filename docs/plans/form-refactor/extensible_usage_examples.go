package examples

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"github.com/petrock/example_module_path/core"
)

// Custom UserRole enum type
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
	RoleGuest UserRole = "guest"
)

// UserRoleConverter converts strings to UserRole enum
type UserRoleConverter struct{}

func (c UserRoleConverter) CanConvert(targetType reflect.Type) bool {
	return targetType == reflect.TypeOf(UserRole(""))
}

func (c UserRoleConverter) Convert(value string, targetType reflect.Type) (interface{}, error) {
	role := UserRole(strings.ToLower(value))
	switch role {
	case RoleAdmin, RoleUser, RoleGuest:
		return role, nil
	default:
		return nil, fmt.Errorf("invalid user role: %s (valid: admin, user, guest)", value)
	}
}

func (c UserRoleConverter) ConvertSlice(values []string, targetType reflect.Type) (interface{}, error) {
	slice := make([]UserRole, len(values))
	for i, value := range values {
		converted, err := c.Convert(value, targetType.Elem())
		if err != nil {
			return nil, fmt.Errorf("error converting role element %d: %w", i, err)
		}
		slice[i] = converted.(UserRole)
	}
	return slice, nil
}

// RegexValidator validates fields against custom patterns
type RegexValidator struct{}

func (v RegexValidator) CanValidate(ctx *core.FieldContext) bool {
	return ctx.FieldType.Kind() == reflect.String && ctx.GetTag("pattern", "") != ""
}

func (v RegexValidator) Validate(ctx *core.FieldContext) []core.ParseError {
	str, ok := ctx.Value.(string)
	if !ok || str == "" {
		return nil
	}

	pattern := ctx.GetTag("pattern", "")
	if pattern == "" {
		return nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return []core.ParseError{{
			Field:   ctx.Name,
			Message: fmt.Sprintf("Invalid pattern configuration: %s", err.Error()),
			Code:    "invalid_pattern",
		}}
	}

	if !regex.MatchString(str) {
		return []core.ParseError{{
			Field:   ctx.Name,
			Message: fmt.Sprintf("Must match pattern: %s", pattern),
			Code:    "pattern_mismatch",
			Meta: map[string]interface{}{
				"pattern":      pattern,
				"actual_value": str,
			},
		}}
	}

	return nil
}

// UniqueUsernameValidator demonstrates business rule validation
type UniqueUsernameValidator struct {
	ExistingUsers []string // Mock database of existing usernames
}

func (v UniqueUsernameValidator) CanValidate(ctx *core.FieldContext) bool {
	return ctx.FieldType.Kind() == reflect.String && 
		   ctx.GetTagBool("unique_username")
}

func (v UniqueUsernameValidator) Validate(ctx *core.FieldContext) []core.ParseError {
	str, ok := ctx.Value.(string)
	if !ok || str == "" {
		return nil
	}

	// Mock database check
	for _, existing := range v.ExistingUsers {
		if strings.EqualFold(existing, str) {
			return []core.ParseError{{
				Field:   ctx.Name,
				Message: "Username is already taken",
				Code:    "username_taken",
				Meta: map[string]interface{}{
					"username": str,
				},
			}}
		}
	}

	return nil
}

// ExtendedTagParser supports custom tag syntax
type ExtendedTagParser struct{}

func (p ExtendedTagParser) ParseTags(field reflect.StructField) map[string]string {
	tags := make(map[string]string)

	// Parse custom business rule tags: "businessrule:unique_username"
	if businessTag := field.Tag.Get("businessrule"); businessTag != "" {
		tags[businessTag] = "true"
	}

	// Parse pattern tags: "pattern:^[A-Z]{2,3}$"
	if patternTag := field.Tag.Get("pattern"); patternTag != "" {
		tags["pattern"] = patternTag
	}

	// Parse role tags: "role:admin,user"
	if roleTag := field.Tag.Get("role"); roleTag != "" {
		tags["allowed_roles"] = roleTag
	}

	return tags
}

// Example usage demonstrating extensibility
func ExampleCustomComponents() {
	// Create a parser with custom components
	parser := core.NewParser()

	// Register custom converter
	parser.RegisterConverter(UserRoleConverter{})

	// Register custom validators
	parser.RegisterValidator(RegexValidator{})
	parser.RegisterValidator(UniqueUsernameValidator{
		ExistingUsers: []string{"admin", "test", "demo"},
	})

	// Register custom tag parser
	parser.RegisterTagParser(ExtendedTagParser{})

	// Example struct with custom validation
	type UserRegistration struct {
		Username string   `json:"username" validate:"required,minlen=3" businessrule:"unique_username"`
		Email    string   `json:"email" validate:"required,email"`
		Role     UserRole `json:"role"`
		Country  string   `json:"country" pattern:"^[A-Z]{2,3}$"`
		Tags     []string `json:"tags"`
	}

	// Test with form data
	data := map[string]interface{}{
		"username": "newuser",
		"email":    "newuser@example.com",
		"role":     "admin",
		"country":  "US",
		"tags":     []string{"developer", "golang"},
	}

	var user UserRegistration
	err := parser.ParseFrom(core.MapSource{Data: data}, &user)

	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	fmt.Printf("User registered: %+v\n", user)
	// Output: User registered: {Username:newuser Email:newuser@example.com Role:admin Country:US Tags:[developer golang]}
}

// Example demonstrating multiple parser instances with different configurations
func ExampleMultipleParsers() {
	// Strict parser for admin operations
	strictParser := core.NewParser()
	strictParser.RegisterValidator(RegexValidator{})
	strictParser.RegisterConverter(UserRoleConverter{})

	// Lenient parser for user input
	lenientParser := core.NewParser()
	lenientParser.RegisterConverter(UserRoleConverter{})

	type AdminConfig struct {
		APIKey string `json:"api_key" validate:"required" pattern:"^[A-Za-z0-9]{32}$"`
		Role   UserRole `json:"role"`
	}

	type UserProfile struct {
		Name string   `json:"name" validate:"required"`
		Role UserRole `json:"role"`
	}

	// Admin data with strict validation
	adminData := map[string]interface{}{
		"api_key": "abc123def456ghi789jkl012mno345pq",
		"role":    "admin",
	}

	var adminConfig AdminConfig
	err := strictParser.ParseFrom(core.MapSource{Data: adminData}, &adminConfig)
	if err != nil {
		fmt.Printf("Admin validation failed: %v\n", err)
	} else {
		fmt.Printf("Admin config valid: %+v\n", adminConfig)
	}

	// User data with lenient validation
	userData := map[string]interface{}{
		"name": "John Doe",
		"role": "user",
	}

	var userProfile UserProfile
	err = lenientParser.ParseFrom(core.MapSource{Data: userData}, &userProfile)
	if err != nil {
		fmt.Printf("User validation failed: %v\n", err)
	} else {
		fmt.Printf("User profile valid: %+v\n", userProfile)
	}
}

// Demonstrate CLI argument parsing with custom validators
func ExampleCLIWithCustomValidation() {
	parser := core.NewParser()
	parser.RegisterConverter(UserRoleConverter{})
	parser.RegisterValidator(RegexValidator{})

	type CLIConfig struct {
		Environment string   `json:"env" validate:"required" pattern:"^(dev|staging|prod)$"`
		Role        UserRole `json:"role"`
		Debug       bool     `json:"debug"`
		Workers     int      `json:"workers" validate:"min=1,max=10"`
	}

	// Simulate CLI arguments
	args := []string{
		"--env=prod",
		"--role=admin", 
		"--debug",
		"--workers=5",
	}

	var config CLIConfig
	err := parser.ParseFrom(core.NewArgsSource(args), &config)

	if err != nil {
		fmt.Printf("CLI validation error: %v\n", err)
		return
	}

	fmt.Printf("CLI config: %+v\n", config)
	// Output: CLI config: {Environment:prod Role:admin Debug:true Workers:5}
}
