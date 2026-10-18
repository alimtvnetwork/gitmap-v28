package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PromptPayload models an inbound prompt dispatch or enqueue request.
type PromptPayload struct {
	ProjectName   string `json:"projectName,omitempty"`
	ProjectTarget string `json:"projectTarget,omitempty"`
	TargetNode    string `json:"targetNode,omitempty"`
	Title         string `json:"title,omitempty"`
	PromptText    string `json:"promptText"`
	TaskCategory  string `json:"taskCategory,omitempty"`
	IsUrgent      bool   `json:"isUrgent,omitempty"`
	IsEnqueue     bool   `json:"isEnqueue,omitempty"`
	DelaySeconds  int    `json:"delaySeconds,omitempty"`
}

// EffectiveTarget returns the target project identifier.
func (p PromptPayload) EffectiveTarget() string {
	if p.ProjectName != "" {
		return p.ProjectName
	}
	if p.ProjectTarget != "" {
		return p.ProjectTarget
	}
	return "."
}

// EffectiveTitle returns the human-readable prompt title.
func (p PromptPayload) EffectiveTitle() string {
	if p.Title != "" {
		return p.Title
	}
	return "Prompt Dispatch"
}

// PromptRecord stores state and metadata for a prompt across its lifecycle.
type PromptRecord struct {
	ID            string `json:"id"`
	Status        string `json:"status"` // "running", "queued", "saved", "completed", "failed"
	QueuePosition int    `json:"queuePosition,omitempty"`
	ProjectName   string `json:"projectName,omitempty"`
	ProjectTarget string `json:"projectTarget,omitempty"`
	TargetNode    string `json:"targetNode,omitempty"`
	Title         string `json:"title,omitempty"`
	PromptText    string `json:"promptText"`
	CreatedAt     string `json:"createdAt"`
	CompletedAt   string `json:"completedAt,omitempty"`
}

// PromptResult represents the execution outcome of a prompt action.
type PromptResult struct {
	PromptID      string `json:"promptId"`
	Status        string `json:"status"` // "delivered", "executing", "completed", "failed", "queued"
	TargetNode    string `json:"targetNode"`
	ProjectName   string `json:"projectName"`
	DispatchedAt  string `json:"dispatchedAt"`
	ExecutionPID  int    `json:"executionPid,omitempty"`
	OutputPreview string `json:"outputPreview,omitempty"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
}

// PromptResendRequest encapsulates parameters to retry or re-route a prompt.
type PromptResendRequest struct {
	PromptID     string `json:"promptId"`
	OverrideNode string `json:"overrideNode,omitempty"`
	IsUrgent     bool   `json:"isUrgent,omitempty"`
}

// PromptTreeResponse provides the nested hierarchy of prompts by node and project.
type PromptTreeResponse struct {
	TotalActive int              `json:"totalActive"`
	TotalQueued int              `json:"totalQueued"`
	Nodes       []NodePromptTree `json:"nodes"`
}

// NodePromptTree represents a node and its nested projects.
type NodePromptTree struct {
	NodeAlias string                 `json:"nodeAlias"`
	IsOnline  bool                   `json:"isOnline"`
	Projects  []ProjectPromptSummary `json:"projects"`
}

// ProjectPromptSummary contains prompt details for a specific project.
type ProjectPromptSummary struct {
	ProjectName   string         `json:"projectName"`
	WorkspacePath string         `json:"workspacePath"`
	ActivePrompts []PromptRecord `json:"activePrompts"`
	QueuedPrompts []PromptRecord `json:"queuedPrompts"`
}

// ActivePromptSummary provides lightweight telemetry for active prompts.
type ActivePromptSummary struct {
	WorkspacePath  string `json:"workspacePath"`
	ProjectName    string `json:"projectName"`
	ConversationID string `json:"conversationId,omitempty"`
	PromptTitle    string `json:"promptTitle"`
	PromptSnippet  string `json:"promptSnippet"`
	Duration       string `json:"duration,omitempty"`
	Status         string `json:"status"`
}

// AGYUIStatusPayload represents the full state returned to the web UI.
type AGYUIStatusPayload struct {
	IsServerRunning bool                   `json:"isServerRunning"`
	NodeAlias       string                 `json:"nodeAlias"`
	RunningProjects []RunningProjectRecord `json:"runningProjects"`
	ActivePrompts   []ActivePromptSummary  `json:"activePrompts"`
	QueuedPrompts   []AgyPromptQueueEntry  `json:"queuedPrompts"`
	SavedPrompts    []PromptRecord         `json:"savedPrompts"`
}

func getTempAgyFilePath(filename string) string {
	return filepath.Join(".ai-memory", "temp", filename)
}

func ensureTempDir() error {
	dir := filepath.Join(".ai-memory", "temp")
	return os.MkdirAll(dir, 0755)
}

func readRecordsFromFile(filePath string) ([]PromptRecord, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []PromptRecord{}, nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read records")
	}
	var records []PromptRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return []PromptRecord{}, nil
	}
	return records, nil
}

func writeRecordsToFile(filePath string, records []PromptRecord) error {
	if err := ensureTempDir(); err != nil {
		return apperror.WrapSimple(err, "ensure temp dir")
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal records")
	}
	return os.WriteFile(filePath, data, 0644)
}

// LoadSavedPrompts loads prompt templates saved by the user.
func LoadSavedPrompts() ([]PromptRecord, error) {
	path := getTempAgyFilePath("agy-saved-prompts.json")
	return readRecordsFromFile(path)
}

// SavePromptTemplate persists a reusable prompt template.
func SavePromptTemplate(record PromptRecord) error {
	path := getTempAgyFilePath("agy-saved-prompts.json")
	records, _ := readRecordsFromFile(path)
	if record.ID == "" {
		record.ID = fmt.Sprintf("tpl-%d", time.Now().UnixNano())
	}
	if record.CreatedAt == "" {
		record.CreatedAt = time.Now().Format(time.RFC3339)
	}
	record.Status = "saved"
	records = append([]PromptRecord{record}, records...)
	return writeRecordsToFile(path, records)
}

// GetPromptHistory loads historical dispatches from disk.
func GetPromptHistory() ([]PromptRecord, error) {
	path := getTempAgyFilePath("agy-prompt-history.json")
	return readRecordsFromFile(path)
}

// AppendPromptHistory adds an execution record to history log.
func AppendPromptHistory(record PromptRecord) error {
	path := getTempAgyFilePath("agy-prompt-history.json")
	records, _ := readRecordsFromFile(path)
	records = append([]PromptRecord{record}, records...)
	if len(records) > 100 {
		records = records[:100]
	}
	return writeRecordsToFile(path, records)
}

// EnqueuePromptPayload adds a prompt to the file-backed queue.
func EnqueuePromptPayload(payload PromptPayload) (*PromptRecord, error) {
	path := getTempAgyFilePath("agy-prompt-queue.json")
	records, _ := readRecordsFromFile(path)
	pos := len(records) + 1
	node := payload.TargetNode
	if node == "" {
		node, _ = os.Hostname()
	}
	rec := PromptRecord{
		ID:            fmt.Sprintf("q-%d", time.Now().UnixNano()),
		ProjectName:   payload.EffectiveTarget(),
		ProjectTarget: payload.EffectiveTarget(),
		TargetNode:    node,
		Title:         payload.EffectiveTitle(),
		PromptText:    payload.PromptText,
		Status:        "queued",
		QueuePosition: pos,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	records = append(records, rec)
	if err := writeRecordsToFile(path, records); err != nil {
		return nil, err
	}
	_, _ = EnqueuePrompt("ui_enqueue", rec.Title, payload.PromptText)
	return &rec, nil
}

// SendPromptImmediate dispatches a prompt immediately to a target project.
func SendPromptImmediate(payload PromptPayload) (AgyInjectionResult, error) {
	target := payload.EffectiveTarget()
	title := payload.EffectiveTitle()
	res := InjectAgyPrompt(target, payload.PromptText, title, "ui_send", false)
	node := payload.TargetNode
	if node == "" {
		node, _ = os.Hostname()
	}
	rec := PromptRecord{
		ID:            fmt.Sprintf("disp-%d", time.Now().UnixNano()),
		ProjectName:   target,
		ProjectTarget: target,
		TargetNode:    node,
		Title:         title,
		PromptText:    payload.PromptText,
		Status:        "running",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	_ = AppendPromptHistory(rec)
	return res, nil
}

func findPromptRecordByID(records []PromptRecord, id string) *PromptRecord {
	for _, rec := range records {
		if rec.ID == id {
			r := rec
			return &r
		}
	}

	return nil
}

func dispatchUrgentResend(payload PromptPayload, targetNode string) (*PromptResult, error) {
	if _, err := SendPromptImmediate(payload); err != nil {
		return nil, err
	}

	return &PromptResult{
		PromptID:      fmt.Sprintf("resend-%d", time.Now().UnixNano()),
		Status:        "delivered",
		TargetNode:    targetNode,
		ProjectName:   payload.EffectiveTarget(),
		DispatchedAt:  time.Now().Format(time.RFC3339),
		OutputPreview: "Prompt re-dispatched immediately to " + targetNode,
	}, nil
}

// ResendPrompt dispatches an existing prompt from history or template.
func ResendPrompt(req PromptResendRequest) (*PromptResult, error) {
	history, _ := GetPromptHistory()
	targetRecord := findPromptRecordByID(history, req.PromptID)
	if targetRecord == nil {
		saved, _ := LoadSavedPrompts()
		targetRecord = findPromptRecordByID(saved, req.PromptID)
	}
	if targetRecord == nil {
		return nil, apperror.NewWithDetails("agy.resend", "E404", "prompt not found", "cmdagy", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}

	targetNode := targetRecord.TargetNode
	if req.OverrideNode != "" {
		targetNode = req.OverrideNode
	}
	if targetNode == "" {
		targetNode, _ = os.Hostname()
	}

	payload := PromptPayload{
		ProjectName:   targetRecord.ProjectName,
		ProjectTarget: targetRecord.ProjectTarget,
		TargetNode:    targetNode,
		Title:         "Resend: " + targetRecord.Title,
		PromptText:    targetRecord.PromptText,
		IsUrgent:      req.IsUrgent,
	}

	if req.IsUrgent {
		return dispatchUrgentResend(payload, targetNode)
	}

	queuedRec, err := EnqueuePromptPayload(payload)
	if err != nil {
		return nil, err
	}
	return &PromptResult{
		PromptID:      queuedRec.ID,
		Status:        "queued",
		TargetNode:    targetNode,
		ProjectName:   payload.EffectiveTarget(),
		DispatchedAt:  time.Now().Format(time.RFC3339),
		OutputPreview: fmt.Sprintf("Prompt enqueued at position %d", queuedRec.QueuePosition),
	}, nil
}

func normalizeWorkspacePath(path string) string {
	if path == "" {
		return "$WORKSPACE_ROOT"
	}
	clean := filepath.ToSlash(path)
	return clean
}

// BuildPromptTree aggregates running and queued prompts grouped by node and project.
func BuildPromptTree() (*PromptTreeResponse, error) {
	nodeAlias, _ := os.Hostname()
	running, _ := DiscoverRunningProjects()
	queuePath := getTempAgyFilePath("agy-prompt-queue.json")
	queuedRecords, _ := readRecordsFromFile(queuePath)

	projectMap := make(map[string]*ProjectPromptSummary)
	totalActive := 0

	for _, p := range running {
		projName := p.ProjectName
		if projName == "" {
			projName = filepath.Base(p.ProjectPath)
		}
		summary, exists := projectMap[projName]
		if !exists {
			summary = &ProjectPromptSummary{
				ProjectName:   projName,
				WorkspacePath: normalizeWorkspacePath(p.ProjectPath),
				ActivePrompts: []PromptRecord{},
				QueuedPrompts: []PromptRecord{},
			}
			projectMap[projName] = summary
		}
		if p.HasActivePrompt {
			totalActive++
			summary.ActivePrompts = append(summary.ActivePrompts, PromptRecord{
				ID:          fmt.Sprintf("p-%s", p.ConversationId),
				ProjectName: projName,
				Title:       p.PromptPreview,
				PromptText:  p.FullPrompt,
				Status:      p.Status,
				CreatedAt:   time.Now().Format(time.RFC3339),
			})
		}
	}

	totalQueued := len(queuedRecords)
	for _, q := range queuedRecords {
		projName := q.ProjectName
		if projName == "" {
			projName = "default"
		}
		summary, exists := projectMap[projName]
		if !exists {
			summary = &ProjectPromptSummary{
				ProjectName:   projName,
				WorkspacePath: "$WORKSPACE_ROOT",
				ActivePrompts: []PromptRecord{},
				QueuedPrompts: []PromptRecord{},
			}
			projectMap[projName] = summary
		}
		summary.QueuedPrompts = append(summary.QueuedPrompts, q)
	}

	var projects []ProjectPromptSummary
	for _, p := range projectMap {
		projects = append(projects, *p)
	}

	return &PromptTreeResponse{
		TotalActive: totalActive,
		TotalQueued: totalQueued,
		Nodes: []NodePromptTree{
			{
				NodeAlias: nodeAlias,
				IsOnline:  true,
				Projects:  projects,
			},
		},
	}, nil
}

func buildActiveSummaryFromProject(p RunningProjectRecord) ActivePromptSummary {
	snip := p.PromptPreview
	if snip == "" {
		snip = p.FullPrompt
	}
	if len(snip) > 80 {
		snip = snip[:77] + "..."
	}
	return ActivePromptSummary{
		WorkspacePath:  p.ProjectPath,
		ProjectName:    p.ProjectName,
		ConversationID: p.ConversationId,
		PromptTitle:    p.PromptPreview,
		PromptSnippet:  snip,
		Status:         p.Status,
	}
}

// CollectActivePrompts extracts summaries from running projects.
func CollectActivePrompts(running []RunningProjectRecord) []ActivePromptSummary {
	var summaries []ActivePromptSummary
	for _, p := range running {
		if p.HasActivePrompt {
			summaries = append(summaries, buildActiveSummaryFromProject(p))
		}
	}
	return summaries
}

// GetAGYUIStatus collects complete telemetry for the AGY web studio.
func GetAGYUIStatus() (*AGYUIStatusPayload, error) {
	node, _ := os.Hostname()
	running, _ := DiscoverRunningProjects()
	activePrompts := CollectActivePrompts(running)
	qStatus, _ := GetQueueStatus()
	saved, _ := LoadSavedPrompts()

	return &AGYUIStatusPayload{
		IsServerRunning: true,
		NodeAlias:       node,
		RunningProjects: running,
		ActivePrompts:   activePrompts,
		QueuedPrompts:   qStatus.Queued,
		SavedPrompts:    saved,
	}, nil
}

var (
	queueWorkerOnce sync.Once
	queueWorkerStop chan struct{}
)

// StartPromptQueueWorker begins continuous background processing of queued prompts.
func StartPromptQueueWorker() {
	queueWorkerOnce.Do(func() {
		queueWorkerStop = make(chan struct{})
		go runPromptQueueWorkerLoop(queueWorkerStop)
	})
}

// StopPromptQueueWorker signals the background queue worker to terminate.
func StopPromptQueueWorker() {
	if queueWorkerStop != nil {
		close(queueWorkerStop)
		queueWorkerStop = nil
	}
}

func runPromptQueueWorkerLoop(stopChan <-chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	concurrencyLimit := make(chan struct{}, 2)

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			processNextQueuedPrompts(concurrencyLimit)
		}
	}
}

func processNextQueuedPrompts(sem chan struct{}) {
	queuePath := getTempAgyFilePath("agy-prompt-queue.json")
	records, err := readRecordsFromFile(queuePath)
	if err != nil || len(records) == 0 {
		return
	}

	var remaining []PromptRecord
	for i, rec := range records {
		if rec.Status == "queued" && len(sem) < cap(sem) {
			sem <- struct{}{}
			go func(r PromptRecord) {
				defer func() { <-sem }()
				executeQueuedPrompt(r)
			}(rec)
		} else if rec.Status == "queued" {
			rec.QueuePosition = len(remaining) + 1
			remaining = append(remaining, rec)
		} else if i > 0 {
			remaining = append(remaining, rec)
		}
	}
	_ = writeRecordsToFile(queuePath, remaining)
}

func executeQueuedPrompt(rec PromptRecord) {
	payload := PromptPayload{
		ProjectName:   rec.ProjectName,
		ProjectTarget: rec.ProjectTarget,
		TargetNode:    rec.TargetNode,
		Title:         rec.Title,
		PromptText:    rec.PromptText,
	}
	_, _ = SendPromptImmediate(payload)
	rec.Status = "completed"
	rec.CompletedAt = time.Now().Format(time.RFC3339)
	_ = AppendPromptHistory(rec)
}
