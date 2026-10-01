package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const (
	// SettingGitIgnoreCheckInterval is the database settings key.
	SettingGitIgnoreCheckInterval = "gitignore_check_interval"

	// DefaultGitIgnoreCheckInterval is the default duration between audits (24 hours).
	DefaultGitIgnoreCheckInterval = "24h"
)

// Settings models user-configurable application preferences.
type Settings struct {
	GitIgnoreCheckInterval string `json:"gitignore_check_interval,omitempty"`
}

// DefaultSettings constructs Settings populated with standard defaults.
func DefaultSettings() Settings {
	return Settings{
		GitIgnoreCheckInterval: DefaultGitIgnoreCheckInterval,
	}
}

// GetGitIgnoreCheckInterval parses the configured interval with a 24-hour default fallback.
func (s *Settings) GetGitIgnoreCheckInterval() time.Duration {
	if s == nil || strings.TrimSpace(s.GitIgnoreCheckInterval) == "" {
		return 24 * time.Hour
	}
	dur, err := ParseIgnoreInterval(s.GitIgnoreCheckInterval)
	if err != nil {
		return 24 * time.Hour
	}
	return dur
}

// ParseIgnoreInterval parses human-readable duration strings into time.Duration.
func ParseIgnoreInterval(val string) (time.Duration, error) {
	trimmed := strings.ToLower(strings.TrimSpace(val))
	if isBypassInterval(trimmed) {
		return 0, nil
	}
	if strings.HasSuffix(trimmed, "d") {
		return parseDayInterval(val, trimmed)
	}
	return parseStandardDuration(val, trimmed)
}

func isBypassInterval(val string) bool {
	return val == "" || val == "0" || val == "off" || val == "never" || val == "disabled"
}

func parseDayInterval(raw, val string) (time.Duration, error) {
	daysStr := strings.TrimSuffix(val, "d")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 0 {
		return 0, fmt.Errorf("invalid day interval %q: must be positive integer followed by 'd'", raw)
	}
	return time.Duration(days) * 24 * time.Hour, nil
}

func parseStandardDuration(raw, val string) (time.Duration, error) {
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid interval %q: expected duration like 24h, 1d, 12h, 30m, or 0/off", raw)
	}
	if d < 0 {
		return 0, fmt.Errorf("interval cannot be negative: %s", raw)
	}
	return d, nil
}

// GetGitIgnoreCheckInterval retrieves the raw interval setting from store.DB or returns default.
func GetGitIgnoreCheckInterval(db *store.DB) string {
	if db == nil {
		return DefaultGitIgnoreCheckInterval
	}
	val := db.GetSetting(SettingGitIgnoreCheckInterval)
	if strings.TrimSpace(val) == "" {
		return DefaultGitIgnoreCheckInterval
	}
	return strings.TrimSpace(val)
}

// GetGitIgnoreTTL calculates the effective time.Duration from store.DB or defaults to 24h.
func GetGitIgnoreTTL(db *store.DB) time.Duration {
	intervalStr := GetGitIgnoreCheckInterval(db)
	ttl, err := ParseIgnoreInterval(intervalStr)
	if err != nil {
		return 24 * time.Hour
	}
	return ttl
}

// SetGitIgnoreCheckInterval validates and updates the interval in the settings table.
func SetGitIgnoreCheckInterval(db *store.DB, interval string) error {
	trimmed := strings.TrimSpace(interval)
	_, err := ParseIgnoreInterval(trimmed)
	if err != nil {
		return err
	}
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return db.SetSetting(SettingGitIgnoreCheckInterval, trimmed)
}
