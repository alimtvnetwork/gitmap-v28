package cmd_test

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/osutil"
)

func TestInstallerMultiOSSuite(t *testing.T) {
	// Verify host OS detection
	info := osutil.DetectHostOS()
	if info.Family == "" {
		t.Fatal("expected non-empty OS family")
	}

	// Verify multi-IP parser
	ips := cmdssh.ParseMultiIPList("10.0.0.1, 10.0.0.2")
	if len(ips) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(ips))
	}
}
