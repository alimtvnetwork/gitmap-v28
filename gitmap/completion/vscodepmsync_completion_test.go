package completion

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestAllCommandsContainsVSCodePMSync asserts that both the long-form
// `vscode-pm-sync` command and its short alias `vpm` are exposed via
// tab-completion. Regression guard for v4.36.0+ where the command was
// added to constants but the generated completion list was stale.
func TestAllCommandsContainsVSCodePMSync(t *testing.T) {
	have := make(map[string]bool, len(AllCommands()))
	for _, v := range AllCommands() {
		have[v] = true
	}

	for _, want := range []string{constants.CmdVSCodePMSync, constants.CmdVSCodePMSyncAlias} {
		if !have[want] {
			t.Fatalf("AllCommands() missing %q (vscode-pm-sync completion drift)", want)
		}
	}
}

// TestVSCodePMSyncFlagCompletionAllShells asserts every generated shell
// completion script references the vpm command branch and surfaces all
// four documented flags (`--dry-run`, `--projects-json`, `--tag`,
// `--mode`) plus the three `--mode` value tokens. Catches the class of
// drift where someone adds a flag to constants_cli.go but forgets to
// thread it through bash/zsh/powershell.
func TestVSCodePMSyncFlagCompletionAllShells(t *testing.T) {
	cases := []struct {
		name  string
		shell string
	}{
		{"bash", constants.ShellBash},
		{"zsh", constants.ShellZsh},
		{"powershell", constants.ShellPowerShell},
	}

	wantTokens := []string{
		"vscode-pm-sync", "vpm",
		"--dry-run", "--projects-json", "--tag", "--mode",
		constants.VSCodePMSyncModeUnion,
		constants.VSCodePMSyncModeReplace,
		constants.VSCodePMSyncModeIntersection,
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			script, err := Generate(tc.shell)
			if err != nil {
				t.Fatalf("Generate(%q) returned error: %v", tc.shell, err)
			}

			for _, tok := range wantTokens {
				if !strings.Contains(script, tok) {
					t.Errorf("%s completion script missing token %q", tc.shell, tok)
				}
			}
		})
	}
}
