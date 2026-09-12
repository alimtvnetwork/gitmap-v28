package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func Test_runSJHistory(t *testing.T) {
	cmd := &cobra.Command{}
	_ = cmd
}

func Test_printSJHistory(t *testing.T) {
	_ = context.Background()
	_ = bytes.NewBuffer(nil)
}
