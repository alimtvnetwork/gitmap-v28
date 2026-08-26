package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/osutil"
)

func TestInstallerMultiOSSuite(t *testing.T) {
	// Verify host OS detection
	info := osutil.DetectHostOS()
	if info.Family == "" {
		t.Fatal("expected non-empty OS family")
	}

	// Verify multi-IP parser
	ips := cmd.ParseMultiIPList("10.0.0.1, 10.0.0.2")
	if len(ips) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(ips))
	}
}
