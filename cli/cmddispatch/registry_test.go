package cmddispatch

import (
	"strings"
	"testing"
)

// TestDispatchRegistryNoCollisions aggregates the real dispatch tables and
// fails on any name claimed by two or more owners. This is the regression net
// for alias hijacks like the `agm`→gitignore collision (task 247).
func TestDispatchRegistryNoCollisions(t *testing.T) {
	names, err := CollectNames()
	if err != nil {
		t.Fatalf("CollectNames failed: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("CollectNames returned zero names; expected the dispatch tables")
	}
	for _, c := range FindCollisions(names) {
		t.Errorf("dispatch collision: %q claimed by %s", c.Name, strings.Join(c.Owners, ", "))
	}
}

// TestFindCollisionsDetectsDuplicate proves the check catches the exact shape
// of the 247 regression: "agm" claimed both by a core alias table and by the
// agy subsystem dispatcher. Synthetic data only — production tables untouched.
func TestFindCollisionsDetectsDuplicate(t *testing.T) {
	names := []DispatchName{
		{Name: "agm", Owner: "cli/cmd/rootcore.go:78", Table: "coreBasicMaintenanceEntries"},
		{Name: "scan", Owner: "cli/cmd/rootcore.go:99", Table: "coreBasicOpEntries"},
		{Name: "agm", Owner: "cli/cmd/root.go:910", Table: "dispatchAgySubsystem"},
	}
	collisions := FindCollisions(names)
	if len(collisions) != 1 {
		t.Fatalf("expected 1 collision, got %d: %v", len(collisions), collisions)
	}
	got := collisions[0]
	if got.Name != "agm" {
		t.Errorf("expected collision on %q, got %q", "agm", got.Name)
	}
	if len(got.Owners) != 2 {
		t.Errorf("expected 2 owners, got %v", got.Owners)
	}
}

// TestFindCollisionsCleanTree is the inverse: disjoint names never collide.
func TestFindCollisionsCleanTree(t *testing.T) {
	names := []DispatchName{
		{Name: "scan", Owner: "cli/cmd/rootcore.go:99", Table: "core"},
		{Name: "agm", Owner: "cli/cmd/root.go:910", Table: "agy"},
	}
	if got := FindCollisions(names); len(got) != 0 {
		t.Errorf("expected no collisions, got %v", got)
	}
}
