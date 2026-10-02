package store

import (
	"context"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// DuplicateRepoGroup represents a collection of repository records sharing an identical remote URL.
type DuplicateRepoGroup struct {
	CleanKey   string
	Count      int
	Keeper     model.ScanRecord
	Duplicates []model.ScanRecord
}

// DeduplicationSummary captures the aggregate results of a database repository deduplication run.
type DeduplicationSummary struct {
	GroupsFound   int
	RowsPurged    int64
	KeeperRepoIDs []int64
	PurgedRepoIDs []int64
}

// cleanRemoteKey computes the normalized clean remote URL matching the SQLite grouping query.
func cleanRemoteKey(httpsURL, sshURL string) string {
	raw := httpsURL
	if raw == "" {
		raw = sshURL
	}
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, ".git", "")
	s = strings.TrimRight(s, "/")
	return s
}

// FindDuplicateRepos queries the SQLite database for duplicate repositories grouped by clean remote URL.
// It executes strictly against the database without touching the filesystem.
func (db *DB) FindDuplicateRepos() ([]DuplicateRepoGroup, error) {
	return findDuplicateReposRunner(db.conn)
}

func findDuplicateReposRunner(runner sqlQueryer) ([]DuplicateRepoGroup, error) {
	rows, err := QueryWrapper(runner, constants.SQLSelectDuplicateRepoRows).Destruct()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	groupMap := make(map[string][]model.ScanRecord)
	var groupKeys []string

	for _, r := range records {
		key := cleanRemoteKey(r.HTTPSUrl, r.SSHUrl)
		if key == "" {
			continue
		}
		if _, exists := groupMap[key]; !exists {
			groupKeys = append(groupKeys, key)
		}
		groupMap[key] = append(groupMap[key], r)
	}

	var groups []DuplicateRepoGroup
	for _, key := range groupKeys {
		list := groupMap[key]
		if len(list) < 2 {
			continue
		}

		keeper := list[0]
		dups := make([]model.ScanRecord, len(list)-1)
		copy(dups, list[1:])

		groups = append(groups, DuplicateRepoGroup{
			CleanKey:   key,
			Count:      len(list),
			Keeper:     keeper,
			Duplicates: dups,
		})
	}

	return groups, nil
}

// DeduplicateRepos transactionally prunes redundant repositories in gitmap.db.
// Child foreign keys (Release, GroupRepo, VersionProbe) are safely remapped to the keeper record.
func (db *DB) DeduplicateRepos(keepNewest bool) (*DeduplicationSummary, error) {
	summary := &DeduplicationSummary{}

	hasRelease := isTablePresent(db.conn, "Release")
	hasGroupRepo := isTablePresent(db.conn, "GroupRepo")
	hasVersionProbe := isTablePresent(db.conn, "VersionProbe")

	wrap, appErr := dbengine.WrapDb(db.conn, dbengine.DbSQLite)
	if appErr != nil {
		return nil, appErr
	}

	appErr = wrap.WithTransaction(context.Background(), func(tx *dbengine.TxWrapper) *apperror.AppError {
		groups, err := findDuplicateReposRunner(tx.Tx())
		if err != nil {
			return apperror.WrapSimple(err, "findDuplicateReposRunner")
		}

		if len(groups) == 0 {
			return nil
		}

		for _, grp := range groups {
			allRecords := append([]model.ScanRecord{grp.Keeper}, grp.Duplicates...)
			var keeper model.ScanRecord
			var duplicates []model.ScanRecord

			if keepNewest {
				maxIdx := 0
				for i, rec := range allRecords {
					if rec.ID > allRecords[maxIdx].ID {
						maxIdx = i
					}
				}
				keeper = allRecords[maxIdx]
				for i, rec := range allRecords {
					if i != maxIdx {
						duplicates = append(duplicates, rec)
					}
				}
			} else {
				minIdx := 0
				for i, rec := range allRecords {
					if rec.ID < allRecords[minIdx].ID {
						minIdx = i
					}
				}
				keeper = allRecords[minIdx]
				for i, rec := range allRecords {
					if i != minIdx {
						duplicates = append(duplicates, rec)
					}
				}
			}

			summary.KeeperRepoIDs = append(summary.KeeperRepoIDs, keeper.ID)

			for _, dup := range duplicates {
				dupID := dup.ID

				// 1. Remap Release.RepoId (handle UNIQUE(RepoId, Tag))
				if hasRelease {
					if _, err := tx.Tx().Exec("UPDATE OR IGNORE Release SET RepoId = ? WHERE RepoId = ?", keeper.ID, dupID); err != nil {
						return apperror.WrapSimple(err, "remapRelease")
					}
					if _, err := tx.Tx().Exec("DELETE FROM Release WHERE RepoId = ?", dupID); err != nil {
						return apperror.WrapSimple(err, "deleteReleaseDups")
					}
				}

				// 2. Remap GroupRepo.RepoId (handle PRIMARY KEY (GroupId, RepoId))
				if hasGroupRepo {
					if _, err := tx.Tx().Exec("INSERT OR IGNORE INTO GroupRepo (GroupId, RepoId) SELECT GroupId, ? FROM GroupRepo WHERE RepoId = ?", keeper.ID, dupID); err != nil {
						return apperror.WrapSimple(err, "remapGroupRepo")
					}
					if _, err := tx.Tx().Exec("DELETE FROM GroupRepo WHERE RepoId = ?", dupID); err != nil {
						return apperror.WrapSimple(err, "deleteGroupRepoDups")
					}
				}

				// 3. Remap VersionProbe.RepoId
				if hasVersionProbe {
					if _, err := tx.Tx().Exec("UPDATE VersionProbe SET RepoId = ? WHERE RepoId = ?", keeper.ID, dupID); err != nil {
						return apperror.WrapSimple(err, "remapVersionProbe")
					}
				}

				// 4. Purge duplicate Repo row
				res, err := tx.Tx().Exec("DELETE FROM Repo WHERE RepoId = ?", dupID)
				if err != nil {
					return apperror.WrapSimple(err, "deleteDuplicateRepo")
				}

				rows, _ := res.RowsAffected()
				summary.RowsPurged += rows
				summary.PurgedRepoIDs = append(summary.PurgedRepoIDs, dupID)
			}

			summary.GroupsFound++
		}

		return nil
	})

	if appErr != nil {
		return nil, appErr
	}

	return summary, nil
}
