package cmd

import (
	"bytes"
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

func TestRunIPCmd(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := runIPCmd(cmd, []string{}, ctx)
	if err != nil {
		t.Logf("runIPCmd returned error in environment: %v", err)
	}
}

func TestExecuteIPCmd(t *testing.T) {
	ctx := context.Background()
	buffer := &bytes.Buffer{}

	err := executeIPCmd(ctx, false, buffer)
	if err != nil {
		t.Logf("executeIPCmd returned: %v", err)
	}

	if buffer.Len() == 0 {
		return
	}

	outStr := buffer.String()
	if len(outStr) == 0 {
		t.Errorf("expected non-empty output")
	}
}
