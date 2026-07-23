// Package cmd — scanprojects.go handles project detection during scan.
package cmd

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/detector"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// detectionWorkerCap bounds concurrent per-repo detection walks. Detection
// is I/O bound (lots of os.ReadDir / os.Stat); 8 workers saturate a typical
// SSD without exhausting file descriptors on Windows where the default
// per-process handle budget is tighter than POSIX.
const detectionWorkerCap = 8

// detectAllProjects runs project detection across all scanned repos in
// parallel. Each repo is walked exactly once by the detector (see
// detector.DetectProjects); we use a small worker pool so 24 repos don't
// take 24× the wall clock of a single repo.
func detectAllProjects(records []model.ScanRecord) []detector.DetectionResult {
	workers := resolveDetectionWorkers(len(records))
	jobs := make(chan model.ScanRecord)
	resultsCh := make(chan []detector.DetectionResult, len(records))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go detectionWorker(jobs, resultsCh, &wg)
	}
	for _, rec := range records {
		jobs <- rec
	}
	close(jobs)
	wg.Wait()
	close(resultsCh)

	all, repoCount := collectDetectionResults(resultsCh)
	fmt.Printf(constants.MsgProjectDetectDone, len(all), repoCount)

	return all
}

// detectionWorker pulls repos off jobs and emits results onto resultsCh.
// Each repo's detector.DetectProjects call is independent, so the only
// shared state is the channel buffer.
func detectionWorker(jobs <-chan model.ScanRecord, resultsCh chan<- []detector.DetectionResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for rec := range jobs {
		resultsCh <- detector.DetectProjects(rec.AbsolutePath, rec.ID, rec.RepoName)
	}
}

// collectDetectionResults flattens per-repo results and counts repos that
// contributed at least one detected project.
func collectDetectionResults(resultsCh <-chan []detector.DetectionResult) ([]detector.DetectionResult, int) {
	var all []detector.DetectionResult
	repoCount := 0
	for results := range resultsCh {
		if len(results) > 0 {
			repoCount++
			all = append(all, results...)
		}
	}

	return all, repoCount
}

// resolveDetectionWorkers picks a worker-pool size bounded by repo count,
// CPU count, and the hard cap. With 1 repo there's no point spinning up 8
// goroutines; with 100 repos on a 4-core box the I/O bound pool can still
// usefully run 8 in flight.
func resolveDetectionWorkers(repoCount int) int {
	if repoCount <= 1 {
		return 1
	}
	n := runtime.NumCPU()
	if n < 1 {
		n = 1
	}
	if n > detectionWorkerCap {
		n = detectionWorkerCap
	}
	if n > repoCount {
		n = repoCount
	}

	return n
}

// upsertProjectsToDB persists detected projects and metadata to SQLite.
func upsertProjectsToDB(results []detector.DetectionResult, records []model.ScanRecord, outputDir string) {
	if len(results) == 0 {
		return
	}
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

		return
	}
	upsertProjectRecords(db, results, records)
}

// upsertProjectRecords inserts projects, metadata, and cleans stale records.
func upsertProjectRecords(db *store.DB, results []detector.DetectionResult, records []model.ScanRecord) {
	count := 0
	repoIDs := collectRepoIDs(results)
	for i := range results {
		r := &results[i]
		err := db.UpsertDetectedProject(r.Project)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

			continue
		}
		if err := resolveDetectedProjectID(db, r); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrProjectUpsert, err)

			continue
		}
		count++
		upsertProjectMetadata(db, *r)
	}
	cleanStaleProjects(db, repoIDs, results)
	fmt.Printf(constants.MsgProjectUpsertDone, count)
}

// resolveDetectedProjectID syncs the project ID with the persisted DB row.
func resolveDetectedProjectID(db *store.DB, r *detector.DetectionResult) error {
	id, err := db.SelectDetectedProjectID(
		r.Project.RepoID,
		r.Project.ProjectTypeID,
		r.Project.RelativePath,
	)
	if err != nil {
		return err
	}
	r.Project.ID = id

	return nil
}

// upsertProjectMetadata persists Go or C# metadata for a detection result.
func upsertProjectMetadata(db *store.DB, r detector.DetectionResult) {
	if r.GoMeta != nil {
		upsertGoProjectMeta(db, r)
	}
	if r.Csharp != nil {
		upsertCsharpProjectMeta(db, r)
	}
}
