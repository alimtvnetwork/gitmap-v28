package constants

import "os"

// Command names, subcommands, and aliases for cmdhooks.
const (
	CmdHooks        = "hooks"
	CmdHooksAlias   = "hkin"
	SubCmdInstall   = "install"
	SubCmdUninstall = "uninstall"
	SubCmdRun       = "run"
	SubCmdCheck     = "check"
	SubCmdList      = "list"
	SubCmdStatus    = "status"
	SubCmdHelp      = "help"
)

// Flag names used for CLI flag parsing.
const (
	FlagAll     = "all"
	FlagWorkdir = "workdir"
	FlagDryRun  = "dry-run"
	FlagYes     = "yes"
	FlagHelp    = "help"
)

// Flag CLI arguments including long and short variants.
const (
	FlagArgAll       = "--all"
	FlagArgWorkdir   = "--workdir"
	FlagArgDryRun    = "--dry-run"
	FlagArgYes       = "--yes"
	FlagArgHelp      = "--help"
	FlagArgHelpShort = "-h"
)

// Flag descriptions for help text and flag sets.
const (
	FlagDescAll     = "target all repos in the workdir"
	FlagDescWorkdir = "repo workdir"
	FlagDescDryRun  = "print the plan without installing"
	FlagDescYes     = "actually install the hooks"
	FlagDescHelp    = "show this help"
)

// FlagSet identifiers.
const (
	FlagSetHooksInstall = "hooks install"
)

// Directory names, file names, and path components.
const (
	DirGit              = ".git"
	DirHooks            = "hooks"
	FileCommitMsg       = "commit-msg"
	DefaultReposDir     = "repos"
	DefaultWorkspaceDir = "workspace"
	SentinelCurrentDir  = "."
	PathHookCommitMsg   = ".git/hooks/commit-msg"
)

// Environment variables.
const (
	EnvGitmapWorkdir = "GITMAP_WORKDIR"
)

// File permissions.
const (
	HookFilePerm os.FileMode = 0755
)

// Guard marker identifying hooks managed by GitMap.
const (
	GuardMarker = "gitmap commit-msg guard"
)

// Shell hook script templates.
const (
	CommitMsgGuardTemplate = `#!/bin/sh
# gitmap commit-msg guard (installed by ` + "`gitmap hooks install --all`" + `)
# House rule: the commit subject line must not contain a colon.
subject=$(head -n1 "$1")
case "$subject" in
*:*)
  echo "commit-msg guard: subject must not contain ':'" >&2
  echo "  got: $subject" >&2
  exit 1
  ;;
esac
exit 0
`
)

// Help and usage documentation.
const (
	HelpTextHooksInstall = `hooks install — install the commit-msg guard hook across repos.

Description:
  Installs the house commit-msg guard (no colons in the subject line)
  into every repo in the workdir. Safe by default: --dry-run prints the
  plan; a real install only happens with --yes.

Usage:
  gitmap hooks install --all [--workdir <dir>] [--dry-run] [--yes]
  gitmap hkin --all [--yes]

Flags:
  --all           required: target all repos in the workdir
  --workdir <dir> repo workdir (default: $GITMAP_WORKDIR or ~/workspace/repos)
  --dry-run       print the install plan without writing (default: true)
  --yes           actually install the hooks
  --help, -h      show this help

Example:
  gitmap hooks install --all --dry-run
`
)

// Plan detail descriptions.
const (
	DetailAlreadyGuarded = "guard already installed"
	DetailHookRelPath    = ".git/hooks/commit-msg"
)

// Informational and formatted log messages.
const (
	MsgDryRunNotice       = "dry-run: no hooks were installed. re-run with --yes to install."
	MsgNoReposFoundFmt    = "hooks install: no git repos found in %s\n"
	MsgPlanHeaderFmt      = "hooks install plan: %d repos (%d would install, %d already guarded, %d skipped)\n"
	MsgPlanRowFmt         = "  %-9s %-32s %s\n"
	MsgInstallSummaryFmt  = "hooks install: installed commit-msg guard in %d/%d repos (%d already had it, %d skipped).\n"
)

// Error messages and wrapped error context prefixes.
const (
	ErrMissingSubcommand    = "hooks: missing subcommand (install)"
	ErrUnknownSubcommandFmt = "hooks: unknown subcommand %q (install)"
	ErrMissingAllFlag       = "missing required flag: --all"
	ErrWrapFlags            = "hooks.install.flags"
	ErrWrapWritePrefix      = "hooks.install.write: "
)
