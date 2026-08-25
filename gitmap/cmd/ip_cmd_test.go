package cmd

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func Test_runSSHJoin(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := runSSHJoin(cmd, []string{}, ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = runSSHJoin(cmd, []string{"add"}, ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func Test_runSJRm(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := runSJRm(cmd, []string{"alias1"}, ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = runSJRm(cmd, []string{}, ctx)
	if err == nil {
		t.Errorf("expected error for missing args, got nil")
	}
}
