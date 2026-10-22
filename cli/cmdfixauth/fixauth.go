package cmdfixauth

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

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunFixAuth parses flags and orchestrates the auth-fix pipeline.
func RunFixAuth(args []string) error {
	user, email, assumeYes, force := parseFixAuthFlags(args)
	userErr := checkUserFlag(user)
	if userErr != nil {
		return userErr
	}

	repoErr := checkGitRepo()
	if repoErr != nil {
		return repoErr
	}

	return executeFixAuthPipeline(
		user,
		email,
		assumeYes,
		force,
	)
}

func checkUserFlag(user string) *apperror.AppError {
	if user == "" {
		fmt.Fprintln(os.Stderr, "✗ --user <github-username> is required")
		fmt.Fprintln(os.Stderr, "  example: gitmap fix-auth --user aukgit --email me@example.com")

		return apperror.NewSimple("--user <github-username> is required", "E_FIXAUTH_USER")
	}

	return nil
}

func checkGitRepo() *apperror.AppError {
	isGitRepo := cmdpull.IsGitRepoCWD()
	if !isGitRepo {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run inside the repo you want to fix)")

		return apperror.NewSimple("not a git repository (run inside the repo you want to fix)", "E_NOT_GIT_REPO")
	}

	return nil
}

func executeFixAuthPipeline(
	user string,
	email string,
	assumeYes bool,
	force bool,
) error {
	keyPath := fixAuthKeyPath(user)
	err := setupSSHDirAndKey(
		keyPath,
		email,
		assumeYes,
		force,
	)
	if err != nil {
		return err
	}

	return bindAndAnnounceAuth(user, keyPath)
}

func setupSSHDirAndKey(
	keyPath string,
	email string,
	assumeYes bool,
	force bool,
) error {
	dirErr := cmdssh.EnsureSSHDir(filepath.Dir(keyPath))
	if dirErr != nil {
		return apperror.WrapSimple(dirErr, "mkdir ~/.ssh")
	}

	return fixAuthGenerate(
		keyPath,
		resolveFixAuthEmail(email),
		assumeYes,
		force,
	)
}

func bindAndAnnounceAuth(user, keyPath string) error {
	bindErr := fixAuthBind(keyPath)
	if bindErr != nil {
		return bindErr
	}

	annErr := fixAuthAnnounce(user, keyPath)
	if annErr != nil {
		return annErr
	}

	_, _ = cmdssh.SanitizeSSHConfigFile("")

	return nil
}

func setupFixAuthFlags(fs *flag.FlagSet) (*string, *string, *bool, *bool) {
	user := fs.String("user", "", "GitHub username (required)")
	fs.StringVar(user, "u", "", "GitHub username (short)")
	email := fs.String("email", "", "Email comment")
	fs.StringVar(email, "e", "", "Email (short)")
	yes := fs.Bool("yes", false, "Skip confirmation prompts")
	fs.BoolVar(yes, "y", false, "Yes (short)")
	force := fs.Bool("force", false, "Overwrite existing key")
	fs.BoolVar(force, "f", false, "Force (short)")

	return user, email, yes, force
}

// parseFixAuthFlags parses fix-auth CLI flags with short aliases.
func parseFixAuthFlags(args []string) (string, string, bool, bool) {
	fs := flag.NewFlagSet(constants.CmdFixAuth, flag.ExitOnError)
	u, e, y, f := setupFixAuthFlags(fs)
	fs.Parse(args)
	userVal := *u
	if userVal == "" && fs.NArg() > 0 {
		userVal = fs.Arg(0)
	}

	return userVal, *e, *y, *f
}

// resolveFixAuthEmail falls back to git global user.email.
func resolveFixAuthEmail(email string) string {
	if email != "" {
		return email
	}

	resolved := cmdssh.ResolveGitEmail()
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
func fixAuthGenerate(
	keyPath string,
	email string,
	assumeYes bool,
	force bool,
) error {
	_, err := os.Stat(keyPath)
	hasKey := err == nil
	isValid := validateKeyGeneration(
		hasKey,
		keyPath,
		assumeYes,
		force,
	)
	if !isValid {
		return nil
	}

	valErr := cmdssh.ValidateSSHKeygen()
	if valErr != nil {
		return apperror.WrapSimple(valErr, constants.ErrSSHKeygenMissing)
	}

	return runSSHKeygenEd25519(keyPath, email)
}

func validateKeyGeneration(
	hasKey bool,
	keyPath string,
	assumeYes bool,
	force bool,
) bool {
	if !hasKey {
		return true
	}

	if !force {
		fmt.Printf("• key already exists, reusing: %s\n", keyPath)

		return false
	}

	return confirmAndRemoveOldKey(keyPath, assumeYes)
}

func confirmAndRemoveOldKey(keyPath string, assumeYes bool) bool {
	isConfirmed := assumeYes || confirmOverwrite(keyPath)
	if !isConfirmed {
		fmt.Println("• aborted; existing key kept")

		return false
	}

	_ = os.Remove(keyPath)
	_ = os.Remove(keyPath + ".pub")

	return true
}

// runSSHKeygenEd25519 invokes ssh-keygen with an empty passphrase.
func runSSHKeygenEd25519(keyPath, email string) error {
	cmd := exec.Command(constants.SSHKeygenBin,
		"-t", "ed25519", "-C", email, "-f", keyPath, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "ssh-keygen")
	}

	return nil
}

// confirmOverwrite prompts y/N before clobbering an existing key.
func confirmOverwrite(keyPath string) bool {
	fmt.Printf("⚠ overwrite existing key %s? [y/N]: ", keyPath)
	var ans string
	_, _ = fmt.Scanln(&ans)

	return strings.EqualFold(strings.TrimSpace(ans), "y")
}

// fixAuthBind pins the given key to the cwd repo's core.sshCommand.
func fixAuthBind(keyPath string) error {
	ref := sshBindRefForKey(keyPath)
	cmdStr := fmt.Sprintf("ssh -i %s -F /dev/null -o IdentitiesOnly=yes", ref)
	out, err := exec.Command("git", "config", "core.sshCommand", cmdStr).CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("git config core.sshCommand: %s", out))
	}

	fmt.Printf("✓ pinned repo → %s\n", cmdStr)

	return nil
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
func fixAuthAnnounce(user, keyPath string) error {
	pub, err := os.ReadFile(keyPath + ".pub")
	if err != nil {
		return apperror.WrapSimple(err, "read public key")
	}

	printPublicKeyBanner(user, strings.TrimSpace(string(pub)))

	return nil
}

func printPublicKeyBanner(user, trimmed string) {
	fmt.Println("\n=== PUBLIC KEY (add this to GitHub) ===")
	fmt.Println(trimmed)
	fmt.Println("=======================================")
	cmdssh.CopyPubKeyAndAnnounce(trimmed)
	fmt.Printf("\nNext steps for %s:\n", user)
	fmt.Println("  1. Open https://github.com/settings/ssh/new")
	fmt.Println("  2. Paste the key above, save.")
	fmt.Println("  3. Run:  git push")
}
