package folder

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// FolderReport represents the aggregated summary payload of a scanned directory hierarchy.
type FolderReport struct {
	Root               string      `json:"root" yaml:"root"`
	TotalFiles         int         `json:"totalFiles" yaml:"totalFiles"`
	TotalLines         int         `json:"totalLines" yaml:"totalLines"`
	TotalSizeBytes     int64       `json:"totalSizeBytes" yaml:"totalSizeBytes"`
	TotalSizeFormatted string      `json:"totalSizeFormatted" yaml:"totalSizeFormatted"`
	Files              []*FileMeta `json:"files" yaml:"files"`
}

// BuildReport aggregates multiple FileMeta records into a structured summary report.
func BuildReport(root string, files []*FileMeta) *FolderReport {
	var totalLines int
	var totalBytes int64
	for _, f := range files {
		totalLines += f.LinesOfCode
		totalBytes += f.SizeBytes
	}
	return &FolderReport{
		Root:               root,
		TotalFiles:         len(files),
		TotalLines:         totalLines,
		TotalSizeBytes:     totalBytes,
		TotalSizeFormatted: formatBytes(totalBytes),
		Files:              files,
	}
}

// RenderJson formats the folder report as indented JSON.
func RenderJson(report *FolderReport) (string, error) {
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// RenderYaml formats the folder report as YAML text.
func RenderYaml(report *FolderReport) (string, error) {
	b, err := yaml.Marshal(report)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// RenderFlat formats files as a flat list of paths with optional inline metadata.
func RenderFlat(files []*FileMeta, isDetailed bool) string {
	var sb strings.Builder
	for _, f := range files {
		if isDetailed {
			metaInfo := formatMetaSuffix(f, true)
			sb.WriteString(fmt.Sprintf("%s%s\n", f.Path, metaInfo))
		} else {
			sb.WriteString(f.Path + "\n")
		}
	}
	return sb.String()
}
