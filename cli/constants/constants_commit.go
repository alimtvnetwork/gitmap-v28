package constants

// CommitVerbType represents shorthand and longhand commit-push verbs.
type CommitVerbType string

const (
	CommitVerbFeature  CommitVerbType = "cpf"
	CommitVerbBug      CommitVerbType = "cbf"
	CommitVerbRefactor CommitVerbType = "cfr"
	CommitVerbChore    CommitVerbType = "cpc"
	CommitVerbRelease  CommitVerbType = "cpr"
	CommitVerbBugAlt   CommitVerbType = "cpb"
	CommitVerbFlat     CommitVerbType = "cm"
)

// IsValid reports whether the verb is a recognized commit-push verb.
func (v CommitVerbType) IsValid() bool {
	switch v {
	case CommitVerbFeature, CommitVerbBug, CommitVerbRefactor, CommitVerbChore, CommitVerbRelease, CommitVerbBugAlt, CommitVerbFlat:
		return true
	default:
		return false
	}
}

// IsFeature reports whether the verb represents a feature commit.
func (v CommitVerbType) IsFeature() bool {
	return v == CommitVerbFeature
}

// IsBug reports whether the verb represents a bug fix commit.
func (v CommitVerbType) IsBug() bool {
	return v == CommitVerbBug || v == CommitVerbBugAlt
}

// IsRefactor reports whether the verb represents a refactoring commit.
func (v CommitVerbType) IsRefactor() bool {
	return v == CommitVerbRefactor
}

// IsChore reports whether the verb represents a chore commit.
func (v CommitVerbType) IsChore() bool {
	return v == CommitVerbChore
}

// IsRelease reports whether the verb represents a release commit.
func (v CommitVerbType) IsRelease() bool {
	return v == CommitVerbRelease
}

// String returns the string representation of the verb.
func (v CommitVerbType) String() string {
	return string(v)
}

// Label returns the human-readable description of the verb.
func (v CommitVerbType) Label() string {
	switch v {
	case CommitVerbFeature:
		return "Feature Commit"
	case CommitVerbBug, CommitVerbBugAlt:
		return "Bug Fix Commit"
	case CommitVerbRefactor:
		return "Refactor Commit"
	case CommitVerbChore:
		return "Chore Commit"
	case CommitVerbRelease:
		return "Release Commit"
	case CommitVerbFlat:
		return "Flat Commit"
	default:
		return "Unknown Verb"
	}
}

// Prefix returns the standard commit title prefix.
func (v CommitVerbType) Prefix() string {
	switch v {
	case CommitVerbFeature:
		return "Feature: "
	case CommitVerbBug, CommitVerbBugAlt:
		return "Bug: "
	case CommitVerbRefactor:
		return "Refactor: "
	case CommitVerbChore:
		return "Chore: "
	case CommitVerbRelease:
		return "Release: "
	default:
		return ""
	}
}

// ConventionalCommitType represents standardized conventional commit type tokens.
type ConventionalCommitType string

const (
	ConventionalFeat     ConventionalCommitType = "FEAT"
	ConventionalFix      ConventionalCommitType = "FIX"
	ConventionalRefactor ConventionalCommitType = "REFACTOR"
	ConventionalChore    ConventionalCommitType = "CHORE"
	ConventionalRelease  ConventionalCommitType = "RELEASE"
)
