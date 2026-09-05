package appfault_test

import (
	"errors"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

func TestAppError_MultiDestinationFormatters(t *testing.T) {
	err := appfault.Wrap(errtype.Validation, errors.New("invalid email format"), "user validation failed").
		WithStatusCode(400).
		WithCaller(appfault.CallerInfo{File: "services/user/validator.go", Line: 42, Function: "ValidateUser"}).
		WithContext("email", "bad@")

	// 1. Stdout Banner Formatter
	stdoutOut := err.FormatStdout()
	if !strings.Contains(stdoutOut, "❌ ERROR [Validation:2] user validation failed (HTTP 400)") {
		t.Fatalf("unexpected stdout banner: %s", stdoutOut)
	}

	if !strings.Contains(stdoutOut, "Caller:  services/user/validator.go:42 (ValidateUser)") {
		t.Fatalf("expected caller in stdout banner: %s", stdoutOut)
	}

	if !strings.Contains(stdoutOut, "Cause:   invalid email format") {
		t.Fatalf("expected cause in stdout banner: %s", stdoutOut)
	}

	// 2. JSON Formatter
	jsonOut := err.FormatJSON()
	if !strings.Contains(jsonOut, `"Type": 2`) {
		t.Fatalf("expected Type 2 in json: %s", jsonOut)
	}

	if !strings.Contains(jsonOut, `"Message": "user validation failed"`) {
		t.Fatalf("expected Message in json: %s", jsonOut)
	}

	if !strings.Contains(jsonOut, `"Function": "ValidateUser"`) {
		t.Fatalf("expected Caller function in json: %s", jsonOut)
	}

	// 3. Text Log Formatter
	textLogOut := err.FormatTextLog()
	if !strings.Contains(textLogOut, "[ERROR] [Validation:2] status=400") {
		t.Fatalf("unexpected text log: %s", textLogOut)
	}

	if !strings.Contains(textLogOut, `caller="services/user/validator.go:42 (ValidateUser)"`) {
		t.Fatalf("expected caller in text log: %s", textLogOut)
	}

	if !strings.Contains(textLogOut, `msg="user validation failed"`) {
		t.Fatalf("expected msg in text log: %s", textLogOut)
	}
}

func TestResult_MultiDestinationFormatters(t *testing.T) {
	// Success case
	okRes := appfault.NewSuccess("operation completed")
	if !strings.Contains(okRes.FormatStdout(), "✅ SUCCESS: operation completed") {
		t.Fatalf("unexpected ok stdout: %s", okRes.FormatStdout())
	}

	if !strings.Contains(okRes.FormatTextLog(), "[INFO] status=200") {
		t.Fatalf("unexpected ok log: %s", okRes.FormatTextLog())
	}

	// Failure case
	failRes := appfault.NewFailure[string](errtype.NotFound, errors.New("record not found"))
	if !strings.Contains(failRes.FormatStdout(), "❌ ERROR [NotFound:3]") {
		t.Fatalf("unexpected fail stdout: %s", failRes.FormatStdout())
	}

	if !strings.Contains(failRes.FormatTextLog(), "[ERROR] [NotFound:3]") {
		t.Fatalf("unexpected fail log: %s", failRes.FormatTextLog())
	}
}
