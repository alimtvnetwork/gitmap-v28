// Package cmd — scanprojectsmeta.go handles Go and C# metadata persistence.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/detector"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// upsertGoProjectMeta persists Go metadata and runnables.
func upsertGoProjectMeta(db *store.DB, r detector.DetectionResult) {
	r.GoMeta.DetectedProjectID = r.Project.ID
	if err := db.UpsertGoMetadata(*r.GoMeta); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoMetadataUpsert, err)

		return
	}
	saved, err := db.SelectGoMetadata(r.Project.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoMetadataUpsert, err)

		return
	}
	r.GoMeta.ID = saved.ID
	runnableIDs := upsertGoRunnables(db, r.GoMeta)
	if err := db.DeleteStaleGoRunnables(r.GoMeta.ID, runnableIDs); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clean stale Go runnables: %v\n", err)
	}
}

// upsertGoRunnables persists all runnable files and returns their IDs.
func upsertGoRunnables(db *store.DB, meta *model.GoProjectMetadata) []int64 {
	var ids []int64
	for _, run := range meta.Runnables {
		run.GoMetadataID = meta.ID
		if err := db.UpsertGoRunnable(run); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrGoRunnableUpsert, err)

			continue
		}
		ids = append(ids, run.ID)
	}

	return ids
}

// upsertCsharpProjectMeta persists C# metadata, project files, and key files.
func upsertCsharpProjectMeta(db *store.DB, r detector.DetectionResult) {
	r.Csharp.DetectedProjectID = r.Project.ID
	if err := db.UpsertCsharpMetadata(*r.Csharp); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCsharpMetaUpsert, err)

		return
	}
	saved, err := db.SelectCsharpMetadata(r.Project.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCsharpMetaUpsert, err)

		return
	}
	r.Csharp.ID = saved.ID
	fileIDs := upsertCsharpFiles(db, r.Csharp)
	keyIDs := upsertCsharpKeyFiles(db, r.Csharp)
	if err := db.DeleteStaleCsharpFiles(r.Csharp.ID, fileIDs); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clean stale C# files: %v\n", err)
	}
	if err := db.DeleteStaleCsharpKeyFiles(r.Csharp.ID, keyIDs); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not clean stale C# key files: %v\n", err)
	}
}

// upsertCsharpFiles persists C# project files and returns their IDs.
func upsertCsharpFiles(db *store.DB, meta *model.CsharpProjectMetadata) []int64 {
	var ids []int64
	for _, f := range meta.ProjectFiles {
		f.CsharpMetadataID = meta.ID
		if err := db.UpsertCsharpProjectFile(f); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCsharpFileUpsert, err)

			continue
		}
		ids = append(ids, f.ID)
	}

	return ids
}

// upsertCsharpKeyFiles persists C# key files and returns their IDs.
func upsertCsharpKeyFiles(db *store.DB, meta *model.CsharpProjectMetadata) []int64 {
	var ids []int64
	for _, f := range meta.KeyFiles {
		f.CsharpMetadataID = meta.ID
		if err := db.UpsertCsharpKeyFile(f); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCsharpKeyUpsert, err)

			continue
		}
		ids = append(ids, f.ID)
	}

	return ids
}

// collectRepoIDs extracts unique repo IDs from detection results.
func collectRepoIDs(results []detector.DetectionResult) map[int64]bool {
	ids := make(map[int64]bool)
	for _, r := range results {
		ids[r.Project.RepoID] = true
	}

	return ids
}

// cleanStaleProjects removes projects no longer detected for each repo.
func cleanStaleProjects(db *store.DB, repoIDs map[int64]bool, results []detector.DetectionResult) {
	for repoID := range repoIDs {
		keepIDs := collectKeepIDs(repoID, results)
		cleaned, err := db.DeleteStaleProjects(repoID, keepIDs)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrProjectCleanup, repoID, err)

			continue
		}
		if cleaned > 0 {
			fmt.Printf(constants.MsgProjectCleanedStale, cleaned)
		}
	}
}

// collectKeepIDs collects project IDs to keep for a given repo.
func collectKeepIDs(repoID int64, results []detector.DetectionResult) []int64 {
	var ids []int64
	for _, r := range results {
		if r.Project.RepoID == repoID {
			ids = append(ids, r.Project.ID)
		}
	}

	return ids
}
