package cmd

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunSSHAlias(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := runSSHAlias(cmd, []string{"192.168.1.1", "my-alias"}, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = runSSHAlias(cmd, []string{"192.168.1.1"}, ctx)
	if err == nil {
		t.Errorf("Expected error for missing arguments")
	}
}
