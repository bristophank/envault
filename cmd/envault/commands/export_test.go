package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envault/envault/cmd/envault/commands"
)

func setupExportEnv(t *testing.T) (dir string) {
	t.Helper()
	dir = t.TempDir()
	t.Setenv("HOME", dir)

	// Reuse setupSealEnv logic: generate a key and seal a file.
	root := commands.Root()
	root.SetArgs([]string{"keys", "generate"})
	if err := root.Execute(); err != nil {
		t.Fatalf("keys generate: %v", err)
	}

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("EXPORT_KEY=hello\nANOTHER=world\n"), 0600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	sealedFile := filepath.Join(dir, ".env.age")
	root2 := commands.Root()
	root2.SetArgs([]string{"seal", "--input", envFile, "--output", sealedFile})
	if err := root2.Execute(); err != nil {
		t.Fatalf("seal: %v", err)
	}
	return dir
}

func TestExportDotenvFormat(t *testing.T) {
	dir := setupExportEnv(t)
	sealedFile := filepath.Join(dir, ".env.age")

	var buf bytes.Buffer
	cmd := commands.Root()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"export", "--format", "dotenv", "--input", sealedFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("export: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "EXPORT_KEY=") {
		t.Errorf("expected EXPORT_KEY in output, got: %s", out)
	}
}

func TestExportShellFormat(t *testing.T) {
	dir := setupExportEnv(t)
	sealedFile := filepath.Join(dir, ".env.age")

	var buf bytes.Buffer
	cmd := commands.Root()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"export", "--format", "shell", "--input", sealedFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("export shell: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export EXPORT_KEY=") {
		t.Errorf("expected shell export syntax, got: %s", out)
	}
}

func TestExportJSONFormat(t *testing.T) {
	dir := setupExportEnv(t)
	sealedFile := filepath.Join(dir, ".env.age")

	var buf bytes.Buffer
	cmd := commands.Root()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"export", "--format", "json", "--input", sealedFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("export json: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}
