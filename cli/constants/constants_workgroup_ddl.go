package constants

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
