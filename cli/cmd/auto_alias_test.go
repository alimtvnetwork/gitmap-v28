package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestAutoAlias_WordBoundarySplitting(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"anti-gravity-manager", "agm"},
		{"git_map_tool", "gmt"},
		{"my awesome project", "map"},
		{"icon-coding-guidelines", "icg"},
		{"common-linux-installer", "cli"},
		{"scripts-fixer", "sf"},
	}
	for _, tc := range cases {
		actual := GenerateAutoAlias(tc.input)
		if actual != tc.expected {
			t.Errorf("GenerateAutoAlias(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestAutoAlias_CompoundWords(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"gitmap", "gm"},
		{"github", "gh"},
		{"gitlab", "gl"},
		{"desktop", "dt"},
		{"pipeline", "pl"},
		{"workflow", "wf"},
		{"database", "db"},
		{"GitMap", "gm"},
	}
	for _, tc := range cases {
		actual := GenerateAutoAlias(tc.input)
		if actual != tc.expected {
			t.Errorf("GenerateAutoAlias(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestAutoAlias_WpPrefix(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"wp-git-log", "wp-gl"},
		{"wp-html-automated", "wp-ha"},
		{"wp_git_log", "wp-gl"},
		{"wp git log", "wp-gl"},
		{"wp-gitmap", "wp-gm"},
	}
	for _, tc := range cases {
		actual := GenerateAutoAlias(tc.input)
		if actual != tc.expected {
			t.Errorf("GenerateAutoAlias(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestAutoAlias_PresentationsPrefix(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"bsrm-presentation-hiltrax", "prep-bsrm"},
		{"bsrm-presentation-hiltrax-v4", "prep-bsrm"},
		{"presentations-repos/bsrm-presentation-hiltrax", "prep-bsrm"},
		{"presentation-docker", "prep-docker"},
		{"sales-presentation", "prep-sales"},
		{"presentation", "prep"},
	}
	for _, tc := range cases {
		actual := GenerateAutoAlias(tc.input)
		if actual != tc.expected {
			t.Errorf("GenerateAutoAlias(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestAutoAlias_GracefulFallbackSingleWord(t *testing.T) {
	cases := []string{
		"core",
		"api",
		"kernel",
		"utils",
		"",
	}
	for _, input := range cases {
		actual := GenerateAutoAlias(input)
		if actual != "" {
			t.Errorf("GenerateAutoAlias(%q) = %q; want empty string", input, actual)
		}
	}
}

func TestAutoAlias_CollisionAvoidance(t *testing.T) {
	taken := map[string]bool{
		"agm":  true,
		"agm2": true,
	}
	isTaken := func(s string) bool {
		return taken[s]
	}

	actual := ResolveUniqueAlias("agm", isTaken)
	if actual != "agm3" {
		t.Errorf("ResolveUniqueAlias(agm) = %q; want agm3", actual)
	}

	fresh := ResolveUniqueAlias("gm", isTaken)
	if fresh != "gm" {
		t.Errorf("ResolveUniqueAlias(gm) = %q; want gm", fresh)
	}
}

func TestAutoAlias_PopulateRepoAliases(t *testing.T) {
	db, err := store.OpenInMemory()
	if err != nil {
		t.Fatalf("OpenInMemory failed: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	repoID := seedTestRepo(t, db, "/path/to/anti-gravity-manager", "agm-slug", "anti-gravity-manager")
	assertAutoAliasPopulated(t, db, repoID, "agm")
}

func seedTestRepo(t *testing.T, db *store.DB, path, slug, name string) int64 {
	t.Helper()
	err := db.UpsertRepos([]model.ScanRecord{
		{AbsolutePath: path, Slug: slug, RepoName: name, Branch: "main", Notes: "note"},
	})
	if err != nil {
		t.Fatalf("UpsertRepos failed: %v", err)
	}

	repos, err := db.FindBySlug(slug)
	if err != nil || len(repos) == 0 {
		t.Fatalf("FindBySlug failed: %v", err)
	}

	return repos[0].ID
}

func assertAutoAliasPopulated(t *testing.T, db *store.DB, wantRepoID int64, wantAlias string) {
	t.Helper()
	count, appErr := PopulateRepoAliasesWithCount(db)
	if appErr != nil {
		t.Fatalf("PopulateRepoAliasesWithCount failed: %v", appErr)
	}
	if count != 1 {
		t.Errorf("count = %d; want 1", count)
	}

	alias, err := db.FindAliasByName(wantAlias)
	if err != nil {
		t.Fatalf("FindAliasByName(%s) failed: %v", wantAlias, err)
	}
	if alias.RepoID != wantRepoID {
		t.Errorf("alias.RepoID = %d; want %d", alias.RepoID, wantRepoID)
	}
}

func TestCmdInstall_EnsureTrackedRepoAliases(t *testing.T) {
	db, err := store.OpenInMemory()
	if err != nil {
		t.Fatalf("OpenInMemory failed: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	seedTestRepo(t, db, "/path/to/wp-git-log", "wpgl-slug", "wp-git-log")
	assertTrackedRepoAliasesEnsured(t, db)
}

func assertTrackedRepoAliasesEnsured(t *testing.T, db *store.DB) {
	t.Helper()
	count, err := cmdinstall.EnsureTrackedRepoAliases(db)
	if err != nil {
		t.Fatalf("EnsureTrackedRepoAliases failed: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d; want 1", count)
	}

	if !db.AliasExists("wp-gl") {
		t.Errorf("expected alias wp-gl to exist in database")
	}
}
