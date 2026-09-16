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

func runCloneListTable(sourcePath, targetDir string, isSSH bool, onlyFilter, excludeFilter string) error {
	records, err := cloner.LoadRecords(sourcePath)
	if err != nil {
		return apperror.WrapSimple(err, "load records from "+sourcePath)
	}
	records = filterRecordsByOnly(records, onlyFilter)
	records = filterRecordsByExclude(records, excludeFilter)
	if len(records) == 0 {
		fmt.Printf("No repositories found in %s\n", sourcePath)
		return nil
	}
	printCloneRecordsTable(records, targetDir, isSSH)
	return nil
}

func printCloneRecordsTable(records []model.ScanRecord, targetDir string, isSSH bool) {
	fmt.Printf("%-5s %-32s %-16s %-8s %-8s %s\n", "ID", "SLUG / REPOSITORY", "BRANCH", "METHOD", "EXISTS", "REMOTE URL")
	fmt.Println(strings.Repeat("-", 100))
	for i, r := range records {
		slug := resolveRecordSlug(r)
		branch := r.Branch
		if branch == "" {
			branch = "-"
		}
		method, remoteURL := resolveRecordMethodAndURL(r, isSSH)
		exists := resolveRecordLocalExists(r, targetDir)
		fmt.Printf("%-5d %-32s %-16s %-8s %-8s %s\n", i+1, slug, branch, method, exists, remoteURL)
	}
}

func resolveRecordMethodAndURL(r model.ScanRecord, isSSH bool) (string, string) {
	rawURL := resolveRecordRemoteURL(r)
	if isSSH {
		return "SSH", resolveRecordSSHURL(rawURL)
	}
	if r.Transport == "ssh" || isSSHCloneURL(rawURL) {
		return "SSH", rawURL
	}
	return "HTTPS", rawURL
}

func resolveRecordSSHURL(rawURL string) string {
	sshURL, ok := ConvertURLToSSH(rawURL)
	if ok {
		return sshURL
	}
	return rawURL
}

func resolveRecordLocalExists(r model.ScanRecord, targetDir string) string {
	if targetDir == "" {
		targetDir = "."
	}
	dest := filepath.Join(targetDir, model.CleanRelativePath(r.RelativePath))
	if isGitRepo(dest) {
		return "YES"
	}
	return "-"
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

func filterRecordsByExclude(records []model.ScanRecord, excludeFilter string) []model.ScanRecord {
	if strings.TrimSpace(excludeFilter) == "" {
		return records
	}
	tokens := splitFilterTokens(excludeFilter)
	var filtered []model.ScanRecord
	for i, r := range records {
		seqID := i + 1
		if isRecordMatchTokens(seqID, r, tokens) {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
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
