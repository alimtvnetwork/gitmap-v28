package dbengine

import (
	"context"
	"testing"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
	"database/sql"
)

func assertItemCount(t *testing.T, wrapper *DbWrapper, expected int) {
	t.Helper()
	row, err := wrapper.QueryRow(context.Background(), "SELECT COUNT(*) FROM TestItem WHERE ItemName = ?", "Delta")
	if err != nil {
		t.Fatalf("query row failed: %v", err)
	}

	var count int
	if scanErr := row.Scan(&count); scanErr != nil {
		t.Fatalf("scan count failed: %v", scanErr)
	}

	if count != expected {
		t.Errorf("expected count %d, got %d", expected, count)
	}
}

func testTxExecutors(t *testing.T, tx *TxWrapper) *apperror.AppError {
	t.Helper()
	var _ SqlExecutor = tx
	if tx.Tx() == nil || tx.Compiler() == nil {
		t.Fatalf("tx or compiler is nil")
	}

	insQ := "INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Tool', 1)"
	if _, err := tx.Exec(context.Background(), insQ); err != nil {
		return err
	}

	res := tx.ExecRowsAffected(context.Background(), "UPDATE TestItem SET IsActive = 0 WHERE ItemName = ?", "Delta")
	if res.IsFailed() || res.Value != 1 {
		t.Fatalf("ExecRowsAffected failed or affected != 1: %v", res.Err)
	}

	return nil
}

func testTxQueries(t *testing.T, tx *TxWrapper) *apperror.AppError {
	t.Helper()
	row, err := tx.QueryRow(context.Background(), "SELECT Category FROM TestItem WHERE ItemName = ?", "Delta")
	if err != nil {
		return err
	}

	var cat string
	if scanErr := row.Scan(&cat); scanErr != nil || cat != "Tool" {
		t.Fatalf("QueryRow scan failed: %v, got %s", scanErr, cat)
	}

	rows, qErr := tx.Query(context.Background(), "SELECT ItemName FROM TestItem WHERE ItemName = ?", "Delta")
	if qErr != nil {
		return qErr
	}

	defer rows.Close()

	return nil
}

func verifySQLiteWrapperPragmas(t *testing.T, wrapper *DbWrapper) {
	t.Helper()
	row, err := wrapper.QueryRow(context.Background(), "PRAGMA busy_timeout;")
	if err != nil {
		t.Fatalf("query busy_timeout failed: %v", err)
	}

	var timeout int
	if scanErr := row.Scan(&timeout); scanErr != nil || timeout != 5000 {
		t.Errorf("expected busy_timeout 5000, got %d (err: %v)", timeout, scanErr)
	}

	rowFk, fkErr := wrapper.QueryRow(context.Background(), "PRAGMA foreign_keys;")
	if fkErr != nil {
		t.Fatalf("query foreign_keys failed: %v", fkErr)
	}

	var fk int
	if scanErr := rowFk.Scan(&fk); scanErr != nil || fk != 1 {
		t.Errorf("expected foreign_keys 1, got %d (err: %v)", fk, scanErr)
	}
}

func assertRepoFirstInsideTx(t *testing.T, ctx context.Context, repo *Repository[TestItem, TestItemFieldType]) *apperror.AppError {
	res := repo.First(ctx, TestItemDb.ItemName, "Alpha")
	if res.IsFailed() {
		return res.Err
	}

	if res.Value.ItemName != "Alpha" {
		t.Errorf("expected Alpha, got %s", res.Value.ItemName)
	}

	return nil
}

func assertRepoCountInsideTx(t *testing.T, ctx context.Context, repo *Repository[TestItem, TestItemFieldType]) *apperror.AppError {
	res := repo.Count(ctx, TestItemDb.Category, "Tool")
	if res.IsFailed() {
		return res.Err
	}

	if res.Value != 2 {
		t.Errorf("expected count 2, got %d", res.Value)
	}

	return nil
}

func testRepoInsideTx(t *testing.T, ctx context.Context, baseRepo *Repository[TestItem, TestItemFieldType], tx *TxWrapper) *apperror.AppError {
	txRepo := baseRepo.WithExecutor(tx)
	if err := assertRepoFirstInsideTx(t, ctx, txRepo); err != nil {
		return err
	}

	return assertRepoCountInsideTx(t, ctx, txRepo)
}

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
