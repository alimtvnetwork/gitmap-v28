// Package cmd — reporeclone.go implements the single-repo
// "wipe-and-re-clone" flow that overlays the existing manifest-based
// `gitmap reclone` command.
//
// Triggers when `reclone` / `rec` / `rc` / `relclone` / `clone-now`
// is invoked AND (the sole positional arg is a path containing
// `.git`) OR (no positional arg + cwd is inside a git repo). In
// every other shape (a manifest path, --manifest flag, or no repo
// in sight) we fall through to runCloneNow's manifest pipeline so
// existing scripts keep working byte-for-byte.
package cmdclone

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// tryRunRepoReclone owns the command and returns true when it
// handled (and exited / completed) the request. Returns false so
// the caller can continue with the manifest pipeline.

func tryRunRepoReclone(args []string) bool {
	yes, positionals := splitRepoRecloneArgs(args)
	target, ok := resolveRepoRecloneTarget(positionals)
	if !ok {
		return false
	}

	if len(positionals) == 1 && isGitURL(positionals[0]) {
		_ = RunRepoRecloneURL(target, positionals[0], yes)
		return true
	}

	RunRepoReclone(target, yes)

	return true
}

// splitRepoRecloneArgs separates the -y / --y / -yes flag from
// positionals without disturbing the broader flagset. Anything we
// don't recognize is forwarded as a positional so the manifest
// pipeline can still flag it.

func splitRepoRecloneArgs(args []string) (bool, []string) {
	var yes bool
	var positionals []string
	for _, a := range args {
		switch a {
		case "-" + constants.FlagRepoRecloneYes,
			"--" + constants.FlagRepoRecloneYes,
			"-yes", "--yes":
			yes = true
		default:
			positionals = append(positionals, a)
		}
	}

	return yes, positionals
}

// resolveRepoRecloneTarget returns an absolute path that is a git
// repo or target path for a repository URL, or ok=false if this invocation
// doesn't match the single-repo shape.

func resolveRepoRecloneTarget(positionals []string) (string, bool) {
	if len(positionals) > 1 {
		return "", false
	}

	if len(positionals) == 1 {
		return resolveRepoFromArg(positionals[0])
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	isNonGitRepoDir := !isGitRepoDir(cwd)
	if isNonGitRepoDir {
		return "", false
	}

	return cwd, true
}

func isGitURL(arg string) bool {
	clean := strings.ToLower(strings.TrimSpace(arg))
	return strings.HasPrefix(clean, "http://") ||
		strings.HasPrefix(clean, "https://") ||
		strings.HasPrefix(clean, "git@") ||
		strings.HasPrefix(clean, "ssh://") ||
		strings.HasSuffix(clean, ".git")
}

func cleanRepoURLName(rawURL string) string {
	clean := strings.TrimSpace(rawURL)
	clean = strings.TrimRight(clean, "/\\")
	return extractRepoName(clean)
}

func resolveRepoFromArg(arg string) (string, bool) {
	if isGitURL(arg) {
		repoName := cleanRepoURLName(arg)
		if repoName == "" {
			return "", false
		}
		cwd, err := os.Getwd()
		if err != nil {
			return "", false
		}
		if filepath.Base(cwd) == repoName && isGitRepoDir(cwd) {
			return cwd, true
		}
		dest := filepath.Join(cwd, repoName)
		return dest, true
	}

	abs, err := filepath.Abs(arg)
	if err != nil {
		return "", false
	}

	info, statErr := os.Stat(abs)
	if statErr != nil || !info.IsDir() {
		return "", false
	}

	isNonGitRepoDir := !isGitRepoDir(abs)
	if isNonGitRepoDir {
		return "", false
	}

	return abs, true
}

func isGitRepoDir(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))

	return err == nil
}

// RunRepoRecloneURL executes re-cloning from an explicit Git repository URL.
// If the destination directory exists, it confirms and wipes it before cloning;
// if it does not exist, it runs clone directly.
func RunRepoRecloneURL(target, origin string, yes bool) error {
	info, statErr := os.Stat(target)
	folderExists := statErr == nil && info.IsDir()

	if folderExists {
		parent := filepath.Dir(target)
		fmt.Printf(constants.MsgRepoReclonePlan, target, origin, parent)

		if !yes && !confirmRepoReclone(target, origin) {
			fmt.Fprint(os.Stderr, constants.MsgRepoRecloneAborted)

			return apperror.NewSimple("reclone aborted", "E9000")
		}

		if _, escapeErr := escapeCwdIfInside(target); escapeErr != nil {
			fmt.Fprintln(os.Stderr, escapeErr.Error())

			return apperror.WrapSimple(escapeErr, "escape cwd")
		}

		fmt.Printf(constants.MsgRepoRecloneRemoving, target)
		if rmErr := os.RemoveAll(target); rmErr != nil {
			return apperror.WrapSimple(rmErr, constants.ErrRepoRecloneRemove)
		}
	}

	dest := target
	origin = coerceURLToStoredTransport(origin)
	fmt.Printf(constants.MsgRepoRecloneCloning, origin, dest)
	params := DirectCloneParams{URL: origin, NoReplace: true}
	if _, cloneErr := runDirectCloneWithProfileFallback(params, origin, dest); cloneErr != nil {
		return apperror.WrapSimple(cloneErr, constants.ErrRepoRecloneClone)
	}

	persistRecloneTransport(origin)

	fmt.Printf(constants.MsgRepoRecloneDone, dest)
	WriteShellHandoff(dest)

	return nil
}

// RunRepoReclone executes the destructive wipe + re-clone after
// confirmation. Each failure path emits a standardized stderr line
// per the zero-swallow error policy and exits non-zero.

func RunRepoReclone(target string, yes bool) error {
	origin, err := currentOriginURL(target)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrRepoRecloneNoOrigin)
	}
	if origin == "" {
		return apperror.NewSimple(constants.ErrRepoRecloneNoOrigin, "E9000")
	}

	return RunRepoRecloneURL(target, origin, yes)
}

// confirmRepoReclone reads y/N from stdin. Non-TTY callers get a
// hard refusal pointing at -y; this prevents `yes | gitmap rc` from
// silently nuking the wrong tree in a CI job.

func confirmRepoReclone(target, origin string) bool {
	isNonStdinTTY := !isStdinTTY()
	if isNonStdinTTY {
		fmt.Fprint(os.Stderr, constants.ErrRepoRecloneNonTTY)

		return false
	}

	fmt.Printf(constants.MsgRepoRecloneConfirm, target, origin)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false
	}

	return strings.EqualFold(strings.TrimSpace(line), "y")
}

func isStdinTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

// parseRepoRecloneFromFlagSet is a thin wrapper used by tests +
// future direct callers that want to register -y on a parent
// flagset rather than rely on splitRepoRecloneArgs.
//
//nolint:unused
func parseRepoRecloneFromFlagSet(fs *flag.FlagSet) *bool {
	return fs.Bool(constants.FlagRepoRecloneYes, false,
		constants.FlagDescRepoRecloneYes)
}
