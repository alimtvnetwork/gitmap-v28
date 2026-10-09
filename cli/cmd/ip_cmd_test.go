package cmd

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/spf13/cobra"
)

func Test_runSSHJoin(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := cmdssh.RunSSHJoin(cmd, []string{}, ctx)
	if err == nil {
		t.Errorf("expected error for empty args, got nil")
	}

	err = cmdssh.RunSSHJoin(cmd, []string{"add"}, ctx)
	if err == nil {
		t.Errorf("expected error for missing add target, got nil")
	}
}

func Test_runSJRm(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := cmdssh.RunSJRm(cmd, []string{"alias1"}, ctx)
	if err != nil {
		t.Logf("runSJRm for unseeded alias1: %v", err)
	}

	err = cmdssh.RunSJRm(cmd, []string{}, ctx)
	if err == nil {
		t.Errorf("expected error for missing args, got nil")
	}
}
