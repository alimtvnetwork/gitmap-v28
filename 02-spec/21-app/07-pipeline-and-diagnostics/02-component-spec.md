# 07-pipeline-and-diagnostics: PE Extractor Components, Heatmap CLI & ETA Specification

- **Spec ID:** `07-pipeline-and-diagnostics/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Pipeline Error Extractor, Heatmap Renderer, ETA Predictor, RCA Generator
- **Dependencies:** `cli/cmdpipeline`, `cli/pipelinedb`, `cli/termtable`, `cli/termpad`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The Pipeline and Diagnostics cluster comprises four primary operational components:

```
07-pipeline-and-diagnostics/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Pipeline Error Extractor** | `cli/cmdpipeline/pipeline_error_extract.go`, `cli/cmdpipeline/pipeline_stacktrace.go` | Log analysis, regex filtering, polyglot traceback recognition. |
| **Traceback Heatmap UI** | `cli/cmdpipeline/pipeline_details.go`, `cli/cmdpipeline/pipeline_ai_errors.go` | Terminal UI cards, ANSI color stages, traceback highlights. |
| **ETA Forecaster** | `cli/cmdpipeline/pipeline_eta.go`, `cli/pipelinedb/timing.go` | Historical run queries, moving-average computation, dynamic countdowns. |
| **4-Part RCA Reporter** | `cli/cmdpipeline/pipeline_rca.go` | Automated synthesis of markdown RCA documents and fix recommendations. |

---

## 2. Polyglot Traceback Extraction Engine

### 2.1 Stack Trace Start Markers
The parser identifies error boundaries using language-specific start markers:
- **Python:** `Traceback (most recent call last):`, `FAIL: test_`, `ERROR: test_`.
- **Go:** `panic:`, `goroutine \d+ \[running\]:`, `--- FAIL: Test`.
- **Node.js:** `Error:`, `TypeError:`, `AssertionError:`, `at Object.<anonymous>`.

### 2.2 Stack Frame Extraction Heuristics
```go
package cmdpipeline

func ExtractStackTrace(logLines []string) []string {
    var frames []string
    inTrace := false
    for _, line := range logLines {
        if isTraceStart(line) {
            inTrace = true
        }
        if inTrace {
            if isTraceTerminator(line) && len(frames) > 0 {
                break
            }
            if !isProgressNoise(line) {
                frames = append(frames, line)
            }
        }
    }
    return frames
}
```

---

## 3. Pipeline Error CLI Command Suite

### 3.1 Commands
- `gitmap pe` (`gitmap pipeline-error`): Displays latest failed workflow run with traceback cards and suggested fixes.
- `gitmap pea` (`gitmap pipeline-error --all`): Displays full failure matrix across all active runs.
- `gitmap pe --heatmap`: Renders color-coded terminal heatmap of workflow stages with execution duration bars.
- `gitmap pe --rca`: Emits complete 4-part RCA markdown report to console or clipboard.

---

## 4. Accurate ETA Forecasting Subsystem

The ETA engine estimates time remaining for in-flight workflows:
1. Queries historical timings from `pipelinedb.StageTiming` for the given workflow and branch.
2. Calculates trimmed mean (excluding outliers $>2\sigma$).
3. Renders live countdown timer in CLI:
   `Running: Step 3/5 [Linting & Tests] — ETA: 42s remaining (Avg: 2m 15s)`

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isTracebackExtractionTested: true
  isHeatmapCardFormattingVerified: true
  isEtaComputationWithinTolerance: true
  isPositiveBooleansUsed: true
```

- [x] Traceback extraction isolated from test runner execution summaries.
- [x] Heatmap cards present file/line coordinates clearly.
- [x] ETA forecasts within $\pm 15\%$ of actual workflow completion.
