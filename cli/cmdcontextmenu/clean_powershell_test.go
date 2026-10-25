package cmdcontextmenu

import (
	"errors"
	"strings"
	"testing"
)

func withMockReg(osName string, fn regRunnerFunc, testFn func()) {
	oldOS := targetOS
	oldRunner := regRunner
	targetOS = osName
	regRunner = fn
	defer func() {
		targetOS = oldOS
		regRunner = oldRunner
	}()
	testFn()
}

func TestCleanDuplicatePowerShell_DryRun(t *testing.T) {
	callCount := 0
	mockRunner := func(args ...string) ([]byte, error) {
		callCount++

		return []byte("ok"), nil
	}
	withMockReg("windows", mockRunner, func() {
		err := CleanDuplicatePowerShell(false, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if callCount != 0 {
			t.Fatalf("expected 0 calls during dry run, got %d", callCount)
		}
	})
}

func TestCleanDuplicatePowerShell_RestoreDryRun(t *testing.T) {
	callCount := 0
	mockRunner := func(args ...string) ([]byte, error) {
		callCount++

		return []byte("ok"), nil
	}
	withMockReg("windows", mockRunner, func() {
		err := CleanDuplicatePowerShell(true, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if callCount != 0 {
			t.Fatalf("expected 0 calls during dry run, got %d", callCount)
		}
	})
}

func TestCleanDuplicatePowerShell_NonWindows(t *testing.T) {
	callCount := 0
	mockRunner := func(args ...string) ([]byte, error) {
		callCount++

		return []byte("ok"), nil
	}
	withMockReg("linux", mockRunner, func() {
		err := CleanDuplicatePowerShell(false, false)
		if err != nil {
			t.Fatalf("unexpected error on non-windows: %v", err)
		}
		if callCount != 0 {
			t.Fatalf("expected 0 calls on non-windows, got %d", callCount)
		}
	})
}

func TestCleanDuplicatePowerShell_CleanSuccess(t *testing.T) {
	var recorded [][]string
	mockRunner := func(args ...string) ([]byte, error) {
		recorded = append(recorded, args)

		return []byte("success"), nil
	}
	withMockReg("windows", mockRunner, func() {
		err := CleanDuplicatePowerShell(false, false)
		if err != nil {
			t.Fatalf("clean failed: %v", err)
		}
		if len(recorded) != 4 {
			t.Fatalf("expected 4 reg add commands, got %d", len(recorded))
		}
	})
}

func TestCleanDuplicatePowerShell_CleanError(t *testing.T) {
	mockRunner := func(args ...string) ([]byte, error) {
		return nil, errors.New("reg failed")
	}
	withMockReg("windows", mockRunner, func() {
		err := CleanDuplicatePowerShell(false, false)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestCleanDuplicatePowerShell_RestoreSuccess(t *testing.T) {
	var deleted []string
	mockRunner := func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "query" {
			return []byte("found"), nil
		}
		if len(args) > 0 && args[0] == "delete" {
			deleted = append(deleted, args[1])
		}

		return []byte("success"), nil
	}
	withMockReg("windows", mockRunner, func() {
		err := CleanDuplicatePowerShell(true, false)
		if err != nil {
			t.Fatalf("restore failed: %v", err)
		}
		if len(deleted) != 2 {
			t.Fatalf("expected 2 keys deleted, got %d", len(deleted))
		}
	})
}

func TestCleanDuplicatePowerShell_RestoreNotPresent(t *testing.T) {
	deleteCount := 0
	mockRunner := func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "query" {
			return nil, errors.New("not found")
		}
		deleteCount++

		return []byte("success"), nil
	}
	withMockReg("windows", mockRunner, func() {
		err := CleanDuplicatePowerShell(true, false)
		if err != nil {
			t.Fatalf("restore failed: %v", err)
		}
		if deleteCount != 0 {
			t.Fatalf("expected 0 deletes when keys absent, got %d", deleteCount)
		}
	})
}

func TestDetectPowerShellStatus(t *testing.T) {
	mockRunner := func(args ...string) ([]byte, error) {
		key := args[1]
		if strings.Contains(key, "PowerShellMenu") {
			return []byte("found"), nil
		}

		return nil, errors.New("not found")
	}
	withMockReg("windows", mockRunner, func() {
		st := DetectPowerShellStatus()
		if !st.HasLegacyMenu {
			t.Errorf("expected HasLegacyMenu true")
		}
		if st.HasPwsh7Menu {
			t.Errorf("expected HasPwsh7Menu false")
		}
		if st.IsDuplicate {
			t.Errorf("expected IsDuplicate false")
		}
	})
}
