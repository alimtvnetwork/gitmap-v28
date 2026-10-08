package dbengine

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

func TestTransaction_CommitAndRollback(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	// Successful transaction
	err := wrapper.WithTransaction(ctx, func(tx *TxWrapper) *apperror.AppError {
		_, execErr := tx.tx.Exec("INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Service', 1)")
		if execErr != nil {
			return apperror.WrapSimple(execErr, "insert delta")
		}

		return nil
	})
	if err != nil {
		t.Fatalf("expected transaction success, got %v", err)
	}

	// Rollback transaction
	_ = wrapper.WithTransaction(ctx, func(tx *TxWrapper) *apperror.AppError {
		_, _ = tx.tx.Exec("INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Echo', 'Service', 1)")

		return apperror.WrapSimple(sql.ErrTxDone, "simulated tx failure")
	})

	repo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)
	itemsRes := repo.FindAll(ctx, 10)
	if itemsRes.IsFailed() {
		t.Fatalf("FindAll failed: %v", itemsRes.Err)
	}

	if len(itemsRes.Value) != 4 {
		t.Errorf("expected 4 items (Delta committed, Echo rolled back), got %d", len(itemsRes.Value))
	}
}

func TestTxWrapper_SqlExecutorMethods(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	appErr := wrapper.WithTransaction(context.Background(), func(tx *TxWrapper) *apperror.AppError {
		if err := testTxExecutors(t, tx); err != nil {
			return err
		}

		return testTxQueries(t, tx)
	})

	if appErr != nil {
		t.Fatalf("WithTransaction failed: %v", appErr)
	}

	assertItemCount(t, wrapper, 1)
}

func TestWithTransaction_RollbackOnError(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	appErr := wrapper.WithTransaction(context.Background(), func(tx *TxWrapper) *apperror.AppError {
		insQ := "INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Tool', 1)"
		if _, err := tx.Exec(context.Background(), insQ); err != nil {
			return err
		}

		return apperror.WrapSimple(errors.New("simulated error"), "test rollback")
	})

	if appErr == nil {
		t.Fatalf("expected error from transaction, got nil")
	}

	assertItemCount(t, wrapper, 0)
}

func runTxPanic(wrapper *DbWrapper) {
	_ = wrapper.WithTransaction(context.Background(), func(tx *TxWrapper) *apperror.AppError {
		insQ := "INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Tool', 1)"
		_, _ = tx.Exec(context.Background(), insQ)

		panic("panic inside tx") // lint-allow: panic
	})
}

func TestWithTransaction_PanicRecovery(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	didPanic := false
	func() {
		defer func() {
			if p := recover(); p != nil {
				didPanic = true
			}
		}()
		runTxPanic(wrapper)
	}()

	if !didPanic {
		t.Fatalf("expected panic to propagate out of WithTransaction")
	}

	assertItemCount(t, wrapper, 0)
}

func TestWithImmediateTransaction_Success(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	appErr := wrapper.WithImmediateTransaction(context.Background(), func(tx *TxWrapper) *apperror.AppError {
		insQ := "INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Tool', 1)"
		_, err := tx.Exec(context.Background(), insQ)

		return err
	})

	if appErr != nil {
		t.Fatalf("WithImmediateTransaction failed: %v", appErr)
	}

	assertItemCount(t, wrapper, 1)
}

func TestWithImmediateTransaction_Rollback(t *testing.T) {
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	appErr := wrapper.WithImmediateTransaction(context.Background(), func(tx *TxWrapper) *apperror.AppError {
		insQ := "INSERT INTO TestItem (ItemName, Category, IsActive) VALUES ('Delta', 'Tool', 1)"
		if _, err := tx.Exec(context.Background(), insQ); err != nil {
			return err
		}

		return apperror.WrapSimple(errors.New("immediate tx fail"), "rollback immediate")
	})

	if appErr == nil {
		t.Fatalf("expected error from WithImmediateTransaction, got nil")
	}

	assertItemCount(t, wrapper, 0)
}

func TestOpenDb_SQLiteConfigured(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "opendb_test.db")
	wrapper, appErr := OpenDb(DbSQLite, dbPath)
	if appErr != nil {
		t.Fatalf("OpenDb failed: %v", appErr)
	}

	defer wrapper.Close()

	if maxOpen := wrapper.Conn().Stats().MaxOpenConnections; maxOpen != 1 {
		t.Errorf("expected MaxOpenConnections 1, got %d", maxOpen)
	}

	verifySQLiteWrapperPragmas(t, wrapper)
}

func TestRepository_WithExecutorInsideTransaction(t *testing.T) {
	ctx := context.Background()
	wrapper := setupInMemoryDb(t)
	defer wrapper.Close()

	baseRepo := NewRepository[TestItem, TestItemFieldType](wrapper, "TestItem", scanTestItem)
	appErr := wrapper.WithTransaction(ctx, func(tx *TxWrapper) *apperror.AppError {
		return testRepoInsideTx(t, ctx, baseRepo, tx)
	})
	if appErr != nil {
		t.Fatalf("WithTransaction with repo.WithExecutor failed: %v", appErr)
	}
}
