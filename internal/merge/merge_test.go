package merge_test

import (
	"testing"

	"github.com/envault/envault/internal/merge"
)

func TestMergeNoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"C": "3"}

	r := merge.Merge(base, incoming, merge.KeepBase)

	if r.Env["A"] != "1" || r.Env["B"] != "2" || r.Env["C"] != "3" {
		t.Fatalf("unexpected env: %v", r.Env)
	}
	if len(r.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", r.Conflicts)
	}
	if len(r.Added) != 1 || r.Added[0] != "C" {
		t.Fatalf("expected Added=[C], got %v", r.Added)
	}
}

func TestMergeConflictKeepBase(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}

	r := merge.Merge(base, incoming, merge.KeepBase)

	if r.Env["KEY"] != "old" {
		t.Fatalf("expected base value 'old', got %q", r.Env["KEY"])
	}
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	if r.Conflicts[0].BaseVal != "old" || r.Conflicts[0].IncomingVal != "new" {
		t.Fatalf("unexpected conflict values: %+v", r.Conflicts[0])
	}
}

func TestMergeConflictTakeIncoming(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}

	r := merge.Merge(base, incoming, merge.TakeIncoming)

	if r.Env["KEY"] != "new" {
		t.Fatalf("expected incoming value 'new', got %q", r.Env["KEY"])
	}
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict recorded, got %d", len(r.Conflicts))
	}
}

func TestMergeRemovedKeys(t *testing.T) {
	base := map[string]string{"KEEP": "1", "GONE": "2"}
	incoming := map[string]string{"KEEP": "1"}

	r := merge.Merge(base, incoming, merge.KeepBase)

	if len(r.Removed) != 1 || r.Removed[0] != "GONE" {
		t.Fatalf("expected Removed=[GONE], got %v", r.Removed)
	}
	if _, ok := r.Env["GONE"]; !ok {
		t.Fatal("removed key should still be present in merged env")
	}
}

func TestMergeEmptyInputs(t *testing.T) {
	r := merge.Merge(map[string]string{}, map[string]string{}, merge.KeepBase)
	if len(r.Env) != 0 || len(r.Conflicts) != 0 {
		t.Fatalf("expected empty result, got %+v", r)
	}
}
