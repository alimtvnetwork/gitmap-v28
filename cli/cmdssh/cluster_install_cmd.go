package cmdssh

import (
	"context"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
)

type clusterInstallOptions struct {
	component string
	target    string
	version   string
	isSudo    bool
	parallel  int
	osType    string
}

// ClusterInstallCmd represents the gitmap cluster install command.
var ClusterInstallCmd = &cobra.Command{
	Use:   "install [component] [target]",
	Short: "Install GitMap on remote cluster nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClusterInstallCLI(args)
	},
}

func isWindowsOS(osType string) bool {
	lower := strings.ToLower(strings.TrimSpace(osType))
	return strings.Contains(lower, "win")
}

func isExplicitVersion(version string) bool {
	clean := strings.TrimSpace(version)
	return clean != "" && clean != "latest"
}

func buildWindowsInstallOneLiner(version string) string {
	base := `powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex"`
	if isExplicitVersion(version) {
		return base + " -Version " + version
	}
	return base
}

func buildUnixInstallOneLiner(version string) string {
	base := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | bash"
	if isExplicitVersion(version) {
		return base + " -s -- --version " + version
	}
	return base
}

// BuildGitmapInstallOneLiner generates the platform-specific one-liner for installing GitMap.
func BuildGitmapInstallOneLiner(osType string, version string) string {
	if isWindowsOS(osType) {
		return buildWindowsInstallOneLiner(version)
	}
	return buildUnixInstallOneLiner(version)
}

func parseVersionPrefix(arg string) (string, int, bool) {
	for _, prefix := range []string{"--version=", "-v="} {
		if strings.HasPrefix(arg, prefix) {
			return strings.TrimPrefix(arg, prefix), 1, true
		}
	}
	return "", 0, false
}

func parseInstallVersionFlag(args []string, idx int) (string, int, bool) {
	if val, c, isMatched := parseVersionPrefix(args[idx]); isMatched {
		return val, c, true
	}
	if (args[idx] == "--version" || args[idx] == "-v") && idx+1 < len(args) {
		return args[idx+1], 2, true
	}
	return "", 0, false
}

func parseParallelPrefix(arg string) (int, int, bool) {
	for _, prefix := range []string{"--parallel=", "-p="} {
		if strings.HasPrefix(arg, prefix) {
			val, err := strconv.Atoi(strings.TrimPrefix(arg, prefix))
			return val, 1, err == nil
		}
	}
	return 0, 0, false
}

func parseInstallParallelFlag(args []string, idx int) (int, int, bool) {
	if val, c, isMatched := parseParallelPrefix(args[idx]); isMatched {
		return val, c, true
	}
	if (args[idx] == "--parallel" || args[idx] == "-p") && idx+1 < len(args) {
		val, err := strconv.Atoi(args[idx+1])
		return val, 2, err == nil
	}
	return 0, 0, false
}

func parseInstallSudoFlag(arg string) (bool, int, bool) {
	if arg == "--sudo" || arg == "-s" || arg == "--sudo=true" {
		return true, 1, true
	}
	if arg == "--sudo=false" || arg == "--no-sudo" {
		return false, 1, true
	}
	return true, 0, false
}

func parseInstallOSFlag(args []string, idx int) (string, int, bool) {
	if strings.HasPrefix(args[idx], "--os=") {
		return strings.TrimPrefix(args[idx], "--os="), 1, true
	}
	if args[idx] == "--os" && idx+1 < len(args) {
		return args[idx+1], 2, true
	}
	return "", 0, false
}

func parseInstallConfigFlag(args []string, idx int, opts *clusterInstallOptions) (int, bool) {
	if isSudo, c, isMatched := parseInstallSudoFlag(args[idx]); isMatched {
		opts.isSudo = isSudo
		return c, true
	}
	if ver, c, isMatched := parseInstallVersionFlag(args, idx); isMatched {
		opts.version = ver
		return c, true
	}
	return 0, false
}

func parseInstallFlagAtIndex(args []string, idx int, opts *clusterInstallOptions) (int, bool) {
	if c, isMatched := parseInstallConfigFlag(args, idx, opts); isMatched {
		return c, true
	}
	if par, c, isMatched := parseInstallParallelFlag(args, idx); isMatched {
		opts.parallel = par
		return c, true
	}
	if osType, c, isMatched := parseInstallOSFlag(args, idx); isMatched {
		opts.osType = osType
		return c, true
	}
	return 0, false
}

func resolveInstallSinglePositional(token string, opts *clusterInstallOptions) {
	if token == "gitmap" {
		opts.component = "gitmap"
		opts.target = "all"
		return
	}
	opts.component = "gitmap"
	opts.target = token
}

func resolveInstallPositional(pos []string, opts *clusterInstallOptions) {
	if len(pos) >= 2 {
		opts.component = pos[0]
		opts.target = pos[1]
		return
	}
	if len(pos) == 1 {
		resolveInstallSinglePositional(pos[0], opts)
	}
}

func newDefaultClusterInstallOptions() *clusterInstallOptions {
	return &clusterInstallOptions{
		component: "gitmap",
		target:    "all",
		version:   "latest",
		isSudo:    true,
		parallel:  4,
		osType:    "linux",
	}
}

func parseClusterInstallArgs(args []string) *clusterInstallOptions {
	opts := newDefaultClusterInstallOptions()
	var pos []string
	idx := 0
	for idx < len(args) {
		if c, isFlag := parseInstallFlagAtIndex(args, idx, opts); isFlag {
			idx += c
			continue
		}
		pos = append(pos, args[idx])
		idx++
	}
	resolveInstallPositional(pos, opts)
	return opts
}

func isClusterInstallHelp(args []string) bool {
	if len(args) == 0 {
		return true
	}
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}
	return false
}

func showClusterInstallHelp() error {
	helptext.Print("cluster-install")
	return nil
}

func executeClusterInstall(ctx context.Context, opts *clusterInstallOptions) error {
	hosts, err := resolveClusterHosts(ctx, opts.target)
	if err != nil {
		return err
	}
	installCmd := BuildGitmapInstallOneLiner(opts.osType, opts.version)
	results := dispatchClusterRunFn(ctx, hosts, installCmd, opts.isSudo, opts.parallel)
	if hasExecutionFailures(results) {
		return apperror.NewExecutionError("one or more cluster nodes failed installation")
	}
	return nil
}

// RunClusterInstallCLI executes GitMap installation across cluster nodes.
func RunClusterInstallCLI(args []string) error {
	if isClusterInstallHelp(args) {
		return showClusterInstallHelp()
	}
	opts := parseClusterInstallArgs(args)
	return executeClusterInstall(context.Background(), opts)
}
