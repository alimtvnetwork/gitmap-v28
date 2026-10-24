package enums

// ConventionalVerbType defines uppercase conventional commit type tokens.
type ConventionalVerbType string

const (
	ConventionalFeat     ConventionalVerbType = "FEAT"
	ConventionalFix      ConventionalVerbType = "FIX"
	ConventionalChore    ConventionalVerbType = "CHORE"
	ConventionalRefactor ConventionalVerbType = "REFACTOR"
	ConventionalRelease  ConventionalVerbType = "RELEASE"
)

// IsValid reports whether the type matches a recognized conventional verb.
func (c ConventionalVerbType) IsValid() bool {
	switch c {
	case ConventionalFeat, ConventionalFix, ConventionalChore,
		ConventionalRefactor, ConventionalRelease:
		return true
	default:
		return false
	}
}

// String returns the string representation of the conventional verb.
func (c ConventionalVerbType) String() string {
	return string(c)
}

// AllConventionalVerbs returns all registered ConventionalVerbType variants.
func AllConventionalVerbs() []ConventionalVerbType {
	return []ConventionalVerbType{
		ConventionalFeat,
		ConventionalFix,
		ConventionalChore,
		ConventionalRefactor,
		ConventionalRelease,
	}
}
