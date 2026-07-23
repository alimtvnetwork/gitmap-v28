package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// checkReleaseRepoIntegrity verifies the Release↔Repo FK relationship:
//   - orphaned Release rows (RepoId pointing to a non-existent Repo)
//   - Repo rows with zero releases (informational, not an error)
func checkReleaseRepoIntegrity() int {
	db, err := store.OpenDefault()
	if err != nil {
		return 0 // checkDatabase already reports this
	}
	defer db.Close()

	orphaned, reposNoRel, err := db.ReleaseRepoIntegrity()
	if err != nil {
		printWarn(fmt.Sprintf(constants.DoctorIntegrityFail, err))

		return 0
	}

	issues := 0
	if orphaned > 0 {
		printIssue(
			fmt.Sprintf(constants.DoctorOrphanedReleases, orphaned),
			constants.DoctorOrphanedDetail,
		)
		printFix(constants.DoctorOrphanedFix)
		issues++
	} else {
		printOK(constants.DoctorNoOrphans)
	}

	if reposNoRel > 0 {
		printWarn(fmt.Sprintf(constants.DoctorReposNoReleases, reposNoRel))
	}

	return issues
}
