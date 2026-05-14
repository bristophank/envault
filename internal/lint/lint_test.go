package lint_test

import (
	"testing"

	"github.com/envault/envault/internal/lint"
)

func TestLintCleanFile(t *testing.T) {
	content := "DB_HOST=localhost\nDB_PORT=5432\nDB_NAME=myapp\n"
	result, err := lint.Lint(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestLintEmptyValue(t *testing.T) {
	content := "API_KEY=\nSECRET=abc\n"
	result, err := lint.Lint(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Key != "API_KEY" {
		t.Errorf("expected issue on API_KEY, got %s", result.Issues[0].Key)
	}
	if result.Issues[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestLintDuplicateKey(t *testing.T) {
	content := "FOO=bar\nBAZ=qux\nFOO=other\n"
	result, err := lint.Lint(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasErrors() {
		t.Error("expected errors due to duplicate key")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "FOO" && issue.Severity == lint.SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected error issue for duplicate FOO key")
	}
}

func TestLintPlaceholderValue(t *testing.T) {
	content := "DB_PASS=changeme\nAPP_TOKEN=FIXME\n"
	result, err := lint.Lint(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(result.Issues))
	}
	for _, issue := range result.Issues {
		if issue.Severity != lint.SeverityWarning {
			t.Errorf("expected warning for placeholder, got %s", issue.Severity)
		}
	}
}

func TestLintNoErrors(t *testing.T) {
	content := "KEY=value\n"
	result, err := lint.Lint(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasErrors() {
		t.Error("expected no errors")
	}
}

func TestIssueString(t *testing.T) {
	issue := lint.Issue{
		Line:     3,
		Key:      "MY_KEY",
		Message:  "empty value",
		Severity: lint.SeverityWarning,
	}
	s := issue.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
