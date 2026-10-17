package cmdspace

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
)

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
				decoded, decErr := base64.StdEncoding.DecodeString(fields[2])
				recordMockAuthPayload(m, decoded, decErr)
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

func recordMockAuthPayload(m *mockSMTPServer, decoded []byte, decErr error) {
	if decErr != nil {
		return
	}
	m.mu.Lock()
	m.authPayload = string(decoded)
	m.mu.Unlock()
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
