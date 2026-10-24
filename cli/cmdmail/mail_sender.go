package cmdmail

import (
	"crypto/tls"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// MailSendOptions configures message transmission and routing.
type MailSendOptions struct {
	AccountName  string   `json:"accountName"`
	From         string   `json:"from"`
	To           []string `json:"to"`
	Cc           []string `json:"cc,omitempty"`
	Bcc          []string `json:"bcc,omitempty"`
	Subject      string   `json:"subject"`
	Body         string   `json:"body"`
	IsHTML       bool     `json:"isHTML"`
	TemplateName string   `json:"templateName,omitempty"`
	Host         string   `json:"host"`
	Port         int      `json:"port"`
	IsTLS        bool     `json:"isTLS"`
	IsDryRun     bool     `json:"isDryRun"`
}

// EmailTransmitter abstracts email dispatch for standard network and testing.
type EmailTransmitter interface {
	SendMail(opts MailSendOptions, rawMessage []byte) error
}

// CurrentTransmitter points to the active email dispatcher.
var CurrentTransmitter EmailTransmitter = &StandardEmailTransmitter{}

// StandardEmailTransmitter implements EmailTransmitter using net/smtp and crypto/tls.
type StandardEmailTransmitter struct{}

// SendMail handles transmission via implicit TLS (465) or STARTTLS (587/standard).
func (t *StandardEmailTransmitter) SendMail(opts MailSendOptions, rawMessage []byte) error {
	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	allRcpts := append(append([]string{}, opts.To...), opts.Cc...)
	allRcpts = append(allRcpts, opts.Bcc...)

	if opts.Port == 465 || opts.IsTLS {
		return sendImplicitTLS(addr, opts.Host, opts.From, allRcpts, rawMessage)
	}
	return sendStartTLSOrPlain(addr, opts.Host, opts.From, allRcpts, rawMessage, opts.Port == 587)
}

func sendImplicitTLS(addr, host, from string, rcpts []string, raw []byte) error {
	tlsCfg := &tls.Config{ServerName: host, InsecureSkipVerify: isLocalHost(host)}
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	return deliverSMTPPayload(client, from, rcpts, raw)
}

func sendStartTLSOrPlain(addr, host, from string, rcpts []string, raw []byte, isStartTLS bool) error {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := maybeStartTLS(client, host, isStartTLS); err != nil {
		return err
	}
	return deliverSMTPPayload(client, from, rcpts, raw)
}

func maybeStartTLS(client *smtp.Client, host string, isStartTLS bool) error {
	if !isStartTLS {
		return nil
	}
	tlsCfg := &tls.Config{ServerName: host, InsecureSkipVerify: isLocalHost(host)}
	return client.StartTLS(tlsCfg)
}

func deliverSMTPPayload(client *smtp.Client, from string, rcpts []string, raw []byte) error {
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, r := range rcpts {
		if err := client.Rcpt(r); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(raw); err != nil {
		_ = w.Close()
		return err
	}
	return w.Close()
}

func isLocalHost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || strings.HasSuffix(host, ".local")
}

// BuildMIMEMessage formats an RFC 822/2822 payload, stripping BCC from headers.
func BuildMIMEMessage(opts MailSendOptions) []byte {
	var sb strings.Builder
	writeHeader(&sb, "From", opts.From)
	writeHeader(&sb, "To", strings.Join(opts.To, ", "))
	if len(opts.Cc) > 0 {
		writeHeader(&sb, "Cc", strings.Join(opts.Cc, ", "))
	}
	writeHeader(&sb, "Subject", opts.Subject)
	writeHeader(&sb, "Date", time.Now().UTC().Format(time.RFC1123Z))
	writeHeader(&sb, "MIME-Version", "1.0")
	if opts.IsHTML {
		writeHeader(&sb, constants.HeaderContentType, "text/html; charset=UTF-8")
	} else {
		writeHeader(&sb, constants.HeaderContentType, "text/plain; charset=UTF-8")
	}
	sb.WriteString("\r\n")
	sb.WriteString(opts.Body)
	return []byte(sb.String())
}

func writeHeader(sb *strings.Builder, key, val string) {
	if val != "" {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", key, val))
	}
}

// MaskEmail transforms raw addresses into privacy format (e.g. u***r@d***n.com).
func MaskEmail(raw string) string {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return "***"
	}
	parts := strings.Split(clean, "@")
	if len(parts) != 2 {
		return maskSimpleString(clean)
	}
	localPart := maskLocalPart(parts[0])
	domainPart := maskDomainPart(parts[1])
	return localPart + "@" + domainPart
}

func maskSimpleString(s string) string {
	if len(s) >= 3 {
		return string(s[0]) + "***" + string(s[len(s)-1])
	}
	return string(s[0]) + "***"
}

func maskLocalPart(local string) string {
	n := len(local)
	if n >= 3 {
		return string(local[0]) + "***" + string(local[n-1])
	}
	if n > 0 {
		return string(local[0]) + "***"
	}
	return "***"
}

func maskDomainPart(domain string) string {
	tokens := strings.Split(domain, ".")
	if len(tokens) <= 1 {
		return maskSimpleString(domain)
	}
	tld := tokens[len(tokens)-1]
	maskedTokens := make([]string, len(tokens)-1)
	for i := 0; i < len(tokens)-1; i++ {
		maskedTokens[i] = maskLabel(tokens[i])
	}
	return strings.Join(maskedTokens, ".") + "." + tld
}

func maskLabel(lbl string) string {
	if len(lbl) >= 3 {
		return string(lbl[0]) + "***" + string(lbl[len(lbl)-1])
	}
	if len(lbl) > 0 {
		return string(lbl[0]) + "***"
	}
	return "***"
}

// MaskEmailList sanitizes a slice of email addresses.
func MaskEmailList(emails []string) []string {
	res := make([]string, len(emails))
	for i, e := range emails {
		res[i] = MaskEmail(e)
	}
	return res
}

// MaskEmailCSV sanitizes comma-separated email list string.
func MaskEmailCSV(csv string) string {
	rawList := strings.Split(csv, ",")
	maskedList := make([]string, 0, len(rawList))
	for _, item := range rawList {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			maskedList = append(maskedList, MaskEmail(trimmed))
		}
	}
	return strings.Join(maskedList, ", ")
}
