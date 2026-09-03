package cliexit

import (
	"errors"
	"testing"
)

func TestHandleSpecializedExitCodes(t *testing.T) {
	var capturedCode int
	prevExit := SetExitFunc(func(code int) {
		capturedCode = code
	})
	defer SetExitFunc(prevExit)

	testCases := []struct {
		name     string
		fn       func()
		expected int
	}{
		{
			name: "validation_error",
			fn: func() {
				HandleValidationError(errors.New("validation failed"))
			},
			expected: int(ExitCodeValidationError),
		},
		{
			name: "usage_error",
			fn: func() {
				HandleUsageError(errors.New("syntax error"))
			},
			expected: int(ExitCodeUsageError),
		},
		{
			name: "general_error",
			fn: func() {
				HandleGeneralError(errors.New("io failure"))
			},
			expected: int(ExitCodeGeneralError),
		},
		{
			name: "not_found",
			fn: func() {
				HandleNotFound(errors.New("missing item"))
			},
			expected: int(ExitCodeNotFound),
		},
		{
			name: "success",
			fn: func() {
				HandleSuccess()
			},
			expected: int(ExitCodeSuccess),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			capturedCode = -1
			tc.fn()
			if capturedCode != tc.expected {
				t.Fatalf("expected exit code %d, got %d", tc.expected, capturedCode)
			}
		})
	}
}
