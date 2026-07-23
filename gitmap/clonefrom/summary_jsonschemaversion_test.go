package clonefrom

// summary_jsonschemaversion_test.go — independently pins the
// `schemaVersion` field of the clone-from JSON report. This test
// is INTENTIONALLY redundant with the canonical/empty golden tests
// in summary_jsongolden_test.go: those tests catch field-set drift
// but a careless contributor could regenerate them and silently
// bump the version. This test asserts the version literal is 1
// and the field exists in BOTH the empty and populated cases, so
// any deliberate bump must be made here AND justified in the PR
// (and CHANGELOG.md, since it's a downstream-visible break).
//
// The version is sourced from constants.CloneFromReportSchemaVersion
// and the literal expectation lives here — they must stay in sync.
// When you bump the constant, this test fails until you update the
// expected literal below, forcing a conscious decision.
//
// See also: spec/05-coding-guidelines/ on stable on-disk schemas.

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// expectedSchemaVersionPinned is the version every clone-from JSON
// report must declare. Bumping requires: (1) editing this literal,
// (2) editing constants.CloneFromReportSchemaVersion to match, (3)
// regenerating both JSON goldens, (4) noting the break in CHANGELOG.md.
const expectedSchemaVersionPinned = 3

// envelopePeek is a minimal decoder used only by this test — we
// deliberately do NOT reuse reportEnvelopeJSON so a rename of the
// production type's JSON tag still trips this test.
type envelopePeek struct {
	SchemaVersion *int              `json:"schemaVersion"`
	Rows          []json.RawMessage `json:"rows"`
}

// TestCloneFromReportJSON_SchemaVersion_ConstantPinned guards the
// constant value itself. A drift here means someone bumped the
// schema without updating downstream consumers / CHANGELOG.
func TestCloneFromReportJSON_SchemaVersion_ConstantPinned(t *testing.T) {
	if constants.CloneFromReportSchemaVersion != expectedSchemaVersionPinned {
		t.Fatalf("CloneFromReportSchemaVersion drifted: got %d, "+
			"expected %d. If this bump is intentional, update "+
			"expectedSchemaVersionPinned in this file, regenerate "+
			"the JSON goldens, and document the break in CHANGELOG.md.",
			constants.CloneFromReportSchemaVersion,
			expectedSchemaVersionPinned)
	}
}

// TestCloneFromReportJSON_SchemaVersion_EmittedEmpty verifies the
// version is present in the empty-results envelope (regression
// guard: a refactor that special-cased the empty path could omit it).
func TestCloneFromReportJSON_SchemaVersion_EmittedEmpty(t *testing.T) {
	assertEnvelopeSchemaVersion(t, nil, 0)
}

// TestCloneFromReportJSON_SchemaVersion_EmittedPopulated verifies
// the version is present in the canonical 3-row envelope.
func TestCloneFromReportJSON_SchemaVersion_EmittedPopulated(t *testing.T) {
	assertEnvelopeSchemaVersion(t, canonicalReportResults(), 3)
}

// assertEnvelopeSchemaVersion runs the production writer, decodes
// the result with a minimal local type, and asserts both the
// schemaVersion literal and the row count. Split out to keep each
// test under the function-length budget.
func assertEnvelopeSchemaVersion(t *testing.T, results []Result, wantRows int) {
	t.Helper()
	var buf bytes.Buffer
	if err := writeReportRowsJSON(&buf, results); err != nil {
		t.Fatalf("writeReportRowsJSON: %v", err)
	}
	var env envelopePeek
	if err := json.Unmarshal(buf.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v\nraw: %s", err, buf.String())
	}
	if env.SchemaVersion == nil {
		t.Fatalf("schemaVersion field missing from envelope. raw: %s",
			buf.String())
	}
	if *env.SchemaVersion != expectedSchemaVersionPinned {
		t.Fatalf("schemaVersion mismatch: got %d, want %d",
			*env.SchemaVersion, expectedSchemaVersionPinned)
	}
	if len(env.Rows) != wantRows {
		t.Fatalf("rows length: got %d, want %d", len(env.Rows), wantRows)
	}
}
