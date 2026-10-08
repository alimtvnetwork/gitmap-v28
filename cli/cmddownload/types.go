package cmddownload

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// DownloadEngine represents the active downloader backend.
type DownloadEngine string

const (
	EngineAria2c DownloadEngine = "aria2c"
	EngineCurl   DownloadEngine = "curl"
	EngineGoHTTP DownloadEngine = "go_http"
	EngineAuto   DownloadEngine = "auto"
)

// DownloadOptions parameterizes the download operation.
type DownloadOptions struct {
	URL              string         `json:"url"`
	OutputPath       string         `json:"outputPath"`
	Engine           DownloadEngine `json:"engine"`
	Threads          int            `json:"threads"`
	Splits           int            `json:"splits"`
	MinSplitSize     string         `json:"minSplitSize"`
	IsForceOverwrite bool           `json:"isForceOverwrite"`
	IsQuiet          bool           `json:"isQuiet"`
	IsJSON           bool           `json:"isJSON"`
}

// DownloadProgress captures instantaneous progress metrics.
type DownloadProgress struct {
	DownloadedBytes int64         `json:"downloadedBytes"`
	TotalBytes      int64         `json:"totalBytes"`
	Percent         float64       `json:"percent"`
	BytesPerSec     int64         `json:"bytesPerSec"`
	Elapsed         time.Duration `json:"elapsed"`
	ETA             time.Duration `json:"eta"`
}

// DownloadResult encapsulates completed download telemetry matching spec JSON contract.
type DownloadResult struct {
	Status              string         `json:"status"`
	URL                 string         `json:"url"`
	DestinationPath     string         `json:"destination_path"`
	FileName            string         `json:"file_name"`
	FileSizeBytes       int64          `json:"file_size_bytes"`
	DownloadedBytes     int64          `json:"downloaded_bytes"`
	DurationMs          int64          `json:"duration_ms"`
	AverageSpeedBps     int64          `json:"average_speed_bytes_sec"`
	EngineUsed          DownloadEngine `json:"engine_used"`
	Tier                int            `json:"tier"`
	HasFallbackOccurred bool           `json:"fallback_occurred"`
	ErrorMessage        *string        `json:"error"`
}

// ResultDownload wraps DownloadResult in monadic Result.
type ResultDownload = result.Result[DownloadResult]
