// Package audit provides a simple append-only audit log for envault
// operations such as seal, unseal, rotate, and key generation.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Target    string    `json:"target,omitempty"`
	User      string    `json:"user,omitempty"`
	Details   string    `json:"details,omitempty"`
}

// Logger writes audit entries to a newline-delimited JSON file.
type Logger struct {
	path string
}

// New returns a Logger that appends to the file at path.
// Parent directories are created if they do not exist.
func New(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("audit: create directory: %w", err)
	}
	return &Logger{path: path}, nil
}

// Log appends an entry to the audit log file.
func (l *Logger) Log(op, target, details string) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Target:    target,
		User:      currentUser(),
		Details:   details,
	}

	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// ReadAll reads and returns all entries from the audit log.
func (l *Logger) ReadAll() ([]Entry, error) {
	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func currentUser() string {
	if u := os.Getenv("USER"); u != "" {
		return u
	}
	return os.Getenv("USERNAME")
}
