package store

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ExportAll gathers the entire database into a single DatabaseExport struct.
func (db *DB) ExportAll() (model.DatabaseExport, error) {
	export := model.DatabaseExport{
		Version:    constants.Version,
		ExportedAt: time.Now().Format(time.RFC3339),
	}

	var err error

	export.Repos, err = db.ListRepos()
	if err != nil {
		return export, err
	}

	export.Groups, err = db.exportGroups()
	if err != nil {
		return export, err
	}

	export.Releases, err = db.ListReleases()
	if err != nil {
		return export, err
	}

	export.History, err = db.ListHistory()
	if err != nil {
		return export, err
	}

	export.Bookmarks, err = db.ListBookmarks()
	if err != nil {
		return export, err
	}

	return export, nil
}

// exportGroups loads all groups with their member repo slugs.
func (db *DB) exportGroups() ([]model.GroupExport, error) {
	groups, err := db.ListGroups()
	if err != nil {
		return nil, err
	}

	var results []model.GroupExport

	for _, g := range groups {
		ge, err := db.buildGroupExport(g)
		if err != nil {
			return nil, err
		}
		results = append(results, ge)
	}

	return results, nil
}

// buildGroupExport creates a GroupExport with repo slugs for one group.
func (db *DB) buildGroupExport(g model.Group) (model.GroupExport, error) {
	repos, err := db.ShowGroup(g.Name)
	if err != nil {
		return model.GroupExport{}, err
	}

	slugs := make([]string, len(repos))
	for i, r := range repos {
		slugs[i] = r.Slug
	}

	return model.GroupExport{Group: g, RepoSlugs: slugs}, nil
}
