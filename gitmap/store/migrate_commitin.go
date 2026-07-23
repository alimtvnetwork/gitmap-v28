package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// commitInDDL returns the ordered list of CREATE / INDEX / SEED
// statements for the commit-in feature. Order matters:
//
//  1. Enum-mirror tables FIRST (FK targets).
//  2. Profile (FK target for CommitInRun.ProfileId).
//  3. Profile children.
//  4. CommitInRun, then InputRepo, then SourceCommit, then file/
//     rewrite/skip/sha-map tables (each FK back into the previous).
//  5. All `INSERT OR IGNORE` seeds (idempotent).
//
// Every statement is independently idempotent: re-running Migrate()
// against an already-migrated DB MUST be a no-op.
func commitInDDL() []string {
	return []string{
		// Enum mirrors.
		constants.SQLCreateCommitInRunStatus,
		constants.SQLCreateCommitInInputKind,
		constants.SQLCreateCommitInOutcome,
		constants.SQLCreateCommitInSkipReason,
		constants.SQLCreateCommitInExclusionKind,
		constants.SQLCreateCommitInMessageRuleKind,
		constants.SQLCreateCommitInLanguage,
		constants.SQLCreateCommitInConflictMode,
		// Profile + children.
		constants.SQLCreateCommitInProfile,
		constants.SQLCreateCommitInProfileDefaultIdx,
		constants.SQLCreateCommitInProfileExclusion,
		constants.SQLCreateCommitInProfileMessageRule,
		// Run + commit chain.
		constants.SQLCreateCommitInRun,
		constants.SQLCreateCommitInRunSourceIdx,
		constants.SQLCreateCommitInInputRepo,
		constants.SQLCreateCommitInSourceCommit,
		constants.SQLCreateCommitInSourceCommitShaIdx,
		constants.SQLCreateCommitInSourceCommitFile,
		constants.SQLCreateCommitInRewritten,
		constants.SQLCreateCommitInRewrittenShaIdx,
		constants.SQLCreateCommitInSkipLog,
		constants.SQLCreateCommitInShaMap,
		constants.SQLCreateCommitInShaMapIdx,
		// Migration 007 — tag replay map (spec §09).
		constants.SQLCreateCommitInTagOutcome,
		constants.SQLCreateCommitInReplayMap,
		constants.SQLCreateCommitInReplayMapTagNameIdx,
		constants.SQLCreateCommitInReplayMapDestShaIdx,
		constants.SQLCreateCommitInReplayMapBranchIdx,
		// Enum-mirror seeds.
		constants.SQLSeedCommitInRunStatus,
		constants.SQLSeedCommitInInputKind,
		constants.SQLSeedCommitInOutcome,
		constants.SQLSeedCommitInSkipReason,
		constants.SQLSeedCommitInExclusionKind,
		constants.SQLSeedCommitInMessageRuleKind,
		constants.SQLSeedCommitInLanguage,
		constants.SQLSeedCommitInConflictMode,
		constants.SQLSeedCommitInTagOutcome,
	}
}

// migrateCommitIn applies the commit-in DDL bundle. Failures abort the
// whole Migrate() pipeline so the schema-version marker is NOT stamped,
// guaranteeing the next run retries.
func (db *DB) migrateCommitIn() error {
	for _, stmt := range commitInDDL() {
		if _, err := db.conn.Exec(stmt); err != nil {
			return fmt.Errorf("commit-in migration failed: %w (statement: %s)", err, stmt)
		}
	}
	return nil
}
