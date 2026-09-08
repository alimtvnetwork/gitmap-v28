package constants

// SQL statements for SplitDatabaseRegistry.
const (
	SQLCreateSplitDatabaseRegistry = `CREATE TABLE IF NOT EXISTS SplitDatabaseRegistry (
    SplitDatabaseRegistryId INTEGER PRIMARY KEY AUTOINCREMENT,
    DatabaseType            TEXT NOT NULL,
    DatabaseKey             TEXT NOT NULL,
    DatabasePath            TEXT NOT NULL,
    SizeBytes               INTEGER NOT NULL DEFAULT 0,
    TableCount              INTEGER NOT NULL DEFAULT 0,
    RecordCount             INTEGER NOT NULL DEFAULT 0,
    SchemaVersion           INTEGER NOT NULL DEFAULT 1,
    Status                  TEXT NOT NULL DEFAULT 'active',
    IsActive                INTEGER NOT NULL DEFAULT 1,
    IsAttached              INTEGER NOT NULL DEFAULT 0,
    Description             TEXT NULL,
    Notes                   TEXT NULL,
    Comments                TEXT NULL,
    LastAccessedAt          INTEGER NOT NULL DEFAULT 0,
    LastSyncedAt            INTEGER NOT NULL DEFAULT 0,
    CreatedAt               INTEGER NOT NULL DEFAULT (unixepoch()),
    UpdatedAt               INTEGER NOT NULL DEFAULT (unixepoch()),
    UNIQUE(DatabaseType, DatabaseKey)
);`

	SQLCreateSplitDatabaseRegistryTypeIndex = `CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_Type ON SplitDatabaseRegistry(DatabaseType);`

	SQLCreateSplitDatabaseRegistryStatusIndex = `CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_Status ON SplitDatabaseRegistry(Status);`

	SQLCreateSplitDatabaseRegistryUpdatedAtIndex = `CREATE INDEX IF NOT EXISTS IdxSplitDatabaseRegistry_UpdatedAt ON SplitDatabaseRegistry(UpdatedAt);`

	SQLUpsertSplitDatabaseRegistry = `INSERT INTO SplitDatabaseRegistry (
    DatabaseType, DatabaseKey, DatabasePath, SizeBytes, TableCount,
    RecordCount, SchemaVersion, Status, IsActive, IsAttached,
    Description, Notes, Comments, LastAccessedAt, LastSyncedAt,
    CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(DatabaseType, DatabaseKey) DO UPDATE SET
    DatabasePath = excluded.DatabasePath,
    SizeBytes = excluded.SizeBytes,
    TableCount = excluded.TableCount,
    RecordCount = excluded.RecordCount,
    SchemaVersion = excluded.SchemaVersion,
    Status = excluded.Status,
    IsActive = excluded.IsActive,
    IsAttached = excluded.IsAttached,
    Description = CASE WHEN excluded.Description IS NOT NULL AND excluded.Description != '' THEN excluded.Description ELSE SplitDatabaseRegistry.Description END,
    Notes = CASE WHEN excluded.Notes IS NOT NULL AND excluded.Notes != '' THEN excluded.Notes ELSE SplitDatabaseRegistry.Notes END,
    Comments = CASE WHEN excluded.Comments IS NOT NULL AND excluded.Comments != '' THEN excluded.Comments ELSE SplitDatabaseRegistry.Comments END,
    LastAccessedAt = CASE WHEN excluded.LastAccessedAt != 0 THEN excluded.LastAccessedAt ELSE SplitDatabaseRegistry.LastAccessedAt END,
    LastSyncedAt = excluded.LastSyncedAt,
    UpdatedAt = excluded.UpdatedAt;`

	SQLSelectSplitDB = `SELECT SplitDatabaseRegistryId, DatabaseType, DatabaseKey, DatabasePath,
    SizeBytes, TableCount, RecordCount, SchemaVersion, Status,
    IsActive, IsAttached, Description, Notes, Comments,
    LastAccessedAt, LastSyncedAt, CreatedAt, UpdatedAt
FROM SplitDatabaseRegistry
WHERE DatabaseType = ? AND DatabaseKey = ?;`

	SQLSelectListSplitDBByType = `SELECT SplitDatabaseRegistryId, DatabaseType, DatabaseKey, DatabasePath,
    SizeBytes, TableCount, RecordCount, SchemaVersion, Status,
    IsActive, IsAttached, Description, Notes, Comments,
    LastAccessedAt, LastSyncedAt, CreatedAt, UpdatedAt
FROM SplitDatabaseRegistry
WHERE DatabaseType = ?
ORDER BY DatabaseType, DatabaseKey;`

	SQLSelectListAllSplitDB = `SELECT SplitDatabaseRegistryId, DatabaseType, DatabaseKey, DatabasePath,
    SizeBytes, TableCount, RecordCount, SchemaVersion, Status,
    IsActive, IsAttached, Description, Notes, Comments,
    LastAccessedAt, LastSyncedAt, CreatedAt, UpdatedAt
FROM SplitDatabaseRegistry
ORDER BY DatabaseType, DatabaseKey;`
)
