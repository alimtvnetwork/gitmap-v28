package constants

// Workgroup and workspace table names (PascalCase).
const (
	TableWorkDirectory     = "WorkDirectory"
	TableFolderWorkDir     = "FolderWorkDir"
	TableFolderGroup       = "FolderGroup"
	TableSshKeyAssociation = "SshKeyAssociation"
	ViewWorkDirectories    = "work_directories"
)

// SSH Key association types.
const (
	AssocTypeRepo    = "repo"
	AssocTypeFolder  = "folder"
	AssocTypeGroup   = "group"
	AssocTypeWorkDir = "workdir"
	AssocTypeDefault = "default"
)

// WorkgroupCmdType defines strongly-typed workgroup subcommands.
type WorkgroupCmdType string

const (
	WorkgroupCmdList    WorkgroupCmdType = "list"
	WorkgroupCmdAdd     WorkgroupCmdType = "add"
	WorkgroupCmdBind    WorkgroupCmdType = "bind"
	WorkgroupCmdUnbind  WorkgroupCmdType = "unbind"
	WorkgroupCmdShow    WorkgroupCmdType = "show"
	WorkgroupCmdTestSSH WorkgroupCmdType = "test-ssh"
	WorkgroupCmdRemove  WorkgroupCmdType = "remove"
	WorkgroupCmdLs      WorkgroupCmdType = "ls"
	WorkgroupCmdTest    WorkgroupCmdType = "test"
	WorkgroupCmdRm      WorkgroupCmdType = "rm"
	WorkgroupCmdDelete  WorkgroupCmdType = "delete"
)

// WorkgroupFlagType defines strongly-typed CLI flags for workgroup commands.
type WorkgroupFlagType string

const (
	WorkgroupFlagSSHKey    WorkgroupFlagType = "--ssh-key"
	WorkgroupFlagDesc      WorkgroupFlagType = "--desc"
	WorkgroupFlagHost      WorkgroupFlagType = "--host"
	WorkgroupFlagUserName  WorkgroupFlagType = "--user-name"
	WorkgroupFlagUserEmail WorkgroupFlagType = "--user-email"
	WorkgroupFlagDefault   WorkgroupFlagType = "--default"
	WorkgroupFlagJSON      WorkgroupFlagType = "--json"
	WorkgroupFlagJSONShort WorkgroupFlagType = "-j"
	WorkgroupFlagHelp      WorkgroupFlagType = "--help"
	WorkgroupFlagHelpShort WorkgroupFlagType = "-h"
)

// SQL: create WorkDirectory table and indexes.
const (
	SQLCreateWorkDirectoryTable = `CREATE TABLE IF NOT EXISTS WorkDirectory (
	WorkDirectoryId    INTEGER PRIMARY KEY AUTOINCREMENT,
	AbsolutePath       TEXT NOT NULL UNIQUE COLLATE NOCASE,
	Label              TEXT NOT NULL DEFAULT '',
	IsDefault          INTEGER NOT NULL DEFAULT 0,
	HasActiveBinding   INTEGER NOT NULL DEFAULT 1,
	CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	SQLCreateWorkDirectoryAbsPathIndex = `CREATE UNIQUE INDEX IF NOT EXISTS IdxWorkDirectory_AbsolutePath 
	ON WorkDirectory(AbsolutePath COLLATE NOCASE);`

	SQLCreateWorkDirectoryIsDefaultIndex = `CREATE INDEX IF NOT EXISTS IdxWorkDirectory_IsDefault 
	ON WorkDirectory(IsDefault);`
)

// SQL: backward-compatibility view mapping legacy work_directories to WorkDirectory.
const SQLCreateWorkDirectoriesCompatView = `CREATE VIEW IF NOT EXISTS work_directories AS
SELECT
	WorkDirectoryId    AS id,
	AbsolutePath       AS absolute_path,
	Label              AS label,
	IsDefault          AS is_default,
	CreatedAt          AS created_at,
	UpdatedAt          AS updated_at
FROM WorkDirectory;`

// SQL: create FolderWorkDir table and indexes.
const (
	SQLCreateFolderWorkDirTable = `CREATE TABLE IF NOT EXISTS FolderWorkDir (
	FolderWorkDirId    INTEGER PRIMARY KEY AUTOINCREMENT,
	FolderPath         TEXT NOT NULL UNIQUE COLLATE NOCASE,
	WorkDirectoryId    INTEGER NOT NULL REFERENCES WorkDirectory(WorkDirectoryId) ON DELETE CASCADE,
	IsRecursive        INTEGER NOT NULL DEFAULT 1,
	Priority           INTEGER NOT NULL DEFAULT 100,
	IsActive           INTEGER NOT NULL DEFAULT 1,
	CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	SQLCreateFolderWorkDirFolderPathIndex = `CREATE UNIQUE INDEX IF NOT EXISTS IdxFolderWorkDir_FolderPath 
	ON FolderWorkDir(FolderPath COLLATE NOCASE);`

	SQLCreateFolderWorkDirWorkDirIdIndex = `CREATE INDEX IF NOT EXISTS IdxFolderWorkDir_WorkDirId 
	ON FolderWorkDir(WorkDirectoryId);`

	SQLCreateFolderWorkDirPriorityIndex = `CREATE INDEX IF NOT EXISTS IdxFolderWorkDir_Priority 
	ON FolderWorkDir(Priority DESC);`
)

// SQL: create FolderGroup table and indexes.
const (
	SQLCreateFolderGroupTable = `CREATE TABLE IF NOT EXISTS FolderGroup (
	FolderGroupId      INTEGER PRIMARY KEY AUTOINCREMENT,
	GroupName          TEXT NOT NULL COLLATE NOCASE,
	FolderPath         TEXT NOT NULL COLLATE NOCASE,
	Description        TEXT NOT NULL DEFAULT '',
	IsDefault          INTEGER NOT NULL DEFAULT 0,
	IsActive           INTEGER NOT NULL DEFAULT 1,
	CreatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt          TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (GroupName COLLATE NOCASE, FolderPath COLLATE NOCASE)
);`

	SQLCreateFolderGroupGroupNameIndex = `CREATE INDEX IF NOT EXISTS IdxFolderGroup_GroupName 
	ON FolderGroup(GroupName COLLATE NOCASE);`

	SQLCreateFolderGroupFolderPathIndex = `CREATE INDEX IF NOT EXISTS IdxFolderGroup_FolderPath 
	ON FolderGroup(FolderPath COLLATE NOCASE);`
)

// SQL: create SshKeyAssociation table and indexes.
const (
	SQLCreateSshKeyAssociationTable = `CREATE TABLE IF NOT EXISTS SshKeyAssociation (
	SshKeyAssociationId INTEGER PRIMARY KEY AUTOINCREMENT,
	SshKeyId            INTEGER NOT NULL REFERENCES SshKey(SshKeyId) ON DELETE CASCADE,
	AssociationType     TEXT NOT NULL,
	TargetIdentifier    TEXT NOT NULL COLLATE NOCASE,
	FolderPath          TEXT NOT NULL DEFAULT '' COLLATE NOCASE,
	GroupName           TEXT NOT NULL DEFAULT '' COLLATE NOCASE,
	WorkDirectoryId     INTEGER NULL REFERENCES WorkDirectory(WorkDirectoryId) ON DELETE SET NULL,
	RepoId              INTEGER NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	GitAuthorName       TEXT NOT NULL DEFAULT '',
	GitAuthorEmail      TEXT NOT NULL DEFAULT '',
	IsActive            INTEGER NOT NULL DEFAULT 1,
	IsAutoInjected      INTEGER NOT NULL DEFAULT 1,
	Priority            INTEGER NOT NULL DEFAULT 100,
	CreatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt           TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (AssociationType, TargetIdentifier COLLATE NOCASE, SshKeyId)
);`

	SQLCreateSshKeyAssocTargetKeyIndex = `CREATE UNIQUE INDEX IF NOT EXISTS IdxSshKeyAssoc_TargetKey 
	ON SshKeyAssociation(AssociationType, TargetIdentifier COLLATE NOCASE, SshKeyId);`

	SQLCreateSshKeyAssocSshKeyIdIndex = `CREATE INDEX IF NOT EXISTS IdxSshKeyAssoc_SshKeyId 
	ON SshKeyAssociation(SshKeyId);`

	SQLCreateSshKeyAssocFolderPathIndex = `CREATE INDEX IF NOT EXISTS IdxSshKeyAssoc_FolderPath 
	ON SshKeyAssociation(FolderPath COLLATE NOCASE);`

	SQLCreateSshKeyAssocGroupNameIndex = `CREATE INDEX IF NOT EXISTS IdxSshKeyAssoc_GroupName 
	ON SshKeyAssociation(GroupName COLLATE NOCASE);`

	SQLCreateSshKeyAssocWorkDirIdIndex = `CREATE INDEX IF NOT EXISTS IdxSshKeyAssoc_WorkDirId 
	ON SshKeyAssociation(WorkDirectoryId);`

	SQLCreateSshKeyAssocRepoIdIndex = `CREATE INDEX IF NOT EXISTS IdxSshKeyAssoc_RepoId 
	ON SshKeyAssociation(RepoId);`

	SQLCreateSshKeyAssocPriorityIndex = `CREATE INDEX IF NOT EXISTS IdxSshKeyAssoc_Priority 
	ON SshKeyAssociation(Priority DESC);`
)

// SQL: Drop statements for reset and rollback.
const (
	SQLDropSshKeyAssociationTable    = "DROP TABLE IF EXISTS SshKeyAssociation"
	SQLDropFolderWorkDirTable        = "DROP TABLE IF EXISTS FolderWorkDir"
	SQLDropFolderGroupTable          = "DROP TABLE IF EXISTS FolderGroup"
	SQLDropWorkDirectoriesCompatView = "DROP VIEW IF EXISTS work_directories"
	SQLDropWorkDirectoryTable        = "DROP TABLE IF EXISTS WorkDirectory"
)

// SQL: WorkDirectory queries.
const (
	SQLUpsertWorkDirectory = `INSERT INTO WorkDirectory (AbsolutePath, Label, IsDefault, HasActiveBinding, UpdatedAt)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(AbsolutePath) DO UPDATE SET
	Label = CASE WHEN excluded.Label != '' THEN excluded.Label ELSE Label END,
	IsDefault = CASE WHEN excluded.IsDefault != 0 THEN excluded.IsDefault ELSE IsDefault END,
	HasActiveBinding = CASE WHEN excluded.HasActiveBinding != 0 THEN excluded.HasActiveBinding ELSE HasActiveBinding END,
	UpdatedAt = CURRENT_TIMESTAMP;`

	SQLSelectAllWorkDirectories = `SELECT WorkDirectoryId, AbsolutePath, Label, IsDefault, HasActiveBinding, CreatedAt, UpdatedAt
FROM WorkDirectory
ORDER BY IsDefault DESC, AbsolutePath ASC;`

	SQLSelectWorkDirectoryByPath = `SELECT WorkDirectoryId, AbsolutePath, Label, IsDefault, HasActiveBinding, CreatedAt, UpdatedAt
FROM WorkDirectory
WHERE AbsolutePath = ? COLLATE NOCASE;`

	SQLSelectWorkDirectoryById = `SELECT WorkDirectoryId, AbsolutePath, Label, IsDefault, HasActiveBinding, CreatedAt, UpdatedAt
FROM WorkDirectory
WHERE WorkDirectoryId = ?;`

	SQLSetDefaultWorkDirectory = `UPDATE WorkDirectory 
SET IsDefault = CASE WHEN WorkDirectoryId = ? OR AbsolutePath = ? COLLATE NOCASE THEN 1 ELSE 0 END,
	UpdatedAt = CURRENT_TIMESTAMP;`

	SQLDeleteWorkDirectory = `DELETE FROM WorkDirectory 
WHERE WorkDirectoryId = ? OR AbsolutePath = ? COLLATE NOCASE;`

	SQLSelectDefaultWorkDirectory = `SELECT WorkDirectoryId, AbsolutePath, Label, IsDefault, HasActiveBinding, CreatedAt, UpdatedAt
FROM WorkDirectory
WHERE IsDefault = 1 AND HasActiveBinding = 1
ORDER BY WorkDirectoryId ASC
LIMIT 1;`

	SQLResolveWorkDirDirectPrefix = `SELECT WorkDirectoryId, AbsolutePath, Label, IsDefault, HasActiveBinding, CreatedAt, UpdatedAt
FROM WorkDirectory
WHERE HasActiveBinding = 1
  AND (
	? = AbsolutePath COLLATE NOCASE
	OR ? LIKE AbsolutePath || '/%'
	OR ? LIKE AbsolutePath || '\%'
  )
ORDER BY LENGTH(AbsolutePath) DESC
LIMIT 1;`
)

// SQL: FolderWorkDir queries.
const (
	SQLUpsertFolderWorkDir = `INSERT INTO FolderWorkDir (FolderPath, WorkDirectoryId, IsRecursive, Priority, IsActive, UpdatedAt)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(FolderPath) DO UPDATE SET
	WorkDirectoryId = excluded.WorkDirectoryId,
	IsRecursive = excluded.IsRecursive,
	Priority = excluded.Priority,
	IsActive = excluded.IsActive,
	UpdatedAt = CURRENT_TIMESTAMP;`

	SQLDeleteFolderWorkDir = `DELETE FROM FolderWorkDir WHERE FolderPath = ? COLLATE NOCASE;`

	SQLSelectAllFolderWorkDirs = `SELECT FolderWorkDirId, FolderPath, WorkDirectoryId, IsRecursive, Priority, IsActive, CreatedAt, UpdatedAt
FROM FolderWorkDir
ORDER BY Priority DESC, FolderPath ASC;`

	SQLResolveWorkDirectoryForPath = `SELECT f.FolderWorkDirId, f.FolderPath, f.WorkDirectoryId, f.IsRecursive, f.Priority, f.IsActive,
       w.WorkDirectoryId, w.AbsolutePath, w.Label, w.IsDefault, w.HasActiveBinding, w.CreatedAt, w.UpdatedAt
FROM FolderWorkDir f
JOIN WorkDirectory w ON f.WorkDirectoryId = w.WorkDirectoryId
WHERE f.IsActive = 1
  AND (
	? = f.FolderPath COLLATE NOCASE
	OR (f.IsRecursive = 1 AND ? LIKE f.FolderPath || '/%')
	OR (f.IsRecursive = 1 AND ? LIKE f.FolderPath || '\%')
  )
ORDER BY LENGTH(f.FolderPath) DESC, f.Priority DESC
LIMIT 1;`
)

// SQL: FolderGroup queries.
const (
	SQLUpsertFolderGroup = `INSERT INTO FolderGroup (GroupName, FolderPath, Description, IsDefault, IsActive, UpdatedAt)
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(GroupName, FolderPath) DO UPDATE SET
	Description = CASE WHEN excluded.Description != '' THEN excluded.Description ELSE Description END,
	IsDefault = excluded.IsDefault,
	IsActive = excluded.IsActive,
	UpdatedAt = CURRENT_TIMESTAMP;`

	SQLDeleteFolderFromGroup = `DELETE FROM FolderGroup WHERE GroupName = ? COLLATE NOCASE AND FolderPath = ? COLLATE NOCASE;`

	SQLSelectFoldersInGroup = `SELECT FolderGroupId, GroupName, FolderPath, Description, IsDefault, IsActive, CreatedAt, UpdatedAt
FROM FolderGroup
WHERE GroupName = ? COLLATE NOCASE AND IsActive = 1
ORDER BY FolderPath ASC;`

	SQLSelectActiveGroupsSummary = `SELECT GroupName, COUNT(FolderPath) AS FolderCount, MIN(CreatedAt) AS CreatedAt
FROM FolderGroup
WHERE IsActive = 1
GROUP BY GroupName COLLATE NOCASE
ORDER BY GroupName ASC;`
)

// SQL: SshKeyAssociation queries.
const (
	SQLUpsertSshKeyAssociation = `INSERT INTO SshKeyAssociation (
	SshKeyId, AssociationType, TargetIdentifier, FolderPath, GroupName,
	WorkDirectoryId, RepoId, GitAuthorName, GitAuthorEmail, IsActive, IsAutoInjected, Priority, UpdatedAt
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(AssociationType, TargetIdentifier, SshKeyId) DO UPDATE SET
	FolderPath = excluded.FolderPath,
	GroupName = excluded.GroupName,
	WorkDirectoryId = excluded.WorkDirectoryId,
	RepoId = excluded.RepoId,
	GitAuthorName = CASE WHEN excluded.GitAuthorName != '' THEN excluded.GitAuthorName ELSE GitAuthorName END,
	GitAuthorEmail = CASE WHEN excluded.GitAuthorEmail != '' THEN excluded.GitAuthorEmail ELSE GitAuthorEmail END,
	IsActive = excluded.IsActive,
	IsAutoInjected = excluded.IsAutoInjected,
	Priority = excluded.Priority,
	UpdatedAt = CURRENT_TIMESTAMP;`

	SQLDeleteSshKeyAssociation = `DELETE FROM SshKeyAssociation WHERE SshKeyAssociationId = ?;`

	SQLSelectSshKeyAssociationsByType = `SELECT SshKeyAssociationId, SshKeyId, AssociationType, TargetIdentifier, FolderPath, GroupName,
       WorkDirectoryId, RepoId, GitAuthorName, GitAuthorEmail, IsActive, IsAutoInjected, Priority, CreatedAt, UpdatedAt
FROM SshKeyAssociation
WHERE AssociationType = ? AND IsActive = 1
ORDER BY Priority DESC;`

	SQLSelectAllSshKeyAssociations = `SELECT SshKeyAssociationId, SshKeyId, AssociationType, TargetIdentifier, FolderPath, GroupName,
       WorkDirectoryId, RepoId, GitAuthorName, GitAuthorEmail, IsActive, IsAutoInjected, Priority, CreatedAt, UpdatedAt
FROM SshKeyAssociation
ORDER BY Priority DESC;`

	SQLResolveEffectiveCredentials = `SELECT a.SshKeyAssociationId, a.AssociationType, a.TargetIdentifier,
       k.SshKeyId, k.Name AS KeyName, k.PrivatePath, k.PublicKey, k.Fingerprint,
       a.GitAuthorName, a.GitAuthorEmail, a.IsAutoInjected, a.Priority
FROM SshKeyAssociation a
JOIN SshKey k ON a.SshKeyId = k.SshKeyId
WHERE a.IsActive = 1
  AND (
	(a.AssociationType = 'repo' AND a.RepoId = ?)
	OR
	(a.AssociationType = 'folder' AND (
		? = a.FolderPath COLLATE NOCASE
		OR ? LIKE a.FolderPath || '/%'
		OR ? LIKE a.FolderPath || '\%'
	))
	OR
	(a.AssociationType = 'group' AND a.GroupName IN (
		SELECT g.GroupName FROM FolderGroup g 
		WHERE g.IsActive = 1 AND (
			? = g.FolderPath COLLATE NOCASE 
			OR ? LIKE g.FolderPath || '/%' 
			OR ? LIKE g.FolderPath || '\%'
		)
	))
	OR
	(a.AssociationType = 'workdir' AND a.WorkDirectoryId = ?)
  )
ORDER BY 
	CASE a.AssociationType
		WHEN 'repo' THEN 1
		WHEN 'folder' THEN 2
		WHEN 'group' THEN 3
		WHEN 'workdir' THEN 4
		ELSE 5
	END ASC,
	LENGTH(a.FolderPath) DESC,
	a.Priority DESC
LIMIT 1;`

	SQLSelectDefaultSSHKeySimple = `SELECT SshKeyId, Name, PrivatePath, PublicKey, Fingerprint
FROM SshKey
WHERE IsDefault = 1
ORDER BY SshKeyId ASC
LIMIT 1;`
)
