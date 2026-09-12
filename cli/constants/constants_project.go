package constants

// Project type IDs matching seed order in SQLSeedProjectTypes.
const (
	ProjectTypeGoID     int64 = 1
	ProjectTypeNodeID   int64 = 2
	ProjectTypeReactID  int64 = 3
	ProjectTypeCppID    int64 = 4
	ProjectTypeCsharpID int64 = 5
)

// Project type keys.
const (
	ProjectKeyGo     = "go"
	ProjectKeyNode   = "node"
	ProjectKeyReact  = "react"
	ProjectKeyCpp    = "cpp"
	ProjectKeyCsharp = "csharp"
)

// Project detection table names (v15: PascalCase singular + {Table}Id PK).
const (
	TableProjectType       = "ProjectType"
	TableDetectedProject   = "DetectedProject"
	TableGoProjectMetadata = "GoProjectMetadata"
	TableGoRunnableFile    = "GoRunnableFile"
	TableCsharpProjectMeta = "CsharpProjectMetadata"
	TableCsharpProjectFile = "CsharpProjectFile"
	TableCsharpKeyFile     = "CsharpKeyFile"
)

// Legacy project detection table names retained ONLY for migration detection
// (do not use in new SQL). Includes both pre-v15 plurals and the
// pre-Csharp-rename "CSharp*" spellings.
const (
	LegacyTableProjectTypes       = "ProjectTypes"
	LegacyTableDetectedProjects   = "DetectedProjects"
	LegacyTableGoRunnableFiles    = "GoRunnableFiles"
	LegacyTableCsharpProjectMeta  = "CSharpProjectMetadata" // pre-Csharp spelling
	LegacyTableCsharpProjectFiles = "CSharpProjectFiles"    // pre-Csharp spelling + plural
	LegacyTableCsharpKeyFiles     = "CSharpKeyFiles"        // pre-Csharp spelling + plural
)

// Project JSON output filenames.
const (
	JSONFileGoProjects     = "go-projects.json"
	JSONFileNodeProjects   = "node-projects.json"
	JSONFileReactProjects  = "react-projects.json"
	JSONFileCppProjects    = "cpp-projects.json"
	JSONFileCsharpProjects = "csharp-projects.json"
)

// Detection indicator files.
const (
	IndicatorGoMod       = "go.mod"
	IndicatorPackageJSON = "package.json"
	IndicatorCMakeLists  = "CMakeLists.txt"
	IndicatorMesonBuild  = "meson.build"
)

// Detection file extensions.
const (
	ExtCsproj  = ".csproj"
	ExtFsproj  = ".fsproj"
	ExtVcxproj = ".vcxproj"
	ExtSln     = ".sln"
)

// Go structural indicators.
const (
	GoCmdDir      = "cmd"
	GoMainFile    = "main.go"
	GoSumFile     = "go.sum"
	CMakeBuildPfx = "cmake-build-"
)

// gitmap:cmd top-level
// Project query commands.
const (
	CmdGoRepos         = "go-repos"
	CmdGoReposAlias    = "gr"
	CmdNodeRepos       = "node-repos"
	CmdNodeReposAlias  = "nr"
	CmdReactRepos      = "react-repos"
	CmdReactReposAlias = "rr"
	CmdCppRepos        = "cpp-repos"
	CmdCppReposAlias   = "cr"
	CmdCsharpRepos     = "csharp-repos"
	CmdCsharpAlias     = "csr"
)

// Project query flags.
const (
	FlagProjectJSON  = "json"
	FlagProjectCount = "count"
)

// Project query help text.
const (
	HelpGoRepos     = "  go-repos (gr)       List repositories containing Go projects"
	HelpNodeRepos   = "  node-repos (nr)     List repositories containing Node.js projects"
	HelpReactRepos  = "  react-repos (rr)    List repositories containing React projects"
	HelpCppRepos    = "  cpp-repos (cr)      List repositories containing C++ projects"
	HelpCsharpRepos = "  csharp-repos (csr)  List repositories containing C# projects"
)

// Project detection messages.
const (
	MsgProjectDetectDone   = "  🧭 Detected %d project(s) across %d repo(s)\n"
	MsgProjectUpsertDone   = "  ✅ Saved %d detected project(s) to database\n"
	MsgProjectJSONWritten  = "  📄 %-22s %d record(s)\n"
	MsgProjectNoDB         = "No database found. Run 'gitmap scan' first.\n"
	MsgProjectNoneFound    = "No %s projects found.\n"
	MsgProjectCount        = "%d\n"
	MsgProjectCleanedStale = "Cleaned %d stale project records\n"
	MsgProjectListCount    = "\n%d projects found.\n"
)

// Project detection error messages.
const (
	ErrProjectDetect       = "failed to detect projects in %s: %v\n"
	ErrProjectUpsert       = "failed to upsert detected project: %v"
	ErrProjectQuery        = "failed to query projects: %v"
	ErrProjectJSONWrite    = "failed to write %s: %v (operation: write)\n"
	ErrProjectParseMod     = "failed to parse go.mod in %s: %v\n"
	ErrProjectParsePkgJSON = "failed to parse package.json in %s: %v\n"
	ErrProjectParseCsproj  = "failed to parse .csproj in %s: %v\n"
	ErrProjectCleanup      = "failed to clean stale projects for repo %d: %v\n"
	ErrGoMetadataUpsert    = "failed to upsert Go metadata: %v"
	ErrGoRunnableUpsert    = "failed to upsert Go runnable: %v"
	ErrCsharpMetaUpsert    = "failed to upsert C# metadata: %v"
	ErrCsharpFileUpsert    = "failed to upsert C# project file: %v"
	ErrCsharpKeyUpsert     = "failed to upsert C# key file: %v"
)

// Legacy data recovery messages.
const (
	MsgLegacyProjectData = "Database contains legacy project data from a previous version.\n" +
		"To fix, run one of:\n\n" +
		"  gitmap rescan          Re-scan repos and rebuild project data\n" +
		"  gitmap db-reset --confirm   Reset the entire database\n"
)

// React indicator dependencies.
var ReactIndicatorDeps = []string{
	"react",
	"@types/react",
	"react-scripts",
	"next",
	"gatsby",
	"remix",
	"@remix-run/react",
}

// C# key file patterns.
var CsharpKeyFilePatterns = []string{
	"global.json",
	"nuget.config",
	"Directory.Build.props",
	"Directory.Build.targets",
	"launchSettings.json",
	"appsettings.json",
}

// Project detection exclusion directories.
//
// These directories are skipped during in-repo project detection walks.
// Skipping noisy/generated trees (especially node_modules, .git, vendored
// caches, and CMS upload trees) is the single most important factor in
// keeping `gitmap scan` fast on real-world projects — a WordPress repo
// can easily contain 100k+ files under wp-content alone.
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
	".next",
	".nuxt",
	".svelte-kit",
	".turbo",
	".parcel-cache",
	".angular",
	"coverage",
	".nyc_output",
	"__pycache__",
	".pytest_cache",
	".mypy_cache",
	".ruff_cache",
	".tox",
	".gradle",
	".idea",
	".vs",
	".vscode-test",
	".terraform",
	".serverless",
	"tmp",
	"temp",
	"logs",
	"Pods",
	"DerivedData",
	"wp-content",
	"wp-admin",
	"wp-includes",
	"uploads",
}
