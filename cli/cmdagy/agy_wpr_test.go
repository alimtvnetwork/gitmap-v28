package cmdagy

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestParseWPROptions(t *testing.T) {
	args := []string{"start", "my-project", "-t", "3m", "-p", "custom-check", "-s", "my-suffix", "-n"}
	opts := parseWPROptions(args)

	if opts.Subcommand != "start" {
		t.Errorf("expected subcommand start, got %s", opts.Subcommand)
	}
	if opts.Target != "my-project" {
		t.Errorf("expected target my-project, got %s", opts.Target)
	}
	if opts.Interval != 3*time.Minute {
		t.Errorf("expected interval 3m, got %v", opts.Interval)
	}
	if opts.PrefixTemplate != "custom-check" {
		t.Errorf("expected prefix custom-check, got %s", opts.PrefixTemplate)
	}
	if opts.SuffixTemplate != "my-suffix" {
		t.Errorf("expected suffix my-suffix, got %s", opts.SuffixTemplate)
	}
	if !opts.IsDryRun {
		t.Errorf("expected isDryRun true")
	}
}

func TestParseWPRDeployArgs(t *testing.T) {
	alias, proj := parseDeployArgs([]string{"node-01", "gitmap"})
	if alias != "node-01" {
		t.Errorf("expected node-01, got %s", alias)
	}
	if proj != "gitmap" {
		t.Errorf("expected gitmap, got %s", proj)
	}

	if !isDeployHelpNeeded("", "") {
		t.Errorf("expected isDeployHelpNeeded true when args are empty")
	}
}

func TestWPRMachineCache(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-wpr-cache-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origOverride := store.BinaryDataDir()
	store.SetBinaryDataDirForTesting(filepath.Join(tempDir, "data"))
	defer store.SetBinaryDataDirForTesting(origOverride)

	conn := db.SSHConnection{
		Alias:     "fleet-box",
		IPAddress: "192.168.1.50",
		OS:        "linux",
	}

	UpdateWPRMachineCache(conn, "fleet-workstation-01", "online")
	cached, isFound := FindCachedMachine("fleet-box")
	if !isFound {
		t.Fatalf("expected to find fleet-box in cache")
	}
	if cached.MachineName != "fleet-workstation-01" {
		t.Errorf("expected machine name fleet-workstation-01, got %s", cached.MachineName)
	}
	if cached.IPAddress != "192.168.1.50" {
		t.Errorf("expected IP 192.168.1.50, got %s", cached.IPAddress)
	}
}

func TestRunWPRCLI_DryRunActions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitmap-wpr-actions-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	origOverride := store.BinaryDataDir()
	store.SetBinaryDataDirForTesting(filepath.Join(tempDir, "data"))
	defer store.SetBinaryDataDirForTesting(origOverride)

	// Test help
	if err := RunWPRCLI([]string{"help"}); err != nil {
		t.Errorf("RunWPRCLI help failed: %v", err)
	}

	// Test status
	if err := RunWPRCLI([]string{"status", "--json"}); err != nil {
		t.Errorf("RunWPRCLI status failed: %v", err)
	}

	// Test ls
	if err := RunWPRCLI([]string{"ls", "--json"}); err != nil {
		t.Errorf("RunWPRCLI ls failed: %v", err)
	}

	// Test disable
	if err := RunWPRCLI([]string{"disable"}); err != nil {
		t.Errorf("RunWPRCLI disable failed: %v", err)
	}

	// Test remove
	if err := RunWPRCLI([]string{"remove", "all"}); err != nil {
		t.Errorf("RunWPRCLI remove failed: %v", err)
	}
}
