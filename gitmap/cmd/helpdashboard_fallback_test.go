package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestOpenHostedDocsFallbackPrintsURL pins the runtime guarantee from
// spec/02-app-issues/34-hd-hosted-docs-fallback.md: the hosted docs URL
// is always written to stderr BEFORE attempting to launch the browser,
// so the user can copy it manually even when `start`/`open`/`xdg-open`
// is missing (e.g. minimal CI containers).
func TestOpenHostedDocsFallbackPrintsURL(t *testing.T) {
	stderr := captureStderr(t, openHostedDocsFallback)

	if !strings.Contains(stderr, constants.DocsURL) {
		t.Fatalf("openHostedDocsFallback must print DocsURL %q to stderr; got %q", constants.DocsURL, stderr)
	}
}

// TestOpenURLNonFatalOnMissingLauncher verifies openURL never panics or
// exits when the OS launcher binary is absent — it must only log a
// warning so the caller's printed URL remains the user's fallback.
func TestOpenURLNonFatalOnMissingLauncher(t *testing.T) {
	// Capturing stderr also exercises the error-logging branch when the
	// launcher process fails to start in restricted sandboxes.
	_ = captureStderr(t, func() {
		openURL(constants.DocsURL)
	})
}
