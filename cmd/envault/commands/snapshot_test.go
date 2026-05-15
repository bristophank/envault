package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupSnapshotEnv(t *testing.T) (dir string, runCmd func(args ...string) (string, error)) {
	t.Helper()

	dir = t.TempDir()

	// Write a minimal .env file
	envContent := "DB_HOST=localhost\nDB_PORT=5432\nSECRET=abc123\n"
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	runCmd = func(args ...string) (string, error) {
		var buf bytes.Buffer
		root := &cobra.Command{Use: "envault"}
		snapshotCmd := NewSnapshotCmd()
		root.AddCommand(snapshotCmd)
		root.SetOut(&buf)
		root.SetErr(&buf)
		root.SetArgs(append([]string{"snapshot"}, args...))
		err := root.Execute()
		return buf.String(), err
	}

	return dir, runCmd
}

func TestSnapshotSaveCreatesFile(t *testing.T) {
	dir, _ := setupSnapshotEnv(t)

	envPath := filepath.Join(dir, ".env")
	snapshotDir := filepath.Join(dir, ".envault", "snapshots")

	var buf bytes.Buffer
	root := &cobra.Command{Use: "envault"}
	root.AddCommand(NewSnapshotCmd())
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"snapshot", "save", "--file", envPath, "--dir", snapshotDir, "--name", "v1"})

	if err := root.Execute(); err != nil {
		t.Fatalf("snapshot save failed: %v\noutput: %s", err, buf.String())
	}

	entries, err := os.ReadDir(snapshotDir)
	if err != nil {
		t.Fatalf("read snapshot dir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one snapshot file, got none")
	}
}

func TestSnapshotListShowsSnapshots(t *testing.T) {
	dir, _ := setupSnapshotEnv(t)

	envPath := filepath.Join(dir, ".env")
	snapshotDir := filepath.Join(dir, ".envault", "snapshots")

	// Save two snapshots
	for _, name := range []string{"snap-a", "snap-b"} {
		var buf bytes.Buffer
		root := &cobra.Command{Use: "envault"}
		root.AddCommand(NewSnapshotCmd())
		root.SetOut(&buf)
		root.SetErr(&buf)
		root.SetArgs([]string{"snapshot", "save", "--file", envPath, "--dir", snapshotDir, "--name", name})
		if err := root.Execute(); err != nil {
			t.Fatalf("snapshot save %s failed: %v", name, err)
		}
	}

	// List snapshots
	var buf bytes.Buffer
	root := &cobra.Command{Use: "envault"}
	root.AddCommand(NewSnapshotCmd())
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"snapshot", "list", "--dir", snapshotDir})
	if err := root.Execute(); err != nil {
		t.Fatalf("snapshot list failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "snap-a") || !strings.Contains(out, "snap-b") {
		t.Errorf("expected both snapshots in output, got: %s", out)
	}
}

func TestSnapshotRestoreWritesFile(t *testing.T) {
	dir, _ := setupSnapshotEnv(t)

	envPath := filepath.Join(dir, ".env")
	snapshotDir := filepath.Join(dir, ".envault", "snapshots")
	restorePath := filepath.Join(dir, ".env.restored")

	// Save snapshot
	{
		var buf bytes.Buffer
		root := &cobra.Command{Use: "envault"}
		root.AddCommand(NewSnapshotCmd())
		root.SetOut(&buf)
		root.SetErr(&buf)
		root.SetArgs([]string{"snapshot", "save", "--file", envPath, "--dir", snapshotDir, "--name", "restore-test"})
		if err := root.Execute(); err != nil {
			t.Fatalf("snapshot save failed: %v", err)
		}
	}

	// Restore snapshot
	{
		var buf bytes.Buffer
		root := &cobra.Command{Use: "envault"}
		root.AddCommand(NewSnapshotCmd())
		root.SetOut(&buf)
		root.SetErr(&buf)
		root.SetArgs([]string{"snapshot", "restore", "--dir", snapshotDir, "--name", "restore-test", "--output", restorePath})
		if err := root.Execute(); err != nil {
			t.Fatalf("snapshot restore failed: %v\noutput: %s", err, buf.String())
		}
	}

	data, err := os.ReadFile(restorePath)
	if err != nil {
		t.Fatalf("read restored file: %v", err)
	}
	if !strings.Contains(string(data), "DB_HOST=localhost") {
		t.Errorf("restored file missing expected content, got: %s", string(data))
	}
}
