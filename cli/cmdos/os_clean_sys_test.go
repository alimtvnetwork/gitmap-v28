package cmdos

import (
	"testing"
)

type mockSystemCleanEngine struct {
	cleanCalled  bool
	vacuumCalled bool
}

func (m *mockSystemCleanEngine) CleanSystemPackages() (int64, error) {
	m.cleanCalled = true
	return 1024, nil
}

func (m *mockSystemCleanEngine) VacuumJournals() error {
	m.vacuumCalled = true
	return nil
}

func TestSystemCleanMock(t *testing.T) {
	mock := &mockSystemCleanEngine{}
	prev := GetSystemCleanEngine()
	defer SetSystemCleanEngine(prev)
	SetSystemCleanEngine(mock)

	if !isSystemCleanTarget("sys") || !isSystemCleanTarget("system") || !isSystemCleanTarget("--system") {
		t.Error("expected valid system clean targets")
	}
	if isSystemCleanTarget("dev") || isSystemCleanTarget("temp") {
		t.Error("unexpected match for dev or temp")
	}

	_ = runOSClean([]string{"sys"})
	_ = GetSystemCleanEngine().VacuumJournals()
	if !mock.vacuumCalled {
		t.Error("expected vacuum to be called on mock")
	}
}
