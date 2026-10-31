package cmdinstall

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestResolveWingetPhpPackage(t *testing.T) {
	if constants.WingetPkgPHP != "PHP.PHP.8.5" {
		t.Fatalf("WingetPkgPHP = %q, want PHP.PHP.8.5", constants.WingetPkgPHP)
	}
	if got := resolvePackageName(constants.PkgMgrWinget, constants.ToolPHP); got != constants.WingetPkgPHP {
		t.Fatalf("winget php package = %q, want %q", got, constants.WingetPkgPHP)
	}
}

func TestResolveChocoPhpPackage(t *testing.T) {
	if got := resolvePackageName(constants.PkgMgrChocolatey, constants.ToolPHP); got != "php" {
		t.Fatalf("choco php package = %q, want php", got)
	}
}

func TestBuildWingetPhpVersionCommand(t *testing.T) {
	cmd := buildWingetCommand(constants.WingetPkgPHP, "8.5.12")
	joined := strings.Join(cmd, " ")
	if !strings.Contains(joined, "PHP.PHP.8.5") || !strings.Contains(joined, "--version 8.5.12") {
		t.Fatalf("winget php command = %q, want PHP.PHP.8.5 with --version 8.5.12", joined)
	}
}

func TestParseAptMadisonVersion(t *testing.T) {
	madison := "php8.5 | 8.5.12-1+ubuntu24.04.1+deb.sury.org+1 | https://ppa.launchpadcontent.net/ondrej/php/ubuntu noble/main amd64 Packages\n" +
		"php8.5 | 8.5.11-1+ubuntu24.04.1+deb.sury.org+1 | https://ppa.launchpadcontent.net/ondrej/php/ubuntu noble/main amd64 Packages\n"

	exact, _ := parseAptMadisonVersion(madison, "8.5.12")
	if exact != "8.5.12-1+ubuntu24.04.1+deb.sury.org+1" {
		t.Fatalf("madison 8.5.12 = %q, want the 8.5.12 PPA build", exact)
	}

	if exact, _ := parseAptMadisonVersion(madison, "8.5.13"); exact != "" {
		t.Fatalf("madison 8.5.13 = %q, want empty (not available)", exact)
	}
}

func TestPhpVersionMatches(t *testing.T) {
	if !phpVersionMatches("8.5.12", "8.5.12") {
		t.Fatal("8.5.12 should match request 8.5.12")
	}
	if phpVersionMatches("8.5.11", "8.5.12") {
		t.Fatal("8.5.11 must not match request 8.5.12")
	}
	if phpVersionMatches("8.3.6", "8.5.12") {
		t.Fatal("8.3.6 must not match request 8.5.12")
	}
	if !phpVersionMatches("8.5.12", "") {
		t.Fatal("8.5.12 should match an empty request (any 8.5.x)")
	}
	if phpVersionMatches("8.3.6", "") {
		t.Fatal("8.3.6 must not match an empty request (8.5 line only)")
	}
}

func TestPhpAptInstallArgs(t *testing.T) {
	args := phpAptInstallArgs("8.5.12-1+ubuntu24.04.1+deb.sury.org+1")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "php8.5=8.5.12-1+ubuntu24.04.1+deb.sury.org+1") {
		t.Fatalf("apt args = %q, want pinned php8.5", joined)
	}
	if !strings.Contains(joined, "php8.5-cli") || !strings.Contains(joined, "php8.5-mbstring") {
		t.Fatalf("apt args = %q, want cli and mbstring extensions", joined)
	}

	unpinned := strings.Join(phpAptInstallArgs(""), " ")
	if strings.Contains(unpinned, "php8.5=") {
		t.Fatalf("unpinned apt args = %q, want no version pin", unpinned)
	}
}

func TestPhpHelpMenuEntry(t *testing.T) {
	menu := buildInstallHelpMenu()
	for _, section := range menu.Sections {
		for _, entry := range section.Entries {
			if entry.Command == "php" {
				return
			}
		}
	}

	t.Fatal("install help menu has no php entry")
}
