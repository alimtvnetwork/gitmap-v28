package cmdrelease

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/scanpipe"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"os"
	"path/filepath"
)

// autoRegisterCurrentRepo registers the cwd as a Repo + ScanFolder so
// release persistence can satisfy the Release.RepoId FK without forcing
// the user to run `gitmap scan` manually first.
//
// Strategy: parent dir is the ScanFolder, cwd is the single Repo. We
// reuse scanpipe.BuildRecords so slug / URLs / branch are populated the
// same way a regular scan would.
func autoRegisterCurrentRepo(db *store.DB, cwd string) error {
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return apperror.WrapSimple(err, "could not resolve cwd")
	}

	parent := filepath.Dir(absCwd)
	repoInfo := scanpipe.RepoInfo{
		AbsolutePath: absCwd,
		RelativePath: filepath.Base(absCwd),
	}

	records := scanpipe.BuildRecords([]scanpipe.RepoInfo{repoInfo}, "https", "")

	if err := db.UpsertRepos(records); err != nil {
		return apperror.WrapSimple(err, "upsert repo failed")
	}

	folder, err := db.EnsureScanFolder(parent, "", "")
	if err != nil {
		return apperror.WrapSimple(err, "ensure scan folder failed")
	}

	if err := db.TagReposByScanFolder(folder.ID, []string{absCwd}); err != nil {
		return apperror.WrapSimple(err, "tag repo failed")
	}

	// Trailing blank line: separates the release summary from the next
	// shell prompt so the terminal output doesn't visually run into PS1.
	fmt.Fprintf(os.Stdout, "  ✓ Auto-registered repo %q under scan folder %q (#%d)\n\n",
		absCwd, parent, folder.ID)

	return nil
}
