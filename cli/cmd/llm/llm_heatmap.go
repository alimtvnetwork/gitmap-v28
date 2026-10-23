// Package llm — llm_heatmap.go: Heatmap training telemetry ingestion, source resolution, and digest formatting.
package llm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/heatmap"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// HeatmapTrainingExample represents a structured cognitive decision trace extracted from heatmap telemetry.
type HeatmapTrainingExample struct {
	TaskSlug         string `json:"taskSlug"`
	RepoSlug         string `json:"repoSlug"`
	ActionType       string `json:"actionType"`
	FilePath         string `json:"filePath"`
	ExpectedOutcome  string `json:"expectedOutcome,omitempty"`
	ActualOutcome    string `json:"actualOutcome,omitempty"`
	IsExpectationMet bool   `json:"isExpectationMet"`
	ReasoningState   string `json:"reasoningState,omitempty"`
	ModelName        string `json:"modelName,omitempty"`
	ToolName         string `json:"toolName,omitempty"`
}

// LlmHeatmapOptions configures heatmap training source resolution and ingestion.
type LlmHeatmapOptions struct {
	Source   string `json:"source"`
	RepoSlug string `json:"repoSlug"`
	TaskSlug string `json:"taskSlug"`
	Limit    int    `json:"limit"`
}

// HeatmapIngestionResult tracks telemetry ingestion metrics.
type HeatmapIngestionResult struct {
	TotalScannedRecords   int  `json:"totalScannedRecords"`
	AcceptedTrainingPairs int  `json:"acceptedTrainingPairs"`
	RejectedPairs         int  `json:"rejectedPairs"`
	HasRemoteSource       bool `json:"hasRemoteSource"`
}

const sqlSelectActionLogExamples = `SELECT 
	COALESCE(TaskSlug, ''), 
	COALESCE(RepoName, ''), 
	COALESCE(ActionType, ''), 
	COALESCE(FilePath, ''), 
	COALESCE(ExpectedOutcome, ''), 
	COALESCE(ActualOutcome, ''), 
	IsExpectationMet, 
	COALESCE(ReasoningState, ''), 
	COALESCE(ModelName, ''), 
	COALESCE(ToolName, '') 
FROM HeatmapActionLog 
ORDER BY HeatmapActionLogId DESC 
LIMIT ?`

// ResolveHeatmapTrainingSources discovers local tracking DBs, split DBs, and companion DBs, or validates explicit path/URL.
func ResolveHeatmapTrainingSources(opts LlmHeatmapOptions) ([]string, error) {
	hasExplicitSource := strings.TrimSpace(opts.Source) != ""
	if hasExplicitSource {
		return resolveExplicitHeatmapSource(strings.TrimSpace(opts.Source))
	}

	return discoverAmbientHeatmapSources(opts)
}

func resolveExplicitHeatmapSource(source string) ([]string, error) {
	isRemote := strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://")
	if isRemote {
		return []string{source}, nil
	}

	cleanPath := filepath.ToSlash(source)

	return []string{cleanPath}, nil
}

func discoverAmbientHeatmapSources(opts LlmHeatmapOptions) ([]string, error) {
	var sources []string
	sources = appendRootDbSource(sources)
	sources = appendTierSources(sources, opts)

	return sources, nil
}

func appendRootDbSource(sources []string) []string {
	rootDb := filepath.ToSlash(filepath.Join(store.EffectiveDataDir(), "heatmap.db"))
	if _, err := os.Stat(rootDb); err == nil {
		return append(sources, rootDb)
	}

	return sources
}

func appendTierSources(sources []string, opts LlmHeatmapOptions) []string {
	dbs, err := heatmap.DiscoverHeatmapDatabases()
	if err != nil {
		return sources
	}

	return filterDiscoveredSources(sources, dbs, opts.RepoSlug)
}

func filterDiscoveredSources(sources []string, dbs []heatmap.DiscoveredHeatmapDB, repoFilter string) []string {
	seen := make(map[string]bool)
	for _, s := range sources {
		seen[s] = true
	}

	for _, db := range dbs {
		if hasRepoMatch(db.RepoSlug, repoFilter) && !seen[db.Path] {
			sources = append(sources, filepath.ToSlash(db.Path))
			seen[db.Path] = true
		}
	}

	return sources
}

func hasRepoMatch(repoSlug, filter string) bool {
	hasNoFilter := strings.TrimSpace(filter) == ""
	if hasNoFilter {
		return true
	}

	return strings.EqualFold(repoSlug, strings.TrimSpace(filter))
}

// IngestHeatmapTrainingData ingests decision traces from a URL or local SQLite database.
func IngestHeatmapTrainingData(source string, limit int) ([]HeatmapTrainingExample, error) {
	effectiveLimit := resolveIngestionLimit(limit)
	isRemote := strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://")
	if isRemote {
		return ingestRemoteHeatmapData(source, effectiveLimit)
	}

	return ingestLocalHeatmapDb(source, effectiveLimit)
}

func resolveIngestionLimit(limit int) int {
	isPositive := limit > 0
	if isPositive {
		return limit
	}

	return 500
}

func ingestRemoteHeatmapData(url string, limit int) ([]HeatmapTrainingExample, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, apperror.WrapSimple(err, "llm.heatmap.get_remote")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apperror.WrapSimple(err, "llm.heatmap.read_body")
	}

	return parseHeatmapJsonPayload(body, limit)
}

func parseHeatmapJsonPayload(data []byte, limit int) ([]HeatmapTrainingExample, error) {
	var rawExamples []HeatmapTrainingExample
	if err := json.Unmarshal(data, &rawExamples); err != nil {
		return nil, apperror.WrapSimple(err, "llm.heatmap.unmarshal_json")
	}

	if len(rawExamples) > limit {
		return rawExamples[:limit], nil
	}

	return rawExamples, nil
}

func ingestLocalHeatmapDb(dbPath string, limit int) ([]HeatmapTrainingExample, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, apperror.WrapSimple(err, "llm.heatmap.stat_db")
	}

	dsn := fmt.Sprintf("file:%s?mode=ro&cache=shared", filepath.ToSlash(dbPath))
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, apperror.WrapSimple(err, "llm.heatmap.open_db")
	}
	defer conn.Close()

	return queryActionLogExamples(conn, limit)
}

func queryActionLogExamples(conn *sql.DB, limit int) ([]HeatmapTrainingExample, error) {
	rows, err := conn.Query(sqlSelectActionLogExamples, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "llm.heatmap.query_action_log")
	}
	defer rows.Close()

	return scanActionLogRows(rows)
}

func scanActionLogRows(rows *sql.Rows) ([]HeatmapTrainingExample, error) {
	var examples []HeatmapTrainingExample
	for rows.Next() {
		ex, err := scanSingleActionLogRow(rows)
		if err != nil {
			return nil, err
		}
		examples = append(examples, ex)
	}

	return examples, nil
}

func scanSingleActionLogRow(rows *sql.Rows) (HeatmapTrainingExample, error) {
	var ex HeatmapTrainingExample
	var isMetInt int
	dest := []any{&ex.TaskSlug, &ex.RepoSlug, &ex.ActionType, &ex.FilePath, &ex.ExpectedOutcome, &ex.ActualOutcome, &isMetInt, &ex.ReasoningState, &ex.ModelName, &ex.ToolName}
	if err := rows.Scan(dest...); err != nil {
		return HeatmapTrainingExample{}, apperror.WrapSimple(err, "llm.heatmap.scan_row")
	}
	ex.IsExpectationMet = isMetInt == 1

	return ex, nil
}

// PartitionHeatmapExamples partitions records into positive reinforcement and negative anti-pattern examples.
func PartitionHeatmapExamples(examples []HeatmapTrainingExample) ([]HeatmapTrainingExample, []HeatmapTrainingExample) {
	var positive []HeatmapTrainingExample
	var negative []HeatmapTrainingExample

	for _, ex := range examples {
		if ex.IsExpectationMet {
			positive = append(positive, ex)
		} else {
			negative = append(negative, ex)
		}
	}

	return positive, negative
}

// FormatHeatmapTrainingDigest formats a structured Markdown digest partitioned into positive and negative traces.
func FormatHeatmapTrainingDigest(examples []HeatmapTrainingExample) string {
	hasExamples := len(examples) > 0
	if !hasExamples {
		return "No heatmap training actions recorded.\n"
	}

	positive, negative := PartitionHeatmapExamples(examples)

	var sb strings.Builder
	writeDigestHeader(&sb, len(examples), len(positive), len(negative))
	writePositiveDigest(&sb, positive)
	writeNegativeDigest(&sb, negative)

	return sb.String()
}

func writeDigestHeader(sb *strings.Builder, total, pos, neg int) {
	sb.WriteString("======================================================================\n")
	sb.WriteString("HEATMAP COGNITIVE DECISION TRACES & TRAINING DIGEST\n")
	sb.WriteString("======================================================================\n")
	sb.WriteString(fmt.Sprintf("Total Telemetry Actions: %d | Positive Reinforcements: %d | Negative RCA: %d\n\n", total, pos, neg))
}

func writePositiveDigest(sb *strings.Builder, positive []HeatmapTrainingExample) {
	hasPositive := len(positive) > 0
	if !hasPositive {
		return
	}

	sb.WriteString("### [Positive Reinforcement Traces (Verified Successes)]\n")
	renderPositiveItems(sb, positive)
	sb.WriteString("\n")
}

func renderPositiveItems(sb *strings.Builder, positive []HeatmapTrainingExample) {
	limit := resolveDigestLimit(len(positive))
	for i := 0; i < limit; i++ {
		p := positive[i]
		sb.WriteString(fmt.Sprintf("- [%s] Task: %s | Action: %s | Tool: %s | File: %s\n", p.ReasoningState, p.TaskSlug, p.ActionType, p.ToolName, p.FilePath))
		if p.ExpectedOutcome != "" {
			sb.WriteString(fmt.Sprintf("  Expected: %s -> Actual: %s\n", p.ExpectedOutcome, p.ActualOutcome))
		}
	}
}

func writeNegativeDigest(sb *strings.Builder, negative []HeatmapTrainingExample) {
	hasNegative := len(negative) > 0
	if !hasNegative {
		return
	}

	sb.WriteString("### [Negative Anti-Pattern Traces (RCA Prevention)]\n")
	renderNegativeItems(sb, negative)
	sb.WriteString("\n")
}

func renderNegativeItems(sb *strings.Builder, negative []HeatmapTrainingExample) {
	limit := resolveDigestLimit(len(negative))
	for i := 0; i < limit; i++ {
		n := negative[i]
		sb.WriteString(fmt.Sprintf("- [%s] Task: %s | Action: %s | Tool: %s | File: %s\n", n.ReasoningState, n.TaskSlug, n.ActionType, n.ToolName, n.FilePath))
		if n.ExpectedOutcome != "" {
			sb.WriteString(fmt.Sprintf("  Divergence: Expected %q, got %q\n", n.ExpectedOutcome, n.ActualOutcome))
		}
	}
}

func resolveDigestLimit(count int) int {
	if count < 10 {
		return count
	}

	return 10
}

// FormatHeatmapTrainingJsonl formats examples as newline-delimited JSON.
func FormatHeatmapTrainingJsonl(examples []HeatmapTrainingExample) (string, error) {
	var sb strings.Builder
	for _, ex := range examples {
		bytes, err := json.Marshal(ex)
		if err != nil {
			return "", apperror.WrapSimple(err, "llm.heatmap.marshal_jsonl")
		}
		sb.Write(bytes)
		sb.WriteByte('\n')
	}

	return sb.String(), nil
}
