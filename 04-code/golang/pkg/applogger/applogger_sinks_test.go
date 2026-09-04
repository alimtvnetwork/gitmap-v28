package applogger_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"coding-guidelines/common/pkg/applogger"
)

type mockZapLogger struct {
	logs []string
}

func (m *mockZapLogger) Debug(args ...any) { m.logs = append(m.logs, "DEBUG") }
func (m *mockZapLogger) Info(args ...any)  { m.logs = append(m.logs, "INFO") }
func (m *mockZapLogger) Warn(args ...any)  { m.logs = append(m.logs, "WARN") }
func (m *mockZapLogger) Error(args ...any) { m.logs = append(m.logs, "ERROR") }
func (m *mockZapLogger) Fatal(args ...any) { m.logs = append(m.logs, "FATAL") }
func (m *mockZapLogger) Sync() error       { return nil }

func TestFileSinkLogger(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")
	cfg := applogger.Config{MinLevel: applogger.LevelInfo, Driver: applogger.DriverFile, FilePath: logFile}
	l, _ := applogger.New(cfg)
	l.Info("file log message")
	_ = l.Close()

	content, err := os.ReadFile(logFile)
	if err != nil || !strings.Contains(string(content), "file log message") {
		t.Fatalf("expected file to contain log message, got %s", string(content))
	}
}

func TestZapAdapterDriver(t *testing.T) {
	mockZap := &mockZapLogger{}
	cfg := applogger.Config{MinLevel: applogger.LevelInfo, Driver: applogger.DriverZap, ZapLogger: mockZap}
	l, _ := applogger.New(cfg)
	l.Info("test zap info")
	l.Error("test zap error")

	if len(mockZap.logs) != 2 {
		t.Fatalf("expected 2 zap log calls, got %d", len(mockZap.logs))
	}
}
