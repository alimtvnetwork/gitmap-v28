package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// escapeNestedGitRepo walks up from the current working directory
// while each ancestor is itself a git repo, chdir'ing into the first
// non-repo ancestor before the clone step runs. This prevents the
// failure mode reported against v6.47.0 where running
// `gitmap cfrp <url>` from inside a working repo nested the freshly
// cloned tree under another repo's git context and aborted with
// `fetch-pack: invalid index-pack output`.
//
// Stops as soon as an ancestor is reached that is NOT a git repo
// (matches the user's stated rule: "go to the parent and see that
// the parent is a Git repo. If not, then it will run that command").
// Bounded to 32 hops as a defensive cap against filesystem loops.
func escapeNestedGitRepo() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	if !clonenext.IsGitRepo(cwd) {
		return
	}
	origin := cwd
	target := cwd
	for hop := 0; hop < 32; hop++ {
		parent := filepath.Dir(target)
		if parent == target {
			break
		}
		target = parent
		if !clonenext.IsGitRepo(target) {
			break
		}
	}
	if target == origin {
		return
	}
	if chErr := os.Chdir(target); chErr != nil {
		fmt.Fprintf(os.Stderr, constants.WarnCFREscapeChdir, target, chErr)
		return
	}
	fmt.Fprintf(os.Stderr, constants.MsgCFREscapeNested, origin, target)
}
