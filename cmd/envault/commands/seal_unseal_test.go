package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envault/cmd/envault/commands"
)

func setupSealEnv(t *testing.T) (dir string, cleanup func()) {
	t.Helper()
	dir = t.TempDir()

	// Write a sample .env file
	envContent := "DB_HOST=localhost\nDB_PASS=secret\nAPI_KEY=abc123\n"
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(envContent), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	// Run init to create a key
	initCmd := commands.Root()
	initCmd.SetArgs([]string{"--dir", dir, "init"})
	if err := initCmd.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	return dir, func() {}
}

func TestSealCreatesEncryptedFile(t *testing.T) {
	dir, cleanup := setupSealEnv(t)
	defer cleanup()

	cmd := commands.Root()
	cmd.SetArgs([]string{
		"--dir", dir,
		"seal",
		"--in", filepath.Join(dir, ".env"),
		"--out", filepath.Join(dir, ".env.age"),
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("seal failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".env.age")); os.IsNotExist(err) {
		t.Fatal("expected .env.age to exist after seal")
	}
}

func TestUnsealRestoresContent(t *testing.T) {
	dir, cleanup := setupSealEnv(t)
	defer cleanup()

	sealCmd := commands.Root()
	sealCmd.SetArgs([]string{
		"--dir", dir,
		"seal",
		"--in", filepath.Join(dir, ".env"),
		"--out", filepath.Join(dir, ".env.age"),
	})
	if err := sealCmd.Execute(); err != nil {
		t.Fatalf("seal failed: %v", err)
	}

	outPath := filepath.Join(dir, ".env.decrypted")
	unsealCmd := commands.Root()
	unsealCmd.SetArgs([]string{
		"--dir", dir,
		"unseal",
		"--in", filepath.Join(dir, ".env.age"),
		"--out", outPath,
	})
	if err := unsealCmd.Execute(); err != nil {
		t.Fatalf("unseal failed: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	want := "DB_HOST=localhost\nDB_PASS=secret\nAPI_KEY=abc123\n"
	if string(got) != want {
		t.Errorf("decrypted content mismatch\ngot:  %q\nwant: %q", string(got), want)
	}
}

func TestSealMissingInputFile(t *testing.T) {
	dir, cleanup := setupSealEnv(t)
	defer cleanup()

	cmd := commands.Root()
	cmd.SetArgs([]string{
		"--dir", dir,
		"seal",
		"--in", filepath.Join(dir, "nonexistent.env"),
		"--out", filepath.Join(dir, ".env.age"),
	})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing input file, got nil")
	}
}
