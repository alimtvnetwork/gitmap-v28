package cmdssh

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClusterInit_GenerateJSON_Defaults(t *testing.T) {
	flags, err := parseClusterInitFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error parsing defaults: %v", err)
	}

	cfg := generateClusterConfigJSON(flags)
	if cfg.Control["k8s-m1"] != "192.168.1.10" {
		t.Fatalf("expected default control IP 192.168.1.10, got %s", cfg.Control["k8s-m1"])
	}

	if cfg.Nodes["worker-1"] != "192.168.1.11" || cfg.Nodes["worker-2"] != "192.168.1.12" {
		t.Fatalf("expected default worker IPs, got %v", cfg.Nodes)
	}

	if cfg.User.Name != "ubuntu" {
		t.Fatalf("expected default user ubuntu, got %s", cfg.User.Name)
	}
}

func TestClusterInit_GenerateJSON_Custom(t *testing.T) {
	args := []string{
		"--control", "10.0.0.1",
		"--workers", "10.0.0.2,10.0.0.3,10.0.0.4",
		"--user", "devadmin",
		"--password", "TopSecret456",
		"custom-cluster.json",
	}

	flags, err := parseClusterInitFlags(args)
	if err != nil {
		t.Fatalf("unexpected error parsing custom flags: %v", err)
	}

	if flags.outFile != "custom-cluster.json" {
		t.Fatalf("expected custom-cluster.json, got %s", flags.outFile)
	}

	cfg := generateClusterConfigJSON(flags)
	if cfg.Control["k8s-m1"] != "10.0.0.1" {
		t.Fatalf("expected 10.0.0.1, got %s", cfg.Control["k8s-m1"])
	}

	if len(cfg.Nodes) != 3 {
		t.Fatalf("expected 3 worker nodes, got %d", len(cfg.Nodes))
	}

	if cfg.User.Name != "devadmin" || cfg.User.Password != "TopSecret456" {
		t.Fatalf("expected custom user credentials, got %+v", cfg.User)
	}
}

func TestClusterInit_WriteFile_Hermetic(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "cluster.json")

	flags, _ := parseClusterInitFlags([]string{"--out", outPath})
	cfg := generateClusterConfigJSON(flags)

	// First write should succeed
	if err := writeClusterConfigFile(outPath, cfg, false); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	// Overwrite without force should fail
	if err := writeClusterConfigFile(outPath, cfg, false); err == nil {
		t.Fatalf("expected error when overwriting without force")
	}

	// Overwrite with force should succeed
	if err := writeClusterConfigFile(outPath, cfg, true); err != nil {
		t.Fatalf("expected overwrite with force to succeed, got %v", err)
	}

	// Verify file is readable as valid ClusterConfigJSON
	loaded, loadErr := LoadClusterConfigFile(outPath)
	if loadErr != nil {
		t.Fatalf("failed to load generated config: %v", loadErr)
	}

	if loaded.Control["k8s-m1"] != "192.168.1.10" {
		t.Fatalf("loaded config mismatch")
	}
}
