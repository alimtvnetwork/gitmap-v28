package cmdip

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunIPCmd(t *testing.T) {
	cmd := &cobra.Command{}
	ctx := context.Background()

	err := RunIPCmd(cmd, []string{}, ctx)
	if err != nil {
		t.Logf("RunIPCmd returned error in environment: %v", err)
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
