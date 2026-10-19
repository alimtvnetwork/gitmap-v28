// Package cmdagy — agy_email_subcommands.go implements Email notification setup and testing.
package cmdagy

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// AgyEmailCmd provides the Cobra entrypoint for gitmap [agy] email.
var AgyEmailCmd = &cobra.Command{
	Use:                "email",
	Aliases:            []string{"mail", "smtp"},
	Short:              "Speed setup, status, and testing for Email notifications",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyEmailCLI(args)
	},
}

// RunAgyEmailCLI routes Email speed setup and test subcommands.
func RunAgyEmailCLI(args []string) error {
	if len(args) == 0 || isTgEmailHelpArg(args[0]) {
		RenderAgyEmailHelp()
		return executeEmailStatus(hasSpeedFlag(args, "--json"))
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	return dispatchEmailSubcommand(sub, args[1:])
}

func dispatchEmailSubcommand(sub string, rest []string) error {
	switch sub {
	case "setup", "config", "configure", "set":
		return executeEmailSetup(rest)
	case "status", "st", "ls", "info":
		return executeEmailStatus(hasSpeedFlag(rest, "--json"))
	case "test", "send":
		return executeEmailTest(rest)
	default:
		RenderAgyEmailHelp()
		return nil
	}
}

func executeEmailSetup(args []string) error {
	smtpHost := extractSpeedFlag(args, "--smtp", "-s")
	fromAddr := extractSpeedFlag(args, "--from", "-f")
	toAddr := extractSpeedFlag(args, "--to", "-t")
	pass := extractSpeedFlag(args, "--password", "-p")
	saveEmailSpeedConfig(smtpHost, fromAddr, toAddr, pass)
	fmt.Printf("%s✔ Configured Email speed settings (smtp=%s, from=%s, to=%s)%s\n",
		constants.ColorGreen, resolveEmailSMTP(), resolveEmailFrom(), resolveEmailTo(), constants.ColorReset)
	return nil
}

func saveEmailSpeedConfig(smtpHost, fromAddr, toAddr, pass string) {
	setGlobalIfNonEmpty("email.smtp_host", "SMTP_HOST", smtpHost)
	setGlobalIfNonEmpty("email.from", "SMTP_FROM", fromAddr)
	setGlobalIfNonEmpty("email.to", "AGM_NOTIFY_EMAIL", toAddr)
	setGlobalIfNonEmpty("email.password", "SMTP_PASSWORD", pass)
}

func executeEmailStatus(isJSON bool) error {
	st := collectSpeedSettingsStatus()
	if isJSON {
		return printSpeedJSON(st)
	}
	fmt.Printf("%s● GitMap Email Speed Notification Settings%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  SMTP Host:Port:    %s\n", st.EmailSMTPHost)
	fmt.Printf("  Sender (From):     %s\n", st.EmailFrom)
	fmt.Printf("  Recipient (To):    %s\n", st.EmailTo)
	fmt.Printf("  Ready:             %v\n", st.EmailReady)
	return nil
}

func executeEmailTest(args []string) error {
	st := collectSpeedSettingsStatus()
	isSent := sendEmailSwitchNotice("[GitMap] Speed Email Test", "GitMap Email Speed Setup verification succeeded.")
	fmt.Printf("%s✔ Test email dispatched to %s via %s (delivered=%v)%s\n",
		constants.ColorGreen, st.EmailTo, st.EmailSMTPHost, isSent, constants.ColorReset)
	return nil
}

func resolveEmailSMTP() string {
	return readSpeedVar("email.smtp_host", defaultIfEmpty(os.Getenv("SMTP_HOST"), "smtp.gmail.com:587"))
}

func resolveEmailFrom() string {
	return readSpeedVar("email.from", defaultIfEmpty(os.Getenv("SMTP_FROM"), "gitmap-bot@localhost"))
}

func resolveEmailTo() string {
	return readSpeedVar("email.to", defaultIfEmpty(os.Getenv("AGM_NOTIFY_EMAIL"), "dev@localhost"))
}

// RenderAgyEmailHelp renders the boxed help menu for gitmap email / gitmap agy email.
func RenderAgyEmailHelp() {
	termout.RenderMenu(termout.HelpMenu{
		Title: "Email Speed Notification Setup (gitmap email)",
		UsageLines: []string{
			"gitmap email <setup|status|test|help> [flags]",
			"gitmap agy email <setup|status|test|help> [flags]",
		},
		Sections: []termout.HelpSection{
			{
				Title: "Email Speed Subcommands",
				Entries: []termout.CommandEntry{
					{Command: "setup --smtp <h:p> --from <f> --to <t>", Description: "Configure SMTP host, sender, recipient, and app password"},
					{Command: "status [--json]", Description: "Display current Email speed notification configuration"},
					{Command: "test", Description: "Send a test verification email using configured SMTP settings"},
					{Command: "help", Description: "Display this boxed Email speed setup help menu"},
				},
			},
		},
	})
}
