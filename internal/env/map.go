package env

import (
	"fmt"
	"sort"
	"strings"
)

// ToMap converts a slice of entries into a key-value map.
// Comment-only entries are ignored.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.Key != "" {
			m[e.Key] = e.Value
		}
	}
	return m
}

// FromMap converts a key-value map into a sorted slice of entries.
func FromMap(m map[string]string) []Entry {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(m))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Value: m[k]})
	}
	return entries
}

// Merge combines base entries with overrides. Keys in overrides
// replace matching keys in base; new keys are appended.
func Merge(base, overrides []Entry) []Entry {
	overrideMap := ToMap(overrides)
	result := make([]Entry, 0, len(base)+len(overrides))
	seen := make(map[string]bool)

	for _, e := range base {
		if e.Key == "" {
			result = append(result, e)
			continue
		}
		if v, ok := overrideMap[e.Key]; ok {
			result = append(result, Entry{Key: e.Key, Value: v})
		} else {
			result = append(result, e)
		}
		seen[e.Key] = true
	}

	for _, e := range overrides {
		if e.Key != "" && !seen[e.Key] {
			result = append(result, e)
		}
	}
	return result
}

// ParseString is a convenience wrapper around Parse for string input.
func ParseString(s string) ([]Entry, error) {
	return Parse(strings.NewReader(s))
}

// MustParseString parses a string and panics on error. Useful in tests.
func MustParseString(s string) []Entry {
	entries, err := ParseString(s)
	if err != nil {
		panic(fmt.Sprintf("env.MustParseString: %v", err))
	}
	return entries
}
