package constants

// Repo-reclone (single-repo wipe + re-clone) constants.
//
// Triggered when `gitmap reclone` / `rc` / `rec` is invoked with NO
// manifest file (positional or auto-discovered) but the cwd, or the
// single positional argument, IS a git repo. The handler deletes
// that working tree and re-clones it from `remote.origin.url` into
// the same parent directory.
//
// Spec: spec/04-generic-cli/32-repo-reclone.md (single-repo flow).
const (
	FlagRepoRecloneYes     = "y"
	FlagDescRepoRecloneYes = "Skip the destructive confirmation prompt"
)

// User-facing strings. Centralised per the no-magic-strings rule.
const (
	MsgRepoReclonePlan    = "repo-reclone: target=%s origin=%s parent=%s\n"
	MsgRepoRecloneConfirm = "repo-reclone: about to DELETE %q and re-clone from %s.\n" +
		"  Type 'y' to continue, anything else aborts: "
	MsgRepoRecloneAborted   = "repo-reclone: aborted by user; nothing was removed\n"
	MsgRepoRecloneRemoving  = "repo-reclone: removing %s\n"
	MsgRepoRecloneCloning   = "repo-reclone: git clone %s → %s\n"
	MsgRepoRecloneDone      = "repo-reclone: done (%s)\n"
	ErrRepoRecloneNotGit    = "repo-reclone: %q is not a git repo (no .git directory)\n"
	ErrRepoRecloneNoOrigin  = "repo-reclone: cannot read remote.origin.url in %s: %v\n"
	ErrRepoRecloneRemove    = "repo-reclone: removing %s failed: %v\n"
	ErrRepoRecloneClone     = "repo-reclone: git clone %s into %s failed: %v\n"
	ErrRepoRecloneNonTTY    = "repo-reclone: refusing to prompt on a non-interactive stream; pass -y to confirm\n"
	ErrRepoRecloneBadTarget = "repo-reclone: target %q does not exist: %v\n"
)
