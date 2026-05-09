package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envault/envault/internal/audit"
)

func newTempLogger(t *testing.T) (*audit.Logger, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "audit", "audit.log")
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("audit.New: %v", err)
	}
	return l, path
}

func TestLogCreatesFile(t *testing.T) {
	l, path := newTempLogger(t)
	if err := l.Log("seal", "secrets.env.age", ""); err != nil {
		t.Fatalf("Log: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("log file not created: %v", err)
	}
}

func TestLogAndReadAll(t *testing.T) {
	l, _ := newTempLogger(t)

	ops := []struct{ op, target, details string }{
		{"seal", "secrets.env.age", "encrypted with 2 recipients"},
		{"unseal", "secrets.env.age", ""},
		{"rotate", "secrets.env.age", "re-encrypted"},
	}
	for _, o := range ops {
		if err := l.Log(o.op, o.target, o.details); err != nil {
			t.Fatalf("Log(%s): %v", o.op, err)
		}
	}

	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != len(ops) {
		t.Fatalf("expected %d entries, got %d", len(ops), len(entries))
	}
	for i, e := range entries {
		if e.Operation != ops[i].op {
			t.Errorf("entry %d: op = %q, want %q", i, e.Operation, ops[i].op)
		}
		if e.Target != ops[i].target {
			t.Errorf("entry %d: target = %q, want %q", i, e.Target, ops[i].target)
		}
		if e.Timestamp.IsZero() {
			t.Errorf("entry %d: timestamp is zero", i)
		}
	}
}

func TestReadAllMissingFileReturnsEmpty(t *testing.T) {
	l, _ := newTempLogger(t)
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestLogAppendsAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l1, _ := audit.New(path)
	if err := l1.Log("seal", "a.env.age", ""); err != nil {
		t.Fatal(err)
	}

	l2, _ := audit.New(path)
	if err := l2.Log("unseal", "a.env.age", ""); err != nil {
		t.Fatal(err)
	}

	entries, err := l2.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
