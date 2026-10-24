package cmdprune

import (
	"errors"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprune/enums"
)

func TestParsePruneFlags(t *testing.T) {
	testParsePruneFlagsDefaults(t)
	testParsePruneFlagsCombinations(t)
}

func testParsePruneFlagsDefaults(t *testing.T) {
	isDryRun, isConfirmed, isRemote := parsePruneFlags([]string{})
	if isDryRun || isConfirmed || isRemote {
		t.Error("expected default flags to be false")
	}
}

func testParsePruneFlagsCombinations(t *testing.T) {
	args := []string{constants.FlagArgDryRun, constants.FlagArgConfirm, constants.FlagArgRemote}
	isDryRun, isConfirmed, isRemote := parsePruneFlags(args)
	if !isDryRun || !isConfirmed || !isRemote {
		t.Error("expected all flags to be parsed as true")
	}
}

func TestResolvePruneAction(t *testing.T) {
	if resolvePruneAction(true, false) != enums.PruneActionDryRun {
		t.Error("expected dry-run action")
	}
	if resolvePruneAction(false, true) != enums.PruneActionDelete {
		t.Error("expected delete action")
	}
	if resolvePruneAction(false, false) != enums.PruneActionPrompt {
		t.Error("expected prompt action")
	}
}

func TestResolvePruneScope(t *testing.T) {
	if resolvePruneScope(true) != enums.PruneScopeAll {
		t.Error("expected all scope when remote is true")
	}
	if resolvePruneScope(false) != enums.PruneScopeLocal {
		t.Error("expected local scope when remote is false")
	}
}

func TestIsPruneConfirmed(t *testing.T) {
	testIsPruneConfirmedDirect(t)
	testIsPruneConfirmedPromptYes(t)
	testIsPruneConfirmedPromptNo(t)
}

func testIsPruneConfirmedDirect(t *testing.T) {
	if !isPruneConfirmed(true, 1) {
		t.Error("expected true when confirmed flag is set")
	}
}

func testIsPruneConfirmedPromptYes(t *testing.T) {
	origReader := promptReaderFn
	defer func() { promptReaderFn = origReader }()

	promptReaderFn = func() string { return constants.PromptAnswerY }
	if !isPruneConfirmed(false, 2) {
		t.Error("expected true when prompt answered 'y'")
	}

	promptReaderFn = func() string { return constants.PromptAnswerYUpper }
	if !isPruneConfirmed(false, 2) {
		t.Error("expected true when prompt answered 'Y'")
	}
}

func testIsPruneConfirmedPromptNo(t *testing.T) {
	origReader := promptReaderFn
	defer func() { promptReaderFn = origReader }()

	promptReaderFn = func() string { return constants.PromptAnswerN }
	if isPruneConfirmed(false, 3) {
		t.Error("expected false when prompt answered 'n'")
	}

	promptReaderFn = func() string { return "invalid" }
	if isPruneConfirmed(false, 3) {
		t.Error("expected false when prompt answered with invalid input")
	}
}

func TestBranchToTag(t *testing.T) {
	if branchToTag("release/v1.0.0") != "v1.0.0" {
		t.Error("failed converting release/v1.0.0")
	}
	if branchToTag("release/v2.10.4") != "v2.10.4" {
		t.Error("failed converting release/v2.10.4")
	}
	if branchToTag("main") != "main" {
		t.Error("unexpected conversion for non-release branch")
	}
}

func TestParseGitBranchOutput(t *testing.T) {
	testParseGitBranchOutputEmpty(t)
	testParseGitBranchOutputValid(t)
}

func testParseGitBranchOutputEmpty(t *testing.T) {
	if parseGitBranchOutput(nil) != nil {
		t.Error("expected nil on nil output")
	}
	if parseGitBranchOutput([]byte("   \n \n\t")) != nil {
		t.Error("expected nil on whitespace output")
	}
}

func testParseGitBranchOutputValid(t *testing.T) {
	raw := []byte("  release/v1.0.0  \n  release/v1.1.0\n")
	branches := parseGitBranchOutput(raw)
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	if branches[0] != "release/v1.0.0" || branches[1] != "release/v1.1.0" {
		t.Error("mismatch in parsed branches")
	}
}

func TestFilterStaleBranches(t *testing.T) {
	origTagExists := tagExistsFn
	defer func() { tagExistsFn = origTagExists }()

	tagExistsFn = func(tag string) bool {
		return tag == "v1.0.0"
	}

	branches := []string{"release/v1.0.0", "release/v2.0.0"}
	stale := filterStaleBranches(branches)
	if len(stale) != 1 {
		t.Fatalf("expected 1 stale branch, got %d", len(stale))
	}
	if stale[0].name != "release/v1.0.0" || stale[0].tag != "v1.0.0" {
		t.Error("mismatch in filtered stale branch")
	}
}

func TestDeleteSingleBranch(t *testing.T) {
	testDeleteSingleBranchSuccess(t)
	testDeleteSingleBranchFailure(t)
}

func testDeleteSingleBranchSuccess(t *testing.T) {
	origDeleteLocal := deleteLocalBranchFn
	origDeleteRemote := deleteRemoteBranchFn
	defer func() {
		deleteLocalBranchFn = origDeleteLocal
		deleteRemoteBranchFn = origDeleteRemote
	}()

	hasRemoteCalled := false
	deleteLocalBranchFn = func(name string) error { return nil }
	deleteRemoteBranchFn = func(name string) { hasRemoteCalled = true }

	sb := staleBranch{name: "release/v1.0.0", tag: "v1.0.0"}
	deleted := deleteSingleBranch(sb, true)
	if deleted != 1 || !hasRemoteCalled {
		t.Error("expected successful branch deletion and remote invocation")
	}
}

func testDeleteSingleBranchFailure(t *testing.T) {
	origDeleteLocal := deleteLocalBranchFn
	origDeleteRemote := deleteRemoteBranchFn
	defer func() {
		deleteLocalBranchFn = origDeleteLocal
		deleteRemoteBranchFn = origDeleteRemote
	}()

	hasRemoteCalled := false
	deleteLocalBranchFn = func(name string) error { return errors.New("delete error") }
	deleteRemoteBranchFn = func(name string) { hasRemoteCalled = true }

	sb := staleBranch{name: "release/v1.0.0", tag: "v1.0.0"}
	deleted := deleteSingleBranch(sb, true)
	if deleted != 0 || hasRemoteCalled {
		t.Error("expected failed deletion and no remote invocation")
	}
}

func TestDeleteStaleBranches(t *testing.T) {
	origDeleteLocal := deleteLocalBranchFn
	defer func() { deleteLocalBranchFn = origDeleteLocal }()

	deleteLocalBranchFn = func(name string) error { return nil }
	stale := []staleBranch{
		{name: "release/v1.0.0", tag: "v1.0.0"},
		{name: "release/v1.1.0", tag: "v1.1.0"},
	}

	deleteStaleBranches(stale, false)
}

func TestRunPrune(t *testing.T) {
	testRunPruneNoStale(t)
	testRunPruneDryRun(t)
	testRunPruneConfirmed(t)
}

func testRunPruneNoStale(t *testing.T) {
	origList := listBranchesFn
	origTag := tagExistsFn
	origCheckHelp := CheckHelpFn
	defer func() {
		listBranchesFn = origList
		tagExistsFn = origTag
		CheckHelpFn = origCheckHelp
	}()

	hasCheckHelpCalled := false
	CheckHelpFn = func(command string, args []string) {
		hasCheckHelpCalled = true
	}
	listBranchesFn = func() []string { return []string{"release/v1.0.0"} }
	tagExistsFn = func(tag string) bool { return false }

	err := RunPrune([]string{})
	if err != nil || !hasCheckHelpCalled {
		t.Errorf("expected clean exit on no stale branches, got %v", err)
	}
}

func testRunPruneDryRun(t *testing.T) {
	origList := listBranchesFn
	origTag := tagExistsFn
	origDeleteLocal := deleteLocalBranchFn
	defer func() {
		listBranchesFn = origList
		tagExistsFn = origTag
		deleteLocalBranchFn = origDeleteLocal
	}()

	hasDeleteCalled := false
	listBranchesFn = func() []string { return []string{"release/v1.0.0"} }
	tagExistsFn = func(tag string) bool { return true }
	deleteLocalBranchFn = func(name string) error {
		hasDeleteCalled = true

		return nil
	}

	err := RunPrune([]string{constants.FlagArgDryRun})
	if err != nil || hasDeleteCalled {
		t.Errorf("expected clean dry run without deletion, got %v", err)
	}
}

func testRunPruneConfirmed(t *testing.T) {
	origList := listBranchesFn
	origTag := tagExistsFn
	origDeleteLocal := deleteLocalBranchFn
	defer func() {
		listBranchesFn = origList
		tagExistsFn = origTag
		deleteLocalBranchFn = origDeleteLocal
	}()

	hasDeleteCalled := false
	listBranchesFn = func() []string { return []string{"release/v1.0.0"} }
	tagExistsFn = func(tag string) bool { return true }
	deleteLocalBranchFn = func(name string) error {
		hasDeleteCalled = true

		return nil
	}

	err := RunPrune([]string{constants.FlagArgConfirm})
	if err != nil || !hasDeleteCalled {
		t.Errorf("expected branch deletion on confirm, got %v", err)
	}
}
