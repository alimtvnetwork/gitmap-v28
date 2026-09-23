package cmdssh

import (
	"errors"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestIsNodesVersionRequest(t *testing.T) {
	cases := []struct {
		args     []string
		expected bool
	}{
		{[]string{"-v"}, true},
		{[]string{"--version"}, true},
		{[]string{"gitmap", "-v"}, true},
		{[]string{"gitmap", "--version"}, true},
		{[]string{"version"}, true},
		{[]string{"ls"}, false},
		{[]string{"clear"}, false},
		{[]string{"rm", "w1"}, false},
	}

	for _, tc := range cases {
		actual := isNodesVersionRequest(tc.args)
		if actual != tc.expected {
			t.Errorf("isNodesVersionRequest(%v) = %v; want %v", tc.args, actual, tc.expected)
		}
	}
}

func TestParseRemoteVersionOutput(t *testing.T) {
	ver, isInstalled := parseRemoteVersionOutput("gitmap v6.315.0\nextra info", nil)
	if !isInstalled || ver != "gitmap v6.315.0" {
		t.Errorf("expected installed gitmap v6.315.0, got ver=%q isInstalled=%v", ver, isInstalled)
	}

	ver2, isInstalled2 := parseRemoteVersionOutput("v6.312.0", nil)
	if !isInstalled2 || ver2 != "v6.312.0" {
		t.Errorf("expected installed v6.312.0, got ver=%q isInstalled=%v", ver2, isInstalled2)
	}

	ver3, isInstalled3 := parseRemoteVersionOutput("'gitmap' is not recognized as an internal or external command", nil)
	if isInstalled3 || ver3 != "not installed" {
		t.Errorf("expected not installed, got ver=%q isInstalled=%v", ver3, isInstalled3)
	}

	ver4, isInstalled4 := parseRemoteVersionOutput("", errors.New("command failed"))
	if isInstalled4 || ver4 != "not installed" {
		t.Errorf("expected not installed on err, got ver=%q isInstalled=%v", ver4, isInstalled4)
	}
}

func TestIsAllTarget(t *testing.T) {
	cases := []struct {
		target   string
		expected bool
	}{
		{"", true},
		{"all", true},
		{"all-nodes", true},
		{"ALL-NODES", true},
		{"allnodes", true},
		{"nodes", true},
		{"w1", false},
		{"192.168.1.3", false},
	}

	for _, tc := range cases {
		actual := isAllTarget(tc.target)
		if actual != tc.expected {
			t.Errorf("isAllTarget(%q) = %v; want %v", tc.target, actual, tc.expected)
		}
	}
}

func TestParseUpdateOptions(t *testing.T) {
	opts := parseUpdateOptions([]string{"all-nodes", "--except", "w1,192.168.1.7,2", "--dry-run"})
	if opts.Target != "all-nodes" {
		t.Errorf("expected target all-nodes, got %q", opts.Target)
	}
	if opts.Except != "w1,192.168.1.7,2" {
		t.Errorf("expected except w1,192.168.1.7,2, got %q", opts.Except)
	}
	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun true, got false")
	}

	opts2 := parseUpdateOptions([]string{"agm", "--remote", "w2"})
	if opts2.Pkg != "agm" {
		t.Errorf("expected pkg agm, got %q", opts2.Pkg)
	}
	if opts2.Target != "w2" {
		t.Errorf("expected target w2, got %q", opts2.Target)
	}
}

func TestIsConnExcluded(t *testing.T) {
	conn := db.SSHConnection{
		Alias:     "w2",
		IPAddress: "192.168.1.7",
		Username:  "Administrator",
	}

	excludeByAlias := []string{"w1", "w2"}
	if !isConnExcluded(conn, excludeByAlias) {
		t.Errorf("expected conn to be excluded by alias w2")
	}

	excludeByIP := []string{"192.168.1.7"}
	if !isConnExcluded(conn, excludeByIP) {
		t.Errorf("expected conn to be excluded by IP")
	}

	excludeOther := []string{"w1", "w3", "99"}
	if isConnExcluded(conn, excludeOther) {
		t.Errorf("did not expect conn to be excluded by other list")
	}
}

func TestSSHMacroSubcommandDispatch(t *testing.T) {
	if err := runSSHMacroCLI([]string{"help"}); err != nil {
		t.Errorf("expected nil from macro help, got %v", err)
	}
	if err := runSSHMacroCLI([]string{"ls"}); err != nil {
		t.Errorf("expected nil from macro ls, got %v", err)
	}
}
