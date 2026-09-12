package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var IPCmd = &cobra.Command{
	Use:   "ip",
	Short: "Print local IP",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runIPCmd(cmd, args, cmd.Context())
	},
}

//nolint:revive
func runIPCmd(cmd *cobra.Command, args []string, ctx context.Context) error {
	return executeIPCmd(ctx, true, os.Stdout)
}

func executeIPCmd(ctx context.Context, skipLoopback bool, writer io.Writer) error {
	ipStr, err := GetLocalIP(ctx, skipLoopback, "")
	if err != nil {
		return apperror.New("executeIPCmd", "E_INTERNAL_ERROR", map[string]any{"err": err.Error()})
	}

	if skipLoopback && (ipStr == "127.0.0.1" || ipStr == "::1") {
		return apperror.New("executeIPCmd", "E_INTERNAL_ERROR", map[string]any{"err": "only loopback found"})
	}

	fmt.Fprintln(writer, ipStr)

	return nil
}
