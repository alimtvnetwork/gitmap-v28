package cmd

// JSON encoder for `gitmap find-next --json`.
//
// Migrated off json.Encoder onto gitmap/stablejson so the top-level
// key order (repo, nextVersionTag, nextVersionNum, method, probedAt)
// becomes a compile-time decision rather than a reflection accident.
// Schema: spec/08-json-schemas/find-next.schema.json.
//
// The nested `repo` value rides on encoding/json via stablejson's
// per-value Marshal — model.ScanRecord's wire shape is pinned by the
// existing tag-based golden in findnextjson_contract_test.go, and
// per stable-json-encoding constraints we don't hand-enumerate
// wide (>12-field) structs here.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// find-next top-level wire keys. Names + order are the contract;
// reordering or renaming here is a consumer-facing break and the
// schema contract test fails on drift in either direction.
const (
	findNextKeyRepo           = "repo"
	findNextKeyNextVersionTag = "nextVersionTag"
	findNextKeyNextVersionNum = "nextVersionNum"
	findNextKeyMethod         = "method"
	findNextKeyProbedAt       = "probedAt"
)

// encodeFindNextJSON writes rows as a stablejson 2-space-indented
// array. Empty input emits `[]\n` so `jq length` works without a
// special case. Split out from CLI dispatch so contract tests can
// capture the bytes into a buffer instead of stdout.
func encodeFindNextJSON(w io.Writer, rows []model.FindNextRow) error {
	return stablejson.WriteArray(w, buildFindNextJSONItems(rows))
}

// buildFindNextJSONItems is the single source of (field name, field
// order, value) for find-next. Centralized so a future column
// rename/reorder is one diff and the contract test catches schema
// drift in the same PR.
func buildFindNextJSONItems(rows []model.FindNextRow) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(rows))
	for _, r := range rows {
		items = append(items, []stablejson.Field{
			{Key: findNextKeyRepo, Value: r.Repo},
			{Key: findNextKeyNextVersionTag, Value: r.NextVersionTag},
			{Key: findNextKeyNextVersionNum, Value: r.NextVersionNum},
			{Key: findNextKeyMethod, Value: r.Method},
			{Key: findNextKeyProbedAt, Value: r.ProbedAt},
		})
	}

	return items
}
