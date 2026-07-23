// Package model defines the core data structures for gitmap.
package model

import "github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

// ScanRecord holds all information about a discovered Git repository.
//
// Depth is the directory level at which the repo was found, counted
// from the scan root (0 = the scan root itself, 1 = its immediate
// children, …). Surfaced in CSV / JSON output so users can verify
// when DefaultMaxDepth (or a custom --max-depth) prevented walking
// into deeper directories: a row with Depth == cap is a candidate
// for a deeper rescan.
//
// Transport is the URL-scheme bucket the repo's discovered remote
// falls into: one of "ssh" | "https" | "other". Surfaced as a CSV
// column / JSON field so users can filter clones by transport with
// a one-liner (`awk -F, '$13=="ssh"' gitmap.csv`, `jq '.[]|
// select(.transport=="ssh")' gitmap.json`). Mirrors the same three-
// bucket collapse that the clone-from terminal summary uses (see
// clonefrom.TransportTally) so the two views stay in lockstep.
type ScanRecord struct {
	ID               int64  `json:"id"                csv:"id"`
	Slug             string `json:"slug"              csv:"slug"`
	RepoID           string `json:"repoId"            csv:"repoId"`
	RepoName         string `json:"repoName"          csv:"repoName"`
	HTTPSUrl         string `json:"httpsUrl"          csv:"httpsUrl"`
	SSHUrl           string `json:"sshUrl"            csv:"sshUrl"`
	DiscoveredURL    string `json:"discoveredUrl"     csv:"discoveredUrl"`
	Branch           string `json:"branch"            csv:"branch"`
	BranchSource     string `json:"branchSource"      csv:"branchSource"`
	RelativePath     string `json:"relativePath"      csv:"relativePath"`
	AbsolutePath     string `json:"absolutePath"      csv:"absolutePath"`
	CloneInstruction string `json:"cloneInstruction"  csv:"cloneInstruction"`
	Notes            string `json:"notes"             csv:"notes"`
	Depth            int    `json:"depth"             csv:"depth"`
	Transport        string `json:"transport"         csv:"transport"`
}

// ReleaseConfig holds release-specific configuration from config.json.
type ReleaseConfig struct {
	Targets   []ReleaseTarget `json:"targets"`
	Checksums bool            `json:"checksums"`
	Compress  bool            `json:"compress"`
}

// ReleaseTarget represents a single GOOS/GOARCH pair in config.json.
type ReleaseTarget struct {
	GOOS   string `json:"goos"`
	GOARCH string `json:"goarch"`
}

// Config holds application configuration loaded from JSON and CLI flags.
type Config struct {
	DefaultMode      string        `json:"defaultMode"`
	DefaultOutput    string        `json:"defaultOutput"`
	OutputDir        string        `json:"outputDir"`
	ExcludeDirs      []string      `json:"excludeDirs"`
	Notes            string        `json:"notes"`
	Release          ReleaseConfig `json:"release"`
	DashboardRefresh int           `json:"dashboardRefresh"`
}

// DefaultConfig returns a Config with sensible built-in defaults.
func DefaultConfig() Config {

	return Config{
		DefaultMode:      constants.ModeHTTPS,
		DefaultOutput:    constants.OutputTerminal,
		OutputDir:        constants.DefaultOutputDir,
		ExcludeDirs:      []string{},
		Notes:            "",
		DashboardRefresh: constants.DefaultDashboardRefresh,
		Release: ReleaseConfig{
			Targets:   []ReleaseTarget{},
			Checksums: false,
			Compress:  false,
		},
	}
}

// CloneResult tracks the outcome of a single clone operation.
//
// Notes carries non-fatal diagnostics about how the clone was performed —
// for example, which branch-selection strategy was applied based on the
// record's BranchSource.
type CloneResult struct {
	Record  ScanRecord
	Success bool
	Error   string
	Notes   string
}

// CloneSummary aggregates results of a batch clone operation.
//
// Skipped tracks repos that were already cloned and up to date according to
// the clone cache; they are also counted in Succeeded since the desired
// state was achieved without performing a clone or pull.
type CloneSummary struct {
	Succeeded int
	Failed    int
	Cloned    []CloneResult
	Errors    []CloneResult
	Skipped   []CloneResult
}

// ScanCache stores the flags used for the last scan so rescan can replay them.
type ScanCache struct {
	Dir           string `json:"dir"`
	ConfigPath    string `json:"configPath"`
	Mode          string `json:"mode"`
	Output        string `json:"output"`
	OutFile       string `json:"outFile"`
	OutputPath    string `json:"outputPath"`
	GithubDesktop bool   `json:"githubDesktop"`
	OpenFolder    bool   `json:"openFolder"`
	Quiet         bool   `json:"quiet"`
}
