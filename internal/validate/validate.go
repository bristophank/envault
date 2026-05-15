// Package validate provides schema-based validation for .env files.
// It checks that required keys are present and that values match expected patterns.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for a single environment variable.
type Rule struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Result holds the outcome of a validation check.
type Result struct {
	Key     string
	Message string
}

func (r Result) Error() string {
	return fmt.Sprintf("%s: %s", r.Key, r.Message)
}

// Schema is a collection of validation rules.
type Schema []Rule

// Validate checks the provided env map against the schema and returns any violations.
func Validate(env map[string]string, schema Schema) []Result {
	var results []Result

	for _, rule := range schema {
		val, exists := env[rule.Key]

		if rule.Required && (!exists || strings.TrimSpace(val) == "") {
			results = append(results, Result{
				Key:     rule.Key,
				Message: "required key is missing or empty",
			})
			continue
		}

		if exists && rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				results = append(results, Result{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				results = append(results, Result{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
			}
		}
	}

	return results
}

// ParseSchema parses a simple schema file where each line is:
//   KEY [required] [pattern=REGEX]
func ParseSchema(lines []string) (Schema, error) {
	var schema Schema
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		rule := Rule{Key: parts[0]}
		for _, part := range parts[1:] {
			switch {
			case part == "required":
				rule.Required = true
			case strings.HasPrefix(part, "pattern="):
				rule.Pattern = strings.TrimPrefix(part, "pattern=")
			default:
				return nil, fmt.Errorf("line %d: unknown directive %q", i+1, part)
			}
		}
		schema = append(schema, rule)
	}
	return schema, nil
}
