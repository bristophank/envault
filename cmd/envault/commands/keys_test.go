package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envault/envault/cmd/envault/commands"
)

func TestKeysGenerateCreatesKey(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cmd := commands.NewKeysCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	cmd.SetArgs([]string{"generate"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	keyPath := filepath.Join(tmpDir, ".envault", "key.txt")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Fatal("expected key file to be created")
	}

	output := buf.String()
	if !strings.Contains(output, "Public key:") {
		t.Errorf("expected output to contain public key, got: %s", output)
	}
}

func TestKeysGenerateFailsIfKeyExists(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// First generation should succeed.
	cmd := commands.NewKeysCmd()
	cmd.SetArgs([]string{"generate"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("first generate failed: %v", err)
	}

	// Second generation without --force should fail.
	cmd2 := commands.NewKeysCmd()
	buf := &bytes.Buffer{}
	cmd2.SetOut(buf)
	cmd2.SetErr(buf)
	cmd2.SetArgs([]string{"generate"})
	if err := cmd2.Execute(); err == nil {
		t.Fatal("expected error when key already exists")
	}
}

func TestKeysGenerateForceOverwrites(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cmd := commands.NewKeysCmd()
	cmd.SetArgs([]string{"generate"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("first generate failed: %v", err)
	}

	cmd2 := commands.NewKeysCmd()
	cmd2.SetArgs([]string{"generate", "--force"})
	if err := cmd2.Execute(); err != nil {
		t.Fatalf("force generate failed: %v", err)
	}
}

func TestKeysShowWithoutKeyFails(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cmd := commands.NewKeysCmd()
	cmd.SetArgs([]string{"show"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when no key exists")
	}
}

func TestKeysShowAfterGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	genCmd := commands.NewKeysCmd()
	genCmd.SetArgs([]string{"generate"})
	if err := genCmd.Execute(); err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	buf := &bytes.Buffer{}
	showCmd := commands.NewKeysCmd()
	showCmd.SetOut(buf)
	showCmd.SetArgs([]string{"show"})
	if err := showCmd.Execute(); err != nil {
		t.Fatalf("show failed: %v", err)
	}

	if !strings.Contains(buf.String(), "Public key:") {
		t.Errorf("expected public key in output, got: %s", buf.String())
	}
}
