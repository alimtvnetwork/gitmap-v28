package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dashboard"
)

// runDashboard handles the "dashboard" subcommand.
func runDashboard(args []string) {
	checkHelp("dashboard", args)
	opts, outDir, openFlag := parseDashboardFlags(args)

	fmt.Println(constants.MsgDashCollecting)

	data, err := dashboard.Collect(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDashCollect, err)
		os.Exit(1)
	}

	jsonPath, err := dashboard.WriteJSON(outDir, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDashWriteJSON, jsonPath, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgDashWriteJSON, dashboard.Summary(jsonPath),
		data.Meta.TotalCommits, len(data.Authors))

	htmlPath, err := dashboard.WriteHTML(outDir, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDashWriteHTML, htmlPath, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgDashWriteHTML, dashboard.Summary(htmlPath))
	fmt.Printf(constants.MsgDashGenerated, outDir)

	if openFlag {
		openDashboard(htmlPath)
	}
}

// parseDashboardFlags parses dashboard-specific CLI flags.
func parseDashboardFlags(args []string) (dashboard.CollectOptions, string, bool) {
	fs := flag.NewFlagSet(constants.CmdDashboard, flag.ExitOnError)
	limit := fs.Int("limit", 0, constants.FlagDescDashLimit)
	since := fs.String("since", "", constants.FlagDescDashSince)
	noMerges := fs.Bool("no-merges", false, constants.FlagDescNoMerges)
	outDir := fs.String("out-dir", constants.DashboardOutDir, constants.FlagDescDashOutDir)
	openFlag := fs.Bool("open", false, constants.FlagDescDashOpen)
	fs.Parse(args)

	opts := dashboard.CollectOptions{
		RepoPath: ".",
		Limit:    *limit,
		Since:    *since,
		NoMerges: *noMerges,
	}

	return opts, *outDir, *openFlag
}

// openDashboard opens the HTML file in the default browser.
func openDashboard(path string) {
	fmt.Println(constants.MsgDashOpening)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case constants.OSWindows:
		cmd = exec.Command(constants.CmdWindowsShell, constants.CmdArgSlashC, constants.CmdArgStart, path)
	case constants.OSDarwin:
		cmd = exec.Command(constants.CmdOpen, path)
	default:
		cmd = exec.Command(constants.CmdXdgOpen, path)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not open dashboard in browser: %v\n", err)
	}
}
