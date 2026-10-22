package cmdcpar

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CollectDirtyRepositories scans all tracked repositories and returns those with dirty states.
func CollectDirtyRepositories() []DirtyRepoSummary {
	records := resolveAllRecords()
	var dirty []DirtyRepoSummary
	for _, r := range records {
		relPath := resolveRelativePath(r)
		diag := gitutil.InspectDirtyState(relPath)
		if diag.IsDirty {
			dirty = appendDirty(dirty, r, relPath, diag)
		}
	}

	return dirty
}

func resolveRelativePath(r model.ScanRecord) string {
	hasRel := len(r.RelativePath) > 0
	if hasRel {
		return r.RelativePath
	}

	hasName := len(r.RepoName) > 0
	if hasName {
		return r.RepoName
	}

	return "."
}

func appendDirty(dirty []DirtyRepoSummary, r model.ScanRecord, relPath string, diag gitutil.DirtyDiagnosis) []DirtyRepoSummary {
	return append(dirty, DirtyRepoSummary{
		RepoName:  r.RepoName,
		RepoPath:  relPath,
		Diagnosis: diag,
	})
}

func resolveAllRecords() []model.ScanRecord {
	records := queryStoreRecords()
	hasRecords := len(records) > 0
	if hasRecords {
		return records
	}

	return []model.ScanRecord{{RelativePath: ".", RepoName: "."}}
}

func queryStoreRecords() []model.ScanRecord {
	s, err := store.OpenDefault()
	hasErr := err != nil
	if hasErr {
		return nil
	}

	defer s.Close()

	records, _ := s.ListRepos()

	return records
}
