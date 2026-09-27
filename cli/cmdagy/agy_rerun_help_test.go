package cmdagy

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestIsRerunHelpToken_Affirmative(t *testing.T) {
	cases := []string{
		"help", "HELP", "Help", "hElp",
		"-h", "-H", "--help", "--HELP", "-help", "--h",
		"?", "/?", "-?", "--?",
		"/h", "/H", "/help", "/HELP",
		"man", "MAN", "info", "INFO",
		"usage", "USAGE", "-usage", "--usage", "/usage", "/USAGE",
		"  help  ", "  --help  ", "  /?  ", "  ?  ",
	}

	for _, token := range cases {
		if !isRerunHelpToken(token) {
			t.Errorf("expected isRerunHelpToken(%q) to be true", token)
		}
	}
}

func TestIsRerunHelpToken_Negative(t *testing.T) {
	cases := []string{
		"1", "2", "3", "gitmap", "wp-exam", "movie-cli",
		"all", "queue", "--restart", "-r", "--dry-run", "-d",
		"--model", "pro", "--prompt", "is-done", "", "   ",
	}

	for _, token := range cases {
		if isRerunHelpToken(token) {
			t.Errorf("expected isRerunHelpToken(%q) to be false", token)
		}
	}
}

func TestIsRerunHelpRequested_Positions(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"empty", []string{}, false},
		{"target only", []string{"1"}, false},
		{"flags only", []string{"--dry-run", "-r"}, false},
		{"sole help", []string{"help"}, true},
		{"sole dash h", []string{"-h"}, true},
		{"sole question", []string{"?"}, true},
		{"sole slash question", []string{"/?"}, true},
		{"trailing help", []string{"1", "help"}, true},
		{"trailing dash help", []string{"all", "--help"}, true},
		{"trailing slash h", []string{"queue", "/h"}, true},
		{"middle help", []string{"--dry-run", "help", "1"}, true},
		{"uppercase help", []string{"HELP"}, true},
		{"uppercase usage", []string{"--USAGE"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isRerunHelpRequested(tc.args)
			if got != tc.expected {
				t.Errorf("isRerunHelpRequested(%v) = %v, expected %v", tc.args, got, tc.expected)
			}
		})
	}
}

func TestIsHelpKeyword_MatchesAllHelpTokens(t *testing.T) {
	helpWords := []string{"help", "INFO", "man", "/?", "?", "--help", "-h", "/h", "usage"}
	for _, hw := range helpWords {
		if !isHelpKeyword(hw) {
			t.Errorf("expected isHelpKeyword(%q) to be true", hw)
		}
	}

	validProjects := []string{"gitmap", "wp-exam", "1", "2", "proj"}
	for _, vp := range validProjects {
		if isHelpKeyword(vp) {
			t.Errorf("expected isHelpKeyword(%q) to be false", vp)
		}
	}
}

func TestRenderAgyRerunHelp_Output(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	RenderAgyRerunHelp()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	out := buf.String()

	if len(out) == 0 {
		t.Fatal("RenderAgyRerunHelp produced no output")
	}
	expectedStrings := []string{
		"ANTIGRAVITY RERUN & PROMPT REPLAY GUIDE",
		"gitmap rerun",
		"gitmap rerun 1",
		"gitmap rerun all",
		"gitmap rerun queue",
		"gitmap rerun help",
		"--dry-run",
		"--restart",
	}
	for _, exp := range expectedStrings {
		if !bytes.Contains([]byte(out), []byte(exp)) {
			t.Errorf("expected output to contain %q", exp)
		}
	}
}
