package enums

// SplitDatabaseType identifies the functional domain of a split SQLite database.
type SplitDatabaseType string

const (
	SplitDBTypeRepodb        SplitDatabaseType = "repodb"
	SplitDBTypeSearch        SplitDatabaseType = "search"
	SplitDBTypeTasks         SplitDatabaseType = "tasks"
	SplitDBTypeAiInstruction SplitDatabaseType = "ai-instruction"
	SplitDBTypeAiAnalysis    SplitDatabaseType = "ai_analysis"
	SplitDBTypeStartup       SplitDatabaseType = "startup"
)

// IsValid reports whether the database type is a registered split database domain.
func (d SplitDatabaseType) IsValid() bool {
	switch d {
	case SplitDBTypeRepodb, SplitDBTypeSearch, SplitDBTypeTasks,
		SplitDBTypeAiInstruction, SplitDBTypeAiAnalysis, SplitDBTypeStartup:
		return true
	default:
		return false
	}
}

// String returns the string representation of the split database type.
func (d SplitDatabaseType) String() string {
	return string(d)
}

// AllSplitDatabaseTypes returns all functional SplitDatabaseType variants.
func AllSplitDatabaseTypes() []SplitDatabaseType {
	return []SplitDatabaseType{
		SplitDBTypeRepodb,
		SplitDBTypeSearch,
		SplitDBTypeTasks,
		SplitDBTypeAiInstruction,
		SplitDBTypeAiAnalysis,
		SplitDBTypeStartup,
	}
}
