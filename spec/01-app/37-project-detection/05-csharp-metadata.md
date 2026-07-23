# Project Type Detection — C# Metadata

## Purpose

C# projects have rich metadata beyond basic detection. Solution files,
project files, target frameworks, and output types are all relevant
for downstream tooling. This metadata is stored in a separate table
because it is disjoint from the core `DetectedProjects` record.

---

## Collected Metadata

### From `.csproj` Files

| Field              | Source                    | Notes                            |
|--------------------|---------------------------|----------------------------------|
| CsprojPath         | `*.csproj` location       | Absolute path to .csproj file    |
| CsprojName         | `*.csproj` filename       | e.g., `MyApp.csproj`            |
| TargetFramework    | `<TargetFramework>` tag   | e.g., `net8.0`, `net6.0`        |
| OutputType         | `<OutputType>` tag        | `Exe`, `Library`, `WinExe`      |
| RootNamespace      | `<RootNamespace>` tag     | e.g., `MyApp.Core`              |
| AssemblyName       | `<AssemblyName>` tag      | e.g., `MyApp`                   |
| Sdk                | `<Project Sdk="...">` attr| e.g., `Microsoft.NET.Sdk.Web`    |

### From `.sln` Files

| Field              | Source                    | Notes                            |
|--------------------|---------------------------|----------------------------------|
| SlnPath            | `*.sln` location          | Absolute path to .sln file       |
| SlnName            | `*.sln` filename          | e.g., `MySolution.sln`          |

### Key File Types

C# projects commonly contain these key files that should be recorded:

| File Pattern        | Purpose                              |
|---------------------|--------------------------------------|
| `*.csproj`          | Project definition                   |
| `*.sln`             | Solution file (multi-project)        |
| `*.fsproj`          | F# project (related .NET)            |
| `global.json`       | SDK version pinning                  |
| `nuget.config`      | NuGet package source config          |
| `Directory.Build.props` | Shared MSBuild properties        |
| `Directory.Build.targets`| Shared MSBuild targets           |
| `*.props`           | Custom MSBuild property files        |
| `launchSettings.json`| Debug/run configuration             |
| `appsettings.json`  | Application configuration            |

---

## CSharpProjectMetadata Table

| Column              | Type    | Constraints                                         | Notes                      |
|---------------------|---------|-----------------------------------------------------|----------------------------|
| Id                  | TEXT    | PRIMARY KEY                                         | UUID                       |
| DetectedProjectId   | TEXT    | NOT NULL, FK → DetectedProjects(Id) ON DELETE CASCADE | Link to parent project   |
| SlnPath             | TEXT    | DEFAULT ''                                          | Absolute path to .sln      |
| SlnName             | TEXT    | DEFAULT ''                                          | Solution filename          |
| GlobalJsonPath      | TEXT    | DEFAULT ''                                          | Path to global.json        |
| SdkVersion          | TEXT    | DEFAULT ''                                          | From global.json `sdk.version` |

**Unique constraint:** `(DetectedProjectId)` — one metadata row per
detected C# project.

---

## CSharpProjectFiles Table

Stores each `.csproj` (or `.fsproj`) discovered within the C# project.

| Column              | Type    | Constraints                                              | Notes                   |
|---------------------|---------|----------------------------------------------------------|-------------------------|
| Id                  | TEXT    | PRIMARY KEY                                              | UUID                    |
| CSharpMetadataId    | TEXT    | NOT NULL, FK → CSharpProjectMetadata(Id) ON DELETE CASCADE | Link to C# metadata  |
| FilePath            | TEXT    | NOT NULL                                                 | Absolute path to file   |
| RelativePath        | TEXT    | NOT NULL                                                 | Path relative to project|
| FileName            | TEXT    | NOT NULL                                                 | e.g., `MyApp.csproj`   |
| ProjectName         | TEXT    | NOT NULL                                                 | Parsed from filename    |
| TargetFramework     | TEXT    | DEFAULT ''                                               | e.g., `net8.0`          |
| OutputType          | TEXT    | DEFAULT ''                                               | `Exe`, `Library`        |
| Sdk                 | TEXT    | DEFAULT ''                                               | e.g., `Microsoft.NET.Sdk.Web` |

**Unique constraint:** `(CSharpMetadataId, RelativePath)` — one entry
per project file per C# project.

---

## CSharpKeyFiles Table

Stores additional key files found in the C# project tree.

| Column              | Type    | Constraints                                              | Notes                   |
|---------------------|---------|----------------------------------------------------------|-------------------------|
| Id                  | TEXT    | PRIMARY KEY                                              | UUID                    |
| CSharpMetadataId    | TEXT    | NOT NULL, FK → CSharpProjectMetadata(Id) ON DELETE CASCADE | Link to C# metadata  |
| FileType            | TEXT    | NOT NULL                                                 | e.g., `nuget.config`   |
| FilePath            | TEXT    | NOT NULL                                                 | Absolute path           |
| RelativePath        | TEXT    | NOT NULL                                                 | Path relative to project|

**Unique constraint:** `(CSharpMetadataId, RelativePath)` — one entry
per key file per C# project.

---

## SQL Statements

### Create CSharpProjectMetadata

```sql
CREATE TABLE IF NOT EXISTS CSharpProjectMetadata (
    Id                TEXT PRIMARY KEY,
    DetectedProjectId TEXT NOT NULL UNIQUE
        REFERENCES DetectedProjects(Id) ON DELETE CASCADE,
    SlnPath           TEXT DEFAULT '',
    SlnName           TEXT DEFAULT '',
    GlobalJsonPath    TEXT DEFAULT '',
    SdkVersion        TEXT DEFAULT ''
)
```

### Create CSharpProjectFiles

```sql
CREATE TABLE IF NOT EXISTS CSharpProjectFiles (
    Id                TEXT PRIMARY KEY,
    CSharpMetadataId  TEXT NOT NULL
        REFERENCES CSharpProjectMetadata(Id) ON DELETE CASCADE,
    FilePath          TEXT NOT NULL,
    RelativePath      TEXT NOT NULL,
    FileName          TEXT NOT NULL,
    ProjectName       TEXT NOT NULL,
    TargetFramework   TEXT DEFAULT '',
    OutputType        TEXT DEFAULT '',
    Sdk               TEXT DEFAULT '',
    UNIQUE(CSharpMetadataId, RelativePath)
)
```

### Create CSharpKeyFiles

```sql
CREATE TABLE IF NOT EXISTS CSharpKeyFiles (
    Id                TEXT PRIMARY KEY,
    CSharpMetadataId  TEXT NOT NULL
        REFERENCES CSharpProjectMetadata(Id) ON DELETE CASCADE,
    FileType          TEXT NOT NULL,
    FilePath          TEXT NOT NULL,
    RelativePath      TEXT NOT NULL,
    UNIQUE(CSharpMetadataId, RelativePath)
)
```

---

## JSON Output Extension

When writing `csharp-projects.json`, each C# project record includes
additional fields:

```json
{
  "id": "uuid",
  "repoId": "uuid",
  "repoName": "my-dotnet-app",
  "projectType": "csharp",
  "projectName": "MyApp",
  "absolutePath": "/home/user/repos/my-dotnet-app",
  "repoPath": "/home/user/repos/my-dotnet-app",
  "relativePath": ".",
  "primaryIndicator": "MyApp.sln",
  "detectedAt": "2026-03-11T09:54:00Z",
  "csharpMetadata": {
    "slnPath": "/home/user/repos/my-dotnet-app/MyApp.sln",
    "slnName": "MyApp.sln",
    "globalJsonPath": "/home/user/repos/my-dotnet-app/global.json",
    "sdkVersion": "8.0.100",
    "projectFiles": [
      {
        "filePath": "/home/user/repos/my-dotnet-app/src/MyApp.Api/MyApp.Api.csproj",
        "relativePath": "src/MyApp.Api/MyApp.Api.csproj",
        "fileName": "MyApp.Api.csproj",
        "projectName": "MyApp.Api",
        "targetFramework": "net8.0",
        "outputType": "Exe",
        "sdk": "Microsoft.NET.Sdk.Web"
      },
      {
        "filePath": "/home/user/repos/my-dotnet-app/src/MyApp.Core/MyApp.Core.csproj",
        "relativePath": "src/MyApp.Core/MyApp.Core.csproj",
        "fileName": "MyApp.Core.csproj",
        "projectName": "MyApp.Core",
        "targetFramework": "net8.0",
        "outputType": "Library",
        "sdk": "Microsoft.NET.Sdk"
      }
    ],
    "keyFiles": [
      {
        "fileType": "nuget.config",
        "filePath": "/home/user/repos/my-dotnet-app/nuget.config",
        "relativePath": "nuget.config"
      },
      {
        "fileType": "Directory.Build.props",
        "filePath": "/home/user/repos/my-dotnet-app/Directory.Build.props",
        "relativePath": "Directory.Build.props"
      }
    ]
  }
}
```

---

## Detection Flow for C#

```
C# project root
  │
  ├─ Find *.sln           → record solution metadata
  ├─ Find global.json     → parse SDK version
  ├─ Walk tree:
  │   ├─ Find *.csproj    → parse XML for framework/output/SDK
  │   ├─ Find *.fsproj    → record as project file
  │   ├─ Find nuget.config     → record as key file
  │   ├─ Find Directory.Build.* → record as key file
  │   ├─ Find appsettings.json → record as key file
  │   └─ Find launchSettings.json → record as key file
  │
  └─ Assemble CSharpProjectMetadata + files + key files
```
