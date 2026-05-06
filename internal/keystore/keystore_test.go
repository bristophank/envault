package keystore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envault/internal/keystore"
)

func newTempStore(t *testing.T) *keystore.Store {
	t.Helper()
	dir := t.TempDir()
	s := keystore.New(dir)
	if err := s.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	return s
}

func TestInitCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	s := keystore.New(dir)
	if err := s.Init(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(dir, keystore.DefaultKeysDir)
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist", expected)
	}
}

func TestSaveLoadPrivateKey(t *testing.T) {
	s := newTempStore(t)
	key := "AGE-SECRET-KEY-1EXAMPLEKEY"
	if err := s.SavePrivateKey(key); err != nil {
		t.Fatalf("SavePrivateKey: %v", err)
	}
	got, err := s.LoadPrivateKey()
	if err != nil {
		t.Fatalf("LoadPrivateKey: %v", err)
	}
	if got != key {
		t.Errorf("got %q, want %q", got, key)
	}
}

func TestLoadPrivateKeyMissing(t *testing.T) {
	s := newTempStore(t)
	_, err := s.LoadPrivateKey()
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestAddAndLoadRecipients(t *testing.T) {
	s := newTempStore(t)
	recipients := []string{
		"age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
		"age1lggyhqrw2nlhcxprm67z43rta597azn8gknawjehu9d9dl0jq3yqqvfafg",
	}
	for _, r := range recipients {
		if err := s.AddRecipient(r); err != nil {
			t.Fatalf("AddRecipient(%q): %v", r, err)
		}
	}
	got, err := s.LoadRecipients()
	if err != nil {
		t.Fatalf("LoadRecipients: %v", err)
	}
	if len(got) != len(recipients) {
		t.Fatalf("got %d recipients, want %d", len(got), len(recipients))
	}
	for i, r := range recipients {
		if got[i] != r {
			t.Errorf("recipient[%d]: got %q, want %q", i, got[i], r)
		}
	}
}

func TestLoadRecipientsEmpty(t *testing.T) {
	s := newTempStore(t)
	got, err := s.LoadRecipients()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty recipients, got %v", got)
	}
}
