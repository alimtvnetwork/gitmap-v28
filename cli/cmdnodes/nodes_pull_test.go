package cmdnodes

import (
	"errors"
	"testing"
)

func TestIsNodesPullAllRequest(t *testing.T) {
	positiveCases := [][]string{
		{"pa"},
		{"PA"},
		{"pull-all"},
		{"pullall"},
		{"pull"},
		{"pas"},
		{"pull-all-ssh"},
		{"pa", "--status"},
		{"pull-all", "--json"},
	}

	for _, tc := range positiveCases {
		if !IsNodesPullAllRequest(tc) {
			t.Errorf("expected IsNodesPullAllRequest(%v) = true, got false", tc)
		}
	}

	negativeCases := [][]string{
		{},
		{"ping"},
		{"clone"},
		{"pull", "errors"},
		{"history"},
		{"status"},
	}

	for _, tc := range negativeCases {
		if IsNodesPullAllRequest(tc) {
			t.Errorf("expected IsNodesPullAllRequest(%v) = false, got true", tc)
		}
	}
}

func TestRunNodesPullAll_DispatchesToHook(t *testing.T) {
	origHook := RunNodesPullAllFn
	defer func() { RunNodesPullAllFn = origHook }()

	var capturedArgs []string
	called := false
	RunNodesPullAllFn = func(args []string) error {
		called = true
		capturedArgs = args
		return nil
	}

	err := RunNodesPullAll([]string{"--status", "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Errorf("expected RunNodesPullAllFn to be called")
	}
	if len(capturedArgs) != 2 || capturedArgs[0] != "--status" || capturedArgs[1] != "--json" {
		t.Errorf("unexpected capturedArgs: %v", capturedArgs)
	}
}

func TestRunNodesPullAll_HookErrorPropagated(t *testing.T) {
	origHook := RunNodesPullAllFn
	defer func() { RunNodesPullAllFn = origHook }()

	expectedErr := errors.New("fleet pull failed")
	RunNodesPullAllFn = func(args []string) error {
		return expectedErr
	}

	err := RunNodesPullAll([]string{})
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
