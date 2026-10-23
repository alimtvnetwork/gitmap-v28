package cmdclone

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
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
	cfg := termout.TableConfig{
		BorderStyle:  termout.BorderRounded,
		HasBorders:   true,
		HeaderColor:  termout.RoleSGRBold(termout.RolePrimary),
		BorderColor:  termout.SafeTokyoStyle(termout.TokyoMuted).Render(""),
		EllipsisText: "…",
		Columns:      buildCloneColumns(),
		Rows:         buildCloneTableRows(records, targetDir, isSSH),
	}

	termout.PrintTable(cfg)
}

func buildCloneColumns() []termout.Column {
	return []termout.Column{
		{Title: "ID", MinWidth: 4, MaxWidth: 6, Align: termout.AlignRight},
		{Title: "SLUG / REPOSITORY", MinWidth: 24, MaxWidth: 36, Align: termout.AlignLeft},
		{Title: "BRANCH", MinWidth: 12, MaxWidth: 18, Align: termout.AlignLeft},
		{Title: "METHOD", MinWidth: 8, MaxWidth: 8, Align: termout.AlignCenter},
		{Title: "EXISTS", MinWidth: 8, MaxWidth: 8, Align: termout.AlignCenter},
		{Title: "REMOTE URL", MinWidth: 28, MaxWidth: 44, Align: termout.AlignLeft},
	}
}

func buildCloneTableRow(idx int, r model.ScanRecord, targetDir string, isSSH bool) termout.Row {
	slug := resolveRecordSlug(r)
	branch := r.Branch
	if branch == "" {
		branch = "-"
	}
	method, remoteURL := resolveRecordMethodAndURL(r, isSSH)
	exists := formatExistsCell(resolveRecordLocalExists(r, targetDir))

	return termout.Row{
		Cells: []string{strconv.Itoa(idx), slug, branch, method, exists, remoteURL},
	}
}

func formatExistsCell(exists string) string {
	if exists == "YES" {
		return termout.TokyoOKPill("YES")
	}

	return termout.SafeTokyoStyle(termout.TokyoMuted).Render("-")
}

func buildCloneTableRows(records []model.ScanRecord, targetDir string, isSSH bool) []termout.Row {
	rows := make([]termout.Row, 0, len(records))
	for i, r := range records {
		rows = append(rows, buildCloneTableRow(i+1, r, targetDir, isSSH))
	}

	return rows
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
