package cmdpushfix

import (
	"errors"
	"flag"
	"reflect"
	"strings"
	"testing"
)

func TestParsePushFixFlagsDefaults(t *testing.T) {
	opts, err := parsePushFixFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error parsing empty flags: %v", err)
	}
	if opts.Remote != "origin" || opts.Branch != "" {
		t.Errorf("expected default origin/'', got %s/%s", opts.Remote, opts.Branch)
	}
	if opts.IsDryRun || opts.IsForce || opts.IsYes || opts.IsSSH || opts.IsHTTPS {
		t.Errorf("expected all booleans false, got %+v", opts)
	}
}

func TestParsePushFixFlagsCustom(t *testing.T) {
	args := []string{"--remote", "upstream", "--branch", "dev", "--dry-run", "--force", "--yes", "--ssh"}
	opts, err := parsePushFixFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Remote != "upstream" || opts.Branch != "dev" {
		t.Errorf("expected remote upstream and branch dev, got %s/%s", opts.Remote, opts.Branch)
	}
	if !opts.IsDryRun || !opts.IsForce || !opts.IsYes || !opts.IsSSH || opts.IsHTTPS {
		t.Errorf("expected custom flags set, got %+v", opts)
	}
}

func TestParsePushFixFlagsShort(t *testing.T) {
	args := []string{"-r", "gitlab", "-b", "feat-auth", "-n", "-f", "-y", "--https"}
	opts, err := parsePushFixFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Remote != "gitlab" || opts.Branch != "feat-auth" {
		t.Errorf("expected gitlab/feat-auth, got %s/%s", opts.Remote, opts.Branch)
	}
	if !opts.IsDryRun || !opts.IsForce || !opts.IsYes || opts.IsSSH || !opts.IsHTTPS {
		t.Errorf("expected short flags set, got %+v", opts)
	}
}

func TestParsePushFixFlagsInvalid(t *testing.T) {
	_, err := parsePushFixFlags([]string{"--unsupported-flag"})
	if err == nil {
		t.Errorf("expected error on unsupported flag, got nil")
	}
}

func TestParsePushFixFlagsHelp(t *testing.T) {
	_, err := parsePushFixFlags([]string{"--help"})
	if !errors.Is(err, flag.ErrHelp) {
		t.Errorf("expected ErrHelp on --help, got %v", err)
	}
}

func TestCheckIsNonFastForward(t *testing.T) {
	assertNonFFCase(t, " ! [rejected] master -> master (fetch first)", true)
	assertNonFFCase(t, " ! [rejected] main -> main (non-fast-forward)", true)
	assertNonFFCase(t, "Permission denied (publickey). fatal: Could not read from remote", false)
	assertNonFFCase(t, "Everything up-to-date", false)
	assertNonFFCase(t, "", false)
	assertNonFFCase(t, "fatal: remote origin not found", false)
}

func assertNonFFCase(t *testing.T, input string, want bool) {
	t.Helper()
	got := CheckIsNonFastForward(input)
	if got != want {
		t.Errorf("CheckIsNonFastForward(%q) = %v, want %v", input, got, want)
	}
}

func TestBuildPushArgsDefaults(t *testing.T) {
	state := &PushFixState{RemoteName: "origin", Branch: "main", HasUpstream: true}
	opts := PushFixOptions{}
	got := buildPushArgs(state, opts)
	want := []string{"push", "origin", "main"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildPushArgs = %v, want %v", got, want)
	}
}

func TestBuildPushArgsNoUpstream(t *testing.T) {
	state := &PushFixState{RemoteName: "origin", Branch: "feature", HasUpstream: false}
	opts := PushFixOptions{}
	got := buildPushArgs(state, opts)
	want := []string{"push", "-u", "origin", "feature"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildPushArgs = %v, want %v", got, want)
	}
}

func TestBuildPushArgsDryRunAndForce(t *testing.T) {
	state := &PushFixState{RemoteName: "upstream", Branch: "main", HasUpstream: true}
	opts := PushFixOptions{IsDryRun: true, IsForce: true}
	got := buildPushArgs(state, opts)
	want := []string{"push", "--dry-run", "--force-with-lease", "upstream", "main"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildPushArgs = %v, want %v", got, want)
	}
}

func TestFormatPushFixBadgeSafe(t *testing.T) {
	assertBadgeSafe(t, BadgeDiagnosing, "[DIAGNOSE]")
	assertBadgeSafe(t, BadgeAuthError, "[AUTH-FAIL]")
	assertBadgeSafe(t, BadgeDiverged, "[NON-FF]")
	assertBadgeSafe(t, BadgeTransport, "[TRANSPORT]")
	assertBadgeSafe(t, BadgeLeaseForce, "[FORCE-LEASE]")
	assertBadgeSafe(t, BadgeDryRun, "[DRY-RUN]")
	assertBadgeSafe(t, BadgeSuccess, "[SUCCESS]")
}

func assertBadgeSafe(t *testing.T, badge PushFixBadgeType, want string) {
	t.Helper()
	got := FormatPushFixBadge(badge, true)
	if got != want {
		t.Errorf("FormatPushFixBadge(%q, true) = %q, want %q", badge, got, want)
	}
}

func TestFormatPushFixBadgeRich(t *testing.T) {
	assertBadgeRich(t, BadgeDiagnosing, "[ 🔍 DIAGNOSING ]")
	assertBadgeRich(t, BadgeAuthError, "[ ✗ AUTH ERROR ]")
	assertBadgeRich(t, BadgeDiverged, "[ ⚠ DIVERGED ]")
	assertBadgeRich(t, BadgeTransport, "[ ⚡ PROTOCOL ]")
	assertBadgeRich(t, BadgeLeaseForce, "[ 🛡️ FORCE-LEASE ]")
	assertBadgeRich(t, BadgeDryRun, "[ ℹ DRY-RUN ]")
	assertBadgeRich(t, BadgeSuccess, "[ ✓ RECOVERED ]")
}

func assertBadgeRich(t *testing.T, badge PushFixBadgeType, want string) {
	t.Helper()
	got := FormatPushFixBadge(badge, false)
	if got != want {
		t.Errorf("FormatPushFixBadge(%q, false) = %q, want %q", badge, got, want)
	}
}

func TestColorizePushFixBadge(t *testing.T) {
	colored := ColorizePushFixBadge(BadgeSuccess, true)
	if !strings.Contains(colored, "[SUCCESS]") {
		t.Errorf("expected colored badge to contain text, got %q", colored)
	}
}

func TestBuildPushFixHelpMenu(t *testing.T) {
	menu := BuildPushFixHelpMenu()
	if menu.Title != "GITMAP PUSH-FIX & AUTH RECOVERY" {
		t.Errorf("unexpected menu title: %q", menu.Title)
	}
	if len(menu.Sections) != 2 {
		t.Errorf("expected 2 sections, got %d", len(menu.Sections))
	}
}
