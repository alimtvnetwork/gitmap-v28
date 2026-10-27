package cmdagent

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

// DispatchTaskSQLite handles gitmap task agent SQLite operations.
func DispatchTaskSQLite(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs {
		return printTaskSQLiteHelp()
	}

	sub := strings.ToLower(args[0])
	rest := args[1:]

	switch sub {
	case "init":
		return runTaskSQLiteInit(rest)
	case "add", "add-subtasks":
		return runTaskSQLiteAdd(rest)
	case "claim":
		return runTaskSQLiteClaim(rest)
	case "complete":
		return runTaskSQLiteComplete(rest)
	case "fail":
		return runTaskSQLiteFail(rest)
	case "log-action":
		return runTaskSQLiteLogAction(rest)
	case "status":
		return runTaskSQLiteStatus(rest)
	case "schema":
		return runTaskSQLiteSchema(rest)
	case "help", "-h", "--help":
		return printTaskSQLiteHelp()
	default:
		return appfault.NewValidationError(fmt.Sprintf("unknown task subcommand: %s", sub))
	}
}

func printTaskSQLiteHelp() error {
	fmt.Println("GitMap Native SQLite Agent Task Engine")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitmap task init --name <title> [--budget <N>] [--db <path>] [--spec-file <path>] [--spec-folder <dir>]")
	fmt.Println("  gitmap task add --db <path> --tasks-json '<json>'")
	fmt.Println("  gitmap task claim --db <path> --agent <agent-role> [--subtask-id <id>]")
	fmt.Println("  gitmap task complete --db <path> --subtask-id <id> [--evidence <summary>]")
	fmt.Println("  gitmap task fail --db <path> --subtask-id <id> [--reason <summary>]")
	fmt.Println("  gitmap task log-action --db <path> --subtask-id <id> --agent <role> --action <act>")
	fmt.Println("  gitmap task status --db <path> [--json]")
	fmt.Println("  gitmap task schema")
	fmt.Println()
	return nil
}

func runTaskSQLiteInit(args []string) error {
	fs := flag.NewFlagSet("task init", flag.ContinueOnError)
	name := fs.String("name", "", "Task name or prompt description")
	budget := fs.Int("budget", 300, "Total steps budget")
	dbFlag := fs.String("db", "", "Path to SQLite database")
	specFile := fs.String("spec-file", "", "Relative path to root spec file")
	fs.StringVar(specFile, "s", "", "Relative path to root spec file (shorthand)")
	specFolder := fs.String("spec-folder", "", "Relative path to spec directory")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	hasName := strings.TrimSpace(*name) != ""
	if !hasName {
		return appfault.NewValidationError("missing required flag: --name")
	}

	repoRoot := FindRepoRoot()
	tempDir := ResolveAgentTempDir("")
	slug := Slugify(*name)

	var dbPath string
	var runDir string

	hasDb := strings.TrimSpace(*dbFlag) != ""
	if hasDb {
		dbPath = resolveTaskDbPath(repoRoot, *dbFlag)
		runDir = filepath.Dir(dbPath)
	} else {
		nextPrefix := resolveNextRunPrefix(tempDir)
		runDirName := fmt.Sprintf("%02d-%s", nextPrefix, slug)
		runDir = filepath.Join(tempDir, runDirName)
		dbPath = filepath.Join(runDir, Tier2TaskDBName)
	}

	mkErr := os.MkdirAll(runDir, 0755)
	if mkErr != nil {
		return appfault.WrapExecution(mkErr, "failed to create run directory")
	}

	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	initErr := InitTier2Schema(db)
	if initErr != nil {
		return initErr
	}

	now := NowISO()
	insertQuery := `
		INSERT INTO ParentTask (
			TaskName, TaskSlug, RunDirectory, SpecFile, SpecFolder, Status, IsActive, HasCompleted,
			TotalStepsBudget, CurrentStep, CreatedAt, UpdatedAt
		) VALUES (?, ?, ?, ?, ?, 'ACTIVE', 1, 0, ?, 1, ?, ?);
	`
	relRunDir := toRelativePath(repoRoot, runDir)
	res, execErr := db.Exec(insertQuery, *name, slug, relRunDir, *specFile, *specFolder, *budget, now, now)
	if execErr != nil {
		return appfault.WrapExecution(execErr, "failed to insert parent task")
	}

	parentId, _ := res.LastInsertId()
	syncTaskToHeatmap(repoRoot, slug, *name, relRunDir, dbPath, *specFile, *specFolder)
	out := buildTaskInitOutput(repoRoot, dbPath, relRunDir, slug, *specFile, *specFolder, parentId, *budget)
	return printJson(out)
}

func syncTaskToHeatmap(repoRoot, slug, name, relRunDir, dbPath, specFile, specFolder string) {
	recordTaskInHeatmapDB(store.ResolveHeatmapRootDbPath(repoRoot), slug, name, relRunDir, dbPath, specFile, specFolder)
	recordTaskInHeatmapDB(store.ResolveHeatmapRootDbPath(""), slug, name, relRunDir, dbPath, specFile, specFolder)
}

func recordTaskInHeatmapDB(targetPath, slug, name, relRunDir, dbPath, specFile, specFolder string) {
	h, err := store.OpenHeatmapDBAtPath(targetPath)
	if err != nil {
		return
	}
	defer h.Close()

	_ = h.RecordParentTask(store.HeatmapParentTaskRecord{
		TaskSlug:         slug,
		TaskName:         name,
		RunDirectory:     relRunDir,
		RootDbPath:       dbPath,
		RootSpecFilePath: specFile,
		TaskFolder:       specFolder,
		Status:           "IN_PROGRESS",
	})
}

func runTaskSQLiteSchema(args []string) error {
	schemaText := `
Schema Specification: agent-task.db (Tier 2 Run-Scoped SQLite Database)

1. ParentTask Table:
   - ParentTaskId: INTEGER PRIMARY KEY AUTOINCREMENT
   - TaskName: TEXT NOT NULL
   - TaskSlug: TEXT NOT NULL
   - RunDirectory: TEXT NOT NULL
   - Status: TEXT NOT NULL DEFAULT 'ACTIVE'
   - IsActive: INTEGER NOT NULL DEFAULT 1
   - HasCompleted: INTEGER NOT NULL DEFAULT 0
   - TotalStepsBudget: INTEGER NOT NULL DEFAULT 300
   - CurrentStep: INTEGER NOT NULL DEFAULT 1
   - Notes: TEXT NULL
   - CreatedAt: TEXT NOT NULL (ISO 8601 UTC)
   - UpdatedAt: TEXT NOT NULL (ISO 8601 UTC)

2. Subtask Table:
   - SubtaskId: INTEGER PRIMARY KEY AUTOINCREMENT
   - ParentTaskId: INTEGER NOT NULL (FK -> ParentTask.ParentTaskId)
   - TaskCode: TEXT NOT NULL (e.g. Task-01)
   - Title: TEXT NOT NULL
   - AssignedAgentRole: TEXT NULL (e.g. Worker 01)
   - OwnedFilesJson: TEXT NOT NULL DEFAULT '[]'
   - Status: TEXT NOT NULL DEFAULT 'PENDING' (PENDING, IN_PROGRESS, DONE, FAILED)
   - IsBlocked: INTEGER NOT NULL DEFAULT 0
   - HasCompleted: INTEGER NOT NULL DEFAULT 0
   - Evidence: TEXT NULL
   - CreatedAt: TEXT NOT NULL (ISO 8601 UTC)
   - UpdatedAt: TEXT NOT NULL (ISO 8601 UTC)

3. AgentActionLog Table:
   - ActionLogId: INTEGER PRIMARY KEY AUTOINCREMENT
   - SubtaskId: INTEGER NOT NULL (FK -> Subtask.SubtaskId)
   - AgentRole: TEXT NOT NULL
   - ActionType: TEXT NOT NULL (CLAIM, COMPLETE, FAIL, EXEC, etc.)
   - TargetFile: TEXT NULL
   - ActionDetails: TEXT NOT NULL
   - Status: TEXT NOT NULL DEFAULT 'IN_PROGRESS'
   - CreatedAt: TEXT NOT NULL (ISO 8601 UTC)
`
	fmt.Println(schemaText)
	return nil
}
