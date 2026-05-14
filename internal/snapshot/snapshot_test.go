package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/snapshot"
)

func newTempSnapshot(t *testing.T) (*snapshot.Manager, string) {
	t.Helper()
	base := t.TempDir()
	snapsDir := filepath.Join(base, "snapshots")
	return snapshot.New(snapsDir), base
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func TestSaveCreatesSnapshot(t *testing.T) {
	m, base := newTempSnapshot(t)
	src := filepath.Join(base, ".env.age")
	writeFile(t, src, "encrypted-content")

	snap, err := m.Save(src)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(snap); err != nil {
		t.Fatalf("snapshot file missing: %v", err)
	}
}

func TestListReturnsSnapshots(t *testing.T) {
	m, base := newTempSnapshot(t)
	src := filepath.Join(base, ".env.age")
	writeFile(t, src, "data")

	if _, err := m.Save(src); err != nil {
		t.Fatalf("Save 1: %v", err)
	}
	if _, err := m.Save(src); err != nil {
		t.Fatalf("Save 2: %v", err)
	}

	paths, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(paths))
	}
}

func TestListMissingDirReturnsEmpty(t *testing.T) {
	m := snapshot.New(filepath.Join(t.TempDir(), "nonexistent"))
	paths, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(paths) != 0 {
		t.Fatalf("expected empty list, got %v", paths)
	}
}

func TestRestoreRecreatesFile(t *testing.T) {
	m, base := newTempSnapshot(t)
	src := filepath.Join(base, ".env.age")
	writeFile(t, src, "original")

	snap, err := m.Save(src)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Overwrite the source, then restore.
	writeFile(t, src, "overwritten")
	if err := m.Restore(snap, src); err != nil {
		t.Fatalf("Restore: %v", err)
	}

	got, _ := os.ReadFile(src)
	if string(got) != "original" {
		t.Fatalf("expected 'original', got %q", got)
	}
}
