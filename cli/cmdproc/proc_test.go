package cmdproc

import (
	"os"
	"runtime"
	"testing"
)

func TestProc_SelfKillRefused(t *testing.T) {
	err := TerminateProcess(os.Getpid(), false)
	if err == nil {
		t.Fatal("expected error when attempting to terminate current GitMap process, got nil")
	}
}

func TestProc_ParentKillRefused(t *testing.T) {
	parentPid := os.Getppid()
	if parentPid > 0 {
		err := TerminateProcess(parentPid, false)
		if err == nil {
			t.Fatal("expected error when attempting to terminate parent process, got nil")
		}
	}
}

func TestProc_SystemKillRefused(t *testing.T) {
	if runtime.GOOS == "windows" {
		err := TerminateProcess(4, false)
		if err == nil {
			t.Fatal("expected error when attempting to terminate System PID 4, got nil")
		}
	} else {
		err := TerminateProcess(1, false)
		if err == nil {
			t.Fatal("expected error when attempting to terminate init PID 1, got nil")
		}
	}
}

func TestProc_InvalidPID(t *testing.T) {
	err := RunKill([]string{"not-a-number"})
	if err == nil {
		t.Fatal("expected validation error for non-numeric PID, got nil")
	}

	err = RunKill([]string{"-5"})
	if err == nil {
		t.Fatal("expected validation error for negative PID, got nil")
	}

	err = RunKill([]string{})
	if err == nil {
		t.Fatal("expected validation error for missing PID, got nil")
	}
}

func TestProc_RunRouting(t *testing.T) {
	if err := Run([]string{}); err != nil {
		t.Fatalf("expected usage print without error, got: %v", err)
	}

	if err := Run([]string{"help"}); err != nil {
		t.Fatalf("expected help print without error, got: %v", err)
	}
}
