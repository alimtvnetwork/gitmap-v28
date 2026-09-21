package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/committransfer"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/movemerge"
	"github.com/alimtvnetwork/gitmap-v28/cli/prdb"
)

// commitTransferSpec describes one of the three commit-transfer commands.
type commitTransferSpec struct {
	Name      string // e.g. constants.CmdCommitLeft
	LogPrefix string // e.g. constants.LogPrefixCommitLeft
}

// runCommitTransfer is the single entry point for commit-left,
// commit-right, and commit-both.
//
// Phase 1 (v3.76.0): commit-right.
// Phase 2 (v3.102.0): commit-left wired through committransfer.RunLeft.
// Phase 3 (v3.102.0): commit-both wired through committransfer.RunBoth.
func runCommitTransfer(spec commitTransferSpec, args []string) error {
	checkHelp(spec.Name, args)
	executeCommitTransfer(spec, args)

	return nil
}

// executeCommitTransfer parses flags, resolves endpoints, and dispatches
// to the directional runner that matches spec.Name. Kept separate from
// runCommitTransfer so the help-check stays a one-liner and the dispatch
// table reads top-to-bottom without nested branching.
func executeCommitTransfer(spec commitTransferSpec, args []string) {
	opts, positional := parseCommitTransferArgs(spec, args)
	if len(positional) != 2 {
		fmt.Fprintf(os.Stderr, constants.ErrCTArgCountFmt, spec.Name, len(positional))
		fmt.Fprintf(os.Stderr, constants.MsgCTUsageFmt, spec.Name, spec.Name)
		cliexit.HandleError(nil, 1)
	}

	if opts.Interleave && spec.Name != constants.CmdCommitBoth {
		fmt.Fprintf(os.Stderr,
			"%s --interleave is only valid for commit-both (got %s)\n",
			opts.LogPrefix, spec.Name)
		cliexit.HandleError(nil, 2)
	}

	cfg, _ := config.LoadFromFile(constants.DefaultConfigPath)
	opts.Message.KeepUrl = cfg.CommitReplayKeepUrl
	opts.Message.Templates = cfg.CommitReplayTemplates
	left, right, resolveErr := resolveCommitEndpoints(positional[0], positional[1], opts)
	if resolveErr != nil {
		fmt.Fprintf(os.Stderr, "%s endpoint resolve failed: %v\n", opts.LogPrefix, resolveErr)
		cliexit.HandleError(nil, 1)
	}

	opts.Message.SourceDisplayName = pickSourceDisplayName(spec.Name, left, right, opts.Message.KeepUrl)
	if err := dispatchDirection(spec.Name, left.WorkingDir, right.WorkingDir, opts); err != nil {
		fmt.Fprintf(os.Stderr, "%s replay failed: %v\n", opts.LogPrefix, err)
		cliexit.HandleError(nil, 1)
	}
}

// dispatchDirection routes to the right RunX function. LEFT/RIGHT
// positional ordering matches the spec — `commit-left LEFT RIGHT`
// writes commits onto LEFT (using RIGHT as source). For commit-both,
// --interleave switches to the author-date merged stream.
func dispatchDirection(name, leftDir, rightDir string, opts committransfer.Options) error {
	switch name {
	case constants.CmdCommitLeft:
		return committransfer.RunLeft(leftDir, rightDir, opts)
	case constants.CmdCommitBoth:
		if opts.Interleave {
			return committransfer.RunBothInterleaved(leftDir, rightDir, opts)
		}

		return committransfer.RunBoth(leftDir, rightDir, opts)
	default:
		// commit-right (and any future direction defaulting to L→R).
		return committransfer.RunRight(leftDir, rightDir, opts)
	}
}

// pickSourceDisplayName labels the provenance footer with the side
// that's being read from. commit-left reads from RIGHT; commit-right
// reads from LEFT; commit-both reads from both — we pick LEFT as the
// canonical label because the RunBoth implementation uses (L→R) first.
func pickSourceDisplayName(name string, left, right movemerge.Endpoint, keepUrl bool) string {
	display := left.DisplayName
	if name == constants.CmdCommitLeft {
		display = right.DisplayName
	}

	if keepUrl {
		return display
	}

	return sanitizeURLBaseName(display)
}

func sanitizeURLBaseName(base string) string {
	if strings.Contains(base, "://") || strings.HasPrefix(base, "git@") {
		parts := strings.Split(base, "/")
		base = parts[len(parts)-1]

		return strings.TrimSuffix(base, ".git")
	}

	return base
}

//nolint:unused
func _deprecated_pickSourceDisplayName(name string, left, right movemerge.Endpoint) string {
	if name == constants.CmdCommitLeft {
		return right.DisplayName
	}

	return left.DisplayName
}

func provisionDestTarget(raw string, isLocal bool) string {
	prov, err := EnsureOrProvisionDestinationRepo(raw, isLocal)
	if err == nil {
		return prov
	}

	return raw
}

func resolveEndpointsResolved(leftTarget, rightTarget string) (movemerge.Endpoint, movemerge.Endpoint, error) {
	mmOpts := movemerge.Options{}
	resolvedLeft := resolveEndpointString(leftTarget)
	left, err := movemerge.ResolveEndpoint(resolvedLeft, true, mmOpts)
	if err != nil {
		return left, movemerge.Endpoint{}, err
	}

	resolvedRight := resolveEndpointString(rightTarget)
	right, err := movemerge.ResolveEndpoint(resolvedRight, false, mmOpts)

	return left, right, err
}

// resolveCommitEndpoints reuses the merge-* endpoint resolver, auto-provisioning
// missing destination endpoints if needed.
func resolveCommitEndpoints(leftRaw, rightRaw string, opts committransfer.Options,
) (movemerge.Endpoint, movemerge.Endpoint, error) {
	leftTarget := leftRaw
	rightTarget := rightRaw
	isLeftDest := opts.CommandName == constants.CmdCommitLeft
	if isLeftDest {
		leftTarget = provisionDestTarget(leftRaw, opts.IsLocal)
	}
	if !isLeftDest {
		rightTarget = provisionDestTarget(rightRaw, opts.IsLocal)
	}

	return resolveEndpointsResolved(leftTarget, rightTarget)
}

func defaultPRMode(cmdName string) string {
	if cmdName == constants.CmdPR || cmdName == constants.CmdPullRequest {
		return "merges"
	}

	return ""
}

// parseCommitTransferArgs builds the Options struct + positional args.
// One function per concern would be cleaner, but the flag.FlagSet API
// keeps us under the per-function line cap as long as helpers extract
// the message-policy block.
func parseCommitTransferArgs(spec commitTransferSpec, args []string,
) (committransfer.Options, []string) {
	fs := flag.NewFlagSet(spec.Name, flag.ExitOnError)
	opts := committransfer.Options{
		CommandName: spec.Name, LogPrefix: spec.LogPrefix,
		IncludeMerges: true, // v6.0.0 default — merge commits preserved
		PRMode:        defaultPRMode(spec.Name),
		Message: committransfer.MessagePolicy{
			DropPatterns: committransfer.DefaultDropPatterns,
			Conventional: true, Provenance: true,
			CommandName: spec.Name,
		},
	}

	registerCommitTransferBools(fs, &opts)
	registerCommitTransferStrings(fs, &opts)
	fs.Parse(reorderFlagsBeforeArgs(args))

	return opts, fs.Args()
}

// registerCommitTransferBools wires every boolean flag from spec §8.
func registerCommitTransferBools(fs *flag.FlagSet, opts *committransfer.Options) {
	fs.BoolVar(&opts.Yes, constants.FlagCTYes, false, constants.FlagDescCTYes)
	fs.BoolVar(&opts.Yes, "y", false, constants.FlagDescCTYes)
	fs.BoolVar(&opts.DryRun, constants.FlagCTDryRun, false, constants.FlagDescCTDryRun)
	fs.BoolVar(&opts.NoPush, constants.FlagCTNoPush, false, constants.FlagDescCTNoPush)
	fs.BoolVar(&opts.NoCommit, constants.FlagCTNoCommit, false, constants.FlagDescCTNoCommit)
	fs.BoolFunc(constants.FlagCTIncludeMerges, constants.FlagDescCTIncludeMerges,
		func(string) error { opts.IncludeMerges = true; return nil })
	fs.BoolFunc(constants.FlagCTNoIncludeMerges, constants.FlagDescCTNoIncludeMerges,
		func(string) error { opts.IncludeMerges = false; return nil })
	fs.BoolVar(&opts.Mirror, constants.FlagCTMirror, false, constants.FlagDescCTMirror)
	fs.BoolVar(&opts.ForceReplay, constants.FlagCTForceReplay, false, constants.FlagDescCTForceReplay)
	fs.BoolVar(&opts.Interleave, constants.FlagCTInterleave, false, constants.FlagDescCTInterleave)
	fs.BoolVar(&opts.IsLocal, "local", false, "Provision destination repository locally only (skip GitHub creation)")
	fs.BoolVar(&opts.IsLocal, "no-remote", false, "Alias for --local")
	registerMessagePolicyToggles(fs, opts)
}

// registerMessagePolicyToggles wires the on/off pairs for §6 stages.
// Uses fs.BoolFunc (Go 1.21+) so negations don't consume a value.
// Order on the command line is the order of effect (last wins).
func registerMessagePolicyToggles(fs *flag.FlagSet, opts *committransfer.Options) {
	fs.BoolFunc(constants.FlagCTConventional, constants.FlagDescCTConventional,
		func(string) error { opts.Message.Conventional = true; return nil })
	fs.BoolFunc(constants.FlagCTNoConventional, constants.FlagDescCTNoConventional,
		func(string) error { opts.Message.Conventional = false; return nil })
	fs.BoolFunc(constants.FlagCTProvenance, constants.FlagDescCTProvenance,
		func(string) error { opts.Message.Provenance = true; return nil })
	fs.BoolVar(&opts.Message.TemplateOverride, constants.FlagCTTemplateOverride, false, constants.FlagDescCTTemplateOverride)
	fs.BoolFunc(constants.FlagCTNoProvenance, constants.FlagDescCTNoProvenance,
		func(string) error { opts.Message.Provenance = false; return nil })
}

// registerCommitTransferStrings wires value-taking flags + repeatable
// regex patterns. --no-strip and --no-drop are BoolFunc (no value).
func registerCommitTransferStrings(fs *flag.FlagSet, opts *committransfer.Options) {
	fs.StringVar(&opts.PRMode, constants.FlagCTPR, opts.PRMode, constants.FlagDescCTPRMode)
	fs.IntVar(&opts.Limit, constants.FlagCTLimit, 0, constants.FlagDescCTLimit)
	fs.StringVar(&opts.Since, constants.FlagCTSince, "", constants.FlagDescCTSince)
	fs.IntVar(&opts.MaxHistoryScan, constants.FlagCTMaxHistoryScan, 0, constants.FlagDescCTMaxHistoryScan)
	fs.Func(constants.FlagCTStrip, constants.FlagDescCTStrip, func(v string) error {
		opts.Message.StripPatterns = append(opts.Message.StripPatterns, v)

		return nil
	})
	fs.Func(constants.FlagCTAppendFooter, constants.FlagDescCTAppendFooter, func(v string) error {
		opts.Message.AppendFooters = append(opts.Message.AppendFooters, v)

		return nil
	})
	fs.Func(constants.FlagCTDrop, constants.FlagDescCTDrop, func(v string) error {
		opts.Message.DropPatterns = append(opts.Message.DropPatterns, v)

		return nil
	})
	fs.BoolFunc(constants.FlagCTNoStrip, constants.FlagDescCTNoStrip, func(string) error {
		opts.Message.StripPatterns = nil

		return nil
	})
	fs.BoolFunc(constants.FlagCTNoDrop, constants.FlagDescCTNoDrop, func(string) error {
		opts.Message.DropPatterns = nil

		return nil
	})
}

// commitTransferSpecFor maps a command name or alias to its spec.
func commitTransferSpecFor(command string) (commitTransferSpec, bool) {
	switch command {
	case constants.CmdCommitLeft, constants.CmdCommitLeftA:
		return commitTransferSpec{
			Name: constants.CmdCommitLeft, LogPrefix: constants.LogPrefixCommitLeft,
		}, true
	case constants.CmdCommitRight, constants.CmdCommitRightA:
		return commitTransferSpec{
			Name: constants.CmdCommitRight, LogPrefix: constants.LogPrefixCommitRight,
		}, true
	case constants.CmdCommitBoth, constants.CmdCommitBothA:
		return commitTransferSpec{
			Name: constants.CmdCommitBoth, LogPrefix: constants.LogPrefixCommitBoth,
		}, true
	case constants.CmdPR, constants.CmdPullRequest:
		return commitTransferSpec{
			Name: constants.CmdPR, LogPrefix: constants.LogPrefixPR,
		}, true
	}

	return commitTransferSpec{}, false
}

func runPRClean(args []string) error {
	checkHelp(constants.CmdPRClean, args)
	isYes, repoPath := parsePRCleanArgs(args)
	repoRoot, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "runPRClean.getwd")
	}
	if repoPath != "" {
		repoRoot = repoPath
	}

	return executePRClean(repoRoot, isYes)
}

func parsePRCleanArgs(args []string) (bool, string) {
	isYes := false
	repoPath := ""
	for _, arg := range args {
		if arg == "-y" || arg == "--yes" {
			isYes = true
			continue
		}
		if !strings.HasPrefix(arg, "-") && repoPath == "" {
			repoPath = arg
		}
	}

	return isYes, repoPath
}

func executePRClean(repoRoot string, isYes bool) error {
	slug := prdb.SanitizeRepoSlug(filepath.Base(repoRoot))
	dbRes := prdb.OpenPrSplitDb(slug, repoRoot)
	if dbRes.IsFailure() {
		return dbRes.Err
	}
	db := dbRes.Value
	defer db.Close()

	branchesRes := db.ListMergedPrBranches()
	if branchesRes.IsFailure() {
		return branchesRes.Err
	}

	return processPRCleanBranches(db, repoRoot, branchesRes.Value, isYes)
}

func processPRCleanBranches(db *prdb.PrSplitDb, repoRoot string, branches []prdb.PrBranchRecord, isYes bool) error {
	if len(branches) == 0 {
		fmt.Println("No merged PR branches to clean.")
		return nil
	}
	if !isYes && !confirmPRClean(len(branches)) {
		fmt.Println("PR clean canceled.")
		return nil
	}

	return pruneMergedBranches(db, repoRoot, branches)
}

func confirmPRClean(count int) bool {
	fmt.Printf("Remove %d closed/merged PR branches? [y/N]: ", count)
	var response string
	fmt.Scanln(&response)
	resp := strings.TrimSpace(strings.ToLower(response))

	return resp == "y" || resp == "yes"
}

func pruneMergedBranches(db *prdb.PrSplitDb, repoRoot string, branches []prdb.PrBranchRecord) error {
	for _, b := range branches {
		cmd := exec.Command("git", "-C", repoRoot, "branch", "-D", b.BranchName)
		_ = cmd.Run()
		db.MarkPrBranchDeleted(b.BranchName)
		fmt.Printf("Deleted PR branch: %s\n", b.BranchName)
	}
	fmt.Printf("Cleaned %d merged PR branches.\n", len(branches))

	return nil
}

func runPRList(args []string) error {
	checkHelp(constants.CmdPRList, args)
	isJSON, repoPath := parsePRListArgs(args)
	repoRoot, err := os.Getwd()
	if err != nil {
		return apperror.WrapSimple(err, "runPRList.getwd")
	}
	if repoPath != "" {
		repoRoot = repoPath
	}

	return executePRList(repoRoot, isJSON)
}

func parsePRListArgs(args []string) (bool, string) {
	isJSON := false
	repoPath := ""
	for _, arg := range args {
		if arg == "--json" {
			isJSON = true
			continue
		}
		if !strings.HasPrefix(arg, "-") && repoPath == "" {
			repoPath = arg
		}
	}

	return isJSON, repoPath
}

func executePRList(repoRoot string, isJSON bool) error {
	slug := prdb.SanitizeRepoSlug(filepath.Base(repoRoot))
	dbRes := prdb.OpenPrSplitDb(slug, repoRoot)
	if dbRes.IsFailure() {
		return dbRes.Err
	}
	db := dbRes.Value
	defer db.Close()

	activeRes := db.ListActivePrBranches()
	mergedRes := db.ListMergedPrBranches()
	if activeRes.IsFailure() {
		return activeRes.Err
	}
	if mergedRes.IsFailure() {
		return mergedRes.Err
	}

	return displayPRList(activeRes.Value, mergedRes.Value, isJSON)
}

func displayPRList(active, merged []prdb.PrBranchRecord, isJSON bool) error {
	if isJSON {
		return outputPRListJSON(active, merged)
	}
	renderPRListTable(active, merged)

	return nil
}

func outputPRListJSON(active, merged []prdb.PrBranchRecord) error {
	payload := map[string]any{
		"active": active,
		"merged": merged,
	}
	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "outputPRListJSON.marshal")
	}
	fmt.Println(string(bytes))

	return nil
}

func renderPRListTable(active, merged []prdb.PrBranchRecord) {
	fmt.Printf("PR Branches (Active: %d, Merged: %d):\n", len(active), len(merged))
	for _, b := range active {
		fmt.Printf("  [active] %s (%s)\n", b.BranchName, b.BranchType)
	}
	for _, b := range merged {
		fmt.Printf("  [merged] %s (%s)\n", b.BranchName, b.BranchType)
	}
}
