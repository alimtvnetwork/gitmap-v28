package constants

// SQL: create ProjectType table (v15: singular + ProjectTypeId PK).
const SQLCreateProjectType = `CREATE TABLE IF NOT EXISTS ProjectType (
	ProjectTypeId INTEGER PRIMARY KEY AUTOINCREMENT,
	Key           TEXT NOT NULL UNIQUE,
	Name          TEXT NOT NULL,
	Description   TEXT DEFAULT ''
)`

// SQL: create DetectedProject table (v15: singular + DetectedProjectId PK).
// FK references v15 Repo(RepoId) and ProjectType(ProjectTypeId).
const SQLCreateDetectedProject = `CREATE TABLE IF NOT EXISTS DetectedProject (
	DetectedProjectId INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoId            INTEGER NOT NULL REFERENCES Repo(RepoId) ON DELETE CASCADE,
	ProjectTypeId     INTEGER NOT NULL REFERENCES ProjectType(ProjectTypeId),
	ProjectName       TEXT NOT NULL,
	AbsolutePath      TEXT NOT NULL,
	RepoPath          TEXT NOT NULL,
	RelativePath      TEXT NOT NULL,
	PrimaryIndicator  TEXT NOT NULL,
	DetectedAt        TEXT DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(RepoId, ProjectTypeId, RelativePath)
)`

// SQL: create GoProjectMetadata table (v15: GoProjectMetadataId PK).
const SQLCreateGoProjectMetadata = `CREATE TABLE IF NOT EXISTS GoProjectMetadata (
	GoProjectMetadataId INTEGER PRIMARY KEY AUTOINCREMENT,
	DetectedProjectId   INTEGER NOT NULL UNIQUE
		REFERENCES DetectedProject(DetectedProjectId) ON DELETE CASCADE,
	GoModPath           TEXT NOT NULL,
	GoSumPath           TEXT DEFAULT '',
	ModuleName          TEXT NOT NULL,
	GoVersion           TEXT DEFAULT ''
)`

// SQL: create GoRunnableFile table (v15: singular + GoRunnableFileId PK).
const SQLCreateGoRunnableFile = `CREATE TABLE IF NOT EXISTS GoRunnableFile (
	GoRunnableFileId INTEGER PRIMARY KEY AUTOINCREMENT,
	GoMetadataId     INTEGER NOT NULL
		REFERENCES GoProjectMetadata(GoProjectMetadataId) ON DELETE CASCADE,
	RunnableName     TEXT NOT NULL,
	FilePath         TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	UNIQUE(GoMetadataId, RelativePath)
)`

// SQL: create CsharpProjectMetadata table (v15: CsharpProjectMetadataId PK
// + Csharp abbreviation per strict v15 PascalCase rule).
const SQLCreateCsharpProjectMeta = `CREATE TABLE IF NOT EXISTS CsharpProjectMetadata (
	CsharpProjectMetadataId INTEGER PRIMARY KEY AUTOINCREMENT,
	DetectedProjectId       INTEGER NOT NULL UNIQUE
		REFERENCES DetectedProject(DetectedProjectId) ON DELETE CASCADE,
	SlnPath                 TEXT DEFAULT '',
	SlnName                 TEXT DEFAULT '',
	GlobalJsonPath          TEXT DEFAULT '',
	SdkVersion              TEXT DEFAULT ''
)`

// SQL: create CsharpProjectFile table (v15: singular + CsharpProjectFileId PK).
const SQLCreateCsharpProjectFile = `CREATE TABLE IF NOT EXISTS CsharpProjectFile (
	CsharpProjectFileId INTEGER PRIMARY KEY AUTOINCREMENT,
	CsharpMetadataId    INTEGER NOT NULL
		REFERENCES CsharpProjectMetadata(CsharpProjectMetadataId) ON DELETE CASCADE,
	FilePath            TEXT NOT NULL,
	RelativePath        TEXT NOT NULL,
	FileName            TEXT NOT NULL,
	ProjectName         TEXT NOT NULL,
	TargetFramework    TEXT DEFAULT '',
	OutputType          TEXT DEFAULT '',
	Sdk                 TEXT DEFAULT '',
	UNIQUE(CsharpMetadataId, RelativePath)
)`

// SQL: create CsharpKeyFile table (v15: singular + CsharpKeyFileId PK).
const SQLCreateCsharpKeyFile = `CREATE TABLE IF NOT EXISTS CsharpKeyFile (
	CsharpKeyFileId  INTEGER PRIMARY KEY AUTOINCREMENT,
	CsharpMetadataId INTEGER NOT NULL
		REFERENCES CsharpProjectMetadata(CsharpProjectMetadataId) ON DELETE CASCADE,
	FileType         TEXT NOT NULL,
	FilePath         TEXT NOT NULL,
	RelativePath     TEXT NOT NULL,
	UNIQUE(CsharpMetadataId, RelativePath)
)`

// SQL: seed project types.
const SQLSeedProjectTypes = `INSERT OR IGNORE INTO ProjectType (Key, Name, Description) VALUES
	('go',     'Go',      'Go modules and packages'),
	('node',   'Node.js', 'Node.js projects'),
	('react',  'React',   'React applications'),
	('cpp',    'C++',     'C and C++ projects'),
	('csharp', 'C#',      '.NET and C# projects')`

// SQL: upsert detected project.
const SQLUpsertDetectedProject = `INSERT INTO DetectedProject
	(RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(RepoId, ProjectTypeId, RelativePath) DO UPDATE SET
		ProjectName=excluded.ProjectName,
		AbsolutePath=excluded.AbsolutePath,
		RepoPath=excluded.RepoPath,
		PrimaryIndicator=excluded.PrimaryIndicator,
		DetectedAt=CURRENT_TIMESTAMP`

// SQL: upsert Go metadata.
const SQLUpsertGoMetadata = `INSERT INTO GoProjectMetadata
	(DetectedProjectId, GoModPath, GoSumPath, ModuleName, GoVersion)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(DetectedProjectId) DO UPDATE SET
		GoModPath=excluded.GoModPath,
		GoSumPath=excluded.GoSumPath,
		ModuleName=excluded.ModuleName,
		GoVersion=excluded.GoVersion`

// SQL: upsert Go runnable file.
const SQLUpsertGoRunnable = `INSERT INTO GoRunnableFile
	(GoMetadataId, RunnableName, FilePath, RelativePath)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(GoMetadataId, RelativePath) DO UPDATE SET
		RunnableName=excluded.RunnableName,
		FilePath=excluded.FilePath`

// SQL: upsert C# metadata.
const SQLUpsertCsharpMetadata = `INSERT INTO CsharpProjectMetadata
	(DetectedProjectId, SlnPath, SlnName, GlobalJsonPath, SdkVersion)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(DetectedProjectId) DO UPDATE SET
		SlnPath=excluded.SlnPath,
		SlnName=excluded.SlnName,
		GlobalJsonPath=excluded.GlobalJsonPath,
		SdkVersion=excluded.SdkVersion`

// SQL: upsert C# project file.
const SQLUpsertCsharpProjectFile = `INSERT INTO CsharpProjectFile
	(CsharpMetadataId, FilePath, RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(CsharpMetadataId, RelativePath) DO UPDATE SET
		FilePath=excluded.FilePath,
		FileName=excluded.FileName,
		ProjectName=excluded.ProjectName,
		TargetFramework=excluded.TargetFramework,
		OutputType=excluded.OutputType,
		Sdk=excluded.Sdk`

// SQL: upsert C# key file.
const SQLUpsertCsharpKeyFile = `INSERT INTO CsharpKeyFile
	(CsharpMetadataId, FileType, FilePath, RelativePath)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(CsharpMetadataId, RelativePath) DO UPDATE SET
		FileType=excluded.FileType,
		FilePath=excluded.FilePath`

// SQL: query detected project ID by identity tuple.
const SQLSelectDetectedProjectID = `SELECT DetectedProjectId
	FROM DetectedProject
	WHERE RepoId = ? AND ProjectTypeId = ? AND RelativePath = ?`

// SQL: query projects by type key (v15: JOIN Repo on RepoId, ProjectType on ProjectTypeId).
const SQLSelectProjectsByTypeKey = `SELECT dp.DetectedProjectId, dp.RepoId, pt.Key, dp.ProjectName,
	dp.AbsolutePath, dp.RepoPath, dp.RelativePath,
	dp.PrimaryIndicator, dp.DetectedAt, r.RepoName
	FROM DetectedProject dp
	JOIN ProjectType pt ON dp.ProjectTypeId = pt.ProjectTypeId
	JOIN Repo r ON dp.RepoId = r.RepoId
	WHERE pt.Key = ?
	ORDER BY r.RepoName, dp.RelativePath`

// SQL: count projects by type key.
const SQLCountProjectsByTypeKey = `SELECT COUNT(*)
	FROM DetectedProject dp
	JOIN ProjectType pt ON dp.ProjectTypeId = pt.ProjectTypeId
	WHERE pt.Key = ?`

// SQL: query Go metadata.
const SQLSelectGoMetadata = `SELECT GoProjectMetadataId, DetectedProjectId, GoModPath, GoSumPath,
	ModuleName, GoVersion
	FROM GoProjectMetadata WHERE DetectedProjectId = ?`

// SQL: query Go runnables.
const SQLSelectGoRunnables = `SELECT GoRunnableFileId, GoMetadataId, RunnableName, FilePath,
	RelativePath
	FROM GoRunnableFile WHERE GoMetadataId = ?
	ORDER BY RunnableName`

// SQL: query C# metadata.
const SQLSelectCsharpMetadata = `SELECT CsharpProjectMetadataId, DetectedProjectId, SlnPath, SlnName,
	GlobalJsonPath, SdkVersion
	FROM CsharpProjectMetadata WHERE DetectedProjectId = ?`

// SQL: query C# project files.
const SQLSelectCsharpProjectFiles = `SELECT CsharpProjectFileId, CsharpMetadataId, FilePath,
	RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk
	FROM CsharpProjectFile WHERE CsharpMetadataId = ?
	ORDER BY RelativePath`

// SQL: query C# key files.
const SQLSelectCsharpKeyFiles = `SELECT CsharpKeyFileId, CsharpMetadataId, FileType, FilePath,
	RelativePath
	FROM CsharpKeyFile WHERE CsharpMetadataId = ?
	ORDER BY RelativePath`

// SQL: stale cleanup (v15 singular tables, {Table}Id PKs in WHERE/IN clauses).
const (
	SQLDeleteStaleProjects       = "DELETE FROM DetectedProject WHERE RepoId = ? AND DetectedProjectId NOT IN (%s)"
	SQLDeleteStaleGoRunnables    = "DELETE FROM GoRunnableFile WHERE GoMetadataId = ? AND GoRunnableFileId NOT IN (%s)"
	SQLDeleteStaleCsharpFiles    = "DELETE FROM CsharpProjectFile WHERE CsharpMetadataId = ? AND CsharpProjectFileId NOT IN (%s)"
	SQLDeleteStaleCsharpKeyFiles = "DELETE FROM CsharpKeyFile WHERE CsharpMetadataId = ? AND CsharpKeyFileId NOT IN (%s)"
)

// SQL: drop project detection tables (v15 names + legacy retained for Reset).
const (
	SQLDropGoRunnableFile          = "DROP TABLE IF EXISTS GoRunnableFile"
	SQLDropGoRunnableFiles         = "DROP TABLE IF EXISTS GoRunnableFiles" // legacy
	SQLDropGoProjectMetadata       = "DROP TABLE IF EXISTS GoProjectMetadata"
	SQLDropCsharpKeyFile           = "DROP TABLE IF EXISTS CsharpKeyFile"
	SQLDropCsharpKeyFiles          = "DROP TABLE IF EXISTS CSharpKeyFiles" // legacy (pre-Csharp + plural)
	SQLDropCsharpProjectFile       = "DROP TABLE IF EXISTS CsharpProjectFile"
	SQLDropCsharpProjectFiles      = "DROP TABLE IF EXISTS CSharpProjectFiles" // legacy
	SQLDropCsharpProjectMeta       = "DROP TABLE IF EXISTS CsharpProjectMetadata"
	SQLDropCsharpProjectMetaLegacy = "DROP TABLE IF EXISTS CSharpProjectMetadata" // legacy
	SQLDropDetectedProject         = "DROP TABLE IF EXISTS DetectedProject"
	SQLDropDetectedProjects        = "DROP TABLE IF EXISTS DetectedProjects" // legacy
	SQLDropProjectType             = "DROP TABLE IF EXISTS ProjectType"
	SQLDropProjectTypes            = "DROP TABLE IF EXISTS ProjectTypes" // legacy
)
