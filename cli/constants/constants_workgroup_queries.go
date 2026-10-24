package constants

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
