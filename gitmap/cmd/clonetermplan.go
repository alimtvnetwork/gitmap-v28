package cmd

// clonetermplan.go — adapters that fan a finite Plan into one
// RepoTermBlock per row using the shared streaming helper.
//
// Used by:
//   - clone-now's DRY-RUN path (printCloneNowTermBlocks): all blocks
//     upfront, since there is no execution to interleave with.
//   - clone-pick (printClonePickTermBlock): always one row, so
//     "upfront" and "streamed" are the same thing.
//
// The EXECUTE paths for clone-now and clone-from no longer call
// these helpers — they stream one block per row via the executor's
// BeforeRow hook (see clonetermrow.go + execute_hooks.go in each
// package). That puts each preview immediately before its own
// `git clone` instead of dumping every block upfront.
//
// Per the locked design:
//
//   - URL-driven commands (clone) interleave one block per URL
//     immediately before that URL's clone runs (see clonetermurl.go).
//   - Plan-driven EXECUTE (clone-now, clone-from) streams via the
//     BeforeRow hook (see clonetermrow.go).
//   - Plan-driven DRY-RUN (clone-now) prints all blocks upfront
//     because there's no clone progress to interleave with.

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonepick"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// printCloneNowTermBlocks emits one block per plan row to stdout
// using ls-remote-detected branches. Always runs in `terminal`
// mode (callers gate by cfg.output == OutputTerminal before
// calling) so the function is a straight loop with no flag check.
func printCloneNowTermBlocks(plan clonenow.Plan) {
	for i, row := range plan.Rows {
		url := row.PickURL(plan.Mode)
		dest := row.RelativePath
		branch := row.Branch
		source := "manifest"
		if len(branch) == 0 {
			branch = detectRemoteHEAD(url)
			source = remoteBranchSource(branch)
		}
		maybePrintCloneTermBlock(constants.OutputTerminal, CloneTermBlockInput{
			Index:        i + 1,
			Name:         pickCloneNowName(row, dest),
			Branch:       branch,
			BranchSource: source,
			OriginalURL:  url,
			TargetURL:    url,
			Dest:         dest,
			CmdBranch:    row.Branch, // executor only passes -b when row.Branch is set
		})
	}
}

// pickCloneNowName picks the most informative repo name for the
// block: explicit RepoName when present, else the dest folder.
func pickCloneNowName(row clonenow.Row, dest string) string {
	if len(row.RepoName) > 0 {
		return row.RepoName
	}

	return dest
}

// printClonePickTermBlock emits a single block for the clone-pick
// plan. clone-pick is always one URL → one block; the destination
// is plan.DestDir (defaults to "."). Branch is taken from the plan
// when the user pinned --branch, else discovered via ls-remote.
//
// Faithfulness: clone-pick's executor (clonepick/sparse.go
// gitClonePartial) shells out to:
//
//	git clone --filter=blob:none --no-checkout \
//	    [--branch X] [--depth N] <url> <dest>
//
// Note the long-form `--branch` (NOT `-b`) and space-separated
// `--depth N` (NOT `--depth=N`). To stay byte-identical we bypass
// the standard `-b` slot (CmdBranch="") and emit ALL pre-positional
// flags via CmdExtraArgsPre, exactly mirroring buildCloneCommandPreview
// in clonepick/render.go (which is what the dry-run already prints).
func printClonePickTermBlock(plan clonepick.Plan) {
	branch := plan.Branch
	source := "manifest"
	if len(branch) == 0 {
		branch = detectRemoteHEAD(plan.RepoUrl)
		source = remoteBranchSource(branch)
	}
	name := plan.Name
	if len(name) == 0 {
		name = repoNameFromURL(plan.RepoUrl)
	}
	in := CloneTermBlockInput{
		Index:           1,
		Name:            name,
		Branch:          branch,
		BranchSource:    source,
		OriginalURL:     plan.RepoUrl,
		TargetURL:       plan.RepoUrl,
		Dest:            plan.DestDir,
		CmdBranch:       "", // explicit opt-out: clone-pick uses --branch (long form)
		CmdExtraArgsPre: clonePickCmdPre(plan),
	}
	maybePrintCloneTermBlock(constants.OutputTerminal, in)
	// Verifier: clonepick.BuildGitArgs is the same builder
	// gitClonePartial uses, so any drift between displayed and
	// executed argv surfaces in the report.
	runCmdFaithfulCheck(in, clonepick.BuildGitArgs(plan, plan.DestDir))
	runCmdPrintArgv(clonepick.BuildGitArgs(plan, plan.DestDir))
}

// clonePickCmdPre assembles the literal pre-positional flag tokens
// for clone-pick's printed cmd, in the exact order
// gitClonePartial appends them. plan.Branch (NOT the detected
// fallback) drives --branch — same rule as the executor.
func clonePickCmdPre(plan clonepick.Plan) []string {
	parts := []string{"--filter=blob:none", "--no-checkout"}
	if len(plan.Branch) > 0 {
		parts = append(parts, "--branch", plan.Branch)
	}
	if plan.Depth > 0 {
		parts = append(parts, "--depth", fmt.Sprintf("%d", plan.Depth))
	}

	return parts
}
