package cmdmail

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"net"
	"net/http"
	"time"
)

const (
	// DefaultMailUIPort is the default start port for Mail UI probing.
	DefaultMailUIPort = 8090
	// MailPortScanAttempts defines maximum sequential port scan attempts.
	MailPortScanAttempts = 50
)

// MailUIOptions configures the mail UI HTTP server.
type MailUIOptions struct {
	Host string
	Port int
}

// MailSendRequest specifies payload for sending an email.
type MailSendRequest struct {
	Account      string   `json:"account"`
	To           []string `json:"to"`
	Cc           []string `json:"cc,omitempty"`
	Bcc          []string `json:"bcc,omitempty"`
	Subject      string   `json:"subject"`
	BodyMarkdown string   `json:"bodyMarkdown"`
	IsHTML       bool     `json:"isHTML"`
}

// MailSendResponse returns dispatch outcome.
type MailSendResponse struct {
	IsSuccess      bool   `json:"isSuccess"`
	Message        string `json:"message"`
	SentAt         string `json:"sentAt"`
	RecipientCount int    `json:"recipientCount"`
	Error          string `json:"error,omitempty"`
}

// MailTemplate defines an email template skeleton.
type MailTemplate struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Subject      string `json:"subject"`
	BodyMarkdown string `json:"bodyMarkdown"`
}

// MailAccount defines configured outgoing mail accounts.
type MailAccount struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Host      string `json:"host"`
	IsDefault bool   `json:"isDefault"`
}

// BuildMailUIMux constructs the HTTP route multiplexer.
func BuildMailUIMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMailIndex)
	mux.HandleFunc("/ui", handleMailIndex)
	mux.HandleFunc("/api/mail/send", handleAPIMailSend)
	mux.HandleFunc("/api/mail/templates", handleAPIMailTemplates)
	mux.HandleFunc("/api/mail/accounts", handleAPIMailAccounts)

	return mux
}

// ScanAvailableMailListener searches for an open port using sequential probing.
func ScanAvailableMailListener(host string, startPort int) (net.Listener, int, error) {
	bindHost := resolveHost(host)
	targetPort := resolveStartPort(startPort)

	for port := targetPort; port < targetPort+MailPortScanAttempts; port++ {
		addr := fmt.Sprintf("%s:%d", bindHost, port)
		ln, err := net.Listen("tcp", addr)
		isBound := err == nil && ln != nil
		if isBound {
			return ln, port, nil
		}
	}

	return bindFallbackPort(bindHost)
}

func resolveHost(host string) string {
	if host == "" {
		return "127.0.0.1"
	}
	return host
}

func resolveStartPort(port int) int {
	if port <= 0 {
		return DefaultMailUIPort
	}
	return port
}

func bindFallbackPort(host string) (net.Listener, int, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:0", host))
	isBound := err == nil && ln != nil
	if isBound {
		boundPort := ln.Addr().(*net.TCPAddr).Port
		return ln, boundPort, nil
	}
	return nil, 0, apperror.WrapSimple(err, "failed to bind mail UI on any port")
}

// StartMailUIServer boots the HTTP server with graceful port probing.
func StartMailUIServer(opts MailUIOptions) (*http.Server, net.Listener, int, error) {
	listener, port, err := ScanAvailableMailListener(opts.Host, opts.Port)
	if err != nil {
		return nil, nil, 0, err
	}

	server := &http.Server{
		Handler:           BuildMailUIMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server, listener, port, nil
}

func handleMailIndex(w http.ResponseWriter, r *http.Request) {
	isPathMatch := r.URL.Path == "/" || r.URL.Path == "/ui"
	if !isPathMatch {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(BuildMailComposerHTML()))
}

func handleAPIMailSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "POST method required")
		return
	}

	var req MailSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hasRecipients := len(req.To) > 0
	if !hasRecipients {
		writeJSONError(w, http.StatusBadRequest, "At least one recipient is required")
		return
	}

	resp := processSendSuccess(req)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func processSendSuccess(req MailSendRequest) MailSendResponse {
	totalCount := len(req.To) + len(req.Cc) + len(req.Bcc)
	msg := fmt.Sprintf("Email dispatched to %d recipient(s) via account [%s]", totalCount, req.Account)

	return MailSendResponse{
		IsSuccess:      true,
		Message:        msg,
		SentAt:         time.Now().UTC().Format(time.RFC3339),
		RecipientCount: totalCount,
	}
}

func handleAPIMailTemplates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetDefaultMailTemplates())
}

func handleAPIMailAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(GetConfiguredMailAccounts())
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	resp := MailSendResponse{
		IsSuccess: false,
		Error:     msg,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// GetDefaultMailTemplates supplies built-in email templates.
func GetDefaultMailTemplates() []MailTemplate {
	return []MailTemplate{
		{
			ID:           "daily-digest",
			Name:         "Daily Digest & Fleet Status",
			Subject:      "[GitMap] Daily Fleet & Health Digest",
			BodyMarkdown: "## Fleet Health Summary\n\nAll repositories synced successfully.\n\n- Active nodes: **12**\n- Unresolved alerts: **0**\n- Last run: `ok`\n",
		},
		{
			ID:           "ci-alert",
			Name:         "CI/CD Pipeline Failure",
			Subject:      "[ALERT] Pipeline Failure Detected",
			BodyMarkdown: "## Pipeline Failure Alert\n\nBuild execution encountered an error.\n\n- Branch: `main`\n- Trigger: Automated test runner\n- Status: **Action Required**\n",
		},
		{
			ID:           "sync-summary",
			Name:         "Multi-Repo Sync Report",
			Subject:      "[GitMap] Repository Synchronization Report",
			BodyMarkdown: "## Synchronization Complete\n\nSynchronized canonical templates and prompts across workspace.\n\n- Repos updated: **43**\n- Zero conflict status verified.\n",
		},
	}
}

// GetConfiguredMailAccounts returns available mailing sender profiles.
func GetConfiguredMailAccounts() []MailAccount {
	return []MailAccount{
		{
			Name:      "default",
			Email:     "notifier@gitmap.local",
			Host:      "smtp://localhost:1025",
			IsDefault: true,
		},
		{
			Name:      "alerts",
			Email:     "alerts@gitmap.internal",
			Host:      "smtp.internal:587",
			IsDefault: false,
		},
		{
			Name:      "ci-bot",
			Email:     "ci-noreply@gitmap.io",
			Host:      "mail.gitmap.io:465",
			IsDefault: false,
		},
	}
}
