package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// recordVersionHistory persists version transition data in the database.
func recordVersionHistory(absPath string, fromVersion, toVersion int, flattenedPath string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not open database for version history: %v\n", err)

		return
	}
	defer db.Close()

	// Find or create the repo record by absolute path.
	repoID, findErr := db.GetRepoIDByPath(absPath)
	if findErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: repo not found in database for version tracking: %v\n", findErr)

		return
	}

	fromTag := fmt.Sprintf("v%d", fromVersion)
	toTag := fmt.Sprintf("v%d", toVersion)

	// Update current version on the Repos row.
	if updateErr := db.UpdateRepoVersion(repoID, toTag, toVersion); updateErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not update repo version: %v\n", updateErr)
	}

	// Stamp clone-recency so `gitmap release` can auto-cd here from a parent dir.
	if markErr := db.MarkCloned(absPath); markErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not stamp clone time: %v\n", markErr)
	}

	// Insert version history record.
	record := model.RepoVersionHistoryRecord{
		RepoID:         repoID,
		FromVersionTag: fromTag,
		FromVersionNum: fromVersion,
		ToVersionTag:   toTag,
		ToVersionNum:   toVersion,
		FlattenedPath:  flattenedPath,
	}

	if _, insertErr := db.InsertVersionHistory(record); insertErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not record version history: %v\n", insertErr)

		return
	}

	fmt.Printf(constants.MsgFlattenVersionDB, fromVersion, toVersion)
}
