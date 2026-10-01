package cmdpull

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// StartPullHeartbeat runs a background ticker emitting progress notices once threshold elapses.
func StartPullHeartbeat(threshold time.Duration, interval time.Duration, progressFn func() string, isJSON bool, out ...io.Writer) func() {
	if isJSON {
		return func() {}
	}
	ticker, stopChan := initHeartbeatTicker(interval)
	var once sync.Once
	w := resolveHeartbeatWriter(out...)
	go runHeartbeatLoop(ticker, stopChan, time.Now(), threshold, progressFn, w)

	return func() {
		once.Do(func() {
			ticker.Stop()
			close(stopChan)
		})
	}
}

func initHeartbeatTicker(interval time.Duration) (*time.Ticker, chan struct{}) {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	return time.NewTicker(interval), make(chan struct{})
}

func resolveHeartbeatWriter(out ...io.Writer) io.Writer {
	if len(out) > 0 && out[0] != nil {
		return out[0]
	}
	return os.Stdout
}

func runHeartbeatLoop(ticker *time.Ticker, stopChan <-chan struct{}, startTime time.Time, threshold time.Duration, progressFn func() string, w io.Writer) {
	for {
		select {
		case <-stopChan:
			return
		case tickTime := <-ticker.C:
			emitHeartbeat(tickTime.Sub(startTime), threshold, progressFn, w)
		}
	}
}

func emitHeartbeat(elapsed, threshold time.Duration, progressFn func() string, w io.Writer) {
	if elapsed < threshold {
		return
	}
	msg := resolveHeartbeatMessage(progressFn)
	elapsedSec := int(elapsed.Seconds())
	fmt.Fprintf(w, "  [⏳ %ds elapsed] %s\n", elapsedSec, msg)
}

func resolveHeartbeatMessage(progressFn func() string) string {
	if progressFn != nil {
		if msg := progressFn(); msg != "" {
			return msg
		}
	}
	return "Pulling repositories..."
}
