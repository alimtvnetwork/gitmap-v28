package clonefrom

// execute_checkout.go — post-clone checkout helpers. Lives in its
// own file so execute.go stays under the 200-line cap and the
// checkout semantics have one obvious home.
//
// Three modes (resolved by EffectiveCheckout from row.Checkout +
// the global default constants.CloneFromCheckoutDefault):
//
//   - "auto"  → no-op. `git clone` (with --branch when set) already
//               materialized the working tree. This is the legacy
//               pre-feature behavior, preserved byte-for-byte.
//   - "skip"  → no-op HERE; the work happened earlier in
//               buildGitArgs which appended --no-checkout. We do
//               NOT also run `git checkout` after — that would
//               defeat the point of skipping.
//   - "force" → if row.Branch is non-empty, run
//               `git -C <dest> checkout <branch>`. On failure
//               (typo'd branch, branch missing on remote, detached
//               HEAD with no target) the row is marked failed with
//               a clear MsgCloneFromBranchMissingFmt detail and a
//               Code Red stderr log, mirroring how MkdirAll
//               failures are reported by execute_dest.go.
//               When row.Branch is empty, force is a no-op (git
//               already checked out the remote's default HEAD).

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// EffectiveCheckout resolves a row's Checkout field to one of the
// three concrete modes. Empty → CloneFromCheckoutDefault ("auto").
// Exported so renderers + faithful-cmd verifiers agree with the
// executor on the resolved mode without re-importing this file.
func EffectiveCheckout(r Row) string {
	if len(r.Checkout) == 0 {
		return constants.CloneFromCheckoutDefault
	}

	return r.Checkout
}

// runPostCloneCheckout performs the post-clone checkout step when
// the resolved mode is "force". Returns (detail, ok) following the
// same convention as runGitClone: empty detail on success, single-
// line summary on failure. Auto/skip modes return ("", true)
// immediately — the executor's hot path stays branchless for them.
func runPostCloneCheckout(r Row, dest, cwd string) (string, bool) {
	if EffectiveCheckout(r) != constants.CloneFromCheckoutForce {
		return "", true
	}
	if len(r.Branch) == 0 {
		// force on a default-HEAD row: git already checked out HEAD.
		return "", true
	}

	return runGitCheckout(r.Branch, resolvePostCheckoutDest(dest, cwd))
}

// resolvePostCheckoutDest computes the absolute path to hand `git
// -C` so the checkout runs inside the freshly-cloned repo regardless
// of whether row.Dest was relative or absolute. Mirrors the join
// resolveDest does for the skip-check.
func resolvePostCheckoutDest(dest, cwd string) string {
	if filepath.IsAbs(dest) {
		return dest
	}

	return filepath.Join(cwd, dest)
}

// runGitCheckout shells out to `git -C <dir> checkout <branch>` and
// translates the result. On failure we emit the project's standard
// Code Red stderr log AND return a short detail string so the per-
// row line + CSV report carry the same diagnosis. Failure cause is
// almost always "branch missing on remote" (the row pinned a typo
// or a deleted branch); we surface that wording in the detail to
// save the user a trip to the git stderr scrollback.
func runGitCheckout(branch, absDest string) (string, bool) {
	cmd := exec.Command(constants.GitBin, "checkout", branch)
	cmd.Dir = absDest
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFromCheckoutFailed,
			absDest, branch, err)
		_ = out

		return fmt.Sprintf(constants.MsgCloneFromBranchMissingFmt, branch), false
	}

	return "", true
}
