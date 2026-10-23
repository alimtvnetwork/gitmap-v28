package cmdupdate

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// Parallel segmented download for release archives (spec 311-04).
//
// downloadHTTP fetches an archive in a single stream. On links where the
// release CDN or the last mile caps throughput per connection, splitting
// the archive into Range segments fetched concurrently is faster; the
// installer's aria2c path exists for the same reason. On proxied or
// per-connection-throttled paths, however, parallel connections measure
// SLOWER than one stream (see the benchmark table in spec 311): on this
// repo's own network path, 4 segments took a median 9.6s vs 3.3s single
// for the 17.6 MiB v6.590.0 archive. The segmented path is therefore
// OPT-IN via GITMAP_DOWNLOAD_SEGMENTS=1, with a range-support probe
// before use and the original single-stream downloadHTTP as the default
// and as the fallback for every doubtful case.

const (
	downloadSegmentCount    = 4
	downloadSegmentMinBytes = 4 << 20 // 4 MiB: below this, one stream is faster than the setup cost.
	downloadEngineSegmented = "go_http_segmented"
	downloadEngineSingle    = "go_http"
)

// segmentRange is one inclusive byte range [start, end] of the archive.
type segmentRange struct {
	start int64
	end   int64
}

// downloadArchiveFast downloads url to dstPath. The segmented path is
// used only when it is explicitly enabled via GITMAP_DOWNLOAD_SEGMENTS,
// the archive is large enough, and the server supports ranges. It
// returns the engine name actually used, for update telemetry.
func downloadArchiveFast(url, dstPath string, totalSize int64, onProg func(w, t int64)) (string, error) {
	if os.Getenv(constants.EnvDownloadSegments) == "" {
		return downloadEngineSingle, downloadHTTP(url, dstPath, onProg)
	}
	if totalSize < downloadSegmentMinBytes {
		return downloadEngineSingle, downloadHTTP(url, dstPath, onProg)
	}
	if !probeRangeSupport(url) {
		return downloadEngineSingle, downloadHTTP(url, dstPath, onProg)
	}
	if err := fetchSegments(url, dstPath, totalSize, onProg); err != nil {
		return downloadEngineSegmented, err
	}

	return downloadEngineSegmented, nil
}

// probeRangeSupport reports whether the server answers a 1-byte range
// request with 206 Partial Content. Redirects are followed by the default
// client, which preserves the Range header for GET requests.
func probeRangeSupport(url string) bool {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("Range", "bytes=0-0")
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // URL is a GitHub release asset from release metadata.
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 16))

	return resp.StatusCode == http.StatusPartialContent
}

// splitSegmentRanges partitions total bytes into count contiguous ranges.
func splitSegmentRanges(total int64, count int) []segmentRange {
	if count < 1 {
		count = 1
	}

	chunk := total / int64(count)
	segments := make([]segmentRange, 0, count)
	for i := range count {
		start := int64(i) * chunk
		end := start + chunk - 1
		if i == count-1 {
			end = total - 1
		}
		segments = append(segments, segmentRange{start: start, end: end})
	}

	return segments
}

// fetchSegments downloads all ranges concurrently into one pre-sized
// output file. Segments write to disjoint offsets via WriteAt, which is
// safe for concurrent use on an *os.File.
func fetchSegments(url, dstPath string, totalSize int64, onProg func(w, t int64)) error {
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()
	if err := out.Truncate(totalSize); err != nil {
		return err
	}

	agg := &progressAggregator{total: totalSize, onProg: onProg}
	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error
	for _, seg := range splitSegmentRanges(totalSize, downloadSegmentCount) {
		wg.Add(1)
		go func(seg segmentRange) {
			defer wg.Done()
			if errSeg := fetchSegment(url, out, seg, agg); errSeg != nil {
				once.Do(func() { firstErr = errSeg })
			}
		}(seg)
	}
	wg.Wait()

	return firstErr
}

// fetchSegment downloads one byte range and writes it at its offset.
func fetchSegment(url string, out *os.File, seg segmentRange, agg *progressAggregator) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", seg.start, seg.end))
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // URL is a GitHub release asset from release metadata.
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusPartialContent {
		return apperror.NewExecutionError(
			fmt.Sprintf("HTTP %d from %s (range %d-%d)", resp.StatusCode, url, seg.start, seg.end))
	}

	_, err = io.Copy(countingWriter{w: io.NewOffsetWriter(out, seg.start), agg: agg}, resp.Body)

	return err
}

// progressAggregator merges per-segment byte counts into throttled
// progress callbacks, mirroring the 100ms cadence of ProgressReader.
type progressAggregator struct {
	mu         sync.Mutex
	written    int64
	total      int64
	lastUpdate time.Time
	onProg     func(w, t int64)
}

func (a *progressAggregator) add(n int64) {
	if a.onProg == nil {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.written += n
	now := time.Now()
	if now.Sub(a.lastUpdate) > 100*time.Millisecond {
		a.lastUpdate = now
		a.onProg(a.written, a.total)
	}
}

// countingWriter forwards writes and reports their size to an aggregator.
type countingWriter struct {
	w   io.Writer
	agg *progressAggregator
}

func (c countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.agg.add(int64(n))

	return n, err
}
