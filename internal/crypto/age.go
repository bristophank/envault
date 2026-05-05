package crypto

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// Encrypt encrypts plaintext using the provided age public key recipients.
// Returns PEM-armored ciphertext.
func Encrypt(plaintext []byte, publicKeys []string) ([]byte, error) {
	var recipients []age.Recipient
	for _, pub := range publicKeys {
		r, err := age.ParseX25519Recipient(strings.TrimSpace(pub))
		if err != nil {
			return nil, fmt.Errorf("invalid public key %q: %w", pub, err)
		}
		recipients = append(recipients, r)
	}

	if len(recipients) == 0 {
		return nil, fmt.Errorf("at least one recipient public key is required")
	}

	var buf bytes.Buffer
	armorWriter := armor.NewWriter(&buf)

	w, err := age.Encrypt(armorWriter, recipients...)
	if err != nil {
		return nil, fmt.Errorf("failed to create age encryptor: %w", err)
	}

	if _, err := w.Write(plaintext); err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize encryption: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close armor writer: %w", err)
	}

	return buf.Bytes(), nil
}

// Decrypt decrypts PEM-armored ciphertext using the provided age private key.
func Decrypt(ciphertext []byte, privateKey string) ([]byte, error) {
	id, err := age.ParseX25519Identity(strings.TrimSpace(privateKey))
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	armorReader := armor.NewReader(bytes.NewReader(ciphertext))
	r, err := age.Decrypt(armorReader, id)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read decrypted data: %w", err)
	}

	return plaintext, nil
}

// GenerateKeyPair generates a new age X25519 key pair.
// Returns (publicKey, privateKey, error).
func GenerateKeyPair() (string, string, error) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key pair: %w", err)
	}
	return id.Recipient().String(), id.String(), nil
}
