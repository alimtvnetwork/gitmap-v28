package cmdpipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// PipelineEtaHistoryRecord captures timing telemetry for an individual completed run.
type PipelineEtaHistoryRecord struct {
	RunId           uint64 `json:"runId"`
	Workflow        string `json:"workflow"`
	Sha             string `json:"sha"`
	DurationSeconds int    `json:"durationSeconds"`
	Conclusion      string `json:"conclusion"`
	CreatedAt       string `json:"createdAt"`
}

// PipelineEtaCache encapsulates rolling workflow timing history and moving averages.
type PipelineEtaCache struct {
	Repo      string                     `json:"repo"`
	UpdatedAt string                     `json:"updatedAt"`
	History   []PipelineEtaHistoryRecord `json:"history"`
	Averages  map[string]int             `json:"averages"`
}

// ResolveEtaCacheFilePath returns the path to eta_history.json for the target repo.
func ResolveEtaCacheFilePath(repo string) string {
	dir := resolvePipelineDirForRepo(repo)
	_ = os.MkdirAll(dir, 0755)

	return filepath.Join(dir, "eta_history.json")
}

// LoadPipelineEtaCache reads the persisted rolling ETA cache from disk.
func LoadPipelineEtaCache(repo string) PipelineEtaCache {
	filePath := ResolveEtaCacheFilePath(repo)
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0 {
		return initEmptyEtaCache(repo)
	}

	var cache PipelineEtaCache
	if jsonErr := json.Unmarshal(data, &cache); jsonErr != nil {
		return initEmptyEtaCache(repo)
	}
	ensureAveragesMap(&cache)

	return cache
}

func initEmptyEtaCache(repo string) PipelineEtaCache {
	return PipelineEtaCache{
		Repo:      repo,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		History:   make([]PipelineEtaHistoryRecord, 0),
		Averages:  make(map[string]int),
	}
}

func ensureAveragesMap(cache *PipelineEtaCache) {
	if cache.Averages == nil {
		cache.Averages = make(map[string]int)
	}
}

// SavePipelineEtaCache writes the rolling ETA cache to disk.
func SavePipelineEtaCache(cache PipelineEtaCache) error {
	filePath := ResolveEtaCacheFilePath(cache.Repo)
	cache.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// UpdatePipelineEtaCache ingest runs, recalculates moving averages, and persists to disk.
func UpdatePipelineEtaCache(repo string, runs []ghRunItem) PipelineEtaCache {
	cache := LoadPipelineEtaCache(repo)
	if len(runs) == 0 {
		return cache
	}

	mergeRunsIntoHistory(&cache, runs)
	sortHistoryRecordsDesc(&cache)
	capHistoryRecords(&cache, 20)
	recomputeMovingAverages(&cache)
	_ = SavePipelineEtaCache(cache)

	return cache
}

func mergeRunsIntoHistory(cache *PipelineEtaCache, runs []ghRunItem) {
	existingMap := indexExistingHistory(cache.History)
	for _, r := range runs {
		if !isCompletedRunWithDuration(r) {
			continue
		}
		record := convertRunToEtaRecord(r)
		existingMap[record.RunId] = record
	}
	cache.History = flattenHistoryMap(existingMap)
}

func isCompletedRunWithDuration(r ghRunItem) bool {
	if r.Status != "completed" {
		return false
	}
	dur := calculateRunDuration(r.CreatedAt, r.UpdatedAt)

	return dur >= 15
}

func convertRunToEtaRecord(r ghRunItem) PipelineEtaHistoryRecord {
	return PipelineEtaHistoryRecord{
		RunId:           r.DatabaseId,
		Workflow:        r.Name,
		Sha:             r.HeadSha,
		DurationSeconds: calculateRunDuration(r.CreatedAt, r.UpdatedAt),
		Conclusion:      r.Conclusion,
		CreatedAt:       r.CreatedAt,
	}
}

func indexExistingHistory(history []PipelineEtaHistoryRecord) map[uint64]PipelineEtaHistoryRecord {
	indexed := make(map[uint64]PipelineEtaHistoryRecord, len(history))
	for _, h := range history {
		indexed[h.RunId] = h
	}

	return indexed
}

func flattenHistoryMap(indexed map[uint64]PipelineEtaHistoryRecord) []PipelineEtaHistoryRecord {
	records := make([]PipelineEtaHistoryRecord, 0, len(indexed))
	for _, rec := range indexed {
		records = append(records, rec)
	}

	return records
}

func sortHistoryRecordsDesc(cache *PipelineEtaCache) {
	sort.SliceStable(cache.History, func(i, j int) bool {
		return cache.History[i].CreatedAt > cache.History[j].CreatedAt
	})
}

func capHistoryRecords(cache *PipelineEtaCache, limit int) {
	if len(cache.History) > limit {
		cache.History = cache.History[:limit]
	}
}

func recomputeMovingAverages(cache *PipelineEtaCache) {
	grouped := groupHistoryByWorkflow(cache.History)
	for wf, records := range grouped {
		cache.Averages[wf] = calculateWorkflowMovingAverage(records)
	}
}

func groupHistoryByWorkflow(records []PipelineEtaHistoryRecord) map[string][]PipelineEtaHistoryRecord {
	grouped := make(map[string][]PipelineEtaHistoryRecord)
	for _, rec := range records {
		grouped[rec.Workflow] = append(grouped[rec.Workflow], rec)
	}

	return grouped
}

func calculateWorkflowMovingAverage(records []PipelineEtaHistoryRecord) int {
	successDurs := collectSuccessDurations(records)
	if len(successDurs) > 0 {
		return computeBaselineDuration(successDurs)
	}

	allValidDurs := collectAllValidDurations(records)
	if len(allValidDurs) > 0 {
		return computeBaselineDuration(allValidDurs)
	}

	return 0
}

func collectSuccessDurations(records []PipelineEtaHistoryRecord) []int {
	var durs []int
	for _, r := range records {
		if r.Conclusion == "success" && r.DurationSeconds >= 15 {
			durs = append(durs, r.DurationSeconds)
		}
	}

	return durs
}

func collectAllValidDurations(records []PipelineEtaHistoryRecord) []int {
	var durs []int
	for _, r := range records {
		if r.DurationSeconds >= 45 {
			durs = append(durs, r.DurationSeconds)
		}
	}

	return durs
}

// GetCachedWorkflowETA returns the rolling moving average duration for a workflow.
func GetCachedWorkflowETA(repo, workflowName string) int {
	cache := LoadPipelineEtaCache(repo)
	if avg, hasAvg := cache.Averages[workflowName]; hasAvg && avg > 0 {
		return avg
	}

	for wf, avg := range cache.Averages {
		if strings.EqualFold(wf, workflowName) && avg > 0 {
			return avg
		}
	}

	return 0
}
