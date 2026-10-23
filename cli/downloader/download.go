package downloader

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DownloadTier indicates the acceleration level utilized.
type DownloadTier int

const (
	Tier1Aria2c       DownloadTier = 1
	Tier2NativeRange  DownloadTier = 2
	Tier3SingleStream DownloadTier = 3
)

// DownloadOptions holds execution configuration for URL retrieval.
type DownloadOptions struct {
	URL             string
	DestinationPath string
	IsForceFresh    bool
	IsQuiet         bool
	IsJSON          bool
	Timeout         time.Duration
	ChunkCount      int
	MinChunkSize    int64
	Headers         map[string]string
}

// DownloadResult summarizes the completed download operation.
type DownloadResult struct {
	URL              string       `json:"url"`
	DestinationPath  string       `json:"destinationPath"`
	ContentLength    int64        `json:"contentLength"`
	DurationMs       int64        `json:"durationMs"`
	TierUsed         DownloadTier `json:"tierUsed"`
	IsRangeSupported bool         `json:"isRangeSupported"`
	IsCacheHit       bool         `json:"isCacheHit"`
	Sha256Checksum   string       `json:"sha256Checksum"`
	PayloadReader    io.Reader    `json:"-"`
}

type resourceMeta struct {
	ContentLength   int64
	HasRangeSupport bool
}

func parseResponseMeta(resp *http.Response) resourceMeta {
	hasRange := resp.Header.Get("Accept-Ranges") == "bytes" || resp.StatusCode == http.StatusPartialContent
	length := resp.ContentLength
	if length < 0 {
		length = 0
	}

	return resourceMeta{
		ContentLength:   length,
		HasRangeSupport: hasRange,
	}
}

func probeViaHead(urlText string) (resourceMeta, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodHead, urlText, nil)
	if err != nil {
		return resourceMeta{}, apperror.WrapSimple(err, "downloader.probeViaHead.req")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resourceMeta{}, apperror.WrapSimple(err, "downloader.probeViaHead.do")
	}
	defer resp.Body.Close()

	return parseResponseMeta(resp), nil
}

func parseContentRangeTotal(cr string) int64 {
	idx := strings.LastIndex(cr, "/")
	if idx == -1 || idx+1 >= len(cr) {
		return 0
	}

	total, err := strconv.ParseInt(cr[idx+1:], 10, 64)
	if err != nil {
		return 0
	}

	return total
}

func parseRangeResponseMeta(resp *http.Response) resourceMeta {
	hasRange := resp.StatusCode == http.StatusPartialContent
	length := parseContentRangeTotal(resp.Header.Get("Content-Range"))
	if length <= 0 {
		length = resp.ContentLength
	}

	return resourceMeta{
		ContentLength:   length,
		HasRangeSupport: hasRange,
	}
}

func executeRangeProbe(req *http.Request) (resourceMeta, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resourceMeta{}, apperror.WrapSimple(err, "downloader.executeRangeProbe.do")
	}
	defer resp.Body.Close()

	return parseRangeResponseMeta(resp), nil
}

func probeViaRange(urlText string) (resourceMeta, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, urlText, nil)
	if err != nil {
		return resourceMeta{}, apperror.WrapSimple(err, "downloader.probeViaRange.req")
	}

	req.Header.Set("Range", "bytes=0-0")

	return executeRangeProbe(req)
}

func probeRemoteResource(urlText string) (resourceMeta, error) {
	meta, err := probeViaHead(urlText)
	if err == nil && meta.ContentLength > 0 {
		return meta, nil
	}

	return probeViaRange(urlText)
}

func copyResponseBodyToFile(dstPath string, resp *http.Response) error {
	defer resp.Body.Close()
	out, err := os.Create(dstPath)
	if err != nil {
		return apperror.WrapSimple(err, "downloader.copyResponseBodyToFile.create")
	}
	defer out.Close()

	_, copyErr := io.Copy(out, resp.Body)

	return apperror.WrapSimple(copyErr, "downloader.copyResponseBodyToFile.copy")
}

func downloadSingleStreamFile(urlText, dstPath string) error {
	resp, err := http.Get(urlText)
	if err != nil {
		return apperror.WrapSimple(err, "downloader.downloadSingleStreamFile.get")
	}

	return copyResponseBodyToFile(dstPath, resp)
}

func tryTier1Aria2c(urlText, dstPath string) error {
	dir := filepath.Dir(dstPath)
	file := filepath.Base(dstPath)

	return DownloadWithAria2c(urlText, dir, file, DefaultSegmentCount)
}

func tryAria2cIfAvailable(urlText, dstPath string) bool {
	if !HasAria2c() {
		return false
	}

	return tryTier1Aria2c(urlText, dstPath) == nil
}

func executeSegmentedDownload(urlText, dstPath string, totalSize int64) error {
	segments := CalculateSegments(totalSize, DefaultSegmentCount)

	return DownloadParallelRange(urlText, dstPath, totalSize, segments)
}

func downloadAccelerated(urlText, dstPath string, totalSize int64) error {
	if tryAria2cIfAvailable(urlText, dstPath) {
		return nil
	}

	return executeSegmentedDownload(urlText, dstPath, totalSize)
}

// DownloadFile downloads a remote URL to dstPath using multi-tier acceleration.
func DownloadFile(urlText, dstPath string, opts DownloadOptions) error {
	meta, err := probeRemoteResource(urlText)
	if err != nil || !meta.HasRangeSupport || meta.ContentLength < MinSegmentThreshold {
		return downloadSingleStreamFile(urlText, dstPath)
	}

	if err := downloadAccelerated(urlText, dstPath, meta.ContentLength); err == nil {
		return nil
	}

	return downloadSingleStreamFile(urlText, dstPath)
}

func resolveDestinationPath(opts DownloadOptions) (string, error) {
	if opts.DestinationPath != "" {
		return opts.DestinationPath, nil
	}

	f, err := os.CreateTemp("", "gitmap-dl-*")
	if err != nil {
		return "", apperror.WrapSimple(err, "downloader.resolveDestinationPath.temp")
	}

	path := f.Name()
	_ = f.Close()

	return path, nil
}

func determineDownloadTier(meta resourceMeta) DownloadTier {
	if !meta.HasRangeSupport || meta.ContentLength < MinSegmentThreshold {
		return Tier3SingleStream
	}

	if HasAria2c() {
		return Tier1Aria2c
	}

	return Tier2NativeRange
}

func buildCompletedResult(opts DownloadOptions, tier DownloadTier, size int64, dur time.Duration) (DownloadResult, error) {
	f, err := os.Open(opts.DestinationPath)
	if err != nil {
		return DownloadResult{}, apperror.WrapSimple(err, "downloader.buildCompletedResult.open")
	}

	return DownloadResult{
		URL:              opts.URL,
		DestinationPath:  opts.DestinationPath,
		ContentLength:    size,
		DurationMs:       dur.Milliseconds(),
		TierUsed:         tier,
		IsRangeSupported: tier != Tier3SingleStream,
		PayloadReader:    f,
	}, nil
}

func executeDownloadPipeline(dstPath string, opts DownloadOptions) (resourceMeta, error) {
	meta, _ := probeRemoteResource(opts.URL)
	if err := DownloadFile(opts.URL, dstPath, opts); err != nil {
		return resourceMeta{}, err
	}

	return meta, nil
}

// Download orchestrates multi-tier URL retrieval returning a DownloadResult.
func Download(opts DownloadOptions) (DownloadResult, error) {
	start := time.Now()
	dstPath, err := resolveDestinationPath(opts)
	if err != nil {
		return DownloadResult{}, err
	}

	opts.DestinationPath = dstPath
	meta, err := executeDownloadPipeline(dstPath, opts)
	if err != nil {
		return DownloadResult{}, err
	}

	tier := determineDownloadTier(meta)

	return buildCompletedResult(opts, tier, meta.ContentLength, time.Since(start))
}

// DownloadSingleStream executes standard single-stream HTTP download.
func DownloadSingleStream(opts DownloadOptions) (DownloadResult, error) {
	start := time.Now()
	dstPath, err := resolveDestinationPath(opts)
	if err != nil {
		return DownloadResult{}, err
	}
	opts.DestinationPath = dstPath

	if err := downloadSingleStreamFile(opts.URL, dstPath); err != nil {
		return DownloadResult{}, err
	}

	return buildCompletedResult(opts, Tier3SingleStream, 0, time.Since(start))
}
