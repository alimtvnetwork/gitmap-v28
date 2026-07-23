package cmd

// JSON encoders for `gitmap list-releases [--json] [--all-repos]`.
//
// Two surfaces are pinned to gitmap/stablejson so key ordering becomes
// a compile-time decision and stops being a reflection accident:
//
//   - encodeListReleasesJSON         → default per-repo view, fields
//     mirror model.ReleaseRecord's `json:` tags exactly so the wire
//     format stays byte-compatible with the legacy json.MarshalIndent
//     output. Schema: spec/08-json-schemas/list-releases.schema.json.
//
//   - encodeListReleasesAllReposJSON → joined --all-repos view. The
//     underlying store.ReleaseAcrossRepos struct historically had NO
//     json tags, so encoding/json emitted PascalCase field names
//     (ReleaseID, RepoSlug, …). We preserve that surface verbatim to
//     avoid silently breaking downstream scripts. Schema:
//     spec/08-json-schemas/list-releases-all-repos.schema.json.
//
// Both encoders take an io.Writer so the contract test in
// listreleases_jsonschema_contract_test.go can capture bytes into a
// buffer; CLI dispatch (printReleasesJSON / printAllReposJSON) passes
// os.Stdout. Empty input encodes as `[]\n` for both — preserves the
// `jq length` ergonomics every other gitmap stable JSON surface
// already guarantees.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// list-releases per-repo wire keys. Names + order mirror the
// `json:` tags on model.ReleaseRecord — KEEP IN SYNC. Reordering or
// renaming here is a consumer-facing break; the contract test fails
// on drift in either direction.
const (
	listReleasesKeyID           = "id"
	listReleasesKeyRepoID       = "repoId"
	listReleasesKeyVersion      = "version"
	listReleasesKeyTag          = "tag"
	listReleasesKeyBranch       = "branch"
	listReleasesKeySourceBranch = "sourceBranch"
	listReleasesKeyCommitSha    = "commitSha"
	listReleasesKeyChangelog    = "changelog"
	listReleasesKeyNotes        = "notes"
	listReleasesKeyIsDraft      = "isDraft"
	listReleasesKeyIsPreRelease = "isPreRelease"
	listReleasesKeyIsLatest     = "isLatest"
	listReleasesKeySource       = "source"
	listReleasesKeyCreatedAt    = "createdAt"
)

// list-releases --all-repos wire keys. PascalCase is intentional:
// store.ReleaseAcrossRepos has no `json:` tags, so the legacy
// json.MarshalIndent output used Go field names verbatim. We pin
// that surface here to avoid a silent rename when migrating off
// MarshalIndent.
const (
	listReleasesAllKeyReleaseID    = "ReleaseID"
	listReleasesAllKeyRepoID       = "RepoID"
	listReleasesAllKeyRepoSlug     = "RepoSlug"
	listReleasesAllKeyRepoPath     = "RepoPath"
	listReleasesAllKeyVersion      = "Version"
	listReleasesAllKeyTag          = "Tag"
	listReleasesAllKeyBranch       = "Branch"
	listReleasesAllKeyCommitSha    = "CommitSha"
	listReleasesAllKeySource       = "Source"
	listReleasesAllKeyIsDraft      = "IsDraft"
	listReleasesAllKeyIsLatest     = "IsLatest"
	listReleasesAllKeyIsPreRelease = "IsPreRelease"
	listReleasesAllKeyCreatedAt    = "CreatedAt"
)

// encodeListReleasesJSON writes the per-repo release list as a
// stablejson 2-space-indented array. Empty input emits `[]\n`.
func encodeListReleasesJSON(w io.Writer, records []model.ReleaseRecord) error {
	return stablejson.WriteArray(w, buildListReleasesJSONItems(records))
}

// buildListReleasesJSONItems is the single source of (field name,
// field order, value) for the per-repo view. Centralized so a future
// column rename or reorder is one diff and the contract test catches
// schema drift in the same PR.
func buildListReleasesJSONItems(records []model.ReleaseRecord) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(records))
	for _, r := range records {
		items = append(items, []stablejson.Field{
			{Key: listReleasesKeyID, Value: r.ID},
			{Key: listReleasesKeyRepoID, Value: r.RepoID},
			{Key: listReleasesKeyVersion, Value: r.Version},
			{Key: listReleasesKeyTag, Value: r.Tag},
			{Key: listReleasesKeyBranch, Value: r.Branch},
			{Key: listReleasesKeySourceBranch, Value: r.SourceBranch},
			{Key: listReleasesKeyCommitSha, Value: r.CommitSha},
			{Key: listReleasesKeyChangelog, Value: r.Changelog},
			{Key: listReleasesKeyNotes, Value: r.Notes},
			{Key: listReleasesKeyIsDraft, Value: r.IsDraft},
			{Key: listReleasesKeyIsPreRelease, Value: r.IsPreRelease},
			{Key: listReleasesKeyIsLatest, Value: r.IsLatest},
			{Key: listReleasesKeySource, Value: r.Source},
			{Key: listReleasesKeyCreatedAt, Value: r.CreatedAt},
		})
	}

	return items
}

// encodeListReleasesAllReposJSON writes the joined --all-repos view
// as a stablejson 2-space-indented array. PascalCase keys preserve
// the legacy MarshalIndent surface byte-for-byte.
func encodeListReleasesAllReposJSON(w io.Writer, records []store.ReleaseAcrossRepos) error {
	return stablejson.WriteArray(w, buildListReleasesAllReposJSONItems(records))
}

// buildListReleasesAllReposJSONItems mirrors buildListReleasesJSONItems
// for the joined view. Field order follows the struct declaration in
// store/releaseacrossrepos.go to match what json.MarshalIndent emitted
// for years.
func buildListReleasesAllReposJSONItems(records []store.ReleaseAcrossRepos) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(records))
	for _, r := range records {
		items = append(items, []stablejson.Field{
			{Key: listReleasesAllKeyReleaseID, Value: r.ReleaseID},
			{Key: listReleasesAllKeyRepoID, Value: r.RepoID},
			{Key: listReleasesAllKeyRepoSlug, Value: r.RepoSlug},
			{Key: listReleasesAllKeyRepoPath, Value: r.RepoPath},
			{Key: listReleasesAllKeyVersion, Value: r.Version},
			{Key: listReleasesAllKeyTag, Value: r.Tag},
			{Key: listReleasesAllKeyBranch, Value: r.Branch},
			{Key: listReleasesAllKeyCommitSha, Value: r.CommitSha},
			{Key: listReleasesAllKeySource, Value: r.Source},
			{Key: listReleasesAllKeyIsDraft, Value: r.IsDraft},
			{Key: listReleasesAllKeyIsLatest, Value: r.IsLatest},
			{Key: listReleasesAllKeyIsPreRelease, Value: r.IsPreRelease},
			{Key: listReleasesAllKeyCreatedAt, Value: r.CreatedAt},
		})
	}

	return items
}
