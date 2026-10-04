package cmdagent

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

const (
	Tier1MasterDBName = "ai_agents.db"
	Tier2TaskDBName   = "agent-task.db"
	Tier2TaskDBAlt    = "task.db"
	AgentsSubdirName  = "agents"
	SettingAgentDir   = "agent.temp_dir"
	EnvAgentTempDir   = "GITMAP_AGENT_TEMP_DIR"
)

// ResolveAgentTempDir implements priority: --dir flag -> env -> db setting -> default.
func ResolveAgentTempDir(dirFlag string) string {
	hasDirFlag := strings.TrimSpace(dirFlag) != ""
	if hasDirFlag {
		return filepath.Clean(strings.TrimSpace(dirFlag))
	}

	envDir := resolveEnvTempDir()
	hasEnvDir := envDir != ""
	if hasEnvDir {
		return envDir
	}

	settingDir := resolveSettingTempDir()
	hasSettingDir := settingDir != ""
	if hasSettingDir {
		return settingDir
	}

	return resolveDefaultTempDir()
}

func resolveEnvTempDir() string {
	envVal := strings.TrimSpace(os.Getenv(EnvAgentTempDir))
	hasEnvVal := envVal != ""
	if hasEnvVal {
		return filepath.Clean(envVal)
	}

	return ""
}

func resolveSettingTempDir() string {
	db, err := store.OpenDefault()
	hasErr := err != nil
	if hasErr {
		return ""
	}
	defer db.Close()

	val := strings.TrimSpace(db.GetSetting(SettingAgentDir))
	hasVal := val != ""
	if hasVal {
		return filepath.Clean(val)
	}

	return ""
}

func resolveDefaultTempDir() string {
	repoRoot := FindRepoRoot()
	defaultPath := filepath.Join(repoRoot, ".ai-memory", "temp-agents")

	return filepath.Clean(defaultPath)
}

// FindRepoRoot traverses upward from the current directory to find repo root.
func FindRepoRoot() string {
	cwd, err := os.Getwd()
	hasErr := err != nil
	if hasErr {
		return "."
	}

	dir := cwd
	for {
		isRoot := isRepoRoot(dir)
		if isRoot {
			return dir
		}

		parent := filepath.Dir(dir)
		isTop := parent == dir
		if isTop {
			break
		}
		dir = parent
	}

	return cwd
}

func isRepoRoot(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	info, err := os.Stat(gitDir)
	hasGit := err == nil && info.IsDir()
	if hasGit {
		return true
	}

	versionFile := filepath.Join(dir, "version.json")
	_, verErr := os.Stat(versionFile)

	return verErr == nil
}

// Slugify generates a deterministic lowercase kebab-case slug.
func Slugify(text string) string {
	clean := cleanSlugRunes(strings.ToLower(strings.TrimSpace(text)))
	words := strings.Fields(clean)
	hasWords := len(words) > 0
	if !hasWords {
		return "agent"
	}

	limit := len(words)
	hasExcess := limit > 6
	if hasExcess {
		limit = 6
	}

	return strings.Join(words[:limit], "-")
}

func cleanSlugRunes(text string) string {
	var builder strings.Builder
	for _, r := range text {
		isLetterOrDigit := unicode.IsLetter(r) || unicode.IsDigit(r)
		if isLetterOrDigit {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}

	return builder.String()
}

// NowISO returns current UTC timestamp in ISO 8601 format.
func NowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// OpenSqliteDB opens an SQLite connection with WAL mode and busy timeout.
func OpenSqliteDB(dbPath string) (*sql.DB, *appfault.AppError) {
	parentDir := filepath.Dir(dbPath)
	mkdirErr := os.MkdirAll(parentDir, 0755)
	hasMkdirErr := mkdirErr != nil
	if hasMkdirErr {
		return nil, appfault.WrapExecution(mkdirErr, "failed to create database directory")
	}

	connStr := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", filepath.ToSlash(dbPath))
	db, err := sql.Open("sqlite", connStr)
	hasErr := err != nil
	if hasErr {
		return nil, appfault.WrapExecution(err, "failed to open sqlite database")
	}

	return db, nil
}

// InitTier1Schema initializes schema for ai_agents.db master registry.
func InitTier1Schema(db *sql.DB) *appfault.AppError {
	const ddl = `
	CREATE TABLE IF NOT EXISTS ParentTaskRegistry (
		ParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
		TaskSlug TEXT NOT NULL,
		TaskName TEXT NOT NULL,
		RunDirectory TEXT NOT NULL,
		RootDbPath TEXT NOT NULL,
		Status TEXT NOT NULL DEFAULT 'ACTIVE',
		TotalStepsBudget INTEGER NOT NULL DEFAULT 300,
		CompletedSteps INTEGER NOT NULL DEFAULT 0,
		SpawnedAgentCount INTEGER NOT NULL DEFAULT 0,
		CreatedAt TEXT NOT NULL,
		UpdatedAt TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS AgentRegistry (
		AgentId INTEGER PRIMARY KEY AUTOINCREMENT,
		ParentTaskId INTEGER NOT NULL,
		AgentRole TEXT NOT NULL,
		AgentSlug TEXT NOT NULL,
		SplitDbPath TEXT NOT NULL,
		Status TEXT NOT NULL DEFAULT 'IDLE',
		LastHeartbeatAt TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS GlobalLifecycleMetrics (
		MetricKey TEXT PRIMARY KEY,
		MetricValue TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_parent_task_status ON ParentTaskRegistry(Status);
	CREATE INDEX IF NOT EXISTS idx_agent_parent ON AgentRegistry(ParentTaskId);
	`
	_, err := db.Exec(ddl)
	hasErr := err != nil
	if hasErr {
		return appfault.WrapExecution(err, "failed to init tier1 master schema")
	}

	return nil
}

// InitTier2Schema initializes schema for run-scoped agent-task.db.
func InitTier2Schema(db *sql.DB) *appfault.AppError {
	const ddl = `
	CREATE TABLE IF NOT EXISTS ParentTask (
		ParentTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
		TaskName TEXT NOT NULL,
		TaskSlug TEXT NOT NULL,
		RunDirectory TEXT NOT NULL,
		Status TEXT NOT NULL DEFAULT 'ACTIVE',
		IsActive INTEGER NOT NULL DEFAULT 1,
		HasCompleted INTEGER NOT NULL DEFAULT 0,
		TotalStepsBudget INTEGER NOT NULL DEFAULT 300,
		CurrentStep INTEGER NOT NULL DEFAULT 1,
		Notes TEXT NULL,
		CreatedAt TEXT NOT NULL,
		UpdatedAt TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS Subtask (
		SubtaskId INTEGER PRIMARY KEY AUTOINCREMENT,
		ParentTaskId INTEGER NOT NULL,
		TaskCode TEXT NOT NULL,
		Title TEXT NOT NULL,
		AssignedAgentRole TEXT NULL,
		OwnedFilesJson TEXT NOT NULL DEFAULT '[]',
		Status TEXT NOT NULL DEFAULT 'PENDING',
		IsBlocked INTEGER NOT NULL DEFAULT 0,
		HasCompleted INTEGER NOT NULL DEFAULT 0,
		Evidence TEXT NULL,
		CreatedAt TEXT NOT NULL,
		UpdatedAt TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS AgentActionLog (
		ActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
		SubtaskId INTEGER NOT NULL,
		AgentRole TEXT NOT NULL,
		ActionType TEXT NOT NULL,
		TargetFile TEXT NULL,
		StartLine INTEGER NULL DEFAULT 0,
		EndLine INTEGER NULL DEFAULT 0,
		QueryOrCommand TEXT NULL,
		ActionDetails TEXT NOT NULL,
		DurationMs INTEGER NOT NULL DEFAULT 0,
		Status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
		ErrorMessage TEXT NULL,
		CreatedAt TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_subtask_status ON Subtask(Status);
	CREATE INDEX IF NOT EXISTS idx_action_subtask ON AgentActionLog(SubtaskId);
	`
	_, err := db.Exec(ddl)
	hasErr := err != nil
	if hasErr {
		return appfault.WrapExecution(err, "failed to init tier2 task schema")
	}

	return nil
}

// InitTier3Schema initializes schema for agent private split-db.
func InitTier3Schema(db *sql.DB) *appfault.AppError {
	const ddl = `
	CREATE TABLE IF NOT EXISTS AgentActionLog (
		ActionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
		SubtaskId INTEGER NOT NULL,
		AgentRole TEXT NOT NULL,
		ActionType TEXT NOT NULL,
		TargetFile TEXT NULL,
		StartLine INTEGER NULL DEFAULT 0,
		EndLine INTEGER NULL DEFAULT 0,
		QueryOrCommand TEXT NULL,
		ActionDetails TEXT NOT NULL,
		DurationMs INTEGER NOT NULL DEFAULT 0,
		Status TEXT NOT NULL DEFAULT 'IN_PROGRESS',
		ErrorMessage TEXT NULL,
		CreatedAt TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_agent_action_subtask ON AgentActionLog(SubtaskId);
	CREATE INDEX IF NOT EXISTS idx_agent_action_type ON AgentActionLog(ActionType);
	CREATE INDEX IF NOT EXISTS idx_agent_action_created ON AgentActionLog(CreatedAt);
	`
	_, err := db.Exec(ddl)
	hasErr := err != nil
	if hasErr {
		return appfault.WrapExecution(err, "failed to init tier3 agent schema")
	}

	return nil
}

// ResolveTaskDir finds task run folder by task identifier or finds latest run.
func ResolveTaskDir(tempDir, taskFlag string) (string, *appfault.AppError) {
	hasTaskFlag := strings.TrimSpace(taskFlag) != ""
	if hasTaskFlag {
		matchedDir, isMatch := findTaskDirByPattern(tempDir, strings.TrimSpace(taskFlag))
		if isMatch {
			return matchedDir, nil
		}

		directPath := filepath.Join(tempDir, strings.TrimSpace(taskFlag))
		info, err := os.Stat(directPath)
		hasDirect := err == nil && info.IsDir()
		if hasDirect {
			return directPath, nil
		}

		return "", appfault.NewNotFoundError("task directory not found for: " + taskFlag)
	}

	return findLatestRunDir(tempDir)
}

func findTaskDirByPattern(tempDir, pattern string) (string, bool) {
	entries, err := os.ReadDir(tempDir)
	hasErr := err != nil
	if hasErr {
		return "", false
	}

	patLower := strings.ToLower(pattern)
	for _, entry := range entries {
		isDir := entry.IsDir()
		if !isDir {
			continue
		}

		nameLower := strings.ToLower(entry.Name())
		isMatch := nameLower == patLower || strings.Contains(nameLower, patLower)
		if isMatch {
			return filepath.Join(tempDir, entry.Name()), true
		}
	}

	return "", false
}

func findLatestRunDir(tempDir string) (string, *appfault.AppError) {
	dirs := scanRunDirs(tempDir)
	hasDirs := len(dirs) > 0
	if !hasDirs {
		return "", appfault.NewNotFoundError("no active or recent task run directories found in: " + tempDir)
	}

	return dirs[len(dirs)-1], nil
}

func scanRunDirs(tempDir string) []string {
	entries, err := os.ReadDir(tempDir)
	hasErr := err != nil
	if hasErr {
		return nil
	}

	var results []string
	for _, entry := range entries {
		isDir := entry.IsDir()
		if !isDir {
			continue
		}
		results = append(results, filepath.Join(tempDir, entry.Name()))
	}
	sort.Strings(results)

	return results
}

// OpenTier1DB opens and initializes the root master ai_agents.db.
func OpenTier1DB(tempDir string) (*sql.DB, *appfault.AppError) {
	dbPath := filepath.Join(tempDir, Tier1MasterDBName)
	db, err := OpenSqliteDB(dbPath)
	hasErr := err != nil
	if hasErr {
		return nil, err
	}

	initErr := InitTier1Schema(db)
	hasInitErr := initErr != nil
	if hasInitErr {
		_ = db.Close()
		return nil, initErr
	}

	return db, nil
}

// OpenTier2DB opens and initializes run-scoped agent-task.db.
func OpenTier2DB(taskDir string) (*sql.DB, *appfault.AppError) {
	dbPath := filepath.Join(taskDir, Tier2TaskDBName)
	db, err := OpenSqliteDB(dbPath)
	hasErr := err != nil
	if hasErr {
		return nil, err
	}

	initErr := InitTier2Schema(db)
	hasInitErr := initErr != nil
	if hasInitErr {
		_ = db.Close()
		return nil, initErr
	}

	return db, nil
}

// OpenTier3DB opens and initializes agent private split-db in agents/<agent-slug>.db.
func OpenTier3DB(taskDir, agentRole string) (*sql.DB, *appfault.AppError) {
	agentSlug := Slugify(agentRole)
	agentsDir := filepath.Join(taskDir, AgentsSubdirName)
	dbPath := filepath.Join(agentsDir, fmt.Sprintf("%s.db", agentSlug))

	db, err := OpenSqliteDB(dbPath)
	hasErr := err != nil
	if hasErr {
		return nil, err
	}

	initErr := InitTier3Schema(db)
	hasInitErr := initErr != nil
	if hasInitErr {
		_ = db.Close()
		return nil, initErr
	}

	return db, nil
}
