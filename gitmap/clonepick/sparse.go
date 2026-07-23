package clonepick

// sparse.go: the actual git-clone + sparse-checkout pipeline.
//
// Each step shells out to git (no go-git dependency) so we get the
// exact behavior users see when they run the same commands by hand.
// Errors are wrapped with constants.ErrClonePickGit* so the cmd-layer
// can render a single-line failure message without exposing exec
// internals.

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runSparseCheckout executes the five-step pipeline documented in the
// spec (clone --no-checkout, sparse-checkout init, set, checkout,
// optionally remove .git). Returns the absolute destination path on
// success so the cmd layer can call WriteShellHandoff with it.
//
// When plan.PreClonedSrc is set, the clone step is skipped and the
// pre-cloned metadata dir is promoted into dest instead -- this is
// the "clone-once" path used by --ask so we don't pay the network
// cost twice.
func runSparseCheckout(plan Plan, progress io.Writer) (string, error) {
	dest, err := prepareDest(plan)
	if err != nil {
		return "", err
	}
	if err := acquireRepoTree(plan, dest, progress); err != nil {
		return dest, err
	}
	if err := gitSparseInit(plan, dest, progress); err != nil {
		return dest, fmt.Errorf(constants.ErrClonePickGitSparseInit, err)
	}
	if err := gitSparseSet(plan, dest, progress); err != nil {
		return dest, fmt.Errorf(constants.ErrClonePickGitSparseSet, err)
	}
	if err := gitCheckout(dest, progress); err != nil {
		return dest, fmt.Errorf(constants.ErrClonePickGitCheckout, err)
	}
	if !plan.KeepGit {
		if err := os.RemoveAll(filepath.Join(dest, ".git")); err != nil {
			return dest, fmt.Errorf(constants.ErrClonePickFsRemoveDotGit, err)
		}
	}

	return dest, nil
}

// acquireRepoTree either promotes a pre-cloned metadata dir into
// dest (clone-once path) or runs a fresh `git clone --no-checkout`.
// Split out so runSparseCheckout stays under the function-length cap.
func acquireRepoTree(plan Plan, dest string, progress io.Writer) error {
	if len(plan.PreClonedSrc) > 0 {
		if err := promotePreClonedSrc(plan.PreClonedSrc, dest); err != nil {
			return fmt.Errorf(constants.ErrClonePickPromoteSrc, err)
		}

		return nil
	}
	if err := gitClonePartial(plan, dest, progress); err != nil {
		return fmt.Errorf(constants.ErrClonePickGitClone, err)
	}

	return nil
}

// prepareDest resolves DestDir to an absolute path and creates it if
// missing. Refuses to write into a non-empty directory unless --force
// was passed -- the partial clone would silently fail if `git clone`
// saw an existing tree.
func prepareDest(plan Plan) (string, error) {
	abs, err := filepath.Abs(plan.DestDir)
	if err != nil {
		return "", fmt.Errorf(constants.ErrClonePickFsCreateDest, err)
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return abs, fmt.Errorf(constants.ErrClonePickFsCreateDest, err)
	}
	if plan.Force {
		return abs, nil
	}
	if err := assertDestEmpty(abs); err != nil {
		return abs, err
	}

	return abs, nil
}

// assertDestEmpty returns an error when abs already contains entries.
// Hidden files count -- a stray .DS_Store still indicates an existing
// tree and we'd rather refuse than half-clone over it.
func assertDestEmpty(abs string) error {
	entries, err := os.ReadDir(abs)
	if err != nil {
		return fmt.Errorf(constants.ErrClonePickFsCreateDest, err)
	}
	if len(entries) > 0 {
		return fmt.Errorf("%s", constants.MsgClonePickDestDirty)
	}

	return nil
}

// gitClonePartial runs `git clone --filter=blob:none --no-checkout`
// with the optional --branch / --depth flags. Cloning into an
// existing empty dir requires the dir to actually exist, which
// prepareDest guarantees -- but `git clone <url> <abs>` itself
// expects the dir to NOT exist OR to be empty, both of which the
// caller guarantees.
func gitClonePartial(plan Plan, dest string, progress io.Writer) error {
	args := []string{"clone", "--filter=blob:none", "--no-checkout"}
	if len(plan.Branch) > 0 {
		args = append(args, "--branch", plan.Branch)
	}
	if plan.Depth > 0 {
		args = append(args, "--depth", strconv.Itoa(plan.Depth))
	}
	args = append(args, plan.RepoUrl, dest)

	return runGit(progress, "", args...)
}

// gitSparseInit enables sparse-checkout. Cone mode is the default --
// non-cone is opt-in via plan.Cone == false (which ParseArgs flips
// when the path list needs file-level patterns).
func gitSparseInit(plan Plan, dest string, progress io.Writer) error {
	args := []string{"sparse-checkout", "init"}
	if plan.Cone {
		args = append(args, "--cone")
	} else {
		args = append(args, "--no-cone")
	}

	return runGit(progress, dest, args...)
}

// gitSparseSet writes the path list into .git/info/sparse-checkout.
// In cone mode each entry is treated as a directory; in non-cone
// mode they're literal gitignore-style patterns.
func gitSparseSet(plan Plan, dest string, progress io.Writer) error {
	args := append([]string{"sparse-checkout", "set"}, plan.Paths...)

	return runGit(progress, dest, args...)
}

// gitCheckout materializes the working tree against the active
// sparse pattern. We use the default branch (whatever HEAD points to
// after the partial clone) -- plan.Branch was already honored by
// gitClonePartial via --branch.
func gitCheckout(dest string, progress io.Writer) error {
	return runGit(progress, dest, "checkout")
}

// runGit invokes git with the given args, streaming output to the
// progress writer. workdir == "" means the process cwd.
func runGit(progress io.Writer, workdir string, args ...string) error {
	cmd := exec.Command("git", args...)
	if len(workdir) > 0 {
		cmd.Dir = workdir
	}
	cmd.Stdout = progress
	cmd.Stderr = progress

	return cmd.Run()
}
