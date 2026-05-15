// Package history tracks seal/unseal operations on vault files,
// allowing users to view a timeline of changes per environment file.
package history

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single recorded operation.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"` // "seal" | "unseal" | "rotate"
	File      string    `json:"file"`
	User      string    `json:"user"`
	Note      string    `json:"note,omitempty"`
}

// History manages a per-file operation log stored in the keystore directory.
type History struct {
	dir string
}

// New returns a History that stores logs under dir.
func New(dir string) *History {
	return &History{dir: dir}
}

func (h *History) logPath(file string) string {
	base := filepath.Base(file) + ".history.jsonl"
	return filepath.Join(h.dir, base)
}

// Record appends an entry for the given file.
func (h *History) Record(op, file, user, note string) error {
	if err := os.MkdirAll(h.dir, 0o700); err != nil {
		return err
	}
	e := Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		File:      file,
		User:      user,
		Note:      note,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(h.logPath(file), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(data, '\n'))
	return err
}

// List returns all recorded entries for the given file, oldest first.
func (h *History) List(file string) ([]Entry, error) {
	data, err := os.ReadFile(h.logPath(file))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
