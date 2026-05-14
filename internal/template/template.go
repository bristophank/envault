// Package template provides functionality for rendering .env files
// from a template with variable substitution and required key validation.
package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/envault/envault/internal/env"
)

// MissingKeyError is returned when a required key has no value.
type MissingKeyError struct {
	Keys []string
}

func (e *MissingKeyError) Error() string {
	return fmt.Sprintf("missing required keys: %s", strings.Join(e.Keys, ", "))
}

// Renderer applies values from a source env map to a template env file.
type Renderer struct{}

// New returns a new Renderer.
func New() *Renderer {
	return &Renderer{}
}

// Render reads a template file (keys with empty or placeholder values),
// substitutes values from src, and returns the rendered content.
// Keys marked with a "required" comment (# required) must be present in src.
func (r *Renderer) Render(templatePath string, src map[string]string) (string, error) {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("read template: %w", err)
	}
	return r.RenderBytes(data, src)
}

// RenderBytes performs substitution on raw template bytes.
func (r *Renderer) RenderBytes(data []byte, src map[string]string) (string, error) {
	entries, err := env.Parse(strings.NewReader(string(data)))
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var missing []string
	for i, entry := range entries {
		if val, ok := src[entry.Key]; ok {
			entries[i].Value = val
		} else if isRequired(entry.Comment) && entry.Value == "" {
			missing = append(missing, entry.Key)
		}
	}

	if len(missing) > 0 {
		return "", &MissingKeyError{Keys: missing}
	}

	return env.Serialize(entries), nil
}

// isRequired returns true if the comment contains the word "required".
func isRequired(comment string) bool {
	return strings.Contains(strings.ToLower(comment), "required")
}
