// Package recipients manages the list of age public keys (recipients)
// that are authorized to decrypt a vault.
package recipients

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// File is the name of the recipients file stored alongside a vault.
const File = ".envault-recipients"

// List holds a deduplicated ordered set of age public keys.
type List struct {
	keys []string
}

// New returns an empty List.
func New() *List {
	return &List{}
}

// Add appends a recipient public key if it is not already present.
// Returns an error if the key is empty or obviously malformed.
func (l *List) Add(pubkey string) error {
	pubkey = strings.TrimSpace(pubkey)
	if pubkey == "" {
		return fmt.Errorf("recipients: public key must not be empty")
	}
	if !strings.HasPrefix(pubkey, "age1") {
		return fmt.Errorf("recipients: invalid age public key %q (must start with age1)", pubkey)
	}
	for _, k := range l.keys {
		if k == pubkey {
			return nil // already present, idempotent
		}
	}
	l.keys = append(l.keys, pubkey)
	return nil
}

// Remove deletes a recipient by public key. Returns an error if not found.
func (l *List) Remove(pubkey string) error {
	pubkey = strings.TrimSpace(pubkey)
	for i, k := range l.keys {
		if k == pubkey {
			l.keys = append(l.keys[:i], l.keys[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("recipients: key %q not found", pubkey)
}

// Keys returns a copy of the current recipient public keys.
func (l *List) Keys() []string {
	out := make([]string, len(l.keys))
	copy(out, l.keys)
	return out
}

// Len returns the number of recipients.
func (l *List) Len() int { return len(l.keys) }

// LoadFile reads a recipients file from disk (one public key per line).
// Lines starting with '#' and blank lines are ignored.
func LoadFile(path string) (*List, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("recipients: open %s: %w", path, err)
	}
	defer f.Close()

	l := New()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := l.Add(line); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("recipients: scan %s: %w", path, err)
	}
	return l, nil
}

// SaveFile writes the recipient list to disk, one key per line.
func SaveFile(path string, l *List) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("recipients: create %s: %w", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, k := range l.keys {
		if _, err := fmt.Fprintln(w, k); err != nil {
			return err
		}
	}
	return w.Flush()
}
