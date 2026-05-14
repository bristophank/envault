package template_test

import (
	"strings"
	"testing"

	"github.com/envault/envault/internal/template"
)

func TestRenderSubstitutesValues(t *testing.T) {
	tmpl := `DB_HOST=localhost
DB_PORT=5432
DB_PASS=
`
	src := map[string]string{
		"DB_HOST": "prod.db.example.com",
		"DB_PASS": "supersecret",
	}

	r := template.New()
	out, err := r.RenderBytes([]byte(tmpl), src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST=prod.db.example.com") {
		t.Errorf("expected substituted DB_HOST, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASS=supersecret") {
		t.Errorf("expected substituted DB_PASS, got:\n%s", out)
	}
	// DB_PORT not in src, should keep original value
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Errorf("expected original DB_PORT, got:\n%s", out)
	}
}

func TestRenderMissingRequiredKey(t *testing.T) {
	tmpl := `API_KEY= # required
OPTIONAL_KEY=
`
	r := template.New()
	_, err := r.RenderBytes([]byte(tmpl), map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	var mke *template.MissingKeyError
	if !isType(err, &mke) {
		t.Fatalf("expected MissingKeyError, got %T", err)
	}
	if len(mke.Keys) != 1 || mke.Keys[0] != "API_KEY" {
		t.Errorf("expected API_KEY in missing keys, got %v", mke.Keys)
	}
}

func TestRenderRequiredKeyProvided(t *testing.T) {
	tmpl := `API_KEY= # required
`
	src := map[string]string{"API_KEY": "abc123"}
	r := template.New()
	out, err := r.RenderBytes([]byte(tmpl), src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "API_KEY=abc123") {
		t.Errorf("expected API_KEY=abc123, got:\n%s", out)
	}
}

func TestRenderEmptyTemplate(t *testing.T) {
	r := template.New()
	out, err := r.RenderBytes([]byte(""), map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}

// isType is a helper to check error type via pointer assignment.
func isType(err error, target interface{}) bool {
	switch t := target.(type) {
	case **template.MissingKeyError:
		if v, ok := err.(*template.MissingKeyError); ok {
			*t = v
			return true
		}
	}
	return false
}
