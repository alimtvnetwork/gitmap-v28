package cmdagy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const batchCursorRelativePath = ".ai-memory/temp/pipeline-fix-batch-cursor.json"

func resolveBatchCursorPath() string {
	return filepath.Join(resolveProjectRootDir(), batchCursorRelativePath)
}

// LoadPipelineFixBatchCursor loads or initializes cursor state for multi-project batching.
func LoadPipelineFixBatchCursor(path string, defaultLimit int) PipelineFixBatchCursor {
	cursor := PipelineFixBatchCursor{BatchLimit: defaultLimit}
	data, err := os.ReadFile(path)
	if err != nil {
		return cursor
	}

	_ = json.Unmarshal(data, &cursor)
	if cursor.BatchLimit <= 0 {
		cursor.BatchLimit = defaultLimit
	}

	return cursor
}

// SavePipelineFixBatchCursor persists the cursor state to disk.
func SavePipelineFixBatchCursor(path string, cursor PipelineFixBatchCursor) error {
	cursor.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_ = os.MkdirAll(filepath.Dir(path), 0755)
	marshaled, err := json.MarshalIndent(cursor, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, marshaled, 0644)
}

// ResetPipelineFixBatchCursor resets the cursor to the beginning.
func ResetPipelineFixBatchCursor(path string) error {
	cursor := PipelineFixBatchCursor{LastIndex: 0, BatchLimit: 3}

	return SavePipelineFixBatchCursor(path, cursor)
}
