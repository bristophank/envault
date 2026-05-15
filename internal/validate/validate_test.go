package validate_test

import (
	"testing"

	"github.com/envault/envault/internal/validate"
)

func TestValidateMissingRequired(t *testing.T) {
	schema := validate.Schema{
		{Key: "DATABASE_URL", Required: true},
	}
	results := validate.Validate(map[string]string{}, schema)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected key %q", results[0].Key)
	}
}

func TestValidateEmptyValueFails(t *testing.T) {
	schema := validate.Schema{
		{Key: "SECRET", Required: true},
	}
	results := validate.Validate(map[string]string{"SECRET": "   "}, schema)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestValidatePatternMatch(t *testing.T) {
	schema := validate.Schema{
		{Key: "PORT", Pattern: `^\d+$`},
	}
	results := validate.Validate(map[string]string{"PORT": "8080"}, schema)
	if len(results) != 0 {
		t.Errorf("expected no errors, got %v", results)
	}
}

func TestValidatePatternMismatch(t *testing.T) {
	schema := validate.Schema{
		{Key: "PORT", Pattern: `^\d+$`},
	}
	results := validate.Validate(map[string]string{"PORT": "not-a-number"}, schema)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestValidateNoErrors(t *testing.T) {
	schema := validate.Schema{
		{Key: "API_KEY", Required: true},
		{Key: "TIMEOUT", Pattern: `^\d+$`},
	}
	env := map[string]string{"API_KEY": "abc123", "TIMEOUT": "30"}
	results := validate.Validate(env, schema)
	if len(results) != 0 {
		t.Errorf("expected no errors, got %v", results)
	}
}

func TestParseSchema(t *testing.T) {
	lines := []string{
		"# comment",
		"DATABASE_URL required",
		"PORT pattern=^\\d+$",
		"DEBUG",
	}
	schema, err := validate.ParseSchema(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(schema))
	}
	if !schema[0].Required {
		t.Error("DATABASE_URL should be required")
	}
	if schema[1].Pattern != `^\d+$` {
		t.Errorf("unexpected pattern: %q", schema[1].Pattern)
	}
}

func TestParseSchemaUnknownDirective(t *testing.T) {
	lines := []string{"KEY unknown_directive"}
	_, err := validate.ParseSchema(lines)
	if err == nil {
		t.Error("expected error for unknown directive")
	}
}
