package cmddb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// PopulateRepoAliasesWithCountFn is wired by cmd/di_hooks.go to the canonical implementation.
var PopulateRepoAliasesWithCountFn func(db *store.DB) (int, *apperror.AppError)

// PopulateRepoAliasesWithCount populates aliases and returns the count of added aliases.
func PopulateRepoAliasesWithCount(db *store.DB) (int, *apperror.AppError) {
	if PopulateRepoAliasesWithCountFn != nil {
		return PopulateRepoAliasesWithCountFn(db)
	}

	return 0, nil
}
