package cmd

// Schema contract for `gitmap list-releases [--all-repos] --json`.
// Pairs the runtime encoders (encodeListReleasesJSON +
// encodeListReleasesAllReposJSON) with the published schemas at
// spec/08-json-schemas/list-releases.schema.json and
// spec/08-json-schemas/list-releases-all-repos.schema.json so a drift
// in either side fails the build. Mirrors the structure of
// startuplist_jsonschema_contract_test.go — same generic helpers
// (findSchemaFile, loadSchemaFile, propertyOrder extraction) live in
// jsonschema_helpers_test.go.

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

const (
	listReleasesSchemaFilename         = "list-releases.schema.json"
	listReleasesAllReposSchemaFilename = "list-releases-all-repos.schema.json"
)

// listReleasesItemSchema descends into the top-level array's `items`
// subschema where the per-record object contract lives. Centralized
// so each assertion stays flat.
func listReleasesItemSchema(t *testing.T, root map[string]any) map[string]any {
	t.Helper()
	items, ok := root["items"].(map[string]any)
	if !ok {
		t.Fatalf("schema has no items object")
	}

	return items
}

// TestListReleasesSchema_TopLevelIsArray pins the most fundamental
// shape decision: empty output is `[]`, NOT `null`. A future encoder
// regression that emitted `null` for empty would break every
// downstream `jq length` consumer.
func TestListReleasesSchema_TopLevelIsArray(t *testing.T) {
	for _, name := range []string{listReleasesSchemaFilename, listReleasesAllReposSchemaFilename} {
		root := loadSchemaFile(t, name)
		if root["type"] != "array" {
			t.Fatalf("%s top-level type = %v, want array", name, root["type"])
		}
	}
}

// TestListReleasesSchema_RequiredKeysMatchEncoder asserts the
// per-repo schema's required-key set equals the encoder's emitted
// keys. A new column added to model.ReleaseRecord without a schema
// update fails here with a clean diff.
func TestListReleasesSchema_RequiredKeysMatchEncoder(t *testing.T) {
	root := loadSchemaFile(t, listReleasesSchemaFilename)
	required := stringSliceFromAny(listReleasesItemSchema(t, root)["required"])
	want := []string{
		"id", "repoId", "version", "tag", "branch", "sourceBranch",
		"commitSha", "changelog", "notes", "isDraft", "isPreRelease",
		"isLatest", "source", "createdAt",
	}
	if !equalStringSlices(required, want) {
		t.Fatalf("schema required = %v, want %v", required, want)
	}
}

// TestListReleasesSchema_PropertyOrderMatchesEncoder is the headline
// contract test: encode a real record, parse the resulting JSON
// preserving key order, and assert the order matches the schema's
// propertyOrder array. The ONLY guard that catches a reordering of
// the stablejson.Field slice in listreleasesrender.go.
func TestListReleasesSchema_PropertyOrderMatchesEncoder(t *testing.T) {
	root := loadSchemaFile(t, listReleasesSchemaFilename)
	want := stringSliceFromAny(listReleasesItemSchema(t, root)["propertyOrder"])
	if len(want) == 0 {
		t.Fatalf("schema item has no propertyOrder array")
	}
	records := []model.ReleaseRecord{{Version: "1.0.0", Tag: "v1.0.0"}}
	var buf bytes.Buffer
	if err := encodeListReleasesJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	got := extractFirstObjectKeyOrder(t, buf.Bytes())
	if !equalStringSlices(got, want) {
		t.Fatalf("emitted key order = %v, schema propertyOrder = %v", got, want)
	}
}

// TestListReleasesSchema_EmptyEncodesAsArray pins empty-input
// behavior end-to-end. Belt-and-suspenders against a future encoder
// regression that emits `null` for empty.
func TestListReleasesSchema_EmptyEncodesAsArray(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeListReleasesJSON(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytes.Equal(bytes.TrimSpace(buf.Bytes()), []byte("[]")) {
		t.Fatalf("empty encoded as %q, want []", buf.Bytes())
	}
}

// TestListReleasesAllReposSchema_RequiredKeysMatchEncoder mirrors the
// per-repo required-key check for the joined --all-repos view. The
// PascalCase keys are intentional — see schema description.
func TestListReleasesAllReposSchema_RequiredKeysMatchEncoder(t *testing.T) {
	root := loadSchemaFile(t, listReleasesAllReposSchemaFilename)
	required := stringSliceFromAny(listReleasesItemSchema(t, root)["required"])
	want := []string{
		"ReleaseID", "RepoID", "RepoSlug", "RepoPath", "Version", "Tag",
		"Branch", "CommitSha", "Source", "IsDraft", "IsLatest",
		"IsPreRelease", "CreatedAt",
	}
	if !equalStringSlices(required, want) {
		t.Fatalf("schema required = %v, want %v", required, want)
	}
}

// TestListReleasesAllReposSchema_PropertyOrderMatchesEncoder pins
// PascalCase key order for the joined view. Drift here means a
// silent rename of the legacy MarshalIndent surface — must NOT
// happen without a deliberate consumer-facing changelog entry.
func TestListReleasesAllReposSchema_PropertyOrderMatchesEncoder(t *testing.T) {
	root := loadSchemaFile(t, listReleasesAllReposSchemaFilename)
	want := stringSliceFromAny(listReleasesItemSchema(t, root)["propertyOrder"])
	if len(want) == 0 {
		t.Fatalf("schema item has no propertyOrder array")
	}
	records := []store.ReleaseAcrossRepos{{ReleaseID: 1, RepoSlug: "a/b", Version: "1.0.0"}}
	var buf bytes.Buffer
	if err := encodeListReleasesAllReposJSON(&buf, records); err != nil {
		t.Fatalf("encode: %v", err)
	}
	got := extractFirstObjectKeyOrder(t, buf.Bytes())
	if !equalStringSlices(got, want) {
		t.Fatalf("emitted key order = %v, schema propertyOrder = %v", got, want)
	}
}

// TestListReleasesAllReposSchema_EmptyEncodesAsArray pins the
// joined view's empty-input shape.
func TestListReleasesAllReposSchema_EmptyEncodesAsArray(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeListReleasesAllReposJSON(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytes.Equal(bytes.TrimSpace(buf.Bytes()), []byte("[]")) {
		t.Fatalf("empty encoded as %q, want []", buf.Bytes())
	}
}

// TestListReleasesJSON_ByteCompatWithLegacyMarshalIndent guards the
// migration's headline promise: the new stablejson encoder produces
// bytes IDENTICAL to the legacy json.MarshalIndent output the
// command emitted before this PR. Any drift here means a downstream
// consumer doing byte-level diffs (CI golden files, content-hash
// caches) would see a spurious change. The legacy `\n` suffix
// matches the old `fmt.Println(string(data))` framing.
func TestListReleasesJSON_ByteCompatWithLegacyMarshalIndent(t *testing.T) {
	rec := []model.ReleaseRecord{{
		ID: 7, RepoID: 3, Version: "1.0.0", Tag: "v1.0.0",
		Branch: "release-v1.0.0", SourceBranch: "main", CommitSha: "abc",
		Changelog: "first", Notes: "n", IsLatest: true, Source: "repo",
		CreatedAt: "2026-05-05T00:00:00Z",
	}}
	legacy, err := jsonMarshalIndentForTest(rec)
	if err != nil {
		t.Fatalf("legacy marshal: %v", err)
	}
	var got bytes.Buffer
	if err := encodeListReleasesJSON(&got, rec); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytes.Equal(append(legacy, '\n'), got.Bytes()) {
		t.Fatalf("byte drift from legacy MarshalIndent:\n--- legacy ---\n%s\n--- new ---\n%s",
			legacy, got.Bytes())
	}
}

// TestListReleasesAllReposJSON_ByteCompatWithLegacyMarshalIndent
// pins the same byte-equality guarantee for the joined --all-repos
// view. PascalCase keys are part of the contract here.
func TestListReleasesAllReposJSON_ByteCompatWithLegacyMarshalIndent(t *testing.T) {
	rec := []store.ReleaseAcrossRepos{{
		ReleaseID: 1, RepoID: 2, RepoSlug: "a/b", RepoPath: "/x",
		Version: "1.0.0", Tag: "v1.0.0", Branch: "release-v1.0.0",
		CommitSha: "abc", Source: "repo", IsLatest: true,
		CreatedAt: "2026-05-05",
	}}
	legacy, err := jsonMarshalIndentForTest(rec)
	if err != nil {
		t.Fatalf("legacy marshal: %v", err)
	}
	var got bytes.Buffer
	if err := encodeListReleasesAllReposJSON(&got, rec); err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !bytes.Equal(append(legacy, '\n'), got.Bytes()) {
		t.Fatalf("byte drift from legacy MarshalIndent:\n--- legacy ---\n%s\n--- new ---\n%s",
			legacy, got.Bytes())
	}
}

// jsonMarshalIndentForTest reproduces the EXACT call the legacy
// printReleasesJSON used (`json.MarshalIndent(v, "", "  ")`) so the
// byte-compat tests above pin the new encoder against that surface.
// Kept inside the contract test file so the encoding/json import
// only lives in test code — production listreleases.go is fully on
// stablejson.
func jsonMarshalIndentForTest(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
