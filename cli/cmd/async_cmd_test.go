package cmd

import (
	"testing"
)

func TestAsyncCmdHelp(t *testing.T) {
	err := RunAsyncCmd([]string{})
	if err != nil {
		t.Fatalf("RunAsyncCmd without args returned error: %v", err)
	}
}

func TestAsyncCmdList(t *testing.T) {
	err := RunAsyncCmd([]string{"ls"})
	if err != nil {
		t.Fatalf("RunAsyncCmd ls returned error: %v", err)
	}
}

func TestAsyncCmdStop(t *testing.T) {
	err := RunAsyncCmd([]string{"stop", "1234"})
	if err != nil {
		t.Fatalf("RunAsyncCmd stop returned error: %v", err)
	}

	errNoArg := RunAsyncCmd([]string{"stop"})
	if errNoArg == nil {
		t.Fatalf("expected error when stopping without ID, got nil")
	}
}

func TestParseAsyncArgs(t *testing.T) {
	opts, err := parseAsyncArgs([]string{"echo", "hello", "-t", "5"})
	if err != nil {
		t.Fatalf("parseAsyncArgs returned error: %v", err)
	}

	if opts.intervalS != 5 {
		t.Errorf("expected interval 5, got %d", opts.intervalS)
	}

	if opts.command != "echo hello" {
		t.Errorf("expected command 'echo hello', got %q", opts.command)
	}

	opts2, err2 := parseAsyncArgs([]string{"-t=10", "gitmap", "status"})
	if err2 != nil {
		t.Fatalf("parseAsyncArgs returned error: %v", err2)
	}

	if opts2.intervalS != 10 {
		t.Errorf("expected interval 10, got %d", opts2.intervalS)
	}

	if opts2.command != "gitmap status" {
		t.Errorf("expected command 'gitmap status', got %q", opts2.command)
	}

	_, errMissing := parseAsyncArgs([]string{"-t", "5"})
	if errMissing == nil {
		t.Errorf("expected error when command is empty, got nil")
	}
}
