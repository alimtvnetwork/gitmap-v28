package cmd

import (
	"os"
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestStripVSCodeTagFlags verifies that every supported flag shape
// (long, short, `=value`, comma-list, repeated) is peeled off argv
// AND merged into the matching env var without losing other tokens.
func TestStripVSCodeTagFlags(t *testing.T) {
	for _, k := range []string{
		constants.EnvVSCodeTagAdd,
		constants.EnvVSCodeTagSkip,
		constants.EnvVSCodeTagMarker,
	} {
		os.Unsetenv(k)
	}
	defer func() {
		for _, k := range []string{
			constants.EnvVSCodeTagAdd,
			constants.EnvVSCodeTagSkip,
			constants.EnvVSCodeTagMarker,
		} {
			os.Unsetenv(k)
		}
	}()

	in := []string{
		"clone",
		"--vscode-tag", "work",
		"-vscode-tag-skip", "git,node",
		"--vscode-tag-marker=Gemfile=ruby",
		"--vscode-tag", "urgent,daily",
		"https://example.com/r.git",
	}
	want := []string{"clone", "https://example.com/r.git"}

	got := stripVSCodeTagFlags(in)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("argv = %#v, want %#v", got, want)
	}

	addEnv := os.Getenv(constants.EnvVSCodeTagAdd)
	wantAdd := "work" + constants.EnvVSCodeTagSeparator + "urgent" + constants.EnvVSCodeTagSeparator + "daily"
	if addEnv != wantAdd {
		t.Errorf("EnvVSCodeTagAdd = %q, want %q", addEnv, wantAdd)
	}

	skipEnv := os.Getenv(constants.EnvVSCodeTagSkip)
	wantSkip := "git" + constants.EnvVSCodeTagSeparator + "node"
	if skipEnv != wantSkip {
		t.Errorf("EnvVSCodeTagSkip = %q, want %q", skipEnv, wantSkip)
	}

	markerEnv := os.Getenv(constants.EnvVSCodeTagMarker)
	if markerEnv != "Gemfile=ruby" {
		t.Errorf("EnvVSCodeTagMarker = %q, want %q", markerEnv, "Gemfile=ruby")
	}
}
