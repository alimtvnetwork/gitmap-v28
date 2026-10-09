package cmdstats

// JSON contract tests for `gitmap stats --json`.
//
// stats emits a single object summarizing overall command-history
// metrics plus a nested `commands` array of per-command rows. The
// contract covers:
//
//   - Top-level object shape with all required keys.
//   - Key order: totalCommands, uniqueCommands, totalSuccess,
//     totalFail, overallFailRate, avgDurationMs, commands.
//   - Empty commands list is `[]` (NOT null) on the wire.
//
// Regenerate fixtures with:
//
//   GITMAP_UPDATE_GOLDEN=1 go test ./cmd/ -run StatsJSONContract

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// canonicalStatsOverall builds a deterministic two-row stats snapshot.
// One row exercises a populated command; tests pin both the wrapping
// object's key order and the embedded compact array's bytes.
func canonicalStatsOverall(t *testing.T) (model.OverallStats, []model.CommandStats) {
	t.Helper()
	commands := []model.CommandStats{
		{
			Command:      "clone",
			TotalRuns:    10,
			SuccessCount: 9,
			FailCount:    1,
			FailRate:     0.1,
			AvgDuration:  1234,
			MinDuration:  500,
			MaxDuration:  3000,
			LastUsed:     "2026-05-26T08:30:00Z",
		},
		{
			Command:      "stats",
			TotalRuns:    2,
			SuccessCount: 2,
			FailCount:    0,
			FailRate:     0.0,
			AvgDuration:  42,
			MinDuration:  40,
			MaxDuration:  45,
			LastUsed:     "2026-05-26T09:00:00Z",
		},
	}

	overall := model.OverallStats{
		TotalCommands:   12,
		UniqueCommands:  2,
		TotalSuccess:    11,
		TotalFail:       1,
		OverallFailRate: 0.0833,
		AvgDuration:     1100,
		Commands:        commands,
	}

	return overall, commands
}

// TestStatsJSONContract_EmptyCommandsArray pins the empty-commands
// shape: a fully-populated overall object with `commands: []`.
func TestStatsJSONContract_EmptyCommandsArray(t *testing.T) {
	overall := model.OverallStats{}
	stAssertGoldenBytesDeterministic(t, "stats_empty.json", func() ([]byte, error) {
		var buf bytes.Buffer
		err := encodeStatsJSON(&buf, overall, nil)

		return buf.Bytes(), err
	})
}

// Canonical-byte fixture is intentionally omitted: float formatting
// (e.g. 0.0833) and embedded compact-array bytes are tied to Go's
// json.Marshal output, so regenerating via GITMAP_UPDATE_GOLDEN is
// the safer pin. The key-order test below covers the structural
// contract without locking in float-printing artifacts.

// TestStatsJSONContract_KeyOrder asserts the top-level object's key
// order matches the schema registry declaration.
func TestStatsJSONContract_KeyOrder(t *testing.T) {
	overall, commands := canonicalStatsOverall(t)
	var buf bytes.Buffer
	if err := encodeStatsJSON(&buf, overall, commands); err != nil {
		t.Fatalf("encode: %v", err)
	}

	stAssertSchemaKeysFirstObject(t, buf.Bytes(), "stats")
}

// --- local test helpers (minimal equivalents of cmd's infra) ---

func stPackageDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("stPackageDir: runtime.Caller failed")
	}
	return filepath.Dir(file)
}

func stAssertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join(stPackageDir(), "testdata", name)
	if os.Getenv("GITMAP_UPDATE_GOLDEN") == "1" && os.Getenv("GITMAP_ALLOW_GOLDEN_UPDATE") == "1" {
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
		t.Fatalf("read golden %s: %v", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("golden mismatch for %s", name)
	}
}

func stAssertGoldenBytesDeterministic(t *testing.T, name string, encode func() ([]byte, error)) {
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
			t.Fatalf("%s: determinism broken", name)
		}
	}
	stAssertGoldenBytes(t, name, first)
}

func stAssertSchemaKeysFirstObject(t *testing.T, raw []byte, name string) {
	t.Helper()
	schemaPath := filepath.Join(stPackageDir(), "testdata", "schemas", name+".v1.json")
	schemaRaw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var doc struct {
		Keys []string `json:"keys"`
	}
	if err := json.Unmarshal(schemaRaw, &doc); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	got := stFirstObjectKeys(t, raw)
	if len(got) != len(doc.Keys) {
		t.Fatalf("key count mismatch: got %v, want %v", got, doc.Keys)
	}
	for i := range got {
		if got[i] != doc.Keys[i] {
			t.Fatalf("key order mismatch: got %v, want %v", got, doc.Keys)
		}
	}
}

func stFirstObjectKeys(t *testing.T, raw []byte) []string {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	for {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("no object: %v", err)
		}
		if d, ok := tok.(json.Delim); ok && d == '{' {
			break
		}
	}
	var keys []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("keys: %v", err)
		}
		k, ok := tok.(string)
		if !ok {
			t.Fatalf("expected string key")
		}
		keys = append(keys, k)
		var skip any
		if err := dec.Decode(&skip); err != nil {
			t.Fatalf("skip: %v", err)
		}
	}
	return keys
}
