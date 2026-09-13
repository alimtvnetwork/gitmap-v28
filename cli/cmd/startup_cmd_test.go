package cmd

import (
	"testing"
)

func TestParseStartupAddArgs(t *testing.T) {
	cases := []struct {
		args     []string
		wantName string
		wantFreq string
		wantType string
		wantIcon string
	}{
		{
			args:     []string{"C:\\scripts\\run.ps1", "--frequency=once-a-week", "--icon=C:\\icons\\app.ico"},
			wantName: "run",
			wantFreq: "once-a-week",
			wantType: "powershell",
			wantIcon: "C:\\icons\\app.ico",
		},
		{
			args:     []string{"macro:sync-all", "--frequency=everytime"},
			wantName: "sync-all",
			wantFreq: "everytime",
			wantType: "macro",
			wantIcon: "",
		},
		{
			args:     []string{"/home/user/clean.sh", "--name=cleaner"},
			wantName: "cleaner",
			wantFreq: "everytime",
			wantType: "bash",
			wantIcon: "",
		},
	}

	for _, tc := range cases {
		opts, err := parseStartupAddArgs(tc.args)
		if err != nil {
			t.Fatalf("parseStartupAddArgs(%v) failed: %v", tc.args, err)
		}
		if opts.name != tc.wantName {
			t.Errorf("name = %q, want %q", opts.name, tc.wantName)
		}
		if opts.frequency != tc.wantFreq {
			t.Errorf("frequency = %q, want %q", opts.frequency, tc.wantFreq)
		}
		targetType := detectTargetType(opts.target)
		if targetType != tc.wantType {
			t.Errorf("targetType = %q, want %q", targetType, tc.wantType)
		}
		if opts.icon != tc.wantIcon {
			t.Errorf("icon = %q, want %q", opts.icon, tc.wantIcon)
		}
	}
}

func TestDeriveStartupName(t *testing.T) {
	if n := deriveStartupName("macro:daily-backup"); n != "daily-backup" {
		t.Errorf("expected daily-backup, got %q", n)
	}
	if n := deriveStartupName("C:\\tools\\checker.ps1"); n != "checker" {
		t.Errorf("expected checker, got %q", n)
	}
	if n := deriveStartupName("/usr/local/bin/agent"); n != "agent" {
		t.Errorf("expected agent, got %q", n)
	}
}
