package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/envault/envault/internal/history"
)

func newTempHistory(t *testing.T) *history.History {
	t.Helper()
	dir, err := os.MkdirTemp("", "history-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return history.New(dir)
}

func TestListMissingFileReturnsEmpty(t *testing.T) {
	h := newTempHistory(t)
	entries, err := h.List(".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty, got %d entries", len(entries))
	}
}

func TestRecordAndList(t *testing.T) {
	h := newTempHistory(t)

	if err := h.Record("seal", ".env", "alice", "initial seal"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := h.Record("rotate", ".env", "bob", ""); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := h.List(".env")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Operation != "seal" {
		t.Errorf("expected 'seal', got %q", entries[0].Operation)
	}
	if entries[1].User != "bob" {
		t.Errorf("expected user 'bob', got %q", entries[1].User)
	}
	if entries[0].Note != "initial seal" {
		t.Errorf("expected note 'initial seal', got %q", entries[0].Note)
	}
}

func TestTimestampIsRecent(t *testing.T) {
	h := newTempHistory(t)
	before := time.Now().UTC().Add(-time.Second)
	if err := h.Record("unseal", ".env", "carol", ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	after := time.Now().UTC().Add(time.Second)

	entries, _ := h.List(".env")
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
}

func TestSeparateLogsPerFile(t *testing.T) {
	h := newTempHistory(t)
	_ = h.Record("seal", ".env", "alice", "")
	_ = h.Record("seal", ".env.staging", "alice", "")
	_ = h.Record("seal", ".env.staging", "alice", "")

	prod, _ := h.List(".env")
	staging, _ := h.List(".env.staging")

	if len(prod) != 1 {
		t.Errorf("expected 1 prod entry, got %d", len(prod))
	}
	if len(staging) != 2 {
		t.Errorf("expected 2 staging entries, got %d", len(staging))
	}
}
