-- Repos, Groups, Aliases, Releases, and Workspaces Schema and Parameterized Queries

-- Schema: Repo Table
CREATE TABLE IF NOT EXISTS Repo (
    RepoId           INTEGER PRIMARY KEY AUTOINCREMENT,
    Slug             TEXT NOT NULL,
    RepoName         TEXT NOT NULL,
    HttpsUrl         TEXT NOT NULL,
    SshUrl           TEXT NOT NULL,
    Branch           TEXT NOT NULL,
    RelativePath     TEXT NOT NULL,
    AbsolutePath     TEXT NOT NULL,
    CloneInstruction TEXT NOT NULL,
    Notes            TEXT DEFAULT '',
    LastInjectedDesktopAt TEXT DEFAULT '',
    LastInjectedVSCodeAt  TEXT DEFAULT '',
    LastClonedAt     TEXT DEFAULT '',
    IdentifiedTransport TEXT NOT NULL DEFAULT '',
    CreatedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt        TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxRepo_AbsolutePath ON Repo(AbsolutePath);

-- Schema: Group Table
CREATE TABLE IF NOT EXISTS "Group" (
    GroupId     INTEGER PRIMARY KEY AUTOINCREMENT,
    Name        TEXT NOT NULL UNIQUE,
    Description TEXT DEFAULT '',
    Color       TEXT DEFAULT '',
    CreatedAt   TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Schema: GroupRepo Join Table
CREATE TABLE IF NOT EXISTS GroupRepo (
    GroupId INTEGER NOT NULL REFERENCES "Group"(GroupId) ON DELETE CASCADE,
    RepoId  INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
    PRIMARY KEY (GroupId, RepoId)
);

-- Schema: Alias Table
CREATE TABLE IF NOT EXISTS Alias (
    AliasId   INTEGER PRIMARY KEY AUTOINCREMENT,
    Alias     TEXT NOT NULL UNIQUE,
    RepoId    INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
    IsPrimary INTEGER DEFAULT 0,
    Source    TEXT DEFAULT 'manual',
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxAlias_RepoId ON Alias(RepoId);

-- Schema: Release Table
CREATE TABLE IF NOT EXISTS Release (
    ReleaseId    INTEGER PRIMARY KEY AUTOINCREMENT,
    RepoId       INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
    Version      TEXT NOT NULL,
    Tag          TEXT NOT NULL,
    Branch       TEXT NOT NULL,
    SourceBranch TEXT NOT NULL,
    CommitSha    TEXT NOT NULL,
    Changelog    TEXT DEFAULT '',
    Notes        TEXT DEFAULT '',
    IsDraft      INTEGER DEFAULT 0,
    IsPreRelease INTEGER DEFAULT 0,
    IsLatest     INTEGER DEFAULT 0,
    Source       TEXT DEFAULT 'release',
    CreatedAt    TEXT DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (RepoId, Tag)
);
CREATE INDEX IF NOT EXISTS IdxRelease_RepoId ON Release(RepoId);

-- Schema: WorkDirectory Table
CREATE TABLE IF NOT EXISTS WorkDirectory (
    WorkDirectoryId INTEGER PRIMARY KEY AUTOINCREMENT,
    AbsolutePath    TEXT UNIQUE NOT NULL,
    Label           TEXT,
    IsDefault       INTEGER DEFAULT 0,
    CreatedAt       TEXT DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt       TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Query: Upsert Repo
INSERT INTO Repo (Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(AbsolutePath) DO UPDATE SET
    Slug=excluded.Slug, RepoName=excluded.RepoName, HttpsUrl=excluded.HttpsUrl,
    SshUrl=excluded.SshUrl, Branch=excluded.Branch, RelativePath=excluded.RelativePath,
    CloneInstruction=excluded.CloneInstruction, Notes=excluded.Notes,
    IdentifiedTransport=excluded.IdentifiedTransport, UpdatedAt=CURRENT_TIMESTAMP;

-- Query: Select All Repos
SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport
FROM Repo ORDER BY Slug;

-- Query: Select Repo By ID
SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport
FROM Repo WHERE RepoId = ?;

-- Query: Select Repo By Slug
SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport
FROM Repo WHERE Slug = ?;

-- Query: Select Repo By Absolute Path
SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport
FROM Repo WHERE AbsolutePath = ?;

-- Query: Suggest Repo By Slug Or Name
SELECT Slug FROM Repo WHERE Slug LIKE ? OR RepoName LIKE ? ORDER BY Slug ASC LIMIT 10;

-- Query: Delete Repo By Path
DELETE FROM Repo WHERE AbsolutePath = ?;

-- Query: Delete Repo By Slug
DELETE FROM Repo WHERE Slug = ?;

-- Query: Delete Repo By ID
DELETE FROM Repo WHERE RepoId = ?;

-- Query: Insert Group
INSERT INTO "Group" (Name, Description, Color) VALUES (?, ?, ?);

-- Query: Select All Groups
SELECT GroupId, Name, Description, Color, CreatedAt FROM "Group" ORDER BY Name;

-- Query: Select Group By Name
SELECT GroupId, Name, Description, Color, CreatedAt FROM "Group" WHERE Name = ?;

-- Query: Select Group By ID
SELECT GroupId, Name, Description, Color, CreatedAt FROM "Group" WHERE GroupId = ?;

-- Query: Delete Group
DELETE FROM "Group" WHERE Name = ?;

-- Query: Insert GroupRepo
INSERT OR IGNORE INTO GroupRepo (GroupId, RepoId) VALUES (?, ?);

-- Query: Delete GroupRepo
DELETE FROM GroupRepo WHERE GroupId = ? AND RepoId = ?;

-- Query: Select Group Repos
SELECT r.RepoId, r.Slug, r.RepoName, r.HttpsUrl, r.SshUrl, r.Branch,
    r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes, r.IdentifiedTransport
FROM Repo r JOIN GroupRepo gr ON r.RepoId = gr.RepoId
WHERE gr.GroupId = ? ORDER BY r.Slug;

-- Query: Count Group Repos
SELECT COUNT(*) FROM GroupRepo WHERE GroupId = ?;

-- Query: Insert Alias
INSERT INTO Alias (Alias, RepoId) VALUES (?, ?);

-- Query: Insert Alias Full
INSERT INTO Alias (Alias, RepoId, IsPrimary, Source) VALUES (?, ?, ?, ?);

-- Query: Update Alias
UPDATE Alias SET RepoId = ? WHERE Alias = ?;

-- Query: Select All Aliases
SELECT a.AliasId, a.Alias, a.RepoId, a.IsPrimary, a.Source, a.CreatedAt
FROM Alias a ORDER BY a.Alias;

-- Query: Select Alias By Name
SELECT a.AliasId, a.Alias, a.RepoId, a.IsPrimary, a.Source, a.CreatedAt
FROM Alias a WHERE a.Alias = ?;

-- Query: Select Aliases By Repo ID
SELECT a.AliasId, a.Alias, a.RepoId, COALESCE(a.IsPrimary, 0), COALESCE(a.Source, 'manual'), a.CreatedAt
FROM Alias a WHERE a.RepoId = ? ORDER BY a.IsPrimary DESC, a.AliasId ASC;

-- Query: Delete Alias
DELETE FROM Alias WHERE Alias = ?;

-- Query: Select Alias With Repo
SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt, r.AbsolutePath, r.Slug
FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId WHERE a.Alias = ?;

-- Query: Select All Aliases With Repo
SELECT a.AliasId, a.Alias, a.RepoId, a.CreatedAt, r.AbsolutePath, r.Slug
FROM Alias a JOIN Repo r ON a.RepoId = r.RepoId ORDER BY a.Alias;

-- Query: Select Unaliased Repos
SELECT r.RepoId, r.Slug, r.RepoName
FROM Repo r LEFT JOIN Alias a ON r.RepoId = a.RepoId
WHERE a.AliasId IS NULL ORDER BY r.Slug;

-- Query: Upsert Release
INSERT INTO Release (RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RepoId, Tag) DO UPDATE SET
    Version=excluded.Version, Branch=excluded.Branch, SourceBranch=excluded.SourceBranch,
    CommitSha=excluded.CommitSha, Changelog=excluded.Changelog, Notes=excluded.Notes,
    IsDraft=excluded.IsDraft, IsPreRelease=excluded.IsPreRelease,
    IsLatest=excluded.IsLatest, Source=excluded.Source;

-- Query: Select All Releases
SELECT ReleaseId, RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt
FROM Release ORDER BY CreatedAt DESC;

-- Query: Select All Releases Across Repos
SELECT R.ReleaseId, R.RepoId, P.Slug, P.AbsolutePath, R.Version, R.Tag, R.Branch, R.CommitSha, R.Source, R.IsDraft, R.IsLatest, R.IsPreRelease, R.CreatedAt
FROM Release R JOIN Repo P ON P.RepoId = R.RepoId ORDER BY R.CreatedAt DESC, P.Slug ASC;

-- Query: Select Release By Tag
SELECT ReleaseId, RepoId, Version, Tag, Branch, SourceBranch, CommitSha, Changelog, Notes, IsDraft, IsPreRelease, IsLatest, Source, CreatedAt
FROM Release WHERE Tag = ?;

-- Query: Clear Latest Release
UPDATE Release SET IsLatest = 0 WHERE IsLatest = 1 AND RepoId = ?;

-- Query: Upsert Work Directory
INSERT INTO WorkDirectory (AbsolutePath, Label, IsDefault, UpdatedAt)
VALUES (?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT(AbsolutePath) DO UPDATE SET
    Label = CASE WHEN ? != '' THEN ? ELSE Label END,
    IsDefault = CASE WHEN ? != 0 THEN ? ELSE IsDefault END,
    UpdatedAt = CURRENT_TIMESTAMP;

-- Query: Select All Work Directories
SELECT WorkDirectoryId, AbsolutePath, Label, IsDefault, CreatedAt, UpdatedAt
FROM WorkDirectory ORDER BY IsDefault DESC, UpdatedAt DESC;

-- Query: Set Default Work Directory
UPDATE WorkDirectory SET IsDefault = CASE WHEN WorkDirectoryId = ? OR AbsolutePath = ? THEN 1 ELSE 0 END;

-- Query: Delete Work Directory
DELETE FROM WorkDirectory WHERE WorkDirectoryId = ? OR AbsolutePath = ?;
