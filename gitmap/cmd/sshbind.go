package cmd

// `gitmap ssh-bind <keyfile>` — pin a specific SSH private key to the
// current repo via `git config core.sshCommand`. Fixes the classic
// "wrong GitHub account is being offered" push failure without any
// global SSH config change. Discover candidate keys with `gitmap whoami`.

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runSSHBind sets core.sshCommand for the cwd repo to the given key.
// Accepts an absolute path, a ~/-relative path, or a bare filename
// (resolved under ~/.ssh). Exits 0 on success, 1 on any failure.
func runSSHBind(args []string) {
	if !isGitRepoCWD() {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap ssh-bind` inside a repo)")
		exitWith(1)

		return
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: gitmap ssh-bind <key-filename-or-path>")
		fmt.Fprintln(os.Stderr, "  tip: run `gitmap whoami` to list available keys under ~/.ssh")
		exitWith(1)

		return
	}
	keyRef, keyPath := resolveSSHKeyPath(args[0])
	if _, err := os.Stat(keyPath); err != nil {
		fmt.Fprintf(os.Stderr, "✗ key not found: %s (%v)\n", keyPath, err)
		exitWith(1)

		return
	}
	cmdStr := fmt.Sprintf("ssh -i %s -F /dev/null -o IdentitiesOnly=yes", keyRef)
	out, err := exec.Command("git", "config", "core.sshCommand", cmdStr).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ git config failed: %v\n%s", err, out)
		exitWith(1)

		return
	}
	fmt.Printf("✓ pinned SSH key for this repo: %s\n", keyPath)
	fmt.Printf("  core.sshCommand = %s\n", cmdStr)
	fmt.Println("  test with: git push")
}

// resolveSSHKeyPath returns (refForSSH, absoluteForStat). Bare names
// map to ~/.ssh/<name>; ~/ paths are expanded for stat but kept as
// ~/ for the ssh -i argument so the config stays portable.
func resolveSSHKeyPath(arg string) (string, string) {
	home, _ := os.UserHomeDir()
	if strings.HasPrefix(arg, "~/") || strings.HasPrefix(arg, "~\\") {
		return arg, home + string(os.PathSeparator) + arg[2:]
	}
	if strings.ContainsAny(arg, "/\\") || (len(arg) > 1 && arg[1] == ':') {
		return arg, arg
	}
	ref := "~/.ssh/" + arg
	abs := home + string(os.PathSeparator) + ".ssh" + string(os.PathSeparator) + arg

	return ref, abs
}
