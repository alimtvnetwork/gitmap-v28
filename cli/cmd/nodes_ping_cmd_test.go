package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNodesPing_ParseOptions(t *testing.T) {
	// Test defaults
	optsDefault := parseNodesPingOptions([]string{})
	if optsDefault.count != 2 {
		t.Errorf("expected default count 2, got %d", optsDefault.count)
	}
	if optsDefault.timeoutMs != 1500 {
		t.Errorf("expected default timeout 1500ms, got %d", optsDefault.timeoutMs)
	}
	if optsDefault.isJSON || optsDefault.isRaw || optsDefault.filterSSH {
		t.Errorf("expected default flags to be false")
	}

	// Test overrides
	args := []string{"--count", "4", "--timeout", "2s", "--raw", "--json", "--ssh", "main"}
	opts := parseNodesPingOptions(args)
	if opts.count != 4 {
		t.Errorf("expected count 4, got %d", opts.count)
	}
	if opts.timeoutMs != 2000 {
		t.Errorf("expected timeout 2000ms, got %d", opts.timeoutMs)
	}
	if !opts.isRaw {
		t.Errorf("expected isRaw true")
	}
	if !opts.isJSON {
		t.Errorf("expected isJSON true")
	}
	if !opts.filterSSH {
		t.Errorf("expected filterSSH true")
	}
	if opts.targetFilter != "main" {
		t.Errorf("expected targetFilter main, got %s", opts.targetFilter)
	}
}

func TestNodesPing_ParseWindowsOutput(t *testing.T) {
	winSuccess := `Pinging 192.168.1.3 with 32 bytes of data:
Reply from 192.168.1.3: bytes=32 time<1ms TTL=128
Reply from 192.168.1.3: bytes=32 time=2ms TTL=128

Ping statistics for 192.168.1.3:
    Packets: Sent = 2, Received = 2, Lost = 0 (0% loss),
Approximate round trip times in milli-seconds:
    Minimum = 0ms, Maximum = 2ms, Average = 1ms
`
	sent, recv, lost, lossPct, minR, avgR, maxR := parseMachinePingOutput(winSuccess, 2)
	if sent != 2 || recv != 2 || lost != 0 || lossPct != 0 {
		t.Errorf("unexpected packet stats: sent=%d recv=%d lost=%d lossPct=%d", sent, recv, lost, lossPct)
	}
	if minR != "0ms" || maxR != "2ms" || avgR != "1ms" {
		t.Errorf("unexpected RTT stats: min=%s max=%s avg=%s", minR, maxR, avgR)
	}

	winTimeout := `Pinging 192.168.1.20 with 32 bytes of data:
Request timed out.
Request timed out.

Ping statistics for 192.168.1.20:
    Packets: Sent = 2, Received = 0, Lost = 2 (100% loss),
`
	sent2, recv2, lost2, lossPct2, _, avgR2, _ := parseMachinePingOutput(winTimeout, 2)
	if sent2 != 2 || recv2 != 0 || lost2 != 2 || lossPct2 != 100 {
		t.Errorf("unexpected timeout packet stats: sent=%d recv=%d lost=%d lossPct=%d", sent2, recv2, lost2, lossPct2)
	}
	if avgR2 != "-" {
		t.Errorf("expected avgR2 to be '-', got %s", avgR2)
	}
}

func TestNodesPing_ParseLinuxOutput(t *testing.T) {
	linuxSuccess := `PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.035 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=64 time=0.041 ms

--- 127.0.0.1 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1024ms
rtt min/avg/max/mdev = 0.035/0.038/0.041/0.003 ms
`
	sent, recv, lost, lossPct, minR, avgR, maxR := parseMachinePingOutput(linuxSuccess, 2)
	if sent != 2 || recv != 2 || lost != 0 || lossPct != 0 {
		t.Errorf("unexpected packet stats: sent=%d recv=%d lost=%d lossPct=%d", sent, recv, lost, lossPct)
	}
	if minR != "0.035ms" || avgR != "0.038ms" || maxR != "0.041ms" {
		t.Errorf("unexpected Linux RTT stats: min=%s max=%s avg=%s", minR, maxR, avgR)
	}
}

func TestNodesPing_ResolveStatus(t *testing.T) {
	st1 := resolvePingStatus(true, true, 0)
	if !strings.Contains(st1, "ONLINE") {
		t.Errorf("expected ONLINE, got %s", st1)
	}

	st2 := resolvePingStatus(false, true, 100)
	if !strings.Contains(st2, "REACHABLE") {
		t.Errorf("expected REACHABLE, got %s", st2)
	}

	st3 := resolvePingStatus(false, false, 100)
	if !strings.Contains(st3, "OFFLINE") {
		t.Errorf("expected OFFLINE, got %s", st3)
	}

	st4 := resolvePingStatus(true, false, 50)
	if !strings.Contains(st4, "DEGRADED") {
		t.Errorf("expected DEGRADED, got %s", st4)
	}
}

func TestNodesPing_RenderTable(t *testing.T) {
	results := []NodePingResult{
		{
			Alias:             "w1",
			Role:              "worker",
			Host:              "192.168.1.3",
			Port:              22,
			PacketsSent:       2,
			PacketsReceived:   2,
			PacketsLost:       0,
			PacketLossPercent: 0,
			AvgRTT:            "1ms",
			TCPOk:             true,
			TCPRTT:            "1ms",
			ICMPOk:            true,
			Status:            "● ONLINE",
			Duration:          100 * time.Millisecond,
		},
		{
			Alias:             "main",
			Role:              "worker",
			Host:              "192.168.1.20",
			Port:              22,
			PacketsSent:       2,
			PacketsReceived:   0,
			PacketsLost:       2,
			PacketLossPercent: 100,
			AvgRTT:            "-",
			TCPOk:             true,
			TCPRTT:            "2ms",
			ICMPOk:            false,
			Status:            "● REACHABLE (TCP)",
			Duration:          150 * time.Millisecond,
		},
	}

	var buf bytes.Buffer
	opts := nodesPingOptions{count: 2, timeoutMs: 1500}
	err := renderPingTable(&buf, results, opts)
	if err != nil {
		t.Fatalf("renderPingTable failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "GITMAP FLEET NODES PING") {
		t.Errorf("expected banner in output, got:\n%s", out)
	}
	if !strings.Contains(out, "w1") || !strings.Contains(out, "main") {
		t.Errorf("expected nodes in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Fleet Ping") {
		t.Errorf("expected summary footer in output, got:\n%s", out)
	}
}

func TestNodesPing_JSONSerialization(t *testing.T) {
	res := []NodePingResult{
		{
			Alias:  "w1",
			Host:   "192.168.1.3",
			Port:   22,
			Status: "● ONLINE",
		},
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	if !strings.Contains(string(data), `"host":"192.168.1.3"`) {
		t.Errorf("expected serialized host in json, got: %s", string(data))
	}
}

func TestNodesPing_Help(t *testing.T) {
	if !isNodesPingHelpRequest([]string{"--help"}) {
		t.Errorf("expected --help to be recognized")
	}
	if !isNodesPingHelpRequest([]string{"-h"}) {
		t.Errorf("expected -h to be recognized")
	}
	if !isNodesPingHelpRequest([]string{"help"}) {
		t.Errorf("expected help to be recognized")
	}
}
