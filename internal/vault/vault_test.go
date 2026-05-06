package vault_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envault/envault/internal/crypto"
	"github.com/envault/envault/internal/vault"
)

func TestSealUnsealRoundtrip(t *testing.T) {
	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	vaultFile := filepath.Join(dir, ".env.vault")
	outFile := filepath.Join(dir, ".env.out")

	original := "DB_HOST=localhost\nDB_PORT=5432\nSECRET=supersecret\n"
	if err := os.WriteFile(envFile, []byte(original), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	v := vault.New(dir)

	if err := v.Seal(envFile, []string{pub}, vaultFile); err != nil {
		t.Fatalf("seal: %v", err)
	}

	if _, err := os.Stat(vaultFile); err != nil {
		t.Fatalf("vault file not created: %v", err)
	}

	plaintext, err := v.Unseal(vaultFile, priv, outFile)
	if err != nil {
		t.Fatalf("unseal: %v", err)
	}

	if plaintext != original {
		t.Errorf("roundtrip mismatch:\ngot:  %q\nwant: %q", plaintext, original)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if string(data) != original {
		t.Errorf("written file mismatch:\ngot:  %q\nwant: %q", string(data), original)
	}
}

func TestSealInvalidEnvFile(t *testing.T) {
	pub, _, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	if err := os.WriteFile(envFile, []byte("!!!invalid"), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	v := vault.New(dir)
	err = v.Seal(envFile, []string{pub}, "")
	if err == nil {
		t.Fatal("expected error for invalid env file, got nil")
	}
}

func TestUnsealNoDstReturnsContent(t *testing.T) {
	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	vaultFile := filepath.Join(dir, ".env.vault")

	original := "API_KEY=abc123\n"
	if err := os.WriteFile(envFile, []byte(original), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	v := vault.New(dir)
	if err := v.Seal(envFile, []string{pub}, vaultFile); err != nil {
		t.Fatalf("seal: %v", err)
	}

	plaintext, err := v.Unseal(vaultFile, priv, "")
	if err != nil {
		t.Fatalf("unseal: %v", err)
	}

	if !strings.Contains(plaintext, "API_KEY=abc123") {
		t.Errorf("expected API_KEY in plaintext, got: %q", plaintext)
	}
}
