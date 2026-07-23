package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

func TestFormatHandoffLogLineIncludesStableFields(t *testing.T) {
	line := formatHandoffLogLine("cleanup", "remove_fail", map[string]string{
		"path": `C:\bin\gitmap.exe.old`,
		"err":  "Access is denied",
	})

	wantParts := []string{
		"phase=cleanup",
		"event=remove_fail",
		`path=C:\bin\gitmap.exe.old`,
		`err="Access is denied"`,
	}
	for _, want := range wantParts {
		if !strings.Contains(line, want) {
			t.Fatalf("log line missing %q\n%s", want, line)
		}
	}
}

func TestBuildCleanupChildArgsForwardsDebugFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"gitmap", constants.CmdUpdateRunner, constants.FlagDebugWindows, constants.FlagDebugWindowsJSON}
	args := buildCleanupChildArgs()
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, constants.CmdUpdateCleanup) ||
		!strings.Contains(joined, constants.FlagDebugWindows) ||
		!strings.Contains(joined, constants.FlagDebugWindowsJSON) {
		t.Fatalf("cleanup child args missing forwarded flags: %#v", args)
	}
}

func TestBuildCleanupChildEnvForwardsDelayAndJSONPath(t *testing.T) {
	oldArgs := os.Args
	oldDebug := os.Getenv(constants.EnvDebugWindows)
	oldJSON := os.Getenv(constants.EnvDebugWindowsJSON)
	defer func() {
		os.Args = oldArgs
		_ = os.Setenv(constants.EnvDebugWindows, oldDebug)
		_ = os.Setenv(constants.EnvDebugWindowsJSON, oldJSON)
	}()

	os.Args = []string{"gitmap", constants.CmdUpdateRunner, constants.FlagDebugWindows}
	_ = os.Setenv(constants.EnvDebugWindowsJSON, `C:\tmp\trace.jsonl`)
	env := strings.Join(buildCleanupChildEnv(), "\n")

	wantParts := []string{
		constants.EnvUpdateCleanupDelayMS + "=1500",
		constants.EnvDebugWindows + "=1",
		constants.EnvDebugWindowsJSON + `=C:\tmp\trace.jsonl`,
	}
	for _, want := range wantParts {
		if !strings.Contains(env, want) {
			t.Fatalf("cleanup child env missing %q\n%s", want, env)
		}
	}
}
