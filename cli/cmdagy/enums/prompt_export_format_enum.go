package enums

import (
	"strings"
)

// PromptExportFormatType defines supported export serialization formats.
type PromptExportFormatType string

const (
	// PromptExportFormatTypeJSON indicates export to JSON format.
	PromptExportFormatTypeJSON PromptExportFormatType = "json"

	// PromptExportFormatTypeSQL indicates export to SQL format.
	PromptExportFormatTypeSQL PromptExportFormatType = "sql"

	// PromptExportFormatTypeDB indicates export to SQLite .db format.
	PromptExportFormatTypeDB PromptExportFormatType = "db"
)

const (
	// PromptExportFormatJSON alias for PromptExportFormatTypeJSON.
	PromptExportFormatJSON = PromptExportFormatTypeJSON

	// PromptExportFormatSQL alias for PromptExportFormatTypeSQL.
	PromptExportFormatSQL = PromptExportFormatTypeSQL

	// PromptExportFormatDB alias for PromptExportFormatTypeDB.
	PromptExportFormatDB = PromptExportFormatTypeDB
)

// ParsePromptExportFormatType parses raw format string into PromptExportFormatType.
// It intentionally returns no error, matching the original cmdagy behavior:
// recognized names and dotted/sqlite spellings map to a known variant, and any
// unknown input passes through verbatim (so IsValid reports false for it)
// instead of being coerced to a default format.
func ParsePromptExportFormatType(raw string) PromptExportFormatType {
	lower := strings.ToLower(strings.TrimSpace(raw))

	switch lower {
	case "json", ".json":
		return PromptExportFormatTypeJSON
	case "sql", ".sql":
		return PromptExportFormatTypeSQL
	case "db", ".db", "sqlite", ".sqlite":
		return PromptExportFormatTypeDB
	default:
		return PromptExportFormatType(lower)
	}
}

// IsValid reports whether format is recognized.
func (f PromptExportFormatType) IsValid() bool {
	return f == PromptExportFormatTypeJSON || f == PromptExportFormatTypeSQL || f == PromptExportFormatTypeDB
}

// String returns the string representation of the export format.
func (f PromptExportFormatType) String() string {
	return string(f)
}

// AllPromptExportFormatTypes returns all valid PromptExportFormatType variants.
func AllPromptExportFormatTypes() []PromptExportFormatType {
	return []PromptExportFormatType{
		PromptExportFormatTypeJSON,
		PromptExportFormatTypeSQL,
		PromptExportFormatTypeDB,
	}
}

// IsJSON reports whether format is JSON.
func (f PromptExportFormatType) IsJSON() bool {
	return f == PromptExportFormatTypeJSON
}

// IsSQL reports whether format is SQL.
func (f PromptExportFormatType) IsSQL() bool {
	return f == PromptExportFormatTypeSQL
}

// IsDB reports whether format is SQLite DB.
func (f PromptExportFormatType) IsDB() bool {
	return f == PromptExportFormatTypeDB
}

// HasExportFormat reports whether the export format is non-empty and valid.
func (f PromptExportFormatType) HasExportFormat() bool {
	return f.IsValid()
}
