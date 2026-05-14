// Package lint provides validation for .env file contents,
// checking for common issues like missing values, duplicate keys,
// and insecure patterns.
package lint

import (
	"fmt"
	"strings"

	"github.com/envault/envault/internal/env"
)

// Severity indicates how serious a lint issue is.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Issue represents a single lint finding.
type Issue struct {
	Line     int
	Key      string
	Message  string
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("%s [line %d] %s: %s", i.Severity, i.Line, i.Key, i.Message)
}

// Result holds all issues found during linting.
type Result struct {
	Issues []Issue
}

// HasErrors returns true if any error-severity issues exist.
func (r *Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Lint parses and validates the given .env content, returning a Result.
func Lint(content string) (*Result, error) {
	entries, err := env.Parse(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	result := &Result{}
	seen := make(map[string]int)

	for i, entry := range entries {
		lineNum := i + 1

		// Duplicate key check
		if prev, ok := seen[entry.Key]; ok {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  fmt.Sprintf("duplicate key (first seen at line %d)", prev),
				Severity: SeverityError,
			})
		}
		seen[entry.Key] = lineNum

		// Empty value warning
		if entry.Value == "" {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  "empty value",
				Severity: SeverityWarning,
			})
		}

		// Detect keys that look like placeholders
		val := strings.ToLower(entry.Value)
		if val == "changeme" || val == "todo" || val == "fixme" || val == "placeholder" {
			result.Issues = append(result.Issues, Issue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  fmt.Sprintf("suspicious placeholder value %q", entry.Value),
				Severity: SeverityWarning,
			})
		}
	}

	return result, nil
}
