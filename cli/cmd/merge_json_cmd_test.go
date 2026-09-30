package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/jsonenvelope"
)

func TestMergeJSON_DeduplicationAndOneBasedIndex(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "nodes1.json")
	file2 := filepath.Join(tempDir, "nodes2.json")
	outFile := filepath.Join(tempDir, "merged-nodes.json")

	env1 := jsonenvelope.NewEnvelope(jsonenvelope.TypeSSHNodes, file1, "gitmap ssh export", "2.0", map[string]any{
		"nodes": []map[string]any{
			{"id": 0, "alias": "w1", "ipAddress": "192.168.1.10", "username": "admin"},
			{"id": 1, "alias": "w2", "ipAddress": "192.168.1.11", "username": "admin"},
		},
	})
	b1, _ := json.MarshalIndent(env1, "", "  ")
	_ = os.WriteFile(file1, b1, 0644)

	env2 := jsonenvelope.NewEnvelope(jsonenvelope.TypeSSHNodes, file2, "gitmap ssh export", "2.0", map[string]any{
		"nodes": []map[string]any{
			{"id": 0, "alias": "w2", "ipAddress": "192.168.1.11", "username": "admin"}, // Duplicate
			{"id": 1, "alias": "w3", "ipAddress": "192.168.1.12", "username": "admin"},
		},
	})
	b2, _ := json.MarshalIndent(env2, "", "  ")
	_ = os.WriteFile(file2, b2, 0644)

	opts := mergeJSONOptions{
		OutputFile: outFile,
		TargetType: jsonenvelope.TypeSSHNodes,
		IsYes:      true,
	}

	records := collectMergeableRecords([]string{file1, file2}, jsonenvelope.TypeSSHNodes)
	if len(records) != 2 {
		t.Fatalf("expected 2 mergeable records, got %d", len(records))
	}

	summary, data, err := executeRecordsMerge(records, opts)
	if err != nil {
		t.Fatalf("executeRecordsMerge failed: %v", err)
	}

	if summary.TotalItemsIn != 4 {
		t.Errorf("expected 4 total items in, got %d", summary.TotalItemsIn)
	}
	if summary.DuplicateItems != 1 {
		t.Errorf("expected 1 duplicate item, got %d", summary.DuplicateItems)
	}
	if summary.FinalItemCount != 3 {
		t.Errorf("expected 3 final items, got %d", summary.FinalItemCount)
	}

	var parsed struct {
		Attributes jsonenvelope.EnvelopeAttributes `json:"attributes"`
		Data       struct {
			Nodes []struct {
				ID        int    `json:"id"`
				Alias     string `json:"alias"`
				IPAddress string `json:"ipAddress"`
			} `json:"nodes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal merged result failed: %v", err)
	}

	if len(parsed.Data.Nodes) != 3 {
		t.Fatalf("expected 3 nodes in output, got %d", len(parsed.Data.Nodes))
	}

	// Verify 1-based indexing
	for idx, n := range parsed.Data.Nodes {
		expectedID := idx + 1
		if n.ID != expectedID {
			t.Errorf("node %d expected 1-based id %d, got %d", idx, expectedID, n.ID)
		}
	}
}
