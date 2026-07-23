package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestShouldUseMultiCloneCoversReportedInvocation pins the exact PowerShell
// invocation the user reported failing with `pending task already exists ...`
// and `Cloning email-creator-v1 into https://github.com/...`. Current source
// MUST route this to runCloneMulti — never to executeDirectClone with the
// second URL as a folder name.
func TestShouldUseMultiCloneCoversReportedInvocation(t *testing.T) {
	cases := [][]string{
		// Comma-glued (PowerShell pastes this as one argv slot).
		{"https://github.com/alimtvnetwork/email-creator-v1,https://github.com/alimtvnetwork/email-reader-v3,https://github.com/alimtvnetwork/account-automator"},
		// PowerShell silent comma-split into separate argv slots.
		{
			"https://github.com/alimtvnetwork/email-creator-v1",
			"https://github.com/alimtvnetwork/email-reader-v3",
			"https://github.com/alimtvnetwork/account-automator",
		},
		// Comma + space mix.
		{
			"https://github.com/alimtvnetwork/email-creator-v1,",
			"https://github.com/alimtvnetwork/email-reader-v3,",
			"https://github.com/alimtvnetwork/account-automator",
		},
	}

	for i, positionals := range cases {
		cf := CloneFlags{Source: positionals[0], Positional: positionals}
		if !shouldUseMultiClone(cf) {
			t.Fatalf("case %d: shouldUseMultiClone returned false for %#v — would route to executeDirectClone with second URL as folder name", i, positionals)
		}
	}
}

// TestIsDirectURLAcceptsAllReportedShapes locks in the URL detector that
// both runClone routing and the stale-binary guard depend on.
func TestIsDirectURLAcceptsAllReportedShapes(t *testing.T) {
	urls := []string{
		"https://github.com/alimtvnetwork/email-creator-v1",
		"https://github.com/alimtvnetwork/email-reader-v3",
		"http://gitlab.example.com/foo/bar",
		"git@github.com:alimtvnetwork/account-automator.git",
	}
	for _, u := range urls {
		if !isDirectURL(u) {
			t.Fatalf("isDirectURL rejected %q — folder-name disambiguation will break", u)
		}
	}
	// And it must reject things that LOOK like folder names.
	notURLs := []string{"my-repo", "C:\\work\\repo", "./repo", ""}
	for _, n := range notURLs {
		if isDirectURL(n) {
			t.Fatalf("isDirectURL accepted %q as URL — would falsely trigger stale-binary guard", n)
		}
	}
}

// TestStaleBinaryGuardMessageMentionsRecoverySteps ensures the guard message
// keeps pointing the user at the exact commands they need. Without these,
// recurrence reports keep coming back with no actionable path forward.
func TestStaleBinaryGuardMessageMentionsRecoverySteps(t *testing.T) {
	msg := constants.ErrCloneStaleBinaryFolderURL
	for _, want := range []string{"gitmap update", "gitmap pending clear", "gitmap doctor", "PATH refreshes"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("stale-binary guard message missing %q\n%s", want, msg)
		}
	}
}
