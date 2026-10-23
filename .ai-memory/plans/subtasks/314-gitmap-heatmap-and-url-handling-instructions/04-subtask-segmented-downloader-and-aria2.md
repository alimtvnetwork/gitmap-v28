# Subtask 04: 10-Split Segmented & Aria2 Multi-Tier Downloader

> **Subtask ID:** `SUBTASK-314-04`  
> **Parent Task:** `314-gitmap-heatmap-and-url-handling-instructions`  
> **Status:** Pending  
> **Target Files:**  
> - `cli/downloader/segmented.go` (new file)  
> - `cli/downloader/aria2.go` (new file)  
> - `cli/downloader/download.go` (new file)  
> **Related Architecture Spec:** `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/02-component-spec.md`  

---

## 1. Objective

Implement the high-throughput multi-tier URL downloader in package `cli/downloader`. Payloads with size >= 1–2 MB dynamically activate 10-connection parallel acceleration. If host `aria2c` is present, Tier 1 executes `aria2c -s 10 -x 10 -k 5K`. If absent, Tier 2 coordinates 10 concurrent HTTP Range workers writing to a pre-allocated file using `io.NewOffsetWriter`. Tier 3 provides a resilient fallback to single-stream HTTP GET for non-range servers or small files (< 1MB).

---

## 2. Action Checklist

- [ ] 1. Author `cli/downloader/aria2.go` (Tier 1):
  - [ ] Implement `HasAria2c() bool` via `exec.LookPath("aria2c")`.
  - [ ] Implement `DownloadWithAria2c(opts DownloadOptions) (DownloadResult, error)`.
  - [ ] Configure aria2c parameters: `-s 10 -x 10 -k 5K --file-allocation=none --allow-overwrite=true`.
  - [ ] Stream progress/telemetry to `output.UIErr()` while suppressing console pollution.
  - [ ] Guard against 0-byte or corrupted outputs on failure.

- [ ] 2. Author `cli/downloader/segmented.go` (Tier 2):
  - [ ] Implement metadata probe `probeRemoteResource(urlText string) (resourceMeta, error)`.
  - [ ] Check HTTP headers: `Accept-Ranges: bytes` and parse `Content-Length`.
  - [ ] Implement chunk range math:
    - [ ] Calculate 10 disjoint byte ranges `[start, end]`.
  - [ ] Implement pre-allocation: create target file and expand via `f.Truncate(contentLength)`.
  - [ ] Implement parallel worker coordinator:
    - [ ] Spawn 10 concurrent goroutines using `sync.WaitGroup` and error channels.
    - [ ] Each worker issues HTTP GET with header `Range: bytes={start}-{end}`.
    - [ ] Verify HTTP response status `206 Partial Content`.
    - [ ] Write directly into `io.NewOffsetWriter(targetFile, startOffset)`.
  - [ ] Worker retry loop: up to 3 retries per chunk on transient network drops.

- [ ] 3. Author `cli/downloader/download.go` (Tier 3 & Orchestrator):
  - [ ] Implement `DownloadSingleStream(opts DownloadOptions) (DownloadResult, error)` for Tier 3.
  - [ ] Implement top-level orchestrator `Download(opts DownloadOptions) (DownloadResult, error)`:
    - [ ] Evaluate remote length and range support.
    - [ ] If length >= 1MB and range supported:
      - [ ] If `HasAria2c()`, try Tier 1.
      - [ ] If Tier 1 fails or `aria2c` missing, execute Tier 2 (Go 10-Split Range).
    - [ ] If range not supported or size < 1MB, route to Tier 3 (Single Stream).
  - [ ] Wrap errors with `*apperror.AppError` via `apperror.WrapSimple(err, "downloader.<func>")`.
  - [ ] Return structured `DownloadResult` with duration, tier used, and byte counts.

- [ ] 4. Quality & Coding Guideline Compliance:
  - [ ] Function size limit: all functions <= 15 lines (preferred <= 8 lines).
  - [ ] Zero nested `if` statements. Early return guard clauses only.
  - [ ] Positive booleans only: `HasAria2c`, `HasRangeSupport`, `IsChunkCompleted`, `IsQuiet`, `IsJSON`.
  - [ ] Blank line after `}` and blank line before `return`.
  - [ ] Zero git commands, zero builds, zero tests.

---

## 3. Concrete Code Specifications

### Specification: `cli/downloader/aria2.go`

```go
package downloader

import (
	"os/exec"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// HasAria2c checks if aria2c binary exists in the system PATH.
func HasAria2c() bool {
	_, err := exec.LookPath("aria2c")

	return err == nil
}

// DownloadWithAria2c executes Tier 1 aria2c multi-connection download.
func DownloadWithAria2c(opts DownloadOptions) (DownloadResult, error) {
	start := time.Now()
	args := buildAria2cArgs(opts)
	cmd := exec.Command("aria2c", args...)
	if err := cmd.Run(); err != nil {
		return DownloadResult{}, apperror.WrapSimple(err, "downloader.DownloadWithAria2c")
	}

	return buildCompletedResult(opts, Tier1Aria2c, time.Since(start))
}

func buildAria2cArgs(opts DownloadOptions) []string {
	return []string{
		"--disable-ipv6=true",
		"-x", "10",
		"-s", "10",
		"-k", "5K",
		"--file-allocation=none",
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		"--console-log-level=error",
		"-o", opts.DestinationPath,
		opts.URL,
	}
}
```

### Specification: `cli/downloader/segmented.go`

```go
package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RangeChunk represents a discrete byte segment for parallel transfer.
type RangeChunk struct {
	Index       int
	StartOffset int64
	EndOffset   int64
}

// DownloadSegmented coordinates Tier 2 10-connection parallel Range download.
func DownloadSegmented(opts DownloadOptions, totalSize int64) (DownloadResult, error) {
	start := time.Now()
	targetFile, err := preparePreallocatedFile(opts.DestinationPath, totalSize)
	if err != nil {
		return DownloadResult{}, err
	}
	defer targetFile.Close()

	if err := executeChunkPool(opts.URL, targetFile, totalSize); err != nil {
		return DownloadResult{}, err
	}

	return buildCompletedResult(opts, Tier2NativeRange, time.Since(start))
}

func preparePreallocatedFile(path string, size int64) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "downloader.preparePreallocatedFile.create")
	}

	if err := f.Truncate(size); err != nil {
		_ = f.Close()
		return nil, apperror.WrapSimple(err, "downloader.preparePreallocatedFile.truncate")
	}

	return f, nil
}

func executeChunkPool(urlText string, file *os.File, totalSize int64) error {
	chunks := calculateTenChunks(totalSize)
	var wg sync.WaitGroup
	errChan := make(chan error, len(chunks))

	for _, chunk := range chunks {
		wg.Add(1)
		go runChunkWorker(urlText, file, chunk, &wg, errChan)
	}

	wg.Wait()
	close(errChan)

	return readFirstChunkError(errChan)
}

func calculateTenChunks(totalSize int64) []RangeChunk {
	chunks := make([]RangeChunk, 10)
	chunkSize := (totalSize + 9) / 10
	for i := 0; i < 10; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if end >= totalSize {
			end = totalSize - 1
		}
		chunks[i] = RangeChunk{Index: i, StartOffset: start, EndOffset: end}
	}

	return chunks
}

func runChunkWorker(urlText string, file *os.File, chunk RangeChunk, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlText, nil)
	if err != nil {
		errChan <- apperror.WrapSimple(err, "downloader.runChunkWorker.req")
		return
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", chunk.StartOffset, chunk.EndOffset))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errChan <- apperror.WrapSimple(err, "downloader.runChunkWorker.do")
		return
	}
	defer resp.Body.Close()

	writer := io.NewOffsetWriter(file, chunk.StartOffset)
	if _, err := io.Copy(writer, resp.Body); err != nil {
		errChan <- apperror.WrapSimple(err, "downloader.runChunkWorker.copy")
	}
}

func readFirstChunkError(errChan <-chan error) error {
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
```

### Specification: `cli/downloader/download.go`

```go
package downloader

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const SplitSizeThreshold int64 = 1024 * 1024 // 1 MB

// Download orchestrates multi-tier URL downloads with transparent fallback.
func Download(opts DownloadOptions) (DownloadResult, error) {
	meta, err := probeRemoteResource(opts.URL)
	if err != nil {
		return DownloadSingleStream(opts)
	}

	hasEligibleRange := meta.HasRangeSupport && meta.ContentLength >= SplitSizeThreshold
	if !hasEligibleRange {
		return DownloadSingleStream(opts)
	}

	if HasAria2c() {
		res, err := DownloadWithAria2c(opts)
		if err == nil {
			return res, nil
		}
	}

	res, err := DownloadSegmented(opts, meta.ContentLength)
	if err == nil {
		return res, nil
	}

	return DownloadSingleStream(opts)
}

// DownloadSingleStream executes standard single-stream HTTP GET download (Tier 3).
func DownloadSingleStream(opts DownloadOptions) (DownloadResult, error) {
	start := time.Now()
	resp, err := http.Get(opts.URL)
	if err != nil {
		return DownloadResult{}, apperror.WrapSimple(err, "downloader.DownloadSingleStream")
	}
	defer resp.Body.Close()

	out, err := os.Create(opts.DestinationPath)
	if err != nil {
		return DownloadResult{}, apperror.WrapSimple(err, "downloader.DownloadSingleStream.create")
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return DownloadResult{}, apperror.WrapSimple(err, "downloader.DownloadSingleStream.copy")
	}

	return buildCompletedResult(opts, Tier3SingleStream, time.Since(start))
}
```

---

## 4. Verification & Done Criteria

1. Dynamic 1MB threshold automatically directs downloads to Tier 1 or Tier 2.
2. `HasAria2c()` safely identifies host `aria2c` without panicking.
3. Native Range downloader divides payloads into 10 discrete segments and writes via `io.NewOffsetWriter`.
4. Resilient cascading fallback to Tier 3 when Range headers are ignored.
5. All functions <= 15 lines with guard clauses and positive booleans.
