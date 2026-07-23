// Package cmd — amend.go handles flag parsing and orchestration for the amend command.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// amendFlags holds parsed flags for the amend command.
type amendFlags struct {
	commitHash string
	name       string
	email      string
	branch     string
	dryRun     bool
	forcePush  bool
}

// runAmend is the entry point for the amend command.
func runAmend(args []string) {
	checkHelp("amend", args)
	flags := parseAmendFlags(args)
	validateAmendFlags(flags)
	executeAmend(flags)
}

// parseAmendFlags parses command-line flags for amend.
func parseAmendFlags(args []string) amendFlags {
	var f amendFlags

	// Extract positional SHA (first arg that doesn't start with --)
	remaining := extractAmendSHA(args, &f)

	fs := flag.NewFlagSet(constants.CmdAmend, flag.ExitOnError)
	fs.StringVar(&f.name, constants.FlagAmendName, "", constants.FlagDescAmendName)
	fs.StringVar(&f.email, constants.FlagAmendEmail, "", constants.FlagDescAmendEmail)
	fs.StringVar(&f.branch, constants.FlagAmendBranch, "", constants.FlagDescAmendBranch)
	fs.BoolVar(&f.dryRun, constants.FlagAmendDryRun, false, constants.FlagDescAmendDryRun)
	fs.BoolVar(&f.forcePush, constants.FlagAmendForcePush, false, constants.FlagDescAmendForcePush)

	_ = fs.Parse(remaining)

	return f
}

// extractAmendSHA pulls the optional SHA from the first positional arg.
func extractAmendSHA(args []string, f *amendFlags) []string {
	if len(args) == 0 {
		return args
	}

	if args[0] == "" || args[0][0] == '-' {
		return args
	}

	f.commitHash = args[0]

	return args[1:]
}

// validateAmendFlags ensures at least one of name/email is provided.
func validateAmendFlags(f amendFlags) {
	if f.name == "" && f.email == "" {
		fmt.Fprint(os.Stderr, constants.ErrAmendNoFlags)
		os.Exit(1)
	}
}

// executeAmend runs the main amend workflow.
func executeAmend(f amendFlags) {
	originalBranch := getCurrentBranch()

	if f.branch != "" {
		switchBranch(f.branch)
	}

	targetBranch := resolveTargetBranch(f)
	mode := resolveAmendMode(f)
	commits := listCommitsForAmend(f)

	if len(commits) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrAmendNoCommits)
		os.Exit(1)
	}

	prevName, prevEmail := detectPreviousAuthor(commits)

	if f.dryRun {
		printAmendDryRun(commits, f, prevName, prevEmail)

		returnToBranch(f, originalBranch)

		return
	}

	fmt.Print(constants.MsgAmendWarnRewrite)
	printAmendHeader(f, commits, targetBranch, prevName, prevEmail)

	runFilterBranch(f, commits)
	printAmendProgress(commits)

	auditPath := writeAmendAudit(f, commits, targetBranch, mode, prevName, prevEmail)
	saveAmendToDB(f, commits, targetBranch, mode, prevName, prevEmail)

	fmt.Printf(constants.MsgAmendDone, len(commits))
	fmt.Printf(constants.MsgAmendAuditFile, auditPath)
	fmt.Print(constants.MsgAmendAuditDB)

	if f.forcePush {
		runForcePush()
	} else {
		fmt.Print(constants.MsgAmendWarnPush)
	}

	returnToBranch(f, originalBranch)
}

// resolveTargetBranch returns the branch being amended.
func resolveTargetBranch(f amendFlags) string {
	if f.branch != "" {
		return f.branch
	}

	return getCurrentBranch()
}

// resolveAmendMode determines the amend mode from flags.
func resolveAmendMode(f amendFlags) string {
	if f.commitHash == "" {
		return constants.AmendModeAll
	}

	if f.commitHash == "HEAD" {
		return constants.AmendModeHead
	}

	return constants.AmendModeRange
}

// returnToBranch switches back to the original branch if needed.
func returnToBranch(f amendFlags, original string) {
	if f.branch == "" {
		return
	}

	current := getCurrentBranch()
	if current == original {
		return
	}

	fmt.Printf(constants.MsgAmendReturn, original)
	switchBranch(original)
}
