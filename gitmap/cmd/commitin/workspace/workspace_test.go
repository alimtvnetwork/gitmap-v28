package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestEnsureWorkspaceCreatesAllDirsIdempotently verifies that two
// successive calls produce the same Paths and never error.
func TestEnsureWorkspaceCreatesAllDirsIdempotently(t *testing.T) {
	root := t.TempDir()
	first, err := EnsureWorkspace(root)
	if err != nil {
		t.Fatalf("first EnsureWorkspace: %v", err)
	}
	second, err := EnsureWorkspace(root)
	if err != nil {
		t.Fatalf("second EnsureWorkspace: %v", err)
	}
	if *first != *second {
		t.Fatalf("paths drifted between calls: %+v vs %+v", first, second)
	}
	for _, dir := range []string{first.GitmapRoot, first.CommitInRoot, first.ProfilesDir, first.TempRoot} {
		if info, statErr := os.Stat(dir); statErr != nil || !info.IsDir() {
			t.Fatalf("expected dir %s to exist: %v", dir, statErr)
		}
	}
}

// TestAcquireLockBlocksDoubleAcquire ensures the second AcquireLock on
// the same workspace returns the spec §2.7 LockBusy error message.
func TestAcquireLockBlocksDoubleAcquire(t *testing.T) {
	p, _ := EnsureWorkspace(t.TempDir())
	first, err := AcquireLock(p)
	if err != nil {
		t.Fatalf("first AcquireLock: %v", err)
	}
	defer first.Release()
	if _, err := AcquireLock(p); err == nil {
		t.Fatalf("expected lock-busy error on second acquire")
	}
}

// TestAcquireLockReleaseAllowsReacquire confirms Release clears state.
func TestAcquireLockReleaseAllowsReacquire(t *testing.T) {
	p, _ := EnsureWorkspace(t.TempDir())
	h, err := AcquireLock(p)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	h.Release()
	if _, err := AcquireLock(p); err != nil {
		t.Fatalf("reacquire failed: %v", err)
	}
}

// TestEnsureSourceInitsMissingDir covers spec §2.3 case 4. Uses fake
// git runner so the test stays hermetic.
func TestEnsureSourceInitsMissingDir(t *testing.T) {
	restore := SetGitRunnerForTest(func(_ string, _ ...string) error { return nil })
	defer restore()
	root := filepath.Join(t.TempDir(), "newrepo")
	h, err := EnsureSource(root)
	if err != nil {
		t.Fatalf("EnsureSource: %v", err)
	}
	if h.Kind != SourceKindCreatedAndInit || !h.IsFreshlyInit {
		t.Fatalf("expected CreatedAndInit, got %+v", h)
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("dir not created: %v", err)
	}
}

// TestEnsureSourceReusesExistingRepo covers spec §2.3 case 2.
func TestEnsureSourceReusesExistingRepo(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	h, err := EnsureSource(root)
	if err != nil {
		t.Fatalf("EnsureSource: %v", err)
	}
	if h.Kind != SourceKindExistingRepo || h.IsFreshlyInit {
		t.Fatalf("expected ExistingRepo, got %+v", h)
	}
}

// TestEnsureSourceClonesUrl covers spec §2.3 case 1 via fake runner.
func TestEnsureSourceClonesUrl(t *testing.T) {
	calls := 0
	restore := SetGitRunnerForTest(func(sub string, args ...string) error {
		calls++
		if sub != "clone" {
			t.Fatalf("expected clone, got %s", sub)
		}
		return nil
	})
	defer restore()
	cwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(cwd) }()
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	h, err := EnsureSource("https://example.com/foo.git")
	if err != nil {
		t.Fatalf("EnsureSource: %v", err)
	}
	if h.Kind != SourceKindCloned || calls != 1 {
		t.Fatalf("expected single clone call, got kind=%v calls=%d", h.Kind, calls)
	}
	if filepath.Base(h.Path) != "foo" {
		t.Fatalf("expected basename foo, got %s", h.Path)
	}
}

// TestExpandInputsKeywordAllSortsAscending covers spec §2.4 ordering.
func TestExpandInputsKeywordAllSortsAscending(t *testing.T) {
	parent := t.TempDir()
	mustMkdir(t, filepath.Join(parent, "demo"))
	mustMkdir(t, filepath.Join(parent, "demo-v1"))
	mustMkdir(t, filepath.Join(parent, "demo-v3"))
	mustMkdir(t, filepath.Join(parent, "demo-v10"))
	mustMkdir(t, filepath.Join(parent, "unrelated"))
	source := filepath.Join(parent, "demo-v3")
	got, err := ExpandInputs(source, nil, constants.CommitInInputKeywordAll, 0)
	if err != nil {
		t.Fatalf("expand: %v", err)
	}
	wantNames := []string{"demo", "demo-v1", "demo-v10"}
	if len(got) != len(wantNames) {
		t.Fatalf("got %d siblings (%v), want %d (%v)", len(got), namesOf(got), len(wantNames), wantNames)
	}
	for i, w := range wantNames {
		if got[i].Original != w {
			t.Fatalf("position %d: got %q want %q", i, got[i].Original, w)
		}
		if got[i].OrderIndex != i+1 {
			t.Fatalf("position %d: order index = %d want %d", i, got[i].OrderIndex, i+1)
		}
	}
}

// TestExpandInputsTailKeywordTruncates covers `-N`.
func TestExpandInputsTailKeywordTruncates(t *testing.T) {
	parent := t.TempDir()
	mustMkdir(t, filepath.Join(parent, "p"))
	mustMkdir(t, filepath.Join(parent, "p-v1"))
	mustMkdir(t, filepath.Join(parent, "p-v2"))
	mustMkdir(t, filepath.Join(parent, "p-v3"))
	source := filepath.Join(parent, "p-v3")
	got, err := ExpandInputs(source, nil, "-2", 2)
	if err != nil {
		t.Fatalf("expand: %v", err)
	}
	if len(got) != 2 || got[0].Original != "p-v1" || got[1].Original != "p-v2" {
		t.Fatalf("tail truncation wrong: %v", namesOf(got))
	}
}

// TestExpandInputsExplicitClassifiesUrlsAndPaths verifies per-token
// kind assignment with a mixed input list.
func TestExpandInputsExplicitClassifiesUrlsAndPaths(t *testing.T) {
	dir := t.TempDir()
	got, err := ExpandInputs("ignored", []string{dir, "https://example.com/x.git"}, "", 0)
	if err != nil {
		t.Fatalf("expand: %v", err)
	}
	if got[0].Kind != constants.CommitInInputKindLocalFolder {
		t.Fatalf("expected LocalFolder, got %q", got[0].Kind)
	}
	if got[1].Kind != constants.CommitInInputKindGitUrl || got[1].URL == "" {
		t.Fatalf("expected GitUrl with URL set, got %+v", got[1])
	}
}

// TestCloneInputsStagesAllThreeKinds runs CloneInputs end-to-end with
// a fake git runner; verifies that local folders are reused in place
// and the other two kinds invoke `git clone` with the right target.
func TestCloneInputsStagesAllThreeKinds(t *testing.T) {
	cloneTargets := []string{}
	restore := SetGitRunnerForTest(func(sub string, args ...string) error {
		if sub == "clone" && len(args) == 2 {
			cloneTargets = append(cloneTargets, args[1])
			return os.MkdirAll(args[1], 0o755)
		}
		return nil
	})
	defer restore()
	wsRoot := t.TempDir()
	p, _ := EnsureWorkspace(wsRoot)
	localDir := t.TempDir()
	siblingDir := t.TempDir()
	inputs := []ResolvedInput{
		{OrderIndex: 1, Original: "local", Kind: constants.CommitInInputKindLocalFolder, AbsPath: localDir, Version: -1},
		{OrderIndex: 2, Original: "https://example.com/r.git", Kind: constants.CommitInInputKindGitUrl, URL: "https://example.com/r.git", Version: -1},
		{OrderIndex: 3, Original: filepath.Base(siblingDir), Kind: constants.CommitInInputKindVersionedSibling, AbsPath: siblingDir, Version: 1},
	}
	staged, err := CloneInputs(p, 42, inputs)
	if err != nil {
		t.Fatalf("CloneInputs: %v", err)
	}
	if len(staged) != 3 {
		t.Fatalf("staged %d, want 3", len(staged))
	}
	if staged[0].WorkPath != localDir || staged[0].IsClone {
		t.Fatalf("local folder should be reused in place, got %+v", staged[0])
	}
	if !staged[1].IsClone || !strings.Contains(staged[1].WorkPath, "2-r") {
		t.Fatalf("url stage path wrong: %+v", staged[1])
	}
	if !staged[2].IsClone || !strings.HasPrefix(staged[2].WorkPath, p.TempRoot) {
		t.Fatalf("sibling stage path wrong: %+v", staged[2])
	}
	if len(cloneTargets) != 2 {
		t.Fatalf("expected 2 clone calls, got %d (%v)", len(cloneTargets), cloneTargets)
	}
}

func mustMkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
}

func namesOf(in []ResolvedInput) []string {
	out := make([]string, len(in))
	for i, r := range in {
		out[i] = r.Original
	}
	return out
}
