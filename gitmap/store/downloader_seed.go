package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/downloaderconfig"
)

// SeedDownloaderConfig applies the Seedable-Config to the database.
//
// Behavior matrix (matches spec §Seedable-Config Requirements):
//
//	┌──────────────────────────────┬──────────────────────────────────────┐
//	│ State                        │ Action                               │
//	├──────────────────────────────┼──────────────────────────────────────┤
//	│ No DB row + no seed file     │ Persist downloaderconfig.Defaults()  │
//	│ No DB row + seed file        │ Persist parsed seed                  │
//	│ DB row, seed unchanged       │ No-op (hash matches)                 │
//	│ DB row, seed changed,        │ Skip (preserve user customization)   │
//	│   OverwriteUserConfig=false  │                                      │
//	│ DB row, seed changed,        │ Overwrite with new seed              │
//	│   OverwriteUserConfig=true   │                                      │
//	└──────────────────────────────┴──────────────────────────────────────┘
//
// Called from Migrate() AFTER schema-version is stamped so a fresh install
// gets seeded on first run. Failures are warned to stderr but never fatal:
// an unseeded DB just means the next subcommand falls back to the
// hard-coded defaults until the user runs `gitmap downloader-config`.
func (db *DB) SeedDownloaderConfig(seedPath string) {
	doc, hash := loadSeedOrDefaults(seedPath)

	existing, hasExisting := db.GetDownloaderConfig()
	prevHash := db.GetDownloaderSeedHash()

	switch {
	case !hasExisting:
		// First run — always seed.
	case prevHash == hash:
		// Seed unchanged since last apply — nothing to do.
		return
	case !existing.DownloaderConfig.OverwriteUserConfig && !doc.DownloaderConfig.OverwriteUserConfig:
		// User has customized config and neither side opts into overwrite.
		// Stamp the new hash so we don't re-evaluate this every run, but
		// keep the user's values intact.
		_ = db.SetDownloaderSeedHash(hash)
		if !isSeedBannerSuppressed() {
			fmt.Fprintln(os.Stdout, constants.MsgDownloaderConfigSeedSkip)
		}

		return
	}

	if err := db.SetDownloaderConfig(doc); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not seed downloader config: %v\n", err)

		return
	}
	if err := db.SetDownloaderSeedHash(hash); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record seed hash: %v\n", err)
	}
	if isSeedBannerSuppressed() {
		return
	}
	fmt.Fprintf(os.Stdout, constants.MsgDownloaderConfigSeeded+"\n", short(hash))
	fmt.Fprintf(os.Stdout, constants.MsgDownloaderConfigDBVersion+"\n", doc.DatabaseVersion.LastKnownVersion)
}

// isSeedBannerSuppressed mirrors cmd.isBannerSuppressed without importing
// cmd (which would create a cycle). The seed banner is informational, not
// an error, so it is silenced for `version`, `--no-banner`, GITMAP_QUIET=1,
// and any non-interactive subcommand probe.
func isSeedBannerSuppressed() bool {
	if os.Getenv(constants.EnvGitMapQuiet) == constants.EnvGitMapQuietTrue {
		return true
	}
	for _, arg := range os.Args[1:] {
		if arg == constants.FlagNoBanner {
			return true
		}
		if arg == "version" || arg == "--version" || arg == "-v" {
			return true
		}
	}

	return false
}

// loadSeedOrDefaults resolves the seed file relative to the active
// binary's data dir. Returns the Defaults() document when the file is
// missing or unparseable so the seeder still persists a usable baseline
// on first run.
func loadSeedOrDefaults(seedPath string) (downloaderconfig.Document, string) {
	resolved := resolveSeedPath(seedPath)
	doc, err := downloaderconfig.LoadFile(resolved)
	if err != nil {
		// A missing optional seed is the normal first-run state on a fresh
		// install — silently fall back to defaults. Anything else is a real
		// I/O / parse failure worth surfacing.
		if !errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintf(os.Stderr, constants.WarnDownloaderSeedRead+"\n", resolved, err)
		}
		fallback := downloaderconfig.Defaults()

		return fallback, downloaderconfig.SeedHash(fallback)
	}

	return doc, downloaderconfig.SeedHash(doc)
}

// resolveSeedPath turns a relative seed path into an absolute path
// anchored to the running binary's directory. Falls back to the
// caller-supplied path on any os.Executable failure.
func resolveSeedPath(seedPath string) string {
	if filepath.IsAbs(seedPath) {
		return seedPath
	}
	exe, err := os.Executable()
	if err != nil {
		return seedPath
	}

	return filepath.Join(filepath.Dir(exe), seedPath)
}

// short trims a hex hash to its first 12 chars for log readability.
func short(hash string) string {
	if len(hash) <= 12 {
		return hash
	}

	return hash[:12]
}
