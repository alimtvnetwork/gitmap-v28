// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import "fmt"

// scheduleSendEmailUsage is printed for
// `gitmap space schedule-send-email --help`. It follows the existing space
// usage-string pattern until the spec 243.3 displayer migration.
const scheduleSendEmailUsage = `Usage: gitmap space schedule-send-email --to <email> --task <task> [--name <label>]

Registers an email trigger: when the watched task completes, gitmap
sends an HTML status email (hostname, user, outbound IP, repo, branch,
outcome, timestamps) to the registered address.

Tasks: release, pull, clone, scan, push, backup

Behavior:
  - Delivery prefers the go-email-reader binary (PATH, then ~/go/bin);
    without it, gitmap falls back to built-in SMTP using the existing
    email.smtp_host / email.from settings (STARTTLS on port 587).
  - The SMTP password is never read from stored config: the fallback
    collects it via a masked runtime prompt. Triggers are
    non-interactive, so with no binary and no TTY the send is skipped
    (recorded on stderr) instead of blocking for input.
  - Trigger emails never fail the operation that fired them; send
    failures are reported on stderr only.

Flags:
  --to <email>     Recipient address (required for register)
  --task <task>    Task to watch: release|pull|clone|scan|push|backup
                   (required for register)
  --name <label>   Trigger label; defaults to the task name
  --list           List all registered triggers
  --remove         Remove the trigger named by --name

Examples:
  gitmap space schedule-send-email --to me@example.com --task release
  gitmap space schedule-send-email --to me@example.com --task pull --name nightly-pull
  gitmap space schedule-send-email --list
  gitmap space schedule-send-email --remove --name nightly-pull
`

// scheduleSendEmailUsageLine is the one-line usage embedded in flag errors.
const scheduleSendEmailUsageLine = "gitmap space schedule-send-email --to <email> --task <release|pull|clone|scan|push|backup> [--name <label>] | --list | --remove --name <label>"

// printScheduleSendEmailHelp prints the schedule-send-email usage text.
func printScheduleSendEmailHelp() {
	fmt.Print(scheduleSendEmailUsage)
}
