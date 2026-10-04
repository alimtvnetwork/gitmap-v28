package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PromptPayload models an inbound prompt dispatch or enqueue request.
type PromptPayload struct {
	ProjectTarget string `json:"projectTarget"`
	Title         string `json:"title,omitempty"`
	PromptText    string `json:"promptText"`
	IsEnqueue     bool   `json:"isEnqueue,omitempty"`
}

// PromptRecord stores state and metadata for a prompt across its lifecycle.
type PromptRecord struct {
	ID            string `json:"id"`
	ProjectTarget string `json:"projectTarget"`
	Title         string `json:"title"`
	PromptText    string `json:"promptText"`
	Status        string `json:"status"` // "running", "queued", "saved", "completed"
	CreatedAt     string `json:"createdAt"`
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
	rec := PromptRecord{
		ID:            fmt.Sprintf("q-%d", time.Now().UnixNano()),
		ProjectTarget: payload.ProjectTarget,
		Title:         payload.Title,
		PromptText:    payload.PromptText,
		Status:        "queued",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	records = append(records, rec)
	if err := writeRecordsToFile(path, records); err != nil {
		return nil, err
	}
	_, _ = EnqueuePrompt("ui_enqueue", payload.Title, payload.PromptText)
	return &rec, nil
}

// SendPromptImmediate dispatches a prompt immediately to a target project.
func SendPromptImmediate(payload PromptPayload) (AgyInjectionResult, error) {
	target := payload.ProjectTarget
	if target == "" {
		target = "."
	}
	title := payload.Title
	if title == "" {
		title = "UI Direct Dispatch"
	}
	res := InjectAgyPrompt(target, payload.PromptText, title, "ui_send", false)
	rec := PromptRecord{
		ID:            fmt.Sprintf("disp-%d", time.Now().UnixNano()),
		ProjectTarget: target,
		Title:         title,
		PromptText:    payload.PromptText,
		Status:        "running",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	_ = AppendPromptHistory(rec)
	return res, nil
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
