// Package cmd — clonefixrepo.go: entry points for `gitmap clone-fix-repo`
// (alias `cfr`) and `gitmap clone-fix-repo-pub` (alias `cfrp`).
//
// These are convenience pipelines that chain three existing commands
// in one shot:
//
//	cfr  : clone <url>  →  cd <folder>  →  fix-repo --all
//	cfrp : clone <url>  →  cd <folder>  →  fix-repo --all  →  make-public --yes
//
// Implementation strategy: the chained commands (runFixRepo,
// runMakePublic) all call os.Exit at the end, which would terminate
// our parent process before the next step runs. To stay decoupled
// and side-effect-clean, we shell out to our own binary (resolved
// via os.Executable) for the fix-repo and make-public steps after
// invoking executeDirectClone in-process. This also keeps each
// step's exit code, stdout, and stderr semantics intact.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// runCloneFixRepo implements `gitmap clone-fix-repo` (alias cfr).
func runCloneFixRepo(args []string) error {
	checkHelp(constants.CmdCloneFixRepo, args)
	runCloneFixRepoPipeline(args, false)
	return nil
}

// runCloneFixRepoPub implements `gitmap clone-fix-repo-pub` (alias cfrp).
func runCloneFixRepoPub(args []string) error {
	checkHelp(constants.CmdCloneFixRepoPub, args)
	runCloneFixRepoPipeline(args, true)
	return nil
}

// runCloneFixRepoPipeline is the shared core. `makePublic` controls
// whether the optional 3rd step (visibility flip) runs.
func runCloneFixRepoPipeline(args []string, makePublic bool) error {
	parallel, args := extractParallelFlag(args)
	modifiers, args := ParseCfrModifiers(args)
	if modifiers.PromotePublic {
		makePublic = true
	}
	url, folder, noVSCodeSync, reqVer, useSSH, useHTTPS, autoYes, dryRun, noCommit, noPush := parseCloneFixRepoArgs(args)
	modifiers.NoCommit = modifiers.NoCommit || noCommit
	modifiers.NoPush = modifiers.NoPush || noPush
	f := cloneFixRepoFlags{url, folder, noVSCodeSync, reqVer, useSSH, useHTTPS, autoYes, dryRun, noCommit, noPush}
	if dispatchCFRMultiURL(f, makePublic, modifiers, parallel) {
		return nil
	}
	runSingleCloneFixRepo(f, makePublic, modifiers)
	return nil
}

func dispatchCFRMultiURL(f cloneFixRepoFlags, makePublic bool, modifiers CfrModifierFlags, parallel int) bool {
	urls := splitCommaURLs(f.url)
	if len(urls) <= 1 {
		return false
	}
	runParallelCloneFixRepo(urls, makePublic, f.noVSCodeSync, f.requireVersion, f.useSSH, f.useHTTPS, f.autoYes, f.dryRun, modifiers, parallel)
	return true
}

func runSingleCloneFixRepo(f cloneFixRepoFlags, makePublic bool, modifiers CfrModifierFlags) error {
	folderName, absPath := validateAndPrepareCFR(&f)
	executeCFRClone(f.url, folderName, absPath, f)
	if f.dryRun {
		printDryRunMessage(makePublic, absPath)
		return nil
	}
	executeCFRPostSteps(absPath, makePublic, f, modifiers)
	return nil
}

func validateAndPrepareCFR(f *cloneFixRepoFlags) (string, string) {
	SetCloneDryRun(f.dryRun)
	SetCloneAssumeYes(f.autoYes)
	applyCloneAssumeYesEnv(f.autoYes)
	if len(f.url) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrCloneFixRepoUsage)
		os.Exit(constants.ExitCloneFixRepoBadFlag)
	}
	f.url = applyCloneFixRepoScheme(f.url, f.useSSH, f.useHTTPS)
	escapeNestedGitRepo()
	folderName := deriveFolderNameForCFR(f.url, f.folder)
	absPath := resolveCloneTargetFolder(f.url, folderName)
	return folderName, absPath
}

func executeCFRClone(url, folderName, absPath string, f cloneFixRepoFlags) {
	url = preferExistingFolderTransport(url, absPath)
	url = coerceURLToStoredTransport(url)
	requireOnline()
	executeDirectClone(url, folderName, true, false, "", f.noVSCodeSync)
	if !f.dryRun {
		persistRecloneTransport(url)
	}
}

func executeCFRPostSteps(absPath string, makePublic bool, f cloneFixRepoFlags, modifiers CfrModifierFlags) {
	if err := os.Chdir(absPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoChdirFmt, absPath, err)
		os.Exit(constants.ExitCloneFixRepoChdir)
	}
	maybeRunFixRepoStep(absPath, f.requireVersion)
	if makePublic {
		runChainedGitmapStep([]string{constants.CmdMakePublic, "--" + constants.FlagVisYes})
		runCFRPPriorVersionPrivatize(absPath, f.autoYes)
	}
	dispatchCodingGuidelinesModifier(absPath, modifiers)
	fmt.Printf(constants.MsgCloneFixRepoDone, absPath)
}

// runParallelCloneFixRepo encapsulates the parallel clone execution to avoid nested ifs.
func runParallelCloneFixRepo(urls []string, makePublic bool, noVSCodeSync bool, requireVersion bool, useSSH bool, useHTTPS bool, autoYes bool, dryRun bool, modifiers CfrModifierFlags, parallel int) error {
	subcmd := constants.CmdCloneFixRepo
	if makePublic {
		subcmd = constants.CmdCloneFixRepoPub
	}
	passthrough := buildCFRPassthroughFlags(noVSCodeSync, requireVersion, useSSH, useHTTPS, autoYes, dryRun, modifiers.NoCommit, modifiers.NoPush)
	leadingMods := buildCFRLeadingModifiers(modifiers)
	failed := runCloneFixRepoParallel(urls, subcmd, leadingMods, passthrough, parallel)
	if failed > 0 {
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
	return nil
}

func printDryRunMessage(makePublic bool, absPath string) {
	suffix := ""
	if makePublic {
		suffix = " → make-public --yes"
	}
	fmt.Printf("  "+constants.MsgCloneDryRunNoop+"\n  would chain: fix-repo --all%s @ %s\n",
		suffix, absPath)
}

// buildCFRLeadingModifiers renders modifier flags back into their
// positional-token form so parallel workers re-parse them via
// ParseCfrModifiers. `p` is intentionally omitted here: the subcmd
// (cfr vs cfrp) already encodes the public-visibility choice for
// workers, so re-passing `p` would be redundant.
func buildCFRLeadingModifiers(m CfrModifierFlags) []string {
	out := make([]string, 0, 1)
	if m.InstallCodingGuidelines {
		out = append(out, constants.CfrModifierCodingGuidelines)
	}
	return out
}

// dispatchCodingGuidelinesModifier invokes the v24 Coding Guidelines
// installer against the freshly cloned working tree when the `cg`
// modifier is present, then auto-commits (and optionally pushes) any
// files the installer produced. Errors are already logged by the
// underlying helpers (zero-swallow policy); we surface a non-zero
// exit so the pipeline halts.
func dispatchCodingGuidelinesModifier(absPath string, m CfrModifierFlags) {
	if !m.InstallCodingGuidelines {
		return
	}
	if err := RunCodingGuidelinesInstall(CodingGuidelinesOpts{WorkingDir: absPath}); err != nil {
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
	commitOpts := CGCommitOpts{WorkingDir: absPath, NoCommit: m.NoCommit, NoPush: m.NoPush}
	if err := CommitCodingGuidelines(commitOpts); err != nil {
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
}

// applyCloneFixRepoScheme honors --ssh / --https (and short aliases
// --sh / --ht) by rewriting the URL before the in-process clone runs.
//
// Mirrors `gitmap clone --ssh` semantics: when both flags are set,
// --ssh wins and a one-line stderr warning is printed. Unrecognized
// URL shapes are returned unchanged so non-URL positionals still flow
// through.
func applyCloneFixRepoScheme(url string, useSSH, useHTTPS bool) string {
	if useSSH && useHTTPS {
		fmt.Fprintln(os.Stderr, "warning: --ssh and --https both set; --ssh wins")
		useHTTPS = false
	}
	if converted, ok := applySSHScheme(url, useSSH); ok {
		return converted
	}
	if converted, ok := applyHTTPSScheme(url, useHTTPS); ok {
		return converted
	}
	return url
}

func applySSHScheme(url string, useSSH bool) (string, bool) {
	if converted, ok := ConvertURLToSSH(url); useSSH && ok {
		if converted != url {
			fmt.Printf("↪ --ssh rewrite: %s → %s\n", url, converted)
		}
		return converted, true
	}
	return url, false
}

func applyHTTPSScheme(url string, useHTTPS bool) (string, bool) {
	if converted, ok := ConvertURLToHTTPS(url); useHTTPS && ok {
		if converted != url {
			fmt.Printf("↪ --https rewrite: %s → %s\n", url, converted)
		}
		return converted, true
	}
	return url, false
}

// maybeRunFixRepoStep runs `fix-repo --all` only when the cloned repo
// identity carries a `-vN` suffix. The identity comes from Git remote
// metadata first, not the flattened local folder name.
// `--require-version` restores the strict (exit-4) failure mode for
// CI pipelines that want the old contract.
func maybeRunFixRepoStep(absPath string, requireVersion bool) {
	repoName := resolveCloneFixRepoName(absPath)
	parsed := clonenext.ParseRepoName(repoName)
	if parsed.HasVersion {
		runChainedGitmapStep([]string{constants.CmdFixRepo, "--" + constants.FixRepoFlagAll})

		return
	}
	if requireVersion {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoNeedVersion, parsed.BaseName)
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
	fmt.Printf(constants.MsgCloneFixRepoSkipNoVer, parsed.BaseName)
}

func resolveCloneFixRepoName(absPath string) string {
	remoteURL, err := gitutil.RemoteURL(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnCloneFixRepoRemoteFmt, absPath, err)

		return filepath.Base(absPath)
	}
	repo := repoNameFromURL(remoteURL)
	if len(repo) > 0 {
		return repo
	}
	fmt.Fprintf(os.Stderr, constants.WarnCloneFixRepoRemoteFmt, remoteURL, constants.ErrCloneFixRepoRemoteParse)

	return filepath.Base(absPath)
}

type cloneFixRepoFlags struct {
	url            string
	folder         string
	noVSCodeSync   bool
	requireVersion bool
	useSSH         bool
	useHTTPS       bool
	autoYes        bool
	dryRun         bool
	noCommit       bool
	noPush         bool
}

func applyCFRFlag(name string, f *cloneFixRepoFlags) bool {
	switch name {
	case constants.FlagNoVSCodeSync:
		f.noVSCodeSync = true
	case constants.FlagRequireVersion:
		f.requireVersion = true
	case "ssh", "sh":
		f.useSSH = true
	case "https", "ht":
		f.useHTTPS = true
	default:
		return applyCFRFlagExtra(name, f)
	}
	return true
}

func applyCFRFlagExtra(name string, f *cloneFixRepoFlags) bool {
	switch name {
	case "y", "yes":
		f.autoYes = true
	case constants.FlagCloneDryRun, constants.FlagCloneDryRunShort:
		f.dryRun = true
	case constants.FlagCGNoCommit:
		f.noCommit = true
	case constants.FlagCGNoPush:
		f.noPush = true
	default:
		return false
	}
	return true
}

func extractCFRPositionals(args []string, f *cloneFixRepoFlags) []string {
	var positional []string
	for _, a := range args {
		if !applyCFRFlag(strings.TrimLeft(a, "-"), f) && len(a) > 0 && a[0] != '-' {
			positional = append(positional, a)
		}
	}
	return positional
}

func assignCFRPositionals(positional []string, f *cloneFixRepoFlags) {
	if len(positional) > 0 {
		f.url = positional[0]
	}
	if len(positional) > 1 {
		f.folder = positional[1]
	}
}

// parseCloneFixRepoArgs returns (url, folderName, noVSCodeSync,
// requireVersion, useSSH, useHTTPS, autoYes, dryRun, noCommit,
// noPush). First non-flag arg is the URL; second non-flag is the
// destination folder. Recognized flags:
// --no-vscode-sync, --require-version, --ssh/-ssh/--sh,
// --https/-https/--ht, --no-commit, --no-push. Single-dash forms are
// accepted to match Go's stdlib `flag` package behavior the user
// expects from `-ssh`.
func parseCloneFixRepoArgs(args []string) (string, string, bool, bool, bool, bool, bool, bool, bool, bool) {
	var f cloneFixRepoFlags
	positional := extractCFRPositionals(args, &f)
	assignCFRPositionals(positional, &f)
	return f.url, f.folder, f.noVSCodeSync, f.requireVersion, f.useSSH, f.useHTTPS, f.autoYes, f.dryRun, f.noCommit, f.noPush
}

// resolveCloneTargetFolder mirrors the folder-naming logic in
// executeDirectClone so we know which directory to cd into after
// the clone step finishes. Versioned URLs auto-flatten to BaseName.
func resolveCloneTargetFolder(url, folderName string) string {
	folderName = deriveFolderNameForCFR(url, folderName)
	abs, err := filepath.Abs(folderName)
	if err != nil {
		return folderName
	}

	return abs
}

func deriveFolderNameForCFR(url string, folderName string) string {
	if len(folderName) > 0 {
		return folderName
	}
	repoName := repoNameFromURL(url)
	parsed := clonenext.ParseRepoName(repoName)
	if parsed.HasVersion {
		return parsed.BaseName
	}
	return repoName
}

// runChainedGitmapStep re-execs the current gitmap binary with the
// given args, streaming stdin/stdout/stderr through. Any non-zero
// exit propagates immediately so the pipeline halts on first failure.
func runChainedGitmapStep(args []string) error {
	bin, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoExecFmt, err)
		os.Exit(constants.ExitCloneFixRepoChainFailed)
	}
	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	handleChainedStepResult(cmd.Run())
	return nil
}

func handleChainedStepResult(runErr error) {
	if runErr == nil {
		return
	}
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}
	fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoExecFmt, runErr)
	os.Exit(constants.ExitCloneFixRepoChainFailed)
}
