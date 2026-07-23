package cmd

// clonetermstream.go — shared "print one RepoTermBlock right before
// a clone runs" helper used by clone, clone-now, clone-pick, and
// clone-from when the user passes `--output terminal`. Keeps every
// clone-related command byte-identical in its per-repo preview.
//
// Stream order (matches the answer locked in chat):
//
//	[block for repo i]   ← printed BEFORE we shell out to git
//	<git clone progress>
//	[block for repo i+1]
//	<git clone progress>
//	...
//
// The block goes to stdout (matches clone-next and clone-from).
// Clone progress stays on stderr per the project's stream contract.

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/render"
)

// CloneTermBlockInput carries the per-repo data every clone command
// already has on hand. Branch/BranchSource may be empty — the
// renderer falls back to "(unknown)" so the block shape is stable.
//
// Faithfulness contract (audited): the printed `cmd:` line MUST be
// byte-identical to the argv the executor passes to exec.Command.
// Each caller controls three override fields to achieve that:
//
//   - CmdBranch:        branch passed to `-b` in the printed cmd.
//     Empty means "no `-b` flag".
//   - CmdExtraArgsPre:  literal tokens inserted between `git clone`
//     and the `-b` slot. Used by clone-pick for
//     `--filter=blob:none --no-checkout` and the
//     long-form `--branch X` / `--depth N`.
//   - CmdExtraArgsPost: literal tokens inserted between the `-b`
//     slot and the positional `<url> <dest>` pair.
//     Used by clone-from for `--depth=N` (its
//     executor places --depth AFTER -b).
type CloneTermBlockInput struct {
	Index        int
	Name         string
	Branch       string
	BranchSource string
	OriginalURL  string
	TargetURL    string
	Dest         string
	// CmdBranch overrides which branch (if any) is rendered as `-b`
	// in the printed cmd. Empty = no `-b` flag, regardless of what
	// Branch (the display field) holds. Defaults to Branch when the
	// caller leaves both CmdBranch AND CmdExtraArgs* unset (legacy
	// fallback for clone-now / clone-pick rows).
	CmdBranch string
	// CmdExtraArgsPre are tokens between `git clone` and `-b`.
	CmdExtraArgsPre []string
	// CmdExtraArgsPost are tokens between `-b <branch>` and the
	// positional URL/dest pair.
	CmdExtraArgsPost []string
}

// maybePrintCloneTermBlock emits the standardized RepoTermBlock to
// stdout when output == "terminal". Any other value (including the
// empty default) is a no-op so existing CI logs and screenshots
// stay byte-identical for callers that didn't opt in.
//
// Returns nothing — write errors on stdout are surfaced by the
// caller's later Println/Printf calls. Per the zero-swallow policy
// we still log to stderr so a closed stdout doesn't silently drop
// the preview.
func maybePrintCloneTermBlock(output string, in CloneTermBlockInput) {
	if output != constants.OutputTerminal {
		return
	}
	block := render.RepoTermBlock{
		Index:        in.Index,
		Name:         in.Name,
		Branch:       in.Branch,
		BranchSource: in.BranchSource,
		OriginalURL:  in.OriginalURL,
		TargetURL:    in.TargetURL,
		CloneCommand: buildCloneCommand(in),
	}
	if err := render.RenderRepoTermBlock(os.Stdout, block); err != nil {
		fmt.Fprintf(os.Stderr,
			"  Warning: could not write terminal block for %s: %v\n",
			in.Name, err)
	}
}

// buildCloneCommand returns the exact `git clone` command the
// executor will run. The output is byte-identical to the argv
// joined with spaces — caller-supplied CmdExtraArgsPre/Post and
// CmdBranch drive the differences between commands (see file-header
// faithfulness contract).
//
// Branch resolution rule: an explicitly-set CmdBranch wins. If
// CmdBranch is empty AND neither CmdExtraArgsPre nor CmdExtraArgsPost
// was passed (both nil), we fall back to in.Branch as the `-b` value
// — this keeps the original clone-now / clone-pick row callers
// working without per-call-site churn. Callers that explicitly want
// NO `-b` (URL-driven clone) set CmdBranch=="" AND pass a non-nil
// CmdExtraArgsPre (even if empty slice) — see pickCmdBranch.
func buildCloneCommand(in CloneTermBlockInput) string {
	parts := []string{constants.GitBin, constants.GitClone}
	parts = append(parts, in.CmdExtraArgsPre...)
	branch := pickCmdBranch(in)
	if len(branch) > 0 {
		parts = append(parts, constants.GitBranchFlag, branch)
	}
	parts = append(parts, in.CmdExtraArgsPost...)
	parts = append(parts, in.TargetURL, in.Dest)

	return strings.Join(parts, " ")
}

// pickCmdBranch resolves which branch (if any) to render as `-b`.
// See buildCloneCommand for the rule. Split out so the rule has one
// home and a future change (e.g. always-explicit) is one edit.
func pickCmdBranch(in CloneTermBlockInput) string {
	if len(in.CmdBranch) > 0 {
		return in.CmdBranch
	}
	// Non-nil CmdExtraArgsPre is the explicit "no -b" sentinel
	// (URL-driven clone). nil = legacy fallback to in.Branch.
	if in.CmdExtraArgsPre != nil {
		return ""
	}

	return in.Branch
}
