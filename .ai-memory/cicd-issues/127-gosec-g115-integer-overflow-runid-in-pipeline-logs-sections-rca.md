# RCA-127: Gosec G115 Integer Overflow uint64 -> int64 in Pipeline Logs Sections

## Part 1: Symptom
CI/CD Pipeline failed on workflow `CI` (Run `#40058123270`), Job `Lint Baseline Guard (unused, gosec G115, misspell, gocritic, exhaustive)`, Step `Gosec G115 diff (baseline-diff, full-path only)`.
Error log output:
```
[CI/CD ERROR REPORT] New gosec-G115 Violations Detected (1 finding(s)):
[gosec-G115] G115: integer overflow conversion uint64 -> int64 (NEW vs baseline)
┌─ Error [1/1]: [gosec-G115]
│ File:     cmdpipeline/pipeline_logs_sections.go:65:41
│ Location: cmdpipeline/pipeline_logs_sections.go at line 65, col 41
│ Linter:   gosec-G115
│ Message:  G115: integer overflow conversion uint64 -> int64
│ Context:  NEW finding vs baseline
│ Script:   .github/scripts/check-single-linter-diff.py
└──────────────────────────────────────────────────────────
❌ cmdpipeline/pipeline_logs_sections.go:65:41: [gosec-G115] G115: integer overflow conversion uint64 -> int64
```

## Part 2: Root Cause
In `cli/cmdpipeline/pipeline_logs_sections.go`, helper function `resolveSectionUrl` declared its second parameter as `runId int64`:
```go
func resolveSectionUrl(url string, runId int64) string
```
However, in `SectionFailure` (defined in `cli/cmdpipeline/pipeline.go`), `RunId` is defined as `uint64`. Line 65 performed a direct type conversion:
```go
url := resolveSectionUrl(sec.Url, int64(sec.RunId))
```
The Gosec G115 rule flags integer conversion from `uint64` to `int64` as an potential integer overflow vulnerability because values above `math.MaxInt64` wrap to negative numbers.

## Part 3: Surgical Fix
Updated `resolveSectionUrl` in `cli/cmdpipeline/pipeline_logs_sections.go` to accept `runId uint64` natively matching `SectionFailure.RunId`:
1. Changed `func resolveSectionUrl(url string, runId int64) string` to `func resolveSectionUrl(url string, runId uint64) string`.
2. Changed the guard `if runId <= 0` to `if runId == 0`.
3. Eliminated the type conversion on line 65: `url := resolveSectionUrl(sec.Url, sec.RunId)`.

## Part 4: Verification
1. Verified Go compilation of `cli`: `go build -o ../bin/gitmap.exe .` exited 0 with no errors.
2. Verified guideline compliance: `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdpipeline/pipeline_logs_sections.go --check-only` passed.
3. Successfully resolved Gosec G115 violation at line 65 col 41 with zero remaining type casts.
