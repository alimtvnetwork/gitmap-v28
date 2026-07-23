package cmd

// JSON-schema contract for `gitmap scan-project` per-type files.
//
// `scan-project` writes 5 sibling JSON files (`go-projects.json`,
// `node-projects.json`, `react-projects.json`, `cpp-projects.json`,
// `csharp-projects.json`) — each is a JSON array of detection
// records produced by `buildJSONRecords`. Top-level record keys are
// PascalCase (`Project`, `GoMeta`, `Csharp`) because
// `detector.DetectionResult` has no `json:` tags; this is the
// contractual on-the-wire shape since v1.
//
// These tests pin:
//   1. The 5 file names emitted by `projectTypeJSONMap` exactly
//      match the registry copy.
//   2. Every key produced by `buildJSONRecords` for both the
//      bare-record path (Node/React/Cpp) and the metadata-wrapped
//      path (Go/Csharp) is declared in the schema's
//      `items.properties` map.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/detector"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

const scanProjectSchemaFilename = "scan-project.schema.json"

// TestScanProject_FileMapMatchesRegistry locks the set of emitted
// filenames against the schema-registry copy so silent additions
// or renames break CI.
func TestScanProject_FileMapMatchesRegistry(t *testing.T) {
	got := []string{
		constants.JSONFileGoProjects,
		constants.JSONFileNodeProjects,
		constants.JSONFileReactProjects,
		constants.JSONFileCppProjects,
		constants.JSONFileCsharpProjects,
	}

	regPath := filepath.Join(resolveSchemaDir(), "scan-project.v1.json")
	raw, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatalf("read registry %s: %v", regPath, err)
	}
	var reg struct {
		Files []string `json:"files"`
	}
	if err := json.Unmarshal(raw, &reg); err != nil {
		t.Fatalf("parse registry: %v", err)
	}

	if !equalStringSlices(got, reg.Files) {
		t.Errorf("scan-project files drift:\n got=%v\nwant=%v\nupdate %s if intentional",
			got, reg.Files, regPath)
	}
}

// TestScanProject_RecordKeysSubsetOfSchema runs the live record
// builder for both DetectionResult shape variants and asserts every
// emitted key is declared in the schema's items.properties map.
func TestScanProject_RecordKeysSubsetOfSchema(t *testing.T) {
	root := loadSchemaFile(t, scanProjectSchemaFilename)
	items, _ := root["items"].(map[string]any)
	props, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("items.properties missing")
	}

	results := []detector.DetectionResult{
		{Project: model.DetectedProject{ID: 1, ProjectType: constants.ProjectKeyNode, ProjectName: "demo-node"}},
		{
			Project: model.DetectedProject{ID: 2, ProjectType: constants.ProjectKeyGo, ProjectName: "demo-go"},
			GoMeta:  &model.GoProjectMetadata{ID: 1, ModuleName: "example.com/demo"},
		},
		{
			Project: model.DetectedProject{ID: 3, ProjectType: constants.ProjectKeyCsharp, ProjectName: "demo-cs"},
			Csharp:  &model.CsharpProjectMetadata{ID: 1, SlnName: "Demo.sln"},
		},
	}

	records := buildJSONRecords(results)
	raw, err := json.Marshal(records)
	if err != nil {
		t.Fatalf("marshal records: %v", err)
	}

	keysPerRecord := readEveryObjectKeys(t, raw)
	if len(keysPerRecord) != len(results) {
		t.Fatalf("expected %d records, got %d", len(results), len(keysPerRecord))
	}
	for i, keys := range keysPerRecord {
		for _, key := range keys {
			if _, allowed := props[key]; !allowed {
				t.Errorf("record[%d]: encoder emitted %q not declared in scan-project schema", i, key)
			}
		}
	}
}
