package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseEnvContent verifies that parseEnvContent correctly extracts
// KEY=VALUE pairs and skips comments / blank lines.
func TestParseEnvContent(t *testing.T) {
	input := `
# This is a comment
DB_HOST=localhost
DB_PORT=5432

SECRET_KEY=supersecret
`
	got := parseEnvContent(input)
	if len(got) != 3 {
		t.Fatalf("expected 3 pairs, got %d: %v", len(got), got)
	}
	expected := []string{"DB_HOST=localhost", "DB_PORT=5432", "SECRET_KEY=supersecret"}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("pair[%d]: want %q, got %q", i, e, got[i])
		}
	}
}

func TestParseEnvContentEmptyInput(t *testing.T) {
	got := parseEnvContent("")
	if len(got) != 0 {
		t.Fatalf("expected 0 pairs, got %d", len(got))
	}
}

func TestParseEnvContentOnlyComments(t *testing.T) {
	input := "# comment one\n# comment two\n"
	got := parseEnvContent(input)
	if len(got) != 0 {
		t.Fatalf("expected 0 pairs, got %d", len(got))
	}
}

// setupEnvCmd creates a temp dir with a sealed .env.age file so that the
// full env command can be exercised in integration-style tests.
func setupEnvCmd(t *testing.T) (keysDir, sealedFile string) {
	t.Helper()
	dir := t.TempDir()

	// Reuse the helpers already defined in seal_unseal_test.go.
	env := setupSealEnv(t)

	// Seal the .env file so we have a valid .env.age to work with.
	sealCmd := NewSealCmd()
	sealCmd.SetArgs([]string{
		"--file", env.envFile,
		"--output", filepath.Join(dir, ".env.age"),
		"--keys-dir", env.keysDir,
	})
	if err := sealCmd.Execute(); err != nil {
		t.Fatalf("seal: %v", err)
	}

	return env.keysDir, filepath.Join(dir, ".env.age")
}

func TestEnvInjectsVariables(t *testing.T) {
	keysDir, sealedFile := setupEnvCmd(t)

	// Write a small shell script that prints a specific var.
	scriptDir := t.TempDir()
	script := filepath.Join(scriptDir, "check.sh")
	_ = os.WriteFile(script, []byte("#!/bin/sh\necho $TEST_VAR\n"), 0o755)

	// Capture stdout by redirecting through a temp file.
	outFile := filepath.Join(scriptDir, "out.txt")

	err := runEnv(keysDir, sealedFile, []string{
		"sh", "-c", "echo $TEST_VAR > " + outFile,
	})
	if err != nil {
		t.Fatalf("runEnv: %v", err)
	}

	out, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(out), "world") {
		t.Errorf("expected TEST_VAR=world in output, got: %q", string(out))
	}
}

func TestEnvMissingSealedFile(t *testing.T) {
	dir := t.TempDir()
	err := runEnv(dir, filepath.Join(dir, "nonexistent.age"), []string{"true"})
	if err == nil {
		t.Fatal("expected error for missing sealed file")
	}
}
