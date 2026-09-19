package cmdautomation

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// TopologyOptions configures polyglot codebase topology discovery.
type TopologyOptions struct {
	Dir       string `json:"dir"`
	Query     string `json:"query"`
	TtlSec    int    `json:"ttlSec"`
	IsJson    bool   `json:"isJson"`
	IsRefresh bool   `json:"isRefresh"`
}

// SubsystemData represents detected roots and entry files for a subsystem.
type SubsystemData struct {
	Roots       []string `json:"roots"`
	Entrypoints []string `json:"entrypoints"`
	SchemaFiles []string `json:"schemaFiles"`
	Workflows   []string `json:"workflows"`
	SpecRoots   []string `json:"specRoots"`
	TestRunners []string `json:"testRunners"`
}

// TopologyResult captures the complete discovered topology of a repository.
type TopologyResult struct {
	Version      string                   `json:"version"`
	GeneratedAt  string                   `json:"generatedAt"`
	ExpiresAt    string                   `json:"expiresAt"`
	TtlSeconds   int                      `json:"ttlSeconds"`
	DurationMs   float64                  `json:"durationMs"`
	TotalFiles   int                      `json:"totalFiles"`
	RootPath     string                   `json:"rootPath"`
	Manifests    map[string][]string      `json:"manifests"`
	Languages    map[string]int           `json:"languages"`
	LangRoots    map[string][]string      `json:"langRoots"`
	Subsystems   map[string]SubsystemData `json:"subsystems"`
	Duration     time.Duration            `json:"duration"`
	IsValid      bool                     `json:"isValid"`
}

// DbGenerateOptions configures Go/TS struct and enum generation from SQLite.
type DbGenerateOptions struct {
	DbPath    string `json:"dbPath"`
	OutDir    string `json:"outDir"`
	Lang      string `json:"lang"`
	IsDryRun  bool   `json:"isDryRun"`
	StructDir string `json:"structDir"`
}

// DbGenerateResult summarizes generated code artifacts.
type DbGenerateResult struct {
	GeneratedFiles []string      `json:"generatedFiles"`
	TableCount     int           `json:"tableCount"`
	StructCount    int           `json:"structCount"`
	EnumCount      int           `json:"enumCount"`
	Duration       time.Duration `json:"duration"`
	IsDryRun       bool          `json:"isDryRun"`
	IsSuccess      bool          `json:"isSuccess"`
}

// DbMigrateOptions configures SQLite database migration execution.
type DbMigrateOptions struct {
	DbPath        string `json:"dbPath"`
	MigrationsDir string `json:"migrationsDir"`
	SqlScript     string `json:"sqlScript"`
	RollbackStep  int    `json:"rollbackStep"`
	IsDryRun      bool   `json:"isDryRun"`
	IsStatus      bool   `json:"isStatus"`
}

// DbMigrateResult summarizes applied database migrations.
type DbMigrateResult struct {
	AppliedMigrations []string      `json:"appliedMigrations"`
	TotalApplied      int           `json:"totalApplied"`
	CurrentVersion    string        `json:"currentVersion"`
	Duration          time.Duration `json:"duration"`
	IsDryRun          bool          `json:"isDryRun"`
	IsSuccess         bool          `json:"isSuccess"`
}

// SchemaAuditOptions configures schema convention auditing for SQLite.
type SchemaAuditOptions struct {
	DbPath   string `json:"dbPath"`
	Dir      string `json:"dir"`
	IsStrict bool   `json:"isStrict"`
	IsJson   bool   `json:"isJson"`
}

// SchemaViolation records a single schema convention rule breach.
type SchemaViolation struct {
	File     string `json:"file"`
	Table    string `json:"table"`
	Issue    string `json:"issue"`
	Severity string `json:"severity"`
}

// SchemaAuditResult aggregates schema audit findings.
type SchemaAuditResult struct {
	ScannedFiles int               `json:"scannedFiles"`
	TableCount   int               `json:"tableCount"`
	Violations   []SchemaViolation `json:"violations"`
	Duration     time.Duration     `json:"duration"`
	IsClean      bool              `json:"isClean"`
}

// TableColumnInfo describes a single column in a SQLite table.
type TableColumnInfo struct {
	Cid       int    `json:"cid"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsNotNull bool   `json:"isNotNull"`
	DfltValue string `json:"dfltValue"`
	IsPk      bool   `json:"isPk"`
}

// TableSchemaInfo describes a table and its columns.
type TableSchemaInfo struct {
	Name    string            `json:"name"`
	Columns []TableColumnInfo `json:"columns"`
	SqlBody string            `json:"sqlBody"`
}

type (
	// TopologyResultMonad wraps TopologyResult with AppError.
	TopologyResultMonad = result.Result[TopologyResult]

	// DbGenerateResultMonad wraps DbGenerateResult with AppError.
	DbGenerateResultMonad = result.Result[DbGenerateResult]

	// DbMigrateResultMonad wraps DbMigrateResult with AppError.
	DbMigrateResultMonad = result.Result[DbMigrateResult]

	// SchemaAuditResultMonad wraps SchemaAuditResult with AppError.
	SchemaAuditResultMonad = result.Result[SchemaAuditResult]
)
