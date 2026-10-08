package dbengine

import (
	"context"
	"database/sql"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

func runReadOnlyCheck(tx *TxWrapper) *apperror.AppError {
	row, err := tx.QueryRow(context.Background(), "SELECT ItemName FROM TestItem WHERE ItemId = 1")
	if err != nil {
		return err
	}

	var name string
	if scanErr := row.Scan(&name); scanErr != nil {
		return apperror.WrapSimple(scanErr, "scan item name")
	}

	return nil
}

func TestWithReadOnlyTransaction(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	appErr := wrapper.WithReadOnlyTransaction(context.Background(), runReadOnlyCheck)
	if appErr != nil {
		t.Fatalf("WithReadOnlyTransaction failed: %v", appErr)
	}
}

func runExclusiveInsert(tx *TxWrapper) *apperror.AppError {
	insQ := "INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Tool', 1)"
	_, err := tx.Exec(context.Background(), insQ)
	if err != nil {
		return err
	}

	return nil
}

func TestWithExclusiveTransaction(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	appErr := wrapper.WithExclusiveTransaction(context.Background(), runExclusiveInsert)
	if appErr != nil {
		t.Fatalf("WithExclusiveTransaction failed: %v", appErr)
	}

	assertItemCount(t, wrapper, 1)
}

func runTxOptionsInsert(tx *TxWrapper) *apperror.AppError {
	insQ := "INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Tool', 1)"
	_, err := tx.Exec(context.Background(), insQ)
	if err != nil {
		return err
	}

	return nil
}

func TestWithTxOptions(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	opts := &sql.TxOptions{Isolation: sql.LevelSerializable}
	appErr := wrapper.WithTxOptions(context.Background(), opts, runTxOptionsInsert)
	if appErr != nil {
		t.Fatalf("WithTxOptions failed: %v", appErr)
	}

	assertItemCount(t, wrapper, 1)
}

func testTxPrepareAndValidate(t *testing.T, tx *TxWrapper) *apperror.AppError {
	t.Helper()
	stmt, prepErr := tx.Prepare(context.Background(), "SELECT ItemName FROM TestItem WHERE ItemId = ?")
	if prepErr != nil {
		return prepErr
	}

	defer stmt.Close()

	return tx.ValidateSql(context.Background(), "SELECT ItemName FROM TestItem")
}

func testTxCallFunction(t *testing.T, tx *TxWrapper) *apperror.AppError {
	t.Helper()
	res := tx.CallFunction(context.Background(), "UPPER", "test string")
	if res.IsFailed() {
		return res.Err
	}

	if res.Value != "TEST STRING" {
		t.Errorf("expected 'TEST STRING', got %s", res.Value)
	}

	return nil
}

func TestTxWrapper_Methods(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	appErr := wrapper.WithTransaction(context.Background(), func(tx *TxWrapper) *apperror.AppError {
		if err := testTxPrepareAndValidate(t, tx); err != nil {
			return err
		}

		return testTxCallFunction(t, tx)
	})
	if appErr != nil {
		t.Fatalf("TxWrapper methods failed: %v", appErr)
	}
}
