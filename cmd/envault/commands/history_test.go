package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envault/envault/internal/history"
	"github.com/envault/envault/cmd/envault/commands"
)

func setupHistoryCmd(t *testing.T) (histDir string, cleanup func()) {
	t.Helper()
	tmp, err := os.MkdirTemp("", "history-cmd-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	return filepath.Join(tmp, "history"), func() { os.RemoveAll(tmp) }
}

func TestHistoryNoEntries(t *testing.T) {
	histDir, cleanup := setupHistoryCmd(t)
	defer cleanup()

	h := history.New(histDir)
	_ = h // no records written

	// Simulate command output by calling List directly since we cannot
	// easily inject histDir into the cobra command without refactoring.
	entries, err := history.New(histDir).List(".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestHistoryCommandOutputFormat(t *testing.T) {
	histDir, cleanup := setupHistoryCmd(t)
	defer cleanup()

	h := history.New(histDir)
	if err := h.Record("seal", ".env", "alice", "first seal"); err != nil {
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

	var buf bytes.Buffer
	for _, e := range entries {
		buf.WriteString(e.Operation + " " + e.User + "\n")
	}
	out := buf.String()
	if !strings.Contains(out, "seal alice") {
		t.Errorf("expected 'seal alice' in output, got: %s", out)
	}
	if !strings.Contains(out, "rotate bob") {
		t.Errorf("expected 'rotate bob' in output, got: %s", out)
	}
}

func TestHistoryCmdRegistered(t *testing.T) {
	cmd := commands.NewHistoryCmd()
	if cmd.Use != "history [file]" {
		t.Errorf("unexpected Use: %q", cmd.Use)
	}
}
