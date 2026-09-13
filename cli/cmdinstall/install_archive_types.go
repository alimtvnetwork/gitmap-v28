package cmdinstall

import "context"

// ArchiveInstallStrategyType represents detected deployment mechanism.
type ArchiveInstallStrategyType string

const (
	StrategyBinaryApp ArchiveInstallStrategyType = "binary_app"
	StrategyScript    ArchiveInstallStrategyType = "install_script"
	StrategySource    ArchiveInstallStrategyType = "source_build"
	StrategySingleGz  ArchiveInstallStrategyType = "single_gzip"
	StrategyUnknown   ArchiveInstallStrategyType = "unknown"
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
	Strategy       ArchiveInstallStrategyType
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
