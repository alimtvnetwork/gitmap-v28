package cmdagy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ProjectPromptTree groups active prompt summaries under a project and node.
type ProjectPromptTree struct {
	NodeAlias     string                `json:"nodeAlias"`
	ProjectName   string                `json:"projectName"`
	WorkspacePath string                `json:"workspacePath"`
	ActivePrompts []ActivePromptSummary `json:"activePrompts"`
	QueuedPrompts []AgyPromptQueueEntry `json:"queuedPrompts"`
}

// PromptResendPayload models an inbound request to resend a past prompt.
type PromptResendPayload struct {
	ID            string `json:"id"`
	ProjectTarget string `json:"projectTarget,omitempty"`
}

// PromptRestorePayload models an inbound request to restore a backup batch.
type PromptRestorePayload struct {
	BatchID string `json:"batchId"`
}

func handleAgyUIPromptResend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var req PromptResendPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		http.Error(w, `{"error":"prompt id required"}`, http.StatusBadRequest)
		return
	}

	history, _ := GetPromptHistory()
	targetRec := findPromptRecordByID(history, req.ID)
	if targetRec == nil {
		saved, _ := LoadSavedPrompts()
		targetRec = findPromptRecordByID(saved, req.ID)
	}

	if targetRec == nil {
		http.Error(w, `{"error":"prompt not found in history"}`, http.StatusNotFound)
		return
	}

	target := targetRec.ProjectTarget
	if req.ProjectTarget != "" {
		target = req.ProjectTarget
	}

	res, err := SendPromptImmediate(PromptPayload{
		ProjectTarget: target,
		Title:         targetRec.Title,
		PromptText:    targetRec.PromptText,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	_ = json.NewEncoder(w).Encode(res)
}

func sanitizeWorkspacePath(path string) string {
	if path == "" || path == "." {
		return "."
	}
	base := filepath.Base(path)
	return filepath.Join("$WORKSPACE_ROOT", base)
}

func handleAgyUIPromptTree(w http.ResponseWriter, _ *http.Request) {
	status, err := GetAGYUIStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}

	treeMap := make(map[string]*ProjectPromptTree)
	for _, p := range status.RunningProjects {
		key := p.ProjectName
		if key == "" {
			key = filepath.Base(p.ProjectPath)
		}
		summary := buildActiveSummaryFromProject(p)
		summary.WorkspacePath = sanitizeWorkspacePath(p.ProjectPath)
		treeMap[key] = &ProjectPromptTree{
			NodeAlias:     status.NodeAlias,
			ProjectName:   p.ProjectName,
			WorkspacePath: summary.WorkspacePath,
			ActivePrompts: []ActivePromptSummary{summary},
			QueuedPrompts: []AgyPromptQueueEntry{},
		}
	}

	for _, q := range status.QueuedPrompts {
		key := q.ProjectName
		if key == "" {
			key = "default"
		}
		if item, exists := treeMap[key]; exists {
			item.QueuedPrompts = append(item.QueuedPrompts, q)
		} else {
			treeMap[key] = &ProjectPromptTree{
				NodeAlias:     status.NodeAlias,
				ProjectName:   key,
				WorkspacePath: "$WORKSPACE_ROOT/" + key,
				ActivePrompts: []ActivePromptSummary{},
				QueuedPrompts: []AgyPromptQueueEntry{q},
			}
		}
	}

	var treeList []ProjectPromptTree
	for _, item := range treeMap {
		treeList = append(treeList, *item)
	}

	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	_ = json.NewEncoder(w).Encode(treeList)
}

// BackupRunningPrompts triggers a snapshot batch of all running and queued prompts.
func BackupRunningPrompts() (store.PromptBackupSummary, error) {
	return RunRunningPromptsBackupSummary("")
}

func handleAgyUIPromptBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	summary, err := BackupRunningPrompts()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"batchId": summary.BatchId,
		"count":   summary.TotalPrompts,
	})
}

func handleAgyUIPromptRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	var payload PromptRestorePayload
	_ = json.NewDecoder(r.Body).Decode(&payload)

	db, err := store.OpenBackupPromptsSplitDB("")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":%q}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	summary, err := fetchLatestOrTargetBackup(db, payload.BatchID)
	if err != nil || summary == nil {
		http.Error(w, `{"error":"no backup batch found to restore"}`, http.StatusNotFound)
		return
	}

	restoredCount := restoreItemsToWorkspaces(summary.Items)
	_ = db.MarkBatchRestored(summary.BatchId, true, restoredCount)

	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"restored": restoredCount,
	})
}

func handleAgyUIPromptReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	emptyQueue := AgyPromptQueueFile{
		Queued:    make([]AgyPromptQueueEntry, 0),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	_ = SavePromptQueue(emptyQueue)
	if activeText := resolveActiveTextPath("."); activeText != "" {
		_ = os.Remove(activeText)
	}
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func handleAgyUIPromptRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, `{"error":"GET or POST required"}`, http.StatusMethodNotAllowed)
		return
	}
	handleAgyUIPromptTree(w, r)
}

func fetchLatestOrTargetBackup(db *store.BackupPromptsSplitDB, batchID string) (*store.PromptBackupSummary, error) {
	if len(batchID) > 0 {
		return db.GetBackupBatch(batchID)
	}
	batches, err := db.ListBackupBatches()
	if err != nil || len(batches) == 0 {
		return nil, err
	}
	return db.GetBackupBatch(batches[0].BatchID)
}
