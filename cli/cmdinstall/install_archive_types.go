package cmdinstall

import "context"

// ArchiveInstallStrategy represents detected deployment mechanism.
type ArchiveInstallStrategy string

const (
	StrategyBinaryApp ArchiveInstallStrategy = "binary_app"
	StrategyScript    ArchiveInstallStrategy = "install_script"
	StrategySource    ArchiveInstallStrategy = "source_build"
	StrategySingleGz  ArchiveInstallStrategy = "single_gzip"
	StrategyUnknown   ArchiveInstallStrategy = "unknown"
)

// ArchiveInstallOptions holds options for archive installation.
type ArchiveInstallOptions struct {
	ArchivePath string
	AppName     string
	DestDir     string
	BinDir      string
	DesktopDir  string
	Verbose     bool
	DryRun      bool
}

// ArchiveInspectionResult holds details discovered during archive analysis.
type ArchiveInspectionResult struct {
	Strategy       ArchiveInstallStrategy
	BinaryPath     string
	ScriptPath     string
	DesktopPath    string
	IconPath       string
	SourceRoot     string
	SuggestedName  string
	HasMakefile    bool
	HasConfigure   bool
	HasDesktopFile bool
}

// ArchiveStepProgressFunc reports step-by-step progress to the terminal.
type ArchiveStepProgressFunc func(step, total int, message string)

// ArchiveDeployParams holds all parameters needed for final deployment.
type ArchiveDeployParams struct {
	Ctx        context.Context
	Inspection ArchiveInspectionResult
	Opts       ArchiveInstallOptions
	Progress   ArchiveStepProgressFunc
}
