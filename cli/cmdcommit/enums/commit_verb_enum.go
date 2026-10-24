package enums

import "strings"

// CommitVerbType defines CLI commit verbs and automated workflow prefixes.
type CommitVerbType string

const (
	CommitVerbFeature  CommitVerbType = "cpf"
	CommitVerbBug      CommitVerbType = "cbf"
	CommitVerbRefactor CommitVerbType = "cfr"
	CommitVerbChore    CommitVerbType = "cpc"
	CommitVerbRelease  CommitVerbType = "cpr"
	CommitVerbFlat     CommitVerbType = "cm"
	CommitVerbPullPush CommitVerbType = "pcp"
)

// IsValid reports whether the verb matches a recognized commit workflow.
func (v CommitVerbType) IsValid() bool {
	switch v {
	case CommitVerbFeature, CommitVerbBug, CommitVerbRefactor,
		CommitVerbChore, CommitVerbRelease, CommitVerbFlat, CommitVerbPullPush:
		return true
	default:
		return false
	}
}

// String returns the raw verb string identifier.
func (v CommitVerbType) String() string {
	return string(v)
}

// Prefix returns the commit message prefix associated with the verb.
func (v CommitVerbType) Prefix() string {
	switch v {
	case CommitVerbFeature:
		return "Feature: "
	case CommitVerbBug:
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

// UppercaseToken returns the conventional uppercase token for the verb.
func (v CommitVerbType) UppercaseToken() string {
	switch v {
	case CommitVerbFeature:
		return "FEAT"
	case CommitVerbBug:
		return "FIX"
	case CommitVerbChore:
		return "CHORE"
	case CommitVerbRefactor:
		return "REFACTOR"
	case CommitVerbRelease:
		return "RELEASE"
	default:
		return v.fallbackToken()
	}
}

func (v CommitVerbType) fallbackToken() string {
	if v == CommitVerbFlat {
		return "CM"
	}
	if v == CommitVerbPullPush {
		return "PCP"
	}

	return strings.ToUpper(string(v))
}

// IsFeature reports whether the verb represents a feature commit.
func (v CommitVerbType) IsFeature() bool {
	return v == CommitVerbFeature
}

// IsBug reports whether the verb represents a bugfix commit.
func (v CommitVerbType) IsBug() bool {
	return v == CommitVerbBug
}
