// Package cmd — install_unit_test.go covers pure helpers in install.go,
// installlist.go, installdetect.go, and uninstall.go. No DB, no exec, no
// network — only deterministic logic.
package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// ─────────────────────────── uninstall: flag parsing ──────────────────

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

// ─────────────────────────── uninstall: command builders ──────────────

func TestBuildUninstallCommand(t *testing.T) {
	cases := []struct {
		name    string
		manager string
		tool    string
		purge   bool
		head    string // first arg
		hasFlag string // optional substring assertion
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

// ─────────────────────────── installlist: status resolution ───────────

func TestResolveToolStatusFromDB(t *testing.T) {
	installed := map[string]string{"node": "20.11.0"}

	status, version := resolveToolStatus("node", installed)
	if status != constants.StatusInstalled {
		t.Fatalf("expected installed glyph for DB hit, got %q", status)
	}
	if version != "20.11.0" {
		t.Fatalf("expected DB version 20.11.0, got %q", version)
	}
}

func TestResolveToolStatusUnknownTool(t *testing.T) {
	// "definitely-not-a-real-binary-xyz" should miss both DB + PATH probe.
	status, version := resolveToolStatus("definitely-not-a-real-binary-xyz", map[string]string{})
	if status != constants.StatusNotInstalled {
		t.Fatalf("expected not-installed glyph for missing tool, got %q", status)
	}
	if version != "—" {
		t.Fatalf("expected dash version for missing tool, got %q", version)
	}
}

func TestPickDisplayVersion(t *testing.T) {
	cases := []struct {
		name string
		in   store.InstalledTool
		want string
	}{
		{"valid", store.InstalledTool{VersionString: "1.2.3"}, "1.2.3"},
		{"empty-falls-back-to-dash", store.InstalledTool{VersionString: ""}, "—"},
		{"zeros-fall-back-to-dash", store.InstalledTool{VersionString: "0.0.0"}, "—"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pickDisplayVersion(tc.in)
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// ─────────────────────────── installlist: category ordering ───────────

func TestSortedCategoryNamesCoreFirst(t *testing.T) {
	got := sortedCategoryNames()
	if len(got) == 0 {
		t.Skip("no categories defined; nothing to assert")
	}
	if got[0] != constants.ToolCategoryCore {
		t.Fatalf("ToolCategoryCore must sort first; got order %v", got)
	}
	// Ensure non-Core tail is stable + alphabetical.
	for i := 2; i < len(got); i++ {
		if got[i] < got[i-1] {
			t.Fatalf("non-Core categories must be alphabetical; got %v", got)
		}
	}
}

// ─────────────────────────── installdetect: override path ─────────────

func TestResolvePackageManagerOverride(t *testing.T) {
	got := resolvePackageManager("brew")
	if got != "brew" {
		t.Fatalf("override must win; got %q", got)
	}
}

func TestResolvePackageManagerEmptyDelegates(t *testing.T) {
	// Empty override delegates to detectPackageManager — result is
	// platform-dependent but must be a non-empty known manager.
	got := resolvePackageManager("")
	if got == "" {
		t.Fatal("empty override must still return a default manager")
	}

	known := []string{
		constants.PkgMgrChocolatey, constants.PkgMgrWinget, constants.PkgMgrBrew,
		constants.PkgMgrApt, constants.PkgMgrDnf, constants.PkgMgrPacman,
	}
	for _, k := range known {
		if got == k {
			return
		}
	}

	t.Fatalf("detected manager %q is not in the known set %v", got, known)
}

// containsToken reports whether any element of args contains substr.
func containsToken(args []string, substr string) bool {
	if substr == "" {
		return true
	}

	for _, a := range args {
		if strings.Contains(a, substr) {
			return true
		}
	}

	return false
}
