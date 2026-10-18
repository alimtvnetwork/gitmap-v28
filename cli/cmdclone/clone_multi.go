// Package cmdclone — clone_multi.go handles parallel/batch multi-repo cloning.
package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// isMultiCloneEnabled returns true when the positional arguments describe two
// or more repos to clone in a single invocation.
func isMultiCloneEnabled(cf CloneFlags) bool {
	for _, p := range cf.Positional {
		if strings.ContainsAny(p, urlListSeparators) {
			return true
		}
	}

	for i, p := range cf.Positional {
		if i >= 1 && isDirectURL(sanitizeURLToken(p)) {
			return true
		}
	}

	flat := []string{}
	if len(cf.Positional) >= 1 {
		flat = flattenURLArgs(cf.Positional[:1])
	}

	urlCount := 0
	for _, u := range flat {
		if isDirectURL(u) {
			urlCount++
		}
	}

	if urlCount >= 2 {
		return true
	}

	return false
}

// runCloneMulti clones every URL in the flattened positional list, continuing
// on per-URL failure. Folder name is ignored in this mode.
func runCloneMulti(cf CloneFlags) error {
	flat := flattenURLArgs(cf.Positional)
	urls, invalid := classifyURLs(flat)

	if len(urls) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrCloneAllInvalid)
		cliexit.HandleError(nil, constants.ExitCloneMultiAllInvalid)
	}

	fmt.Printf(constants.MsgCloneMultiBegin, len(urls))

	succeeded := 0
	failed := 0
	pmPairs := make([]vscodepm.Pair, 0, len(urls))

	for idx, url := range urls {
		fmt.Printf(constants.MsgCloneMultiItem, idx+1, len(urls), url)
		printCloneTermBlockForURL(cf.Output, idx+1, url, "")
		maybeCleanMultiFolder(url, cf.Clean)

		if err := executeDirectCloneOne(url, "", cf.GHDesktop, cf.NoReplace); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCloneMultiFailedFmt, idx+1, len(urls), url, err)
			failed++

			continue
		}

		succeeded++
		repoName := repoNameFromURL(url)
		folder := resolveCloneFolder(repoName, "")

		if abs, absErr := filepath.Abs(folder); absErr == nil {
			pmPairs = append(pmPairs, buildClonePMPair(abs, repoName))
		}
	}

	failed += len(invalid)

	fmt.Printf(constants.MsgCloneSummaryMultiFmt, succeeded, failed, len(urls)+len(invalid))
	syncClonedReposToVSCodePM(pmPairs, cf.NoVSCodeSync)

	if failed > 0 {
		cliexit.HandleError(nil, constants.ExitCloneMultiPartialFail)
	}

	return nil
}
