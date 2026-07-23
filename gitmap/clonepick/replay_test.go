package clonepick

// replay_test.go: covers the --replay branch of the cmd dispatcher
// without needing a live SQLite handle. The fakeLoader records every
// call so each test can assert routing (numeric -> ByID, non-numeric
// -> ByName) and bump-on-success semantics.

import (
	"errors"
	"strings"
	"testing"
)

type fakeLoader struct {
	planByID    Plan
	idByID      int64
	errByID     error
	planByName  Plan
	idByName    int64
	errByName   error
	touchCalls  []int64
	touchErr    error
	idArgSeen   int64
	nameArgSeen string
}

func (f *fakeLoader) LoadClonePickByID(id int64) (Plan, int64, error) {
	f.idArgSeen = id

	return f.planByID, f.idByID, f.errByID
}

func (f *fakeLoader) LoadClonePickByName(name string) (Plan, int64, error) {
	f.nameArgSeen = name

	return f.planByName, f.idByName, f.errByName
}

func (f *fakeLoader) TouchClonePickCreatedAt(id int64) error {
	f.touchCalls = append(f.touchCalls, id)

	return f.touchErr
}

func TestLoadFromDBNumericRefRoutesToByID(t *testing.T) {
	loader := &fakeLoader{
		planByID: Plan{Name: "docs", RepoUrl: "u"},
		idByID:   42,
	}
	plan, id, err := LoadFromDB(loader, "42")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != 42 || plan.Name != "docs" {
		t.Fatalf("got id=%d name=%q, want 42/docs", id, plan.Name)
	}
	if loader.idArgSeen != 42 {
		t.Fatalf("ByID called with %d, want 42", loader.idArgSeen)
	}
	if len(loader.nameArgSeen) > 0 {
		t.Fatal("ByName must not be called for numeric ref")
	}
}

func TestLoadFromDBNonNumericRefRoutesToByName(t *testing.T) {
	loader := &fakeLoader{
		planByName: Plan{Name: "release-bundle"},
		idByName:   7,
	}
	plan, id, err := LoadFromDB(loader, "release-bundle")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != 7 || plan.Name != "release-bundle" {
		t.Fatalf("got id=%d name=%q, want 7/release-bundle", id, plan.Name)
	}
	if loader.nameArgSeen != "release-bundle" {
		t.Fatalf("ByName called with %q", loader.nameArgSeen)
	}
}

func TestLoadFromDBMissingRefReturnsUserFacingMessage(t *testing.T) {
	loader := &fakeLoader{errByID: errors.New("sql: no rows")}
	_, _, err := LoadFromDB(loader, "999")
	if err == nil {
		t.Fatal("want error for missing id")
	}
	if !strings.Contains(err.Error(), "no saved selection") {
		t.Fatalf("err missing user-facing prefix: %v", err)
	}
}

func TestLoadFromDBNilLoaderRejected(t *testing.T) {
	_, _, err := LoadFromDB(nil, "1")
	if err == nil {
		t.Fatal("want error when loader is nil")
	}
}

func TestLoadFromDBEmptyRefRejected(t *testing.T) {
	_, _, err := LoadFromDB(&fakeLoader{}, "   ")
	if err == nil {
		t.Fatal("want error for blank ref")
	}
}

func TestTouchAfterReplaySkippedOnDryRun(t *testing.T) {
	loader := &fakeLoader{}
	if err := TouchAfterReplay(loader, 5, true); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(loader.touchCalls) != 0 {
		t.Fatalf("dry-run must not touch DB, got %v", loader.touchCalls)
	}
}

func TestTouchAfterReplayBumpsRow(t *testing.T) {
	loader := &fakeLoader{}
	if err := TouchAfterReplay(loader, 5, false); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(loader.touchCalls) != 1 || loader.touchCalls[0] != 5 {
		t.Fatalf("want touch(5), got %v", loader.touchCalls)
	}
}
