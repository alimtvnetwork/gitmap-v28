package cmdnodes

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseNodesFilterOptions(t *testing.T) {
	opts := parseNodesFilterOptions([]string{"--json", "--fast", "--ssh", "main"})
	if !opts.isJSON {
		t.Errorf("expected isJSON to be true")
	}
	if !opts.isFast {
		t.Errorf("expected isFast to be true")
	}
	if !opts.filterSSH {
		t.Errorf("expected filterSSH to be true")
	}
	if opts.targetFilter != "main" {
		t.Errorf("expected targetFilter 'main', got '%s'", opts.targetFilter)
	}
}

func TestUnifiedNodes_Help(t *testing.T) {
	err := RunUnifiedNodesCLI([]string{"--help"})
	if err != nil {
		t.Fatalf("expected nil error for help, got: %v", err)
	}
}

func TestUnifiedNodes_RenderTable(t *testing.T) {
	sampleNodes := []UnifiedFleetNode{
		{
			Alias:      "main",
			Role:       "worker",
			Host:       "192.168.1.20",
			Port:       22,
			User:       "administrator",
			Status:     "● ready",
			Subsystems: []string{"SSH", "Cluster", "SC"},
			EnrolledAt: "2026-09-29 12:00:00",
			IsOnline:   true,
		},
		{
			Alias:      "w1",
			Role:       "worker",
			Host:       "192.168.1.3",
			Port:       22,
			User:       "Administrator",
			Status:     "○ offline",
			Subsystems: []string{"SSH", "SC"},
			EnrolledAt: "2026-09-29 12:00:00",
			IsOnline:   false,
		},
	}

	var buf bytes.Buffer
	err := renderUnifiedNodesTable(&buf, sampleNodes)
	if err != nil {
		t.Fatalf("unexpected error rendering table: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALIAS") || !strings.Contains(out, "HOST (IP:PORT)") {
		t.Errorf("table missing header columns")
	}
	if !strings.Contains(out, "main") || !strings.Contains(out, "192.168.1.20") {
		t.Errorf("table missing main node")
	}
	if !strings.Contains(out, "SUPPORTED COMMANDS & CLUSTERS") {
		t.Errorf("table missing command matrix header")
	}
	matrixIdx := strings.Index(out, "SUPPORTED COMMANDS & CLUSTERS")
	primaryIdx := strings.Index(out, "HOST (IP:PORT)")
	if matrixIdx == -1 || primaryIdx == -1 {
		t.Fatalf("expected both matrix and primary headers to be present")
	}
	if matrixIdx > primaryIdx {
		t.Errorf("expected commands matrix (index %d) to be rendered above primary nodes table (index %d)", matrixIdx, primaryIdx)
	}

	totalIdx := strings.Index(out, "Total: 2 registered node(s)")
	if totalIdx == -1 || totalIdx < primaryIdx {
		t.Errorf("expected total count summary after primary table")
	}

	nodeAliasIdx := strings.Index(out, "NODE (ALIAS)")
	subsystemsIdx := strings.Index(out, "SUBSYSTEMS")
	hasValidIndices := nodeAliasIdx != -1 && subsystemsIdx > nodeAliasIdx
	if hasValidIndices && strings.Contains(out[nodeAliasIdx:subsystemsIdx], "ROLE") {
		t.Errorf("expected no ROLE column between NODE (ALIAS) and SUBSYSTEMS")
	}
}

func TestUnifiedNodes_FilterSubsystems(t *testing.T) {
	sampleNodes := []UnifiedFleetNode{
		{
			Alias:      "ssh-only",
			Host:       "10.0.0.1",
			Subsystems: []string{"SSH"},
		},
		{
			Alias:      "cluster-only",
			Host:       "10.0.0.2",
			Subsystems: []string{"Cluster"},
		},
		{
			Alias:      "sc-node",
			Host:       "10.0.0.3",
			Subsystems: []string{"SC", "Cluster"},
		},
	}

	optsSSH := nodesFilterOptions{filterSSH: true}
	resSSH := applyNodesFilter(sampleNodes, optsSSH)
	if len(resSSH) != 1 || resSSH[0].Alias != "ssh-only" {
		t.Errorf("expected 1 SSH node, got %d", len(resSSH))
	}

	optsSC := nodesFilterOptions{filterSC: true}
	resSC := applyNodesFilter(sampleNodes, optsSC)
	if len(resSC) != 1 || resSC[0].Alias != "sc-node" {
		t.Errorf("expected 1 SC node, got %d", len(resSC))
	}
}

func TestUnifiedNodes_JSONFormat(t *testing.T) {
	sampleNodes := []UnifiedFleetNode{
		{
			Alias:      "node-a",
			Role:       "worker",
			Host:       "192.168.1.50",
			Port:       22,
			User:       "root",
			Status:     "ready",
			Subsystems: []string{"SSH", "Cluster", "SC"},
			EnrolledAt: "2026-09-29 10:00:00",
			IsOnline:   true,
		},
	}

	data, err := json.Marshal(sampleNodes)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var parsed []UnifiedFleetNode
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}

	if len(parsed) != 1 || parsed[0].Alias != "node-a" {
		t.Errorf("parsed JSON mismatch: %+v", parsed)
	}
}
