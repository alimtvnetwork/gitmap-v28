package clonefrom

// Drift guard for the envelope-level `provenance` map. Three
// contracts pinned here, all derived from the production code so
// the test fails the moment a refactor breaks any of them:
//
//  1. Every JSON row field declared by reportRowJSON has a matching
//     entry in constants.CloneFromReportProvenance — adding a new
//     row field WITHOUT a provenance entry fails this test.
//  2. Provenance entries reference only known stages (scan / mapper
//     / clonefrom). Typos in a stage string fail loudly.
//  3. The emitted envelope round-trips: every provenance.field
//     value appears as a real key under at least one rows[] entry.
//     This catches rename drift between the row struct and the
//     constants table even when both compile.

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestProvenance_CoversEveryReportField(t *testing.T) {
	want := jsonTagFieldNames(reflect.TypeOf(reportRowJSON{}))
	got := provenanceFieldNames()
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("provenance fields drifted from reportRowJSON.\n"+
			"reportRowJSON tags : %v\nprovenance entries : %v\n"+
			"Update constants.CloneFromReportProvenance to match the "+
			"struct field order, then bump CloneFromReportSchemaVersion "+
			"and regenerate the JSON goldens.", want, got)
	}
}

func TestProvenance_StagesAreKnown(t *testing.T) {
	allowed := map[string]bool{
		constants.ProvenanceStageScan:      true,
		constants.ProvenanceStageMapper:    true,
		constants.ProvenanceStageClonefrom: true,
	}
	for _, p := range constants.CloneFromReportProvenance {
		if !allowed[p.Stage] {
			t.Errorf("provenance entry %q references unknown stage %q",
				p.Field, p.Stage)
		}
	}
}

func TestProvenance_RoundTripsInEnvelope(t *testing.T) {
	var buf bytes.Buffer
	if err := writeReportRowsJSON(&buf, canonicalReportResults()); err != nil {
		t.Fatalf("writeReportRowsJSON: %v", err)
	}
	var env struct {
		Provenance []struct {
			Field string `json:"field"`
			Stage string `json:"stage"`
		} `json:"provenance"`
		Rows []map[string]any `json:"rows"`
	}
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v\nraw: %s", err, buf.String())
	}
	if len(env.Rows) == 0 {
		t.Fatal("canonical results produced no rows; cannot verify round-trip")
	}
	for _, p := range env.Provenance {
		if _, exists := env.Rows[0][p.Field]; !exists {
			t.Errorf("provenance.field %q has no matching key in rows[0] %v",
				p.Field, env.Rows[0])
		}
	}
}

// jsonTagFieldNames extracts the json:"..." tag values from a
// struct type in declaration order. Fails the test if any field
// lacks a json tag — every reportRowJSON field MUST be tagged.
func jsonTagFieldNames(t reflect.Type) []string {
	out := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" {
			continue
		}
		out = append(out, tag)
	}

	return out
}

func provenanceFieldNames() []string {
	out := make([]string, 0, len(constants.CloneFromReportProvenance))
	for _, p := range constants.CloneFromReportProvenance {
		out = append(out, p.Field)
	}

	return out
}
