package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// appendPRHelp appends PR and Commit Transfer documentation to the provided writer.
func appendPRHelp(cmd *cobra.Command, args []string, buf io.Writer) error {
	_, err := fmt.Fprint(buf, `
## PR (Pull Request) & Commit Transfer Commands

The 'pr' command family provides automated commit replay, merge detection, feature branch simulation,
and SQLite-backed PR lifecycle management across repositories.

Commands:
- gitmap pr LEFT RIGHT
- gitmap pr in <target> <inputs...>
- gitmap pr-clean [repo] [--yes]
- gitmap pr-list [repo] [--json]

### Automated Merge & PR Simulation
When commits are replayed, merges detected in the source repository are automatically isolated
into dedicated feature branches (feature/<slug> or pr/<id>), generate rich Markdown PR descriptions,
and merge cleanly into the mainline target (--no-ff) while saving PR records into SQLite
(.gitmap/data/pr/<slug>/sql.db).

### Final Snapshot Synchronization
At the conclusion of replay, gitmap synchronizes the entire target working tree with the source's
latest commit/release snapshot, pruning any stray target-only files and ensuring byte-for-byte fidelity.

### JSON Example:
{
  "command": "pr",
  "source": "./upstream-repo",
  "target": "./main-repo",
  "prMode": "merges",
  "since": "HEAD~20",
  "finalSnapshotSync": true
}
`)
	if err != nil {
		return apperror.WrapSimple(err, "appendPRHelp")
	}

	return nil
}
