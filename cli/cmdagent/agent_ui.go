package cmdagent

import (
	"github.com/spf13/cobra"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// UIOptions holds configuration flags for the agent web visualizer dashboard.
type UIOptions struct {
	Port     int
	IsBrowse bool
	TaskId   string
	Dir      string
}

var (
	uiOpts = UIOptions{
		Port:     8095,
		IsBrowse: true,
	}

	uiCmd = &cobra.Command{
		Use:     "ui",
		Aliases: []string{"web", "dashboard", "visualizer"},
		Short:   "Launch interactive web visualizer dashboard for agent fleets and tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return toError(RunAgentUI(uiOpts))
		},
	}
)

func init() {
	AgentCmd.AddCommand(uiCmd)
	initUIFlags()
}

func initUIFlags() {
	uiCmd.Flags().IntVarP(&uiOpts.Port, "port", "p", 8095, "Starting HTTP port for visualizer dashboard (auto-scans 50 ports)")
	uiCmd.Flags().BoolVarP(&uiOpts.IsBrowse, "browse", "b", true, "Automatically launch default web browser on startup")
	uiCmd.Flags().StringVarP(&uiOpts.TaskId, "task-id", "t", "", "Target parent task ID or slug to visualize")
	uiCmd.Flags().StringVarP(&uiOpts.Dir, "dir", "d", "", "Agent temporary directory override")
}

// RunAgentUI resolves options and launches the embedded visualizer server.
func RunAgentUI(opts UIOptions) *appfault.AppError {
	resolvedOpts := sanitizeUIOptions(opts)
	tempDir := ResolveAgentTempDir(resolvedOpts.Dir)
	resolvedOpts.Dir = tempDir

	return StartAgentUIServer(resolvedOpts)
}

func sanitizeUIOptions(opts UIOptions) UIOptions {
	hasValidPort := opts.Port > 0 && opts.Port <= 65535
	if !hasValidPort {
		opts.Port = 8095
	}

	return opts
}
