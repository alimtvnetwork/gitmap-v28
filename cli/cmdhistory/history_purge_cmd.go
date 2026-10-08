package cmdhistory

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpurge"
)

type purgeFlagCollector struct {
	folder        string
	file          string
	commit        string
	repo          string
	isRelease     bool
	isAutoConfirm bool
	isDryRun      bool
	isJson        bool
	isVerbose     bool
	noBackup      bool
}

// NewHistoryPurgeCmd constructs the 'gitmap history purge' Cobra command.
func NewHistoryPurgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "purge [flags]",
		Aliases: []string{"clean", "p"},
		Short:   "Purge files, folders, or commits from Git history with pre-flight and SplitDB undo journal",
		RunE: func(c *cobra.Command, args []string) error {
			return RunHistoryPurgeCLI(args)
		},
	}
	bindPurgeFlags(cmd)

	return cmd
}

func bindPurgeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("folder", "d", "", "Relative directory path to purge from history")
	cmd.Flags().StringP("file", "f", "", "Relative file path to purge from history")
	cmd.Flags().StringP("commit", "c", "", "Target commit SHA or comma-separated list of commit SHAs")
	cmd.Flags().BoolP("release", "r", false, "Scrub matching artifacts and notes from GitHub releases")
	cmd.Flags().BoolP("yes", "y", false, "Auto-confirm pre-flight prompt without interactive input")
	cmd.Flags().BoolP("dry-run", "n", false, "Preview affected commits without mutating Git objects")
	cmd.Flags().String("repo", ".", "Target repository directory")
	cmd.Flags().BoolP("json", "j", false, "Output machine-readable JSON envelope")
	cmd.Flags().BoolP("verbose", "v", false, "Show detailed Git plumbing execution logs")
	cmd.Flags().Bool("no-backup", false, "Skip staging blobs in temporary directory backup")
}

// NewHistoryUndoCmd constructs the 'gitmap history undo' Cobra command.
func NewHistoryUndoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "undo [operation-id] [flags]",
		Aliases: []string{"restore", "u"},
		Short:   "Undo a previous history purge operation and restore original commit history",
		RunE: func(c *cobra.Command, args []string) error {
			return RunHistoryUndoCLI(args)
		},
	}
	bindUndoFlags(cmd)

	return cmd
}

func bindUndoFlags(cmd *cobra.Command) {
	cmd.Flags().String("repo", ".", "Target repository directory")
	cmd.Flags().BoolP("json", "j", false, "Output machine-readable JSON envelope")
	cmd.Flags().BoolP("verbose", "v", false, "Show detailed restoration execution logs")
}

// RunHistoryPurgeCLI is the main entry point for history purge CLI invocations.
func RunHistoryPurgeCLI(args []string) error {
	opts, err := parsePurgeCLIOptions(args)
	if err != nil {
		return err
	}

	report, err := cmdpurge.ExecutePurgeEngineExtended(opts)
	if err != nil {
		return err
	}
	if opts.IsJson && report != nil {
		return printPurgeJsonEnvelope(report)
	}

	return nil
}

func parsePurgeCLIOptions(args []string) (cmdpurge.PurgeExtendedOptions, error) {
	collector := parsePurgeArgsLoop(args)
	tType, tPath, tCommits, err := resolvePurgeTarget(collector)
	if err != nil {
		return cmdpurge.PurgeExtendedOptions{}, err
	}

	return cmdpurge.PurgeExtendedOptions{
		PurgeOptions: cmdpurge.PurgeOptions{
			TargetType:    tType,
			TargetPath:    tPath,
			TargetCommits: tCommits,
			IsAutoConfirm: collector.isAutoConfirm,
			IsDryRun:      collector.isDryRun,
		},
		RepoDir:   collector.repo,
		IsRelease: collector.isRelease,
		IsJson:    collector.isJson,
		IsVerbose: collector.isVerbose,
		NoBackup:  collector.noBackup,
	}, nil
}

func parsePurgeArgsLoop(args []string) *purgeFlagCollector {
	col := &purgeFlagCollector{repo: "."}
	for i := 0; i < len(args); i++ {
		a := args[i]
		i = consumePurgeFlag(a, args, i, col)
	}

	return col
}

func consumePurgeFlag(arg string, args []string, idx int, col *purgeFlagCollector) int {
	switch {
	case (arg == "--folder" || arg == "-d") && idx+1 < len(args):
		idx++
		col.folder = args[idx]
	case (arg == "--file" || arg == "-f") && idx+1 < len(args):
		idx++
		col.file = args[idx]
	case (arg == "--commit" || arg == "-c") && idx+1 < len(args):
		idx++
		col.commit = args[idx]
	case (arg == "--repo") && idx+1 < len(args):
		idx++
		col.repo = args[idx]
	default:
		consumePurgeBooleanFlag(arg, col)
	}

	return idx
}

func consumePurgeBooleanFlag(arg string, col *purgeFlagCollector) {
	switch arg {
	case "-y", "--yes", "--confirm":
		col.isAutoConfirm = true
	case "-r", "--release":
		col.isRelease = true
	case "-n", "--dry-run":
		col.isDryRun = true
	case "-j", "--json":
		col.isJson = true
	case "-v", "--verbose":
		col.isVerbose = true
	case "--no-backup":
		col.noBackup = true
	}
}

func resolvePurgeTarget(col *purgeFlagCollector) (cmdpurge.TargetType, string, []string, error) {
	if len(col.folder) > 0 {
		return cmdpurge.TargetTypeFolder, col.folder, nil, nil
	}
	if len(col.file) > 0 {
		return cmdpurge.TargetTypeFile, col.file, nil, nil
	}
	if len(col.commit) > 0 {
		commits := strings.Split(col.commit, ",")
		return cmdpurge.TargetTypeCommit, col.commit, commits, nil
	}

	return "", "", nil, apperror.NewSimple("ERR_PURGE_NO_TARGET", "no target specified; pass --folder <path>, --file <path>, or --commit <sha>")
}

// RunHistoryUndoCLI is the main entry point for history undo CLI invocations.
func RunHistoryUndoCLI(args []string) error {
	opId, restArgs := parseUndoCLIOpId(args)

	return cmdpurge.RunPurgeUndo(opId, restArgs)
}

func parseUndoCLIOpId(args []string) (int64, []string) {
	if len(args) == 0 {
		return 0, args
	}

	first := args[0]
	if strings.HasPrefix(first, "-") {
		return 0, args
	}

	id, err := strconv.ParseInt(first, 10, 64)
	if err == nil {
		return id, args[1:]
	}

	return 0, args
}

func printPurgeJsonEnvelope(report *cmdpurge.PurgeExecutionReport) error {
	envelope := map[string]interface{}{
		"version":   "2.0",
		"status":    "success",
		"command":   "history purge",
		"timestamp": time.Now().Unix(),
		"data":      report,
		"error":     nil,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(envelope)
}
