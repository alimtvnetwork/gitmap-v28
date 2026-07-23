// Package gitutil — centralized date display formatting.
package gitutil

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// FormatDisplayDate converts a time.Time to the local time zone
// and returns a human-friendly string: DD-Mon-YYYY hh:mm AM/PM.
func FormatDisplayDate(t time.Time) string {
	utc := t.UTC()
	local := utc.Local()

	return local.Format(constants.DateDisplayLayout)
}

// FormatDisplayDateUTC converts a time.Time to UTC
// and returns a human-friendly string: DD-Mon-YYYY hh:mm AM/PM (UTC).
func FormatDisplayDateUTC(t time.Time) string {
	utc := t.UTC()

	return utc.Format(constants.DateDisplayLayout) + constants.DateUTCSuffix
}
