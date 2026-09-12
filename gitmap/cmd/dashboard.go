package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dashboard"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func openDashboardIfRequested(outDir string, openFlag bool) {
	if openFlag {
		openDashboard(filepath.Join(outDir, constants.DashboardHTMLFile))
	}
}

// runDashboard handles the "dashboard" subcommand.
func runDashboard(args []string) error {
	checkHelp("dashboard", args)
	opts, outDir, openFlag := parseDashboardFlags(args)
	fmt.Println(constants.MsgDashCollecting)
	data, appErr := collectDashboardData(opts)
	if appErr != nil {
		return appErr
	}

	if appErr := emitDashboardOutputs(outDir, data); appErr != nil {
		return appErr
	}

	openDashboardIfRequested(outDir, openFlag)

	return nil
}

func collectDashboardData(opts dashboard.CollectOptions) (model.DashboardData, *apperror.AppError) {
	data, err := dashboard.Collect(opts)
	if err != nil {
		return model.DashboardData{}, apperror.WrapSimple(err, constants.ErrDashCollect)
	}

	return data, nil
}

func emitDashboardOutputs(outDir string, data model.DashboardData) *apperror.AppError {
	if appErr := writeDashboardJSON(outDir, data); appErr != nil {
		return appErr
	}

	if appErr := writeDashboardHTML(outDir, data); appErr != nil {
		return appErr
	}

	fmt.Printf(constants.MsgDashGenerated, outDir)

	return nil
}

func writeDashboardJSON(outDir string, data model.DashboardData) *apperror.AppError {
	jsonPath, err := dashboard.WriteJSON(outDir, data)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDashWriteJSON)
	}

	fmt.Printf(constants.MsgDashWriteJSON, dashboard.Summary(jsonPath),
		data.Meta.TotalCommits, len(data.Authors))

	return nil
}

func writeDashboardHTML(outDir string, data model.DashboardData) *apperror.AppError {
	htmlPath, err := dashboard.WriteHTML(outDir, data)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrDashWriteHTML)
	}

	fmt.Printf(constants.MsgDashWriteHTML, dashboard.Summary(htmlPath))

	return nil
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
