package cmd

// JSON encoder for `gitmap diff-profiles --json`.
//
// Migrated off json.MarshalIndent(map[string]any) onto gitmap/stablejson
// so key order becomes a compile-time decision rather than a reflection
// accident. Nested arrays (onlyInA, onlyInB, different) are pre-rendered
// in compact mode and embedded as json.RawMessage so their key order
// is also stable.
//
// Schema: spec/08-json-schemas/diff-profiles.schema.json.

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/stablejson"
)

// diff-profiles top-level wire keys.
const (
	dpKeyProfileA  = "profileA"
	dpKeyProfileB  = "profileB"
	dpKeyOnlyInA   = "onlyInA"
	dpKeyOnlyInB   = "onlyInB"
	dpKeyDifferent = "different"
	dpKeySame      = "same"
)

// diff-profiles repo summary wire keys (onlyInA / onlyInB items).
const (
	dpSummaryKeyName = "name"
	dpSummaryKeyPath = "path"
)

// diff-profiles difference wire keys.
const (
	dpDiffKeyName  = "name"
	dpDiffKeyPathA = "pathA"
	dpDiffKeyPathB = "pathB"
	dpDiffKeyModeA = "modeA"
	dpDiffKeyModeB = "modeB"
)

// encodeDiffProfilesJSON writes a single diff-profiles result as
// stable JSON with 2-space indentation.
func encodeDiffProfilesJSON(w io.Writer, profileA, profileB string, result dpResult) error {
	onlyInARaw, err := renderDPRepoSummariesRaw(result.onlyInA)
	if err != nil {
		return err
	}

	onlyInBRaw, err := renderDPRepoSummariesRaw(result.onlyInB)
	if err != nil {
		return err
	}

	differentRaw, err := renderDPDiffsRaw(result.different)
	if err != nil {
		return err
	}

	return stablejson.WriteObject(w, []stablejson.Field{
		{Key: dpKeyProfileA, Value: profileA},
		{Key: dpKeyProfileB, Value: profileB},
		{Key: dpKeyOnlyInA, Value: onlyInARaw},
		{Key: dpKeyOnlyInB, Value: onlyInBRaw},
		{Key: dpKeyDifferent, Value: differentRaw},
		{Key: dpKeySame, Value: len(result.same)},
	})
}

// renderDPRepoSummariesRaw pre-renders a repo summary array in
// compact mode for embedding.
func renderDPRepoSummariesRaw(records []model.ScanRecord) (json.RawMessage, error) {
	items := make([][]stablejson.Field, 0, len(records))
	for _, r := range records {
		items = append(items, []stablejson.Field{
			{Key: dpSummaryKeyName, Value: r.RepoName},
			{Key: dpSummaryKeyPath, Value: r.AbsolutePath},
		})
	}

	var buf bytes.Buffer
	if err := stablejson.WriteArrayIndent(&buf, items, ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}

// renderDPDiffsRaw pre-renders the differences array in compact
// mode for embedding.
func renderDPDiffsRaw(diffs []dpDiff) (json.RawMessage, error) {
	items := make([][]stablejson.Field, 0, len(diffs))
	for _, d := range diffs {
		items = append(items, []stablejson.Field{
			{Key: dpDiffKeyName, Value: d.Name},
			{Key: dpDiffKeyPathA, Value: d.PathA},
			{Key: dpDiffKeyPathB, Value: d.PathB},
			{Key: dpDiffKeyModeA, Value: d.ModeA},
			{Key: dpDiffKeyModeB, Value: d.ModeB},
		})
	}

	var buf bytes.Buffer
	if err := stablejson.WriteArrayIndent(&buf, items, ""); err != nil {
		return nil, err
	}

	return json.RawMessage(bytes.TrimSuffix(buf.Bytes(), []byte{'\n'})), nil
}
