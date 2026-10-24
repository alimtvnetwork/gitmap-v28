package cmdrevert

import (
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestRunRevertEmptyArgs(t *testing.T) {
	err := RunRevert([]string{})
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}

	if !strings.Contains(err.Error(), "gitmap revert <version>") {
		t.Errorf("expected usage message, got %q", err.Error())
	}
}

func TestValidateRevertVersionTagNotFound(t *testing.T) {
	err := validateRevertVersion("v999.999.999-nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent tag, got nil")
	}

	if !strings.Contains(err.Error(), "not found locally") {
		t.Errorf("expected not found message, got %q", err.Error())
	}
}

func TestHasRevertFlag(t *testing.T) {
	args := []string{"--list-txn", "--force=true", "v1.0.0"}
	if !hasRevertFlag(args, constants.FlagListTxn) {
		t.Error("expected hasRevertFlag to find --list-txn")
	}

	if !hasRevertFlag(args, constants.FlagForce) {
		t.Error("expected hasRevertFlag to find --force")
	}

	if hasRevertFlag(args, constants.FlagPruneTxn) {
		t.Error("expected hasRevertFlag to return false for --prune-txn")
	}
}

func TestFlagValue(t *testing.T) {
	testFlagValueCutPrefix(t)
	testFlagValueNextArg(t)
	testFlagValueMissing(t)
}

func testFlagValueCutPrefix(t *testing.T) {
	args := []string{"--txn=123"}
	val, hasVal := flagValue(args, constants.FlagTxn)
	if !hasVal || val != "123" {
		t.Errorf("expected val 123, got %q (hasVal=%v)", val, hasVal)
	}
}

func testFlagValueNextArg(t *testing.T) {
	args := []string{"--txn", "456"}
	val, hasVal := flagValue(args, constants.FlagTxn)
	if !hasVal || val != "456" {
		t.Errorf("expected val 456, got %q (hasVal=%v)", val, hasVal)
	}
}

func testFlagValueMissing(t *testing.T) {
	args := []string{"--other", "val"}
	_, hasVal := flagValue(args, constants.FlagTxn)
	if hasVal {
		t.Error("expected hasVal false for missing flag")
	}
}

func TestParseTxnID(t *testing.T) {
	id, err := parseTxnID("42")
	if err != nil || id != 42 {
		t.Fatalf("expected id 42, got %d, err=%v", id, err)
	}

	prevExit := cliexit.SetExitFunc(func(int) {})
	defer cliexit.SetExitFunc(prevExit)

	if _, err := parseTxnID("invalid"); err == nil {
		t.Error("expected error parsing non-numeric id")
	}

	if _, err := parseTxnID("-5"); err == nil {
		t.Error("expected error parsing negative id")
	}

	if _, err := parseTxnID("0"); err == nil {
		t.Error("expected error parsing zero id")
	}
}

func TestMustParseTxnID(t *testing.T) {
	if mustParseTxnID("100") != 100 {
		t.Errorf("expected 100, got %d", mustParseTxnID("100"))
	}
}

func TestParseLastN(t *testing.T) {
	count, err := parseLastN("5")
	if err != nil || count != 5 {
		t.Fatalf("expected count 5, got %d, err=%v", count, err)
	}

	prevExit := cliexit.SetExitFunc(func(int) {})
	defer cliexit.SetExitFunc(prevExit)

	if _, err := parseLastN("abc"); err == nil {
		t.Error("expected error parsing invalid count")
	}

	if _, err := parseLastN("-1"); err == nil {
		t.Error("expected error parsing negative count")
	}

	if mustParseLastN("7") != 7 {
		t.Errorf("expected 7, got %d", mustParseLastN("7"))
	}
}

func TestBuildRevertScript(t *testing.T) {
	script := buildRevertScript("C:\\test\\repo", "C:\\test\\repo\\run.ps1")
	if !strings.Contains(script, "C:\\test\\repo") {
		t.Error("script missing repo path")
	}

	if !strings.Contains(script, "run.ps1") {
		t.Error("script missing run.ps1 reference")
	}
}

func TestFilterCommittedTxns(t *testing.T) {
	records := []model.TransactionRecord{
		{ID: 1, Status: constants.TxnStatusCommitted},
		{ID: 2, Status: constants.TxnStatusReverted},
		{ID: 3, Status: constants.TxnStatusCommitted},
		{ID: 4, Status: constants.TxnStatusCommitted},
	}

	filtered := filterCommittedTxns(records, 2)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 filtered records, got %d", len(filtered))
	}

	if filtered[0].ID != 1 || filtered[1].ID != 3 {
		t.Errorf("unexpected filtered records: %+v", filtered)
	}
}

func TestIsAffirmativeResponse(t *testing.T) {
	if !isAffirmativeResponse("yes") {
		t.Error("expected 'yes' to be affirmative")
	}

	if !isAffirmativeResponse("YES") {
		t.Error("expected 'YES' to be affirmative")
	}

	if isAffirmativeResponse("no") {
		t.Error("expected 'no' not to be affirmative")
	}

	if isAffirmativeResponse("other") {
		t.Error("expected 'other' not to be affirmative")
	}
}

func TestGetRepoPath(t *testing.T) {
	prevFn := RepoPathFn
	defer func() { RepoPathFn = prevFn }()

	RepoPathFn = nil
	if getRepoPath() != "" {
		t.Errorf("expected empty string when RepoPathFn is nil, got %q", getRepoPath())
	}

	RepoPathFn = func() string { return "D:\\test\\custom" }
	if getRepoPath() != "D:\\test\\custom" {
		t.Errorf("expected custom path, got %q", getRepoPath())
	}
}

func TestRevertErrorConstructors(t *testing.T) {
	checkoutErr := newRepoPathError("v1.2.3")
	if checkoutErr.Code != constants.ErrCodeE2012 {
		t.Errorf("expected code E2012, got %s", checkoutErr.Code)
	}

	runnerErr := newRunnerRepoPathError()
	if runnerErr.Code != constants.ErrCodeE2013 {
		t.Errorf("expected code E2013, got %s", runnerErr.Code)
	}

	manyErr := newRevertManyError(errors.New("fail"), 1, "test", true)
	if manyErr.Code != constants.ErrCodeE2015 {
		t.Errorf("expected code E2015, got %s", manyErr.Code)
	}
}

func TestRunRevertRunnerMissingRepoPath(t *testing.T) {
	prevFn := RepoPathFn
	defer func() { RepoPathFn = prevFn }()

	RepoPathFn = func() string { return "" }
	err := RunRevertRunner()
	if err == nil {
		t.Fatal("expected error when repo path is empty")
	}

	if !strings.Contains(err.Error(), "repo path is not set") {
		t.Errorf("expected repo path error message, got %q", err.Error())
	}
}

func TestCheckoutRevertTagMissingRepoPath(t *testing.T) {
	prevFn := RepoPathFn
	defer func() { RepoPathFn = prevFn }()

	RepoPathFn = func() string { return "" }
	err := checkoutRevertTag("v1.0.0")
	if err == nil {
		t.Fatal("expected error when repo path is empty")
	}

	if !strings.Contains(err.Error(), "repo path is not set") {
		t.Errorf("expected repo path error message, got %q", err.Error())
	}
}
