package cmd

// `gitmap fix-auth --user <github-user> [--email x@y] [--yes] [--force]`
// end-to-end fix for the classic "Permission denied to <wrong-user>"
// SSH push failure. Mirrors the PowerShell / Bash recipe from the
// support doc but runs cross-platform from Go:
//
//   1. Ensure ~/.ssh exists with 0700.
//   2. Generate ed25519 key at ~/.ssh/id_ed25519_<user> (unless present
//      and --force not passed).
//   3. Pin the current repo to that key via
//      `git config core.sshCommand "ssh -i <path> -F /dev/null -o IdentitiesOnly=yes"`.
//   4. Copy the public key to the OS clipboard and print it with a
//      GitHub "add SSH key" URL so the user can paste + push.

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runFixAuth parses flags and orchestrates the auth-fix pipeline.
func runFixAuth(args []string) {
	user, email, assumeYes, force := parseFixAuthFlags(args)
	if user == "" {
		fmt.Fprintln(os.Stderr, "✗ --user <github-username> is required")
		fmt.Fprintln(os.Stderr, "  example: gitmap fix-auth --user aukgit --email me@example.com")
		exitWith(1)

		return
	}
	if !isGitRepoCWD() {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run inside the repo you want to fix)")
		exitWith(1)

		return
	}
	keyPath := fixAuthKeyPath(user)
	if err := ensureSSHDir(filepath.Dir(keyPath)); err != nil {
		fmt.Fprintf(os.Stderr, "✗ mkdir ~/.ssh failed: %v\n", err)
		exitWith(1)

		return
	}
	fixAuthGenerate(keyPath, resolveFixAuthEmail(email), assumeYes, force)
	fixAuthBind(keyPath)
	fixAuthAnnounce(user, keyPath)
}

// parseFixAuthFlags parses fix-auth CLI flags with short aliases.
func parseFixAuthFlags(args []string) (user, email string, assumeYes, force bool) {
	fs := flag.NewFlagSet(constants.CmdFixAuth, flag.ExitOnError)
	userFlag := fs.String("user", "", "GitHub username (required)")
	fs.StringVar(userFlag, "u", "", "GitHub username (short)")
	emailFlag := fs.String("email", "", "Email comment (defaults to git config user.email)")
	fs.StringVar(emailFlag, "e", "", "Email (short)")
	yesFlag := fs.Bool("yes", false, "Skip confirmation prompts")
	fs.BoolVar(yesFlag, "y", false, "Yes (short)")
	forceFlag := fs.Bool("force", false, "Overwrite existing key")
	fs.BoolVar(forceFlag, "f", false, "Force (short)")
	fs.Parse(args)
	// First positional token acts as --user when flag is empty.
	if *userFlag == "" && fs.NArg() > 0 {
		*userFlag = fs.Arg(0)
	}

	return *userFlag, *emailFlag, *yesFlag, *forceFlag
}

// resolveFixAuthEmail falls back to git global user.email.
func resolveFixAuthEmail(email string) string {
	if email != "" {
		return email
	}
	resolved := resolveGitEmail()
	if resolved == "" {
		fmt.Fprintln(os.Stderr, "  ⚠ git user.email not set; using placeholder — pass --email to override")

		return "gitmap-fix-auth@localhost"
	}

	return resolved
}

// fixAuthKeyPath returns the absolute path for the per-user key.
func fixAuthKeyPath(user string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return filepath.Join(home, ".ssh", "id_ed25519_"+user)
}

// fixAuthGenerate runs ssh-keygen unless the key already exists.
func fixAuthGenerate(keyPath, email string, assumeYes, force bool) {
	if _, err := os.Stat(keyPath); err == nil {
		if !force {
			fmt.Printf("• key already exists, reusing: %s\n", keyPath)

			return
		}
		if !assumeYes && !confirmOverwrite(keyPath) {
			fmt.Println("• aborted; existing key kept")
			exitWith(0)

			return
		}
		_ = os.Remove(keyPath)
		_ = os.Remove(keyPath + ".pub")
	}
	if err := validateSSHKeygen(); err != nil {
		fmt.Fprint(os.Stderr, constants.ErrSSHKeygenMissing)
		exitWith(1)

		return
	}
	runSSHKeygenEd25519(keyPath, email)
}

// runSSHKeygenEd25519 invokes ssh-keygen with an empty passphrase.
func runSSHKeygenEd25519(keyPath, email string) {
	cmd := exec.Command(constants.SSHKeygenBin,
		"-t", "ed25519", "-C", email, "-f", keyPath, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ ssh-keygen failed: %v\n", err)
		exitWith(1)
	}
}

// confirmOverwrite prompts y/N before clobbering an existing key.
func confirmOverwrite(keyPath string) bool {
	fmt.Printf("⚠ overwrite existing key %s? [y/N]: ", keyPath)
	var ans string
	_, _ = fmt.Scanln(&ans)

	return strings.EqualFold(strings.TrimSpace(ans), "y")
}

// fixAuthBind pins the given key to the cwd repo's core.sshCommand.
func fixAuthBind(keyPath string) {
	ref := sshBindRefForKey(keyPath)
	cmdStr := fmt.Sprintf("ssh -i %s -F /dev/null -o IdentitiesOnly=yes", ref)
	out, err := exec.Command("git", "config", "core.sshCommand", cmdStr).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ git config core.sshCommand failed: %v\n%s", err, out)
		exitWith(1)

		return
	}
	fmt.Printf("✓ pinned repo → %s\n", cmdStr)
}

// sshBindRefForKey shortens $HOME/.ssh/... to ~/.ssh/... for portability.
func sshBindRefForKey(keyPath string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return keyPath
	}
	if strings.HasPrefix(keyPath, home) {
		return "~" + strings.ReplaceAll(keyPath[len(home):], "\\", "/")
	}

	return keyPath
}

// fixAuthAnnounce prints the public key + next-step guidance and
// pushes the key into the OS clipboard.
func fixAuthAnnounce(user, keyPath string) {
	pub, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ read public key failed: %v\n", err)
		exitWith(1)

		return
	}
	trimmed := strings.TrimSpace(string(pub))
	fmt.Println("\n=== PUBLIC KEY (add this to GitHub) ===")
	fmt.Println(trimmed)
	fmt.Println("=======================================")
	copyPubKeyAndAnnounce(trimmed)
	fmt.Printf("\nNext steps for %s:\n", user)
	fmt.Println("  1. Open https://github.com/settings/ssh/new")
	fmt.Println("  2. Paste the key above, save.")
	fmt.Println("  3. Run:  git push")
}
