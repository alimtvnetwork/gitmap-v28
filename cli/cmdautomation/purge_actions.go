package cmdautomation

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type rawArtifactResponse struct {
	TotalCount int `json:"total_count"`
	Artifacts  []struct {
		Id          int64  `json:"id"`
		Name        string `json:"name"`
		SizeInBytes int64  `json:"size_in_bytes"`
		CreatedAt   string `json:"created_at"`
	} `json:"artifacts"`
}

// RunPurgeActions queries GitHub Actions API and purges obsolete artifacts.
func RunPurgeActions(opts PurgeActionsOptions) PurgeActionsResultMonad {
	start := time.Now()
	repo := resolveTargetRepo(opts.Repo)
	artifacts := fetchArtifacts(repo)
	purged := purgeArtifactsPool(repo, artifacts, opts.Workers, opts.IsDryRun)
	res := aggregatePurgeResult(repo, purged, start)
	return result.Ok(res)
}

func resolveTargetRepo(custom string) string {
	if len(custom) > 0 {
		return custom
	}
	return "alimtvnetwork/gitmap-v28"
}

func fetchArtifacts(repo string) []PurgeArtifactItem {
	endpoint := fmt.Sprintf("repos/%s/actions/artifacts?per_page=100", repo)
	cmd := exec.Command("gh", "api", endpoint)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	return parseArtifactsJson(out)
}

func parseArtifactsJson(data []byte) []PurgeArtifactItem {
	var resp rawArtifactResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil
	}
	var items []PurgeArtifactItem
	for _, a := range resp.Artifacts {
		items = append(items, PurgeArtifactItem{
			Id:        a.Id,
			Name:      a.Name,
			SizeBytes: a.SizeInBytes,
			CreatedAt: a.CreatedAt,
			IsPurged:  false,
		})
	}
	return items
}

func purgeArtifactsPool(repo string, items []PurgeArtifactItem, workers int, isDryRun bool) []PurgeArtifactItem {
	if isDryRun || len(items) == 0 {
		return items
	}
	w := resolveWorkerCount(workers, len(items))
	jobs := make(chan int, len(items))
	var wg sync.WaitGroup
	for i := 0; i < w; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				items[idx].IsPurged = deleteArtifact(repo, items[idx].Id)
			}
		}()
	}
	for i := range items {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return items
}

func deleteArtifact(repo string, id int64) bool {
	endpoint := fmt.Sprintf("repos/%s/actions/artifacts/%d", repo, id)
	cmd := exec.Command("gh", "api", "-X", "DELETE", endpoint)
	err := cmd.Run()
	return err == nil
}

func aggregatePurgeResult(repo string, items []PurgeArtifactItem, start time.Time) PurgeActionsResult {
	purgedCount := 0
	var freedBytes int64
	for _, it := range items {
		if it.IsPurged {
			purgedCount++
			freedBytes += it.SizeBytes
		}
	}
	return PurgeActionsResult{
		Repo:            repo,
		TotalArtifacts:  len(items),
		PurgedArtifacts: purgedCount,
		FreedBytes:      freedBytes,
		FreedMB:         float64(freedBytes) / (1024 * 1024),
		Artifacts:       items,
		Duration:        time.Since(start),
		IsPass:          true,
	}
}
