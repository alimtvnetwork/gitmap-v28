package constants

// Git commands and arguments.
const (
	GitBin                     = "git"
	SSHBin                     = "ssh"
	GitClone                   = "clone"
	GitPull                    = "pull"
	GitRebase                  = "rebase"
	GitStatus                  = "status"
	GitPorcelainFlag           = "--porcelain"
	GitRemoteCmd               = "remote"
	GitSetURLCmd               = "set-url"
	GitProgressFlag            = "--progress"
	GitNoEditFlag              = "--no-edit"
	GitAutostashFlag           = "--autostash"
	GitLimitOneFlag            = "-1"
	GitFormatRelative          = "--format=%ar"
	GHBin                      = "gh"
	GHRun                      = "run"
	GHList                     = "list"
	GitBranchFlag              = "-b"
	GitDirFlag                 = "-C"
	GitFFOnlyFlag              = "--ff-only"
	GitPullRebaseFlag          = "--rebase"
	GitNoRebaseFlag            = "--no-rebase"
	GitRebaseAbortFlag         = "--abort"
	GitConfigCmd               = "config"
	GitGetFlag                 = "--get"
	GitRemoteOrigin            = "remote.origin.url"
	GitRevParse                = "rev-parse"
	GitAbbrevRef               = "--abbrev-ref"
	GitHEAD                    = "HEAD"
	GitTag                     = "tag"
	GitCheckout                = "checkout"
	GitPush                    = "push"
	GitLsRemote                = "ls-remote"
	GitLsRemoteTags            = "--tags"
	GitOrigin                  = "origin"
	GitOriginPrefix            = "origin/"
	GitCommitPrefix            = "commit:"
	GitTagAnnotateFlag         = "-a"
	GitTagMessageFlag          = "-m"
	GitTagListFlag             = "--list"
	GitBranchListFlag          = "--list"
	GitCatFile                 = "cat-file"
	GitCatFileTypeFlag         = "-t"
	GitCommitType              = "commit"
	GitTagGlob                 = "v*"
	GitRestore                 = "restore"
	GitRm                      = "rm"
	GitRestoreSourceHEAD       = "--source=HEAD"
	GitPathspecRoot            = ":/"
	GitCachedFlag              = "--cached"
	GitQuietFlag               = "-q"
	DotGit                     = ".git"
	GitConfigFile              = "config"
	GitdirPrefix               = "gitdir:"
	GitRefPrefix               = "ref:"
	GitConfigRemoteOrigin      = `[remote "origin"]`
	GitConfigSectionStart      = "["
	GitConfigURLPrefix         = "url"
	GitConfigKeyValueSeparator = "="
	GitConfigSplitPartsCount   = 2
	GitLineSeparator           = "\n"
	EnvGitLFSSkipSmudge        = "GIT_LFS_SKIP_SMUDGE=1"
	// GitSymbolicRef + GitRefsRemotesOriginHEAD power
	// gitutil.ResolveDefaultBranchName, which backs `gitmap branch
	// default` / `b def`. We ask `git symbolic-ref refs/remotes/origin/HEAD`
	// to learn the upstream's chosen default (commonly main, sometimes
	// master, occasionally trunk/develop) instead of guessing.
	GitSymbolicRef           = "symbolic-ref"
	GitRefsRemotesOriginHEAD = "refs/remotes/origin/HEAD"
)

// GitRefCommandType represents Git subcommands operating on references and symbolic pointers.
type GitRefCommandType string

const (
	GitRefCmdSymbolicRef GitRefCommandType = "symbolic-ref"
	GitRefCmdRevParse    GitRefCommandType = "rev-parse"
	GitRefCmdForEachRef  GitRefCommandType = "for-each-ref"
)

// GitRevCommandType represents Git subcommands operating on revision lists and ranges.
type GitRevCommandType string

const (
	GitRevList GitRevCommandType = "rev-list"
)

// Git revision query flags.
const (
	GitCountFlag = "--count"
)

// GitInspectCommandType represents Git subcommands inspecting histories, diffs, and relationships.
type GitInspectCommandType string

const (
	GitCmdCherry    GitInspectCommandType = "cherry"
	GitCmdDiff      GitInspectCommandType = "diff"
	GitCmdMergeBase GitInspectCommandType = "merge-base"
)

const (
	GitCherry    = "cherry"
	GitDiff      = "diff"
	GitMergeBase = "merge-base"
)

// GitFlagType represents command-line flags supplied to Git subprocess executions.
type GitFlagType string

const (
	GitFlagStat GitFlagType = "--stat"
)

const (
	GitStatFlag = "--stat"
)

// GitRemoteNameType represents standard remote identifiers across Git repositories.
type GitRemoteNameType string

const (
	GitRemoteNameOrigin   GitRemoteNameType = "origin"
	GitRemoteNameUpstream GitRemoteNameType = "upstream"
)

// SSH environment and option constants used by git clone wrappers.
const (
	TransportSSH              = "ssh"
	TransportHTTPS            = "https"
	EnvGitSSHCommand          = "GIT_SSH_COMMAND"
	SSHOptionFlag             = "-o"
	SSHStrictHostKeyChecking  = "StrictHostKeyChecking"
	SSHStrictHostKeyAcceptNew = "StrictHostKeyChecking=accept-new"
	EnvAssignmentSeparator    = "="
	EnvGitTerminalPromptZero  = "GIT_TERMINAL_PROMPT=0"
	EnvGitAskpassEmpty        = "GIT_ASKPASS="
	EnvSSHAskpassEmpty        = "SSH_ASKPASS="
	EnvGCMInteractiveNever    = "GCM_INTERACTIVE=never"
	EnvGitSSHCommandBatchYes  = "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=accept-new -o BatchMode=yes"
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

// CredentialStoreType represents the storage backend for Git credentials.
type CredentialStoreType string

const (
	CredentialStoreDPAPI      CredentialStoreType = "dpapi"
	CredentialStoreWinCredMan CredentialStoreType = "wincredman"
	CredentialStoreCache      CredentialStoreType = "cache"
)

// Git credential store environment variables and headless SSH flags.
const (
	EnvGitMapHeadlessSSH       = "GITMAP_HEADLESS_SSH"
	EnvSSHConnection           = "SSH_CONNECTION"
	EnvSSHClient               = "SSH_CLIENT"
	EnvSSHTTY                  = "SSH_TTY"
	EnvGCMCredentialStoreDPAPI = "GCM_CREDENTIAL_STORE=dpapi"
	EnvGCMCredentialStoreWinCM = "GCM_CREDENTIAL_STORE=wincredman"
	EnvGCMCredentialStoreCache = "GCM_CREDENTIAL_STORE=cache"
	EnvGCMNoPersistOne         = "GCM_NO_PERSIST=1"
	EnvValueOne                = "1"
)

// GitProviderType defines recognized Git hosting providers.
type GitProviderType string

const (
	GitProviderGitHub      GitProviderType = "github"
	GitProviderGitLab      GitProviderType = "gitlab"
	GitProviderBitbucket   GitProviderType = "bitbucket"
	GitProviderAzureDevOps GitProviderType = "azuredevops"
	GitProviderUnknown     GitProviderType = "unknown"
)

// GitHostType defines normalized domain names for Git hosts.
type GitHostType string

const (
	GitHostGitHub      GitHostType = "github.com"
	GitHostGitLab      GitHostType = "gitlab.com"
	GitHostBitbucket   GitHostType = "bitbucket.org"
	GitHostAzureDevOps GitHostType = "dev.azure.com"
)

// GitURLPrefixType represents canonical protocol and scheme prefixes for Git remotes.
type GitURLPrefixType string

const (
	// GitHub remote prefixes.
	GitHubHTTPSPrefix     GitURLPrefixType = "https://github.com/"
	GitHubHTTPPrefix      GitURLPrefixType = "http://github.com/"
	GitHubSSHPrefix       GitURLPrefixType = "git@github.com:"
	GitHubSSHSchemePrefix GitURLPrefixType = "ssh://git@github.com/"

	// GitLab remote prefixes.
	GitLabHTTPSPrefix     GitURLPrefixType = "https://gitlab.com/"
	GitLabHTTPPrefix      GitURLPrefixType = "http://gitlab.com/"
	GitLabSSHPrefix       GitURLPrefixType = "git@gitlab.com:"
	GitLabSSHSchemePrefix GitURLPrefixType = "ssh://git@gitlab.com/"

	// Bitbucket remote prefixes.
	BitbucketHTTPSPrefix     GitURLPrefixType = "https://bitbucket.org/"
	BitbucketHTTPPrefix      GitURLPrefixType = "http://bitbucket.org/"
	BitbucketSSHPrefix       GitURLPrefixType = "git@bitbucket.org:"
	BitbucketSSHSchemePrefix GitURLPrefixType = "ssh://git@bitbucket.org/"

	// Azure DevOps remote prefixes.
	AzureHTTPSPrefix GitURLPrefixType = "https://dev.azure.com/"
	AzureSSHPrefix   GitURLPrefixType = "git@ssh.dev.azure.com:v3/"

	// Generic transport scheme prefixes.
	GitPrefixHTTPS     GitURLPrefixType = "https://"
	GitPrefixHTTP      GitURLPrefixType = "http://"
	GitPrefixSSH       GitURLPrefixType = "git@"
	GitPrefixSSHScheme GitURLPrefixType = "ssh://"
	GitPrefixGit       GitURLPrefixType = "git://"
	GitPrefixFile      GitURLPrefixType = "file://"
	GitPrefixSSHGitAt  GitURLPrefixType = "ssh://git@"
)

// Git URL formats, delimiters, and default values.
const (
	GitHubSSHRepoOwnerFmt = "git@github.com:%s/"
	GitHubSSHRepoFmt      = "git@github.com:%s.git"
	GitHubHTTPSRepoFmt    = "https://github.com/%s"
	GitHubHTTPSRepoGitFmt = "https://github.com/%s.git"
	DefaultGitHubOwner    = "alimtvnetwork"

	// GitDotGitSuffix represents the canonical repository extension for remote URLs.
	GitDotGitSuffix = ".git"

	// Delimiters for remote URL construction.
	GitPathSeparator  = "/"
	GitColonSeparator = ":"

	// Host path tokens for boundary tokenization.
	GitHostTokenGitHubSlash = "github.com/"
	GitHostTokenGitHubColon = "github.com:"
)
