package cmdspace

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"os"
)

// spaceUsage is printed for `gitmap space --help` or bare `gitmap space`.
const spaceUsage = `Usage: gitmap space <subcommand> [flags]

Space operations namespace. Subcommands:

  common      Apply the curated common baselines (.gitignore,
              .gitattributes, .prettierignore, .prettierrc) and run
              'git lfs install --local' — same logic as 'gitmap commons'.
  backup-branch
              Create a backup branch (backup/<slug>) from the current
              HEAD for the given task string, then push it to origin.
  PR-merge      Merge the current branch into the target branch
              (default branch when omitted). Regular merge commit;
              --squash embeds PR details plus a changed-files heatmap.
  send-email    Send an email. Delegates to the go-email-reader binary
              when installed; otherwise uses the built-in SMTP sender.
  schedule-send-email
              Register an email trigger fired when a task (release,
              pull, clone, scan, push, backup) completes.
  email ui    Interactive terminal UI for email drafting, account
              overview, and schedule trigger registration.

Flags (with common):
  --dry-run, -n    Print planned additions without touching disk
  --force,   -f    Overwrite conflicting JSON values in .prettierrc

Flags (with backup-branch):
  --no-push       Skip pushing the new branch to origin
  --force         Recreate the branch at current HEAD if it already exists

Flags (with PR-merge):
  --squash        Squash-merge with PR title/body + file heatmap
  --single-commit, --flat
                  Flat single-commit merge with a one-line message
  --no-push       Merge locally only; do not push the target

Flags (with send-email):
  --to <addr>     Recipient(s), comma-separated (required)
  --subject <s>   Subject (required unless --file provides it)
  --body <b>      Plain-text body
  --file <path>   Email file (.yaml/.yml, .json, .md, or .txt)
  --html          Send the body as HTML

Flags (with schedule-send-email):
  --to <addr>     Recipient address (required for register)
  --task <task>   Task to watch: release|pull|clone|scan|push|backup
                  (required for register)
  --name <label>  Trigger label; defaults to the task name
  --list          List all registered triggers
  --remove        Remove the trigger named by --name

Examples:
  gitmap space common
  gitmap space common --dry-run
  gitmap space backup-branch "CLI help displayer overhaul"
  gitmap space backup-branch "hotfix" --no-push
  gitmap space PR-merge main
  gitmap space PR-merge
  gitmap space PR-merge main --squash
  gitmap space send-email --to a@x.com --subject "Hi" --body "Hello"
  gitmap space schedule-send-email --to a@x.com --task release
  gitmap space --help
`

// spaceCommonUsage is printed for `gitmap space common --help`.
const spaceCommonUsage = `Usage: gitmap space common [flags]

Applies the curated common baselines (.gitignore, .gitattributes,
.prettierignore, .prettierrc) and runs 'git lfs install --local',
all in one pass. Same logic as 'gitmap commons'.

Behavior:
  - Line-based targets append MISSING lines only; existing entries
    are preserved verbatim. Safe to re-run.
  - .prettierrc is JSON key-union: missing keys added, existing kept
    unless --force is passed.
  - Idempotent — a second run on an unchanged repo writes nothing.

Flags:
  --dry-run, -n    Print planned additions without touching disk
  --force,   -f    Overwrite conflicting JSON values in .prettierrc

Examples:
  gitmap space common
  gitmap space common --dry-run
  gitmap space common --force
`

// spaceBackupBranchUsage is printed for `gitmap space backup-branch --help`.
const spaceBackupBranchUsage = `Usage: gitmap space backup-branch "<task string>" [flags]

Creates a backup branch named backup/<slug> from the current HEAD and
pushes it to origin. The slug is derived from the task string:
lowercase, spaces/underscores become hyphens, only [a-z0-9-] kept,
repeats collapsed, max 60 chars.

Behavior:
  - The working tree must be clean; a dirty tree is always refused.
  - Refuses to overwrite an existing branch unless --force is passed.
  - Prints backup/<slug> @ <short-sha> on success.

Flags:
  --no-push       Skip pushing the new branch to origin
  --force         Recreate the branch at current HEAD if it already exists

Examples:
  gitmap space backup-branch "CLI help displayer overhaul"
  gitmap space backup-branch "hotfix rollout" --no-push
`

// spaceSendEmailUsage is printed for `gitmap space send-email --help`.
const spaceSendEmailUsage = `Usage: gitmap space send-email --to <addr[,addr]> [flags]

Sends an email. Delegates to the go-email-reader binary when it is
installed (PATH, ~/go/bin, or the go_email_reader_bin setting);
otherwise sends through the built-in SMTP sender using the configured
email settings (see 'gitmap email setup').

Behavior:
  - --to is required (comma-separated addresses allowed).
  - --subject is required unless --file provides it.
  - --file parses YAML (.yaml/.yml), JSON (.json), Markdown (.md),
    or plain text: keys to, cc, bcc, subject, body, html.
  - Explicit flags override the corresponding --file values.
  - The SMTP password is always prompted interactively; it is never
    stored, logged, or passed on a command line.

Flags:
  --to <addr>        Recipient(s), comma-separated (required)
  --subject <text>   Subject (required unless --file provides it)
  --body <text>      Plain-text body
  --file <path>      Email file: .yaml/.yml, .json, .md, or .txt
  --html             Send the body as HTML
  --expand-links     Expand shortened/redirecting URLs in email body/subject
  --json             Output send execution result as JSON

Examples:
  gitmap space send-email --to a@x.com --subject "Hi" --body "Hello"
  gitmap space send-email --file invite.yaml --expand-links --json
  gitmap space send-email --to a@x.com,b@x.com --file note.md --html
`

// dispatchSpace routes `gitmap space <subcommand>`; the `common`
// subcommand reuses the same baseline logic as `gitmap commons`, and
// `backup-branch` delegates to the cmdspace implementation.
func DispatchSpace(command string) (bool, error) {
	if command != constants.CmdSpace {
		return false, nil
	}

	rest := os.Args[2:]
	if len(rest) == 0 {
		fmt.Print(spaceUsage)

		return true, nil
	}

	sub := rest[0]
	if sub == "--help" || sub == "-h" {
		fmt.Print(spaceUsage)

		return true, nil
	}

	if sub == constants.CmdSpaceBackupBranch {
		return true, runSpaceBackupBranch(rest[1:])
	}

	if sub == constants.CmdSpacePRMerge {
		return true, runSpacePRMerge(rest[1:])
	}

	if sub == constants.CmdSpaceSendEmail && hasUIFlag(rest[1:]) {
		return true, runSpaceSendEmailUI(filterUIFlag(rest[1:]))
	}

	if sub == constants.CmdSpaceSendEmail {
		return true, runSpaceSendEmail(rest[1:])
	}

	if sub == constants.CmdSpaceScheduleSendEmail {
		return true, runSpaceScheduleSendEmail(rest[1:])
	}

	if sub == "email" {
		return true, runSpaceEmail(rest[1:])
	}

	if sub != "common" {
		fmt.Print(spaceUsage)

		return true, apperror.NewValidationError(fmt.Sprintf("unknown space subcommand: %s", sub))
	}

	subRest := rest[1:]
	for _, a := range subRest {
		if a == "--help" || a == "-h" {
			fmt.Print(spaceCommonUsage)

			return true, nil
		}
	}

	dry, force := cmdsync.ParseCommonFlags(subRest)

	cmdsync.RunCommonLines(constants.GitIgnoreFile, cmdsync.DefaultGitignoreBaseline, dry)
	cmdsync.RunCommonLines(".gitattributes", cmdsync.DefaultGitattributesBaseline, dry)
	cmdsync.RunCommonLFSInstall(dry)
	cmdsync.RunCommonLines(".prettierignore", cmdsync.DefaultPrettierignoreBaseline, dry)
	cmdsync.RunCommonPrettierRC(dry, force)

	return true, nil
}

// runSpaceBackupBranch prints backup-branch help on --help/-h, else
// delegates to the cmdspace implementation.
func runSpaceBackupBranch(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(spaceBackupBranchUsage)

			return nil
		}
	}

	return RunBackupBranch(args)
}

// runSpacePRMerge prints PR-merge help on --help/-h, else delegates to
// the cmdspace implementation.
func runSpacePRMerge(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printPRMergeHelp()

			return nil
		}
	}

	return RunPRMerge(args)
}

// runSpaceSendEmail prints send-email help on --help/-h, else delegates
// to the cmdspace implementation.
func runSpaceSendEmail(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(spaceSendEmailUsage)

			return nil
		}
	}

	return RunSendEmail(args)
}

// runSpaceScheduleSendEmail prints schedule-send-email help on --help/-h,
// else delegates to the cmdspace implementation.
func runSpaceScheduleSendEmail(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			printScheduleSendEmailHelp()

			return nil
		}
	}

	return RunScheduleSendEmail(args)
}
