package cmdrevert

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrevert/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/txn"
)

// handleRevertTxnFlags dispatches the transaction-journal sub-modes of `gitmap revert`.
func handleRevertTxnFlags(args []string) (bool, error) {
	if hasRevertFlag(args, constants.FlagListTxn) {
		return true, runListTxn()
	}

	if id, hasVal := flagValue(args, constants.FlagShowTxn); hasVal {
		return true, runShowTxn(id)
	}

	if hasRevertFlag(args, constants.FlagPruneTxn) {
		return true, runPruneTxn()
	}

	return dispatchRevertTxn(args)
}

// dispatchRevertTxn handles --txn, --last-txn, and --last-n-txn.
func dispatchRevertTxn(args []string) (bool, error) {
	if id, hasVal := flagValue(args, constants.FlagTxn); hasVal {
		return true, runRevertTxn(id, hasRevertFlag(args, constants.FlagForce))
	}

	if hasRevertFlag(args, constants.FlagLastTxn) {
		return true, runRevertLastTxn(hasRevertFlag(args, constants.FlagForce))
	}

	if raw, hasVal := flagValue(args, constants.FlagLastNTxn); hasVal {
		return true, runRevertLastN(raw, hasRevertFlag(args, constants.FlagForce))
	}

	return false, nil
}

// hasRevertFlag returns true when --name appears anywhere in args.
func hasRevertFlag(args []string, name string) bool {
	target := "--" + name
	for _, a := range args {
		if a == target || strings.HasPrefix(a, target+"=") {
			return true
		}
	}

	return false
}

// flagValue extracts the value of --name <v> or --name=<v>.
func flagValue(args []string, name string) (string, bool) {
	prefix := "--" + name
	for i, a := range args {
		if v, hasCut := strings.CutPrefix(a, prefix+"="); hasCut {
			return v, true
		}

		if a == prefix && i+1 < len(args) {
			return args[i+1], true
		}
	}

	return "", false
}

// runListTxn prints the most recent transactions and exits.
func runListTxn() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.ListTransactions(constants.TxnRetentionCap)
	if err != nil {
		return err
	}

	printTxnRows(rows)

	return nil
}

// printTxnRows renders one TransactionRecord per line.
func printTxnRows(rows []model.TransactionRecord) {
	if len(rows) == 0 {
		fmt.Print(constants.MsgTxnNoCommitted)

		return
	}

	for _, r := range rows {
		fmt.Printf(constants.MsgTxnListRowFmt,
			r.ID, r.Kind, r.Status,
			time.Unix(r.CreatedAt, 0).Format(time.RFC3339),
			r.ReverseSummary)
	}
}

// runShowTxn prints one transaction (header + every captured file) by id.
func runShowTxn(raw string) error {
	id, err := parseTxnID(raw)
	if err != nil {
		return err
	}

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rec, err := db.FindTransactionByID(id)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrTxnRowNotFound)
	}

	files, err := db.ListTransactionFiles(id)
	if err != nil {
		return err
	}

	printTxnDetail(rec, files)

	return nil
}

// printTxnDetail renders one record + its file list.
func printTxnDetail(r model.TransactionRecord, files []model.TransactionFileRecord) {
	fmt.Printf(constants.MsgTxnDetailHeaderFmt,
		r.ID, r.Kind, r.Status, r.Argv, r.Cwd)
	fmt.Printf(constants.MsgTxnDetailReverseFmt, r.ReverseSummary, len(files))
	for _, f := range files {
		fmt.Printf(constants.MsgTxnDetailFileRowFmt, f.Action, f.AbsPath)
	}
}

// runRevertTxn reverts the named transaction id.
func runRevertTxn(raw string, isForce bool) error {
	id, err := parseTxnID(raw)
	if err != nil {
		return err
	}

	return revertOne(id, isForce)
}

// runRevertLastTxn reverts the newest committed transaction.
func runRevertLastTxn(isForce bool) error {
	db, err := openDB()
	if err != nil {
		return err
	}

	id, err := db.LastCommittedTransactionID()
	db.Close()
	if err != nil {
		return err
	}

	if id == 0 {
		fmt.Print(constants.MsgTxnNoCommitted)

		return nil
	}

	return revertOne(id, isForce)
}

// revertOne is the shared confirm + apply path for --txn / --last-txn.
func revertOne(id int64, isForce bool) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rec, err := db.FindTransactionByID(id)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrTxnRowNotFound)
	}

	files, _ := db.ListTransactionFiles(id)
	if !isForce && !hasUserConfirmedRevert(rec, len(files)) {
		fmt.Print(constants.MsgTxnAbortedByUser)

		return nil
	}

	return applyRevert(db, id, rec, isForce)
}

// applyRevert calls txn.Revert and prints the success line.
func applyRevert(db *store.DB, id int64, rec model.TransactionRecord, isForce bool) error {
	if err := txn.Revert(db, id, txn.RevertOptions{Force: isForce}); err != nil {
		return apperror.WrapWithDetails(
			err,
			constants.OpTxnRevert,
			constants.ErrCodeE2014,
			err.Error(),
			constants.DomainCmdRevertTxn,
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"id": id, "kind": rec.Kind, "force": isForce},
		)
	}

	fmt.Printf(constants.MsgTxnReverted, id, rec.Kind)

	return nil
}

// applyRevertOrExit maintains compatibility for external callers.
func applyRevertOrExit(db *store.DB, id int64, rec model.TransactionRecord, isForce bool) {
	if err := applyRevert(db, id, rec, isForce); err != nil {
		cliexit.HandleError(err, int(enums.ExitCodeGeneralError))
	}
}

// hasUserConfirmedRevert prompts the user; returns true on affirmative response.
func hasUserConfirmedRevert(r model.TransactionRecord, fileCount int) bool {
	fmt.Printf(constants.MsgTxnConfirmRevert, r.ID, r.Kind, fileCount)
	fmt.Print(constants.MsgTxnConfirmPrompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	return isAffirmativeResponse(scanner.Text())
}

// isAffirmativeResponse reports whether user input represents an affirmative answer.
func isAffirmativeResponse(text string) bool {
	resp, err := enums.ParseRevertPromptResponseType(text)
	if err != nil {
		return false
	}

	return resp.IsYes()
}

// runPruneTxn forces an immediate prune cycle.
func runPruneTxn() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	dropped, err := db.PruneOldestTransactions(constants.TxnRetentionCap)
	if err != nil {
		return err
	}

	fmt.Printf(constants.MsgTxnPruned, len(dropped))

	return nil
}

// mustOpenForTxn opens the gitmap database or exits with general error.
func mustOpenForTxn() *store.DB {
	db, err := openDB()
	if err != nil {
		cliexit.HandleError(err, int(enums.ExitCodeGeneralError))
	}

	return db
}

// parseTxnID parses a positive int64 transaction id or reports usage error.
func parseTxnID(raw string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrInvalidTxnIDFmt, raw, err)
		cliexit.HandleUsageError(err)

		return 0, err
	}
	if id <= 0 {
		fmt.Fprintf(os.Stderr, constants.ErrInvalidTxnIDNotPosFmt, raw)
		cliexit.HandleUsageError(nil)

		return 0, fmt.Errorf(constants.ErrInvalidTxnIDNotPosFmt, raw)
	}

	return id, nil
}

// mustParseTxnID parses a positive int64 transaction id or exits.
func mustParseTxnID(raw string) int64 {
	id, _ := parseTxnID(raw)

	return id
}
