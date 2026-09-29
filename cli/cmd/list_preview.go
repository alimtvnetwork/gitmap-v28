package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// ListPreviewOptions holds display options for list preview.
type ListPreviewOptions struct {
	ShowNumbers bool
	ShowEmoji   bool
	UseRelative bool
}

// DefaultListPreviewOptions returns default preview options.
func DefaultListPreviewOptions() ListPreviewOptions {
	return ListPreviewOptions{
		ShowNumbers: true,
		ShowEmoji:   true,
	}
}

// PrintListPreview renders repositories in numbered one-line gap format.
func PrintListPreview(records []model.ScanRecord) {
	PrintListPreviewWithOptions(records, DefaultListPreviewOptions())
}

// PrintListPreviewWithOptions renders repositories in one-line gap format with options.
func PrintListPreviewWithOptions(records []model.ScanRecord, opts ListPreviewOptions) {
	if len(records) == 0 {
		fmt.Println(constants.MsgListEmpty)
		return
	}
	for i, r := range records {
		name := resolveRecordDisplayName(r)
		branchInfo := resolveBranchInfo(r.Branch)
		header := formatListPreviewHeader(i+1, name, branchInfo, opts)
		pathStr := resolveRecordPath(r, opts.UseRelative)
		fmt.Printf("%s\n%s\n\n", header, pathStr)
	}
}

func resolveRecordPath(r model.ScanRecord, useRelative bool) string {
	if useRelative && r.RelativePath != "" {
		return r.RelativePath
	}
	return r.AbsolutePath
}

func formatListPreviewHeader(seq int, name, branchInfo string, opts ListPreviewOptions) string {
	prefix := ""
	if opts.ShowNumbers {
		prefix = fmt.Sprintf("%d. ", seq)
	}
	icon := ""
	if opts.ShowEmoji {
		icon = "📦 "
	}
	return fmt.Sprintf("%s%s%s %s", prefix, icon, name, branchInfo)
}

func resolveRecordDisplayName(r model.ScanRecord) string {
	if r.RepoName != "" {
		return r.RepoName
	}
	return r.Slug
}

func resolveBranchInfo(branch string) string {
	if branch == "" {
		return "[git]"
	}
	return fmt.Sprintf("[git: %s]", branch)
}
