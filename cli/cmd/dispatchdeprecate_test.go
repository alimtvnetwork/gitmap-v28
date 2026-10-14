package cmd

// Demonstration of the deprecation mechanism (docs/deprecation-policy.md):
// a deprecated command name warns on stderr naming its replacement and
// removal version, then still executes its handler normally. The demo
// command is registered only inside these tests — no production command
// is deprecated by this file.

import (
	"strconv"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestDeprecatedCommandWarnsThenExecutes(t *testing.T) {
	const oldName = "demo-old-cmd"
	markDeprecated(oldName, "demo-new-cmd", removalTwoReleasesOut())
	t.Cleanup(func() { delete(deprecatedCommands, oldName) })

	executed := false
	entries := []dispatchEntry{
		{[]string{oldName}, func() error { executed = true; return nil }},
	}

	stderr := captureStderr(t, func() {
		found, err := runDispatchTable(oldName, entries)
		if !found {
			t.Fatal("deprecated command was not dispatched")
		}
		if err != nil {
			t.Fatalf("handler error: %v", err)
		}
	})

	if !executed {
		t.Fatal("deprecated command handler did not execute")
	}
	for _, want := range []string{
		"· " + oldName + " is deprecated",
		"use demo-new-cmd instead",
		"removal planned in v",
	} {
		if !strings.Contains(stderr, want) {
			t.Errorf("warning missing %q, got %q", want, stderr)
		}
	}
}

func TestNonDeprecatedCommandIsSilent(t *testing.T) {
	entries := []dispatchEntry{
		{[]string{"demo-fresh-cmd"}, func() error { return nil }},
	}

	stderr := captureStderr(t, func() {
		found, _ := runDispatchTable("demo-fresh-cmd", entries)
		if !found {
			t.Fatal("command was not dispatched")
		}
	})
	if stderr != "" {
		t.Errorf("non-deprecated command printed to stderr: %q", stderr)
	}
}

func TestRemovalTwoReleasesOut(t *testing.T) {
	parts := strings.Split(constants.Version, ".")
	if len(parts) != 3 {
		t.Skipf("constants.Version %q is not major.minor.patch", constants.Version)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		t.Fatalf("minor version %q is not numeric: %v", parts[1], err)
	}

	want := parts[0] + "." + strconv.Itoa(minor+2) + "." + parts[2]
	if got := removalTwoReleasesOut(); got != want {
		t.Errorf("removalTwoReleasesOut() = %q, want %q", got, want)
	}
}
