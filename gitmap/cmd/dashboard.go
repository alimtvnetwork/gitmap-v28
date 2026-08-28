package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dashboard"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runDashboard handles the "dashboard" subcommand.
func runDashboard(args []string) error {
	checkHelp("dashboard", args)
	opts, outDir, openFlag := parseDashboardFlags(args)
	fmt.Println(constants.MsgDashCollecting)
	data := collectDashboardData(opts)
	emitDashboardOutputs(outDir, data)
	if openFlag {
		openDashboard(filepath.Join(outDir, constants.DashboardHTMLFile))
	}
	return nil
}

func collectDashboardData(opts dashboard.CollectOptions) model.DashboardData {
	data, err := dashboard.Collect(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDashCollect, err)
		os.Exit(1)
	}
	return data
}

func emitDashboardOutputs(outDir string, data model.DashboardData) {
	writeDashboardJSON(outDir, data)
	writeDashboardHTML(outDir, data)
	fmt.Printf(constants.MsgDashGenerated, outDir)
}

func writeDashboardJSON(outDir string, data model.DashboardData) {
	jsonPath, err := dashboard.WriteJSON(outDir, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDashWriteJSON, jsonPath, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgDashWriteJSON, dashboard.Summary(jsonPath),
		data.Meta.TotalCommits, len(data.Authors))
}

func writeDashboardHTML(outDir string, data model.DashboardData) {
	htmlPath, err := dashboard.WriteHTML(outDir, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrDashWriteHTML, htmlPath, err)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgDashWriteHTML, dashboard.Summary(htmlPath))
}

type dashFlagSet struct {
	limit    *int
	since    *string
	noMerges *bool
	outDir   *string
	openFlag *bool
	recent   *bool
}

func setupDashboardFlagSet(fs *flag.FlagSet) dashFlagSet {
	return dashFlagSet{
		limit:    fs.Int("limit", 0, constants.FlagDescDashLimit),
		since:    fs.String("since", "", constants.FlagDescDashSince),
		noMerges: fs.Bool("no-merges", false, constants.FlagDescNoMerges),
		outDir:   fs.String("out-dir", constants.DashboardOutDir, constants.FlagDescDashOutDir),
		openFlag: fs.Bool("open", false, constants.FlagDescDashOpen),
		recent:   fs.Bool(constants.FlagRecent, false, constants.FlagDescDashRecent),
	}
}

// parseDashboardFlags parses dashboard-specific CLI flags.
func parseDashboardFlags(args []string) (dashboard.CollectOptions, string, bool) {
	fs := flag.NewFlagSet(constants.CmdDashboard, flag.ExitOnError)
	flags := setupDashboardFlagSet(fs)
	fs.Parse(args)
	opts := buildDashboardCollectOpts(flags)
	return opts, *flags.outDir, *flags.openFlag
}

func buildDashboardCollectOpts(f dashFlagSet) dashboard.CollectOptions {
	return dashboard.CollectOptions{
		RepoPath: ".",
		Limit:    *f.limit,
		Since:    *f.since,
		NoMerges: *f.noMerges,
		Recent:   *f.recent,
	}
}

// openDashboard opens the HTML file in the default browser.
func openDashboard(path string) {
	fmt.Println(constants.MsgDashOpening)
	cmd := buildBrowserCmd(path)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not open dashboard in browser: %v\n", err)
	}
}

func buildBrowserCmd(path string) *exec.Cmd {
	switch runtime.GOOS {
	case constants.OSWindows:
		return exec.Command(constants.CmdWindowsShell, constants.CmdArgSlashC, constants.CmdArgStart, path)
	case constants.OSDarwin:
		return exec.Command(constants.CmdOpen, path)
	default:
		return exec.Command(constants.CmdXdgOpen, path)
	}
}
