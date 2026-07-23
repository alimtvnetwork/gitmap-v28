// Package cmd — projectreposrender.go is the stablejson encoder for
// `gitmap <type>-repos --json`. Migrated off encoding/json so wire-key
// order is a compile-time decision pinned by the schema, not a
// reflection accident on model.DetectedProject.
//
// Schema: spec/08-json-schemas/project-repos.schema.json.
package cmd

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// project-repos wire keys. Names + order are the contract.
const (
	projectReposKeyID               = "id"
	projectReposKeyRepoID           = "repoId"
	projectReposKeyRepoName         = "repoName"
	projectReposKeyProjectTypeID    = "projectTypeId"
	projectReposKeyProjectType      = "projectType"
	projectReposKeyProjectName      = "projectName"
	projectReposKeyAbsolutePath     = "absolutePath"
	projectReposKeyRepoPath         = "repoPath"
	projectReposKeyRelativePath     = "relativePath"
	projectReposKeyPrimaryIndicator = "primaryIndicator"
	projectReposKeyDetectedAt       = "detectedAt"
)

// encodeProjectReposJSON writes projects as a stablejson 2-space-indented
// array. Empty input emits `[]\n`.
func encodeProjectReposJSON(w io.Writer, projects []model.DetectedProject) error {
	return stablejson.WriteArray(w, buildProjectReposJSONItems(projects))
}

// buildProjectReposJSONItems is the single source of (field name,
// field order, value) for project-repos.
func buildProjectReposJSONItems(projects []model.DetectedProject) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(projects))
	for _, p := range projects {
		items = append(items, []stablejson.Field{
			{Key: projectReposKeyID, Value: p.ID},
			{Key: projectReposKeyRepoID, Value: p.RepoID},
			{Key: projectReposKeyRepoName, Value: p.RepoName},
			{Key: projectReposKeyProjectTypeID, Value: p.ProjectTypeID},
			{Key: projectReposKeyProjectType, Value: p.ProjectType},
			{Key: projectReposKeyProjectName, Value: p.ProjectName},
			{Key: projectReposKeyAbsolutePath, Value: p.AbsolutePath},
			{Key: projectReposKeyRepoPath, Value: p.RepoPath},
			{Key: projectReposKeyRelativePath, Value: p.RelativePath},
			{Key: projectReposKeyPrimaryIndicator, Value: p.PrimaryIndicator},
			{Key: projectReposKeyDetectedAt, Value: p.DetectedAt},
		})
	}

	return items
}
