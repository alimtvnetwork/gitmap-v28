package cmdspace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var uiInputReader io.Reader = os.Stdin

type EmailAccountInfo struct {
	Alias, Email, ImapHost, SmtpHost, SmtpTLSMode string
	ImapPort, SmtpPort                            int
	ImapTLS                                       bool
}

type sendEmailUIOptions struct {
	to, subject, body, task         string
	isHTML, isDryRun, isInteractive bool
}

type secretAccountEntry struct {
	Alias string `json:"alias"`
	Email string `json:"email"`
	Imap  struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		TLS  bool   `json:"tls"`
	} `json:"imap"`
	Smtp struct {
		Host    string `json:"host"`
		TLSMode string `json:"tls_mode"`
		Port    int    `json:"port"`
	} `json:"smtp"`
}

const sendEmailUIUsage = `Usage: gitmap space email ui [flags] (or: gitmap space send-email --ui [flags])

Interactive terminal UI for email drafting, account overview, and schedule trigger registration.

Flags:
  --to <addr>          Recipient email address (default: default notifier)
  --subject <text>     Email subject line
  --body <text>        Plain-text or HTML email body
  --html               Send email formatted as HTML
  --task <task>        Optional trigger task (release|pull|clone|scan|push|backup)
  --dry-run, -n        Preview drafted email and triggers without executing
  --non-interactive    Run without interactive terminal prompts
  --help, -h           Show this help message
`

func printSendEmailUIHelp() { fmt.Print(sendEmailUIUsage) }

func RunSpaceSendEmailUI(args []string) error {
	if hasHelpFlag(args) {
		printSendEmailUIHelp()
		return nil
	}
	opts, err := parseSendEmailUIArgs(args)
	if err != nil {
		return err
	}
	return runSendEmailUIWorkflow(opts)
}

func runSendEmailUIWorkflow(opts sendEmailUIOptions) error {
	fmt.Print(RenderEmailUIBanner())
	fmt.Print(RenderAccountTable(LoadEmailAccounts()))
	fmt.Print(RenderNotifierLine(LoadDefaultNotifier()))
	if opts.isInteractive {
		promptMissingFields(&opts, LoadDefaultNotifier())
	}
	fmt.Println(RenderDraftPreview(opts))
	if opts.isDryRun {
		fmt.Printf("  %s✓ [dry-run] Draft preview complete. No action taken.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}
	return executeSendEmailUI(opts)
}

func parseSendEmailUIArgs(args []string) (sendEmailUIOptions, error) {
	opts := sendEmailUIOptions{isInteractive: true}
	for i := 0; i < len(args); i++ {
		next, err := applySendEmailUIArg(args, i, &opts)
		if err != nil {
			return sendEmailUIOptions{}, err
		}
		i = next
	}
	return opts, nil
}

func applySendEmailUIArg(args []string, i int, opts *sendEmailUIOptions) (int, error) {
	arg := args[i]
	switch {
	case arg == "--to" || strings.HasPrefix(arg, "--to="):
		opts.to, i = nextUIFlagValue(args, i, arg, "--to")
	case arg == "--subject" || strings.HasPrefix(arg, "--subject="):
		opts.subject, i = nextUIFlagValue(args, i, arg, "--subject")
	case arg == "--body" || strings.HasPrefix(arg, "--body="):
		opts.body, i = nextUIFlagValue(args, i, arg, "--body")
	case arg == "--task" || strings.HasPrefix(arg, "--task="):
		opts.task, i = nextUIFlagValue(args, i, arg, "--task")
	case arg == "--html":
		opts.isHTML = true
	case slices.Contains([]string{"--dry-run", "-n"}, arg):
		opts.isDryRun = true
	case arg == "--non-interactive":
		opts.isInteractive = false
	default:
		return i, apperror.NewValidationError(fmt.Sprintf("unknown UI flag %q: run 'gitmap space email ui --help'", arg))
	}
	return i, nil
}

func nextUIFlagValue(args []string, i int, arg, flag string) (string, int) {
	if strings.HasPrefix(arg, flag+"=") {
		return strings.TrimPrefix(arg, flag+"="), i
	}
	if i+1 < len(args) {
		return args[i+1], i + 1
	}
	return "", i
}

func promptMissingFields(opts *sendEmailUIOptions, defaultTo string) {
	reader := bufio.NewReader(uiInputReader)
	if len(opts.to) == 0 {
		opts.to = promptUserValue(reader, "Recipient (--to)", defaultTo)
	}
	if len(opts.subject) == 0 {
		opts.subject = promptUserValue(reader, "Subject (--subject)", "")
	}
	if len(opts.body) == 0 {
		opts.body = promptUserValue(reader, "Body (--body)", "")
	}
	if len(opts.task) == 0 {
		opts.task = promptUserValue(reader, "Trigger task", "")
	}
}

func promptUserValue(reader *bufio.Reader, label, defaultVal string) string {
	if len(defaultVal) > 0 {
		fmt.Printf("  Enter %s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("  Enter %s: ", label)
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		return defaultVal
	}
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return defaultVal
	}
	return trimmed
}

func executeSendEmailUI(opts sendEmailUIOptions) error {
	if err := handleTaskPersistence(opts.task, opts.to); err != nil {
		return err
	}
	if len(opts.to) > 0 && len(opts.subject) > 0 {
		return RunSendEmail(buildSendEmailArgsFromUI(opts))
	}
	return nil
}

func handleTaskPersistence(taskStr, to string) error {
	if len(taskStr) == 0 {
		return nil
	}
	task, err := ParseEmailTaskType(taskStr)
	if err != nil {
		return err
	}
	return persistScheduleTrigger(EmailTrigger{
		Name: scheduleTriggerName("", task),
		To:   to,
		Task: task,
	})
}

func buildSendEmailArgsFromUI(opts sendEmailUIOptions) []string {
	args := []string{"--to", opts.to, "--subject", opts.subject, "--body", opts.body}
	if opts.isHTML {
		args = append(args, "--html")
	}
	return args
}
