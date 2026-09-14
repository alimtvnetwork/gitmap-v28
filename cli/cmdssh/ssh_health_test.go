package cmdssh

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestResolveHealthDefaults(t *testing.T) {
	if resolveHealthPort(0) != 22 || resolveHealthPort(2222) != 2222 {
		t.Error("unexpected health port resolution")
	}
	if resolveHealthTimeout(0) != 1500*time.Millisecond {
		t.Error("unexpected health timeout resolution")
	}
}

func TestClassifyProbeError_Known(t *testing.T) {
	if classifyProbeError(nil) != "reachable" {
		t.Errorf("expected reachable, got: %s", classifyProbeError(nil))
	}
	refused := errors.New("dial tcp 127.0.0.1:22: connectex: No connection could be made because the target machine actively refused it.")
	if classifyProbeError(refused) != "connection refused" {
		t.Errorf("expected connection refused, got: %s", classifyProbeError(refused))
	}
	timeout := errors.New("i/o timeout")
	if classifyProbeError(timeout) != "connection timed out" {
		t.Errorf("expected connection timed out, got: %s", classifyProbeError(timeout))
	}
}

func TestClassifyProbeError_NoRoute(t *testing.T) {
	noRoute := errors.New("dial tcp: no route to host")
	if classifyProbeError(noRoute) != "no route to host" {
		t.Errorf("expected no route to host, got: %s", classifyProbeError(noRoute))
	}
	unknown := errors.New("some custom error")
	if classifyProbeError(unknown) != "some custom error" {
		t.Errorf("expected unchanged message, got: %s", classifyProbeError(unknown))
	}
}

func TestProbeHostHealth_OnlineListener(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to bind test listener: %v", err)
	}
	defer ln.Close()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	host := store.SSHHost{Alias: "local", IP: "127.0.0.1", Username: "tester"}
	res := probeHostHealth(context.Background(), host, port, 500*time.Millisecond)
	if !res.IsOnline || res.Status != "ONLINE" || res.Details != "reachable" {
		t.Errorf("unexpected health result for online listener: %+v", res)
	}
}

func TestProbeHostHealth_OfflineTarget(t *testing.T) {
	host := store.SSHHost{Alias: "down", IP: "127.0.0.1", Username: "tester"}
	res := probeHostHealth(context.Background(), host, 59999, 200*time.Millisecond)
	if res.IsOnline || res.Status != "OFFLINE" {
		t.Errorf("expected offline status, got: %+v", res)
	}
}

func TestFindHostByTarget_Alias(t *testing.T) {
	hosts := []store.SSHHost{
		{Alias: "prod", IP: "10.0.0.1", Username: "ubuntu"},
		{Alias: "dev", IP: "10.0.0.2", Username: "root"},
	}
	byAlias, isFoundAlias := findHostByTarget(hosts, "prod")
	if !isFoundAlias || byAlias.IP != "10.0.0.1" {
		t.Errorf("failed to find by alias: %+v", byAlias)
	}
}

func TestFindHostByTarget_IPAndMissing(t *testing.T) {
	hosts := []store.SSHHost{
		{Alias: "prod", IP: "10.0.0.1", Username: "ubuntu"},
		{Alias: "dev", IP: "10.0.0.2", Username: "root"},
	}
	byIP, isFoundIP := findHostByTarget(hosts, "10.0.0.2")
	if !isFoundIP || byIP.Alias != "dev" {
		t.Errorf("failed to find by IP: %+v", byIP)
	}
	_, isFoundMissing := findHostByTarget(hosts, "unknown")
	if isFoundMissing {
		t.Error("expected missing target to return false")
	}
}

func createSampleOnlineResult() SSHHealthResult {
	return SSHHealthResult{
		Status:   "ONLINE",
		IsOnline: true,
		Alias:    "devbox",
		IP:       "192.168.1.50",
		User:     "ubuntu",
		Port:     22,
		Latency:  3 * time.Millisecond,
		Details:  "reachable",
	}
}

func createSampleOfflineResult() SSHHealthResult {
	return SSHHealthResult{
		Status:   "OFFLINE",
		IsOnline: false,
		Alias:    "backup",
		IP:       "192.168.1.99",
		User:     "admin",
		Port:     22,
		Latency:  0,
		Details:  "connection refused",
	}
}

func createSampleHealthResults() []SSHHealthResult {
	return []SSHHealthResult{createSampleOnlineResult(), createSampleOfflineResult()}
}

func TestPrintHealthTable_Output(t *testing.T) {
	var buf bytes.Buffer
	results := createSampleHealthResults()
	if err := printHealthTable(&buf, results); err != nil {
		t.Fatalf("unexpected print error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ONLINE") || !strings.Contains(out, "OFFLINE") {
		t.Errorf("missing status in table output: %s", out)
	}
	if !strings.Contains(out, "devbox") || !strings.Contains(out, "connection refused") {
		t.Errorf("missing alias or details in table output: %s", out)
	}
}

func TestHandleEmptyHostsList_Output(t *testing.T) {
	var buf bytes.Buffer
	results, err := handleEmptyHostsList(&buf)
	if err != nil || len(results) != 0 {
		t.Errorf("unexpected return from handleEmptyHostsList: results=%v, err=%v", results, err)
	}
	if !strings.Contains(buf.String(), "No registered SSH machines") {
		t.Errorf("expected guidance text, got: %s", buf.String())
	}
}
