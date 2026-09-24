// Package cmd — lowercasefix_confirm.go handles pre-flight verification and user confirmation.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var lcfStdinReader io.Reader = os.Stdin

func promptPreflightConfirmation(pairs []RenamePair, opts LowerCaseFixOptions) (bool, error) {
	return promptPreflightConfirmationWithReader(lcfStdinReader, pairs, opts)
}

func promptPreflightConfirmationWithReader(r io.Reader, pairs []RenamePair, opts LowerCaseFixOptions) (bool, error) {
	renderPreflightBox(pairs, opts)
	fmt.Printf("%sType 'confirm' or 'yes' (or 'y') to proceed: %s", constants.ColorYellow, constants.ColorReset)

	reader := bufio.NewReader(r)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false, nil
	}

	ans := strings.ToLower(strings.TrimSpace(line))
	isConfirmed := ans == "confirm" || ans == "yes" || ans == "y"

	return isConfirmed, nil
}

func renderPreflightBox(pairs []RenamePair, opts LowerCaseFixOptions) {
	fmt.Printf("\n%s════════════════════════════════════════════════════════════════%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s⚡ Pre-Flight Verification: Git Lowercase Renamer%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s════════════════════════════════════════════════════════════════%s\n", constants.ColorCyan, constants.ColorReset)
	renderPreflightFound(pairs)
	renderPreflightChanges()
	renderPreflightAfterEffects(opts)
	fmt.Printf("%s────────────────────────────────────────────────────────────────%s\n", constants.ColorCyan, constants.ColorReset)
}

func renderPreflightFound(pairs []RenamePair) {
	fmt.Printf("  ● Found: %d uppercase file(s) to rename:\n", len(pairs))
	limit := 8
	for i, p := range pairs {
		if i >= limit {
			fmt.Printf("    ... and %d more file(s)\n", len(pairs)-limit)
			break
		}
		fmt.Printf("    • %s → %s\n", p.OldBase, p.NewBase)
	}
}

func renderPreflightChanges() {
	fmt.Printf("\n  ● Planned Change: Safe 2-step atomic git mv\n")
	fmt.Printf("    1. Move to temp:   git mv <file> <file>.tmp-lcf (breaks case-collision)\n")
	fmt.Printf("    2. Move to target: git mv <file>.tmp-lcf <file_lowercase>\n")
}

func renderPreflightAfterEffects(opts LowerCaseFixOptions) {
	branch := resolveCurrentBranchName()
	fmt.Printf("\n  ● After-Effects:\n")
	fmt.Printf("    • Working tree and Git index synchronized (git add -A)\n")
	if opts.IsNoCommit {
		fmt.Printf("    • Files staged in index (commit skipped due to --no-commit)\n")
		return
	}
	fmt.Printf("    • Automatically committed to branch: %s\n", branch)
}

func resolveCurrentBranchName() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "current branch"
	}
	name := strings.TrimSpace(string(out))
	if name == "" {
		return "current branch"
	}

	return name
}

func renderCanceledMessage() {
	fmt.Printf("\n%s✖ Operation canceled by user. No files were modified or committed.%s\n\n",
		constants.ColorYellow, constants.ColorReset)
}
