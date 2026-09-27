package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestExtractJSONPayload(t *testing.T) {
	raw := "Some log before\n{\"total\":10,\"pulledCount\":10}\nSome log after"
	extracted := extractJSONPayload(raw)
	expected := "{\"total\":10,\"pulledCount\":10}"
	if extracted != expected {
		t.Fatalf("expected %q, got %q", expected, extracted)
	}
}

func TestFilterActiveFleetStates(t *testing.T) {
	items := []fleetPullRepoItem{
		{RepoName: "repo-1", Status: "up-to-date", Changes: "up-to-date"},
		{RepoName: "repo-2", Status: "updated", Changes: "+10/-2 (1)"},
		{RepoName: "repo-3", Status: "failed", Changes: "conflict"},
		{RepoName: "repo-4", Status: "synced", Changes: "synced"},
		{RepoName: "repo-5", Status: "dirty", Changes: "dirty"},
	}
	active := filterActiveFleetStates(items)
	if len(active) != 3 {
		t.Fatalf("expected 3 active items, got %d", len(active))
	}
	if active[0].RepoName != "repo-2" || active[1].RepoName != "repo-3" || active[2].RepoName != "repo-5" {
		t.Fatalf("unexpected active items: %+v", active)
	}
}

func TestHasFleetJSONFlag(t *testing.T) {
	if !hasFleetJSONFlag([]string{"pa", "--json"}) {
		t.Fatal("expected hasFleetJSONFlag=true for --json")
	}
	if !hasFleetJSONFlag([]string{"pa", "-json"}) {
		t.Fatal("expected hasFleetJSONFlag=true for -json")
	}
	if hasFleetJSONFlag([]string{"pa", "--status"}) {
		t.Fatal("expected hasFleetJSONFlag=false without json")
	}
}

func TestResolveTargetNodeOS(t *testing.T) {
	conn := db.SSHConnection{OS: "linux"}
	osType := resolveTargetNodeOS(nil, conn)
	if osType != "linux" {
		t.Fatalf("expected linux, got %s", osType)
	}
}
