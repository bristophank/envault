package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envault/internal/config"
)

func newTempManager(t *testing.T) *config.Manager {
	t.Helper()
	dir := t.TempDir()
	m, err := config.New(dir)
	if err != nil {
		t.Fatalf("config.New: %v", err)
	}
	return m
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	m := newTempManager(t)
	cfg, err := m.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultIdentityFile != "" || cfg.DefaultRecipientsFile != "" {
		t.Errorf("expected empty config, got %+v", cfg)
	}
}

func TestSaveAndLoad(t *testing.T) {
	m := newTempManager(t)
	want := &config.Config{
		DefaultIdentityFile:   "/home/user/.envault/identity.age",
		DefaultRecipientsFile: "/home/user/.envault/recipients.txt",
	}
	if err := m.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := m.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.DefaultIdentityFile != want.DefaultIdentityFile {
		t.Errorf("DefaultIdentityFile: got %q, want %q", got.DefaultIdentityFile, want.DefaultIdentityFile)
	}
	if got.DefaultRecipientsFile != want.DefaultRecipientsFile {
		t.Errorf("DefaultRecipientsFile: got %q, want %q", got.DefaultRecipientsFile, want.DefaultRecipientsFile)
	}
}

func TestSaveCreatesParentDirectories(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "envault")
	m, err := config.New(dir)
	if err != nil {
		t.Fatalf("config.New: %v", err)
	}
	if err := m.Save(&config.Config{}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(m.Path()); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}

func TestConfigFilePermissions(t *testing.T) {
	m := newTempManager(t)
	if err := m.Save(&config.Config{DefaultIdentityFile: "key.age"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	info, err := os.Stat(m.Path())
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}

func TestPathReturnsConfigFilePath(t *testing.T) {
	dir := t.TempDir()
	m, _ := config.New(dir)
	if m.Path() == "" {
		t.Error("Path() returned empty string")
	}
}
