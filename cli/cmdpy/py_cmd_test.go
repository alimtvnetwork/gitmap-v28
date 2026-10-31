package cmdpy

import (
	"errors"
	"testing"
)

func TestIsPyHelpRequest(t *testing.T) {
	cases := map[string]bool{"": false, "help": true, "--help": true, "-h": true, "-c": false}
	for arg, expected := range cases {
		var args []string
		if arg != "" {
			args = []string{arg}
		}
		if actual := isPyHelpRequest(args); actual != expected {
			t.Errorf("isPyHelpRequest(%s) = %v, expected %v", arg, actual, expected)
		}
	}
}

func TestResolvePythonBinary(t *testing.T) {
	bin, err := resolvePythonBinary()
	if err != nil {
		t.Logf("resolvePythonBinary notice: %v", err)
		return
	}
	if bin == "" {
		t.Errorf("expected non-empty binary string")
	}
}

func TestResolveExitCode(t *testing.T) {
	if code := resolveExitCode(nil); code != 0 {
		t.Errorf("expected code 0, got %d", code)
	}
	if code := resolveExitCode(errors.New("generic error")); code != 1 {
		t.Errorf("expected code 1, got %d", code)
	}
}

func TestResolveErrMsg(t *testing.T) {
	if msg := resolveErrMsg(nil); msg != "" {
		t.Errorf("expected empty string for nil error")
	}
	if msg := resolveErrMsg(errors.New("sample")); msg != "sample" {
		t.Errorf("expected 'sample', got '%s'", msg)
	}
}

func TestWarnAdHocAntiPatterns(t *testing.T) {
	warnAdHocAntiPatterns([]string{"-c", "open('file.txt').splitlines()"})
	warnAdHocAntiPatterns([]string{"-c", "import socket; s = socket.socket(); s.connect(('127.0.0.1', 9222))"})
	warnAdHocAntiPatterns([]string{"-c", "import json; open('transcript_full.jsonl')"})
	warnAdHocAntiPatterns([]string{"-c", "import sqlite3; conn = sqlite3.connect('agent-task.db')"})
	warnAdHocAntiPatterns([]string{"-c", "import sqlite3; conn = sqlite3.connect('conversation_summaries.db')"})
}
