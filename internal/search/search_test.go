package search_test

import (
	"strings"
	"testing"

	"github.com/envault/envault/internal/search"
)

func TestSearchKeysFindsMatch(t *testing.T) {
	s := search.New()
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "secret",
		"PORT":         "8080",
	}
	results := s.SearchKeys(env, "test.env", "api")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", results[0].Key)
	}
}

func TestSearchKeysNoMatch(t *testing.T) {
	s := search.New()
	env := map[string]string{"PORT": "8080"}
	results := s.SearchKeys(env, "test.env", "database")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchValuesFindMatch(t *testing.T) {
	s := search.New()
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"REDIS_URL":    "redis://localhost:6379",
		"PORT":         "8080",
	}
	results := s.SearchValues(env, "test.env", "localhost")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchIsCaseInsensitive(t *testing.T) {
	s := search.New()
	env := map[string]string{"API_KEY": "MySecretValue"}
	results := s.SearchKeys(env, "", "api_key")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestFormatNoResults(t *testing.T) {
	out := search.Format(nil)
	if out != "no matches found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatWithResults(t *testing.T) {
	results := []search.Result{
		{Key: "API_KEY", Value: "secret", File: ".env"},
	}
	out := search.Format(results)
	if !strings.Contains(out, "API_KEY=secret") {
		t.Errorf("expected API_KEY=secret in output, got %q", out)
	}
	if !strings.Contains(out, ".env:") {
		t.Errorf("expected file prefix in output, got %q", out)
	}
}
