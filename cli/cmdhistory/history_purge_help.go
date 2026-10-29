package cmdhistory

import "fmt"

const purgeHelpText = `gitmap history purge — remove a path or commits from history

  --repo <path>            Repository directory. Default is the current directory.
  --file <path>            Relative file or directory path to drop from history.
  --folder <path>          Relative directory to drop from history.
  --commit <sha[,sha]>     One commit SHA, or a comma-separated list.
  --range <a..b>           Commit range, expanded with git rev-list.
  --dry-run, -n            Preview the commits. Does not rewrite or push.
  --yes, -y                Skip the confirm prompt. Required for a real rewrite.
  --json                   Print a JSON result.
  --release                Also scrub matching GitHub release notes.
  --no-backup              Skip the blob backup.

Examples:
  gitmap history purge --file .cursor/skills/execute-parent-task-with-n-steps-v7 --dry-run
  gitmap history purge --commit abc123,def456 --dry-run
  gitmap history purge --range abc123..def456 --dry-run
`

func printPurgeHelp() {
	fmt.Print(purgeHelpText)
}
