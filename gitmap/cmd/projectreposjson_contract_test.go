package cmd

// JSON contract tests for `gitmap <type>-repos --json`.
//
// project-repos emits an array of model.DetectedProject. The contract
// covers:
//
//   - Top-level array shape (empty must be `[]\n`).
//   - Key order: id, repoId, repoName, projectTypeId, projectType,
//     projectName, absolutePath, repoPath, relativePath,
//     primaryIndicator, detectedAt.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run ProjectReposJSONContract

import (
	"bytes"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// TestProjectReposJSONContract_EmptyIsArrayNotNull is the jq-compat guarantee.
func TestProjectReposJSONContract_EmptyIsArrayNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "project_repos_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeProjectReposJSON(&buf, nil)

		return buf.Bytes(), err
	})
}

// canonicalDetectedProject builds a deterministic single row.
func canonicalDetectedProject() model.DetectedProject {
	return model.DetectedProject{
		ID:               42,
		RepoID:           7,
		RepoName:         "gitmap-v27",
		ProjectTypeID:    1,
		ProjectType:      "go",
		ProjectName:      "gitmap",
		AbsolutePath:     "/home/user/code/gitmap-v27/gitmap",
		RepoPath:         "/home/user/code/gitmap-v27",
		RelativePath:     "gitmap",
		PrimaryIndicator: "go.mod",
		DetectedAt:       "2025-01-01T12:00:00Z",
	}
}

// TestProjectReposJSONContract_CanonicalRow_KeyOrder asserts the key
// order of the emitted object matches the schema declaration.
func TestProjectReposJSONContract_CanonicalRow_KeyOrder(t *testing.T) {
	projects := []model.DetectedProject{canonicalDetectedProject()}
	var buf bytes.Buffer
	if err := encodeProjectReposJSON(&buf, projects); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "project-repos")
}
