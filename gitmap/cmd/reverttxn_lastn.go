package cmd

// `gitmap revert --last-n-txn <N>`: replay the stored reverse-operation
// for the most recent N committed transactions, newest first. Each
// transaction is reverted via the same txn.Revert path used by
// --txn/--last-txn, so semantics (sha verification, --force bypass,
// row status flip) are identical.

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/txn"
)

// runRevertLastN parses N, confirms with the user (unless --force), and
// reverts the N most recent committed transactions newest-first.
func runRevertLastN(raw string, force bool) error {
	count := mustParseLastN(raw)
	rows := loadLastCommittedTxns(count)
	if len(rows) == 0 {
		fmt.Print(constants.MsgTxnLastNNoneFound)

		return nil
	}
	if !force && !confirmRevertLastN(rows) {
		fmt.Print(constants.MsgTxnAbortedByUser)

		return nil
	}
	revertManyOrExit(rows, force)
	return nil
}

// mustParseLastN validates the --last-n-txn count argument.
func mustParseLastN(raw string) int {
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n <= 0 {
		fmt.Fprintf(os.Stderr, constants.ErrRevertLastNBadCount, constants.FlagRevertLastN, raw)
		cliexit.HandleError(nil, 2)
	}

	return n
}

// loadLastCommittedTxns returns up to `want` newest-first committed
// transactions. We over-scan the retention window then filter so we
// never depend on a separate SQL query.
func loadLastCommittedTxns(want int) []model.TransactionRecord {
	db := mustOpenForTxn()
	defer db.Close()
	all, err := db.ListTransactions(constants.TxnRetentionCap)
	if err != nil {
		cliexit.HandleError(err, 1)
	}
	out := make([]model.TransactionRecord, 0, want)
	for _, r := range all {
		if r.Status != constants.TxnStatusCommitted {
			continue
		}
		out = append(out, r)
		if len(out) == want {
			break
		}
	}

	return out
}

// confirmRevertLastN previews the planned reverts and waits for `yes`.
func confirmRevertLastN(rows []model.TransactionRecord) bool {
	fmt.Printf(constants.MsgTxnLastNHeader, len(rows))
	for _, r := range rows {
		fmt.Printf(constants.MsgTxnLastNRow, r.ID, r.Kind,
			time.Unix(r.CreatedAt, 0).Format(time.RFC3339), r.ReverseSummary)
	}
	fmt.Print(constants.MsgTxnConfirmPrompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	return strings.TrimSpace(strings.ToLower(scanner.Text())) == "yes"
}

// revertManyOrExit applies txn.Revert to each row in newest-first order.
// Any single failure aborts immediately; already-reverted rows are
// preserved (we do NOT roll forward).
func revertManyOrExit(rows []model.TransactionRecord, force bool) {
	db := mustOpenForTxn()
	defer db.Close()
	for _, r := range rows {
		if err := txn.Revert(db, r.ID, txn.RevertOptions{Force: force}); err != nil {
			appErr := apperror.WrapWithDetails(
				err,
				"txn.revertMany",
				"E2015",
				err.Error(),
				"cmd.reverttxn_lastn",
				apperror.ErrorTypeExecution,
				apperror.SeverityError,
				map[string]any{"id": r.ID, "kind": r.Kind, "force": force},
			)
			cliexit.HandleError(appErr, 1)
		}
		fmt.Printf(constants.MsgTxnReverted, r.ID, r.Kind)
	}
	fmt.Printf(constants.MsgTxnLastNDone, len(rows))
}
