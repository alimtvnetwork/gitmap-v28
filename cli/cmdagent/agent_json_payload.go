package cmdagent

import (
	"encoding/json"
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AgentLogJSONPayload models direct JSON telemetry input for `gitmap agent log`.
type AgentLogJSONPayload struct {
	SubtaskId    int64  `json:"subtaskId,omitempty"`
	Subtask      any    `json:"subtask,omitempty"`
	AgentRole    string `json:"agentRole,omitempty"`
	Role         string `json:"role,omitempty"`
	Agent        string `json:"agent,omitempty"`
	Action       string `json:"action,omitempty"`
	ActionType   string `json:"actionType,omitempty"`
	TargetFile   string `json:"targetFile,omitempty"`
	File         string `json:"file,omitempty"`
	StartLine    int    `json:"startLine,omitempty"`
	EndLine      int    `json:"endLine,omitempty"`
	Query        string `json:"query,omitempty"`
	Command      string `json:"command,omitempty"`
	Details      string `json:"details,omitempty"`
	Reason       string `json:"reason,omitempty"`
	DurationMs   int64  `json:"durationMs,omitempty"`
	Status       string `json:"status,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Error        string `json:"error,omitempty"`
	DirFlag      string `json:"dir,omitempty"`
	TaskIdFlag   string `json:"taskId,omitempty"`
	DbFlag       string `json:"db,omitempty"`
}

// IsJSONPayloadToken reports whether a string represents a JSON object or array.
func IsJSONPayloadToken(raw string) bool {
	trimmed := strings.TrimSpace(raw)

	return strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
}

// ExtractJSONPayload finds the first JSON object or array in a slice of command arguments.
func ExtractJSONPayload(args []string) (string, bool) {
	for i := 0; i < len(args); i++ {
		payload, hasPayload := checkArgumentJSON(args, i)
		if hasPayload {
			return payload, true
		}
	}

	return "", false
}

func checkArgumentJSON(args []string, i int) (string, bool) {
	arg := args[i]
	if IsJSONPayloadToken(arg) {
		return arg, true
	}

	if payload, hasFlagVal := extractFlagValueJSON(arg); hasFlagVal {
		return payload, true
	}

	return checkNextArgJSON(args, i, arg)
}

func extractFlagValueJSON(arg string) (string, bool) {
	prefixes := []string{"--json=", "--data=", "--input="}
	for _, p := range prefixes {
		if strings.HasPrefix(arg, p) {
			val := strings.TrimPrefix(arg, p)
			return val, IsJSONPayloadToken(val)
		}
	}

	return "", false
}

func checkNextArgJSON(args []string, i int, arg string) (string, bool) {
	isMatch := arg == "--json" || arg == "-j" || arg == "--data" || arg == "-d" || arg == "--input"
	hasNext := i+1 < len(args) && IsJSONPayloadToken(args[i+1])
	if isMatch && hasNext {
		return args[i+1], true
	}

	return "", false
}

func resolveSubtaskNumber(val any) int64 {
	switch v := val.(type) {
	case float64:
		return int64(v)
	case string:
		return parseSubtaskId(v)
	default:
		return 0
	}
}

func (p AgentLogJSONPayload) ResolveSubtaskId() int64 {
	if p.SubtaskId > 0 {
		return p.SubtaskId
	}

	return resolveSubtaskNumber(p.Subtask)
}

func (p AgentLogJSONPayload) ResolveRole() string {
	if p.AgentRole != "" {
		return p.AgentRole
	}
	if p.Role != "" {
		return p.Role
	}

	return p.Agent
}

func (p AgentLogJSONPayload) ResolveAction() string {
	if p.Action != "" {
		return p.Action
	}

	return p.ActionType
}

func (p AgentLogJSONPayload) ResolveFile() string {
	if p.TargetFile != "" {
		return p.TargetFile
	}

	return p.File
}

func (p AgentLogJSONPayload) ResolveQuery() string {
	if p.Query != "" {
		return p.Query
	}

	return p.Command
}

func (p AgentLogJSONPayload) ResolveDetails() string {
	if p.Details != "" {
		return p.Details
	}

	return p.Reason
}

func (p AgentLogJSONPayload) ResolveError() string {
	if p.ErrorMessage != "" {
		return p.ErrorMessage
	}

	return p.Error
}

func applyPayloadFields(opts *LogOptions, p AgentLogJSONPayload) {
	if subID := p.ResolveSubtaskId(); subID > 0 {
		opts.SubtaskId = subID
	}
	if role := p.ResolveRole(); role != "" {
		opts.AgentRole = role
	}
	if act := p.ResolveAction(); act != "" {
		opts.Action = NormalizeActionType(act)
	}
	if file := p.ResolveFile(); file != "" {
		opts.TargetFile = file
	}
	applyPayloadDetails(opts, p)
}

func applyPayloadDetails(opts *LogOptions, p AgentLogJSONPayload) {
	if p.StartLine > 0 {
		opts.StartLine = p.StartLine
	}
	if p.EndLine > 0 {
		opts.EndLine = p.EndLine
	}
	if q := p.ResolveQuery(); q != "" {
		opts.Query = q
	}
	if d := p.ResolveDetails(); d != "" {
		opts.Details = d
	}
	applyPayloadOutcome(opts, p)
}

func applyPayloadOutcome(opts *LogOptions, p AgentLogJSONPayload) {
	if p.DurationMs > 0 {
		opts.DurationMs = p.DurationMs
	}
	if p.Status != "" {
		opts.Status = strings.ToUpper(p.Status)
	}
	if errStr := p.ResolveError(); errStr != "" {
		opts.ErrorMessage = errStr
	}
	if p.DirFlag != "" {
		opts.DirFlag = p.DirFlag
	}
	if p.TaskIdFlag != "" {
		opts.TaskIdFlag = p.TaskIdFlag
	}
}

// ApplyAgentLogJSON parses a JSON payload and populates LogOptions fields.
func ApplyAgentLogJSON(opts *LogOptions, raw string) (bool, error) {
	if !IsJSONPayloadToken(raw) {
		return false, nil
	}

	var payload AgentLogJSONPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return false, appfault.WrapSimple(err, fmt.Sprintf("invalid agent log JSON: %s", raw))
	}

	applyPayloadFields(opts, payload)

	return true, nil
}

// InterceptAgentLogJSON extracts and applies a JSON payload from CLI arguments if present.
func InterceptAgentLogJSON(args []string, opts *LogOptions) bool {
	raw, hasPayload := ExtractJSONPayload(args)
	if !hasPayload {
		return false
	}

	hasApplied, _ := ApplyAgentLogJSON(opts, raw)

	return hasApplied
}

// InterceptSubtaskAddJSON extracts a JSON payload from arguments for `gitmap agent subtask add`.
func InterceptSubtaskAddJSON(args []string, opts *SubtaskAddOptions) bool {
	if len(opts.JsonInput) > 0 {
		return true
	}

	raw, hasPayload := ExtractJSONPayload(args)
	if hasPayload {
		opts.JsonInput = raw
		return true
	}

	return false
}
