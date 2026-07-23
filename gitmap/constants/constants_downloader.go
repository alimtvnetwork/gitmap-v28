package constants

// gitmap:cmd top-level
// Downloader CLI commands and shorthand aliases.
//
// Slice 1 of the downloader feature only wires `downloader-config` (and its
// `dc` shorthand). The `download` / `download-unzip` / `dl` / `du` aliases
// are reserved here so a later slice can add them without touching the
// constants surface twice.
const (
	CmdDownloaderConfig      = "downloader-config"
	CmdDownloaderConfigAlias = "dc"

	// Reserved for Slice 3 — kept here so the dispatch table grows in
	// one file and shorthand collisions can be detected by the existing
	// constants-collision check.
	CmdDownload           = "download"       // gitmap:cmd skip
	CmdDownloadAlias      = "dl"             // gitmap:cmd skip
	CmdDownloadUnzip      = "download-unzip" // gitmap:cmd skip
	CmdDownloadUnzipAlias = "du"             // gitmap:cmd skip
)

// SettingType enum (stored as a TEXT discriminator on each Setting key).
//
// The current Setting table is keyed by Key (TEXT PK) and stores values as
// TEXT, so we don't need a separate SettingTypes table — the type lives in
// code as the SettingType const block. Each new key MUST declare a type so
// downstream tooling (db inspect, future migrations) can route by category
// without parsing the JSON value.
type SettingType string

const (
	SettingTypeDownloaderConfig SettingType = "DownloaderConfig"
	SettingTypeDatabaseVersion  SettingType = "DatabaseVersion"
	SettingTypeSystemConfig     SettingType = "SystemConfig"
)

// Settings keys for the downloader feature.
//
// SettingDownloaderConfig holds the full PascalCase JSON blob produced by
// the Seedable-Config (data/downloader-config.json). Stored as TEXT so the
// existing Setting upsert path is reused unchanged.
//
// SettingDownloaderConfigSeedHash holds the SHA-256 of the seed file the
// last time it was applied, so the seeder can detect upstream changes
// without overwriting user-customized configs (unless OverwriteUserConfig
// is enabled in the seed).
const (
	SettingDownloaderConfig         = "DownloaderConfig"
	SettingDownloaderConfigSeedHash = "DownloaderConfigSeedHash"
	SettingDatabaseVersion          = "DatabaseVersion"
)

// Default Seedable-Config path, relative to the gitmap binary's data dir.
const DefaultDownloaderConfigSeedPath = "./data/downloader-config.json"

// Downloader hard-coded fallbacks. Only used when both the DB and the seed
// file are unavailable (e.g., first-run race before Migrate completes).
// Numbers mirror the spec: 15 parallel / 15 splits / 800KB normal / 2MB
// large / 100MB threshold. Tiny-file profile (<2MB → 8 splits @ 100KB)
// matches the user clarification.
const (
	DownloaderDefaultParallel       = 15
	DownloaderDefaultSplits         = 15
	DownloaderDefaultSplitSize      = "800K"
	DownloaderDefaultLargeSplitSize = "2M"
	DownloaderDefaultLargeThreshold = "100M"
	DownloaderDefaultTinyThreshold  = "2M"
	DownloaderDefaultTinySplitSize  = "100K"
	DownloaderDefaultTinySplits     = 8
	DownloaderDefaultPreferred      = "Aria2C"
	DownloaderDefaultFallback       = "GoDownloader"
	DownloaderDefaultAllowFallback  = true
	DownloaderDefaultOverwriteUser  = false
)

// User-facing messages for downloader-config (Slice 1).
const (
	MsgDownloaderConfigBanner       = "▶ gitmap downloader-config v%s"
	MsgDownloaderConfigLoaded       = "  ✓ Loaded downloader config from %s"
	MsgDownloaderConfigSaved        = "  ✓ Saved downloader config to Setting[%s]"
	MsgDownloaderConfigPromptHeader = "  ▸ No file path supplied — entering interactive mode. Press <Enter> to keep the shown default."
	MsgDownloaderConfigSeeded       = "  ✓ Seeded downloader defaults (seed hash %s)"
	MsgDownloaderConfigSeedSkip     = "  ◦ Downloader config already customized (OverwriteUserConfig=false) — skipping seed"
	MsgDownloaderConfigDBVersion    = "  ✓ Recorded LastKnownVersion=%s"
	WarnDownloaderSeedRead          = "  ⚠ Could not read downloader seed at %s: %v"
	WarnDownloaderSeedParse         = "  ⚠ Could not parse downloader seed: %v"
	ErrDownloaderConfigPathRequired = "downloader-config: file path %q does not exist"
	ErrDownloaderConfigInvalidJSON  = "downloader-config: invalid JSON: %w"
	ErrDownloaderConfigMissingKey   = "downloader-config: missing required PascalCase key %q"
	ErrDownloaderConfigBadParallel  = "downloader-config: ParallelDownloads must be 1..64, got %d"
	ErrDownloaderConfigBadSplits    = "downloader-config: SplitConnections must be 1..64, got %d"
)

// HelpDownloaderConfig is shown in the root usage table once the command is
// listed by the help generator.
const HelpDownloaderConfig = "  downloader-config (dc) [path]  Set/seed downloader config (aria2c parallel, splits, thresholds)"
