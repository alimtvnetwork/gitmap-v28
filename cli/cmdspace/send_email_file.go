package cmdspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/linkexpander"
)

// EmailMessage is one parsed email: recipients plus content.
type EmailMessage struct {
	To            []string
	Cc            []string
	Bcc           []string
	Subject       string
	Body          string
	IsHTML        bool
	IsExpanded    bool
	ExpandedCount int
	ExpandedLinks []linkexpander.ExpandedLink
}

// Email-file extensions dispatched by parseEmailFile.
const (
	emailExtYAML = ".yaml"
	emailExtYML  = ".yml"
	emailExtJSON = ".json"
	emailExtMD   = ".md"
)

// parseEmailFile parses an email file by extension: YAML (.yaml/.yml)
// and JSON (.json) read the to/cc/bcc/subject/body/html keys,
// Markdown (.md) takes the first `# ` heading as the subject, and
// anything else is plain text (first non-empty line is the subject,
// an optional `Subject:` prefix is stripped, blank lines preserved).
func parseEmailFile(path string) (EmailMessage, error) {
	var msg EmailMessage

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return msg, apperror.WrapSimple(readErr, "read email file "+path+":")
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		return msg, apperror.NewSimple("empty email file "+path, "E_SEND_EMAIL_EMPTY_FILE")
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case emailExtYAML, emailExtYML:
		return parseEmailYAML(data, path)
	case emailExtJSON:
		return parseEmailJSON(data, path)
	case emailExtMD:
		return parseEmailMarkdown(data)
	default:
		return parseEmailText(data), nil
	}
}

// emailFileRaw is the YAML/JSON email-file shape. to/cc/bcc accept a
// comma-separated string or a list of strings.
type emailFileRaw struct {
	To      interface{} `yaml:"to" json:"to"`
	Cc      interface{} `yaml:"cc" json:"cc"`
	Bcc     interface{} `yaml:"bcc" json:"bcc"`
	Subject string      `yaml:"subject" json:"subject"`
	Body    string      `yaml:"body" json:"body"`
	IsHTML  bool        `yaml:"html" json:"html"`
}

// parseEmailYAML parses a YAML email file.
func parseEmailYAML(data []byte, path string) (EmailMessage, error) {
	var raw emailFileRaw
	if yamlErr := yaml.Unmarshal(data, &raw); yamlErr != nil {
		return EmailMessage{}, apperror.WrapSimple(yamlErr, "parse email yaml "+path+":")
	}

	return emailMessageFromRaw(raw), nil
}

// parseEmailJSON parses a JSON email file.
func parseEmailJSON(data []byte, path string) (EmailMessage, error) {
	var raw emailFileRaw
	if jsonErr := json.Unmarshal(data, &raw); jsonErr != nil {
		return EmailMessage{}, apperror.WrapSimple(jsonErr, "parse email json "+path+":")
	}

	return emailMessageFromRaw(raw), nil
}

// emailMessageFromRaw converts the parsed raw shape into an EmailMessage.
func emailMessageFromRaw(raw emailFileRaw) EmailMessage {
	return EmailMessage{
		To:      normalizeAddresses(raw.To),
		Cc:      normalizeAddresses(raw.Cc),
		Bcc:     normalizeAddresses(raw.Bcc),
		Subject: strings.TrimSpace(raw.Subject),
		Body:    raw.Body,
		IsHTML:  raw.IsHTML,
	}
}

// normalizeAddresses accepts a comma-separated string or a string list.
func normalizeAddresses(val interface{}) []string {
	switch typed := val.(type) {
	case string:
		return splitAddresses(typed)
	case []string:
		return typed
	case []interface{}:
		return collectAddressList(typed)
	default:
		return nil
	}
}

// collectAddressList trims every string item of a YAML/JSON address list.
func collectAddressList(items []interface{}) []string {
	var out []string
	for _, item := range items {
		addr, isStr := item.(string)
		if isStr && len(strings.TrimSpace(addr)) > 0 {
			out = append(out, strings.TrimSpace(addr))
		}
	}

	return out
}

// splitAddresses splits a comma-separated address list.
func splitAddresses(list string) []string {
	var out []string
	for _, part := range strings.Split(list, ",") {
		trimmed := strings.TrimSpace(part)
		if len(trimmed) > 0 {
			out = append(out, trimmed)
		}
	}

	return out
}

// parseEmailMarkdown parses a Markdown email file: the first `# `
// heading is the subject and everything after it is the body.
func parseEmailMarkdown(data []byte) (EmailMessage, error) {
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return EmailMessage{
				Subject: strings.TrimSpace(strings.TrimPrefix(line, "# ")),
				Body:    strings.TrimLeft(strings.Join(lines[i+1:], "\n"), "\n"),
			}, nil
		}
	}

	return EmailMessage{}, apperror.NewSimple("email markdown file has no '# ' heading for the subject", "E_SEND_EMAIL_NO_SUBJECT")
}

// parseEmailText parses a plain-text email file: the first non-empty
// line is the subject and the rest is the body (blank lines preserved).
func parseEmailText(data []byte) EmailMessage {
	lines := strings.Split(string(data), "\n")
	subjectIdx := firstNonEmptyLine(lines)

	subject := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[subjectIdx]), "Subject:"))

	body := ""
	if subjectIdx+1 < len(lines) {
		body = strings.Join(lines[subjectIdx+1:], "\n")
	}

	return EmailMessage{Subject: subject, Body: body}
}

// firstNonEmptyLine returns the index of the first non-blank line.
func firstNonEmptyLine(lines []string) int {
	for i, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			return i
		}
	}

	return 0
}
