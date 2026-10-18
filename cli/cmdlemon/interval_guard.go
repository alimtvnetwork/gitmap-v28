package cmdlemon

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	DefaultPollingInterval = 8 * time.Second
	MinPollingInterval     = 5 * time.Second
	NoticeClampedMessage   = "Notice: interval clamped to minimum 5s safety threshold\n"
)

// IntervalGuard validates and clamps polling intervals to ensure safety.
type IntervalGuard struct {
	out io.Writer
}

// NewIntervalGuard creates an IntervalGuard with the specified writer.
func NewIntervalGuard(out io.Writer) *IntervalGuard {
	targetOut := out
	if targetOut == nil {
		targetOut = os.Stderr
	}

	return &IntervalGuard{out: targetOut}
}

// ResolveInterval enforces the minimum 5s safety threshold on the interval.
func (g *IntervalGuard) ResolveInterval(requested time.Duration) time.Duration {
	if requested <= 0 {
		return DefaultPollingInterval
	}

	if requested < MinPollingInterval {
		fmt.Fprint(g.out, NoticeClampedMessage)

		return MinPollingInterval
	}

	return requested
}

// ParseAndResolveInterval parses a duration string and applies safety clamping.
func (g *IntervalGuard) ParseAndResolveInterval(raw string) (time.Duration, error) {
	if raw == "" {
		return DefaultPollingInterval, nil
	}

	dur, err := time.ParseDuration(raw)
	if err != nil {
		return DefaultPollingInterval, err
	}

	return g.ResolveInterval(dur), nil
}
