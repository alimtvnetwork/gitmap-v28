package cmdpull

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// AnnounceNonGitPullFallback prints notice that gitmap pull defaulted to pull all.
func AnnounceNonGitPullFallback() {
	arrow := resolveSubArrow()
	fmt.Printf("    %s%s%s %scwd is not a git repo — defaulting to gitmap pull all%s\n\n",
		constants.ColorCyan, arrow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
}

// ShouldFallbackToPullAll checks if pull should default to multi-repo pull all.
func ShouldFallbackToPullAll(opts pullOptions) bool {
	if opts.slug != "" || opts.group != "" || opts.all || HasAlias() {
		return false
	}

	return !isGitRepoCWD()
}
