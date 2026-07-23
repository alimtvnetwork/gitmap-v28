// Package store — migrate_v15phase3.go performs the Phase 1.3 v15 renames:
//
//	Amendments     → Amendment      (Id → AmendmentId)
//	CommitTemplates → CommitTemplate (Id → CommitTemplateId)
//	Settings       → Setting        (Key PK preserved — no rename needed,
//	                                   but the table is renamed to singular)
//	SSHKeys        → SshKey         (Id → SshKeyId; abbreviation fix)
//	InstalledTools → InstalledTool  (Id → InstalledToolId)
//	TempReleases   → TempRelease    (Id → TempReleaseId)
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// migrateV15Phase3 runs all Phase 1.3 table rebuilds.
func (db *DB) migrateV15Phase3() error {
	specs := []v15RebuildSpec{
		{
			OldTable:      "Amendments",
			NewTable:      "Amendment",
			NewCreateSQL:  constants.SQLCreateAmendment,
			OldColumnList: "Id, Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed, CreatedAt",
			NewColumnList: "AmendmentId, Branch, FromCommit, ToCommit, TotalCommits, PreviousName, PreviousEmail, NewName, NewEmail, Mode, ForcePushed, CreatedAt",
			StartMsg:      "→ Migrating Amendments → Amendment (AmendmentId PK)...",
			DoneMsg:       "✓ Migrated Amendments → Amendment.",
		},
		{
			OldTable:      "CommitTemplates",
			NewTable:      "CommitTemplate",
			NewCreateSQL:  constants.SQLCreateCommitTemplate,
			OldColumnList: "Id, Kind, Template, CreatedAt",
			NewColumnList: "CommitTemplateId, Kind, Template, CreatedAt",
			StartMsg:      "→ Migrating CommitTemplates → CommitTemplate (CommitTemplateId PK)...",
			DoneMsg:       "✓ Migrated CommitTemplates → CommitTemplate.",
		},
		{
			OldTable:      "Settings",
			NewTable:      "Setting",
			NewCreateSQL:  constants.SQLCreateSetting,
			OldColumnList: "Key, Value",
			NewColumnList: "Key, Value",
			StartMsg:      "→ Migrating Settings → Setting...",
			DoneMsg:       "✓ Migrated Settings → Setting.",
		},
		{
			OldTable:      "SSHKeys",
			NewTable:      "SshKey",
			NewCreateSQL:  constants.SQLCreateSshKey,
			OldColumnList: "Id, Name, PrivatePath, PublicKey, Fingerprint, Email, CreatedAt",
			NewColumnList: "SshKeyId, Name, PrivatePath, PublicKey, Fingerprint, Email, CreatedAt",
			StartMsg:      "→ Migrating SSHKeys → SshKey (SshKeyId PK; v15 abbreviation fix)...",
			DoneMsg:       "✓ Migrated SSHKeys → SshKey.",
		},
		{
			OldTable:      "InstalledTools",
			NewTable:      "InstalledTool",
			NewCreateSQL:  constants.SQLCreateInstalledTool,
			OldColumnList: "Id, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt",
			NewColumnList: "InstalledToolId, Tool, VersionMajor, VersionMinor, VersionPatch, VersionBuild, VersionString, PackageManager, InstallPath, InstalledAt, UpdatedAt",
			StartMsg:      "→ Migrating InstalledTools → InstalledTool (InstalledToolId PK)...",
			DoneMsg:       "✓ Migrated InstalledTools → InstalledTool.",
		},
		{
			OldTable:      "TempReleases",
			NewTable:      "TempRelease",
			NewCreateSQL:  constants.SQLCreateTempRelease,
			OldColumnList: "Id, Branch, VersionPrefix, SequenceNumber, CommitSha, CommitMessage, CreatedAt",
			NewColumnList: "TempReleaseId, Branch, VersionPrefix, SequenceNumber, CommitSha, CommitMessage, CreatedAt",
			StartMsg:      "→ Migrating TempReleases → TempRelease (TempReleaseId PK)...",
			DoneMsg:       "✓ Migrated TempReleases → TempRelease.",
		},
	}

	for _, spec := range specs {
		if err := db.runV15Rebuild(spec); err != nil {
			return fmt.Errorf("phase 1.3 %s: %w", spec.OldTable, err)
		}
	}

	return nil
}
