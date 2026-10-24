package enums

import "testing"

func TestAuthMethodType(t *testing.T) {
	if !AuthMethodToken.IsValid() {
		t.Fatal("expected token to be valid")
	}
	if !AuthMethodToken.IsToken() {
		t.Fatal("expected isToken true")
	}
	if AuthMethodToken.IsBrowser() {
		t.Fatal("expected isBrowser false")
	}
}

func TestCloneStrategyType(t *testing.T) {
	if !StrategyParallel.IsValid() {
		t.Fatal("expected parallel to be valid")
	}
	if !StrategyParallel.IsParallel() {
		t.Fatal("expected isParallel true")
	}
	if StrategySequential.IsParallel() {
		t.Fatal("expected isParallel false")
	}
}

func TestBatchStatusType(t *testing.T) {
	if !BatchStatusOK.IsValid() {
		t.Fatal("expected ok to be valid")
	}
	if !BatchStatusOK.IsOK() {
		t.Fatal("expected isOK true")
	}
	if !BatchStatusFailed.IsFailed() {
		t.Fatal("expected isFailed true")
	}
}
