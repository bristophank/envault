// Package config manages envault's local configuration file (~/.envault/config.toml).
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	defaultConfigDir  = ".envault"
	defaultConfigFile = "config.json"
)

// Config holds envault's user-level configuration.
type Config struct {
	// DefaultIdentityFile is the path to the age private key used by default.
	DefaultIdentityFile string `json:"default_identity_file,omitempty"`
	// DefaultRecipientsFile is the path to the recipients list used by default.
	DefaultRecipientsFile string `json:"default_recipients_file,omitempty"`
}

// Manager handles reading and writing the config file.
type Manager struct {
	path string
}

// New returns a Manager rooted at the given directory.
// If dir is empty, it defaults to $HOME/.envault.
func New(dir string) (*Manager, error) {
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, defaultConfigDir)
	}
	return &Manager{path: filepath.Join(dir, defaultConfigFile)}, nil
}

// Load reads the config file from disk. If the file does not exist, an empty
// Config is returned without error.
func (m *Manager) Load() (*Config, error) {
	data, err := os.ReadFile(m.path)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the config to disk, creating parent directories as needed.
func (m *Manager) Save(cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o600)
}

// Path returns the absolute path to the config file.
func (m *Manager) Path() string {
	return m.path
}
