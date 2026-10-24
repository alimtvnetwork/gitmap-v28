package enums

import (
	"testing"
)

func TestRevertModeType(t *testing.T) {
	modes := AllRevertModeTypes()
	if len(modes) != 8 {
		t.Fatalf("expected 8 revert modes, got %d", len(modes))
	}

	for _, m := range modes {
		if !m.IsValid() || !m.HasMode() {
			t.Errorf("mode %s should be valid", m)
		}

		if m.String() != string(m) {
			t.Errorf("String() mismatch for %s", m)
		}

		parsed, err := ParseRevertModeType(string(m))
		if err != nil || parsed != m {
			t.Errorf("ParseRevertModeType(%q) = %v, err: %v", m, parsed, err)
		}
	}

	if RevertModeType("invalid").IsValid() {
		t.Error("invalid mode should not be valid")
	}

	if _, err := ParseRevertModeType("unknown"); err == nil {
		t.Error("expected error parsing unknown mode")
	}

	if !RevertModeVersion.IsVersion() || RevertModeRunner.IsVersion() {
		t.Error("IsVersion check failed")
	}

	if !RevertModeRunner.IsRunner() || RevertModeVersion.IsRunner() {
		t.Error("IsRunner check failed")
	}

	if !RevertModeTxn.IsTxn() || RevertModeVersion.IsTxn() {
		t.Error("IsTxn check failed")
	}

	if !RevertModeLastTxn.IsLastTxn() || RevertModeVersion.IsLastTxn() {
		t.Error("IsLastTxn check failed")
	}

	if !RevertModeLastNTxn.IsLastNTxn() || RevertModeVersion.IsLastNTxn() {
		t.Error("IsLastNTxn check failed")
	}

	if !RevertModeListTxn.IsListTxn() || RevertModeVersion.IsListTxn() {
		t.Error("IsListTxn check failed")
	}

	if !RevertModeShowTxn.IsShowTxn() || RevertModeVersion.IsShowTxn() {
		t.Error("IsShowTxn check failed")
	}

	if !RevertModePruneTxn.IsPruneTxn() || RevertModeVersion.IsPruneTxn() {
		t.Error("IsPruneTxn check failed")
	}
}

func TestRevertPromptResponseType(t *testing.T) {
	responses := AllRevertPromptResponseTypes()
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}

	for _, r := range responses {
		if !r.IsValid() || !r.HasResponse() {
			t.Errorf("response %s should be valid", r)
		}

		if r.String() != string(r) {
			t.Errorf("String() mismatch for %s", r)
		}

		parsed, err := ParseRevertPromptResponseType(string(r))
		if err != nil || parsed != r {
			t.Errorf("ParseRevertPromptResponseType(%q) = %v, err: %v", r, parsed, err)
		}
	}

	if RevertPromptResponseType("maybe").IsValid() {
		t.Error("maybe should not be valid")
	}

	if _, err := ParseRevertPromptResponseType("invalid"); err == nil {
		t.Error("expected error for invalid response")
	}

	if !RevertPromptResponseYes.IsYes() || RevertPromptResponseNo.IsYes() {
		t.Error("IsYes check failed")
	}

	if !RevertPromptResponseNo.IsNo() || RevertPromptResponseYes.IsNo() {
		t.Error("IsNo check failed")
	}
}

func TestRevertExitCodeType(t *testing.T) {
	codes := AllRevertExitCodeTypes()
	if len(codes) != 3 {
		t.Fatalf("expected 3 exit codes, got %d", len(codes))
	}

	for _, c := range codes {
		if !c.IsValid() || !c.HasExitCode() {
			t.Errorf("exit code %d should be valid", c)
		}

		if len(c.String()) == 0 {
			t.Errorf("String() should not be empty for %d", c)
		}
	}

	if RevertExitCodeType(99).IsValid() {
		t.Error("exit code 99 should not be valid")
	}

	if RevertExitCodeType(99).String() != "exit_code(99)" {
		t.Errorf("unexpected string for 99: %s", RevertExitCodeType(99).String())
	}

	parsed, err := ParseRevertExitCodeType("0")
	if err != nil || parsed != ExitCodeSuccess {
		t.Errorf("parse 0 failed: %v, %v", parsed, err)
	}

	parsed, err = ParseRevertExitCodeType("general")
	if err != nil || parsed != ExitCodeGeneralError {
		t.Errorf("parse general failed: %v, %v", parsed, err)
	}

	parsed, err = ParseRevertExitCodeType("usage")
	if err != nil || parsed != ExitCodeUsageError {
		t.Errorf("parse usage failed: %v, %v", parsed, err)
	}

	if _, err := ParseRevertExitCodeType("unknown"); err == nil {
		t.Error("expected error parsing unknown exit code")
	}

	if !ExitCodeSuccess.IsSuccess() || ExitCodeGeneralError.IsSuccess() {
		t.Error("IsSuccess check failed")
	}

	if !ExitCodeGeneralError.IsGeneralError() || ExitCodeSuccess.IsGeneralError() {
		t.Error("IsGeneralError check failed")
	}

	if !ExitCodeUsageError.IsUsageError() || ExitCodeSuccess.IsUsageError() {
		t.Error("IsUsageError check failed")
	}

	if ExitCodeSuccess.Code() != 0 || ExitCodeGeneralError.Code() != 1 || ExitCodeUsageError.Code() != 2 {
		t.Error("Code() integer values mismatch")
	}
}
