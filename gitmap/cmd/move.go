package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/movemerge"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/txn"
)

// runMove implements `gitmap mv <source> <dest>`.
func runMove(args []string) error {
	checkHelp(constants.CmdMv, args)
	mOpts, positional := parseMoveFlags(args)
	if len(positional) == constants.ExpectedMoveArgsCount && handleRepoMove(positional[0], positional[1], mOpts) {
		return nil
	}
	runMoveMerge(args)
	return nil
}

func handleRepoMove(srcTarget, destTarget string, opts moveOpts) bool {
	db, err := openDB()
	if err != nil {
		return false
	}
	defer db.Close()
	rec, destPath, ok := prepareRepoMove(db, srcTarget, destTarget)
	if !ok {
		return false
	}
	return executeRepoMove(db, *rec, destPath, opts)
}

func prepareRepoMove(db *store.DB, src, dest string) (*model.ScanRecord, string, bool) {
	rec, err := ResolveRepo(db, src)
	if err != nil || rec == nil {
		return nil, "", false
	}
	destPath, err := calculateDestPath(rec.AbsolutePath, dest)
	if err != nil || preflightMove(rec.AbsolutePath, destPath) != nil {
		return nil, "", false
	}
	return rec, destPath, true
}

func executeRepoMove(db *store.DB, rec model.ScanRecord, destPath string, opts moveOpts) bool {
	if opts.dryRun {
		printMoveDryRun(rec.AbsolutePath, destPath)
		return true
	}
	if opts.yes || confirmMovePrompt(rec.Slug, rec.AbsolutePath, destPath) {
		executeMove(db, rec, destPath, opts)
		return true
	}
	fmt.Println(constants.MsgMoveAborted)
	return true
}

func runMoveMerge(args []string) error {
	left, right, opts := parseMoveArgs(args)
	leftEP := mustResolve(left, true, opts)
	rightEP := mustResolve(right, false, opts)
	logResolved(leftEP, rightEP, opts)
	j := beginMoveTxn(leftEP, rightEP)
	if err := movemerge.RunMove(leftEP, rightEP, opts); err != nil {
		_ = j.Abort()
		cliexit.Fail(constants.CmdMv, constants.OpMove, leftEP.DisplayName+" -> "+rightEP.DisplayName, err, constants.ExitCodeError)
	}
	finalizeMoveTxn(j, leftEP, rightEP)
	return nil
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
	return createMoveTxnJournal(db, left, right)
}

func createMoveTxnJournal(db *store.DB, left, right movemerge.Endpoint) *txn.Journal {
	cwd, _ := os.Getwd()
	j, _ := txn.Begin(db, txn.Meta{
		Kind:           constants.TxnKindMv,
		Argv:           os.Args,
		Cwd:            cwd,
		ReverseSummary: fmt.Sprintf(constants.TxnSummaryRenameFmt, left.WorkingDir, right.WorkingDir),
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
	fs, mf := newMoveFlagSet()
	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(constants.ExitCodeUsage)
	}
	rest := fs.Args()
	if len(rest) != constants.ExpectedMoveArgsCount {
		fmt.Fprintf(os.Stderr, constants.ErrMMUsageFmt, constants.CmdMv)
		os.Exit(constants.ExitCodeUsage)
	}
	opts := mf.toOptions(constants.CmdMv, constants.LogPrefixMv, constants.CommitMsgMv)

	return rest[0], rest[1], opts
}

func newMoveFlagSet() (*flag.FlagSet, *movemergeFlagSet) {
	fs := flag.NewFlagSet(constants.CmdMv, flag.ExitOnError)
	mf := &movemergeFlagSet{}
	mf.bindFlags(fs)
	return fs, mf
}

func mustResolve(raw string, isLeft bool, opts movemerge.Options) movemerge.Endpoint {
	ep, err := movemerge.ResolveEndpoint(raw, isLeft, opts)
	if err != nil {
		if db, dbErr := openDB(); dbErr == nil {
			PrintRepoSuggestions(db, raw)
			db.Close()
		}
		cliexit.Fail(constants.CmdMv, constants.OpResolveEndpoint, raw, err, constants.ExitCodeError)
	}

	return ep
}

// logResolved emits the [cmd] resolving LEFT/RIGHT lines.
func logResolved(l, r movemerge.Endpoint, opts movemerge.Options) {
	fmt.Printf(constants.LogResolvedLeftFmt, opts.LogPrefix, l.DisplayName)
	fmt.Printf(constants.LogResolvedRightFmt, opts.LogPrefix, r.DisplayName)
}
