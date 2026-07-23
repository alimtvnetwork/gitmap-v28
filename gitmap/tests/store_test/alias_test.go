// Package store_test — integration tests for Alias CRUD operations.
package store_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// seedRepo inserts a test repo and returns its auto-generated ID.
func seedRepo(t *testing.T, db *store.DB, slug, repoName, absPath string) int64 {
	t.Helper()

	records := []model.ScanRecord{{
		Slug:         slug,
		RepoName:     repoName,
		AbsolutePath: absPath,
	}}

	if err := db.UpsertRepos(records); err != nil {
		t.Fatalf("failed to seed repo %s: %v", slug, err)
	}

	// Look up the auto-generated ID.
	repos, err := db.FindBySlug(slug)
	if err != nil || len(repos) == 0 {
		t.Fatalf("failed to find seeded repo %s", slug)
	}

	return repos[0].ID
}

// TestCreateAlias_Success verifies a new alias is created and returned.
func TestCreateAlias_Success(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repoID := seedRepo(t, db, "github/user/api", "api", "/home/user/repos/api")

	alias, err := db.CreateAlias("api", repoID)
	if err != nil {
		t.Fatalf("CreateAlias failed: %v", err)
	}

	if alias.Alias != "api" {
		t.Errorf("expected alias=api, got %q", alias.Alias)
	}

	if alias.RepoID != repoID {
		t.Errorf("expected repoID=%d, got %d", repoID, alias.RepoID)
	}

	if alias.ID == 0 {
		t.Error("expected non-zero alias ID")
	}

	if alias.CreatedAt == "" {
		t.Error("expected CreatedAt to be populated")
	}
}

// TestCreateAlias_DuplicateFails verifies duplicate alias names are rejected.
func TestCreateAlias_DuplicateFails(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repoID := seedRepo(t, db, "github/user/api", "api", "/home/user/repos/api")

	_, err := db.CreateAlias("api", repoID)
	if err != nil {
		t.Fatalf("first CreateAlias failed: %v", err)
	}

	_, err = db.CreateAlias("api", repoID)
	if err == nil {
		t.Error("expected error on duplicate alias, got nil")
	}
}

// TestResolveAlias_Success verifies alias resolves to correct repo path and slug.
func TestResolveAlias_Success(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repoID := seedRepo(t, db, "github/user/web", "web", "/home/user/repos/web")
	t.Logf("seeded repo ID: %d", repoID)

	alias, err := db.CreateAlias("web", repoID)
	if err != nil {
		t.Fatalf("CreateAlias failed: %v", err)
	}
	t.Logf("created alias ID: %d, RepoID: %d", alias.ID, alias.RepoID)

	// Verify alias exists via direct lookup.
	if !db.AliasExists("web") {
		t.Fatal("alias 'web' should exist after creation")
	}

	resolved, err := db.ResolveAlias("web")
	if err != nil {
		t.Fatalf("ResolveAlias failed: %v", err)
	}

	if resolved.Alias.Alias != "web" {
		t.Errorf("expected alias=web, got %q", resolved.Alias.Alias)
	}

	if resolved.AbsolutePath != "/home/user/repos/web" {
		t.Errorf("expected path=/home/user/repos/web, got %q", resolved.AbsolutePath)
	}

	if resolved.Slug != "github/user/web" {
		t.Errorf("expected slug=github/user/web, got %q", resolved.Slug)
	}
}

// TestResolveAlias_NotFound verifies error for non-existent alias.
func TestResolveAlias_NotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	_, err := db.ResolveAlias("nonexistent")
	if err == nil {
		t.Error("expected error for missing alias, got nil")
	}
}

// TestListAliasesWithRepo_Empty verifies empty list when no aliases exist.
func TestListAliasesWithRepo_Empty(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	aliases, err := db.ListAliasesWithRepo()
	if err != nil {
		t.Fatalf("ListAliasesWithRepo failed: %v", err)
	}

	if len(aliases) != 0 {
		t.Errorf("expected 0 aliases, got %d", len(aliases))
	}
}

// TestListAliasesWithRepo_Multiple verifies all aliases are returned with repo details.
func TestListAliasesWithRepo_Multiple(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	id1 := seedRepo(t, db, "github/user/api", "api", "/repos/api")
	id2 := seedRepo(t, db, "github/user/web", "web", "/repos/web")
	t.Logf("seeded repo IDs: %d, %d", id1, id2)

	if _, err := db.CreateAlias("api", id1); err != nil {
		t.Fatalf("CreateAlias(api) failed: %v", err)
	}
	if _, err := db.CreateAlias("web", id2); err != nil {
		t.Fatalf("CreateAlias(web) failed: %v", err)
	}

	aliases, err := db.ListAliasesWithRepo()
	if err != nil {
		t.Fatalf("ListAliasesWithRepo failed: %v", err)
	}

	if len(aliases) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(aliases))
	}

	// Results are ordered by alias name.
	if aliases[0].Alias.Alias != "api" {
		t.Errorf("expected first alias=api, got %q", aliases[0].Alias.Alias)
	}

	if aliases[1].Alias.Alias != "web" {
		t.Errorf("expected second alias=web, got %q", aliases[1].Alias.Alias)
	}

	if aliases[0].AbsolutePath != "/repos/api" {
		t.Errorf("expected path=/repos/api, got %q", aliases[0].AbsolutePath)
	}
}

// TestDeleteAlias_Success verifies alias is removed.
func TestDeleteAlias_Success(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repoID := seedRepo(t, db, "github/user/api", "api", "/repos/api")

	db.CreateAlias("api", repoID)

	err := db.DeleteAlias("api")
	if err != nil {
		t.Fatalf("DeleteAlias failed: %v", err)
	}

	if db.AliasExists("api") {
		t.Error("alias should not exist after deletion")
	}
}

// TestDeleteAlias_NotFound verifies error when deleting non-existent alias.
func TestDeleteAlias_NotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	err := db.DeleteAlias("ghost")
	if err == nil {
		t.Error("expected error deleting non-existent alias, got nil")
	}
}

// TestAliasExists_TrueAndFalse verifies existence check.
func TestAliasExists_TrueAndFalse(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repoID := seedRepo(t, db, "github/user/infra", "infra", "/repos/infra")

	if db.AliasExists("infra") {
		t.Error("alias should not exist before creation")
	}

	db.CreateAlias("infra", repoID)

	if !db.AliasExists("infra") {
		t.Error("alias should exist after creation")
	}
}

// TestUpdateAlias_ReassignsRepo verifies alias can be reassigned to a different repo.
func TestUpdateAlias_ReassignsRepo(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	id1 := seedRepo(t, db, "github/user/old-api", "old-api", "/repos/old-api")
	id2 := seedRepo(t, db, "github/user/new-api", "new-api", "/repos/new-api")
	t.Logf("seeded repo IDs: old=%d, new=%d", id1, id2)

	if _, err := db.CreateAlias("api", id1); err != nil {
		t.Fatalf("CreateAlias failed: %v", err)
	}

	err := db.UpdateAlias("api", id2)
	if err != nil {
		t.Fatalf("UpdateAlias failed: %v", err)
	}

	// Verify the alias still exists after update.
	if !db.AliasExists("api") {
		t.Fatal("alias 'api' should still exist after update")
	}

	resolved, err := db.ResolveAlias("api")
	if err != nil {
		t.Fatalf("ResolveAlias after update failed: %v", err)
	}

	if resolved.AbsolutePath != "/repos/new-api" {
		t.Errorf("expected path=/repos/new-api, got %q", resolved.AbsolutePath)
	}
}

// TestListUnaliasedRepos_ReturnsOnlyUnaliased verifies filtering.
func TestListUnaliasedRepos_ReturnsOnlyUnaliased(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	id1 := seedRepo(t, db, "github/user/api", "api", "/repos/api")
	seedRepo(t, db, "github/user/web", "web", "/repos/web")

	db.CreateAlias("api", id1)

	unaliased, err := db.ListUnaliasedRepos()
	if err != nil {
		t.Fatalf("ListUnaliasedRepos failed: %v", err)
	}

	if len(unaliased) != 1 {
		t.Fatalf("expected 1 unaliased repo, got %d", len(unaliased))
	}

	if unaliased[0].RepoName != "web" {
		t.Errorf("expected unaliased repo=web, got %q", unaliased[0].RepoName)
	}
}

// TestFindAliasByRepoID_Success verifies lookup by repo ID.
func TestFindAliasByRepoID_Success(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repoID := seedRepo(t, db, "github/user/svc", "svc", "/repos/svc")

	db.CreateAlias("svc", repoID)

	found, err := db.FindAliasByRepoID(repoID)
	if err != nil {
		t.Fatalf("FindAliasByRepoID failed: %v", err)
	}

	if found.Alias != "svc" {
		t.Errorf("expected alias=svc, got %q", found.Alias)
	}
}
