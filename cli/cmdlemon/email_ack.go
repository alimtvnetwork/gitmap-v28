package cmdlemon

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AckReplier handles immediate receipt acknowledgment dispatch for commands.
type AckReplier struct {
	transmitter EmailTransmitter
}

// NewAckReplier creates a new receipt acknowledgment replier.
func NewAckReplier(transmitter EmailTransmitter) *AckReplier {
	return &AckReplier{transmitter: transmitter}
}

// DispatchAck transmits the immediate ACK confirmation back to the sender.
func (r *AckReplier) DispatchAck(env *CommandEnvelope) error {
	if r.transmitter == nil {
		return apperror.NewValidationError("nil email transmitter")
	}

	if env == nil {
		return apperror.NewValidationError("nil command envelope")
	}

	subject := FormatAckSubject(env)
	body := FormatAckBody(env)

	if err := r.transmitter.Transmit(env.Sender, subject, body); err != nil {
		return apperror.WrapSimple(err, "failed to transmit acknowledgment")
	}

	return nil
}

// FormatAckSubject builds the email subject line for receipt acknowledgment.
func FormatAckSubject(env *CommandEnvelope) string {
	return fmt.Sprintf("ACK Received: %s (Request ID: %s)", env.Command, env.RequestID)
}

// FormatAckBody renders the mandated receipt acknowledgment email body text.
func FormatAckBody(env *CommandEnvelope) string {
	ts := env.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	cmdDisplay := env.Command
	hasArgs := env.HasArgs()
	if hasArgs {
		cmdDisplay = fmt.Sprintf("%s %s", env.Command, strings.Join(env.Args, " "))
	}

	return fmt.Sprintf("ACK Received: Command '%s' received at %s (Request ID: %s). Executing...",
		cmdDisplay, ts.Format(time.RFC3339), env.RequestID)
}

// SendImmediateAck dispatches receipt confirmation directly using the provided transmitter.
func SendImmediateAck(transmitter EmailTransmitter, env *CommandEnvelope) error {
	replier := NewAckReplier(transmitter)

	return replier.DispatchAck(env)
}
