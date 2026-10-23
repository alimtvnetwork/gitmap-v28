package cmdscan

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func buildSampleRepos() []store.RepoRecord {
	return []store.RepoRecord{
		{RepoId: 101, Slug: "repo-alpha", RepoName: "Alpha", PrimaryAlias: "alpha"},
		{RepoId: 102, Slug: "repo-beta", RepoName: "Beta", PrimaryAlias: "beta"},
		{RepoId: 103, Slug: "repo-gamma", RepoName: "Gamma", PrimaryAlias: "gamma"},
		{RepoId: 104, Slug: "repo-delta", RepoName: "Delta", PrimaryAlias: "delta"},
		{RepoId: 105, Slug: "repo-epsilon", RepoName: "Epsilon", PrimaryAlias: "epsilon"},
	}
}

func TestScanDisableRepos_RangeParser(t *testing.T) {
	repos := buildSampleRepos()

	ids, err := ResolveTargetRepoIDs([]string{"1-3"}, repos)
	if err != nil {
		t.Fatalf("unexpected error resolving range 1-3: %v", err)
	}

	if len(ids) != 3 || ids[0] != 101 || ids[1] != 102 || ids[2] != 103 {
		t.Errorf("expected IDs [101, 102, 103], got %v", ids)
	}
}

func TestScanDisableRepos_CommaList(t *testing.T) {
	repos := buildSampleRepos()

	ids, err := ResolveTargetRepoIDs([]string{"2,4"}, repos)
	if err != nil {
		t.Fatalf("unexpected error resolving 2,4: %v", err)
	}

	if len(ids) != 2 || ids[0] != 102 || ids[1] != 104 {
		t.Errorf("expected IDs [102, 104], got %v", ids)
	}
}

func TestScanDisableRepos_SlugAndAlias(t *testing.T) {
	repos := buildSampleRepos()

	ids, err := ResolveTargetRepoIDs([]string{"gamma", "repo-epsilon"}, repos)
	if err != nil {
		t.Fatalf("unexpected error resolving slug/alias: %v", err)
	}

	if len(ids) != 2 || ids[0] != 103 || ids[1] != 105 {
		t.Errorf("expected IDs [103, 105], got %v", ids)
	}
}
