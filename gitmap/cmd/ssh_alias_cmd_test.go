package cmd

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func Test_saveAliasCommand(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	// Just a simple validation to avoid db requirement during basic run
	err := runSSHAlias(cmd, []string{"192.168.1.1"}, ctx)
	if err == nil {
		t.Errorf("Expected error for missing arguments")
	}
}
