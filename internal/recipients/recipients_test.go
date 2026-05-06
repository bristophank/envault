package recipients_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envault/internal/recipients"
)

const (
	key1 = "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"
	key2 = "age1cy0su9fwf3gf9mw868g5yut09p6nytfmmnktexz9yxy2pdg5e6as8q2hhd"
)

func TestAddAndList(t *testing.T) {
	l := recipients.New()
	if err := l.Add(key1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := l.Add(key2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Len() != 2 {
		t.Fatalf("expected 2 recipients, got %d", l.Len())
	}
}

func TestAddDuplicateIsIdempotent(t *testing.T) {
	l := recipients.New()
	_ = l.Add(key1)
	_ = l.Add(key1)
	if l.Len() != 1 {
		t.Fatalf("expected 1 recipient after duplicate add, got %d", l.Len())
	}
}

func TestAddInvalidKey(t *testing.T) {
	l := recipients.New()
	if err := l.Add(""); err == nil {
		t.Fatal("expected error for empty key")
	}
	if err := l.Add("notanagekey"); err == nil {
		t.Fatal("expected error for malformed key")
	}
}

func TestRemove(t *testing.T) {
	l := recipients.New()
	_ = l.Add(key1)
	_ = l.Add(key2)

	if err := l.Remove(key1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Len() != 1 {
		t.Fatalf("expected 1 recipient after remove, got %d", l.Len())
	}
	if l.Keys()[0] != key2 {
		t.Fatalf("expected remaining key to be key2")
	}
}

func TestRemoveNotFound(t *testing.T) {
	l := recipients.New()
	if err := l.Remove(key1); err == nil {
		t.Fatal("expected error when removing non-existent key")
	}
}

func TestSaveAndLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, recipients.File)

	l := recipients.New()
	_ = l.Add(key1)
	_ = l.Add(key2)

	if err := recipients.SaveFile(path, l); err != nil {
		t.Fatalf("SaveFile: %v", err)
	}

	loaded, err := recipients.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if loaded.Len() != 2 {
		t.Fatalf("expected 2 recipients after load, got %d", loaded.Len())
	}
	for i, k := range loaded.Keys() {
		if k != l.Keys()[i] {
			t.Errorf("key[%d] mismatch: got %q, want %q", i, k, l.Keys()[i])
		}
	}
}

func TestLoadFileMissing(t *testing.T) {
	_, err := recipients.LoadFile("/nonexistent/path/.envault-recipients")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFileIgnoresComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, recipients.File)

	content := "# this is a comment\n" + key1 + "\n\n# another comment\n" + key2 + "\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	l, err := recipients.LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if l.Len() != 2 {
		t.Fatalf("expected 2 recipients, got %d", l.Len())
	}
}
