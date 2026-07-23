package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/movemerge"
)

// mergeSpec describes one merge variant for the dispatcher.
type mergeSpec struct {
	cmd       string
	prefix    string
	msgFmt    string
	direction movemerge.Direction
}

// runMerge implements merge-both / merge-left / merge-right.
//
// Spec: spec/01-app/97-move-and-merge.md
func runMerge(spec mergeSpec, args []string) {
	checkHelp(spec.cmd, args)
	left, right, opts := parseMergeArgs(spec, args)
	leftEP := mustResolve(left, true, opts)
	rightEP := mustResolve(right, false, opts)
	logResolved(leftEP, rightEP, opts)
	if err := movemerge.RunMerge(leftEP, rightEP, spec.direction, opts); err != nil {
		cliexit.Fail(spec.cmd, "merge", leftEP.DisplayName+" <-> "+rightEP.DisplayName, err, 1)
	}
}

// parseMergeArgs parses positional + flag arguments for any merge-*.
func parseMergeArgs(spec mergeSpec, args []string) (string, string, movemerge.Options) {
	fs := flag.NewFlagSet(spec.cmd, flag.ExitOnError)
	mf := &movemergeFlagSet{}
	mf.bindFlags(fs)
	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(2)
	}
	rest := fs.Args()
	if len(rest) != 2 {
		fmt.Fprintf(os.Stderr, constants.ErrMMUsageFmt, spec.cmd)
		os.Exit(2)
	}
	opts := mf.toOptions(spec.cmd, spec.prefix, spec.msgFmt)

	return rest[0], rest[1], opts
}
