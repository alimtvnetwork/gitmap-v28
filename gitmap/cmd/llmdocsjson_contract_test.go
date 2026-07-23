package cmd

// JSON contract tests for `gitmap llm-docs --format=json`.
//
// llm-docs emits a single top-level object whose keys are conditionally
// appended based on the --sections filter. Contract covers:
//
//   - Empty filter (no sections) emits `{}\n`.
//   - When all sections are included, the top-level key order matches
//     the schema registry declaration verbatim.
//   - Nested `commands` group + per-command key order is contractual.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run LLMDocsJSONContract

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TestLLMDocsJSONContract_EmptyIsObjectNotNull pins the empty-sections
// shape so downstream `jq '.commands // []'` style consumers never see
// `null`.
func TestLLMDocsJSONContract_EmptyIsObjectNotNull(t *testing.T) {
	assertGoldenBytesDeterministic(t, "llm_docs_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeLLMDocsJSON(&buf, map[string]bool{})

		return buf.Bytes(), err
	})
}

// TestLLMDocsJSONContract_TopLevelKeyOrder asserts the top-level key
// order with all sections enabled matches the schema registry.
func TestLLMDocsJSONContract_TopLevelKeyOrder(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeLLMDocsJSON(&buf, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	assertSchemaKeysFirstObject(t, buf.Bytes(), "llm-docs")
}

// TestLLMDocsJSONContract_CommandGroupKeyOrder asserts every nested
// command group emits keys in (title, commands) order, and every
// per-command record starts with (name, alias, description) — the
// optional `example` is appended only when non-empty.
func TestLLMDocsJSONContract_CommandGroupKeyOrder(t *testing.T) {
	var buf bytes.Buffer
	if err := encodeLLMDocsJSON(&buf, map[string]bool{"commands": true}); err != nil {
		t.Fatalf("encode: %v", err)
	}
	var doc struct {
		Commands json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(doc.Commands) == 0 {
		t.Fatalf("expected non-empty commands array")
	}

	groupKeys := readEveryObjectKeys(t, doc.Commands)
	if len(groupKeys) == 0 {
		t.Fatalf("expected at least one command group")
	}
	wantGroupKeys := []string{"title", "commands"}
	for i, got := range groupKeys {
		if !equalStringSlices(got, wantGroupKeys) {
			t.Fatalf("group[%d] keys = %v, want %v", i, got, wantGroupKeys)
		}
	}

	// Drill into the first group's commands array and verify each
	// per-command record's prefix is (name, alias, description) with
	// optional trailing `example`.
	var firstGroup struct {
		Commands json.RawMessage `json:"commands"`
	}
	// Re-parse the first element of the groups array.
	var groupsArr []json.RawMessage
	if err := json.Unmarshal(doc.Commands, &groupsArr); err != nil {
		t.Fatalf("unmarshal groups: %v", err)
	}
	if err := json.Unmarshal(groupsArr[0], &firstGroup); err != nil {
		t.Fatalf("unmarshal first group: %v", err)
	}
	cmdKeys := readEveryObjectKeys(t, firstGroup.Commands)
	wantPrefix := []string{"name", "alias", "description"}
	for i, got := range cmdKeys {
		if len(got) < len(wantPrefix) {
			t.Fatalf("command[%d] only has keys %v", i, got)
		}
		for j, k := range wantPrefix {
			if got[j] != k {
				t.Fatalf("command[%d] key[%d] = %q, want %q (full=%v)", i, j, got[j], k, got)
			}
		}
		if len(got) == 4 && got[3] != "example" {
			t.Fatalf("command[%d] optional 4th key = %q, want \"example\"", i, got[3])
		}
		if len(got) > 4 {
			t.Fatalf("command[%d] has unexpected extra keys: %v", i, got)
		}
	}
}
