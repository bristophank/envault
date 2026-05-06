package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/envault/envault/internal/crypto"
	"github.com/envault/envault/internal/env"
)

const defaultVaultExt = ".vault"

// Vault manages encrypted .env files on disk.
type Vault struct {
	Dir string
}

// New returns a Vault rooted at dir.
func New(dir string) *Vault {
	return &Vault{Dir: dir}
}

// Seal encrypts the plaintext .env file at src and writes the ciphertext to
// dst (defaults to src + ".vault") using the provided age recipient public keys.
func (v *Vault) Seal(src string, recipients []string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("vault: read source: %w", err)
	}

	// Validate the env file is parseable before encrypting.
	if _, err := env.Parse(string(data)); err != nil {
		return fmt.Errorf("vault: invalid env file: %w", err)
	}

	ciphertext, err := crypto.Encrypt(data, recipients)
	if err != nil {
		return fmt.Errorf("vault: encrypt: %w", err)
	}

	if dst == "" {
		dst = src + defaultVaultExt
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return fmt.Errorf("vault: mkdir: %w", err)
	}

	if err := os.WriteFile(dst, ciphertext, 0o600); err != nil {
		return fmt.Errorf("vault: write vault file: %w", err)
	}

	return nil
}

// Unseal decrypts the vault file at src using the age identity (private key)
// and writes the plaintext .env to dst. If dst is empty the content is
// returned as a string without writing a file.
func (v *Vault) Unseal(src string, identity string, dst string) (string, error) {
	ciphertext, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("vault: read vault file: %w", err)
	}

	plaintext, err := crypto.Decrypt(ciphertext, identity)
	if err != nil {
		return "", fmt.Errorf("vault: decrypt: %w", err)
	}

	if dst != "" {
		if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
			return "", fmt.Errorf("vault: mkdir: %w", err)
		}
		if err := os.WriteFile(dst, plaintext, 0o600); err != nil {
			return "", fmt.Errorf("vault: write env file: %w", err)
		}
	}

	return string(plaintext), nil
}
