// Package diff provides utilities for comparing two sets of environment
// variables and producing a human-readable summary of what changed.
package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Change represents a single modification between two env maps.
type Change struct {
	Key  string
	Kind Kind
	Old  string
	New  string
}

// Kind classifies the type of change.
type Kind int

const (
	Added Kind = iota
	Removed
	Modified
)

// String returns a human-readable representation of the change.
func (c Change) String() string {
	switch c.Kind {
	case Added:
		return fmt.Sprintf("+ %s=%s", c.Key, c.New)
	case Removed:
		return fmt.Sprintf("- %s=%s", c.Key, c.Old)
	case Modified:
		return fmt.Sprintf("~ %s: %s → %s", c.Key, c.Old, c.New)
	default:
		return ""
	}
}

// Compare returns an ordered list of Changes between the old and new env maps.
func Compare(oldEnv, newEnv map[string]string) []Change {
	var changes []Change

	for key, newVal := range newEnv {
		if oldVal, exists := oldEnv[key]; !exists {
			changes = append(changes, Change{Key: key, Kind: Added, New: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: key, Kind: Modified, Old: oldVal, New: newVal})
		}
	}

	for key, oldVal := range oldEnv {
		if _, exists := newEnv[key]; !exists {
			changes = append(changes, Change{Key: key, Kind: Removed, Old: oldVal})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})
	return changes
}

// Summary returns a multi-line string summarising all changes.
func Summary(changes []Change) string {
	if len(changes) == 0 {
		return "no changes"
	}
	lines := make([]string, len(changes))
	for i, c := range changes {
		lines[i] = c.String()
	}
	return strings.Join(lines, "\n")
}
