// Package export provides functionality for exporting decrypted env
// variables to various output formats (shell, JSON, dotenv).
package export

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents an output format for exported variables.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatJSON   Format = "json"
)

// Exporter converts a map of env variables to a specific format.
type Exporter struct{}

// New returns a new Exporter.
func New() *Exporter {
	return &Exporter{}
}

// Export serialises vars into the requested format.
func (e *Exporter) Export(vars map[string]string, format Format) (string, error) {
	switch format {
	case FormatDotenv:
		return e.toDotenv(vars), nil
	case FormatShell:
		return e.toShell(vars), nil
	case FormatJSON:
		return e.toJSON(vars)
	default:
		return "", fmt.Errorf("unknown format %q: must be one of dotenv, shell, json", format)
	}
}

func (e *Exporter) toDotenv(vars map[string]string) string {
	keys := sortedKeys(vars)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%q\n", k, vars[k])
	}
	return sb.String()
}

func (e *Exporter) toShell(vars map[string]string) string {
	keys := sortedKeys(vars)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, vars[k])
	}
	return sb.String()
}

func (e *Exporter) toJSON(vars map[string]string) (string, error) {
	b, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b) + "\n", nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
