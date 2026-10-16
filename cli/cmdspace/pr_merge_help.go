package cmdspace

import "fmt"

// printPRMergeHelp prints usage for `gitmap space PR-merge --help`.
func printPRMergeHelp() {
	fmt.Print(`Usage: gitmap space PR-merge [target] [flags]

Merges the current branch into the target branch. With no target, the
default branch is resolved automatically (origin/HEAD, then main, then
master).

Behavior:
  - The working tree must be clean; a dirty tree is always refused.
  - Fetches origin, checks out the target, fast-forwards it, then merges.
  - Default is a regular merge commit: all commits are preserved.
  - On conflicts the merge is aborted and the conflicted files are listed.
    Resolve them manually and merge again - nothing is auto-resolved.
  - Pushes the target to origin unless --no-push is passed.
  - History is never rewritten: --force is rejected.

Flags:
  --squash         Squash-merge. The commit message embeds the PR title
                   and markdown body when a PR exists for the branch,
                   plus a heatmap of the top changed files.
  --single-commit  Flat single-commit merge with a one-line message.
  --flat           Alias of --single-commit.
  --no-push        Merge locally only; do not push the target.

Examples:
  gitmap space PR-merge main
  gitmap space PR-merge
  gitmap space PR-merge main --squash
  gitmap space PR-merge main --no-push
`)
}
