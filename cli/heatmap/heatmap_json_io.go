package heatmap

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// JsonActionPayload is the JSON deserialization target supporting flexible field names.
type JsonActionPayload struct {
	Task             string `json:"task"`
	TaskSlug         string `json:"taskSlug"`
	Action           string `json:"action"`
	ActionType       string `json:"actionType"`
	File             string `json:"file"`
	FilePath         string `json:"filePath"`
	Lines            string `json:"lines"`
	StartLine        int    `json:"startLine"`
	EndLine          int    `json:"endLine"`
	Code             string `json:"code"`
	CodeSnippet      string `json:"codeSnippet"`
	Cmd              string `json:"cmd"`
	CommandText      string `json:"commandText"`
	State            string `json:"state"`
	ReasoningState   string `json:"reasoningState"`
	Reason           string `json:"reason"`
	ReasoningNotes   string `json:"reasoningNotes"`
	Expected         string `json:"expected"`
	ExpectedOutcome  string `json:"expectedOutcome"`
	Actual           string `json:"actual"`
	ActualOutcome    string `json:"actualOutcome"`
	IsExpectationMet *bool  `json:"isExpectationMet"`
	Met              *bool  `json:"met"`
	ParentID         *int64 `json:"parentId"`
	ParentLogId      *int64 `json:"parentLogId"`
	Repo             string `json:"repo"`
	RepoSlug         string `json:"repoSlug"`
	Ai               string `json:"ai"`
	AiName           string `json:"aiName"`
	Model            string `json:"model"`
	ModelName        string `json:"modelName"`
	Tool             string `json:"tool"`
	ToolName         string `json:"toolName"`
	Notes            string `json:"notes"`
	Comments         string `json:"comments"`
}

// ParseJsonActionPayload unmarshals raw JSON into a HeatmapActionLogRecord.
func ParseJsonActionPayload(raw string) (*store.HeatmapActionLogRecord, *apperror.AppError) {
	trimmed := strings.TrimSpace(raw)
	hasValidBounds := len(trimmed) >= 2 && strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")
	if !hasValidBounds {
		return nil, apperror.NewValidation("heatmap.json", "E1001", "invalid json payload: must be enclosed in braces")
	}

	var p JsonActionPayload
	if err := json.Unmarshal([]byte(trimmed), &p); err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.json.unmarshal")
	}

	rec := convertPayloadToRecord(p)

	return &rec, nil
}

func convertPayloadToRecord(p JsonActionPayload) store.HeatmapActionLogRecord {
	rec := buildPrimaryActionFields(p)
	populateActionDetails(&rec, p)

	return rec
}

func buildPrimaryActionFields(p JsonActionPayload) store.HeatmapActionLogRecord {
	start, end := resolveLines(p.Lines, p.StartLine, p.EndLine)

	return store.HeatmapActionLogRecord{
		ParentLogId: resolveParentID(p.ParentID, p.ParentLogId),
		TaskSlug:    resolveFirstNonEmpty(p.Task, p.TaskSlug),
		RepoName:    resolveFirstNonEmpty(p.Repo, p.RepoSlug),
		ActionType:  strings.ToUpper(resolveDefault(resolveFirstNonEmpty(p.Action, p.ActionType), "ANALYZE")),
		FilePath:    resolveFirstNonEmpty(p.File, p.FilePath),
		StartLine:   start,
		EndLine:     end,
	}
}

func populateActionDetails(r *store.HeatmapActionLogRecord, p JsonActionPayload) {
	r.ReasoningState = strings.ToUpper(resolveDefault(resolveFirstNonEmpty(p.State, p.ReasoningState), "EXPLORATION"))
	r.ReasoningNotes = resolveFirstNonEmpty(p.Reason, p.ReasoningNotes)
	r.ExpectedOutcome = resolveFirstNonEmpty(p.Expected, p.ExpectedOutcome)
	r.ActualOutcome = resolveFirstNonEmpty(p.Actual, p.ActualOutcome)
	r.CodeSnippet = resolveFirstNonEmpty(p.Code, p.CodeSnippet)
	r.CommandText = resolveFirstNonEmpty(p.Cmd, p.CommandText)
	r.Notes = resolveFirstNonEmpty(p.Notes, p.Comments)
	r.IsExpectationMet = resolveIsMet(p.IsExpectationMet, p.Met)
}

func resolveParentID(a, b *int64) *int64 {
	if a != nil {
		return a
	}

	return b
}

func resolveIsMet(a, b *bool) bool {
	if a != nil {
		return *a
	}
	if b != nil {
		return *b
	}

	return true
}

func resolveFirstNonEmpty(a, b string) string {
	if strings.TrimSpace(a) != "" {
		return strings.TrimSpace(a)
	}

	return strings.TrimSpace(b)
}

func resolveDefault(val, def string) string {
	if strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}

	return def
}

func resolveLines(lineStr string, start, end int) (int, int) {
	if lineStr != "" {
		return parseLineRange(lineStr)
	}

	return start, end
}

// RenderRcJsonActionResponse emits structured JSON envelope for an action.
func RenderRcJsonActionResponse(out io.Writer, rec store.HeatmapActionLogRecord, dbPath, modelName, toolName string) error {
	resp := map[string]any{
		"status":     "success",
		"actionId":   rec.HeatmapActionLogId,
		"taskSlug":   rec.TaskSlug,
		"repoSlug":   rec.RepoName,
		"actionType": rec.ActionType,
		"filePath":   rec.FilePath,
		"modelName":  modelName,
		"toolName":   toolName,
		"database":   dbPath,
		"isMet":      rec.IsExpectationMet,
		"createdAt":  rec.CreatedAt,
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")

	return enc.Encode(resp)
}

// ReadJsonInputFromSource resolves json string from flag, positional argument, or stdin.
func ReadJsonInputFromSource(flags map[string]string, pos []string) (string, bool) {
	rawInput := flags["json-input"]
	if rawInput == "" {
		rawInput = flags["input"]
	}
	if rawInput == "-" {
		return readStdinJson()
	}
	if rawInput != "" {
		return rawInput, true
	}
	if len(pos) > 0 && strings.HasPrefix(strings.TrimSpace(pos[0]), "{") {
		return pos[0], true
	}

	return "", false
}

func readStdinJson() (string, bool) {
	stdinBytes, err := io.ReadAll(os.Stdin)
	if err != nil || len(stdinBytes) == 0 {
		return "", false
	}

	return string(stdinBytes), true
}
