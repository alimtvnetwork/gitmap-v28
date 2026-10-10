package cmdpull

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
	"os"
	"strconv"
	"strings"
)

// NormalizePullArgs converts positional pull-all and table arguments into flags.
func NormalizePullArgs(args []string) []string {
	normalized := make([]string, 0, len(args)+1)
	for i := 0; i < len(args); i++ {
		if isPullAllTableSeq(args, i) {
			normalized = append(normalized, "--all", "--status")
			i++
			continue
		}
		if isHandShortFlagWithVal(args, i) {
			normalized = append(normalized, "--hand", args[i+1])
			i++
			continue
		}
		normalized = appendNormalizedPullToken(normalized, args[i])
	}

	return reorderPullFlags(normalized)
}

func isHandShortFlagWithVal(args []string, i int) bool {
	if args[i] != "-h" || i+1 >= len(args) {
		return false
	}
	_, err := strconv.Atoi(args[i+1])
	return err == nil
}

func reorderPullFlags(args []string) []string {
	flags := make([]string, 0, len(args))
	pos := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isPullFlagTakingValue(arg) {
			flags = appendValuedFlag(flags, args, &i)
			continue
		}
		flags, pos = categorizePullArg(arg, flags, pos)
	}

	return append(flags, pos...)
}

func categorizePullArg(arg string, flags, pos []string) ([]string, []string) {
	if strings.HasPrefix(arg, "-") {
		return append(flags, arg), pos
	}

	return flags, append(pos, arg)
}

func appendValuedFlag(flags []string, args []string, idx *int) []string {
	flags = append(flags, args[*idx])
	if *idx+1 < len(args) {
		*idx++
		flags = append(flags, args[*idx])
	}
	return flags
}

func isPullFlagTakingValue(arg string) bool {
	if strings.Contains(arg, "=") {
		return false
	}

	return arg == "-g" || arg == "--group" || arg == "-p" || arg == "--parallel" ||
		arg == "-w" || arg == "--w" || arg == "--worker" || arg == "--workers" ||
		arg == "--hand" || arg == "--hands" || arg == "--h" || arg == "--concurrency"
}

func isPullAllTableSeq(args []string, i int) bool {
	isAllToken := args[i] == "all" || args[i] == "--all" || args[i] == "pa" || args[i] == "ta" || args[i] == "pull-all"

	return isAllToken && i+1 < len(args) && args[i+1] == "table"
}

func appendNormalizedPullToken(normalized []string, token string) []string {
	if isTableToken(token) {
		return append(normalized, "--all", "--status")
	}
	if isSSHFleetToken(token) {
		return append(normalized, "--all", "--ssh")
	}
	if isAllToken(token) {
		return append(normalized, "--all")
	}

	return append(normalized, token)
}

func isTableToken(token string) bool {
	lower := strings.ToLower(token)
	return lower == "pat" || lower == "pull-all-table"
}

func isAllToken(token string) bool {
	lower := strings.ToLower(token)
	return lower == "all" || lower == "pa" || lower == "ta" || lower == "pull-all"
}

func initPullFlagSet() (*flag.FlagSet, *pullFlagHolders) {
	fs := flag.NewFlagSet(constants.CmdPull, flag.ExitOnError)
	h := &pullFlagHolders{}
	registerPullCoreFlags(fs, h)
	registerPullRemediationFlags(fs, h)
	registerPullOutputFlags(fs, h)

	return fs, h
}

func registerPullCoreFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	registerPullBasicFlags(fs, h)
	registerPullWorkerFlags(fs, h)
	registerPullConcurrencyPresetFlags(fs, h)
}

func registerPullBasicFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.vFlag = fs.Bool("verbose", false, constants.FlagDescVerbose)
	h.gFlag = fs.String("group", "", constants.FlagDescGroup)
	h.aFlag = fs.Bool("all", false, constants.FlagDescAll)
	h.sFlag = fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	h.pFlag = fs.Int("parallel", 0, constants.FlagDescPullParallel)
	h.oFlag = fs.Bool("only-available", false, constants.FlagDescPullOnlyAvailable)
	h.rawFlag = fs.Bool("raw", false, constants.FlagDescPullRaw)
	h.sshFlag = fs.Bool("ssh", false, "Pull using SSH transport")
	h.httpsFlag = fs.Bool("https", false, "Pull using HTTPS transport")
	fs.StringVar(h.gFlag, "g", "", constants.FlagDescGroup)
	fs.IntVar(h.pFlag, "p", 0, constants.FlagDescPullParallel)
	fs.IntVar(h.pFlag, "concurrency", 0, constants.FlagDescPullParallel)
}

func registerPullWorkerFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.wFlag = fs.Int("w", 0, "Worker pool size")
	fs.IntVar(h.wFlag, "worker", 0, "Worker pool size")
	fs.IntVar(h.wFlag, "workers", 0, "Worker pool size")
	h.handFlag = fs.Int("hand", 0, "Hands per worker")
	fs.IntVar(h.handFlag, "hands", 0, "Hands per worker")
	fs.IntVar(h.handFlag, "h", 0, "Hands per worker")
	h.wwohFlag = fs.Bool("wwoh", false, "Worker with one hand")
	fs.BoolVar(h.wwohFlag, "worker-with-one-hand", false, "Worker with one hand")
}

func registerPullConcurrencyPresetFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.autoScaleFlag = fs.Bool("auto-scale", false, "Auto-scale pull concurrency based on CPU availability")
	h.highPerfFlag = fs.Bool("high-perf", false, "High-performance pull concurrency preset for low CPU pressure")
	fs.BoolVar(h.highPerfFlag, "turbo", false, "High-performance pull concurrency preset for low CPU pressure")
	h.lowCPUFlag = fs.Bool("low-cpu", false, "Conservative pull concurrency preset for high CPU pressure")
	fs.BoolVar(h.lowCPUFlag, "conservative", false, "Conservative pull concurrency preset for high CPU pressure")
}

func registerPullRemediationFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.fixFlag = fs.Bool("fix", false, "Auto-remediate dirty repos")
	h.yFlag = fs.Bool("yes", false, "Remediate without prompt")
	h.noFixFlag = fs.Bool("no-fix", false, "Skip remediation prompt")
	fs.BoolVar(h.yFlag, "y", false, "Remediate without prompt")
}

func registerPullOutputFlags(fs *flag.FlagSet, h *pullFlagHolders) {
	h.statusFlag = fs.Bool("status", false, constants.FlagDescPullStatus)
	fs.BoolVar(h.statusFlag, "table", false, constants.FlagDescPullStatus)
	fs.BoolVar(h.statusFlag, "status-table", false, constants.FlagDescPullStatus)
	h.jsonFlag = fs.Bool("json", false, constants.FlagDescPullJSON)
	h.probeFlag = fs.Bool("probe", false, constants.FlagDescPullProbe)
	fs.BoolVar(h.probeFlag, "probe-repos", false, constants.FlagDescPullProbe)
}

func initBasePullOptions(h *pullFlagHolders) pullOptions {
	return pullOptions{
		group: *h.gFlag, all: *h.aFlag, verbose: *h.vFlag,
		stopOnFail: *h.sFlag, parallel: *h.pFlag,
		workers: *h.wFlag, hands: *h.handFlag, isWWOH: *h.wwohFlag,
		isAutoScale: *h.autoScaleFlag, isHighPerf: *h.highPerfFlag, isLowCPU: *h.lowCPUFlag,
		onlyAvailable: *h.oFlag,
		autoFix:       *h.fixFlag, yes: *h.yFlag, noFix: *h.noFixFlag,
		isRaw: *h.rawFlag, useSSH: *h.sshFlag, useHTTPS: *h.httpsFlag,
		showStatus: *h.statusFlag, isJSON: *h.jsonFlag, isProbe: *h.probeFlag,
	}
}

func resolvePullCPUProfile(h *pullFlagHolders) cloneconcurrency.CPUProfile {
	if *h.highPerfFlag {
		return cloneconcurrency.CPUProfileHighPerf
	}
	if *h.lowCPUFlag {
		return cloneconcurrency.CPUProfileLowCPU
	}
	return cloneconcurrency.CPUProfileAuto
}

func buildPullOptions(h *pullFlagHolders) pullOptions {
	opts := initBasePullOptions(h)
	profile := resolvePullCPUProfile(h)
	if opts.parallel <= 0 {
		opts.parallel = cloneconcurrency.ResolveAdaptivePullConcurrency(profile, 0)
	}

	return opts
}

func parsePullFlags(args []string) pullOptions {
	fs, h := initPullFlagSet()
	fs.Parse(args)
	opts := buildPullOptions(h)
	if fs.NArg() > 0 {
		opts.slug = fs.Arg(0)
	}

	return opts
}

func initVerboseLog() {
	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}
	log.Close()
}

func ExtractTransportFlags(args []string) (bool, bool, []string) {
	var useSSH, useHTTPS bool
	rest := make([]string, 0, len(args))
	for _, a := range args {
		useSSH, useHTTPS, rest = classifyTransportToken(a, useSSH, useHTTPS, rest)
	}

	return useSSH, useHTTPS, rest
}

func classifyTransportToken(a string, useSSH, useHTTPS bool, rest []string) (bool, bool, []string) {
	switch strings.ToLower(a) {
	case "--ssh", "-ssh", "--sh", "-sh", "ssh":
		return true, useHTTPS, rest
	case "--https", "-https", "--ht", "-ht", "https":
		return useSSH, true, rest
	default:
		return useSSH, useHTTPS, append(rest, a)
	}
}
