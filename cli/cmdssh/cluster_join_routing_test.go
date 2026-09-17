package cmdssh

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestClusterHelpGroups_ContainsSSHJoin(t *testing.T) {
	if !strings.Contains(constants.CompactCluster, "ssh-join") && !strings.Contains(constants.CompactCluster, "sj") {
		t.Fatalf("expected CompactCluster to mention ssh-join / sj, got: %s", constants.CompactCluster)
	}
}

func TestRunClusterJoinCLI_HelpFlag(t *testing.T) {
	err := RunClusterJoinCLI([]string{"--help"})
	if err != nil {
		t.Fatalf("expected nil error on --help, got: %v", err)
	}

	errShort := RunClusterJoinCLI([]string{"-h"})
	if errShort != nil {
		t.Fatalf("expected nil error on -h, got: %v", errShort)
	}
}

func TestRunClusterJoinCLI_EmptyArgs(t *testing.T) {
	err := RunClusterJoinCLI([]string{})
	if err == nil {
		t.Fatalf("expected validation error on empty args, got nil")
	}
	if !strings.Contains(err.Error(), "missing target host") {
		t.Fatalf("expected missing target host message, got: %v", err)
	}
}
