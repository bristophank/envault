// Package search provides functionality to search for keys across
// sealed vault files and plaintext .env files.
package search

import (
	"fmt"
	"strings"
)

// Result represents a single match found during a search.
type Result struct {
	Key   string
	Value string
	File  string
}

// Searcher searches for keys or values within env maps.
type Searcher struct{}

// New returns a new Searcher.
func New() *Searcher {
	return &Searcher{}
}

// SearchKeys returns all entries whose key contains the given query (case-insensitive).
func (s *Searcher) SearchKeys(env map[string]string, file, query string) []Result {
	query = strings.ToLower(query)
	var results []Result
	for k, v := range env {
		if strings.Contains(strings.ToLower(k), query) {
			results = append(results, Result{Key: k, Value: v, File: file})
		}
	}
	return results
}

// SearchValues returns all entries whose value contains the given query (case-insensitive).
func (s *Searcher) SearchValues(env map[string]string, file, query string) []Result {
	query = strings.ToLower(query)
	var results []Result
	for k, v := range env {
		if strings.Contains(strings.ToLower(v), query) {
			results = append(results, Result{Key: k, Value: v, File: file})
		}
	}
	return results
}

// Format formats a slice of results for display.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no matches found"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.File != "" {
			sb.WriteString(fmt.Sprintf("%s: %s=%s\n", r.File, r.Key, r.Value))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s\n", r.Key, r.Value))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
