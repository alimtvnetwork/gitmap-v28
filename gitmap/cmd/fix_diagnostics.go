// Package cmd — fix_diagnostics.go delivers blunt git remediation failure analysis.
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

type gitDiagnosticInfo struct {
	RCA       string
	Solutions []string
}

func printBluntRemediationFailure(repoName string, step gitutil.RemediationStep, output string, err error) {
	diag := analyzeGitErrorOutput(output, err)
	cmdStr := step.Name + " " + strings.Join(step.Args, " ")

	fmt.Printf("\n%s Remediation Step Failed on %s: %v\n", constants.ColorRed+"✖"+constants.ColorReset, repoName, err)
	fmt.Printf("  Command:   %s\n", cmdStr)
	if len(strings.TrimSpace(output)) > 0 {
		fmt.Printf("  Output:\n    %s\n", strings.ReplaceAll(strings.TrimSpace(output), "\n", "\n    "))
	}

	fmt.Printf("  RCA (Root Cause): %s\n", diag.RCA)
	fmt.Printf("  Known Solutions:\n")
	for _, sol := range diag.Solutions {
		fmt.Printf("    • %s\n", sol)
	}

	fmt.Println()
}

func analyzeGitErrorOutput(output string, err error) gitDiagnosticInfo {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "pathspec") {
		return pathspecDiagnostic()
	}

	if strings.Contains(lower, "would be overwritten by merge") || strings.Contains(lower, "would be overwritten by rebase") {
		return overwrittenDiagnostic()
	}

	if strings.Contains(lower, "conflict") {
		return conflictDiagnostic()
	}

	if strings.Contains(lower, "permission denied") || strings.Contains(lower, "could not read username") {
		return authDiagnostic()
	}

	if strings.Contains(lower, "could not resolve host") || strings.Contains(lower, "unable to access") {
		return networkDiagnostic()
	}

	return defaultDiagnostic(err)
}

func pathspecDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Command arguments were misquoted or split incorrectly by shell wrappers, causing git to interpret words in commit messages as file paths.",
		Solutions: []string{
			"Run native command: git -C <repo> commit -m \"wip: local changes\"",
			"Verify argument quoting to ensure spaces do not fragment commit messages into pathspecs",
		},
	}
}

func overwrittenDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Incoming remote commits conflict with local modified or untracked files.",
		Solutions: []string{
			"Stash untracked files: git -C <repo> stash -u",
			"Or discard untracked changes: git -C <repo> clean -fd && git -C <repo> reset --hard HEAD",
		},
	}
}

func conflictDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Merge or rebase conflict detected between local commits and remote branch.",
		Solutions: []string{
			"Check status: git -C <repo> status",
			"Resolve conflicts in editor, then: git -C <repo> rebase --continue (or git rebase --abort)",
		},
	}
}

func authDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Authentication failure connecting to remote repository. SSH key or credentials missing.",
		Solutions: []string{
			"Verify SSH agent keys: ssh-add -l",
			"Verify remote URL: git -C <repo> remote -v",
		},
	}
}

func networkDiagnostic() gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: "Network connectivity failure. Unable to reach remote Git host.",
		Solutions: []string{
			"Verify internet connection and DNS resolution",
			"Retry pull operation once network connectivity is restored",
		},
	}
}

func defaultDiagnostic(err error) gitDiagnosticInfo {
	return gitDiagnosticInfo{
		RCA: fmt.Sprintf("Git process exited with failure: %v", err),
		Solutions: []string{
			"Inspect working tree status: git -C <repo> status",
			"Check git log: git -C <repo> log -n 3 --oneline",
		},
	}
}
