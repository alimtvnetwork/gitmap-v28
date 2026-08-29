package ignore

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RunIgnore appends a pattern to .gitignore idempotently.
func RunIgnore(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("RunIgnore", "E_IGNORE_MISSING_PATTERN")
	}
	pattern := args[0]
	return addToGitignore(pattern)
}

// RunIgnoreRm removes matching files from git history then ignores them.
func RunIgnoreRm(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("RunIgnoreRm", "E_IGNORERM_MISSING_PATTERN")
	}
	pattern := args[0]

	pterm.Info.Printf("Removing '%s' from git history...\n", pattern)

	// Single quotes are safer for filter-branch.
	filterCmd := fmt.Sprintf("git rm --cached --ignore-unmatch -r '%s'", pattern)

	cmd := exec.Command("git", "filter-branch", "--force", "--index-filter", filterCmd, "--prune-empty", "--tag-name-filter", "cat", "--", "--all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "RunIgnoreRm: history rewrite failed")
	}

	exec.Command("git", "for-each-ref", "--format=%(refname)", "refs/original/").Run()
	exec.Command("git", "reflog", "expire", "--expire=now", "--all").Run()
	exec.Command("git", "gc", "--prune=now", "--aggressive").Run()

	pterm.Success.Println("Git history successfully rewritten.")

	return addToGitignore(pattern)
}

func addToGitignore(pattern string) error {
	b, err := os.ReadFile(".gitignore")
	if err == nil && containsGitignorePattern(string(b), pattern) {
		pterm.Info.Printf("Pattern '%s' already in .gitignore\n", pattern)
		return nil
	}

	content := formatGitignoreContent(string(b), pattern)
	if err := os.WriteFile(".gitignore", []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "addToGitignore: failed to write file")
	}
	pterm.Success.Printf("Added '%s' to .gitignore\n", pattern)
	return nil
}

func containsGitignorePattern(content, pattern string) bool {
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		if strings.TrimSpace(l) == pattern {
			return true
		}
	}
	return false
}

func formatGitignoreContent(content, pattern string) string {
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + pattern + "\n"
}
