package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/txn"
)

// handleRevertTxnFlags dispatches the transaction-journal sub-modes of
// `gitmap revert` (--list-txn / --show-txn / --txn / --last-txn /
// --prune-txn). Returns true when one of those flags was handled, so the
// legacy version-tag path in runRevert is skipped.
func handleRevertTxnFlags(args []string) bool {
	if hasRevertFlag(args, constants.FlagRevertListTxn) {
		runListTxn()

		return true
	}
	if id, ok := flagValue(args, constants.FlagRevertShowTxn); ok {
		runShowTxn(id)

		return true
	}
	if hasRevertFlag(args, constants.FlagRevertPruneTxn) {
		runPruneTxn()

		return true
	}

	return dispatchRevertTxn(args)
}

// dispatchRevertTxn handles --txn <id> and --last-txn separately so the
// outer dispatcher stays under the line cap.
func dispatchRevertTxn(args []string) bool {
	if id, ok := flagValue(args, constants.FlagRevertTxn); ok {
		runRevertTxn(id, hasRevertFlag(args, constants.FlagRevertForce))

		return true
	}
	if hasRevertFlag(args, constants.FlagRevertLastTxn) {
		runRevertLastTxn(hasRevertFlag(args, constants.FlagRevertForce))

		return true
	}
	if raw, ok := flagValue(args, constants.FlagRevertLastN); ok {
		runRevertLastN(raw, hasRevertFlag(args, constants.FlagRevertForce))

		return true
	}

	return false
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
		if v, ok := strings.CutPrefix(a, prefix+"="); ok {
			return v, true
		}
		if a == prefix && i+1 < len(args) {
			return args[i+1], true
		}
	}

	return "", false
}

// runListTxn prints the most recent transactions and exits.
func runListTxn() {
	db := mustOpenForTxn()
	defer db.Close()
	rows, err := db.ListTransactions(constants.TxnRetentionCap)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTxnDBWrite, err)
		os.Exit(1)
	}
	printTxnRows(rows)
}

// printTxnRows renders one TransactionRecord per line.
func printTxnRows(rows []model.TransactionRecord) {
	if len(rows) == 0 {
		fmt.Print(constants.MsgTxnNoCommitted)

		return
	}
	for _, r := range rows {
		fmt.Printf("#%-4d %-10s %-9s %s  %s\n",
			r.ID, r.Kind, r.Status,
			time.Unix(r.CreatedAt, 0).Format(time.RFC3339),
			r.ReverseSummary)
	}
}

// runShowTxn prints one transaction (header + every captured file) by id.
func runShowTxn(raw string) {
	id := mustParseTxnID(raw)
	db := mustOpenForTxn()
	defer db.Close()
	rec, err := db.FindTransactionByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTxnRowNotFound, id)
		os.Exit(1)
	}
	files, err := db.ListTransactionFiles(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTxnDBWrite, err)
		os.Exit(1)
	}
	printTxnDetail(rec, files)
}

// printTxnDetail renders one record + its file list.
func printTxnDetail(r model.TransactionRecord, files []model.TransactionFileRecord) {
	fmt.Printf("Transaction #%d\n  kind   : %s\n  status : %s\n  argv   : %s\n  cwd    : %s\n",
		r.ID, r.Kind, r.Status, r.Argv, r.Cwd)
	fmt.Printf("  reverse: %s\n  files  : %d\n", r.ReverseSummary, len(files))
	for _, f := range files {
		fmt.Printf("    [%s] %s\n", f.Action, f.AbsPath)
	}
}

// runRevertTxn reverts the named transaction id.
func runRevertTxn(raw string, force bool) {
	id := mustParseTxnID(raw)
	revertOne(id, force)
}

// runRevertLastTxn reverts the newest committed transaction.
func runRevertLastTxn(force bool) {
	db := mustOpenForTxn()
	id, err := db.LastCommittedTransactionID()
	db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTxnDBWrite, err)
		os.Exit(1)
	}
	if id == 0 {
		fmt.Print(constants.MsgTxnNoCommitted)

		return
	}
	revertOne(id, force)
}

// revertOne is the shared confirm + apply path for --txn / --last-txn.
func revertOne(id int64, force bool) {
	db := mustOpenForTxn()
	defer db.Close()
	rec, err := db.FindTransactionByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTxnRowNotFound, id)
		os.Exit(1)
	}
	files, _ := db.ListTransactionFiles(id)
	if !force && !confirmRevert(rec, len(files)) {
		fmt.Print(constants.MsgTxnAbortedByUser)

		return
	}
	applyRevertOrExit(db, id, rec, force)
}

// applyRevertOrExit calls txn.Revert and prints the success line.
func applyRevertOrExit(db *store.DB, id int64, rec model.TransactionRecord, force bool) {
	if err := txn.Revert(db, id, txn.RevertOptions{Force: force}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Printf(constants.MsgTxnReverted, id, rec.Kind)
}

// confirmRevert prompts the user; returns true on a "yes" answer.
func confirmRevert(r model.TransactionRecord, fileCount int) bool {
	fmt.Printf(constants.MsgTxnConfirmRevert, r.ID, r.Kind, fileCount)
	fmt.Print(constants.MsgTxnConfirmPrompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	return strings.TrimSpace(strings.ToLower(scanner.Text())) == "yes"
}

// runPruneTxn forces an immediate prune cycle.
func runPruneTxn() {
	db := mustOpenForTxn()
	defer db.Close()
	dropped, err := db.PruneOldestTransactions(constants.TxnRetentionCap)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTxnDBWrite, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgTxnPruned, len(dropped))
}

// mustOpenForTxn opens the gitmap database or exits.
func mustOpenForTxn() *store.DB {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTxnDBOpen, err)
		os.Exit(1)
	}

	return db
}

// mustParseTxnID parses a positive int64 transaction id or exits with usage.
func mustParseTxnID(raw string) int64 {
	id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || id <= 0 {
		fmt.Fprintf(os.Stderr, "revert: invalid transaction id %q\n", raw)
		os.Exit(2)
	}

	return id
}
