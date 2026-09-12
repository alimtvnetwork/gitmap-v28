package constants

// LFS smudge detection and auto-fix constants for clone-from.
const (
	LFSSmudgeFatalPattern           = `(?m)^fatal: (.*?): smudge filter lfs failed`
	LFSSmudgeErrorPattern           = `(?m)^Error downloading object: ([^\s]+)`
	LFSSmudgeFilterFailedSignature  = "smudge filter lfs failed"
	LFSServerObjectMissingSignature = "Object does not exist on the server"
	LFSSubmatchLength               = 2
	LFSFixCommitMessage             = "chore(lfs): remove pointer for missing LFS object"
	ErrCloneFromLFSRestore          = "git restore failed: %v\n%s"
	ErrCloneFromLFSRmCached         = "git rm --cached failed: %v\n%s"
	ErrCloneFromLFSCommit           = "git commit failed: %v\n%s"
	ErrCloneFromLFSPush             = "git push failed: %v\n%s"
)
