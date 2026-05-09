package commands_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envault/cmd/envault/commands"
	"github.com/yourusername/envault/internal/crypto"
	"github.com/yourusername/envault/internal/keystore"
)

func setupDiffEnv(t *testing.T) (dir string, pubKey string) {
	t.Helper()
	dir = t.TempDir()
	t.Setenv("ENVAULT_DIR", dir)

	ks := keystore.New(dir)
	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}
	if err := ks.SavePrivateKey(priv); err != nil {
		t.Fatalf("save private key: %v", err)
	}
	return dir, pub
}

func sealContent(t *testing.T, dir, filename, content, pubKey string) string {
	t.Helper()
	ciphertext, err := crypto.Encrypt(content, []string{pubKey})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(ciphertext), 0600); err != nil {
		t.Fatalf("write sealed file: %v", err)
	}
	return path
}

func TestDiffShowsAddedKey(t *testing.T) {
	dir, pub := setupDiffEnv(t)

	pathA := sealContent(t, dir, "a.env.age", "FOO=bar\n", pub)
	pathB := sealContent(t, dir, "b.env.age", "FOO=bar\nNEW=val\n", pub)

	cmd := commands.Root()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"diff", pathA, pathB})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("diff command failed: %v", err)
	}
	if !strings.Contains(buf.String(), "NEW") {
		t.Errorf("expected NEW in diff output, got: %s", buf.String())
	}
}

func TestDiffNoChanges(t *testing.T) {
	dir, pub := setupDiffEnv(t)

	pathA := sealContent(t, dir, "a.env.age", "FOO=bar\n", pub)
	pathB := sealContent(t, dir, "b.env.age", "FOO=bar\n", pub)

	cmd := commands.Root()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"diff", pathA, pathB})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("diff command failed: %v", err)
	}
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got: %s", buf.String())
	}
}
