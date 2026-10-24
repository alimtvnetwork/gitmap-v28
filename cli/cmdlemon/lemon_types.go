package cmdlemon

import (
	"fmt"
	"strings"
	"time"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdlemon/enums"

// EmailFormatType represents the encoding of an incoming email command.
// EmailFormatType aliases the strongly-typed enum from the enums package.
type EmailFormatType = enums.EmailFormatType

const (
	EmailFormatTypeText    = enums.EmailFormatTypeText
	EmailFormatTypeJSON    = enums.EmailFormatTypeJSON
	EmailFormatTypeYAML    = enums.EmailFormatTypeYAML
	EmailFormatTypeUnknown = enums.EmailFormatTypeUnknown
)

// EmailMessage models an incoming or outgoing raw email message.
type EmailMessage struct {
	From       string    `json:"from"`
	To         string    `json:"to"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	ReceivedAt time.Time `json:"receivedAt"`
}

// CommandEnvelope represents parsed command instructions and metadata.
type CommandEnvelope struct {
	Command   string          `json:"command"`
	Args      []string        `json:"args"`
	RequestID string          `json:"requestId"`
	Sender    string          `json:"sender"`
	Format    EmailFormatType `json:"format"`
	Timestamp time.Time       `json:"timestamp"`
	RawBody   string          `json:"rawBody,omitempty"`
}

// HasArgs returns true if the command envelope contains one or more arguments.
func (env *CommandEnvelope) HasArgs() bool {
	if env == nil {
		return false
	}

	return len(env.Args) > 0
}

// EmailTransmitter dispatches emails in-memory or over SMTP.
type EmailTransmitter interface {
	Transmit(to, subject, body string) error
}

// EmailTransmitterFunc adapts a plain function to EmailTransmitter.
type EmailTransmitterFunc func(to, subject, body string) error

// Transmit calls the underlying function adapter.
func (f EmailTransmitterFunc) Transmit(to, subject, body string) error {
	return f(to, subject, body)
}

// ExecutionOutcome encapsulates the completion state of a command execution.
type ExecutionOutcome struct {
	Command     string        `json:"command"`
	RequestID   string        `json:"requestId"`
	ExitCode    int           `json:"exitCode"`
	Duration    time.Duration `json:"duration"`
	Stdout      string        `json:"stdout"`
	Stderr      string        `json:"stderr"`
	CompletedAt time.Time     `json:"completedAt"`
	IsSuccess   bool          `json:"isSuccess"`
}

// HasSuccess returns whether the execution outcome succeeded.
func (o ExecutionOutcome) HasSuccess() bool {
	return o.IsSuccess
}

// PopulateDefaults fills missing envelope fields with defaults.
func PopulateDefaults(env *CommandEnvelope, msg EmailMessage, fmtType EmailFormatType) {
	if env == nil {
		return
	}

	env.Format = fmtType
	env.RawBody = msg.Body
	if len(env.Sender) == 0 {
		env.Sender = msg.From
	}

	if len(env.RequestID) == 0 {
		env.RequestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
	}

	if env.Timestamp.IsZero() {
		env.Timestamp = time.Now().UTC()
	}
}

// ExtractFencedBlock extracts fenced code content with the specified markdown tag.
func ExtractFencedBlock(body, tag string) string {
	start := strings.Index(body, tag)
	hasTag := start >= 0
	if !hasTag {
		return body
	}

	rest := body[start+len(tag):]
	end := strings.Index(rest, "```")
	hasClosingFence := end >= 0
	if hasClosingFence {
		return strings.TrimSpace(rest[:end])
	}

	return strings.TrimSpace(rest)
}
