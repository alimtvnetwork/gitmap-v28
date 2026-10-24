package enums

import (
	"testing"
)

func TestUndoModeType(t *testing.T) {
	testUndoModeAllAndString(t)
	testUndoModeParse(t)
	testUndoModePredicates(t)
}

func testUndoModeAllAndString(t *testing.T) {
	all := AllUndoModeTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 mode variants, got %d", len(all))
	}

	for _, mode := range all {
		if !mode.IsValid() {
			t.Errorf("expected valid for %s", mode)
		}
		if mode.String() != string(mode) {
			t.Errorf("expected string representation to match %s", mode)
		}
		if !mode.HasMode() {
			t.Errorf("expected HasMode to be true for %s", mode)
		}
	}
}

func testUndoModeParse(t *testing.T) {
	parsed, err := ParseUndoModeType("  DRY-RUN ")
	if err != nil || parsed != UndoModeDryRun {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseUndoModeType("unknown-mode"); err == nil {
		t.Error("expected error parsing invalid undo mode")
	}

	if UndoModeType("invalid").IsValid() {
		t.Error("expected invalid mode to not be valid")
	}
}

func testUndoModePredicates(t *testing.T) {
	if !UndoModeRestore.IsRestore() || UndoModeList.IsRestore() {
		t.Error("IsRestore predicate failure")
	}
	if !UndoModeList.IsList() || UndoModeDryRun.IsList() {
		t.Error("IsList predicate failure")
	}
	if !UndoModeDryRun.IsDryRun() || UndoModeFlatten.IsDryRun() {
		t.Error("IsDryRun predicate failure")
	}
	if !UndoModeFlatten.IsFlatten() || UndoModeRestore.IsFlatten() {
		t.Error("IsFlatten predicate failure")
	}
}

func TestUndoActionType(t *testing.T) {
	testUndoActionAllAndString(t)
	testUndoActionParse(t)
	testUndoActionPredicates(t)
}

func testUndoActionAllAndString(t *testing.T) {
	all := AllUndoActionTypes()
	if len(all) != 3 {
		t.Errorf("expected 3 action variants, got %d", len(all))
	}

	for _, action := range all {
		if !action.IsValid() {
			t.Errorf("expected valid for %s", action)
		}
		if action.String() != string(action) {
			t.Errorf("expected string representation to match %s", action)
		}
		if !action.HasAction() {
			t.Errorf("expected HasAction to be true for %s", action)
		}
	}
}

func testUndoActionParse(t *testing.T) {
	parsed, err := ParseUndoActionType("  RESTORED ")
	if err != nil || parsed != UndoActionRestored {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	if _, err := ParseUndoActionType("unknown-action"); err == nil {
		t.Error("expected error parsing invalid undo action")
	}

	if UndoActionType("invalid").IsValid() {
		t.Error("expected invalid action to not be valid")
	}
}

func testUndoActionPredicates(t *testing.T) {
	if !UndoActionRestored.IsRestored() || UndoActionDryRun.IsRestored() {
		t.Error("IsRestored predicate failure")
	}
	if !UndoActionDryRun.IsDryRun() || UndoActionFailed.IsDryRun() {
		t.Error("IsDryRun predicate failure")
	}
	if !UndoActionFailed.IsFailed() || UndoActionRestored.IsFailed() {
		t.Error("IsFailed predicate failure")
	}
}

func TestUndoExitCodeType(t *testing.T) {
	testUndoExitCodeAllAndString(t)
	testUndoExitCodeParse(t)
	testUndoExitCodePredicates(t)
}

func testUndoExitCodeAllAndString(t *testing.T) {
	all := AllUndoExitCodeTypes()
	if len(all) != 5 {
		t.Errorf("expected 5 exit code variants, got %d", len(all))
	}

	for _, code := range all {
		if !code.IsValid() {
			t.Errorf("expected valid for %v", code)
		}
		if !code.HasExitCode() {
			t.Errorf("expected HasExitCode to be true for %v", code)
		}
		if code.Code() != int(code) {
			t.Errorf("expected Code() to return %d, got %d", int(code), code.Code())
		}
	}

	if UndoExitCodeType(99).String() != "99" {
		t.Errorf("expected 99 string, got %s", UndoExitCodeType(99).String())
	}
}

func testUndoExitCodeParse(t *testing.T) {
	cases := []struct {
		input    string
		expected UndoExitCodeType
	}{
		{"0", UndoExitSuccess},
		{"success", UndoExitSuccess},
		{"1", UndoExitGeneral},
		{"general", UndoExitGeneral},
		{"6", UndoExitBadFlag},
		{"bad-flag", UndoExitBadFlag},
		{"7", UndoExitWriteFailed},
		{"write-failed", UndoExitWriteFailed},
		{"8", UndoExitBadConfig},
		{"bad-config", UndoExitBadConfig},
	}

	for _, c := range cases {
		parsed, err := ParseUndoExitCodeType(c.input)
		if err != nil || parsed != c.expected {
			t.Errorf("ParseUndoExitCodeType(%q) = %v, err %v; want %v", c.input, parsed, err, c.expected)
		}
	}

	if _, err := ParseUndoExitCodeType("999"); err == nil {
		t.Error("expected error parsing invalid exit code")
	}
}

func testUndoExitCodePredicates(t *testing.T) {
	if !UndoExitSuccess.IsSuccess() || UndoExitGeneral.IsSuccess() {
		t.Error("IsSuccess predicate failure")
	}
	if !UndoExitGeneral.IsGeneral() || UndoExitSuccess.IsGeneral() {
		t.Error("IsGeneral predicate failure")
	}
	if !UndoExitBadFlag.IsBadFlag() || UndoExitSuccess.IsBadFlag() {
		t.Error("IsBadFlag predicate failure")
	}
	if !UndoExitWriteFailed.IsWriteFailed() || UndoExitSuccess.IsWriteFailed() {
		t.Error("IsWriteFailed predicate failure")
	}
	if !UndoExitBadConfig.IsBadConfig() || UndoExitSuccess.IsBadConfig() {
		t.Error("IsBadConfig predicate failure")
	}
}
