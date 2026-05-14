// Package snapshot provides functionality to capture and restore
// point-in-time snapshots of sealed vault files.
package snapshot

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Manager handles saving and listing snapshots of sealed env files.
type Manager struct {
	dir string
}

// New returns a Manager that stores snapshots under dir.
func New(dir string) *Manager {
	return &Manager{dir: dir}
}

// Save copies src into the snapshot directory, tagging it with the
// current timestamp. It returns the path of the created snapshot.
func (m *Manager) Save(src string) (string, error) {
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return "", fmt.Errorf("snapshot: create dir: %w", err)
	}

	base := filepath.Base(src)
	tag := time.Now().UTC().Format("20060102T150405Z")
	dst := filepath.Join(m.dir, fmt.Sprintf("%s.%s", base, tag))

	if err := copyFile(src, dst); err != nil {
		return "", fmt.Errorf("snapshot: save: %w", err)
	}
	return dst, nil
}

// List returns all snapshot paths stored in the directory, sorted
// oldest-first (lexicographic on the timestamp suffix).
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return "", fmt.Errorf("snapshot: list: %w", err)
	}

	var paths []string
	for _, e := range entries {
		if !e.IsDir() {
			paths = append(paths, filepath.Join(m.dir, e.Name()))
		}
	}
	return paths, nil
}

// Restore copies the snapshot at snapshotPath back to dst, overwriting it.
func (m *Manager) Restore(snapshotPath, dst string) error {
	if err := copyFile(snapshotPath, dst); err != nil {
		return fmt.Errorf("snapshot: restore: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
