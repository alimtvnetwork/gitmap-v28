package cmdssh

import (
	"bytes"
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestExpandSubnetIPs_ValidSlash24(t *testing.T) {
	ips, err := expandSubnetIPs("192.168.1.0/24")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(ips) != 254 {
		t.Fatalf("expected 254 IPs, got: %d", len(ips))
	}
	if ips[0] != "192.168.1.1" || ips[253] != "192.168.1.254" {
		t.Errorf("unexpected IP boundary: first=%s, last=%s", ips[0], ips[253])
	}
}

func TestExpandSubnetIPs_ValidSlash30(t *testing.T) {
	ips, err := expandSubnetIPs("10.0.0.0/30")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(ips) != 2 {
		t.Fatalf("expected 2 IPs, got: %d", len(ips))
	}
	if ips[0] != "10.0.0.1" || ips[1] != "10.0.0.2" {
		t.Errorf("unexpected IP values: %v", ips)
	}
}

func TestExpandSubnetIPs_InvalidCIDR(t *testing.T) {
	_, err := expandSubnetIPs("invalid-cidr-string")
	if err == nil {
		t.Fatal("expected error for invalid CIDR, got nil")
	}
}

func TestExpandSubnetIPs_SubnetTooLarge(t *testing.T) {
	_, err := expandSubnetIPs("10.0.0.0/16")
	if err == nil {
		t.Fatal("expected error for subnet larger than 1024, got nil")
	}
}

func TestNormalizeSubnetCIDR_PlainIP(t *testing.T) {
	cidr, err := normalizeSubnetCIDR("192.168.5.10")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cidr != "192.168.5.0/24" {
		t.Errorf("expected 192.168.5.0/24, got: %s", cidr)
	}
}

func TestNormalizeSubnetCIDR_EmptyAutoDetect(t *testing.T) {
	cidr, err := normalizeSubnetCIDR("")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.HasSuffix(cidr, "/24") {
		t.Errorf("expected auto-detected /24 subnet, got: %s", cidr)
	}
}

func TestResolveScanDefaults(t *testing.T) {
	if resolveScanPort(0) != 22 || resolveScanPort(2222) != 2222 {
		t.Error("unexpected port resolution")
	}
	if resolveScanTimeout(0) != 800*time.Millisecond {
		t.Error("unexpected timeout resolution")
	}
	if resolveScanWorkers(0) != 40 || resolveScanWorkers(10) != 10 {
		t.Error("unexpected workers resolution")
	}
}

func TestProbeHostTCP_OnlineListener(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to bind test listener: %v", err)
	}
	defer ln.Close()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	ctx := context.Background()
	isOnline, latency := probeHostTCP(ctx, "127.0.0.1", port, 500*time.Millisecond)
	if !isOnline || latency <= 0 {
		t.Errorf("expected online host with positive latency, got online=%v, latency=%v", isOnline, latency)
	}
}

func TestFormatEnrollmentStatus(t *testing.T) {
	newHost := SSHScanHost{IsEnrolled: false}
	if formatEnrollmentStatus(newHost) != "[NEW]" {
		t.Errorf("expected [NEW], got: %s", formatEnrollmentStatus(newHost))
	}
	enrolledHost := SSHScanHost{IsEnrolled: true, EnrolledAlias: "devbox"}
	if formatEnrollmentStatus(enrolledHost) != "[ENROLLED: devbox]" {
		t.Errorf("expected [ENROLLED: devbox], got: %s", formatEnrollmentStatus(enrolledHost))
	}
}

func createSampleScanHost() SSHScanHost {
	return SSHScanHost{
		IP:            "192.168.1.10",
		Port:          22,
		IsOnline:      true,
		Latency:       4 * time.Millisecond,
		Hostname:      "test-box",
		IsEnrolled:    true,
		EnrolledAlias: "devbox",
		Suggestion:    "gitmap ssh devbox",
	}
}

func TestPrintScanResultTable_WithHosts(t *testing.T) {
	var buf bytes.Buffer
	hosts := []SSHScanHost{createSampleScanHost()}
	if err := printScanResultTable(&buf, hosts, "192.168.1.0/24"); err != nil {
		t.Fatalf("unexpected print error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "192.168.1.10") || !strings.Contains(out, "[ENROLLED: devbox]") {
		t.Errorf("unexpected table output: %s", out)
	}
}

func TestPrintScanResultTable_Empty(t *testing.T) {
	var buf bytes.Buffer
	err := printScanResultTable(&buf, nil, "192.168.1.0/24")
	if err != nil {
		t.Fatalf("unexpected print error: %v", err)
	}
	if !strings.Contains(buf.String(), "No active SSH machines") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
