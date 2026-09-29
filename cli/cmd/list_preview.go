package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// PrintListPreview renders repositories in numbered one-line gap format.
func PrintListPreview(records []model.ScanRecord) {
	if len(records) == 0 {
		return
	}
	for i, r := range records {
		name := resolveRecordDisplayName(r)
		fmt.Printf("%d. %s\n%s\n\n", i+1, name, r.AbsolutePath)
	}
}

func resolveRecordDisplayName(r model.ScanRecord) string {
	if r.RepoName != "" {
		return r.RepoName
	}
	return r.Slug
}

// PrintListTree renders repositories in an emoji tree grouped by parent directory.
func PrintListTree(records []model.ScanRecord) {
	if len(records) == 0 {
		return
	}
	groups := groupRecordsByParent(records)
	keys := sortedListParentKeys(groups)
	seq := 1
	for _, parent := range keys {
		recs := groups[parent]
		fmt.Printf("📁 %s\n", parent)
		seq = renderTreeGroupRows(recs, seq)
		fmt.Println()
	}
}

func renderTreeGroupRows(recs []model.ScanRecord, startSeq int) int {
	seq := startSeq
	for j, r := range recs {
		branchInfo := resolveBranchInfo(r.Branch)
		prefix := "├──"
		if j == len(recs)-1 {
			prefix = "└──"
		}
		fmt.Printf("%s %d. 📦 %s %s\n", prefix, seq, r.RepoName, branchInfo)
		seq++
	}
	return seq
}

func resolveBranchInfo(branch string) string {
	if branch == "" {
		return "[git]"
	}
	return fmt.Sprintf("[git: %s]", branch)
}

func groupRecordsByParent(records []model.ScanRecord) map[string][]model.ScanRecord {
	groups := make(map[string][]model.ScanRecord)
	for _, r := range records {
		parent := filepath.Dir(r.AbsolutePath)
		if parent == "" || parent == "." {
			parent = "Repositories"
		}
		groups[parent] = append(groups[parent], r)
	}
	return groups
}

func sortedListParentKeys(groups map[string][]model.ScanRecord) []string {
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
