package keystore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultKeysDir  = ".envault"
	PrivateKeyFile  = "identity.txt"
	PublicKeysFile  = "recipients.txt"
)

// Store manages age keys on disk.
type Store struct {
	BaseDir string
}

// New returns a Store rooted at baseDir (e.g. the user's home dir).
func New(baseDir string) *Store {
	return &Store{BaseDir: filepath.Join(baseDir, DefaultKeysDir)}
}

// Init creates the key directory with restricted permissions.
func (s *Store) Init() error {
	return os.MkdirAll(s.BaseDir, 0700)
}

// SavePrivateKey writes the private key (identity) to disk.
func (s *Store) SavePrivateKey(privateKey string) error {
	path := filepath.Join(s.BaseDir, PrivateKeyFile)
	return os.WriteFile(path, []byte(privateKey+"\n"), 0600)
}

// LoadPrivateKey reads the private key from disk.
func (s *Store) LoadPrivateKey() (string, error) {
	path := filepath.Join(s.BaseDir, PrivateKeyFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("no identity found at %s: run 'envault init' first", path)
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// AddRecipient appends a public key to the recipients file.
func (s *Store) AddRecipient(publicKey string) error {
	path := filepath.Join(s.BaseDir, PublicKeysFile)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, strings.TrimSpace(publicKey))
	return err
}

// LoadRecipients returns all public keys from the recipients file.
func (s *Store) LoadRecipients() ([]string, error) {
	path := filepath.Join(s.BaseDir, PublicKeysFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var keys []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			keys = append(keys, line)
		}
	}
	return keys, nil
}

// PrivateKeyPath returns the absolute path to the identity file.
func (s *Store) PrivateKeyPath() string {
	return filepath.Join(s.BaseDir, PrivateKeyFile)
}
