# Project Type Detection — Package Structure

## New Package

| Package    | Responsibility                              |
|------------|---------------------------------------------|
| `detector` | Walk repo trees and classify project types  |

---

## New Files

| File                              | Contents                                    |
|-----------------------------------|---------------------------------------------|
| `detector/detector.go`            | Walk repo, collect detected projects        |
| `detector/rules.go`               | Detection rules per project type            |
| `detector/parser.go`              | Parse manifest files (go.mod, package.json) |
| `detector/goparser.go`            | Parse Go metadata and find runnables        |
| `detector/csharpparser.go`        | Parse .csproj XML and collect key files     |
| `cmd/projectrepos.go`             | Query commands (go-repos, node-repos, etc.) |
| `cmd/projectreposoutput.go`       | Terminal and JSON output for queries        |
| `store/project.go`                | DetectedProject CRUD operations             |
| `store/projecttype.go`            | ProjectTypes seed and query                 |
| `store/gometadata.go`             | GoProjectMetadata + GoRunnableFiles CRUD    |
| `store/csharpmetadata.go`         | CSharpProjectMetadata + files CRUD          |
| `model/project.go`                | DetectedProject struct                      |
| `model/projecttype.go`            | ProjectType struct                          |
| `model/gometadata.go`             | GoProjectMetadata + GoRunnableFile structs  |
| `model/csharpmetadata.go`         | CSharpProjectMetadata + related structs     |
| `constants/constants_project.go`  | IDs, keys, table names, indicators, commands, messages, errors |
| `constants/constants_project_sql.go` | All SQL statements (create, upsert, query, cleanup, drop) |

---

## Modified Files

| File                              | Change                                      |
|-----------------------------------|---------------------------------------------|
| `cmd/scan.go`                     | Call detector after BuildRecords            |
| `cmd/scanoutput.go`               | Write project JSON files                    |
| `cmd/root.go`                     | Register query commands in dispatch         |
| `store/store.go`                  | Add project detection table migrations + drops |

---

## Model Structs

### ProjectType

```go
type ProjectType struct {
    ID          string `json:"id"`
    Key         string `json:"key"`
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

### DetectedProject

```go
type DetectedProject struct {
    ID               string `json:"id"`
    RepoID           string `json:"repoId"`
    RepoName         string `json:"repoName"`
    ProjectTypeID    string `json:"projectTypeId"`
    ProjectType      string `json:"projectType"`
    ProjectName      string `json:"projectName"`
    AbsolutePath     string `json:"absolutePath"`
    RepoPath         string `json:"repoPath"`
    RelativePath     string `json:"relativePath"`
    PrimaryIndicator string `json:"primaryIndicator"`
    DetectedAt       string `json:"detectedAt"`
}
```

### GoProjectMetadata

```go
type GoProjectMetadata struct {
    ID                string           `json:"id"`
    DetectedProjectID string           `json:"detectedProjectId"`
    GoModPath         string           `json:"goModPath"`
    GoSumPath         string           `json:"goSumPath"`
    ModuleName        string           `json:"moduleName"`
    GoVersion         string           `json:"goVersion"`
    Runnables         []GoRunnableFile `json:"runnables"`
}
```

### GoRunnableFile

```go
type GoRunnableFile struct {
    ID           string `json:"id"`
    GoMetadataID string `json:"goMetadataId"`
    RunnableName string `json:"runnableName"`
    FilePath     string `json:"filePath"`
    RelativePath string `json:"relativePath"`
}
```

### CSharpProjectMetadata

```go
type CSharpProjectMetadata struct {
    ID                string              `json:"id"`
    DetectedProjectID string              `json:"detectedProjectId"`
    SlnPath           string              `json:"slnPath"`
    SlnName           string              `json:"slnName"`
    GlobalJsonPath    string              `json:"globalJsonPath"`
    SdkVersion        string              `json:"sdkVersion"`
    ProjectFiles      []CSharpProjectFile `json:"projectFiles"`
    KeyFiles          []CSharpKeyFile     `json:"keyFiles"`
}
```

### CSharpProjectFile

```go
type CSharpProjectFile struct {
    ID               string `json:"id"`
    CSharpMetadataID string `json:"csharpMetadataId"`
    FilePath         string `json:"filePath"`
    RelativePath     string `json:"relativePath"`
    FileName         string `json:"fileName"`
    ProjectName      string `json:"projectName"`
    TargetFramework  string `json:"targetFramework"`
    OutputType       string `json:"outputType"`
    Sdk              string `json:"sdk"`
}
```

### CSharpKeyFile

```go
type CSharpKeyFile struct {
    ID               string `json:"id"`
    CSharpMetadataID string `json:"csharpMetadataId"`
    FileType         string `json:"fileType"`
    FilePath         string `json:"filePath"`
    RelativePath     string `json:"relativePath"`
}
```
