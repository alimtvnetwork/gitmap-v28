package cmddownload

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// HasAria2c returns true if aria2c is found in PATH.
func HasAria2c() bool {
	_, err := exec.LookPath("aria2c")

	return err == nil
}

// downloadWithAria2c executes Tier 1 aria2c multi-connection download.
func downloadWithAria2c(opts DownloadOptions, destDir, destFile string) (DownloadResult, error) {
	fullPath := filepath.Join(destDir, destFile)
	startTime := time.Now()

	threads := opts.Threads
	if threads <= 0 {
		threads = 16
	}

	splits := opts.Splits
	if splits <= 0 {
		splits = 80
	}

	minSplit := opts.MinSplitSize
	if minSplit == "" {
		minSplit = "1M"
	}

	args := []string{
		"--disable-ipv6=true",
		"-x", strconv.Itoa(threads),
		"-s", strconv.Itoa(splits),
		"-j", "16",
		"-k", minSplit,
		"--file-allocation=none",
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		"--summary-interval=0",
		"--console-log-level=error",
		"--show-console-readout=false",
		"--dir=" + destDir,
		"-o", destFile,
		opts.URL,
	}

	cmd := exec.Command("aria2c", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	bar := NewProgressBar(destFile, 0, opts.IsQuiet, opts.IsJSON)
	if !opts.IsQuiet && !opts.IsJSON {
		bar.Update(0)
	}

	err := cmd.Run()
	duration := time.Since(startTime)

	if err != nil {
		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineAria2c,
			Tier:            1,
			DurationMs:      duration.Milliseconds(),
		}, apperror.WrapSimple(fmt.Errorf("aria2c failed: %w, stderr: %s", err, stderr.String()), "download.aria2c")
	}

	info, statErr := os.Stat(fullPath)
	if statErr != nil || info.Size() == 0 {
		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineAria2c,
			Tier:            1,
			DurationMs:      duration.Milliseconds(),
		}, apperror.NewSimple("aria2c produced zero-byte or missing file", "E1203")
	}

	var avgSpeed int64
	if duration.Seconds() > 0 {
		avgSpeed = int64(float64(info.Size()) / duration.Seconds())
	}

	bar.Finish(fullPath, info.Size(), duration)

	return DownloadResult{
		Status:              "success",
		URL:                 opts.URL,
		DestinationPath:     fullPath,
		FileName:            destFile,
		FileSizeBytes:       info.Size(),
		DownloadedBytes:     info.Size(),
		DurationMs:          duration.Milliseconds(),
		AverageSpeedBps:     avgSpeed,
		EngineUsed:          EngineAria2c,
		Tier:                1,
		HasFallbackOccurred: false,
	}, nil
}
