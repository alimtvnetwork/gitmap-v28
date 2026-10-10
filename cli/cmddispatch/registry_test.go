package cmddispatch

import (
	"strings"
	"testing"
)

// knownLegacyCollisions records grandfathered dispatch overlaps from legacy dispatch tables.
// Any NEW collision outside this baseline will immediately fail the test.
var knownLegacyCollisions = map[string]bool{
	"aef": true, "agy-errors-fix": true, "al": true, "c": true, "cat": true,
	"cron": true, "crontab": true, "dp": true, "exec": true, "export-all": true,
	"fc": true, "ff": true, "fix-agy": true, "fix-pipeline": true, "folder": true,
	"history": true, "ignore": true, "import-all": true, "ip": true,
	"join-common": true, "lc-fix": true, "lcf": true, "lcr": true, "lower": true,
	"lower-case-fix": true, "lower-case-readme": true, "lowercase": true,
	"lowercase-fix": true, "lowercase-readme": true, "migrate": true, "o": true,
	"open": true, "os": true, "pf": true, "pipeline-fix": true, "pr": true,
	"py": true, "python": true, "readme-lower": true, "readme-lowercase": true,
	"rec": true, "recreate": true, "replace": true, "rp": true, "sb": true,
	"sc": true, "schedule": true, "shutdown-until": true, "shutdown-until-green": true,
	"sj": true, "sjc": true, "ssh-join": true, "ssh-join-c": true,
	"ssh-join-common": true, "ssh-joined": true, "ssh-joiner": true,
	"sug": true, "vscode": true,
}

// TestDispatchRegistryNoCollisions aggregates the real dispatch tables and
// fails on any new name claimed by two or more owners. This is the regression net
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
		if knownLegacyCollisions[c.Name] {
			continue
		}
		t.Errorf("new dispatch collision: %q claimed by %s", c.Name, strings.Join(c.Owners, ", "))
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
