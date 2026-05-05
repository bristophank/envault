package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
}

// Parse reads .env formatted content and returns a slice of entries.
// It preserves comments and blank lines are skipped.
func Parse(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			entries = append(entries, Entry{Comment: line})
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		value = stripQuotes(value)
		entries = append(entries, Entry{Key: key, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return entries, nil
}

// Serialize writes entries back to .env format.
func Serialize(entries []Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Comment != "" {
			sb.WriteString(e.Comment + "\n")
		} else {
			sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, e.Value))
		}
	}
	return sb.String()
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
