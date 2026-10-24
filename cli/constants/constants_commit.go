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
)

// ConventionalCommitType represents standardized conventional commit type tokens.
type ConventionalCommitType string

const (
	ConventionalFeat     ConventionalCommitType = "FEAT"
	ConventionalFix      ConventionalCommitType = "FIX"
	ConventionalRefactor ConventionalCommitType = "REFACTOR"
	ConventionalChore    ConventionalCommitType = "CHORE"
	ConventionalRelease  ConventionalCommitType = "RELEASE"
)
