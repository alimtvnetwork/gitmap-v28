package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/movemerge"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/txn"
)

// runMove implements `gitmap mv LEFT RIGHT`.
//
// Spec: spec/01-app/97-move-and-merge.md
func runMove(args []string) {
	checkHelp(constants.CmdMv, args)
	left, right, opts := parseMoveArgs(args)
	leftEP := mustResolve(left, true, opts)
	rightEP := mustResolve(right, false, opts)
	logResolved(leftEP, rightEP, opts)
	j := beginMoveTxn(leftEP, rightEP)
	if err := movemerge.RunMove(leftEP, rightEP, opts); err != nil {
		_ = j.Abort()
		cliexit.Fail(constants.CmdMv, "move", leftEP.DisplayName+" -> "+rightEP.DisplayName, err, 1)
	}
	finalizeMoveTxn(j, leftEP, rightEP)
}

// beginMoveTxn opens a journal row for a folder→folder mv. Returns a no-op
// journal (id == 0) when either endpoint is a remote URL or when the db is
// unavailable — the move itself must never be blocked by journaling.
func beginMoveTxn(left, right movemerge.Endpoint) *txn.Journal {
	if left.Kind != movemerge.EndpointFolder || right.Kind != movemerge.EndpointFolder {
		return &txn.Journal{}
	}
	db, err := openDB()
	if err != nil {
		return &txn.Journal{}
	}
	cwd, _ := os.Getwd()
	j, _ := txn.Begin(db, txn.Meta{
		Kind:           constants.TxnKindMv,
		Argv:           os.Args,
		Cwd:            cwd,
		ReverseSummary: fmt.Sprintf("rename %q ← %q", left.WorkingDir, right.WorkingDir),
	})

	return j
}

// finalizeMoveTxn records the rename inverse and commits the journal row.
func finalizeMoveTxn(j *txn.Journal, left, right movemerge.Endpoint) {
	if j.ID() == 0 {
		return
	}
	_ = j.RecordRename(left.WorkingDir, right.WorkingDir)
	_ = j.Commit()
}

// parseMoveArgs parses positional + flag arguments for mv.
func parseMoveArgs(args []string) (string, string, movemerge.Options) {
	fs := flag.NewFlagSet(constants.CmdMv, flag.ExitOnError)
	mf := &movemergeFlagSet{}
	mf.bindFlags(fs)
	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(2)
	}
	rest := fs.Args()
	if len(rest) != 2 {
		fmt.Fprintf(os.Stderr, constants.ErrMMUsageFmt, constants.CmdMv)
		os.Exit(2)
	}
	opts := mf.toOptions(constants.CmdMv, constants.LogPrefixMv, constants.CommitMsgMv)

	return rest[0], rest[1], opts
}

// mustResolve resolves an endpoint or exits with code 1 on failure.
func mustResolve(raw string, isLeft bool, opts movemerge.Options) movemerge.Endpoint {
	ep, err := movemerge.ResolveEndpoint(raw, isLeft, opts)
	if err != nil {
		cliexit.Fail(constants.CmdMv, "resolve-endpoint", raw, err, 1)
	}

	return ep
}

// logResolved emits the [cmd] resolving LEFT/RIGHT lines.
func logResolved(l, r movemerge.Endpoint, opts movemerge.Options) {
	fmt.Printf("%s resolving LEFT  : %s\n", opts.LogPrefix, l.DisplayName)
	fmt.Printf("%s resolving RIGHT : %s\n", opts.LogPrefix, r.DisplayName)
}
