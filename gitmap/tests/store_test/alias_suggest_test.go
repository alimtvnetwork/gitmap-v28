// Package store_test — tests for alias suggestion and conflict detection.
package store_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// --- Suggestion: uses RepoName as proposed alias ---

func TestSuggestAlias_UnaliasedRepos(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	seedRepo(t, db, "github/user/api", "api", "/repos/api")
	seedRepo(t, db, "github/user/web", "web", "/repos/web")

	unaliased, err := db.ListUnaliasedRepos()
	if err != nil {
		t.Fatalf("ListUnaliasedRepos failed: %v", err)
	}

	if len(unaliased) != 2 {
		t.Fatalf("expected 2 unaliased, got %d", len(unaliased))
	}

	// Simulate suggest: propose RepoName as alias
	for _, r := range unaliased {
		if db.AliasExists(r.RepoName) {
			t.Errorf("alias %q should not exist before suggestion", r.RepoName)
		}
	}
}

// TestSuggestAlias_SkipsExistingAlias verifies conflict detection.
func TestSuggestAlias_SkipsExistingAlias(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	id1 := seedRepo(t, db, "github/user/api", "api", "/repos/api")
	seedRepo(t, db, "github/org/api", "api", "/repos/org-api")

	// First repo takes the "api" alias
	_, err := db.CreateAlias("api", id1)
	if err != nil {
		t.Fatalf("CreateAlias failed: %v", err)
	}

	// Second repo with same RepoName should be skipped
	unaliased, err := db.ListUnaliasedRepos()
	if err != nil {
		t.Fatalf("ListUnaliasedRepos failed: %v", err)
	}

	for _, r := range unaliased {
		suggestion := r.RepoName
		if db.AliasExists(suggestion) {
			// This is the conflict — suggestion should be skipped
			continue
		}
		t.Errorf("expected conflict for %q but AliasExists returned false", suggestion)
	}
}

// TestSuggestAlias_AutoApplyCreatesAll simulates --apply flag.
func TestSuggestAlias_AutoApplyCreatesAll(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	seedRepo(t, db, "github/user/alpha", "alpha", "/repos/alpha")
	seedRepo(t, db, "github/user/beta", "beta", "/repos/beta")
	seedRepo(t, db, "github/user/gamma", "gamma", "/repos/gamma")

	unaliased, err := db.ListUnaliasedRepos()
	if err != nil {
		t.Fatalf("ListUnaliasedRepos failed: %v", err)
	}

	created := 0
	for _, r := range unaliased {
		if db.AliasExists(r.RepoName) {
			continue
		}
		_, err := db.CreateAlias(r.RepoName, r.ID)
		if err != nil {
			t.Fatalf("CreateAlias(%s) failed: %v", r.RepoName, err)
		}
		created++
	}

	if created != 3 {
		t.Errorf("expected 3 created, got %d", created)
	}

	// All repos should now be aliased
	remaining, _ := db.ListUnaliasedRepos()
	if len(remaining) != 0 {
		t.Errorf("expected 0 unaliased after apply, got %d", len(remaining))
	}
}

// TestSuggestAlias_PartialConflict verifies mixed scenario.
func TestSuggestAlias_PartialConflict(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	id1 := seedRepo(t, db, "github/user/api", "api", "/repos/api")
	seedRepo(t, db, "github/user/web", "web", "/repos/web")
	seedRepo(t, db, "github/org/api", "api", "/repos/org-api")

	// "api" alias already taken by first repo
	db.CreateAlias("api", id1)

	unaliased, err := db.ListUnaliasedRepos()
	if err != nil {
		t.Fatalf("ListUnaliasedRepos failed: %v", err)
	}

	created := 0
	skipped := 0
	for _, r := range unaliased {
		if db.AliasExists(r.RepoName) {
			skipped++
			continue
		}
		db.CreateAlias(r.RepoName, r.ID)
		created++
	}

	if created != 1 {
		t.Errorf("expected 1 created (web), got %d", created)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped (api conflict), got %d", skipped)
	}
}

// TestSuggestAlias_AfterDeleteReopens verifies deleted alias frees the name.
func TestSuggestAlias_AfterDeleteReopens(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	id1 := seedRepo(t, db, "github/user/api", "api", "/repos/api")
	id2 := seedRepo(t, db, "github/org/api", "api", "/repos/org-api")

	if _, err := db.CreateAlias("api", id1); err != nil {
		t.Fatalf("CreateAlias(api, id1) failed: %v", err)
	}

	// "api" conflicts for second repo
	if !db.AliasExists("api") {
		t.Fatal("expected api alias to exist")
	}

	// Delete frees the name
	if err := db.DeleteAlias("api"); err != nil {
		t.Fatalf("DeleteAlias failed: %v", err)
	}

	if db.AliasExists("api") {
		t.Fatal("expected api alias to be gone after delete")
	}

	// Now second repo can claim it
	_, err := db.CreateAlias("api", id2)
	if err != nil {
		t.Fatalf("CreateAlias after delete failed: %v", err)
	}

	resolved, err := db.ResolveAlias("api")
	if err != nil {
		t.Fatalf("ResolveAlias after re-create failed: %v", err)
	}
	if resolved.AbsolutePath != "/repos/org-api" {
		t.Errorf("expected /repos/org-api, got %q", resolved.AbsolutePath)
	}
}

// TestSuggestAlias_EmptyDB verifies no suggestions on empty database.
func TestSuggestAlias_EmptyDB(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	unaliased, err := db.ListUnaliasedRepos()
	if err != nil {
		t.Fatalf("ListUnaliasedRepos failed: %v", err)
	}

	if len(unaliased) != 0 {
		t.Errorf("expected 0 unaliased on empty DB, got %d", len(unaliased))
	}
}

// helperSuggest simulates the suggest loop and returns created/skipped counts.
func helperSuggest(db *store.DB, repos []store.UnaliasedRepo) (created, skipped int) {
	for _, r := range repos {
		if db.AliasExists(r.RepoName) {
			skipped++
			continue
		}
		_, err := db.CreateAlias(r.RepoName, r.ID)
		if err != nil {
			skipped++
			continue
		}
		created++
	}

	return
}

// TestSuggestAlias_HelperMultiplePasses verifies idempotency.
func TestSuggestAlias_HelperMultiplePasses(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	seedRepo(t, db, "github/user/svc", "svc", "/repos/svc")

	unaliased, _ := db.ListUnaliasedRepos()
	c1, s1 := helperSuggest(db, unaliased)

	if c1 != 1 || s1 != 0 {
		t.Errorf("pass 1: expected created=1 skipped=0, got %d/%d", c1, s1)
	}

	// Second pass: all already aliased
	unaliased2, _ := db.ListUnaliasedRepos()
	if len(unaliased2) != 0 {
		t.Errorf("pass 2: expected 0 unaliased, got %d", len(unaliased2))
	}
}
