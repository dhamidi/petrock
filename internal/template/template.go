package template

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"text/template"

	"petrock/internal/utils" // Assuming utils package is at this path
)

//go:embed all:templates
var Templates embed.FS

// RenderTemplate parses a template from the embedded FS, executes it with data,
// and writes the output to targetPath.
func RenderTemplate(fs embed.FS, targetPath string, templateName string, data interface{}) error {
	// Ensure template name uses forward slashes, compatible with embed.FS
	templatePath := filepath.ToSlash(templateName)

	tmpl, err := template.ParseFS(fs, templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s from embedded FS: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
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
