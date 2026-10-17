package cmdspace

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// writeTempEmailFile writes content to dir/name and returns the path.
func writeTempEmailFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	return path
}

func TestParseSendEmailArgsTable(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    sendEmailOptions
		wantErr string
	}{
		{
			name: "space separated",
			args: []string{"--to", "a@x.com", "--subject", "Hi", "--body", "Hello"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi", Body: "Hello"},
		},
		{
			name: "equals form",
			args: []string{"--to=a@x.com", "--subject=Hi"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi"},
		},
		{
			name: "html bool",
			args: []string{"--to", "a@x.com", "--subject", "Hi", "--html"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi", IsHTML: true},
		},
		{
			name:    "html rejects value",
			args:    []string{"--to", "a@x.com", "--html=x"},
			wantErr: "takes no value",
		},
		{
			name:    "unknown flag first",
			args:    []string{"--bogus", "--to", "a@x.com"},
			wantErr: "unknown flag: --bogus",
		},
		{
			name:    "trailing unknown flag rejected before value resolution",
			args:    []string{"--to", "a@x.com", "--subject", "Hi", "--bogus"},
			wantErr: "unknown flag: --bogus",
		},
		{
			name:    "equals unknown flag",
			args:    []string{"--to", "a@x.com", "--bogus=x"},
			wantErr: "unknown flag: --bogus",
		},
		{
			name:    "positional rejected",
			args:    []string{"--to", "a@x.com", "oops"},
			wantErr: "unexpected positional argument",
		},
		{
			name:    "missing value",
			args:    []string{"--to"},
			wantErr: "requires a value",
		},
		{
			name: "file flag",
			args: []string{"--file", "note.yaml", "--to", "b@x.com"},
			want: sendEmailOptions{File: "note.yaml", To: "b@x.com"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := parseSendEmailArgs(tc.args)
			if len(tc.wantErr) > 0 {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
				}

				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts != tc.want {
				t.Fatalf("got %+v, want %+v", opts, tc.want)
			}
		})
	}
}

func TestRunSendEmailMissingTo(t *testing.T) {
	if err := RunSendEmail([]string{"--subject", "Hi"}); err == nil ||
		!strings.Contains(err.Error(), "missing --to") {
		t.Fatalf("expected missing --to error, got %v", err)
	}
}

func TestRunSendEmailMissingSubject(t *testing.T) {
	if err := RunSendEmail([]string{"--to", "a@x.com"}); err == nil ||
		!strings.Contains(err.Error(), "missing subject") {
		t.Fatalf("expected missing subject error, got %v", err)
	}
}

func TestRunSendEmailUnknownFlag(t *testing.T) {
	err := RunSendEmail([]string{"--to", "a@x.com", "--subject", "Hi", "--nope"})
	if err == nil || !strings.Contains(err.Error(), "unknown flag: --nope") {
		t.Fatalf("expected unknown flag error, got %v", err)
	}
}

func TestParseEmailFileYAML(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.yaml",
		"to: a@x.com, b@x.com\ncc:\n  - c@x.com\nsubject:  Hello \nbody: YAML body\nhtml: true\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 2 || msg.To[0] != "a@x.com" || msg.To[1] != "b@x.com" {
		t.Fatalf("to wrong: %v", msg.To)
	}
	if len(msg.Cc) != 1 || msg.Cc[0] != "c@x.com" {
		t.Fatalf("cc wrong: %v", msg.Cc)
	}
	if msg.Subject != "Hello" || msg.Body != "YAML body" || !msg.IsHTML {
		t.Fatalf("content wrong: %+v", msg)
	}
}

func TestParseEmailFileYAMLListTo(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.yml", "to:\n  - a@x.com\n  - b@x.com\nsubject: S\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 2 {
		t.Fatalf("to wrong: %v", msg.To)
	}
}

func TestParseEmailFileJSON(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.json",
		`{"to": ["a@x.com"], "bcc": "d@x.com", "subject": "J", "body": "json body"}`)

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 1 || len(msg.Bcc) != 1 || msg.Subject != "J" || msg.Body != "json body" {
		t.Fatalf("wrong: %+v", msg)
	}
}

func TestParseEmailFileMarkdown(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.md", "# Meeting notes\n\nFirst line.\nSecond line.\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Subject != "Meeting notes" {
		t.Fatalf("subject wrong: %q", msg.Subject)
	}
	if msg.Body != "First line.\nSecond line.\n" {
		t.Fatalf("body wrong: %q", msg.Body)
	}
}

func TestParseEmailFileMarkdownNoHeading(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.md", "no heading here\n")

	if _, err := parseEmailFile(path); err == nil ||
		!strings.Contains(err.Error(), "E_SEND_EMAIL_NO_SUBJECT") {
		t.Fatalf("expected no-subject error, got %v", err)
	}
}

func TestParseEmailFileText(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.txt", "\n\nSubject:  Greetings \n\nBody keeps\n\nblank lines.\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Subject != "Greetings" {
		t.Fatalf("subject wrong: %q", msg.Subject)
	}
	if msg.Body != "\nBody keeps\n\nblank lines.\n" {
		t.Fatalf("body wrong: %q", msg.Body)
	}
}

func TestParseEmailFileErrors(t *testing.T) {
	dir := t.TempDir()

	if _, err := parseEmailFile(filepath.Join(dir, "missing.yaml")); err == nil {
		t.Fatal("expected error for missing file")
	}

	empty := writeTempEmailFile(t, dir, "empty.txt", "   \n  \n")
	if _, err := parseEmailFile(empty); err == nil ||
		!strings.Contains(err.Error(), "E_SEND_EMAIL_EMPTY_FILE") {
		t.Fatalf("expected empty-file error, got %v", err)
	}

	badYAML := writeTempEmailFile(t, dir, "bad.yaml", "to: [unclosed\n")
	if _, err := parseEmailFile(badYAML); err == nil {
		t.Fatal("expected error for bad YAML")
	}

	badJSON := writeTempEmailFile(t, dir, "bad.json", `{"to": `)
	if _, err := parseEmailFile(badJSON); err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestBuildEmailMessageOverrides(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.yaml",
		"to: file@x.com\nsubject: FileSubj\nbody: FileBody\n")

	msg, err := buildEmailMessage(sendEmailOptions{File: path, Subject: "CliSubj", To: "cli@x.com"})
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 1 || msg.To[0] != "cli@x.com" {
		t.Fatalf("to override wrong: %v", msg.To)
	}
	if msg.Subject != "CliSubj" {
		t.Fatalf("subject override wrong: %q", msg.Subject)
	}
	if msg.Body != "FileBody" {
		t.Fatalf("body should come from file: %q", msg.Body)
	}
}

func TestBuildEmailMessageFileMissingSubject(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.json", `{"to": "a@x.com", "body": "no subject key"}`)

	if _, err := buildEmailMessage(sendEmailOptions{File: path, To: "a@x.com"}); err == nil {
		t.Fatal("expected error when neither file nor flag provides a subject")
	}
}

func TestResolveGoEmailReaderBinEmpty(t *testing.T) {
	empty := t.TempDir()
	t.Setenv("PATH", empty)
	t.Setenv("HOME", empty)

	if got := resolveGoEmailReaderBin(); len(got) > 0 {
		t.Fatalf("expected empty resolution, got %q", got)
	}
}

// stubGoEmailReader installs an executable stub that records its argv.
func stubGoEmailReader(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if mkErr := os.Mkdir(binDir, 0755); mkErr != nil {
		t.Fatal(mkErr)
	}
	out := filepath.Join(dir, "argv.txt")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"" + out + "\"\n"
	stub := filepath.Join(binDir, "go-email-reader")
	if writeErr := os.WriteFile(stub, []byte(script), 0755); writeErr != nil {
		t.Fatal(writeErr)
	}
	t.Setenv("PATH", binDir)

	return out
}

func TestDelegateSendEmailArgv(t *testing.T) {
	out := stubGoEmailReader(t)

	args := []string{"--to", "a@x.com", "--subject", "Hi", "--body", "Hello", "--html"}
	if err := RunSendEmail(args); err != nil {
		t.Fatal(err)
	}

	raw, readErr := os.ReadFile(out)
	if readErr != nil {
		t.Fatal(readErr)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	want := []string{"send", "--to", "a@x.com", "--subject", "Hi", "--body", "Hello", "--html"}
	if len(lines) != len(want) {
		t.Fatalf("argv wrong: %v", lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestDelegateSendEmailArgvFile(t *testing.T) {
	out := stubGoEmailReader(t)

	dir := t.TempDir()
	yml := writeTempEmailFile(t, dir, "m.yaml", "to: file@x.com\nsubject: FileSubj\n")
	if err := RunSendEmail([]string{"--file", yml}); err != nil {
		t.Fatal(err)
	}

	raw, readErr := os.ReadFile(out)
	if readErr != nil {
		t.Fatal(readErr)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	want := []string{"send", "--to", "file@x.com", "--file", yml}
	if len(lines) != len(want) {
		t.Fatalf("argv wrong: %v", lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("argv[%d] = %q, want %q", i, lines[i], want[i])
		}
	}
}

func TestBuildSMTPMessageHeaders(t *testing.T) {
	msg := EmailMessage{
		To:      []string{"a@x.com"},
		Cc:      []string{"b@x.com"},
		Bcc:     []string{"secret@x.com"},
		Subject: "Subj",
		Body:    "Body text",
	}

	plain := buildSMTPMessage("sender@example.com", msg)
	for _, want := range []string{"Subject: Subj", "To: a@x.com", "Cc: b@x.com", "Body text", "text/plain"} {
		if !strings.Contains(plain, want) {
			t.Fatalf("missing %q in:\n%s", want, plain)
		}
	}
	if strings.Contains(plain, "secret@x.com") || strings.Contains(plain, "Bcc:") {
		t.Fatalf("Bcc leaked into headers:\n%s", plain)
	}

	html := buildSMTPMessage("sender@example.com", EmailMessage{To: []string{"a@x.com"}, Subject: "S", Body: "<b>x</b>", IsHTML: true})
	if !strings.Contains(html, "text/html") {
		t.Fatalf("expected text/html content type:\n%s", html)
	}
}

// mockSMTPServer is a minimal AUTH-advertising SMTP server for tests.
type mockSMTPServer struct {
	ln          net.Listener
	done        chan struct{}
	mu          sync.Mutex
	authPayload string
	mailFrom    string
	rcpts       []string
	data        string
}

func startMockSMTPServer(t *testing.T) *mockSMTPServer {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	m := &mockSMTPServer{ln: ln, done: make(chan struct{})}
	go m.serve()

	return m
}

func (m *mockSMTPServer) addr() string {
	return m.ln.Addr().String()
}

func (m *mockSMTPServer) serve() {
	defer close(m.done)

	conn, err := m.ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	write := func(s string) { fmt.Fprintf(conn, "%s\r\n", s) }
	readLine := func() string {
		line, _ := reader.ReadString('\n')

		return strings.TrimRight(line, "\r\n")
	}

	write("220 mock ESMTP ready")
	for {
		line := readLine()
		upper := strings.ToUpper(line)

		switch {
		case strings.HasPrefix(upper, "EHLO"):
			write("250-mock")
			write("250 AUTH PLAIN")
		case strings.HasPrefix(upper, "AUTH PLAIN"):
			fields := strings.SplitN(line, " ", 3)
			if len(fields) == 3 {
				if decoded, decErr := base64.StdEncoding.DecodeString(fields[2]); decErr == nil {
					m.mu.Lock()
					m.authPayload = string(decoded)
					m.mu.Unlock()
				}
			}
			write("235 2.7.0 ok")
		case strings.HasPrefix(upper, "MAIL FROM:"):
			m.mu.Lock()
			m.mailFrom = strings.Trim(line[len("MAIL FROM:"):], "<> ")
			m.mu.Unlock()
			write("250 ok")
		case strings.HasPrefix(upper, "RCPT TO:"):
			m.mu.Lock()
			m.rcpts = append(m.rcpts, strings.Trim(line[len("RCPT TO:"):], "<> "))
			m.mu.Unlock()
			write("250 ok")
		case strings.HasPrefix(upper, "DATA"):
			write("354 go ahead")
			var sb strings.Builder
			for {
				bodyLine := readLine()
				if bodyLine == "." {
					break
				}
				sb.WriteString(bodyLine + "\r\n")
			}
			m.mu.Lock()
			m.data = sb.String()
			m.mu.Unlock()
			write("250 queued")
		case strings.HasPrefix(upper, "QUIT"):
			write("221 bye")

			return
		case strings.HasPrefix(upper, "RSET"):
			write("250 ok")
		default:
			write("502 unimplemented")
		}
	}
}

func (m *mockSMTPServer) snapshot() (auth, from string, rcpts []string, data string) {
	<-m.done
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.authPayload, m.mailFrom, append([]string{}, m.rcpts...), m.data
}

func TestDeliverViaSMTPMock(t *testing.T) {
	server := startMockSMTPServer(t)

	host, portStr, splitErr := net.SplitHostPort(server.addr())
	if splitErr != nil {
		t.Fatal(splitErr)
	}
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	msg := EmailMessage{
		To:      []string{"a@x.com"},
		Cc:      []string{"b@x.com"},
		Bcc:     []string{"c@x.com", "d@x.com"},
		Subject: "Test subj",
		Body:    "Hello body",
	}
	if err := deliverViaSMTP(host, port, "sender@example.com", "s3cr3t-pw", msg); err != nil {
		t.Fatal(err)
	}

	auth, from, rcpts, data := server.snapshot()
	if auth != "\x00sender@example.com\x00s3cr3t-pw" {
		t.Fatalf("AUTH PLAIN payload wrong: %q", auth)
	}
	if from != "sender@example.com" {
		t.Fatalf("MAIL FROM wrong: %q", from)
	}
	wantRcpts := []string{"a@x.com", "b@x.com", "c@x.com", "d@x.com"}
	if len(rcpts) != len(wantRcpts) {
		t.Fatalf("RCPT count wrong: %v", rcpts)
	}
	for i := range wantRcpts {
		if rcpts[i] != wantRcpts[i] {
			t.Fatalf("RCPT %d wrong: %v", i, rcpts)
		}
	}
	for _, want := range []string{"Subject: Test subj", "Hello body"} {
		if !strings.Contains(data, want) {
			t.Fatalf("DATA missing %q:\n%s", want, data)
		}
	}
	if strings.Contains(data, "d@x.com") {
		t.Fatalf("Bcc leaked into DATA:\n%s", data)
	}
}

func TestSendEmailViaSMTPPropagatesPromptError(t *testing.T) {
	old := promptEmailPassword
	promptEmailPassword = func() (string, error) {
		return "", errors.New("no tty")
	}
	defer func() { promptEmailPassword = old }()

	msg := EmailMessage{To: []string{"a@x.com"}, Subject: "S", Body: "B"}
	if err := sendEmailViaSMTP(msg); err == nil ||
		!strings.Contains(err.Error(), "no tty") {
		t.Fatalf("expected prompt error, got %v", err)
	}
}
