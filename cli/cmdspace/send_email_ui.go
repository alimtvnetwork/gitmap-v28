package cmdspace

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

var uiInputReader io.Reader = os.Stdin

type EmailAccountInfo struct {
	Alias, Email, ImapHost, SmtpHost, SmtpTLSMode string
	ImapPort, SmtpPort                            int
	ImapTLS                                       bool
}

type sendEmailUIOptions struct {
	to, subject, body, task          string
	isHTML, isDryRun, isInteractive bool
}

type secretAccountEntry struct {
	Alias, Email string `json:"alias"`
	Imap         struct { Host string `json:"host"`; Port int `json:"port"`; TLS bool `json:"tls"` } `json:"imap"`
	Smtp         struct { Host, TLSMode string `json:"host,tls_mode"`; Port int `json:"port"` } `json:"smtp"`
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
	case arg == "--dry-run" || arg == "-n":
		opts.isDryRun = true
	case arg == "--non-interactive":
		opts.isInteractive = false
	default:
		return i, fmt.Errorf("unknown UI flag %q: run 'gitmap space email ui --help'", arg)
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

func LoadEmailAccounts() []EmailAccountInfo {
	if accounts := queryGoEmailReaderAccounts(); len(accounts) > 0 {
		return accounts
	}
	data, err := os.ReadFile(findRepoSecretFile(filepath.Join("repo-secrets", "11-go-email-reader", "01-email-setup.json")))
	if err != nil {
		return nil
	}
	var doc struct { Accounts []secretAccountEntry `json:"accounts"` }
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil
	}
	return mapSecretAccounts(doc.Accounts)
}

func queryGoEmailReaderAccounts() []EmailAccountInfo {
	bin := resolveGoEmailReaderBin()
	if len(bin) == 0 {
		return nil
	}
	out, err := exec.Command(bin, "list").Output()
	if err != nil {
		return nil
	}
	return parseGoEmailReaderListOutput(string(out))
}

func parseGoEmailReaderListOutput(output string) []EmailAccountInfo {
	if strings.Contains(output, "no accounts stored") {
		return nil
	}
	var accounts []EmailAccountInfo
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") {
			accounts = append(accounts, parseAccountLine(trimmed))
		}
	}
	return accounts
}

func parseAccountLine(line string) EmailAccountInfo {
	content := strings.TrimPrefix(line, "- ")
	if parts := strings.SplitN(content, "(", 2); len(parts) == 2 {
		return EmailAccountInfo{Email: strings.TrimSpace(parts[0]), Alias: strings.TrimSuffix(strings.TrimSpace(parts[1]), ")")}
	}
	return EmailAccountInfo{Email: strings.TrimSpace(content), Alias: "-"}
}

func mapSecretAccounts(entries []secretAccountEntry) []EmailAccountInfo {
	res := make([]EmailAccountInfo, 0, len(entries))
	for _, a := range entries {
		res = append(res, EmailAccountInfo{
			Alias: a.Alias, Email: a.Email,
			ImapHost: a.Imap.Host, ImapPort: a.Imap.Port, ImapTLS: a.Imap.TLS,
			SmtpHost: a.Smtp.Host, SmtpPort: a.Smtp.Port, SmtpTLSMode: a.Smtp.TLSMode,
		})
	}
	return res
}

func LoadDefaultNotifier() string {
	if notifier := queryGoEmailReaderNotifier(); len(notifier) > 0 {
		return notifier
	}
	data, err := os.ReadFile(findRepoSecretFile(filepath.Join("repo-secrets", "11-go-email-reader", "02-default-notifier.json")))
	if err != nil {
		return ""
	}
	var doc struct { NotifierEmail string `json:"notifier_email"` }
	_ = json.Unmarshal(data, &doc)
	return strings.TrimSpace(doc.NotifierEmail)
}

func queryGoEmailReaderNotifier() string {
	bin := resolveGoEmailReaderBin()
	if len(bin) == 0 {
		return ""
	}
	out, err := exec.Command(bin, "notifier", "get").Output()
	res := strings.TrimSpace(string(out))
	if err != nil || strings.Contains(res, "no notifier") {
		return ""
	}
	return res
}

func findRepoSecretFile(relPath string) string {
	for _, cand := range []string{relPath, filepath.Join("..", relPath), filepath.Join(".", relPath)} {
		if isSecretFile(cand) {
			return cand
		}
	}
	return ""
}

func isSecretFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func RenderEmailUIBanner() string {
	return fmt.Sprintf("\n%s  ╔══════════════════════════════════════════════════════════╗\n  ║          GitMap Space — Email & Trigger Studio           ║\n  ╚══════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func RenderAccountTable(accounts []EmailAccountInfo) string {
	return termout.RenderTable(termout.TableConfig{
		HeaderColor: constants.ColorYellow,
		BorderColor: constants.ColorDim,
		Columns: []termout.Column{
			{Title: "Alias", MinWidth: 14},
			{Title: "Email Address", MinWidth: 28},
			{Title: "IMAP Server", MinWidth: 24},
			{Title: "SMTP Server", MinWidth: 24},
		},
		Rows: buildAccountTableRows(accounts),
	})
}

func buildAccountTableRows(accounts []EmailAccountInfo) []termout.Row {
	if len(accounts) == 0 {
		return []termout.Row{{Cells: []string{"-", "No accounts stored", "-", "-"}, Color: constants.ColorDim}}
	}
	rows := make([]termout.Row, 0, len(accounts))
	for _, a := range accounts {
		rows = append(rows, termout.Row{
			Cells: []string{a.Alias, a.Email, fmt.Sprintf("%s:%d (TLS=%t)", a.ImapHost, a.ImapPort, a.ImapTLS), fmt.Sprintf("%s:%d (%s)", a.SmtpHost, a.SmtpPort, a.SmtpTLSMode)},
			Color: constants.ColorWhite,
		})
	}
	return rows
}

func RenderNotifierLine(notifier string) string {
	if len(notifier) == 0 {
		notifier = "none"
	}
	return fmt.Sprintf("  Default Notifier: %s%s%s\n\n", constants.ColorGreen, notifier, constants.ColorReset)
}

func RenderDraftPreview(opts sendEmailUIOptions) string {
	format := "Plain Text"
	if opts.isHTML {
		format = "HTML"
	}
	return fmt.Sprintf("  Draft Preview:\n    • Recipient: %s\n    • Subject:   %s\n    • Body:      %s\n    • Format:    %s\n    • Trigger:   %s\n",
		defaultIfBlank(opts.to, "(none)"), defaultIfBlank(opts.subject, "(none)"),
		defaultIfBlank(opts.body, "(none)"), format, defaultIfBlank(opts.task, "(none)"))
}

func defaultIfBlank(val, fallback string) string {
	if len(strings.TrimSpace(val)) == 0 {
		return fallback
	}
	return val
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
	if len(opts.task) > 0 {
		task, err := ParseEmailTaskType(opts.task)
		if err != nil {
			return err
		}
		if pErr := persistScheduleTrigger(EmailTrigger{Name: scheduleTriggerName("", task), To: opts.to, Task: task}); pErr != nil {
			return pErr
		}
	}
	if len(opts.to) > 0 && len(opts.subject) > 0 {
		return RunSendEmail(buildSendEmailArgsFromUI(opts))
	}
	return nil
}

func buildSendEmailArgsFromUI(opts sendEmailUIOptions) []string {
	args := []string{"--to", opts.to, "--subject", opts.subject, "--body", opts.body}
	if opts.isHTML {
		args = append(args, "--html")
	}
	return args
}
