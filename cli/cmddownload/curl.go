package cmddownload

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// HasCurl returns true if curl is found in PATH.
func HasCurl() bool {
	_, err := exec.LookPath("curl")

	return err == nil
}

// downloadWithCurl executes Tier 2 curl fallback download.
func downloadWithCurl(opts DownloadOptions, destDir, destFile string) (DownloadResult, error) {
	fullPath := filepath.Join(destDir, destFile)
	startTime := time.Now()

	args := []string{
		"-fL",
		"--retry", "3",
		"--retry-delay", "2",
		"-o", fullPath,
		opts.URL,
	}

	if opts.IsQuiet || opts.IsJSON {
		args = append([]string{"-s"}, args...)
	} else {
		args = append([]string{"--progress-bar"}, args...)
	}

	cmd := exec.Command("curl", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if !opts.IsQuiet && !opts.IsJSON {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Run()
	duration := time.Since(startTime)

	if err != nil {
		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineCurl,
			Tier:            2,
			DurationMs:      duration.Milliseconds(),
		}, apperror.WrapSimple(fmt.Errorf("curl failed: %w, stderr: %s", err, stderr.String()), "download.curl")
	}

	info, statErr := os.Stat(fullPath)
	if statErr != nil || info.Size() == 0 {
		return DownloadResult{
			Status:          "error",
			URL:             opts.URL,
			DestinationPath: fullPath,
			FileName:        destFile,
			EngineUsed:      EngineCurl,
			Tier:            2,
			DurationMs:      duration.Milliseconds(),
		}, apperror.NewSimple("curl produced zero-byte or missing file", "E1203")
	}

	var avgSpeed int64
	if duration.Seconds() > 0 {
		avgSpeed = int64(float64(info.Size()) / duration.Seconds())
	}

	return DownloadResult{
		Status:              "success",
		URL:                 opts.URL,
		DestinationPath:     fullPath,
		FileName:            destFile,
		FileSizeBytes:       info.Size(),
		DownloadedBytes:     info.Size(),
		DurationMs:          duration.Milliseconds(),
		AverageSpeedBps:     avgSpeed,
		EngineUsed:          EngineCurl,
		Tier:                2,
		HasFallbackOccurred: true,
	}, nil
}
