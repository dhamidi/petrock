// This go.mod is minimal and primarily for validating the template code.
// It will NOT be copied to the generated feature directory.
module petrock_internal/feature_template

go 1.23

require (
	// Placeholder for the core package dependency of the target project.
	// ReplaceInFiles will update this path in the generated .go files,
	// but this require directive itself is ignored during generation.
	petrock_example_module_path/core v0.0.0

	// Add other direct dependencies the template itself needs for compilation, if any.
	// Example: maragu.dev/gomponents v0.20.1
)

// Replace directive to satisfy the core dependency locally within the skeleton context.
// This assumes the main skeleton has a core package.
replace petrock_example_module_path/core => ../core
```

**2. Create `internal/skeleton/feature_template/register.go`**

```go
internal/skeleton/feature_template/register.go
<<<<<<< SEARCH
