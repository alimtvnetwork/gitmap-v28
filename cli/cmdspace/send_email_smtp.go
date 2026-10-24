package cmdspace

import (
	"crypto/tls"
	"fmt"
	"net"
	netmail "net/mail"
	netsmtp "net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// Built-in SMTP sender constants.
const (
	smtpDefaultPort  = 587
	smtpImplicitPort = 465
	smtpDialTimeout  = 15 * time.Second
)

// promptEmailPassword is the SMTP password prompt seam; tests stub it.
// The product path always prompts interactively — never env, never stored.
var promptEmailPassword = readEmailPassword

// sendEmailViaSMTP delivers the message through the built-in SMTP sender.
func sendEmailViaSMTP(msg EmailMessage) error {
	host, port := resolveSMTPHostPort()
	password, passErr := promptEmailPassword()
	if passErr != nil {
		return passErr
	}

	from := resolveSMTPFrom()
	if deliverErr := deliverViaSMTP(host, port, from, password, msg); deliverErr != nil {
		return deliverErr
	}
	fmt.Printf("sent via %s:%d to %s\n", host, port, strings.Join(msg.To, ", "))

	return nil
}

// deliverViaSMTP connects, authenticates with PLAIN auth, and transmits.
func deliverViaSMTP(host string, port int, from, password string, msg EmailMessage) error {
	client, connErr := dialSMTPClient(host, port)
	if connErr != nil {
		return connErr
	}
	defer client.Close()

	if authErr := authSMTPClient(client, host, from, password); authErr != nil {
		return authErr
	}

	return transmitSMTPMessage(client, from, msg)
}

// resolveSMTPHostPort splits the configured `host:port` SMTP endpoint.
func resolveSMTPHostPort() (string, int) {
	raw := readEmailSetting("email.smtp_host", "SMTP_HOST", "smtp.gmail.com:587")
	splitHost, splitPort, splitErr := net.SplitHostPort(raw)
	if splitErr != nil {
		return raw, smtpDefaultPort
	}

	return splitHost, parseSMTPPort(splitPort, smtpDefaultPort)
}

func parseSMTPPort(rawPort string, fallback int) int {
	parsed, err := strconv.Atoi(rawPort)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

// resolveSMTPFrom reads the configured sender address.
func resolveSMTPFrom() string {
	return readEmailSetting("email.from", "SMTP_FROM", "gitmap-bot@localhost")
}

// readEmailSetting resolves an email setting with the same precedence the
// agy email status uses: split-DB setting, then global config variable,
// then environment, then the fallback default.
func readEmailSetting(key, envName, fallback string) string {
	if dbVal := readStoreSetting(key); len(strings.TrimSpace(dbVal)) > 0 {
		return strings.TrimSpace(dbVal)
	}

	if cfgVal, isFound := config.GetVariable("global", key); isFound && len(strings.TrimSpace(cfgVal)) > 0 {
		return strings.TrimSpace(cfgVal)
	}

	if envVal := strings.TrimSpace(os.Getenv(envName)); len(envVal) > 0 {
		return envVal
	}

	return fallback
}

// readStoreSetting reads one key from the default split-DB; empty on any failure.
func readStoreSetting(key string) string {
	dbConn, openErr := store.OpenDefault()
	if openErr != nil {
		return ""
	}
	defer dbConn.Close()

	return dbConn.GetSetting(key)
}

// readEmailPassword prompts for the SMTP password without echoing it.
// Fail-closed on non-TTY stdin; the password is never stored or logged.
func readEmailPassword() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", apperror.NewExecutionError("cannot prompt for the SMTP password: stdin is not a terminal")
	}

	fmt.Print("SMTP password: ")
	raw, readErr := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if readErr != nil {
		return "", apperror.WrapSimple(readErr, "read SMTP password:")
	}

	return string(raw), nil
}

// dialSMTPClient opens the SMTP connection: implicit TLS on port 465,
// STARTTLS-upgraded plain TCP when the server offers it otherwise.
func dialSMTPClient(host string, port int) (*netsmtp.Client, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	if port == smtpImplicitPort {
		return dialImplicitTLS(addr, host)
	}

	return dialStartTLS(addr, host)
}

// dialImplicitTLS opens an implicit-TLS SMTP connection (port 465).
func dialImplicitTLS(addr, host string) (*netsmtp.Client, error) {
	dialer := &net.Dialer{Timeout: smtpDialTimeout}
	conn, dialErr := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: host})
	if dialErr != nil {
		return nil, apperror.WrapSimple(dialErr, "smtp tls dial "+addr+":")
	}

	return newSMTPClient(conn, host)
}

// dialStartTLS opens a plain SMTP connection and upgrades via STARTTLS when offered.
func dialStartTLS(addr, host string) (*netsmtp.Client, error) {
	conn, dialErr := net.DialTimeout("tcp", addr, smtpDialTimeout)
	if dialErr != nil {
		return nil, apperror.WrapSimple(dialErr, "smtp dial "+addr+":")
	}

	client, clientErr := newSMTPClient(conn, host)
	if clientErr != nil {
		return nil, clientErr
	}

	if startErr := maybeStartTLS(client, host); startErr != nil {
		return nil, startErr
	}

	return client, nil
}

// newSMTPClient wraps a connection in an SMTP client; closes conn on failure.
func newSMTPClient(conn net.Conn, host string) (*netsmtp.Client, error) {
	client, clientErr := netsmtp.NewClient(conn, host)
	if clientErr != nil {
		_ = conn.Close()

		return nil, apperror.WrapSimple(clientErr, "smtp client init:")
	}

	return client, nil
}

// maybeStartTLS upgrades the connection when the server advertises STARTTLS.
func maybeStartTLS(client *netsmtp.Client, host string) error {
	offered, _ := client.Extension("STARTTLS")
	if !offered {
		return nil
	}

	if startErr := client.StartTLS(&tls.Config{ServerName: host}); startErr != nil {
		_ = client.Close()

		return apperror.WrapSimple(startErr, "smtp starttls:")
	}

	return nil
}

// authSMTPClient authenticates with PLAIN auth; refuses when the server
// does not advertise AUTH (auth is mandatory, never nil).
func authSMTPClient(client *netsmtp.Client, host, from, password string) error {
	offered, _ := client.Extension("AUTH")
	if !offered {
		return apperror.NewSimple("smtp server does not advertise AUTH", "E_SEND_EMAIL_NO_AUTH")
	}

	auth := netsmtp.PlainAuth("", from, password, host)
	if authErr := client.Auth(auth); authErr != nil {
		return apperror.WrapSimple(authErr, "smtp auth for "+from+":")
	}

	return nil
}

// transmitSMTPMessage runs the MAIL FROM / RCPT TO / DATA envelope.
func transmitSMTPMessage(client *netsmtp.Client, from string, msg EmailMessage) error {
	fromAddr, parseErr := netmail.ParseAddress(from)
	if parseErr != nil {
		return apperror.WrapSimple(parseErr, "invalid from address "+from+":")
	}

	recipients := emailRecipients(msg)
	if envelopeErr := sendEnvelope(client, fromAddr.Address, recipients); envelopeErr != nil {
		return envelopeErr
	}

	return sendData(client, buildSMTPMessage(fromAddr.String(), msg))
}

// emailRecipients flattens To/Cc/Bcc into the SMTP envelope recipient list.
func emailRecipients(msg EmailMessage) []string {
	recipients := append([]string{}, msg.To...)
	recipients = append(recipients, msg.Cc...)
	recipients = append(recipients, msg.Bcc...)

	return recipients
}

// sendEnvelope runs MAIL FROM plus one RCPT TO per recipient.
func sendEnvelope(client *netsmtp.Client, from string, recipients []string) error {
	if mailErr := client.Mail(from); mailErr != nil {
		return apperror.WrapSimple(mailErr, "smtp MAIL FROM:")
	}

	for _, rcpt := range recipients {
		if rcptErr := client.Rcpt(rcpt); rcptErr != nil {
			return apperror.WrapSimple(rcptErr, "smtp RCPT TO "+rcpt+":")
		}
	}

	return nil
}

// sendData writes the DATA payload and quits the session.
func sendData(client *netsmtp.Client, payload string) error {
	writer, dataErr := client.Data()
	if dataErr != nil {
		return apperror.WrapSimple(dataErr, "smtp DATA:")
	}

	if _, writeErr := writer.Write([]byte(payload)); writeErr != nil {
		return apperror.WrapSimple(writeErr, "smtp write:")
	}

	if closeErr := writer.Close(); closeErr != nil {
		return apperror.WrapSimple(closeErr, "smtp data close:")
	}

	return client.Quit()
}

// buildSMTPMessage renders the RFC 5322 message; Bcc stays out of the headers.
func buildSMTPMessage(from string, msg EmailMessage) string {
	contentType := "text/plain; charset=UTF-8"
	if msg.IsHTML {
		contentType = "text/html; charset=UTF-8"
	}

	var sb strings.Builder
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + strings.Join(msg.To, ", ") + "\r\n")
	if len(msg.Cc) > 0 {
		sb.WriteString("Cc: " + strings.Join(msg.Cc, ", ") + "\r\n")
	}
	sb.WriteString("Subject: " + msg.Subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString(constants.HeaderContentType + ": " + contentType + "\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(msg.Body)

	return sb.String()
}
