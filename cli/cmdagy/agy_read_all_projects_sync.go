package cmdagy

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

type localRepo struct {
	Name string
	Path string
}

func discoverLocalRepos() []localRepo {
	scanFile := filepath.Join(".gitmap", "output", "gitmap.json")
	if records, err := model.LoadStatusRecords(scanFile); err == nil && len(records) > 0 {
		return reposFromScanRecords(records)
	}

	return discoverCwdRepos()
}

func reposFromScanRecords(records []model.ScanRecord) []localRepo {
	var repos []localRepo
	for _, r := range records {
		if r.AbsolutePath != "" {
			repos = append(repos, localRepo{Name: r.RepoName, Path: r.AbsolutePath})
		}
	}

	return repos
}

func discoverCwdRepos() []localRepo {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	entries, err := os.ReadDir(cwd)
	if err != nil {
		return nil
	}

	return filterGitDirEntries(cwd, entries)
}

func filterGitDirEntries(cwd string, entries []os.DirEntry) []localRepo {
	var repos []localRepo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		sub := filepath.Join(cwd, e.Name())
		if _, statErr := os.Stat(filepath.Join(sub, ".git")); statErr == nil {
			repos = append(repos, localRepo{Name: e.Name(), Path: sub})
		}
	}

	return repos
}
