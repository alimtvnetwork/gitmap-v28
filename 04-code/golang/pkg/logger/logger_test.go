package logger_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/logger"
)

func TestLoggerConsoleOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := logger.DefaultOptions().WithOutput(buf).WithJson(false)
	log := logger.New(opts)

	log.Info("test info message")
	if !strings.Contains(buf.String(), "test info message") {
		t.Fatalf("expected log output to contain message, got %s", buf.String())
	}
}

func TestLoggerAppErrorLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New(logger.DefaultOptions().WithOutput(buf).WithJson(true))
	appErr := appfault.New(errtype.NotFound, "record missing").WithOp("db.find").WithSiteId(101)
	log.LogError(appErr)

	output := buf.String()
	if !strings.Contains(output, "record missing") || !strings.Contains(output, "NotFound") {
		t.Fatalf("expected JSON log to contain error details, got %s", output)
	}
}

func TestLoggerLevelFilterIgnored(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New(logger.DefaultOptions().WithOutput(buf).WithLevel(logger.LevelWarn))
	log.Debug("debug message")
	log.Info("info message")
	if buf.Len() > 0 {
		t.Fatalf("expected no output for levels below WARN, got %s", buf.String())
	}
}

func TestLoggerLevelFilterMatched(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New(logger.DefaultOptions().WithOutput(buf).WithLevel(logger.LevelWarn))
	log.Warn("warning message")
	if !strings.Contains(buf.String(), "warning message") {
		t.Fatalf("expected warning message in log output, got %s", buf.String())
	}
}

func TestLogLevelEnumAndJSON(t *testing.T) {
	lvl := logger.LevelWarn
	data, err := json.Marshal(lvl)
	if err != nil || string(data) != "\"Warn\"" || lvl.Name() != "Warn" {
		t.Fatalf("expected \"Warn\" JSON, got %s", string(data))
	}

	var parsed logger.LogLevel
	if err := json.Unmarshal([]byte("\"Error\""), &parsed); err != nil || parsed != logger.LevelError {
		t.Fatalf("expected LevelError, got %v", parsed)
	}
}
