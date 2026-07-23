// Package store — migrate_v15phase4.go performs the Phase 1.4 v15 renames:
//
//	ZipGroups               → ZipGroup               (Id → ZipGroupId)
//	ZipGroupItems           → ZipGroupItem           (composite PK preserved)
//	ProjectTypes            → ProjectType            (Id → ProjectTypeId)
//	DetectedProjects        → DetectedProject        (Id → DetectedProjectId)
//	GoRunnableFiles         → GoRunnableFile         (Id → GoRunnableFileId)
//	GoProjectMetadata       → (singular kept)        (Id → GoProjectMetadataId)
//	CSharpProjectMetadata   → CsharpProjectMetadata  (Id → CsharpProjectMetadataId; abbreviation fix)
//	CSharpProjectFiles      → CsharpProjectFile      (Id → CsharpProjectFileId; abbreviation + singular)
//	CSharpKeyFiles          → CsharpKeyFile          (Id → CsharpKeyFileId; abbreviation + singular)
//	CommandHistory          → (singular kept)        (Id → CommandHistoryId)
//	RepoVersionHistory      → (singular kept)        (Id → RepoVersionHistoryId)
//	TaskType                → (singular kept)        (Id → TaskTypeId)
//	PendingTask             → (singular kept)        (Id → PendingTaskId)
//	CompletedTask           → (singular kept)        (Id → CompletedTaskId)
//
// Tables already in v15 form are still rebuilt to convert the `Id` PK into
// `{Table}Id`. Legacy CSharp-spelled tables are detected via their old name.
package store

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// migrateV15Phase4 runs all Phase 1.4 table rebuilds.
func (db *DB) migrateV15Phase4() error {
	specs := phase4Specs()

	for _, spec := range specs {
		if err := db.runV15Rebuild(spec); err != nil {
			return fmt.Errorf("phase 1.4 %s: %w", spec.OldTable, err)
		}
	}

	return nil
}

// phase4Specs returns the rebuild specs in FK-safe order (children first
// where the new schema has REFERENCES). Foreign keys are disabled during
// each rebuild so order matters less, but we still group logically.
func phase4Specs() []v15RebuildSpec {
	return append(append(append(append(
		zipGroupSpecs(),
		projectFamilySpecs()...),
		csharpFamilySpecs()...),
		taskFamilySpecs()...),
		historyFamilySpecs()...)
}

// zipGroupSpecs covers ZipGroups + ZipGroupItems.
func zipGroupSpecs() []v15RebuildSpec {
	return []v15RebuildSpec{
		{
			OldTable:      "ZipGroups",
			NewTable:      "ZipGroup",
			NewCreateSQL:  constants.SQLCreateZipGroup,
			OldColumnList: "Id, Name, ArchiveName, CreatedAt",
			NewColumnList: "ZipGroupId, Name, ArchiveName, CreatedAt",
			StartMsg:      "→ Migrating ZipGroups → ZipGroup (ZipGroupId PK)...",
			DoneMsg:       "✓ Migrated ZipGroups → ZipGroup.",
		},
		{
			OldTable: "ZipGroupItems",
			NewTable: "ZipGroupItem",
			// New CREATE references ZipGroup(ZipGroupId) — must run AFTER
			// the ZipGroups rebuild above. Composite PK preserved.
			NewCreateSQL:  constants.SQLCreateZipGroupItem,
			OldColumnList: "ZipGroupId, RepoPath, RelativePath, FullPath, IsFolder",
			NewColumnList: "ZipGroupId, RepoPath, RelativePath, FullPath, IsFolder",
			StartMsg:      "→ Migrating ZipGroupItems → ZipGroupItem...",
			DoneMsg:       "✓ Migrated ZipGroupItems → ZipGroupItem.",
		},
	}
}

// projectFamilySpecs covers ProjectTypes, DetectedProjects, GoProjectMetadata,
// GoRunnableFiles.
func projectFamilySpecs() []v15RebuildSpec {
	return []v15RebuildSpec{
		{
			OldTable:      "ProjectTypes",
			NewTable:      "ProjectType",
			NewCreateSQL:  constants.SQLCreateProjectType,
			OldColumnList: "Id, Key, Name, Description",
			NewColumnList: "ProjectTypeId, Key, Name, Description",
			StartMsg:      "→ Migrating ProjectTypes → ProjectType (ProjectTypeId PK)...",
			DoneMsg:       "✓ Migrated ProjectTypes → ProjectType.",
		},
		{
			OldTable:      "DetectedProjects",
			NewTable:      "DetectedProject",
			NewCreateSQL:  constants.SQLCreateDetectedProject,
			OldColumnList: "Id, RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator, DetectedAt",
			NewColumnList: "DetectedProjectId, RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator, DetectedAt",
			StartMsg:      "→ Migrating DetectedProjects → DetectedProject (DetectedProjectId PK)...",
			DoneMsg:       "✓ Migrated DetectedProjects → DetectedProject.",
		},
		// GoProjectMetadata is already singular; rebuild only to rename Id PK.
		// Legacy column was `DetectedProjectId` referencing DetectedProjects(Id);
		// since we just rebuilt that table preserving PK values, the integer
		// references stay valid.
		{
			OldTable:      "GoProjectMetadata",
			NewTable:      "GoProjectMetadata_v15",
			NewCreateSQL:  `CREATE TABLE IF NOT EXISTS GoProjectMetadata_v15 (` + goMetaBody() + `)`,
			OldColumnList: "Id, DetectedProjectId, GoModPath, GoSumPath, ModuleName, GoVersion",
			NewColumnList: "GoProjectMetadataId, DetectedProjectId, GoModPath, GoSumPath, ModuleName, GoVersion",
			StartMsg:      "→ Migrating GoProjectMetadata → (rebuild with GoProjectMetadataId PK)...",
			DoneMsg:       "✓ Rebuilt GoProjectMetadata with GoProjectMetadataId PK.",
		},
		{
			// Final rename of staging table back to canonical name.
			OldTable:      "GoProjectMetadata_v15",
			NewTable:      "GoProjectMetadata",
			NewCreateSQL:  constants.SQLCreateGoProjectMetadata,
			OldColumnList: "GoProjectMetadataId, DetectedProjectId, GoModPath, GoSumPath, ModuleName, GoVersion",
			NewColumnList: "GoProjectMetadataId, DetectedProjectId, GoModPath, GoSumPath, ModuleName, GoVersion",
		},
		{
			OldTable:      "GoRunnableFiles",
			NewTable:      "GoRunnableFile",
			NewCreateSQL:  constants.SQLCreateGoRunnableFile,
			OldColumnList: "Id, GoMetadataId, RunnableName, FilePath, RelativePath",
			NewColumnList: "GoRunnableFileId, GoMetadataId, RunnableName, FilePath, RelativePath",
			StartMsg:      "→ Migrating GoRunnableFiles → GoRunnableFile (GoRunnableFileId PK)...",
			DoneMsg:       "✓ Migrated GoRunnableFiles → GoRunnableFile.",
		},
	}
}

// goMetaBody is the inline schema for the GoProjectMetadata staging table.
// We can't use the canonical SQLCreateGoProjectMetadata because its REFERENCES
// clause names DetectedProject(DetectedProjectId), but during the rebuild the
// staging table must reference whatever exists. Since FKs are disabled inside
// runV15Rebuild, we still emit the same body for consistency.
func goMetaBody() string {
	return `
	GoProjectMetadataId INTEGER PRIMARY KEY AUTOINCREMENT,
	DetectedProjectId   INTEGER NOT NULL UNIQUE,
	GoModPath           TEXT NOT NULL,
	GoSumPath           TEXT DEFAULT '',
	ModuleName          TEXT NOT NULL,
	GoVersion           TEXT DEFAULT ''`
}

// csharpFamilySpecs covers the three C# tables, including the
// CSharp→Csharp abbreviation fix (legacy table names use the old spelling).
func csharpFamilySpecs() []v15RebuildSpec {
	return []v15RebuildSpec{
		{
			OldTable:      "CSharpProjectMetadata",
			NewTable:      "CsharpProjectMetadata",
			NewCreateSQL:  constants.SQLCreateCsharpProjectMeta,
			OldColumnList: "Id, DetectedProjectId, SlnPath, SlnName, GlobalJsonPath, SdkVersion",
			NewColumnList: "CsharpProjectMetadataId, DetectedProjectId, SlnPath, SlnName, GlobalJsonPath, SdkVersion",
			StartMsg:      "→ Migrating CSharpProjectMetadata → CsharpProjectMetadata (abbreviation + CsharpProjectMetadataId PK)...",
			DoneMsg:       "✓ Migrated CSharpProjectMetadata → CsharpProjectMetadata.",
		},
		{
			OldTable:      "CSharpProjectFiles",
			NewTable:      "CsharpProjectFile",
			NewCreateSQL:  constants.SQLCreateCsharpProjectFile,
			OldColumnList: "Id, CSharpMetadataId, FilePath, RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk",
			NewColumnList: "CsharpProjectFileId, CsharpMetadataId, FilePath, RelativePath, FileName, ProjectName, TargetFramework, OutputType, Sdk",
			StartMsg:      "→ Migrating CSharpProjectFiles → CsharpProjectFile (abbreviation + CsharpProjectFileId PK)...",
			DoneMsg:       "✓ Migrated CSharpProjectFiles → CsharpProjectFile.",
		},
		{
			OldTable:      "CSharpKeyFiles",
			NewTable:      "CsharpKeyFile",
			NewCreateSQL:  constants.SQLCreateCsharpKeyFile,
			OldColumnList: "Id, CSharpMetadataId, FileType, FilePath, RelativePath",
			NewColumnList: "CsharpKeyFileId, CsharpMetadataId, FileType, FilePath, RelativePath",
			StartMsg:      "→ Migrating CSharpKeyFiles → CsharpKeyFile (abbreviation + CsharpKeyFileId PK)...",
			DoneMsg:       "✓ Migrated CSharpKeyFiles → CsharpKeyFile.",
		},
	}
}

// taskFamilySpecs covers TaskType + PendingTask + CompletedTask.
// All three are already singular; rebuild only to rename Id → {Table}Id.
func taskFamilySpecs() []v15RebuildSpec {
	return []v15RebuildSpec{
		{
			OldTable:      "TaskType",
			NewTable:      "TaskType_v15",
			NewCreateSQL:  `CREATE TABLE IF NOT EXISTS TaskType_v15 (TaskTypeId INTEGER PRIMARY KEY AUTOINCREMENT, Name TEXT NOT NULL UNIQUE)`,
			OldColumnList: "Id, Name",
			NewColumnList: "TaskTypeId, Name",
			StartMsg:      "→ Migrating TaskType → (rebuild with TaskTypeId PK)...",
			DoneMsg:       "✓ Rebuilt TaskType with TaskTypeId PK.",
		},
		{
			OldTable:      "TaskType_v15",
			NewTable:      "TaskType",
			NewCreateSQL:  constants.SQLCreateTaskType,
			OldColumnList: "TaskTypeId, Name",
			NewColumnList: "TaskTypeId, Name",
		},
		{
			OldTable:      "PendingTask",
			NewTable:      "PendingTask_v15",
			NewCreateSQL:  `CREATE TABLE IF NOT EXISTS PendingTask_v15 (` + pendingTaskBody() + `)`,
			OldColumnList: "Id, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, FailureReason, CreatedAt, UpdatedAt",
			NewColumnList: "PendingTaskId, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, FailureReason, CreatedAt, UpdatedAt",
			StartMsg:      "→ Migrating PendingTask → (rebuild with PendingTaskId PK)...",
			DoneMsg:       "✓ Rebuilt PendingTask with PendingTaskId PK.",
		},
		{
			OldTable:      "PendingTask_v15",
			NewTable:      "PendingTask",
			NewCreateSQL:  constants.SQLCreatePendingTask,
			OldColumnList: "PendingTaskId, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, FailureReason, CreatedAt, UpdatedAt",
			NewColumnList: "PendingTaskId, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, FailureReason, CreatedAt, UpdatedAt",
		},
		{
			OldTable:      "CompletedTask",
			NewTable:      "CompletedTask_v15",
			NewCreateSQL:  `CREATE TABLE IF NOT EXISTS CompletedTask_v15 (` + completedTaskBody() + `)`,
			OldColumnList: "Id, OriginalTaskId, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, CompletedAt, CreatedAt",
			NewColumnList: "CompletedTaskId, OriginalTaskId, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, CompletedAt, CreatedAt",
			StartMsg:      "→ Migrating CompletedTask → (rebuild with CompletedTaskId PK)...",
			DoneMsg:       "✓ Rebuilt CompletedTask with CompletedTaskId PK.",
		},
		{
			OldTable:      "CompletedTask_v15",
			NewTable:      "CompletedTask",
			NewCreateSQL:  constants.SQLCreateCompletedTask,
			OldColumnList: "CompletedTaskId, OriginalTaskId, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, CompletedAt, CreatedAt",
			NewColumnList: "CompletedTaskId, OriginalTaskId, TaskTypeId, TargetPath, WorkingDirectory, SourceCommand, CommandArgs, CompletedAt, CreatedAt",
		},
	}
}

func pendingTaskBody() string {
	return `
	PendingTaskId    INTEGER PRIMARY KEY AUTOINCREMENT,
	TaskTypeId       INTEGER NOT NULL,
	TargetPath       TEXT    NOT NULL,
	WorkingDirectory TEXT    DEFAULT '',
	SourceCommand    TEXT    NOT NULL,
	CommandArgs      TEXT    DEFAULT '',
	FailureReason    TEXT    DEFAULT '',
	CreatedAt        TEXT    DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt        TEXT    DEFAULT CURRENT_TIMESTAMP`
}

func completedTaskBody() string {
	return `
	CompletedTaskId  INTEGER PRIMARY KEY AUTOINCREMENT,
	OriginalTaskId   INTEGER NOT NULL,
	TaskTypeId       INTEGER NOT NULL,
	TargetPath       TEXT    NOT NULL,
	WorkingDirectory TEXT    DEFAULT '',
	SourceCommand    TEXT    NOT NULL,
	CommandArgs      TEXT    DEFAULT '',
	CompletedAt      TEXT    DEFAULT CURRENT_TIMESTAMP,
	CreatedAt        TEXT    NOT NULL`
}

// historyFamilySpecs covers CommandHistory + RepoVersionHistory. Both are
// already singular; rebuild only to rename Id → {Table}Id.
func historyFamilySpecs() []v15RebuildSpec {
	return []v15RebuildSpec{
		{
			OldTable:      "CommandHistory",
			NewTable:      "CommandHistory_v15",
			NewCreateSQL:  `CREATE TABLE IF NOT EXISTS CommandHistory_v15 (` + commandHistoryBody() + `)`,
			OldColumnList: "Id, Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount, CreatedAt",
			NewColumnList: "CommandHistoryId, Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount, CreatedAt",
			StartMsg:      "→ Migrating CommandHistory → (rebuild with CommandHistoryId PK)...",
			DoneMsg:       "✓ Rebuilt CommandHistory with CommandHistoryId PK.",
		},
		{
			OldTable:      "CommandHistory_v15",
			NewTable:      "CommandHistory",
			NewCreateSQL:  constants.SQLCreateCommandHistory,
			OldColumnList: "CommandHistoryId, Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount, CreatedAt",
			NewColumnList: "CommandHistoryId, Command, Alias, Args, Flags, StartedAt, FinishedAt, DurationMs, ExitCode, Summary, RepoCount, CreatedAt",
		},
		{
			OldTable:      "RepoVersionHistory",
			NewTable:      "RepoVersionHistory_v15",
			NewCreateSQL:  `CREATE TABLE IF NOT EXISTS RepoVersionHistory_v15 (` + repoVersionHistoryBody() + `)`,
			OldColumnList: "Id, RepoId, FromVersionTag, FromVersionNum, ToVersionTag, ToVersionNum, FlattenedPath, CreatedAt",
			NewColumnList: "RepoVersionHistoryId, RepoId, FromVersionTag, FromVersionNum, ToVersionTag, ToVersionNum, FlattenedPath, CreatedAt",
			StartMsg:      "→ Migrating RepoVersionHistory → (rebuild with RepoVersionHistoryId PK)...",
			DoneMsg:       "✓ Rebuilt RepoVersionHistory with RepoVersionHistoryId PK.",
		},
		{
			OldTable:      "RepoVersionHistory_v15",
			NewTable:      "RepoVersionHistory",
			NewCreateSQL:  constants.SQLCreateRepoVersionHistory,
			OldColumnList: "RepoVersionHistoryId, RepoId, FromVersionTag, FromVersionNum, ToVersionTag, ToVersionNum, FlattenedPath, CreatedAt",
			NewColumnList: "RepoVersionHistoryId, RepoId, FromVersionTag, FromVersionNum, ToVersionTag, ToVersionNum, FlattenedPath, CreatedAt",
		},
	}
}

func commandHistoryBody() string {
	return `
	CommandHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
	Command          TEXT NOT NULL,
	Alias            TEXT DEFAULT '',
	Args             TEXT DEFAULT '',
	Flags            TEXT DEFAULT '',
	StartedAt        TEXT NOT NULL,
	FinishedAt       TEXT DEFAULT '',
	DurationMs       INTEGER DEFAULT 0,
	ExitCode         INTEGER DEFAULT 0,
	Summary          TEXT DEFAULT '',
	RepoCount        INTEGER DEFAULT 0,
	CreatedAt        TEXT DEFAULT CURRENT_TIMESTAMP`
}

func repoVersionHistoryBody() string {
	return `
	RepoVersionHistoryId INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoId               INTEGER NOT NULL,
	FromVersionTag       TEXT NOT NULL,
	FromVersionNum       INTEGER NOT NULL,
	ToVersionTag         TEXT NOT NULL,
	ToVersionNum         INTEGER NOT NULL,
	FlattenedPath        TEXT DEFAULT '',
	CreatedAt            TEXT DEFAULT CURRENT_TIMESTAMP`
}
