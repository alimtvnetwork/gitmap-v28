package constants

// Git commands and arguments.
const (
	GitBin             = "git"
	SSHBin             = "ssh"
	GitClone           = "clone"
	GitPull            = "pull"
	GitRebase          = "rebase"
	GitBranchFlag      = "-b"
	GitDirFlag         = "-C"
	GitFFOnlyFlag      = "--ff-only"
	GitPullRebaseFlag  = "--rebase"
	GitNoRebaseFlag    = "--no-rebase"
	GitRebaseAbortFlag = "--abort"
	GitConfigCmd       = "config"
	GitGetFlag         = "--get"
	GitRemoteOrigin    = "remote.origin.url"
	GitRevParse        = "rev-parse"
	GitAbbrevRef       = "--abbrev-ref"
	GitHEAD            = "HEAD"
	GitTag             = "tag"
	GitCheckout        = "checkout"
	GitPush            = "push"
	GitLsRemote        = "ls-remote"
	GitLsRemoteTags    = "--tags"
	GitOrigin          = "origin"
	GitOriginPrefix    = "origin/"
	GitCommitPrefix    = "commit:"
	GitTagAnnotateFlag = "-a"
	GitTagMessageFlag  = "-m"
	GitTagListFlag     = "--list"
	GitBranchListFlag  = "--list"
	GitCatFile         = "cat-file"
	GitCatFileTypeFlag = "-t"
	GitCommitType      = "commit"
	GitTagGlob         = "v*"
	// GitSymbolicRef + GitRefsRemotesOriginHEAD power
	// gitutil.ResolveDefaultBranchName, which backs `gitmap branch
	// default` / `b def`. We ask `git symbolic-ref refs/remotes/origin/HEAD`
	// to learn the upstream's chosen default (commonly main, sometimes
	// master, occasionally trunk/develop) instead of guessing.
	GitSymbolicRef           = "symbolic-ref"
	GitRefsRemotesOriginHEAD = "refs/remotes/origin/HEAD"
)

// SSH environment and option constants used by git clone wrappers.
const (
	EnvGitSSHCommand          = "GIT_SSH_COMMAND"
	SSHOptionFlag             = "-o"
	SSHStrictHostKeyChecking  = "StrictHostKeyChecking"
	SSHStrictHostKeyAcceptNew = "StrictHostKeyChecking=accept-new"
	EnvAssignmentSeparator    = "="
)

// Git arguments for latest-branch operations.
const (
	GitFetch              = "fetch"
	GitBranch             = "branch"
	GitLog                = "log"
	GitForEachRef         = "for-each-ref"
	GitArgAll             = "--all"
	GitArgPrune           = "--prune"
	GitArgRemote          = "-r"
	GitArgContains        = "--contains"
	GitArgInsideWorkTree  = "--is-inside-work-tree"
	GitLogTipFormat       = "--format=%cI|%H|%s"
	GitLogDelimiter       = "|"
	GitLogFieldCount      = 3
	GitPointsAtFmt        = "--points-at=%s"
	GitRefsRemotesFmt     = "refs/remotes/%s"
	GitFormatRefnameShort = "--format=%(refname:short)"
	GitForEachRefTagFmt   = "--format=%(refname:short)|%(creatordate:iso-strict)"
	GitRefsTagsPrefix     = "refs/tags/"
	HeadPointer           = " -> "
	ShaDisplayLength      = 7
)

// Clone instruction format.
const (
	CloneInstructionFmt = "git clone -b %s %s %s"
	HTTPSFromSSHFmt     = "https://%s/%s"
	SSHFromHTTPSFmt     = "git@%s:%s"
)

// GitStderrNoisePatterns lists substrings whose presence in a git stderr
// line marks the line as harmless cosmetic chatter that the release
// pipeline should suppress. Add ONLY substrings that are guaranteed to
// not appear in a real error — these lines are silently dropped.
var GitStderrNoisePatterns = []string{
	"LF will be replaced by CRLF the next time Git touches it",
}
