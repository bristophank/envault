package env

import (
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	input := `# Database config
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	if entries[1].Key != "DB_HOST" || entries[1].Value != "localhost" {
		t.Errorf("unexpected entry: %+v", entries[1])
	}
}

func TestParseQuotedValues(t *testing.T) {
	input := `SECRET="my secret value"
TOKEN='another token'
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "my secret value" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
	if entries[1].Value != "another token" {
		t.Errorf("expected unquoted value, got %q", entries[1].Value)
	}
}

func TestParseInvalidLine(t *testing.T) {
	input := "INVALID_LINE_NO_EQUALS\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseEmptyKey(t *testing.T) {
	input := "=value\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSerializeRoundtrip(t *testing.T) {
	original := `# Comment
KEY1=value1
KEY2=value2
`
	entries, err := Parse(strings.NewReader(original))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	result := Serialize(entries)
	if result != original {
		t.Errorf("roundtrip mismatch:\nwant: %q\ngot:  %q", original, result)
	}
}

func TestParseSkipsBlankLines(t *testing.T) {
	input := "\nKEY=val\n\n"
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}
