package applogger

import (
	"coding-guidelines/common/pkg/enum/logleveltype"
)

const (
	LevelUnknown = logleveltype.Unknown
	LevelInvalid = logleveltype.Invalid
	LevelDebug   = logleveltype.Debug
	LevelInfo    = logleveltype.Info
	LevelWarn    = logleveltype.Warn
	LevelError   = logleveltype.Error
	LevelFatal   = logleveltype.Fatal
)

type DriverType byte

const (
	DriverConsole DriverType = iota
	DriverFile
	DriverSQLite
	DriverZap
	DriverComposite
)

const createLogsTableSQL = `CREATE TABLE IF NOT EXISTS app_logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp TEXT NOT NULL,
	level TEXT NOT NULL,
	message TEXT NOT NULL,
	caller TEXT,
	fields_json TEXT,
	stack_trace TEXT
);`

type EntryFormatter func(entry LogEntry) string

type EntryFilterFunc func(entry LogEntry) bool
