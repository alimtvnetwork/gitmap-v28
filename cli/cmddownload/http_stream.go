package cmddownload

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// downloadWithGoHTTP streams the resource using native Go net/http.
func downloadWithGoHTTP(opts DownloadOptions, destDir, destFile string) (DownloadResult, error) {
	fullPath := filepath.Join(destDir, destFile)
	startTime := time.Now()

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   0, // No global timeout for large streaming files
	}

	req, errReq := http.NewRequestWithContext(context.Background(), http.MethodGet, opts.URL, nil)
	if errReq != nil {
		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineGoHTTP,
			Tier:            3,
			DurationMs:      0,
		}, apperror.WrapSimple(errReq, "download.http.req")
	}

	req.Header.Set("User-Agent", "gitmap-download/1.0")

	resp, errResp := client.Do(req)
	if errResp != nil {
		duration := time.Since(startTime)

		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineGoHTTP,
			Tier:            3,
			DurationMs:      duration.Milliseconds(),
		}, apperror.WrapSimple(errResp, "download.http.do")
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		duration := time.Since(startTime)

		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineGoHTTP,
			Tier:            3,
			DurationMs:      duration.Milliseconds(),
		}, apperror.NewSimple(fmt.Sprintf("HTTP error %d: %s", resp.StatusCode, resp.Status), "E1203")
	}

	totalBytes := resp.ContentLength

	outFile, errOut := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if errOut != nil {
		duration := time.Since(startTime)

		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineGoHTTP,
			Tier:            3,
			DurationMs:      duration.Milliseconds(),
		}, apperror.WrapSimple(errOut, "download.http.create")
	}

	defer outFile.Close()

	bar := NewProgressBar(destFile, totalBytes, opts.IsQuiet, opts.IsJSON)
	reader := &ProgressReader{
		Reader:     resp.Body,
		Bar:        bar,
		Downloaded: 0,
	}

	copiedBytes, errCopy := io.Copy(outFile, reader)
	duration := time.Since(startTime)

	if errCopy != nil {
		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineGoHTTP,
			Tier:            3,
			DurationMs:      duration.Milliseconds(),
		}, apperror.WrapSimple(errCopy, "download.http.copy")
	}

	var avgSpeed int64
	if duration.Seconds() > 0 {
		avgSpeed = int64(float64(copiedBytes) / duration.Seconds())
	}

	bar.Finish(fullPath, copiedBytes, duration)

	return DownloadResult{
		Status:              "success",
		URL:                 opts.URL,
		DestinationPath:     fullPath,
		FileName:            destFile,
		FileSizeBytes:       copiedBytes,
		DownloadedBytes:     copiedBytes,
		DurationMs:          duration.Milliseconds(),
		AverageSpeedBps:     avgSpeed,
		EngineUsed:          EngineGoHTTP,
		Tier:                3,
		HasFallbackOccurred: true,
	}, nil
}
