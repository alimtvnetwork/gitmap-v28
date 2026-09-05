package applogger_test

import (
	"bytes"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
	"coding-guidelines/common/pkg/applogger"
	"coding-guidelines/common/pkg/errtype"
)

func newTestConsoleLogger(buf *bytes.Buffer, useJSON bool) applogger.Logger {
	sink := applogger.NewConsoleSink(buf, useJSON)
	cfg := applogger.Config{
		MinLevel: applogger.LevelDebug,
		Driver:   applogger.DriverComposite,
		Sinks:    []applogger.LogSink{sink},
	}

	l, _ := applogger.New(cfg)

	return l
}

func TestConsoleSinkLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	l := newTestConsoleLogger(buf, false)
	l.Infof("hello %s", "world")

	if !strings.Contains(buf.String(), "hello world") {
		t.Fatalf("expected log output, got %s", buf.String())
	}
}

func TestLogErrorAndLogFaults(t *testing.T) {
	buf := &bytes.Buffer{}
	l := newTestConsoleLogger(buf, true)
	err := appfault.New(errtype.Database, "connection reset").WithOp("db.connect")

	l.LogError(err)
	l.LogFaults(appfaults.New().Add(err))

	if !strings.Contains(buf.String(), "connection reset") {
		t.Fatalf("expected error message in log output, got %s", buf.String())
	}
}

func TestLoggerChaining(t *testing.T) {
	buf := &bytes.Buffer{}
	l := newTestConsoleLogger(buf, false)
	err := appfault.New(errtype.Database, "connection reset")

	ret := l.Debug("chain debug").
		Info("chain info").
		Warn("chain warn").
		Error("chain error").
		Debugf("chain debugf %d", 1).
		Infof("chain infof %d", 2).
		Warnf("chain warnf %d", 3).
		Errorf("chain errorf %d", 4).
		WithContext("key", "val").
		WithFields(map[string]any{"env": "test"}).
		LogError(err).
		LogFaults(appfaults.New().Add(err))

	if ret == nil {
		t.Fatalf("expected non-nil logger from chain")
	}

	out := buf.String()
	if !strings.Contains(out, "chain info") || !strings.Contains(out, "chain errorf 4") {
		t.Fatalf("expected chained log output, got %s", out)
	}
}
