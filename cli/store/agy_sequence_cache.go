// Package store — agy_sequence_cache.go manages the 24-hour SQLite sequence cache for Antigravity projects and prompts.
package store

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	// SequenceCacheTTLSeconds defines the 24-hour TTL for cached project/prompt sequences.
	SequenceCacheTTLSeconds = 86400

	sqlCreateAgySequenceCache = `CREATE TABLE IF NOT EXISTS AgySequenceCache (
    SeqId TEXT PRIMARY KEY,
    SeqNum INTEGER NOT NULL,
    EntryType TEXT NOT NULL,
    ProjectId TEXT NOT NULL,
    ProjectAlias TEXT NOT NULL,
    ProjectPath TEXT NOT NULL,
    ConversationId TEXT NOT NULL DEFAULT '',
    PromptSnippet TEXT NOT NULL DEFAULT '',
    CreatedAt INTEGER NOT NULL,
    ExpiresAt INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_agy_seq_project ON AgySequenceCache(ProjectId);
CREATE INDEX IF NOT EXISTS idx_agy_seq_conv ON AgySequenceCache(ConversationId);`
)

// AgySequenceRecord represents a cached project or conversation prompt sequence row.
type AgySequenceRecord struct {
	SeqId          string `json:"seqId"`
	SeqNum         int    `json:"seqNum"`
	EntryType      string `json:"entryType"`
	ProjectId      string `json:"projectId"`
	ProjectAlias   string `json:"projectAlias"`
	ProjectPath    string `json:"projectPath"`
	ConversationId string `json:"conversationId,omitempty"`
	PromptSnippet  string `json:"promptSnippet,omitempty"`
	CreatedAt      int64  `json:"createdAt"`
	ExpiresAt      int64  `json:"expiresAt"`
}

func openSequenceCacheDB() (*BackupPromptsSplitDB, error) {
	db, err := OpenBackupPromptsSplitDB("")
	if err != nil {
		return nil, err
	}
	if _, execErr := db.Conn().Exec(sqlCreateAgySequenceCache); execErr != nil {
		_ = db.Close()
		return nil, apperror.WrapSimple(execErr, "init AgySequenceCache table")
	}
	return db, nil
}

// EnsureAndGetSequenceCache checks 24h TTL validity and persists or refreshes project and prompt sequences.
func EnsureAndGetSequenceCache(entries []AgySequenceRecord, isForceRefresh bool) ([]AgySequenceRecord, error) {
	db, err := openSequenceCacheDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	now := time.Now().Unix()
	cached := queryValidSequenceRows(db.Conn(), now)
	if !isForceRefresh && len(cached) > 0 && len(entries) == 0 {
		return cached, nil
	}
	return refreshSequenceCacheRows(db.Conn(), entries, cached, now, isForceRefresh)
}

func refreshSequenceCacheRows(conn *sql.DB, entries, cached []AgySequenceRecord, now int64, isForce bool) ([]AgySequenceRecord, error) {
	if isForce || len(cached) == 0 {
		_, _ = conn.Exec("DELETE FROM AgySequenceCache WHERE ExpiresAt <= ? OR 1=1", now)
	}
	if len(entries) == 0 {
		return cached, nil
	}
	normalized := normalizeSequenceEntries(entries, now)
	if err := upsertSequenceRecords(conn, normalized); err != nil {
		return nil, err
	}
	return normalized, nil
}

func normalizeSequenceEntries(entries []AgySequenceRecord, now int64) []AgySequenceRecord {
	projSeq := 0
	promptSeq := 0
	out := make([]AgySequenceRecord, len(entries))
	for i, rec := range entries {
		out[i] = assignRecordDefaults(rec, &projSeq, &promptSeq, now)
	}
	return out
}

func assignRecordDefaults(rec AgySequenceRecord, projSeq, promptSeq *int, now int64) AgySequenceRecord {
	if rec.EntryType == "" {
		rec.EntryType = "project"
	}
	if rec.EntryType == "prompt" {
		*promptSeq++
		return fillSequenceTimestamps(rec, fmt.Sprintf("P%d", *promptSeq), *promptSeq, now)
	}
	*projSeq++
	return fillSequenceTimestamps(rec, strconv.Itoa(*projSeq), *projSeq, now)
}

func fillSequenceTimestamps(rec AgySequenceRecord, defaultSeqId string, defaultNum int, now int64) AgySequenceRecord {
	if rec.SeqId == "" {
		rec.SeqId = defaultSeqId
	}
	if rec.SeqNum <= 0 {
		rec.SeqNum = defaultNum
	}
	if rec.CreatedAt <= 0 {
		rec.CreatedAt = now
	}
	if rec.ExpiresAt <= now {
		rec.ExpiresAt = now + SequenceCacheTTLSeconds
	}
	return rec
}

func upsertSequenceRecords(conn *sql.DB, records []AgySequenceRecord) error {
	for _, r := range records {
		if err := upsertSingleSequenceRecord(conn, r); err != nil {
			return err
		}
	}
	return nil
}

func upsertSingleSequenceRecord(conn *sql.DB, r AgySequenceRecord) error {
	query := `INSERT OR REPLACE INTO AgySequenceCache
(SeqId, SeqNum, EntryType, ProjectId, ProjectAlias, ProjectPath, ConversationId, PromptSnippet, CreatedAt, ExpiresAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := conn.Exec(query, r.SeqId, r.SeqNum, r.EntryType, r.ProjectId, r.ProjectAlias,
		r.ProjectPath, r.ConversationId, r.PromptSnippet, r.CreatedAt, r.ExpiresAt)
	if err != nil {
		return apperror.WrapSimple(err, "upsert AgySequenceCache")
	}
	return nil
}

func queryValidSequenceRows(conn *sql.DB, now int64) []AgySequenceRecord {
	_, _ = conn.Exec("DELETE FROM AgySequenceCache WHERE ExpiresAt <= ?", now)
	rows, err := conn.Query(`SELECT SeqId, SeqNum, EntryType, ProjectId, ProjectAlias, ProjectPath,
ConversationId, PromptSnippet, CreatedAt, ExpiresAt FROM AgySequenceCache WHERE ExpiresAt > ? ORDER BY SeqNum ASC`, now)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanSequenceRows(rows)
}

func scanSequenceRows(rows *sql.Rows) []AgySequenceRecord {
	var list []AgySequenceRecord
	for rows.Next() {
		var r AgySequenceRecord
		err := rows.Scan(&r.SeqId, &r.SeqNum, &r.EntryType, &r.ProjectId, &r.ProjectAlias,
			&r.ProjectPath, &r.ConversationId, &r.PromptSnippet, &r.CreatedAt, &r.ExpiresAt)
		if err == nil {
			list = append(list, r)
		}
	}
	return list
}

// ResolveSequenceEntry matches a sequence ID ("1", "#1", "P1"), project ID, alias, path, or conversation ID.
func ResolveSequenceEntry(token string) (*AgySequenceRecord, error) {
	clean := cleanSequenceToken(token)
	if clean == "" {
		return nil, apperror.NewSimple("empty sequence or target token", "E9040")
	}
	db, err := openSequenceCacheDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows := queryValidSequenceRows(db.Conn(), time.Now().Unix())
	if matched, isFound := matchRecordFromList(rows, clean); isFound {
		return &matched, nil
	}
	return nil, apperror.NewSimple(fmt.Sprintf("sequence or target %q not found in 24h cache", token), "E9041")
}

func cleanSequenceToken(token string) string {
	trimmed := strings.TrimSpace(token)
	trimmed = strings.TrimPrefix(trimmed, "#")
	return strings.TrimSpace(trimmed)
}

func matchRecordFromList(records []AgySequenceRecord, clean string) (AgySequenceRecord, bool) {
	if rec, isSeq := matchExactSeqID(records, clean); isSeq {
		return rec, true
	}
	for _, r := range records {
		if isFlexibleRecordMatch(r, clean) {
			return r, true
		}
	}
	return AgySequenceRecord{}, false
}

func matchExactSeqID(records []AgySequenceRecord, clean string) (AgySequenceRecord, bool) {
	for _, r := range records {
		if strings.EqualFold(r.SeqId, clean) {
			return r, true
		}
	}
	return AgySequenceRecord{}, false
}

func isFlexibleRecordMatch(r AgySequenceRecord, clean string) bool {
	if isIDOrConvMatch(r, clean) {
		return true
	}
	if strings.EqualFold(r.ProjectAlias, clean) {
		return true
	}
	return isPathTokenMatch(r.ProjectPath, clean)
}

func isIDOrConvMatch(r AgySequenceRecord, clean string) bool {
	if r.ProjectId != "" && (strings.EqualFold(r.ProjectId, clean) || strings.HasPrefix(strings.ToLower(r.ProjectId), strings.ToLower(clean))) {
		return true
	}
	return r.ConversationId != "" && (strings.EqualFold(r.ConversationId, clean) || strings.HasPrefix(strings.ToLower(r.ConversationId), strings.ToLower(clean)))
}

func isPathTokenMatch(projectPath, clean string) bool {
	if projectPath == "" {
		return false
	}
	normProj := filepath.Clean(projectPath)
	normTok := filepath.Clean(clean)
	return strings.EqualFold(normProj, normTok) || strings.EqualFold(filepath.Base(normProj), clean)
}
