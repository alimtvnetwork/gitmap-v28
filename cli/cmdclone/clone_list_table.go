package cmdclone

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func runCloneListTable(sourcePath string) error {
	records, err := cloner.LoadRecords(sourcePath)
	if err != nil {
		return apperror.WrapSimple(err, "load records from "+sourcePath)
	}
	if len(records) == 0 {
		fmt.Printf("No repositories found in %s\n", sourcePath)
		return nil
	}
	printCloneRecordsTable(records)
	return nil
}

func printCloneRecordsTable(records []model.ScanRecord) {
	fmt.Printf("%-5s %-32s %-16s %s\n", "ID", "SLUG / REPOSITORY", "BRANCH", "REMOTE URL")
	fmt.Println(strings.Repeat("-", 80))
	for i, r := range records {
		slug := resolveRecordSlug(r)
		branch := r.Branch
		if branch == "" {
			branch = "-"
		}
		fmt.Printf("%-5d %-32s %-16s %s\n", i+1, slug, branch, resolveRecordRemoteURL(r))
	}
}

func resolveRecordRemoteURL(r model.ScanRecord) string {
	if r.DiscoveredURL != "" {
		return r.DiscoveredURL
	}
	if r.HTTPSUrl != "" {
		return r.HTTPSUrl
	}
	return r.SSHUrl
}

func resolveRecordSlug(r model.ScanRecord) string {
	if r.RelativePath != "" {
		return filepath.ToSlash(r.RelativePath)
	}
	remoteURL := resolveRecordRemoteURL(r)
	parts := strings.Split(strings.TrimSuffix(remoteURL, ".git"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return "-"
}

func filterRecordsByOnly(records []model.ScanRecord, onlyFilter string) []model.ScanRecord {
	if strings.TrimSpace(onlyFilter) == "" {
		return records
	}
	tokens := splitFilterTokens(onlyFilter)
	var matched []model.ScanRecord
	for i, r := range records {
		seqID := i + 1
		if isRecordMatchTokens(seqID, r, tokens) {
			matched = append(matched, r)
		}
	}
	return matched
}

func splitFilterTokens(filter string) []string {
	parts := strings.Split(filter, ",")
	var tokens []string
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			tokens = append(tokens, strings.ToLower(t))
		}
	}
	return tokens
}

func isRecordMatchTokens(seqID int, r model.ScanRecord, tokens []string) bool {
	seqStr := strconv.Itoa(seqID)
	slug := strings.ToLower(resolveRecordSlug(r))
	rel := strings.ToLower(filepath.ToSlash(r.RelativePath))
	url := strings.ToLower(resolveRecordRemoteURL(r))

	for _, token := range tokens {
		if token == seqStr || isPatternMatch(slug, rel, url, token) {
			return true
		}
	}
	return false
}

func isPatternMatch(slug, rel, url, token string) bool {
	if strings.HasSuffix(token, "*") {
		prefix := strings.TrimSuffix(token, "*")
		return strings.HasPrefix(slug, prefix) || strings.HasPrefix(rel, prefix)
	}
	return slug == token || rel == token || strings.HasPrefix(slug, token) || strings.Contains(url, token)
}
