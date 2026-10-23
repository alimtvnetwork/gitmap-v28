# Subtask Execution Plan: LLM Training Heatmap Telemetry Linking (ST-06)

## Document Metadata
- **Subtask ID**: ST-06 (`06-subtask-llm-train-heatmap-linking.md`)
- **Parent Plan**: `.ai-memory/plans/pending/314-gitmap-heatmap-and-url-handling-instructions.md`
- **Specification Reference**:
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/00-overview.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/01-architecture-spec.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/02-component-spec.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/99-consistency-report.md`
- **Status**: `QUEUED / READY FOR PHASE 2 WAVE 2`
- **Assigned Worker**: Swarm Worker 3 (Wave 2)
- **Disjoint File Ownership**:
  - `cli/cmd/llm/llm_heatmap.go` (new)
  - `cli/cmd/llm/llm_train.go` (update)
  - `cli/cmdai/ai_train.go` (update)

---

## 1. Objectives & Architectural Intent

In autonomous AI agent development, multi-agent runs generate rich cognitive decision traces, tool call trajectories, line ranges, and outcome verifications recorded inside SQLite heatmap databases:
- **Multi-AI Companion DBs**: `repo-cache/.heatmap/<seq>-<ai>/<repo>/heatmap.db`
- **Root Telemetry DB**: `EffectiveDataDir()/heatmap.db`
- **Remote REST Endpoints**: e.g. `https://gitmap.internal/api/heatmap/<task_slug>`

When fine-tuning or training LLM models via `gitmap llm train`, `gitmap llm train muse`, or `gitmap ai analysis llm-train`, this telemetry provides direct supervisor signal:
1. Which tools succeeded or failed on specific files (`ActionType`, `StartLine`, `EndLine`).
2. Exact model reasoning traces and user prompts (`ModelName`, `Expectation`, `ActualOutcome`, `IsExpectationMet`).
3. Root cause analyses and remediation patterns.

Subtask 06 implements the ingestion bridge connecting local heatmap companion databases and remote REST endpoints directly into GitMap's LLM training pipeline.

---

## 2. Technical Contracts & Schema

### 2.1 CLI Flags Contract
The training CLI commands (`gitmap llm train`, `gitmap llm train muse`, `gitmap ai analysis llm-train`) receive additive flags:
- `--heatmap <path|url>`: Path to local SQLite heatmap database OR remote HTTP(S) URL.
- `--heatmap-repo <name>`: Optional repository filter (defaults to all).
- `--heatmap-task <slug>`: Optional task slug filter.
- `--heatmap-min-quality`: Minimum rating or verified expectation filter (`isExpectationMet = 1`).

### 2.2 Ingestion Data Structures (`cli/cmd/llm/llm_heatmap.go`)

```go
package llm

import (
	"time"
)

// HeatmapTrainingExample represents a structured training record extracted from heatmap telemetry.
type HeatmapTrainingExample struct {
	TaskSlug       string    `json:"taskSlug"`
	RepoName       string    `json:"repoName"`
	ActionType     string    `json:"actionType"`
	FilePath       string    `json:"filePath"`
	StartLine      int       `json:"startLine"`
	EndLine        int       `json:"endLine"`
	Expectation    string    `json:"expectation"`
	ActualOutcome  string    `json:"actualOutcome"`
	IsMet          bool      `json:"isMet"`
	ModelName      string    `json:"modelName"`
	PromptContext  string    `json:"promptContext"`
	CompletionText string    `json:"completionText"`
	RecordedAt     time.Time `json:"recordedAt"`
}

// HeatmapIngestionResult tracks telemetry ingestion metrics.
type HeatmapIngestionResult struct {
	TotalScannedRecords   int `json:"totalScannedRecords"`
	AcceptedTrainingPairs int `json:"acceptedTrainingPairs"`
	RejectedPairs         int `json:"rejectedPairs"`
	HasRemoteSource       bool `json:"hasRemoteSource"`
}
```

---

## 3. Implementation Steps & Micro-Batches

### Step 3.1: Create Heatmap Source Resolver (`cli/cmd/llm/llm_heatmap.go`)
- Implement `ResolveHeatmapSource(source string) ([]HeatmapTrainingExample, *apperror.AppError)`.
- Inspect whether `source` begins with `http://` or `https://`:
  - If remote: delegate to `cmdurl.FetchCachedURL(source, false)` which uses the URL Split-DB caching engine to download and store the payload. If payload is JSON, parse directly; if binary SQLite, write to temporary cache file and open.
  - If local: resolve relative file path (e.g. `repo-cache/.heatmap/...`), verify file exists, and open via `sql.Open("sqlite", ...)`.
- Ensure function lengths $\le 15$ lines using helper functions `parseRemoteHeatmapPayload` and `queryLocalHeatmapDB`.

### Step 3.2: Query & Transform Action Logs to Training Pairs
- Implement `ExtractTrainingPairs(db *sql.DB, repoFilter, taskFilter string) ([]HeatmapTrainingExample, error)`.
- Query `HeatmapActionLog` records where `IsExpectationMet = 1` or valid action annotations exist.
- Map records to prompt/completion pairs formatted for instruction tuning or DPO (Direct Preference Optimization).

### Step 3.3: Wire CLI Dispatcher & Training Driver
- Update `cli/cmd/llm/llm_train.go` to bind `--heatmap` flag.
- Pass extracted examples to training dataset assembler:
  ```go
  if hasHeatmapSource {
      examples, err := ResolveHeatmapSource(heatmapPathOrURL)
      if err != nil {
          return apperror.WrapSimple(err, "llm.resolveHeatmap")
      }
      trainingSet = append(trainingSet, convertHeatmapToCurriculum(examples)...)
  }
  ```
- Output clean terminal summary table showing ingested examples and source database tier.

---

## 4. Coding Guidelines Compliance Checklist

- [x] **Positive Booleans**: `hasRemoteSource`, `hasHeatmapSource`, `isMet`, `isValidSource`.
- [x] **Function Limits**: All functions strictly $\le 15$ lines (preferred $\le 8$ lines).
- [x] **Zero Nested Ifs**: Guard clauses with early returns on all error conditions.
- [x] **Error Sovereignty**: Leaf functions return `*apperror.AppError` via `apperror.WrapSimple(err, "<scope>")`.
- [x] **Relative Paths**: Database and telemetry paths resolved using relative references.
- [x] **Zero Git Commands**: No subagent git invocations.
- [x] **Zero Builds**: No compilation during implementation turns.

---

## 5. Verification & Acceptance Criteria

1. **Local Companion Ingestion**:
   - Running `gitmap llm train --heatmap repo-cache/.heatmap/01-cline/gitmap/heatmap.db` reads action logs and includes them in training curriculum.
2. **Remote REST Ingestion**:
   - Running `gitmap llm train --heatmap https://example.internal/heatmap.json` downloads via URL cache Split-DB and parses training records.
3. **Filtering**:
   - Flags `--heatmap-repo` and `--heatmap-task` correctly restrict ingested records.
4. **Unit Test**:
   - `cli/cmd/llm/llm_heatmap_test.go` verifies parsing logic with mock SQLite database and mock HTTP payload.
