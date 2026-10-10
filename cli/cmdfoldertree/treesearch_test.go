package cmdfoldertree

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTreeFile_JSONFolderReport(t *testing.T) {
	jsonContent := `{
		"root": ".",
		"totalFiles": 2,
		"files": [
			{
				"path": "cli/cmd/root.go",
				"filename": "root.go",
				"directory": "cli/cmd",
				"extension": ".go",
				"sizeBytes": 2048
			},
			{
				"path": "cli/cmd/roottooling.go",
				"filename": "roottooling.go",
				"directory": "cli/cmd",
				"extension": ".go",
				"sizeBytes": 4096
			}
		]
	}`

	tmpFile := filepath.Join(t.TempDir(), "report.json")
	if err := os.WriteFile(tmpFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	items, err := ParseTreeFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseTreeFile failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	if items[0].RelPath != "cli/cmd/root.go" || items[0].FileName != "root.go" {
		t.Errorf("unexpected first item: %+v", items[0])
	}
	if items[0].IsDir {
		t.Errorf("expected isDir=false, got true")
	}
	if items[1].SizeBytes != 4096 {
		t.Errorf("expected sizeBytes=4096, got %d", items[1].SizeBytes)
	}
}

func TestParseTreeFile_JSONExportDoc(t *testing.T) {
	jsonContent := `{
		"rootPath": ".",
		"rootName": "gitmap",
		"totalNodes": 3,
		"tree": {
			"name": "cli",
			"relPath": "cli",
			"isDir": true,
			"children": [
				{
					"name": "cmd",
					"relPath": "cli/cmd",
					"isDir": true,
					"children": [
						{
							"name": "root.go",
							"relPath": "cli/cmd/root.go",
							"isDir": false
						}
					]
				}
			]
		}
	}`

	tmpFile := filepath.Join(t.TempDir(), "export.json")
	if err := os.WriteFile(tmpFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	items, err := ParseTreeFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseTreeFile failed: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("expected 3 items (cli, cli/cmd, cli/cmd/root.go), got %d", len(items))
	}

	// Verify items contain directory and file
	hasDir := false
	hasFile := false
	for _, it := range items {
		if it.IsDir && it.RelPath == "cli" {
			hasDir = true
		}
		if !it.IsDir && it.RelPath == "cli/cmd/root.go" {
			hasFile = true
		}
	}
	if !hasDir || !hasFile {
		t.Errorf("expected both dir and file parsed, got %+v", items)
	}
}

func TestParseTreeFile_YAML(t *testing.T) {
	yamlContent := `
root: .
files:
  - path: cli/cmd/root.go
    filename: root.go
    directory: cli/cmd
    extension: .go
    sizeBytes: 1024
`

	tmpFile := filepath.Join(t.TempDir(), "report.yaml")
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	items, err := ParseTreeFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseTreeFile failed: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].RelPath != "cli/cmd/root.go" {
		t.Errorf("expected cli/cmd/root.go, got %s", items[0].RelPath)
	}
}

func TestParseTreeFile_IndentedAsciiText(t *testing.T) {
	txtContent := `
cli/
├── cmd/
│   ├── root.go
│   └── roottooling.go
└── cmdfoldertree/
    └── types.go
`

	tmpFile := filepath.Join(t.TempDir(), "tree.txt")
	if err := os.WriteFile(tmpFile, []byte(txtContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	items, err := ParseTreeFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseTreeFile failed: %v", err)
	}

	seenPaths := make(map[string]bool)
	for _, it := range items {
		seenPaths[it.RelPath] = true
	}

	if !seenPaths["cli/cmd/root.go"] {
		t.Errorf("missing cli/cmd/root.go in parsed items: %+v", seenPaths)
	}
	if !seenPaths["cli/cmdfoldertree/types.go"] {
		t.Errorf("missing cli/cmdfoldertree/types.go in parsed items: %+v", seenPaths)
	}
}

func TestParseTreeFile_FlatTextList(t *testing.T) {
	txtContent := `
cli/cmd/root.go
cli/cmd/roottooling.go
docs/readme.md
`

	tmpFile := filepath.Join(t.TempDir(), "flat.txt")
	if err := os.WriteFile(tmpFile, []byte(txtContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	items, err := ParseTreeFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseTreeFile failed: %v", err)
	}

	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[2].RelPath != "docs/readme.md" || items[2].Extension != ".md" {
		t.Errorf("unexpected item: %+v", items[2])
	}
}

func TestMatchTreeItem_Modes(t *testing.T) {
	itemFile := TreeItem{
		RelPath:   "cli/cmdfoldertree/treesearch_cmd.go",
		FileName:  "treesearch_cmd.go",
		DirPath:   "cli/cmdfoldertree",
		Extension: ".go",
		IsDir:     false,
	}

	itemDir := TreeItem{
		RelPath:   "cli/cmdfoldertree",
		FileName:  "cmdfoldertree",
		DirPath:   "cli",
		Extension: "",
		IsDir:     true,
	}

	// 1. Wildcard
	if !MatchTreeItem(itemFile, ModeWildcard, "*search*cmd*") {
		t.Errorf("wildcard *search*cmd* should match %s", itemFile.RelPath)
	}
	if !MatchTreeItem(itemFile, ModeWildcard, "*.go") {
		t.Errorf("wildcard *.go should match %s", itemFile.RelPath)
	}
	if MatchTreeItem(itemFile, ModeWildcard, "*.txt") {
		t.Errorf("wildcard *.txt should NOT match %s", itemFile.RelPath)
	}

	// 2. StartsWith (with folder path preservation)
	if !MatchTreeItem(itemFile, ModeStartsWith, "cli/cmdfoldertree") {
		t.Errorf("startsWith cli/cmdfoldertree should match %s", itemFile.RelPath)
	}
	if !MatchTreeItem(itemFile, ModeStartsWith, "tree") {
		t.Errorf("startsWith tree should match %s", itemFile.FileName)
	}
	if MatchTreeItem(itemFile, ModeStartsWith, "other") {
		t.Errorf("startsWith other should NOT match %s", itemFile.RelPath)
	}

	// 3. Contains
	if !MatchTreeItem(itemFile, ModeContains, "search") {
		t.Errorf("contains search should match %s", itemFile.RelPath)
	}
	if !MatchTreeItem(itemDir, ModeContains, "folder") {
		t.Errorf("contains folder should match %s", itemDir.RelPath)
	}
	if MatchTreeItem(itemFile, ModeContains, "nonexistent") {
		t.Errorf("contains nonexistent should NOT match")
	}

	// 4. EndsWith
	if !MatchTreeItem(itemFile, ModeEndsWith, ".go") {
		t.Errorf("endsWith .go should match %s", itemFile.RelPath)
	}
	if !MatchTreeItem(itemFile, ModeEndsWith, "cmd.go") {
		t.Errorf("endsWith cmd.go should match %s", itemFile.RelPath)
	}
	if MatchTreeItem(itemFile, ModeEndsWith, ".json") {
		t.Errorf("endsWith .json should NOT match")
	}

	// 5. Grep
	if !MatchTreeItem(itemFile, ModeGrep, "^cli/cmdfoldertree/.*_cmd\\.go$") {
		t.Errorf("grep regex should match %s", itemFile.RelPath)
	}
	if MatchTreeItem(itemFile, ModeGrep, "^docs/") {
		t.Errorf("grep ^docs/ should NOT match %s", itemFile.RelPath)
	}
}

func TestSQLiteSplitDB_IngestAndQuery(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "treedb", "test.db")

	items := []TreeItem{
		{
			RootPath:  ".",
			RelPath:   "cli/cmd/root.go",
			FileName:  "root.go",
			DirPath:   "cli/cmd",
			Extension: ".go",
			SizeBytes: 1024,
			IsDir:     false,
			Depth:     2,
		},
		{
			RootPath:  ".",
			RelPath:   "cli/cmdfoldertree/treesearch.go",
			FileName:  "treesearch.go",
			DirPath:   "cli/cmdfoldertree",
			Extension: ".go",
			SizeBytes: 2048,
			IsDir:     false,
			Depth:     2,
		},
		{
			RootPath:  ".",
			RelPath:   "cli/cmdfoldertree",
			FileName:  "cmdfoldertree",
			DirPath:   "cli",
			Extension: "",
			SizeBytes: 0,
			IsDir:     true,
			Depth:     1,
		},
	}

	ingested, err := IngestTreeItems(dbPath, items)
	if err != nil {
		t.Fatalf("IngestTreeItems failed: %v", err)
	}
	if ingested != 3 {
		t.Fatalf("expected 3 ingested, got %d", ingested)
	}

	// Query startsWith
	results, err := QueryTreeDb(dbPath, ModeStartsWith, "cli/cmdfoldertree")
	if err != nil {
		t.Fatalf("QueryTreeDb failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 items for cli/cmdfoldertree, got %d", len(results))
	}

	// Query contains
	results, err = QueryTreeDb(dbPath, ModeContains, "search")
	if err != nil {
		t.Fatalf("QueryTreeDb failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 item for search, got %d", len(results))
	}
	if results[0].RelPath != "cli/cmdfoldertree/treesearch.go" {
		t.Errorf("unexpected query result: %+v", results[0])
	}
}

func TestParseTreeSearchArgs(t *testing.T) {
	args := []string{"-f", "tree.json", "*root*.go", "--json", "-n", "50", "--dirs"}
	parsed := parseTreeSearchArgs(args)

	if parsed.filePath != "tree.json" {
		t.Errorf("expected filePath tree.json, got %s", parsed.filePath)
	}
	if parsed.pattern != "*root*.go" {
		t.Errorf("expected pattern *root*.go, got %s", parsed.pattern)
	}
	if !parsed.isJSON {
		t.Errorf("expected isJSON=true")
	}
	if parsed.limit != 50 {
		t.Errorf("expected limit=50, got %d", parsed.limit)
	}
	if !parsed.dirsOnly {
		t.Errorf("expected dirsOnly=true")
	}
}

func TestDispatchFolderTree_Recognition(t *testing.T) {
	commands := []string{
		"tree-search", "treesearch", "ts",
		"tree-search-startsWith", "tree-search-startswith", "tss",
		"tree-search-contains", "tsc",
		"tree-search-endsWith", "tse",
		"tree-search-grep", "tsg",
		"tree-learn", "tl",
	}

	for _, cmd := range commands {
		if !isTreeSearchCommand(cmd) && !isTreeLearnCommand(cmd) {
			t.Errorf("command %s should be recognized", cmd)
		}
	}
}
