package export_test

import (
	"strings"
	"testing"

	"github.com/envault/envault/internal/export"
)

func TestExportDotenv(t *testing.T) {
	e := export.New()
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := e.Export(vars, export.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=\"bar\"") {
		t.Errorf("expected FOO=\"bar\" in output, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=\"qux\"") {
		t.Errorf("expected BAZ=\"qux\" in output, got: %s", out)
	}
}

func TestExportShell(t *testing.T) {
	e := export.New()
	vars := map[string]string{"DB_URL": "postgres://localhost/db"}
	out, err := e.Export(vars, export.FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "export ") {
		t.Errorf("expected shell export prefix, got: %s", out)
	}
	if !strings.Contains(out, "DB_URL=") {
		t.Errorf("expected DB_URL in output, got: %s", out)
	}
}

func TestExportJSON(t *testing.T) {
	e := export.New()
	vars := map[string]string{"KEY": "value"}
	out, err := e.Export(vars, export.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"KEY"`) {
		t.Errorf("expected KEY in JSON output, got: %s", out)
	}
	if !strings.Contains(out, `"value"`) {
		t.Errorf("expected value in JSON output, got: %s", out)
	}
}

func TestExportSortedOutput(t *testing.T) {
	e := export.New()
	vars := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	out, err := e.Export(vars, export.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected Z_KEY last, got: %s", lines[2])
	}
}

func TestExportUnknownFormat(t *testing.T) {
	e := export.New()
	_, err := e.Export(map[string]string{"K": "v"}, export.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}
