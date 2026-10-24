// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const (
	// emailReaderBinaryName is the go-email-reader CLI used for delivery
	// when present on the machine.
	emailReaderBinaryName = "go-email-reader"
	// emailReaderHTMLFlag marks the body as HTML (Content-Type: text/html)
	// on `go-email-reader send`.
	emailReaderHTMLFlag = "--html"
	// smtpDefaultHostPort mirrors the agy email settings default.
	smtpDefaultHostPort = "smtp.gmail.com:587"
	// smtpDefaultFrom mirrors the agy email settings default.
	smtpDefaultFrom = "gitmap-bot@localhost"
	// settingSMTPHost and settingSMTPFrom are the existing gitmap email
	// settings keys, read the same way executeEmailStatus reads them.
	settingSMTPHost = "email.smtp_host"
	settingSMTPFrom = "email.from"
	// envSMTPHost and envSMTPFrom are the matching environment overrides.
	envSMTPHost = "SMTP_HOST"
	envSMTPFrom = "SMTP_FROM"
)

// defaultSendTriggerEmail delivers through the go-email-reader binary when
// available (PATH, then ~/go/bin), else through the built-in net/smtp
// sender. It is self-contained: it never imports the sibling worker's
// `space send-email` code.
func defaultSendTriggerEmail(to, subject, htmlBody string) error {
	if bin := resolveEmailReaderBinary(); bin != "" {
		return sendViaEmailReader(bin, to, subject, htmlBody)
	}

	return sendViaBuiltinSMTP(to, subject, htmlBody)
}

// resolveEmailReaderBinary locates the go-email-reader CLI: PATH first,
// then ~/go/bin, else "" (caller falls back to built-in SMTP).
func resolveEmailReaderBinary() string {
	if bin, err := exec.LookPath(emailReaderBinaryName); err == nil {
		return bin
	}

	homeBin := filepath.Join(os.Getenv("HOME"), "go", "bin", emailReaderBinaryName)
	if isRegularFile(homeBin) {
		return homeBin
	}

	return ""
}

// isRegularFile reports whether path exists and is a regular file.
func isRegularFile(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.Mode().IsRegular()
}

// sendViaEmailReader delegates to `go-email-reader send` with the body
// marked as HTML via --html and passed via --body.
func sendViaEmailReader(bin, to, subject, htmlBody string) error {
	cmd := exec.Command(bin, "send", "--to", to, "--subject", subject, emailReaderHTMLFlag, "--body", htmlBody)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.WrapSimple(err, "go-email-reader send: "+string(out))
	}

	return nil
}

// sendViaBuiltinSMTP delivers through net/smtp with STARTTLS using the
// gitmap email settings (email.smtp_host / email.from). The SMTP password
// is ALWAYS collected via a masked runtime prompt — never from stored
// config — and triggers are non-interactive: with no TTY the send is
// SKIPPED (recorded on stderr) instead of blocking on input.
func sendViaBuiltinSMTP(to, subject, htmlBody string) error {
	hasTTY := term.IsTerminal(int(os.Stdin.Fd()))
	if hasTTY {
		return sendViaBuiltinSMTPInteractive(to, subject, htmlBody)
	}

	fmt.Fprintf(os.Stderr, "  ⚠ email trigger to %s skipped: no go-email-reader binary and stdin is not a TTY — the SMTP password cannot be prompted non-interactively\n", to)

	return nil
}

// sendViaBuiltinSMTPInteractive prompts for the SMTP password (masked)
// then dispatches the message.
func sendViaBuiltinSMTPInteractive(to, subject, htmlBody string) error {
	hostPort, from := builtinSMTPSettings()

	fmt.Fprintf(os.Stderr, "SMTP password for %s: ", from)
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return apperror.WrapSimple(err, "read SMTP password:")
	}

	return dispatchBuiltinSMTP(hostPort, from, to, subject, htmlBody, string(pass))
}

// builtinSMTPSettings resolves SMTP host:port and from address from the
// trigger settings DB, then the config file, then env, then defaults —
// the same precedence executeEmailStatus uses.
func builtinSMTPSettings() (hostPort, from string) {
	db, err := openTriggerStore()
	if err == nil {
		defer db.Close()

		return triggerSMTPVar(db, settingSMTPHost, envSMTPHost, smtpDefaultHostPort),
			triggerSMTPVar(db, settingSMTPFrom, envSMTPFrom, smtpDefaultFrom)
	}

	return triggerSMTPVar(nil, settingSMTPHost, envSMTPHost, smtpDefaultHostPort),
		triggerSMTPVar(nil, settingSMTPFrom, envSMTPFrom, smtpDefaultFrom)
}

// triggerSMTPVar resolves one SMTP setting: DB setting, then the global
// config variable, then the environment, then the fallback.
func triggerSMTPVar(db *store.DB, key, envName, fallback string) string {
	if dbVal := readDbSetting(db, key); dbVal != "" {
		return dbVal
	}

	if v, isFound := config.GetVariable("global", key); isFound && len(strings.TrimSpace(v)) > 0 {
		return strings.TrimSpace(v)
	}

	if v := strings.TrimSpace(os.Getenv(envName)); v != "" {
		return v
	}

	return fallback
}

func readDbSetting(db *store.DB, key string) string {
	if db == nil {
		return ""
	}

	return strings.TrimSpace(db.GetSetting(key))
}

// dispatchBuiltinSMTP sends one HTML message via net/smtp (STARTTLS is
// negotiated by SendMail when the server advertises it).
func dispatchBuiltinSMTP(hostPort, from, to, subject, htmlBody, password string) error {
	auth := smtp.PlainAuth("", from, password, smtpHostOnly(hostPort))
	msg := buildMIMEMessage(from, to, subject, htmlBody)
	if err := smtp.SendMail(hostPort, auth, from, []string{to}, msg); err != nil {
		return apperror.WrapSimple(err, "smtp send:")
	}

	return nil
}

// smtpHostOnly strips the port for PlainAuth's server-name check.
func smtpHostOnly(hostPort string) string {
	host, _, err := net.SplitHostPort(hostPort)
	if err != nil {
		return hostPort
	}

	return host
}

// buildMIMEMessage assembles the raw RFC 5322 message with an HTML body.
func buildMIMEMessage(from, to, subject, htmlBody string) []byte {
	var sb strings.Builder
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + to + "\r\n")
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString(constants.HeaderContentType + ": text/html; charset=UTF-8\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(htmlBody)

	return []byte(sb.String())
}
