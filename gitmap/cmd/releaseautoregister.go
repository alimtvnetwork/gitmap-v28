package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/mapper"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/scanner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// autoRegisterCurrentRepo registers the cwd as a Repo + ScanFolder so
// release persistence can satisfy the Release.RepoId FK without forcing
// the user to run `gitmap scan` manually first.
//
// Strategy: parent dir is the ScanFolder, cwd is the single Repo. We
// reuse mapper.BuildRecords so slug / URLs / branch are populated the
// same way a regular scan would.
func autoRegisterCurrentRepo(db *store.DB, cwd string) error {
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return fmt.Errorf("could not resolve cwd: %w", err)
	}

	parent := filepath.Dir(absCwd)
	repoInfo := scanner.RepoInfo{
		AbsolutePath: absCwd,
		RelativePath: filepath.Base(absCwd),
	}
	records := mapper.BuildRecords([]scanner.RepoInfo{repoInfo}, "https", "")

	if err := db.UpsertRepos(records); err != nil {
		return fmt.Errorf("upsert repo failed: %w", err)
	}

	folder, err := db.EnsureScanFolder(parent, "", "")
	if err != nil {
		return fmt.Errorf("ensure scan folder failed: %w", err)
	}

	if err := db.TagReposByScanFolder(folder.ID, []string{absCwd}); err != nil {
		return fmt.Errorf("tag repo failed: %w", err)
	}

	// Trailing blank line: separates the release summary from the next
	// shell prompt so the terminal output doesn't visually run into PS1.
	fmt.Fprintf(os.Stdout, "  ✓ Auto-registered repo %q under scan folder %q (#%d)\n\n",
		absCwd, parent, folder.ID)

	return nil
}
