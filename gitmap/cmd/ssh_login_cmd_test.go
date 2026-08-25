package cmd

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunSSHLogin(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := runSSHLogin(cmd, []string{"my-target"}, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = runSSHLogin(cmd, []string{}, ctx)
	if err == nil {
		t.Errorf("Expected error for missing arguments")
	}
}
