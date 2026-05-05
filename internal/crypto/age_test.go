package crypto_test

import (
	"bytes"
	"testing"

	"github.com/user/envault/internal/crypto"
)

func TestGenerateKeyPair(t *testing.T) {
	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}
	if pub == "" {
		t.Error("expected non-empty public key")
	}
	if priv == "" {
		t.Error("expected non-empty private key")
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	pub, priv, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v", err)
	}

	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123\n")

	ciphertext, err := crypto.Encrypt(plaintext, []string{pub})
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	if len(ciphertext) == 0 {
		t.Fatal("expected non-empty ciphertext")
	}
	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := crypto.Decrypt(ciphertext, priv)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptMultipleRecipients(t *testing.T) {
	pub1, priv1, _ := crypto.GenerateKeyPair()
	pub2, priv2, _ := crypto.GenerateKeyPair()

	plaintext := []byte("SECRET=shared")
	ciphertext, err := crypto.Encrypt(plaintext, []string{pub1, pub2})
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	for _, priv := range []string{priv1, priv2} {
		decrypted, err := crypto.Decrypt(ciphertext, priv)
		if err != nil {
			t.Errorf("Decrypt() with recipient key error = %v", err)
			continue
		}
		if !bytes.Equal(decrypted, plaintext) {
			t.Errorf("Decrypt() = %q, want %q", decrypted, plaintext)
		}
	}
}

func TestEncryptNoRecipients(t *testing.T) {
	_, err := crypto.Encrypt([]byte("data"), []string{})
	if err == nil {
		t.Error("expected error when no recipients provided")
	}
}

func TestDecryptInvalidKey(t *testing.T) {
	pub, _, _ := crypto.GenerateKeyPair()
	ciphertext, _ := crypto.Encrypt([]byte("secret"), []string{pub})

	_, wrongPriv, _ := crypto.GenerateKeyPair()
	_, err := crypto.Decrypt(ciphertext, wrongPriv)
	if err == nil {
		t.Error("expected error when decrypting with wrong private key")
	}
}
