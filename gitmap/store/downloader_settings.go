package store

import (
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/downloaderconfig"
)

// GetDownloaderConfig returns the persisted downloader Document, or
// ok=false when no SettingDownloaderConfig row exists yet. Callers that
// need a guaranteed-non-empty config should fall back to
// downloaderconfig.Defaults() when ok is false.
//
// Stored as a JSON TEXT blob under the existing Setting(Key,Value) table —
// no new schema is required, which is why Migrate's SchemaVersionCurrent
// did not need to bump for this feature.
func (db *DB) GetDownloaderConfig() (downloaderconfig.Document, bool) {
	raw := db.GetSetting(constants.SettingDownloaderConfig)
	if raw == "" {
		return downloaderconfig.Document{}, false
	}

	var doc downloaderconfig.Document
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		// Corrupt blob → behave as "absent" so the seeder can replace it
		// on the next Migrate without dragging the user into a parse error
		// for an unrelated subcommand.
		return downloaderconfig.Document{}, false
	}

	return doc, true
}

// SetDownloaderConfig upserts the full Document as a JSON blob and
// records the LastKnownVersion under SettingDatabaseVersion so future
// runs can see which gitmap build last touched the DB.
func (db *DB) SetDownloaderConfig(doc downloaderconfig.Document) error {
	raw, err := downloaderconfig.Marshal(doc)
	if err != nil {
		return err
	}

	if err := db.SetSetting(constants.SettingDownloaderConfig, string(raw)); err != nil {
		return err
	}

	return db.SetSetting(constants.SettingDatabaseVersion, doc.DatabaseVersion.LastKnownVersion)
}

// SetDownloaderSeedHash records the canonical hash of the seed that was
// last applied. Compared on the next Migrate to detect upstream seed
// changes without re-parsing the file every run.
func (db *DB) SetDownloaderSeedHash(hash string) error {
	return db.SetSetting(constants.SettingDownloaderConfigSeedHash, hash)
}

// GetDownloaderSeedHash returns the hash recorded by the last successful
// seed, or "" when no seed has been applied yet.
func (db *DB) GetDownloaderSeedHash() string {
	return db.GetSetting(constants.SettingDownloaderConfigSeedHash)
}
