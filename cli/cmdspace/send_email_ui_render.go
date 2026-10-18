package cmdspace

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// LoadEmailAccounts loads configured accounts from go-email-reader or repo-secrets fallback.
func LoadEmailAccounts() []EmailAccountInfo {
	if accounts := queryGoEmailReaderAccounts(); len(accounts) > 0 {
		return accounts
	}
	data, err := os.ReadFile(findRepoSecretFile(filepath.Join("repo-secrets", "11-go-email-reader", "01-email-setup.json")))
	if err != nil {
		return nil
	}
	var doc struct {
		Accounts []secretAccountEntry `json:"accounts"`
	}
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

// LoadDefaultNotifier loads default notifier email from go-email-reader or secrets.
func LoadDefaultNotifier() string {
	if notifier := queryGoEmailReaderNotifier(); len(notifier) > 0 {
		return notifier
	}
	data, err := os.ReadFile(findRepoSecretFile(filepath.Join("repo-secrets", "11-go-email-reader", "02-default-notifier.json")))
	if err != nil {
		return ""
	}
	var doc struct {
		NotifierEmail string `json:"notifier_email"`
	}
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
	candidates := []string{
		relPath,
		filepath.Join("..", relPath),
		filepath.Join("..", "..", relPath),
		filepath.Join(`D:\work`, relPath),
	}
	for _, cand := range candidates {
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

// RenderEmailUIBanner renders the terminal header for the email UI.
func RenderEmailUIBanner() string {
	return fmt.Sprintf("\n%s  ╔══════════════════════════════════════════════════════════╗\n  ║          GitMap Space — Email & Trigger Studio           ║\n  ╚══════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}

// RenderAccountTable renders the table of available email accounts.
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

// RenderNotifierLine formats the active default notifier summary.
func RenderNotifierLine(notifier string) string {
	if len(notifier) == 0 {
		notifier = "none"
	}
	return fmt.Sprintf("  Default Notifier: %s%s%s\n\n", constants.ColorGreen, notifier, constants.ColorReset)
}

// RenderDraftPreview formats the draft preview before confirmation.
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
