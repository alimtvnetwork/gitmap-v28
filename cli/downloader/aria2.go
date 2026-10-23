package downloader

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// HasAria2c checks if the aria2c executable is present in PATH.
func HasAria2c() bool {
	_, err := exec.LookPath("aria2c")

	return err == nil
}

func buildAria2cArgs(url, destDir, destFile string, segments int) []string {
	s := strconv.Itoa(segments)

	return []string{
		"-s", s, "-x", s, "-k", "5K",
		"--disable-ipv6=true",
		"--file-allocation=none",
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		"--summary-interval=0",
		"--console-log-level=error",
		"--show-console-readout=false",
		"--dir=" + destDir, "-o", destFile, url,
	}
}

// DownloadWithAria2c executes aria2c with multi-connection acceleration.
func DownloadWithAria2c(url, destDir, destFile string, segments int) error {
	if segments < 1 {
		segments = DefaultSegmentCount
	}

	args := buildAria2cArgs(url, destDir, destFile, segments)
	cmd := exec.Command("aria2c", args...)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "downloader.DownloadWithAria2c")
	}

	return nil
}

// DownloadWithAria2cOptions executes aria2c using DownloadOptions.
func DownloadWithAria2cOptions(opts DownloadOptions) (DownloadResult, error) {
	start := time.Now()
	destDir := filepath.Dir(opts.DestinationPath)
	destFile := filepath.Base(opts.DestinationPath)
	if err := DownloadWithAria2c(opts.URL, destDir, destFile, DefaultSegmentCount); err != nil {
		return DownloadResult{}, err
	}

	return buildCompletedResult(opts, Tier1Aria2c, 0, time.Since(start))
}
