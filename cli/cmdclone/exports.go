package cmdclone

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdtask"
)

// RunClone handles the "clone" subcommand.
func RunClone(args []string) error {
	err := runClone(args)
	recordCloneAudit("clone", args, err)
	return err
}

// RunCloneFixRepo handles the "clone-fix-repo" subcommand.
func RunCloneFixRepo(args []string) error {
	err := runCloneFixRepo(args)
	recordCloneAudit("clone-fix-repo", args, err)
	return err
}

// RunCloneFixRepoPub handles the "clone-fix-repo-pub" subcommand.
func RunCloneFixRepoPub(args []string) error {
	err := runCloneFixRepoPub(args)
	recordCloneAudit("clone-fix-repo-pub", args, err)
	return err
}

func recordCloneAudit(action string, args []string, err error) {
	status := "completed"
	if err != nil {
		status = "failed"
	}
	target := strings.Join(args, " ")
	cmdtask.RecordTaskAudit("clone", action, target, "", status)
}

// RunCloneFrom handles the "clone-from" subcommand.
func RunCloneFrom(args []string) error {
	return runCloneFrom(args)
}

// RunCloneNext handles the "clone-next" subcommand.
func RunCloneNext(args []string) error {
	return runCloneNext(args)
}

// RunCloneNow handles the "clone-now" subcommand.
func RunCloneNow(args []string) error {
	return runCloneNow(args)
}

// RunClonePick handles the "clone-pick" subcommand.
func RunClonePick(args []string) error {
	return runClonePick(args)
}

// RunCloneSync handles the "clone-sync" subcommand.
func RunCloneSync() error {
	return runCloneSync()
}

// ExtractRepoName extracts repository name from URL.
func ExtractRepoName(rawURL string) string {
	return extractRepoName(rawURL)
}

// RepoNameFromURL extracts repository name from URL.
func RepoNameFromURL(rawURL string) string {
	return repoNameFromURL(rawURL)
}

// ParseCloneFixRepoArgs parses clone fix repo args.
func ParseCloneFixRepoArgs(args []string) (string, string, bool, bool, bool, bool, bool, bool, bool, bool) {
	return parseCloneFixRepoArgs(args)
}

// OpenInVSCode opens a target directory in VSCode.
func OpenInVSCode(target string) {
	openInVSCode(target)
}

// RegisterSingleDesktop registers desktop file for a repository.
func RegisterSingleDesktop(name, absPath string) {
	registerSingleDesktop(name, absPath)
}

// ResolveCloneNextFolder resolves clone next folder.
func ResolveCloneNextFolder(token string) (string, error) {
	return resolveCloneNextFolder(token)
}

// ResolveCloneShorthand resolves clone shorthand.
func ResolveCloneShorthand(arg string) string {
	return resolveCloneShorthand(arg)
}

// UpsertDirectClone upserts a direct clone record.
func UpsertDirectClone(url, repoName, folderName, absPath string) {
	upsertDirectClone(url, repoName, folderName, absPath)
}

// IsDirectURL reports whether source looks like a git URL.
func IsDirectURL(source string) bool {
	return isDirectURL(source)
}
