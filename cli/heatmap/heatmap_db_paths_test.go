package heatmap

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveDbPathDetailsInternal(t *testing.T) {
	root := findRepoRoot()
	target := filepath.Join(root, ".ai-memory", "heatmap.db")
	absPath, relPath := resolveDbPathDetails(target)

	assertNoBackslash(t, absPath)
	assertNoBackslash(t, relPath)
	if !strings.HasPrefix(relPath, ".ai-memory/") {
		t.Fatalf("expected relPath starting with .ai-memory/, got %s", relPath)
	}
}

func TestResolveDbPathDetailsExternal(t *testing.T) {
	root := findRepoRoot()
	target := filepath.Join(root, "..", "repo-cache", ".heatmap", "root-heatmap.db")
	absPath, relPath := resolveDbPathDetails(target)

	assertNoBackslash(t, absPath)
	assertNoBackslash(t, relPath)
	if !strings.HasPrefix(relPath, "../repo-cache/") {
		t.Fatalf("expected relPath starting with ../repo-cache/, got %s", relPath)
	}
}

func TestResolveRepoSecretsURLNoPanic(t *testing.T) {
	url := ResolveRepoSecretsURL()
	hasValidReturn := url == url
	if !hasValidReturn {
		t.Fatal("unexpected condition")
	}
}

func TestHeatmapTablesJSONPaths(t *testing.T) {
	var buf bytes.Buffer
	tables := []TableInfo{{Name: "TestTable", Type: "table", ColumnCount: 2, RowCount: 5}}
	err := formatTablesOutput(&buf, ".ai-memory/heatmap.db", tables, "json")
	if err != nil {
		t.Fatal(err)
	}
	assertDbEnvelopeKeys(t, buf.Bytes())
}

func TestHeatmapViewJSONPaths(t *testing.T) {
	var buf bytes.Buffer
	cols := []string{"id", "title"}
	records := []map[string]any{{"id": 1, "title": "Alpha"}}
	err := formatViewOutput(&buf, ".ai-memory/heatmap.db", "Tasks", 1, 10, cols, records, "json")
	if err != nil {
		t.Fatal(err)
	}
	assertDbEnvelopeKeys(t, buf.Bytes())
}

func TestHeatmapQueryJSONPaths(t *testing.T) {
	var buf bytes.Buffer
	cols := []string{"val"}
	rows := [][]any{{1}}
	err := renderQueryJSON(".ai-memory/heatmap.db", "SELECT 1", cols, rows, 0.5, &buf)
	if err != nil {
		t.Fatal(err)
	}
	assertDbEnvelopeKeys(t, buf.Bytes())
}

func TestHeatmapQueryYAMLPaths(t *testing.T) {
	var buf bytes.Buffer
	cols := []string{"val"}
	rows := [][]any{{1}}
	err := renderQueryYAML(".ai-memory/heatmap.db", "SELECT 1", cols, rows, 0.5, &buf)
	if err != nil {
		t.Fatal(err)
	}
	assertYamlContainsKeys(t, buf.String())
}

func assertDbEnvelopeKeys(t *testing.T, data []byte) {
	t.Helper()
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	assertMapKeyExists(t, parsed, "database")
	assertMapKeyExists(t, parsed, "relativeDbPath")
	assertMapKeyExists(t, parsed, "absoluteDbPath")
	assertMapKeyExists(t, parsed, "repoSecretsUrl")
}

func assertMapKeyExists(t *testing.T, m map[string]any, key string) {
	t.Helper()
	_, isFound := m[key]
	if !isFound {
		t.Fatalf("missing required key in JSON envelope: %s", key)
	}
}

func assertNoBackslash(t *testing.T, path string) {
	t.Helper()
	hasBackslash := strings.Contains(path, `\`)
	if hasBackslash {
		t.Fatalf("path contains backslash: %s", path)
	}
}

func assertYamlContainsKeys(t *testing.T, content string) {
	t.Helper()
	hasRelative := strings.Contains(content, "relativeDbPath:")
	if !hasRelative {
		t.Fatal("YAML missing relativeDbPath")
	}
	hasSecrets := strings.Contains(content, "repoSecretsUrl:")
	if !hasSecrets {
		t.Fatal("YAML missing repoSecretsUrl")
	}
}
