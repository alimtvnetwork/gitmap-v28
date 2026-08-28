// Package cmd — release_tools.go: release-notes, release-dry, tag-rename.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func runReleaseNotes(args []string) error {
	// Delegates to the flag-aware implementation in release_notes_opts.go,
	// which still accepts the legacy bare-range positional form.
	runReleaseNotesV2(args)
	return nil
}

func runReleaseDry(args []string) error {
	tag := ""
	if len(args) > 0 {
		tag = args[0]
	}
	fmt.Println("\033[1;96m▸ release-dry\033[0m  build + local tag rehearsal (no push)")
	if err := runStep("go build ./...", "go", "build", "./..."); err != nil {
		return apperror.New("fatal error", "E9000", nil)
	}
	if tag == "" {
		fmt.Println("\033[1;92m✓ dry release complete\033[0m  nothing pushed")
		return nil
	}
	if err := runStep("git tag "+tag, "git", "tag", tag); err != nil {
		return apperror.New("fatal error", "E9000", nil)
	}
	fmt.Printf("\033[1;94mnotes for %s\033[0m\n", tag)
	_ = runStep("git log -10 --oneline", "git", "log", "-10", "--oneline")
	fmt.Printf("\n\033[2;37mundo:  \033[0m \033[1;96mgit tag -d %s\033[0m\n", tag)
	fmt.Println("\033[1;92m✓ dry release complete\033[0m  nothing pushed")
	return nil
}

func runStep(label string, name string, args ...string) error {
	fmt.Printf("  \033[2;37m$\033[0m %s\n", label)
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

func runTagRename(args []string) error {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "tag-rename: ERROR usage: gitmap tag-rename <old> <new>")
		os.Exit(2)
	}
	oldTag, newTag := args[0], args[1]
	steps := [][]string{
		{"git", "tag", newTag, oldTag},
		{"git", "tag", "-d", oldTag},
		{"git", "push", "origin", newTag},
		{"git", "push", "origin", ":refs/tags/" + oldTag},
	}
	for _, s := range steps {
		if err := runStep(strings.Join(s, " "), s[0], s[1:]...); err != nil {
			return apperror.Wrap(err, "tag-rename: ERROR step failed:", nil)
		}
	}
	fmt.Printf("\033[1;92m✓ renamed\033[0m %s → %s (local + origin)\n", oldTag, newTag)
	return nil
}
