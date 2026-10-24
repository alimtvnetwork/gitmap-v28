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
	"github.com/alimtvnetwork/gitmap-v28/cli/txn"
)

// runRevertLastN parses N, confirms with the user, and reverts the N newest committed transactions.
func runRevertLastN(raw string, isForce bool) error {
	count, err := parseLastN(raw)
	if err != nil {
		return err
	}

	rows, err := loadLastCommittedTxns(count)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		fmt.Print(constants.MsgTxnLastNNoneFound)

		return nil
	}

	if !isForce && !isLastNConfirmed(rows) {
		fmt.Print(constants.MsgTxnAbortedByUser)

		return nil
	}

	return revertMany(rows, isForce)
}

// parseLastN validates the --last-n-txn count argument.
func parseLastN(raw string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRevertLastNBadCount, constants.FlagLastNTxn, raw)
		cliexit.HandleUsageError(err)

		return 0, err
	}
	if n <= 0 {
		fmt.Fprintf(os.Stderr, constants.ErrRevertLastNBadCount, constants.FlagLastNTxn, raw)
		cliexit.HandleUsageError(nil)

		return 0, fmt.Errorf(constants.ErrRevertLastNBadCount, constants.FlagLastNTxn, raw)
	}

	return n, nil
}

// mustParseLastN provides backwards compatibility for parsing last-n count.
func mustParseLastN(raw string) int {
	n, _ := parseLastN(raw)

	return n
}

// loadLastCommittedTxns returns up to `want` newest-first committed transactions.
func loadLastCommittedTxns(want int) ([]model.TransactionRecord, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	all, err := db.ListTransactions(constants.TxnRetentionCap)
	if err != nil {
		return nil, err
	}

	return filterCommittedTxns(all, want), nil
}

// filterCommittedTxns filters transactions matching committed status up to want cap.
func filterCommittedTxns(all []model.TransactionRecord, want int) []model.TransactionRecord {
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

// isLastNConfirmed previews the planned reverts and waits for affirmative confirmation.
func isLastNConfirmed(rows []model.TransactionRecord) bool {
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

	return isAffirmativeResponse(scanner.Text())
}

// revertMany applies txn.Revert to each row in newest-first order.
func revertMany(rows []model.TransactionRecord, isForce bool) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, r := range rows {
		if err := txn.Revert(db, r.ID, txn.RevertOptions{Force: isForce}); err != nil {
			return newRevertManyError(err, r.ID, r.Kind, isForce)
		}

		fmt.Printf(constants.MsgTxnReverted, r.ID, r.Kind)
	}

	fmt.Printf(constants.MsgTxnLastNDone, len(rows))

	return nil
}

func newRevertManyError(err error, id int64, kind string, isForce bool) *apperror.AppError {
	return apperror.WrapWithDetails(
		err,
		constants.OpTxnRevertMany,
		constants.ErrCodeE2015,
		err.Error(),
		constants.DomainCmdRevertTxnLastN,
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		map[string]any{"id": id, "kind": kind, "force": isForce},
	)
}

// revertManyOrExit maintains compatibility for external callers.
func revertManyOrExit(rows []model.TransactionRecord, isForce bool) {
	if err := revertMany(rows, isForce); err != nil {
		cliexit.HandleError(err, int(enums.ExitCodeGeneralError))
	}
}
