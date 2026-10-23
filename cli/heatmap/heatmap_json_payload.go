package heatmap

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// HeatmapActionPayload models a direct JSON action telemetry payload.
type HeatmapActionPayload struct {
	TaskSlug         string `json:"taskSlug,omitempty"`
	RepoName         string `json:"repoName,omitempty"`
	Action           string `json:"action,omitempty"`
	File             string `json:"file,omitempty"`
	State            string `json:"state,omitempty"`
	Expected         string `json:"expected,omitempty"`
	Actual           string `json:"actual,omitempty"`
	IsExpectationMet bool   `json:"isExpectationMet"`
	AiName           string `json:"aiName,omitempty"`
	ModelName        string `json:"modelName,omitempty"`
	ToolName         string `json:"toolName,omitempty"`
	ToolParameters   string `json:"toolParameters,omitempty"`
	DurationMs       int64  `json:"durationMs,omitempty"`
}

// ParseActionPayloads parses either a single JSON object or an array of JSON objects.
func ParseActionPayloads(raw string) ([]HeatmapActionPayload, error) {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, apperror.NewValidationError("empty JSON payload")
	}

	if strings.HasPrefix(trimmed, "[") {
		return parseActionPayloadArray(trimmed)
	}

	if strings.HasPrefix(trimmed, "{") {
		return parseSingleActionPayload(trimmed)
	}

	return nil, apperror.NewValidationError("invalid JSON payload: expected '{' or '['")
}

func parseActionPayloadArray(trimmed string) ([]HeatmapActionPayload, error) {
	var payloads []HeatmapActionPayload
	if err := json.Unmarshal([]byte(trimmed), &payloads); err != nil {
		return nil, apperror.WrapSimple(err, "parse action payloads array")
	}

	return payloads, nil
}

func parseSingleActionPayload(trimmed string) ([]HeatmapActionPayload, error) {
	var single HeatmapActionPayload
	if err := json.Unmarshal([]byte(trimmed), &single); err != nil {
		return nil, apperror.WrapSimple(err, "parse single action payload")
	}

	return []HeatmapActionPayload{single}, nil
}

func resolveActionModel(p HeatmapActionPayload) string {
	if p.ModelName != "" {
		return p.ModelName
	}

	return p.AiName
}

func resolveActionType(p HeatmapActionPayload) string {
	if p.Action != "" {
		return strings.ToUpper(p.Action)
	}

	return "ANALYZE"
}

func resolveActionState(p HeatmapActionPayload) string {
	if p.State != "" {
		return strings.ToUpper(p.State)
	}

	return "EXPLORATION"
}

// ToRecord converts HeatmapActionPayload to a store.HeatmapActionLogRecord.
func (p HeatmapActionPayload) ToRecord() store.HeatmapActionLogRecord {
	return store.HeatmapActionLogRecord{
		TaskSlug:         p.TaskSlug,
		RepoName:         p.RepoName,
		ActionType:       resolveActionType(p),
		FilePath:         p.File,
		ReasoningState:   resolveActionState(p),
		ExpectedOutcome:  p.Expected,
		ActualOutcome:    p.Actual,
		IsExpectationMet: p.IsExpectationMet,
		ModelName:        resolveActionModel(p),
		ToolName:         p.ToolName,
		CommandText:      p.ToolParameters,
	}
}

// ToFlags converts HeatmapActionPayload into a flags map for heatmap log execution.
func (p HeatmapActionPayload) ToFlags() map[string]string {
	flags := initActionFlags(p)
	flags["met"] = resolveExpectationFlag(p.IsExpectationMet)

	return flags
}

func resolveExpectationFlag(isMet bool) string {
	if isMet {
		return "true"
	}

	return "false"
}

func initActionFlags(p HeatmapActionPayload) map[string]string {
	return map[string]string{
		"task": p.TaskSlug, "repo": p.RepoName, "action": p.Action,
		"file": p.File, "state": p.State, "expected": p.Expected,
		"actual": p.Actual, "ai": p.AiName, "model": p.ModelName,
		"tool": p.ToolName, "cmd": p.ToolParameters,
	}
}
