package cmdinstall

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall/enums"
)

// ArchiveInstallStrategyType is an alias for the canonical enum in the
// package's enums subpackage, kept so existing callers keep compiling.
type ArchiveInstallStrategyType = enums.ArchiveInstallStrategyType

// Archive install strategy constants, aliased from the enums subpackage.
const (
	StrategyBinaryApp = enums.StrategyBinaryApp
	StrategyScript    = enums.StrategyScript
	StrategySource    = enums.StrategySource
	StrategySingleGz  = enums.StrategySingleGz
	StrategyUnknown   = enums.StrategyUnknown
)

// ArchiveInstallOptions holds options for archive installation.
type ArchiveInstallOptions struct {
	ArchivePath    string
	AppName        string
	DestDir        string
	BinDir         string
	DesktopDir     string
	Verbose        bool
	DryRun         bool
	IsDownloadMust bool
}

// ArchiveDownloadParams holds parameters for fetching or caching an archive.
type ArchiveDownloadParams struct {
	URL            string
	DestPath       string
	IsDownloadMust bool
	Verbose        bool
}

// ArchiveCacheValidationResult holds validation outcome of a cached archive file.
type ArchiveCacheValidationResult struct {
	IsValid       bool
	Size          int64
	HasValidMagic bool
	ErrorMessage  string
}

// ArchiveInspectionResult holds details discovered during archive analysis.
type ArchiveInspectionResult struct {
	Strategy       enums.ArchiveInstallStrategyType
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
