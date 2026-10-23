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

const (
	MinSegmentThreshold int64 = 1024 * 1024
	DefaultSegmentCount int   = 10
	MinChunkSize        int64 = 5 * 1024
)

// SegmentRange represents a discrete byte segment for parallel transfer.
type SegmentRange struct {
	Start       int64
	End         int64
	StartOffset int64
	EndOffset   int64
	Index       int
}

// RangeChunk is an alias for SegmentRange.
type RangeChunk = SegmentRange

func maxSegmentsFromSize(totalSize int64) int {
	maxSegments := int(totalSize / MinChunkSize)
	if maxSegments < 1 {
		return 1
	}

	return maxSegments
}

func normalizeSegmentCount(totalSize int64, targetSegments int) int {
	maxSegs := maxSegmentsFromSize(totalSize)
	if targetSegments < DefaultSegmentCount {
		targetSegments = DefaultSegmentCount
	}

	if targetSegments > maxSegs {
		return maxSegs
	}

	return targetSegments
}

func makeSegment(index int, start, end int64) SegmentRange {
	return SegmentRange{
		Start:       start,
		End:         end,
		StartOffset: start,
		EndOffset:   end,
		Index:       index,
	}
}

// CalculateSegments partitions totalSize into contiguous SegmentRanges.
func CalculateSegments(totalSize int64, targetSegments int) []SegmentRange {
	if totalSize <= 0 {
		return nil
	}

	n := normalizeSegmentCount(totalSize, targetSegments)
	segments := make([]SegmentRange, n)
	for i := 0; i < n; i++ {
		start := int64(i) * totalSize / int64(n)
		end := int64(i+1)*totalSize/int64(n) - 1
		segments[i] = makeSegment(i, start, end)
	}

	return segments
}

func preparePreallocatedFile(dstPath string, totalSize int64) (*os.File, error) {
	file, err := os.Create(dstPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "downloader.preparePreallocatedFile.create")
	}

	if err := file.Truncate(totalSize); err != nil {
		_ = file.Close()

		return nil, apperror.WrapSimple(err, "downloader.preparePreallocatedFile.truncate")
	}

	return file, nil
}

func fetchSegmentResponse(urlText string, seg SegmentRange) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlText, nil)
	if err != nil {
		return nil, apperror.WrapSimple(err, "downloader.fetchSegmentResponse.req")
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", seg.Start, seg.End))

	return http.DefaultClient.Do(req)
}

func writeSegmentBody(file *os.File, seg SegmentRange, resp *http.Response) error {
	defer resp.Body.Close()
	writer := io.NewOffsetWriter(file, seg.Start)
	_, err := io.Copy(writer, resp.Body)
	if err != nil {
		return apperror.WrapSimple(err, "downloader.writeSegmentBody.copy")
	}

	return nil
}

func downloadSegmentChunk(urlText string, file *os.File, seg SegmentRange) error {
	resp, err := fetchSegmentResponse(urlText, seg)
	if err != nil {
		return apperror.WrapSimple(err, "downloader.downloadSegmentChunk.fetch")
	}

	return writeSegmentBody(file, seg, resp)
}

func downloadSegmentChunkWithRetry(urlText string, file *os.File, seg SegmentRange) error {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		err := downloadSegmentChunk(urlText, file, seg)
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return lastErr
}

func runSegmentWorker(urlText string, file *os.File, seg SegmentRange, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	if err := downloadSegmentChunkWithRetry(urlText, file, seg); err != nil {
		errChan <- err
	}
}

func readFirstSegmentError(errChan <-chan error) error {
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func executeSegmentWorkers(urlText string, file *os.File, segments []SegmentRange) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(segments))
	for _, seg := range segments {
		wg.Add(1)
		go runSegmentWorker(urlText, file, seg, &wg, errChan)
	}

	wg.Wait()
	close(errChan)

	return readFirstSegmentError(errChan)
}

// DownloadParallelRange downloads segments concurrently into a pre-allocated file.
func DownloadParallelRange(urlText, dstPath string, totalSize int64, segments []SegmentRange) error {
	file, err := preparePreallocatedFile(dstPath, totalSize)
	if err != nil {
		return err
	}
	defer file.Close()

	return executeSegmentWorkers(urlText, file, segments)
}

// DownloadSegmented coordinates Tier 2 10-connection parallel Range download.
func DownloadSegmented(opts DownloadOptions, totalSize int64) (DownloadResult, error) {
	start := time.Now()
	segs := CalculateSegments(totalSize, DefaultSegmentCount)
	if err := DownloadParallelRange(opts.URL, opts.DestinationPath, totalSize, segs); err != nil {
		return DownloadResult{}, err
	}

	return buildCompletedResult(opts, Tier2NativeRange, totalSize, time.Since(start))
}
