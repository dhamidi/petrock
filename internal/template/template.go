package template

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/dhamidi/petrock/internal/utils" // Use correct module path
)

//go:embed all:templates
var Templates embed.FS

// RenderTemplate parses a template from the embedded FS, executes it with data,
// and writes the output to targetPath.
func RenderTemplate(fs embed.FS, targetPath string, templateName string, data interface{}) error {
	// Construct the full path within the embedded FS, always using forward slashes.
	// The embed directive was `all:templates`, so paths start with `templates/`.
	fullTemplatePath := filepath.Join("templates", templateName)
	fullTemplatePath = filepath.ToSlash(fullTemplatePath) // Ensure forward slashes for embed FS

	tmpl, err := template.ParseFS(fs, fullTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s from embedded FS: %w", fullTemplatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", fullTemplatePath, err)
	}

	// Ensure the directory for the target file exists before writing
	if err := utils.EnsureDir(filepath.Dir(targetPath)); err != nil {
		return fmt.Errorf("failed to ensure directory for %s: %w", targetPath, err)
	}

	if err := utils.WriteFile(targetPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write rendered template to %s: %w", targetPath, err)
	}

	return nil
}
