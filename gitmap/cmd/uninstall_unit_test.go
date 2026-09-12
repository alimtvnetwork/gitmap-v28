package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestHasPositionalToolArg(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{"empty", []string{}, false},
		{"only-bool-flags", []string{"--dry-run", "--force"}, false},
		{"tool-name", []string{"vscode"}, true},
		{"tool-with-flags", []string{"--force", "node", "--purge"}, true},
		{"shell-mode-consumes-value", []string{"--shell-mode", "bash"}, false},
		{"shell-mode-then-tool", []string{"--shell-mode", "zsh", "git"}, true},
		{"keep-data-passthrough", []string{"--confirm", "--keep-data"}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := hasPositionalToolArg(tc.args)
			if got != tc.want {
				t.Fatalf("args=%v: got %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestBuildUninstallCommand(t *testing.T) {
	cases := []struct {
		name    string
		manager string
		tool    string
		purge   bool
		head    string
		hasFlag string
	}{
		{"choco-no-purge", constants.PkgMgrChocolatey, "vscode", false, "choco", "-y"},
		{"choco-purge", constants.PkgMgrChocolatey, "vscode", true, "choco", "-x"},
		{"winget", constants.PkgMgrWinget, "git", false, "winget", "uninstall"},
		{"apt-remove", constants.PkgMgrApt, "node", false, "sudo", "remove"},
		{"apt-purge", constants.PkgMgrApt, "node", true, "sudo", "purge"},
		{"brew", constants.PkgMgrBrew, "go", false, "brew", "uninstall"},
		{"snap", constants.PkgMgrSnap, "code", false, "sudo", "remove"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildUninstallCommand(tc.manager, tc.tool, tc.purge)
			if len(got) == 0 || got[0] != tc.head {
				t.Fatalf("head: got %v, want first=%q", got, tc.head)
			}

			if !containsToken(got, tc.hasFlag) {
				t.Fatalf("missing token %q in %v", tc.hasFlag, got)
			}
		})
	}
}

func TestBuildChocoUninstall(t *testing.T) {
	plain := buildChocoUninstall("vscode", false)
	if containsToken(plain, "-x") {
		t.Fatalf("plain choco uninstall should not include -x: %v", plain)
	}

	purged := buildChocoUninstall("vscode", true)
	if !containsToken(purged, "-x") {
		t.Fatalf("purge choco uninstall must include -x: %v", purged)
	}
}

func TestBuildAptUninstall(t *testing.T) {
	rm := buildAptUninstall("node", false)
	if !containsToken(rm, "remove") || containsToken(rm, "purge") {
		t.Fatalf("non-purge apt should use 'remove': %v", rm)
	}

	pg := buildAptUninstall("node", true)
	if !containsToken(pg, "purge") || containsToken(pg, "remove") {
		t.Fatalf("purge apt should use 'purge': %v", pg)
	}
}
