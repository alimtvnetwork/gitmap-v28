# Project Type Detection — Constants

## File

`gitmap-v28/constants/constants_project.go`

This file contains all constants for the project detection feature:
table names, project type seed data, SQL statements, JSON output
filenames, exclusion directories, detection indicators, command
metadata, and error messages.

---

## Project Type IDs

Stable UUIDs used to seed the `ProjectTypes` table. These are
referenced by `DetectedProjects.ProjectTypeId`.

```go
const (
    ProjectTypeGoID     = "b3f1a2c4-5d6e-4f7a-8b9c-0d1e2f3a4b5c"
    ProjectTypeNodeID   = "c4d2b3e5-6f7a-4e8b-9c0d-1e2f3a4b5c6d"
    ProjectTypeReactID  = "d5e3c4f6-7a8b-4f9c-0d1e-2f3a4b5c6d7e"
    ProjectTypeCppID    = "e6f4d5a7-8b9c-4a0d-1e2f-3a4b5c6d7e8f"
    ProjectTypeCSharpID = "f7a5e6b8-9c0d-4b1e-2f3a-4b5c6d7e8f9a"
)
```

## Project Type Keys

```go
const (
    ProjectKeyGo     = "go"
    ProjectKeyNode   = "node"
    ProjectKeyReact  = "react"
    ProjectKeyCpp    = "cpp"
    ProjectKeyCSharp = "csharp"
)
```

---

## Table Names

```go
const (
    TableProjectTypes         = "ProjectTypes"
    TableDetectedProjects     = "DetectedProjects"
    TableGoProjectMetadata    = "GoProjectMetadata"
    TableGoRunnableFiles      = "GoRunnableFiles"
    TableCSharpProjectMeta    = "CSharpProjectMetadata"
    TableCSharpProjectFiles   = "CSharpProjectFiles"
    TableCSharpKeyFiles       = "CSharpKeyFiles"
)
```

---

## JSON Output Filenames

```go
const (
    JSONFileGoProjects     = "go-projects.json"
    JSONFileNodeProjects   = "node-projects.json"
    JSONFileReactProjects  = "react-projects.json"
    JSONFileCppProjects    = "cpp-projects.json"
    JSONFileCSharpProjects = "csharp-projects.json"
)
```

---

## Detection Indicators

### Primary Indicators

```go
const (
    IndicatorGoMod       = "go.mod"
    IndicatorPackageJSON = "package.json"
    IndicatorCMakeLists  = "CMakeLists.txt"
    IndicatorMesonBuild  = "meson.build"
)
```

### File Extensions

```go
const (
    ExtCsproj  = ".csproj"
    ExtFsproj  = ".fsproj"
    ExtVcxproj = ".vcxproj"
    ExtSln     = ".sln"
)
```

### Go Structural Indicators

```go
const (
    GoCmdDir   = "cmd"
    GoMainFile = "main.go"
    GoSumFile  = "go.sum"
)
```

### React Detection Dependencies

```go
var ReactIndicatorDeps = []string{
    "react",
    "@types/react",
    "react-scripts",
    "next",
    "gatsby",
    "remix",
    "@remix-run/react",
}
```

### C# Key File Patterns

```go
var CSharpKeyFilePatterns = []string{
    "global.json",
    "nuget.config",
    "Directory.Build.props",
    "Directory.Build.targets",
    "launchSettings.json",
    "appsettings.json",
}
```

---

## Exclusion Directories

Directories skipped during project detection tree walk.

```go
var ProjectExcludeDirs = []string{
    "node_modules",
    "vendor",
    ".git",
    "dist",
    "build",
    "target",
    "bin",
    "obj",
    "out",
    "testdata",
    "packages",
    ".venv",
    ".cache",
}
```

`cmake-build-*` is matched by prefix, not exact string.

```go
const CMakeBuildPrefix = "cmake-build-"
```

---

## SQL: Create Tables

```go
const SQLCreateProjectTypes = `CREATE TABLE IF NOT EXISTS ProjectTypes (
    Id          TEXT PRIMARY KEY,
    Key         TEXT NOT NULL UNIQUE,
    Name        TEXT NOT NULL,
    Description TEXT DEFAULT ''
)`

const SQLCreateDetectedProjects = `CREATE TABLE IF NOT EXISTS DetectedProjects (
    Id               TEXT PRIMARY KEY,
    RepoId           TEXT NOT NULL REFERENCES Repos(Id) ON DELETE CASCADE,
    ProjectTypeId    TEXT NOT NULL REFERENCES ProjectTypes(Id),
    ProjectName      TEXT NOT NULL,
    AbsolutePath     TEXT NOT NULL,
    RepoPath         TEXT NOT NULL,
    RelativePath     TEXT NOT NULL,
    PrimaryIndicator TEXT NOT NULL,
    DetectedAt       TEXT DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(RepoId, ProjectTypeId, RelativePath)
)`

const SQLCreateGoProjectMetadata = `CREATE TABLE IF NOT EXISTS GoProjectMetadata (
    Id                TEXT PRIMARY KEY,
    DetectedProjectId TEXT NOT NULL UNIQUE
        REFERENCES DetectedProjects(Id) ON DELETE CASCADE,
    GoModPath         TEXT NOT NULL,
    GoSumPath         TEXT DEFAULT '',
    ModuleName        TEXT NOT NULL,
    GoVersion         TEXT DEFAULT ''
)`

const SQLCreateGoRunnableFiles = `CREATE TABLE IF NOT EXISTS GoRunnableFiles (
    Id           TEXT PRIMARY KEY,
    GoMetadataId TEXT NOT NULL
        REFERENCES GoProjectMetadata(Id) ON DELETE CASCADE,
    RunnableName TEXT NOT NULL,
    FilePath     TEXT NOT NULL,
    RelativePath TEXT NOT NULL,
    UNIQUE(GoMetadataId, RelativePath)
)`

const SQLCreateCSharpProjectMetadata = `CREATE TABLE IF NOT EXISTS CSharpProjectMetadata (
    Id                TEXT PRIMARY KEY,
    DetectedProjectId TEXT NOT NULL UNIQUE
        REFERENCES DetectedProjects(Id) ON DELETE CASCADE,
    SlnPath           TEXT DEFAULT '',
    SlnName           TEXT DEFAULT '',
    GlobalJsonPath    TEXT DEFAULT '',
    SdkVersion        TEXT DEFAULT ''
)`

const SQLCreateCSharpProjectFiles = `CREATE TABLE IF NOT EXISTS CSharpProjectFiles (
    Id               TEXT PRIMARY KEY,
    CSharpMetadataId TEXT NOT NULL
        REFERENCES CSharpProjectMetadata(Id) ON DELETE CASCADE,
    FilePath         TEXT NOT NULL,
    RelativePath     TEXT NOT NULL,
    FileName         TEXT NOT NULL,
    ProjectName      TEXT NOT NULL,
    TargetFramework  TEXT DEFAULT '',
    OutputType       TEXT DEFAULT '',
    Sdk              TEXT DEFAULT '',
    UNIQUE(CSharpMetadataId, RelativePath)
)`

const SQLCreateCSharpKeyFiles = `CREATE TABLE IF NOT EXISTS CSharpKeyFiles (
    Id               TEXT PRIMARY KEY,
    CSharpMetadataId TEXT NOT NULL
        REFERENCES CSharpProjectMetadata(Id) ON DELETE CASCADE,
    FileType         TEXT NOT NULL,
    FilePath         TEXT NOT NULL,
    RelativePath     TEXT NOT NULL,
    UNIQUE(CSharpMetadataId, RelativePath)
)`
```

---

## SQL: Seed ProjectTypes

```go
const SQLSeedProjectTypes = `INSERT OR IGNORE INTO ProjectTypes (Id, Key, Name, Description) VALUES
    ('b3f1a2c4-5d6e-4f7a-8b9c-0d1e2f3a4b5c', 'go',     'Go',      'Go modules and packages'),
    ('c4d2b3e5-6f7a-4e8b-9c0d-1e2f3a4b5c6d', 'node',   'Node.js', 'Node.js projects'),
    ('d5e3c4f6-7a8b-4f9c-0d1e-2f3a4b5c6d7e', 'react',  'React',   'React applications'),
    ('e6f4d5a7-8b9c-4a0d-1e2f-3a4b5c6d7e8f', 'cpp',    'C++',     'C and C++ projects'),
    ('f7a5e6b8-9c0d-4b1e-2f3a-4b5c6d7e8f9a', 'csharp', 'C#',      '.NET and C# projects')`
```

---

## SQL: Upsert Operations

```go
const SQLUpsertDetectedProject = `INSERT INTO DetectedProjects
    (Id, RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(RepoId, ProjectTypeId, RelativePath) DO UPDATE SET
        ProjectName=excluded.ProjectName,
        AbsolutePath=excluded.AbsolutePath,
        RepoPath=excluded.RepoPath,
        PrimaryIndicator=excluded.PrimaryIndicator,
        DetectedAt=CURRENT_TIMESTAMP`

const SQLUpsertGoMetadata = `INSERT INTO GoProjectMetadata
    (Id, DetectedProjectId, GoModPath, GoSumPath, ModuleName, GoVersion)
    VALUES (?, ?, ?, ?, ?, ?)
    ON CONFLICT(DetectedProjectId) DO UPDATE SET
        GoModPath=excluded.GoModPath,
        GoSumPath=excluded.GoSumPath,
        ModuleName=excluded.ModuleName,
        GoVersion=excluded.GoVersion`

const SQLUpsertGoRunnable = `INSERT INTO GoRunnableFiles
    (Id, GoMetadataId, RunnableName, FilePath, RelativePath)
    VALUES (?, ?, ?, ?, ?)
    ON CONFLICT(GoMetadataId, RelativePath) DO UPDATE SET
        RunnableName=excluded.RunnableName,
        FilePath=excluded.FilePath`

const SQLUpsertCSharpMetadata = `INSERT INTO CSharpProjectMetadata
    (Id, DetectedProjectId, SlnPath, SlnName, GlobalJsonPath, SdkVersion)
    VALUES (?, ?, ?, ?, ?, ?)
    ON CONFLICT(DetectedProjectId) DO UPDATE SET
        SlnPath=excluded.SlnPath,
        SlnName=excluded.SlnName,
        GlobalJsonPath=excluded.GlobalJsonPath,
        SdkVersion=excluded.SdkVersion`

const SQLUpsertCSharpProjectFile = `INSERT INTO CSharpProjectFiles
    (Id, CSharpMetadataId, FilePath, RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    ON CONFLICT(CSharpMetadataId, RelativePath) DO UPDATE SET
        FilePath=excluded.FilePath,
        FileName=excluded.FileName,
        ProjectName=excluded.ProjectName,
        TargetFramework=excluded.TargetFramework,
        OutputType=excluded.OutputType,
        Sdk=excluded.Sdk`

const SQLUpsertCSharpKeyFile = `INSERT INTO CSharpKeyFiles
    (Id, CSharpMetadataId, FileType, FilePath, RelativePath)
    VALUES (?, ?, ?, ?, ?)
    ON CONFLICT(CSharpMetadataId, RelativePath) DO UPDATE SET
        FileType=excluded.FileType,
        FilePath=excluded.FilePath`
```

---

## SQL: Query Operations

```go
const SQLSelectProjectsByTypeKey = `SELECT dp.Id, dp.RepoId, pt.Key AS ProjectType,
    dp.ProjectName, dp.AbsolutePath, dp.RepoPath, dp.RelativePath,
    dp.PrimaryIndicator, dp.DetectedAt, r.RepoName
    FROM DetectedProjects dp
    JOIN ProjectTypes pt ON dp.ProjectTypeId = pt.Id
    JOIN Repos r ON dp.RepoId = r.Id
    WHERE pt.Key = ?
    ORDER BY r.RepoName, dp.RelativePath`

const SQLCountProjectsByTypeKey = `SELECT COUNT(*)
    FROM DetectedProjects dp
    JOIN ProjectTypes pt ON dp.ProjectTypeId = pt.Id
    WHERE pt.Key = ?`

const SQLSelectGoMetadata = `SELECT Id, DetectedProjectId, GoModPath, GoSumPath,
    ModuleName, GoVersion
    FROM GoProjectMetadata WHERE DetectedProjectId = ?`

const SQLSelectGoRunnables = `SELECT Id, GoMetadataId, RunnableName, FilePath,
    RelativePath
    FROM GoRunnableFiles WHERE GoMetadataId = ?
    ORDER BY RunnableName`

const SQLSelectCSharpMetadata = `SELECT Id, DetectedProjectId, SlnPath, SlnName,
    GlobalJsonPath, SdkVersion
    FROM CSharpProjectMetadata WHERE DetectedProjectId = ?`

const SQLSelectCSharpProjectFiles = `SELECT Id, CSharpMetadataId, FilePath,
    RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk
    FROM CSharpProjectFiles WHERE CSharpMetadataId = ?
    ORDER BY RelativePath`

const SQLSelectCSharpKeyFiles = `SELECT Id, CSharpMetadataId, FileType, FilePath,
    RelativePath
    FROM CSharpKeyFiles WHERE CSharpMetadataId = ?
    ORDER BY RelativePath`
```

---

## SQL: Stale Cleanup

```go
const SQLDeleteStaleProjects = `DELETE FROM DetectedProjects
    WHERE RepoId = ? AND Id NOT IN (%s)`

const SQLDeleteStaleGoRunnables = `DELETE FROM GoRunnableFiles
    WHERE GoMetadataId = ? AND Id NOT IN (%s)`

const SQLDeleteStaleCSharpFiles = `DELETE FROM CSharpProjectFiles
    WHERE CSharpMetadataId = ? AND Id NOT IN (%s)`

const SQLDeleteStaleCSharpKeyFiles = `DELETE FROM CSharpKeyFiles
    WHERE CSharpMetadataId = ? AND Id NOT IN (%s)`
```

**Note:** `%s` is replaced at runtime with a comma-separated list of
`?` placeholders matching the number of upserted IDs.

---

## SQL: Drop Tables

Drop order respects foreign key dependencies (children first).

```go
const (
    SQLDropGoRunnableFiles      = "DROP TABLE IF EXISTS GoRunnableFiles"
    SQLDropGoProjectMetadata    = "DROP TABLE IF EXISTS GoProjectMetadata"
    SQLDropCSharpKeyFiles       = "DROP TABLE IF EXISTS CSharpKeyFiles"
    SQLDropCSharpProjectFiles   = "DROP TABLE IF EXISTS CSharpProjectFiles"
    SQLDropCSharpProjectMeta    = "DROP TABLE IF EXISTS CSharpProjectMetadata"
    SQLDropDetectedProjects     = "DROP TABLE IF EXISTS DetectedProjects"
    SQLDropProjectTypes         = "DROP TABLE IF EXISTS ProjectTypes"
)
```

---

## Command Constants

```go
const (
    CmdGoRepos     = "go-repos"
    CmdNodeRepos   = "node-repos"
    CmdReactRepos  = "react-repos"
    CmdCppRepos    = "cpp-repos"
    CmdCSharpRepos = "csharp-repos"
)
```

### Command Aliases

```go
const (
    AliasGoRepos     = "gr"
    AliasNodeRepos   = "nr"
    AliasReactRepos  = "rr"
    AliasCppRepos    = "cr"
    AliasCSharpRepos = "sr"
)
```

### Command Help Text

```go
const (
    HelpGoRepos     = "List repositories containing Go projects"
    HelpNodeRepos   = "List repositories containing Node.js projects"
    HelpReactRepos  = "List repositories containing React projects"
    HelpCppRepos    = "List repositories containing C++ projects"
    HelpCSharpRepos = "List repositories containing C# projects"
)
```

### Command Flag Names

```go
const (
    FlagProjectJSON  = "json"
    FlagProjectCount = "count"
)
```

---

## Messages

```go
const (
    MsgProjectDetectDone   = "Detected %d projects across %d repos\n"
    MsgProjectUpsertDone   = "Saved %d detected projects to database\n"
    MsgProjectJSONWritten  = "Wrote %s (%d records)\n"
    MsgProjectNoDB         = "No database found. Run 'gitmap-v28 scan' first.\n"
    MsgProjectNoneFound    = "No %s projects found.\n"
    MsgProjectCount        = "%d\n"
    MsgProjectCleanedStale = "Cleaned %d stale project records\n"
)
```

---

## Error Messages

```go
const (
    ErrProjectDetect       = "failed to detect projects in %s: %v\n"
    ErrProjectUpsert       = "failed to upsert detected project: %v"
    ErrProjectQuery        = "failed to query projects: %v"
    ErrProjectJSONWrite    = "failed to write %s: %v\n"
    ErrProjectParseMod     = "failed to parse go.mod in %s: %v\n"
    ErrProjectParsePkgJSON = "failed to parse package.json in %s: %v\n"
    ErrProjectParseCsproj  = "failed to parse .csproj in %s: %v\n"
    ErrProjectParseSln     = "failed to parse .sln in %s: %v\n"
    ErrProjectCleanup      = "failed to clean stale projects for repo %s: %v\n"
    ErrGoMetadataUpsert    = "failed to upsert Go metadata: %v"
    ErrGoRunnableUpsert    = "failed to upsert Go runnable: %v"
    ErrCSharpMetaUpsert    = "failed to upsert C# metadata: %v"
    ErrCSharpFileUpsert    = "failed to upsert C# project file: %v"
    ErrCSharpKeyUpsert     = "failed to upsert C# key file: %v"
)
```

---

## File Size Note

This file will exceed 200 lines. Split into two files if needed:

| File                          | Contents                                |
|-------------------------------|-----------------------------------------|
| `constants_project.go`       | IDs, keys, table names, indicators, exclusions, commands, messages, errors |
| `constants_project_sql.go`   | All SQL statements (create, upsert, query, cleanup, drop) |
