package cmdssh

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func assertTableContains(t *testing.T, output, expected string) {
	t.Helper()
	hasText := strings.Contains(output, expected)
	if !hasText {
		t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
	}
}

func sampleTestHost(alias, role, ip string, port int, user string) store.SSHHost {
	return store.SSHHost{
		Alias:       alias,
		ClusterRole: role,
		IP:          ip,
		Port:        port,
		Username:    user,
		CreatedAt:   time.Date(2026, 9, 17, 10, 0, 0, 0, time.UTC),
	}
}

func renderTableToString(t *testing.T, hosts []store.SSHHost) string {
	t.Helper()
	var buf bytes.Buffer
	err := RenderSSHHostsTable(&buf, hosts)
	if err != nil {
		t.Fatalf("RenderSSHHostsTable returned error: %v", err)
	}
	return buf.String()
}

func assertTableHeaders(t *testing.T, out string) {
	t.Helper()
	assertTableContains(t, out, "ALIAS")
	assertTableContains(t, out, "ROLE")
	assertTableContains(t, out, "HOST (IP:PORT)")
	assertTableContains(t, out, "USER")
	assertTableContains(t, out, "STATUS")
	assertTableContains(t, out, "ENROLLED")
	assertTableContains(t, out, strings.Repeat("-", 95))
}

func TestRenderSSHHostsTable_Empty(t *testing.T) {
	out := renderTableToString(t, nil)
	assertTableContains(t, out, "No nodes registered")
	assertTableContains(t, out, "gitmap sj add <user@ip|ip> [alias]")
}

func TestRenderSSHHostsTable_SingleNode(t *testing.T) {
	h := sampleTestHost("k8s-m1", "control", "192.168.0.100", 22, "kube")
	out := renderTableToString(t, []store.SSHHost{h})
	assertTableHeaders(t, out)
	assertTableContains(t, out, "k8s-m1")
	assertTableContains(t, out, "192.168.0.100:22")
	assertTableContains(t, out, "Total: 1 registered node(s)")
}

func TestRenderSSHHostsTable_MultiNode(t *testing.T) {
	h1 := sampleTestHost("k8s-m1", "control", "192.168.0.100", 22, "kube")
	h2 := sampleTestHost("k8s-w1", "worker", "192.168.0.101", 22, "ubuntu")
	out := renderTableToString(t, []store.SSHHost{h1, h2})
	assertTableContains(t, out, "k8s-m1")
	assertTableContains(t, out, "k8s-w1")
	assertTableContains(t, out, "Total: 2 registered node(s)")
}

func TestRenderSSHHostsTable_NonStandardPort(t *testing.T) {
	h := sampleTestHost("staging", "worker", "10.0.0.50", 2222, "deploy")
	out := renderTableToString(t, []store.SSHHost{h})
	assertTableContains(t, out, "10.0.0.50:2222")
	assertTableContains(t, out, "staging")
}

func TestRenderSSHHostsTable_DefaultFallbacks(t *testing.T) {
	h := store.SSHHost{IP: "10.0.0.99"}
	out := renderTableToString(t, []store.SSHHost{h})
	assertTableContains(t, out, "-")
	assertTableContains(t, out, "worker")
	assertTableContains(t, out, "root")
	assertTableContains(t, out, "10.0.0.99:22")
	assertTableContains(t, out, "Total: 1 registered node(s)")
}

func TestRenderSSHHostsTable_ColorsAndStyles(t *testing.T) {
	h1 := sampleTestHost("cp-node", "control-plane", "192.168.1.10", 22, "root")
	h2 := sampleTestHost("worker-node", "worker", "192.168.1.11", 22, "root")
	out := renderTableToString(t, []store.SSHHost{h1, h2})

	assertTableContains(t, out, constants.ColorCyan)
	assertTableContains(t, out, constants.ColorGreen)
	assertTableContains(t, out, constants.ColorYellow)
	assertTableContains(t, out, constants.ColorReset)
	assertTableContains(t, out, "● ready")
	assertTableContains(t, out, "cp-node")
	assertTableContains(t, out, "worker-node")
}
