package pipelinedb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func computeRelativeDbPath(absPath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.ToSlash(absPath)
	}
	rel, relErr := filepath.Rel(cwd, absPath)
	if relErr != nil || strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(absPath)
	}

	return formatRelativePrefix(filepath.ToSlash(rel))
}

func formatRelativePrefix(slashRel string) string {
	if strings.HasPrefix(slashRel, "./") || strings.HasPrefix(slashRel, "../") {
		return slashRel
	}

	return "./" + slashRel
}

func formatUnitSize(val float64, unit string) string {
	if val == float64(int64(val)) {
		return fmt.Sprintf("%d %s", int64(val), unit)
	}

	return fmt.Sprintf("%.1f %s", val, unit)
}

// FormatHumanSize formats byte counts into human-readable strings (B, KB, MB, GB).
func FormatHumanSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return formatUnitSize(float64(bytes)/1024.0, "KB")
	}
	if bytes < 1024*1024*1024 {
		return formatUnitSize(float64(bytes)/(1024.0*1024.0), "MB")
	}

	return formatUnitSize(float64(bytes)/(1024.0*1024.0*1024.0), "GB")
}

func formatHumanSize(bytes int64) string {
	return FormatHumanSize(bytes)
}

func resolveSummaryPath(info PipelineDatabaseInfo) string {
	if info.RelativePath != "" {
		return info.RelativePath
	}

	return info.Path
}

func resolveRunCountLabel(totalRuns int) string {
	if totalRuns == 1 {
		return "run"
	}

	return "runs"
}

// FormatPipelineDbSummary formats a concise summary string for pipeline database telemetry.
func FormatPipelineDbSummary(info PipelineDatabaseInfo) string {
	pathStr := resolveSummaryPath(info)
	lbl := resolveRunCountLabel(info.TotalRuns)
	sizeStr := info.HumanSize
	if sizeStr == "" {
		sizeStr = formatHumanSize(info.SizeBytes)
	}

	return fmt.Sprintf("%s (%s, %d %s)", pathStr, sizeStr, info.TotalRuns, lbl)
}

func isRegularDbFile(fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}

	return true
}

func (p *PipelineSplitDb) populateFileMetadata(info *PipelineDatabaseInfo) {
	fi, err := os.Stat(p.Path)
	if err != nil {
		return
	}
	info.SizeBytes = fi.Size()
	info.HumanSize = formatHumanSize(fi.Size())
	info.IsExisting = isRegularDbFile(fi)
}

func (p *PipelineSplitDb) populateCounts(info *PipelineDatabaseInfo) {
	if p.conn == nil {
		return
	}
	info.TotalRuns, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun;")
	info.FailedRuns, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineRun WHERE IsSuccess = 0;")
	info.JobCount, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineJob;")
	info.SegmentCount, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineSegment;")
	info.ErrorCount, _ = countQuery(p.conn, "SELECT COUNT(*) FROM PipelineErrorLog;")
	last, _ := queryLastUpdated(p.conn)
	info.LastUpdated = last
}

// GetDatabaseInfo returns diagnostic telemetry and metadata for the database.
func (p *PipelineSplitDb) GetDatabaseInfo() PipelineDatabaseInfo {
	var info PipelineDatabaseInfo
	info.Path = filepath.ToSlash(p.Path)
	info.RelativePath = filepath.ToSlash(computeRelativeDbPath(p.Path))
	p.populateFileMetadata(&info)
	p.populateCounts(&info)

	return info
}

// Vacuum executes SQLite VACUUM and returns reclaimed bytes.
func (p *PipelineSplitDb) Vacuum() (int64, error) {
	if p.conn == nil {
		return 0, nil
	}
	beforeSize := resolveFileSize(p.Path)
	if _, err := p.conn.Exec("VACUUM;"); err != nil {
		return 0, apperror.WrapSimple(err, "vacuum pipeline db")
	}

	return calculateFreedBytes(beforeSize, resolveFileSize(p.Path)), nil
}

func resolveFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return fi.Size()
}

func calculateFreedBytes(before, after int64) int64 {
	if before > after {
		return before - after
	}

	return 0
}
