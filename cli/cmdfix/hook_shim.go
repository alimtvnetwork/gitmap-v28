package cmdfix

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// CheckHelpFn is wired by cmd/di_hooks.go.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// HasAliasFn is wired by cmd/di_hooks.go.
var HasAliasFn func() bool

func HasAlias() bool {
	if HasAliasFn != nil {
		return HasAliasFn()
	}
	return false
}

// GetAliasSlugFn is wired by cmd/di_hooks.go.
var GetAliasSlugFn func() string

func GetAliasSlug() string {
	if GetAliasSlugFn != nil {
		return GetAliasSlugFn()
	}
	return ""
}

// GetAliasPathFn is wired by cmd/di_hooks.go.
var GetAliasPathFn func() string

func GetAliasPath() string {
	if GetAliasPathFn != nil {
		return GetAliasPathFn()
	}
	return ""
}

// LoadAllRecordsDBFn is wired by cmd/di_hooks.go.
var LoadAllRecordsDBFn func() []model.ScanRecord

func loadAllRecordsDB() []model.ScanRecord {
	if LoadAllRecordsDBFn != nil {
		return LoadAllRecordsDBFn()
	}
	return nil
}

// LoadRecordsJSONFallbackFn is wired by cmd/di_hooks.go.
var LoadRecordsJSONFallbackFn func() []model.ScanRecord

func loadRecordsJSONFallback() []model.ScanRecord {
	if LoadRecordsJSONFallbackFn != nil {
		return LoadRecordsJSONFallbackFn()
	}
	return nil
}
