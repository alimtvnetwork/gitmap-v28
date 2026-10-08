package cmdstatus

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// GetAliasPathFn is wired by cmd/di_hooks.go to the canonical implementation.
var GetAliasPathFn func() string

func GetAliasPath() string {
	if GetAliasPathFn != nil {
		return GetAliasPathFn()
	}
	return ""
}

// GetAliasSlugFn is wired by cmd/di_hooks.go to the canonical implementation.
var GetAliasSlugFn func() string

func GetAliasSlug() string {
	if GetAliasSlugFn != nil {
		return GetAliasSlugFn()
	}
	return ""
}

// HasAliasFn is wired by cmd/di_hooks.go to the canonical implementation.
var HasAliasFn func() bool

func HasAlias() bool {
	if HasAliasFn != nil {
		return HasAliasFn()
	}
	return false
}

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// LoadAllRecordsDB loads all repos from the database.
// Exported so cmd/di_hooks.go can wire it into other packages' LoadAllRecordsDBFn hooks.
func LoadAllRecordsDB() []model.ScanRecord {
	return loadAllRecordsDB()
}
