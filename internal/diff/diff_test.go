package diff_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envault/internal/diff"
)

func TestCompareAdded(t *testing.T) {
	old := map[string]string{"A": "1"}
	new := map[string]string{"A": "1", "B": "2"}

	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Added || changes[0].Key != "B" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompareRemoved(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1"}

	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Removed || changes[0].Key != "B" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompareModified(t *testing.T) {
	old := map[string]string{"A": "old"}
	new := map[string]string{"A": "new"}

	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != diff.Modified || c.Old != "old" || c.New != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestCompareNoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	changes := diff.Compare(env, env)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestSummaryNoChanges(t *testing.T) {
	s := diff.Summary(nil)
	if s != "no changes" {
		t.Errorf("expected 'no changes', got %q", s)
	}
}

func TestSummaryContainsKeys(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "baz", "NEW": "val"}

	changes := diff.Compare(old, new)
	s := diff.Summary(changes)

	if !strings.Contains(s, "FOO") {
		t.Errorf("summary missing FOO: %s", s)
	}
	if !strings.Contains(s, "NEW") {
		t.Errorf("summary missing NEW: %s", s)
	}
}

func TestChangesSortedByKey(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"Z": "1", "A": "2", "M": "3"}

	changes := diff.Compare(old, new)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes")
	}
	if changes[0].Key != "A" || changes[1].Key != "M" || changes[2].Key != "Z" {
		t.Errorf("changes not sorted: %v", changes)
	}
}
