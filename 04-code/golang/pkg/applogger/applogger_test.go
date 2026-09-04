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
