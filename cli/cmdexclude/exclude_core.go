package cmdexclude

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func init() {
	fsutil.RegisterCustomExclusionChecker(IsCustomExcluded)
}

// IsCustomExcluded tests whether a path matches any saved custom exclusion pattern.
func IsCustomExcluded(path string) bool {
	excls, err := LoadCustomExclusions()
	if err != nil || len(excls) == 0 {
		return false
	}
	norm := strings.ToLower(strings.ReplaceAll(path, "\\", "/"))
	for _, excl := range excls {
		target := strings.ToLower(strings.ReplaceAll(excl, "\\", "/"))
		if strings.Contains(norm, target) || strings.EqualFold(norm, target) {
			return true
		}
	}
	return false
}

// LoadCustomExclusions reads user-configured exclusion patterns from Split-DB.
func LoadCustomExclusions() ([]string, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	val := db.GetSetting(constants.SettingCustomRepoExclusions)
	if strings.TrimSpace(val) == "" {
		return nil, nil
	}

	var patterns []string
	if errJSON := json.Unmarshal([]byte(val), &patterns); errJSON != nil {
		return nil, errJSON
	}
	return patterns, nil
}

// SaveCustomExclusions saves user-configured exclusion patterns into Split-DB.
func SaveCustomExclusions(patterns []string) error {
	db, err := store.OpenDefault()
	if err != nil {
		return err
	}
	defer db.Close()

	data, errJSON := json.Marshal(patterns)
	if errJSON != nil {
		return errJSON
	}

	return db.SetSetting(constants.SettingCustomRepoExclusions, string(data))
}
