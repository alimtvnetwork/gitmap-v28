package cmdllm

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
//   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 go test ./cmdllm/ -run LLMDocsJSONContract
//
// NOTE: this test lives in cmdllm (not cmd) because the encoder
// (encodeLLMDocsJSON) is unexported here. The golden helpers below
// are local, minimal equivalents of the cmd package's golden-test
// infrastructure.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/goldenguard"
)

// TestLLMDocsJSONContract_EmptyIsObjectNotNull pins the empty-sections
// shape so downstream `jq '.commands // []'` style consumers never see
// `null`.
func TestLLMDocsJSONContract_EmptyIsObjectNotNull(t *testing.T) {
	llmAssertGoldenBytesDeterministic(t, "llm_docs_empty.json", func() ([]byte, error) {
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

	llmAssertSchemaKeysFirstObject(t, buf.Bytes(), "llm-docs")
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

	groupKeys := llmReadEveryObjectKeys(t, doc.Commands)
	if len(groupKeys) == 0 {
		t.Fatalf("expected at least one command group")
	}

	wantGroupKeys := []string{"title", "commands"}
	for i, got := range groupKeys {
		if !llmEqualStringSlices(got, wantGroupKeys) {
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

	cmdKeys := llmReadEveryObjectKeys(t, firstGroup.Commands)
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

// --- local golden-test helpers (minimal equivalents of cmd's infra) ---
// llmPackageDir is defined in llmdocs_jsonschema_contract_test.go.

func llmResolveGoldenPath(name string) string {
	return filepath.Join(llmPackageDir(), "testdata", name)
}

func llmAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := llmResolveGoldenPath(name)
	trigger := os.Getenv("GITMAP_UPDATE_GOLDEN") == "1"
	if goldenguard.AllowUpdate(t, trigger) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}

		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}

		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with GITMAP_UPDATE_GOLDEN=1 and GITMAP_ALLOW_GOLDEN_UPDATE=1 to create)", path, err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s\n--- want\n%s\n--- got\n%s", name, string(want), string(got))
	}
}

func llmAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
	t.Helper()
	first, err := encode()
	if err != nil {
		t.Fatalf("%s: encode run 0: %v", name, err)
	}

	for i := 1; i < 3; i++ {
		got, err := encode()
		if err != nil {
			t.Fatalf("%s: encode run %d: %v", name, i, err)
		}

		if !bytes.Equal(got, first) {
			t.Fatalf("%s: determinism broken — run %d differs from run 0", name, i)
		}
	}

	llmAssertGoldenBytes(t, name, first)
}

// llmAssertSchemaKeysFirstObject checks the top-level key order of the
// emitted object against testdata/schemas/<name>.v1.json.
func llmAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(llmPackageDir(), "testdata", "schemas", name+".v1.json")
	schemaRaw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema %s: %v", schemaPath, err)
	}

	var schemaDoc struct {
		Keys []string `json:"keys"`
	}

	if err := json.Unmarshal(schemaRaw, &schemaDoc); err != nil {
		t.Fatalf("parse schema %s: %v", schemaPath, err)
	}

	got := llmFirstObjectKeys(t, raw)
	if len(got) != len(schemaDoc.Keys) {
		t.Fatalf("schema key count mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
	}

	for i := range got {
		if got[i] != schemaDoc.Keys[i] {
			t.Fatalf("schema key order mismatch for %q: got %v, want %v", name, got, schemaDoc.Keys)
		}
	}
}

func llmEqualStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// llmReadEveryObjectKeys streams the entire top-level array and returns
// one []string of keys per object.
func llmReadEveryObjectKeys(t *testing.T, raw []byte) [][]string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := llmExpectDelim(dec, '['); err != nil {
		t.Fatalf("expected top-level array: %v", err)
	}

	var out [][]string
	for dec.More() {
		if err := llmExpectDelim(dec, '{'); err != nil {
			t.Fatalf("expected object at index %d: %v", len(out), err)
		}

		out = append(out, llmCollectObjectKeys(t, dec))
	}

	return out
}

func llmExpectDelim(dec *json.Decoder, want byte) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}

	delim, isDelim := tok.(json.Delim)
	if !isDelim || delim != json.Delim(want) {
		return fmt.Errorf("want delim %q, got %v (%T)", want, tok, tok)
	}

	return nil
}

func llmCollectObjectKeys(t *testing.T, dec *json.Decoder) []string {
	t.Helper()
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("reading object key: %v", err)
		}

		key, ok := tok.(string)
		if !ok {
			t.Fatalf("expected string key, got %v (%T)", tok, tok)
		}

		keys = append(keys, key)

		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			t.Fatalf("skipping value for key %q: %v", key, err)
		}
	}

	if _, err := dec.Token(); err != nil {
		t.Fatalf("expected closing '}': %v", err)
	}

	return keys
}
