package cmdbranch

import (
	"errors"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdbranch/constants"
)

func TestRunBranch_NoArgs(t *testing.T) {
	hasCheckHelpCalled := false
	hasPrintHelpCalled := false

	CheckHelpFn = func(command string, args []string) {
		if command == constants.CmdBranch {
			hasCheckHelpCalled = true
		}
	}
	PrintHelpAndExitFn = func(command string, args []string) {
		if command == constants.CmdBranch {
			hasPrintHelpCalled = true
		}
	}
	t.Cleanup(func() {
		CheckHelpFn = nil
		PrintHelpAndExitFn = nil
	})

	err := RunBranch([]string{})
	if err != nil {
		t.Errorf("expected nil error on empty args, got %v", err)
	}
	if !hasCheckHelpCalled {
		t.Error("expected CheckHelpFn to be called")
	}
	if !hasPrintHelpCalled {
		t.Error("expected PrintHelpAndExitFn to be called")
	}
}

func TestRunBranch_UnknownSubcommand(t *testing.T) {
	err := RunBranch([]string{"nonexistent-subcommand"})
	if err == nil {
		t.Fatal("expected error on unknown subcommand, got nil")
	}
	if err.Code != constants.ErrCodeE9000 {
		t.Errorf("expected error code %s, got %s", constants.ErrCodeE9000, err.Code)
	}
	if !strings.Contains(err.Op, "nonexistent-subcommand") {
		t.Errorf("expected error op to contain subcommand, got: %s", err.Op)
	}
}

func TestRunBranch_DefaultSuccess(t *testing.T) {
	testRunBranchDefaultVariants(t, "default")
	testRunBranchDefaultVariants(t, "def")
}

func testRunBranchDefaultVariants(t *testing.T, sub string) {
	hasCheckoutCalled := false
	IsInsideWorkTreeFn = func() bool {
		return true
	}
	ResolveDefaultBranchNameFn = func(repoPath string) string {
		return "main"
	}
	CheckoutBranchFn = func(repoPath, branch string) (string, error) {
		hasCheckoutCalled = true

		return "Switched to branch 'main'", nil
	}
	t.Cleanup(func() {
		IsInsideWorkTreeFn = nil
		ResolveDefaultBranchNameFn = nil
		CheckoutBranchFn = nil
	})

	err := RunBranch([]string{sub})
	if err != nil {
		t.Errorf("expected nil error on valid %s subcommand, got %v", sub, err)
	}
	if !hasCheckoutCalled {
		t.Errorf("expected checkout to be called for %s", sub)
	}
}

func TestRunBranch_NotInsideRepo(t *testing.T) {
	IsInsideWorkTreeFn = func() bool {
		return false
	}
	t.Cleanup(func() {
		IsInsideWorkTreeFn = nil
	})

	err := RunBranch([]string{"default"})
	if err == nil {
		t.Fatal("expected error when not inside repo, got nil")
	}
	if err.Code != constants.ErrCodeE9000 {
		t.Errorf("expected error code %s, got %s", constants.ErrCodeE9000, err.Code)
	}
	if err.Op != constants.ErrFatalMessage {
		t.Errorf("expected fatal error op, got %s", err.Op)
	}
}

func TestRunBranch_CheckoutError(t *testing.T) {
	IsInsideWorkTreeFn = func() bool {
		return true
	}
	ResolveDefaultBranchNameFn = func(repoPath string) string {
		return "main"
	}
	CheckoutBranchFn = func(repoPath, branch string) (string, error) {
		return "", errors.New("git checkout failed")
	}
	t.Cleanup(func() {
		IsInsideWorkTreeFn = nil
		ResolveDefaultBranchNameFn = nil
		CheckoutBranchFn = nil
	})

	err := RunBranch([]string{"default"})
	if err == nil {
		t.Fatal("expected error when checkout fails, got nil")
	}
}

func TestHelpMenu(t *testing.T) {
	menu := buildBranchHelpMenu()
	if menu.Title != constants.HelpTitleBranch {
		t.Errorf("unexpected menu title: %s", menu.Title)
	}
	if len(menu.UsageLines) != 3 {
		t.Errorf("expected 3 usage lines, got %d", len(menu.UsageLines))
	}
	if len(menu.FooterFlags) != 2 {
		t.Errorf("expected 2 footer flags, got %d", len(menu.FooterFlags))
	}
	if len(menu.Sections) != 1 {
		t.Errorf("expected 1 section, got %d", len(menu.Sections))
	}
	if len(menu.Tips) != 2 {
		t.Errorf("expected 2 tips, got %d", len(menu.Tips))
	}

	RenderBranchHelp()
}

func TestHookShimDirect(t *testing.T) {
	testCheckHelpHook(t)
	testPrintHelpHook(t)
	testWorkTreeAndCheckoutHooks(t)
}

func testCheckHelpHook(t *testing.T) {
	checkHelp(constants.CmdBranch, nil)
	hasInvoked := false
	CheckHelpFn = func(cmd string, args []string) {
		hasInvoked = true
	}
	checkHelp(constants.CmdBranch, nil)
	CheckHelpFn = nil
	if !hasInvoked {
		t.Error("expected CheckHelpFn to be invoked")
	}
}

func testPrintHelpHook(t *testing.T) {
	printHelpAndExit(constants.CmdBranch, nil)
	hasInvoked := false
	PrintHelpAndExitFn = func(cmd string, args []string) {
		hasInvoked = true
	}
	printHelpAndExit(constants.CmdBranch, nil)
	PrintHelpAndExitFn = nil
	if !hasInvoked {
		t.Error("expected PrintHelpAndExitFn to be invoked")
	}
}

func testWorkTreeAndCheckoutHooks(t *testing.T) {
	IsInsideWorkTreeFn = func() bool {
		return true
	}
	if !isInsideWorkTree() {
		t.Error("expected isInsideWorkTree to return true from hook")
	}
	IsInsideWorkTreeFn = nil

	CheckoutBranchFn = func(path, branch string) (string, error) {
		return "mocked", nil
	}
	out, err := checkoutBranch(".", "feat")
	if out != "mocked" || err != nil {
		t.Errorf("unexpected checkoutBranch result: %s, %v", out, err)
	}
	CheckoutBranchFn = nil

	ResolveDefaultBranchNameFn = func(path string) string {
		return "trunk"
	}
	if resolveDefaultBranchName(".") != "trunk" {
		t.Error("expected trunk from ResolveDefaultBranchNameFn")
	}
	ResolveDefaultBranchNameFn = nil
}
