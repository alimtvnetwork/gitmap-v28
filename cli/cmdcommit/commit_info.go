package cmdcommit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

var (
	// ErrHelpRequested is returned when help flags are passed.
	ErrHelpRequested = errors.New("help requested")

	// ErrMissingCommitQuery is returned when no commit sha or reference is provided.
	ErrMissingCommitQuery = errors.New("missing commit query")
)

// RunCommitInfo is the primary entry point for gitmap commit-info, ci, and commit info.
func RunCommitInfo(args []string) error {
	cleanArgs := stripLeadingInfoSubcommand(args)
	opts, err := parseCommitInfoFlags(cleanArgs)
	if errors.Is(err, ErrHelpRequested) {
		renderCommitInfoHelp(os.Stdout)

		return nil
	}
	if err != nil {
		return apperror.WrapSimple(err, "parseCommitInfoFlags")
	}

	return dispatchCommitInfo(opts)
}

func stripLeadingInfoSubcommand(args []string) []string {
	if len(args) > 0 && strings.EqualFold(args[0], "info") {
		return args[1:]
	}

	return args
}

func dispatchCommitInfo(opts CommitInfoOptions) error {
	if len(opts.Query) == 0 {
		renderCommitInfoHelp(os.Stderr)

		return apperror.WrapSimple(ErrMissingCommitQuery, "commit sha or reference is required")
	}

	return executeCommitInfo(opts)
}

func executeCommitInfo(opts CommitInfoOptions) error {
	repoPath, err := resolveCommitInfoRepoRoot(opts.RepoPath)
	if err != nil {
		return err
	}

	fullSha, err := ResolveUniversalCommitSha(repoPath, opts.Query)
	if err != nil {
		return err
	}

	report, err := BuildCommitInfoReport(repoPath, fullSha)
	if err != nil {
		return err
	}

	return RenderCommitInfoReport(os.Stdout, report, opts.IsJsonFormat)
}

func resolveCommitInfoRepoRoot(repoFlag string) (string, error) {
	if len(repoFlag) > 0 {
		return checkRepoRoot(repoFlag)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", apperror.WrapSimple(err, "os.Getwd")
	}

	return checkRepoRoot(cwd)
}

func checkRepoRoot(dir string) (string, error) {
	res := gitutil.RepoRoot(dir)
	if res.IsFailure() {
		return "", apperror.WrapSimple(res.AppError(), "find repo root for: "+dir)
	}

	return res.Data, nil
}

func parseCommitInfoFlags(args []string) (CommitInfoOptions, error) {
	opts := CommitInfoOptions{}
	for i := 0; i < len(args); i++ {
		consumed, err := evaluateInfoFlagToken(args, i, &opts)
		if err != nil {
			return opts, err
		}
		i += consumed
	}

	return opts, nil
}

func evaluateInfoFlagToken(args []string, idx int, opts *CommitInfoOptions) (int, error) {
	arg := args[idx]
	if isHelpToken(arg) {
		return 0, ErrHelpRequested
	}
	if consumed, hasHandled := handleInfoOptionFlags(args, idx, opts); hasHandled {
		return consumed, nil
	}
	if len(opts.Query) == 0 && !strings.HasPrefix(arg, "-") {
		opts.Query = arg
	}

	return 0, nil
}

func handleInfoOptionFlags(args []string, idx int, opts *CommitInfoOptions) (int, bool) {
	arg := args[idx]
	if arg == "--json" {
		opts.IsJsonFormat = true
		return 0, true
	}
	if arg == "-v" || arg == "--verbose" {
		opts.IsVerbose = true
		return 0, true
	}

	return handleRepoFlagToken(args, idx, opts)
}

func handleRepoFlagToken(args []string, idx int, opts *CommitInfoOptions) (int, bool) {
	arg := args[idx]
	if strings.HasPrefix(arg, "--repo=") {
		opts.RepoPath = strings.TrimPrefix(arg, "--repo=")
		return 0, true
	}
	if (arg == "--repo" || arg == "-r") && idx+1 < len(args) {
		opts.RepoPath = args[idx+1]
		return 1, true
	}

	return 0, false
}

func isHelpToken(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func renderCommitInfoHelp(w io.Writer) {
	c := constants.ColorCyan
	r := constants.ColorReset
	b := constants.ColorBold
	fmt.Fprintf(w, "%sgitmap commit info%s (aliases: %sci%s, %scommit-info%s)\n\n", b, r, c, r, c, r)
	fmt.Fprintf(w, "Inspect commit details, AI provenance, cross-DB telemetry, and file mutations.\n\n")
	fmt.Fprintf(w, "%sUsage:%s\n", b, r)
	fmt.Fprintf(w, "  gitmap commit info <sha|ref> [flags]\n")
	fmt.Fprintf(w, "  gitmap ci <sha|ref> [flags]\n")
	fmt.Fprintf(w, "  gitmap commit-info <sha|ref> [flags]\n\n")
	fmt.Fprintf(w, "%sFlags:%s\n", b, r)
	fmt.Fprintf(w, "  --json          Output machine-readable JSON\n")
	fmt.Fprintf(w, "  --repo, -r      Target repository path\n")
	fmt.Fprintf(w, "  -v, --verbose   Verbose diagnostic output\n")
	fmt.Fprintf(w, "  -h, --help      Show this help menu\n")
}
