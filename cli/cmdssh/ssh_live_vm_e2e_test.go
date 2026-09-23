//go:build e2e

package cmdssh

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"
)

func TestLiveVM_SSHSubnetScan(t *testing.T) {
	var buf bytes.Buffer
	opts := SSHScanOptions{
		Subnet:  "192.168.1.0/24",
		Port:    22,
		Timeout: 500 * time.Millisecond,
		Workers: 50,
	}

	hosts, err := ExecuteSubnetScan(context.Background(), &buf, opts)
	if err != nil {
		t.Fatalf("subnet scan failed: %v", err)
	}

	t.Logf("Discovered %d active hosts on 192.168.1.0/24", len(hosts))
	for _, h := range hosts {
		t.Logf("Host %s:%d (%s) - %s", h.IP, h.Port, h.Hostname, h.Latency)
	}
}

func TestLiveVM_SSHCommonJoin(t *testing.T) {
	pass := os.Getenv("GITMAP_VM_PASS")
	if pass == "" {
		t.Skip("skipping live VM common join test: GITMAP_VM_PASS not set")
	}

	var buf bytes.Buffer
	opts := SSHCommonJoinOptions{
		Username: "administrator",
		RawIPs:   []string{"192.168.1.3(w1),7(w2),12(w3)"},
		Password: pass,
		Port:     22,
		IsJSON:   true,
	}

	res, err := ExecuteCommonJoin(context.Background(), &buf, opts)
	if err != nil {
		t.Fatalf("common join execution failed: %v", err)
	}

	t.Logf("Common join completed: %d/%d succeeded", res.SuccessCount, res.TotalCount)
	for _, target := range res.Targets {
		t.Logf("Target: %s (%s) - OS: %s, Ver: %s, Success: %v, Err: %s",
			target.FullIP, target.Alias, target.DetectedOS, target.OSVersion, target.IsSuccess, target.ErrorMsg)
	}
}
