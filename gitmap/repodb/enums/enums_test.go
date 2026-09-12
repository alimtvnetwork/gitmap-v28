package enums

import (
	"encoding/json"
	"testing"
)

func TestEnums_TableConstants(t *testing.T) {
	if RepoFileTable != "RepoFile" {
		t.Errorf("expected RepoFileTable 'RepoFile', got %s", RepoFileTable)
	}

	if SearchCacheTable != "SearchCache" {
		t.Errorf("expected SearchCacheTable 'SearchCache', got %s", SearchCacheTable)
	}

	if FileSequenceTable != "FileSequence" {
		t.Errorf("expected FileSequenceTable 'FileSequence', got %s", FileSequenceTable)
	}

	if SequenceHistoryTable != "SequenceHistory" {
		t.Errorf("expected SequenceHistoryTable 'SequenceHistory', got %s", SequenceHistoryTable)
	}

	if RepoScanLogTable != "RepoScanLog" {
		t.Errorf("expected RepoScanLogTable 'RepoScanLog', got %s", RepoScanLogTable)
	}

	if IndexedRepoTable != "IndexedRepo" {
		t.Errorf("expected IndexedRepoTable 'IndexedRepo', got %s", IndexedRepoTable)
	}
}

func testFieldReceivers(t *testing.T, field RepoFileFieldType) {
	t.Helper()
	if field.Name() != "RepoFileId" {
		t.Errorf("expected Name 'RepoFileId', got %s", field.Name())
	}

	if field.String() != "RepoFileId" {
		t.Errorf("expected String 'RepoFileId', got %s", field.String())
	}

	if field.Value() != "RepoFileId" {
		t.Errorf("expected Value 'RepoFileId', got %s", field.Value())
	}

	if !field.IsCompare(RepoFileDb.RepoFileId) {
		t.Errorf("expected IsCompare true")
	}

	if !field.IsEnum() {
		t.Errorf("expected IsEnum true")
	}
}

func TestEnums_FieldReceivers(t *testing.T) {
	field := RepoFileDb.RepoFileId
	testFieldReceivers(t, field)
	if !field.IsRepoFileId() {
		t.Errorf("expected IsRepoFileId true")
	}

	if field.IsRelativePath() {
		t.Errorf("expected IsRelativePath false")
	}
}

func TestEnums_JSONMarshaling(t *testing.T) {
	field := RepoFileDb.RepoFileId
	data, err := json.Marshal(field)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var unmarshaled RepoFileFieldType
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if unmarshaled != field {
		t.Errorf("expected %s, got %s", field, unmarshaled)
	}
}

func TestEnums_Registries(t *testing.T) {
	if SearchCacheDb.SearchCacheId != "SearchCacheId" {
		t.Errorf("expected SearchCacheId")
	}

	if FileSequenceDb.FileSequenceId != "FileSequenceId" {
		t.Errorf("expected FileSequenceId")
	}

	if SequenceHistoryDb.SequenceHistoryId != "SequenceHistoryId" {
		t.Errorf("expected SequenceHistoryId")
	}

	if RepoScanLogDb.RepoScanLogId != "RepoScanLogId" {
		t.Errorf("expected RepoScanLogId")
	}

	if IndexedRepoDb.IndexedRepoId != "IndexedRepoId" {
		t.Errorf("expected IndexedRepoId")
	}
}
