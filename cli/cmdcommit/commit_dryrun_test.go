package cmdcommit

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// assertNilCommitErr accepts both a genuine nil error and the typed-nil
// *apperror.AppError that executeCommitPush returns on success (the CLI
// dispatch unwraps it the same way in cli/cmd/root.go).
func assertNilCommitErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	if appErr, ok := err.(*apperror.AppError); ok && appErr == nil {
		return
	}
	t.Fatalf("expected nil commit error, got: %v", err)
}

func writeFixtureFile(t *testing.T, dir, name, content string) error {
	t.Helper()

	return os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
}

func appendFixtureFile(t *testing.T, dir, name, content string) error {
	t.Helper()
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)

	return err
}

// Fixture repo state: HEAD, porcelain status, and staged diff. A dry-run
// must leave all three byte-identical, and must not move a remote.
type fixtureState struct {
	head      string
	porcelain string
	cached    string
	bareHead  string
}

func runGitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}

	return strings.TrimSpace(string(out))
}

func newDryRunFixture(t *testing.T, withRemote bool) (string, string) {
	t.Helper()
	dir := t.TempDir()
	runGitIn(t, dir, "init", "-q", "-b", "main")
	runGitIn(t, dir, "config", "user.email", "dryrun@example.com")
	runGitIn(t, dir, "config", "user.name", "Dry Run")
	runGitIn(t, dir, "config", "commit.gpgsign", "false")
	if err := writeFixtureFile(t, dir, "a.txt", "alpha\n"); err != nil {
		t.Fatal(err)
	}
	if err := writeFixtureFile(t, dir, "c.txt", "charlie\n"); err != nil {
		t.Fatal(err)
	}
	runGitIn(t, dir, "add", "-A")
	runGitIn(t, dir, "commit", "-qm", "initial")

	bare := ""
	if withRemote {
		bare = t.TempDir()
		runGitIn(t, bare, "init", "--bare", "-q")
		runGitIn(t, dir, "remote", "add", "origin", bare)
		runGitIn(t, dir, "push", "-q", "-u", "origin", "main")
	}

	// Dirty the tree three ways: unstaged edit, staged edit, untracked file.
	if err := appendFixtureFile(t, dir, "a.txt", "alpha-changed\n"); err != nil {
		t.Fatal(err)
	}
	if err := appendFixtureFile(t, dir, "c.txt", "charlie-changed\n"); err != nil {
		t.Fatal(err)
	}
	runGitIn(t, dir, "add", "c.txt")
	if err := writeFixtureFile(t, dir, "untracked.txt", "new\n"); err != nil {
		t.Fatal(err)
	}

	return dir, bare
}

func snapshotFixture(t *testing.T, dir, bare string) fixtureState {
	t.Helper()
	state := fixtureState{
		head:      runGitIn(t, dir, "rev-parse", "HEAD"),
		porcelain: runGitIn(t, dir, "status", "--porcelain"),
		cached:    runGitIn(t, dir, "diff", "--cached"),
	}
	if bare != "" {
		state.bareHead = runGitIn(t, bare, "rev-parse", "HEAD")
	}

	return state
}

func assertFixtureUnchanged(t *testing.T, before, after fixtureState) {
	t.Helper()
	if before != after {
		t.Errorf("dry-run mutated the fixture:\nbefore: %+v\nafter:  %+v", before, after)
	}
}

func TestRunCommitPushFeatureDryRunZeroMutation(t *testing.T) {
	dir, bare := newDryRunFixture(t, true)
	t.Chdir(dir)
	before := snapshotFixture(t, dir, bare)

	assertNilCommitErr(t, RunCommitPushFeature([]string{"probe - dry run", "--dry-run"}))

	assertFixtureUnchanged(t, before, snapshotFixture(t, dir, bare))
}

func TestRunCommitPushBugDryRunShortFlagZeroMutation(t *testing.T) {
	dir, _ := newDryRunFixture(t, false)
	t.Chdir(dir)
	before := snapshotFixture(t, dir, "")

	assertNilCommitErr(t, RunCommitPushBug([]string{"probe - dry run", "-n"}))

	assertFixtureUnchanged(t, before, snapshotFixture(t, dir, ""))
}

func TestRunPullCommitPushDryRunDoesNotPull(t *testing.T) {
	dir, _ := newDryRunFixture(t, false)
	// A remote a real pull could never reach: if dry-run pulled first,
	// this would error out instead of returning nil.
	runGitIn(t, dir, "remote", "add", "origin", filepath.Join(t.TempDir(), "no-such-remote.git"))
	t.Chdir(dir)
	before := snapshotFixture(t, dir, "")

	assertNilCommitErr(t, RunPullCommitPush([]string{"probe - dry run", "--dry-run"}))

	assertFixtureUnchanged(t, before, snapshotFixture(t, dir, ""))
}

func TestRunCommitDryRunZeroMutation(t *testing.T) {
	dir, _ := newDryRunFixture(t, false)
	t.Chdir(dir)
	before := snapshotFixture(t, dir, "")

	if err := RunCommit([]string{"probe - plain dry run", "--dry-run", "-p"}); err != nil {
		t.Fatalf("plain commit dry-run returned error: %v", err)
	}

	assertFixtureUnchanged(t, before, snapshotFixture(t, dir, ""))
}

func TestRunCommitDryRunCleanTree(t *testing.T) {
	dir := t.TempDir()
	runGitIn(t, dir, "init", "-q", "-b", "main")
	runGitIn(t, dir, "config", "user.email", "dryrun@example.com")
	runGitIn(t, dir, "config", "user.name", "Dry Run")
	if err := writeFixtureFile(t, dir, "a.txt", "alpha\n"); err != nil {
		t.Fatal(err)
	}
	runGitIn(t, dir, "add", "-A")
	runGitIn(t, dir, "commit", "-qm", "initial")
	t.Chdir(dir)
	before := snapshotFixture(t, dir, "")

	assertNilCommitErr(t, RunCommitPushFeature([]string{"probe - clean dry run", "--dry-run"}))

	assertFixtureUnchanged(t, before, snapshotFixture(t, dir, ""))
}

func TestParsePorcelainStatus(t *testing.T) {
	got := parsePorcelainStatus(" M a.txt\nM  c.txt\n?? untracked.txt\nR  old.txt -> new.txt\n")
	want := []dryRunFile{
		{Path: "a.txt", Status: " M"},
		{Path: "c.txt", Status: "M "},
		{Path: "untracked.txt", Status: "??"},
		{Path: "new.txt", Status: "R "},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parsePorcelainStatus = %v, want %v", got, want)
	}
	if len(parsePorcelainStatus("")) != 0 {
		t.Error("empty porcelain output should parse to no files")
	}
}

func TestScopeWouldStage(t *testing.T) {
	changed := []dryRunFile{{Path: "a.txt", Status: " M"}, {Path: "b.txt", Status: "??"}}
	staged, unchanged := scopeWouldStage(changed, []string{"a.txt", "missing.txt"})
	if len(staged) != 1 || staged[0].Path != "a.txt" {
		t.Errorf("staged = %v, want only a.txt", staged)
	}
	if len(unchanged) != 1 || unchanged[0] != "missing.txt" {
		t.Errorf("unchanged = %v, want [missing.txt]", unchanged)
	}
}
