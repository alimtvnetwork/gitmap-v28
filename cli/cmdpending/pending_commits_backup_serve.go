package cmdpending

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func handleBackupMode(backupConn *sql.DB, errBackup error, opts PendingCommitsOptions) error {
	if errBackup != nil {
		return errBackup
	}
	if backupConn == nil {
		return apperror.NewSimple("backup database is not available", "E9004")
	}
	return serveBackupPendingCommits(backupConn, opts)
}

func serveBackupPendingCommits(backupConn *sql.DB, opts PendingCommitsOptions) error {
	backupRecords, errList := ListAllBackupPendingStatuses(backupConn)
	if errList != nil {
		return errList
	}
	inspected := convertBackupToPendingRecords(backupRecords, opts.TargetRepo)
	payload := buildPendingCommitsPayload(inspected, opts)
	if opts.IsJSON {
		return emitPendingCommitsJSON(payload)
	}
	renderLocalPendingCommitsTerminal(payload, len(inspected), inspected, opts)
	return nil
}

func convertBackupToPendingRecords(backupRecords []PendingCommitBackupRecord, targetRepo string) []RepoPendingCommitRecord {
	var inspected []RepoPendingCommitRecord
	for _, b := range backupRecords {
		if !isTargetRepoMatch(b, targetRepo) {
			continue
		}
		item, ok := parseBackupPayloadJSON(b.PayloadJSON)
		if ok {
			inspected = append(inspected, item)
			continue
		}
		item = RepoPendingCommitRecord{
			RepoName:             b.RepoName,
			RelativePath:         b.RepoPath,
			CurrentBranch:        b.CurrentBranch,
			Version:              b.Version,
			ShortVersionBranch:   b.ShortVersionBranch,
			IsDirty:              b.IsDirty,
			IsClean:              !b.IsDirty && b.UnpushedCount == 0,
			HasUncommitted:       b.UncommittedCount > 0,
			HasUnpushed:          b.UnpushedCount > 0,
			TotalUncommitted:     b.UncommittedCount,
			UnpushedCommitsCount: b.UnpushedCount,
		}
		if item.IsDirty {
			item.RemediationOptions = buildRepoRemediationOptions(item)
		}
		inspected = append(inspected, item)
	}
	return inspected
}

func isTargetRepoMatch(b PendingCommitBackupRecord, targetRepo string) bool {
	if targetRepo == "" || strings.EqualFold(targetRepo, "all") {
		return true
	}
	return strings.EqualFold(b.RepoName, targetRepo) || strings.EqualFold(b.RepoPath, targetRepo)
}

func parseBackupPayloadJSON(rawJSON string) (RepoPendingCommitRecord, bool) {
	if rawJSON == "" || rawJSON == "{}" {
		return RepoPendingCommitRecord{}, false
	}
	var item RepoPendingCommitRecord
	if err := json.Unmarshal([]byte(rawJSON), &item); err != nil || item.RepoName == "" {
		return RepoPendingCommitRecord{}, false
	}
	if item.IsDirty && len(item.RemediationOptions) == 0 {
		item.RemediationOptions = buildRepoRemediationOptions(item)
	}
	return item, true
}
