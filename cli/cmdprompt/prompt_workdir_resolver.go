// Package cmdprompt — prompt_workdir_resolver.go extracts targets from registered work directories.
package cmdprompt

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func ResolveAllWorkDirPromptTargets() PromptTargetSliceResult {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return result.FailSlice[string](apperror.WrapSimple(errDB, "open store for workdir targets"))
	}

	defer db.Close()

	dirs, errList := db.ListWorkDirs()
	if errList != nil {
		return result.FailSlice[string](apperror.WrapSimple(errList, "list workdirs for prompt targets"))
	}

	var allTargets []string
	for _, d := range dirs {
		childRes := DiscoverPromptChildRepos(d.AbsolutePath)
		if childRes.HasRecord() {
			allTargets = append(allTargets, childRes.Data...)
		} else {
			allTargets = append(allTargets, d.AbsolutePath)
		}
	}

	return result.OkSlice(allTargets)
}
