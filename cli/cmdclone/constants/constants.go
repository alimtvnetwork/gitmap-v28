package constants

// Clone flags and UI banners for cmdclone.
const (
	FlagMaxConcurrency = "max-concurrency"
	FlagWorkers        = "workers"
	FlagParallel       = "parallel"
	FlagTargetDir      = "target-dir"
	FlagOnly           = "only"
	FlagExclude        = "exclude"
	FlagToken          = "token"
	FlagSSHKey         = "ssh-key"
	FlagClean          = "clean"
	FlagForce          = "force"
	FlagReclone        = "reclone"
	FlagSafePull       = "safe-pull"
	FlagRepoCache      = "repo-cache"

	MsgMultiCloneNoURLs      = "  ⚠ No valid repository URLs or slugs found in input.\n  Run 'gitmap mc --help' to see supported input formats.\n"
	MsgMultiCloneDryRunTitle = "\n%s[GitMap MultiClone — Dry Run Preview]%s\n"
)
