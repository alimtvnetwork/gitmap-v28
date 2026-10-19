package cmdlemon

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const StateFilePath = ".gitmap/lemon_state.json"

type (
	DaemonState struct {
		LastPollTime        time.Time `json:"lastPollTime"`
		LastExecutedCommand string    `json:"lastExecutedCommand"`
		PollCount           int64     `json:"pollCount"`
		Pid                 int       `json:"pid"`
		IsRunning           bool      `json:"isRunning"`
	}
	DaemonConfig struct {
		PollFunc          func(ctx context.Context) error
		Interval          time.Duration
		IsOnce, IsVerbose bool
	}
	LemonDaemon struct {
		cfg   DaemonConfig
		guard *IntervalGuard
		mu    sync.RWMutex
		state DaemonState
	}
	GitmapDaemon = LemonDaemon
)

func NewLemonDaemon(cfg DaemonConfig) *LemonDaemon {
	guard := NewIntervalGuard(os.Stderr)
	cfg.Interval = guard.ResolveInterval(cfg.Interval)
	return &LemonDaemon{cfg: cfg, guard: guard, state: DaemonState{Pid: os.Getpid()}}
}

func NewGitmapDaemon(cfg DaemonConfig) *GitmapDaemon {
	return NewLemonDaemon(cfg)
}

func (d *LemonDaemon) GetState() DaemonState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

func (d *LemonDaemon) Run(ctx context.Context) error {
	if !IsDaemonEnabled() && !d.cfg.IsOnce {
		return nil
	}
	d.setRunning(true)
	defer d.setRunning(false)
	if d.cfg.IsOnce {
		return d.executeSinglePass(ctx)
	}
	return d.runLoop(ctx)
}

func (d *LemonDaemon) setRunning(running bool) {
	d.mu.Lock()
	d.state.IsRunning = running
	data, _ := json.Marshal(d.state)
	d.mu.Unlock()
	_ = os.MkdirAll(filepath.Dir(StateFilePath), 0o755)
	_ = os.WriteFile(StateFilePath, data, 0o644)
}

func (d *LemonDaemon) runLoop(ctx context.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sig)
	t := time.NewTicker(d.cfg.Interval)
	defer t.Stop()
	for {
		if err := d.executeSinglePass(ctx); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-sig:
			return nil
		case <-t.C:
		}
	}
}

func (d *LemonDaemon) executeSinglePass(ctx context.Context) error {
	if !IsDaemonEnabled() {
		return nil
	}
	d.mu.Lock()
	d.state.LastPollTime, d.state.PollCount = time.Now(), d.state.PollCount+1
	data, _ := json.Marshal(d.state)
	d.mu.Unlock()
	_ = os.MkdirAll(filepath.Dir(StateFilePath), 0o755)
	_ = os.WriteFile(StateFilePath, data, 0o644)
	if d.cfg.PollFunc != nil {
		return d.cfg.PollFunc(ctx)
	}
	return nil
}
