// Package verbose provides a shared debug logger that writes to a timestamped
// log file when --verbose is enabled. All output is also echoed to stderr.
package verbose

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Logger writes verbose debug output to a file and optionally to stderr.
type Logger struct {
	file    *os.File
	enabled bool
}

// enabled tracks whether verbose mode is active globally.
var global *Logger

// Init creates the log file and enables verbose logging.
// Call once at startup when --verbose is set.
func Init() (*Logger, error) {
	logDir := constants.DefaultOutputFolder
	_ = os.MkdirAll(logDir, constants.DirPermission)

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logDir, fmt.Sprintf(constants.VerboseLogFileFmt, timestamp))

	file, err := os.Create(logPath)
	if err != nil {
		return nil, err
	}

	l := &Logger{file: file, enabled: true}
	global = l
	fmt.Printf(constants.MsgVerboseLogFile, logPath)

	return l, nil
}

// Close flushes and closes the log file.
func (l *Logger) Close() {
	if l != nil && l.file != nil {
		l.file.Close()
	}
}

// Log writes a formatted message to the log file and stderr.
func (l *Logger) Log(format string, args ...interface{}) {
	if l == nil {
		return
	}
	if l.enabled {
		writeLogEntry(l, format, args...)
	}
}

// writeLogEntry writes a timestamped log entry to file and stderr.
func writeLogEntry(l *Logger, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	ts := time.Now().Format("15:04:05.000")
	entry := fmt.Sprintf("[%s] %s\n", ts, line)
	_, _ = l.file.WriteString(entry)
	fmt.Fprint(os.Stderr, constants.ColorDim+entry+constants.ColorReset)
}

// IsEnabled returns true if verbose mode is active.
func IsEnabled() bool {
	return global != nil && global.enabled
}

// Get returns the global logger (may be nil).
func Get() *Logger {
	return global
}
