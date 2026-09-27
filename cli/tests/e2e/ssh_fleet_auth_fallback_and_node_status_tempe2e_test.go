//go:build tempe2e

package e2e

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func isAuthTempE2ESkipped() bool {
	return os.Getenv("RUN_TEMP_E2E") != "1"
}

func TestSSHFleetAuthFallback_TempE2E(t *testing.T) {
	if isAuthTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	conns, err := cmdssh.FetchAllSSHConnectionsForTest()
	if err != nil {
		t.Fatalf("failed to load SSH connections: %v", err)
	}

	for _, c := range conns {
		if c.Alias != "w1" && c.Alias != "w2" && c.Alias != "w3" {
			continue
		}
		t.Run(c.Alias, func(t *testing.T) {
			client, isOk := cmdssh.ConnectSSHClientForTest(c)
			if !isOk || client == nil {
				t.Fatalf("expected successful SSH connection with auth fallback for %s (%s)", c.Alias, c.IPAddress)
			}
			defer client.Close()
		})
	}
}

func TestSSHFleetNodeStatusLiveProbe_TempE2E(t *testing.T) {
	if isAuthTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	hosts := []store.SSHHost{
		{Alias: "offline-1", IP: "10.20.0.11", Port: 22, Username: "admin"},
		{Alias: "w2", IP: "192.168.1.7", Port: 22, Username: "administrator"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var buf bytes.Buffer
	err := cmdssh.RenderSSHHostsTable(&buf, hosts)
	if err != nil {
		t.Fatalf("RenderSSHHostsTable failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "w2") {
		t.Errorf("expected table to contain w2, got:\n%s", out)
	}
	_ = ctx
}

func TestSSHAggregateFleetUpdate_TempE2E(t *testing.T) {
	if isAuthTempE2ESkipped() {
		t.Skip("skipping temporary e2e test; run on-demand with RUN_TEMP_E2E=1 and -tags=tempe2e")
	}

	err := cmdssh.RunSSHUpdateCLI([]string{"agm", "--dry-run"})
	if err != nil {
		t.Fatalf("RunSSHUpdateCLI failed: %v", err)
	}
}
