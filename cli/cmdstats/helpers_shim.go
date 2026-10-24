package cmdstats

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstats/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

func isLegacyDataError(err error) bool {
	hasScanErr := strings.Contains(err.Error(), constants.ErrPatternScanError)
	hasDriverValErr := strings.Contains(err.Error(), constants.ErrPatternDriverValue)

	return hasScanErr || hasDriverValErr
}
