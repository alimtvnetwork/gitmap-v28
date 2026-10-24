package enums

// CommitFormatType designates the detected format style of a commit message.
type CommitFormatType byte

const (
	FormatStandardBracket CommitFormatType = iota + 1
	FormatHyphenSeparated
	FormatConventional
)

// IsValid reports whether the format type is a recognized format option.
func (f CommitFormatType) IsValid() bool {
	switch f {
	case FormatStandardBracket, FormatHyphenSeparated, FormatConventional:
		return true
	default:
		return false
	}
}

// String returns a descriptive name for the commit format.
func (f CommitFormatType) String() string {
	switch f {
	case FormatStandardBracket:
		return "StandardBracket"
	case FormatHyphenSeparated:
		return "HyphenSeparated"
	case FormatConventional:
		return "Conventional"
	default:
		return "Unknown"
	}
}

// IsHyphenSeparated reports whether the format is hyphen-separated syntax.
func (f CommitFormatType) IsHyphenSeparated() bool {
	return f == FormatHyphenSeparated
}

// IsConventional reports whether the format is conventional commit syntax.
func (f CommitFormatType) IsConventional() bool {
	return f == FormatConventional
}
