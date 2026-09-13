package cmd

import (
	"os"
	"testing"
)

func TestStorageCmdDriveReport(t *testing.T) {
	err := RunStorageCmd([]string{})
	if err != nil {
		t.Fatalf("RunStorageCmd without args returned error: %v", err)
	}
}

func TestStorageCmdList(t *testing.T) {
	err := RunStorageCmd([]string{"ls"})
	if err != nil {
		t.Fatalf("RunStorageCmd ls returned error: %v", err)
	}
}

func TestGetDiskSpaceMetrics(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd failed: %v", err)
	}

	metrics, err := getDiskSpaceMetrics(cwd)
	if err != nil {
		t.Fatalf("getDiskSpaceMetrics failed: %v", err)
	}

	if metrics.TotalBytes == 0 {
		t.Errorf("expected total bytes > 0, got 0")
	}

	if metrics.DrivePath == "" {
		t.Errorf("expected non-empty drive path")
	}
}

func TestIsStorageListSubcommand(t *testing.T) {
	if !isStorageListSubcommand("ls") {
		t.Errorf("expected ls to be storage list subcommand")
	}

	if !isStorageListSubcommand("list") {
		t.Errorf("expected list to be storage list subcommand")
	}

	if !isStorageListSubcommand("db") {
		t.Errorf("expected db to be storage list subcommand")
	}

	if isStorageListSubcommand("unknown") {
		t.Errorf("expected unknown NOT to be storage list subcommand")
	}
}
